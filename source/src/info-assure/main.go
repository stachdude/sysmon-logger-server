package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
	"github.com/robfig/cron"
	"github.com/voxelbrain/goptions"
	util "github.com/woanware/goutil"
	"gopkg.in/mgutz/dat.v1"
	"gopkg.in/mgutz/dat.v1/sqlx-runner"
	"gopkg.in/yaml.v2"
	"os"
	"runtime"
	"time"
)

// ##### Constants  ###########################################################

const APP_TITLE string = "SysMon Logger"
const APP_NAME string = "sml-server"
const APP_VERSION string = "1.0.1"

// ##### Variables ###########################################################

var (
	logger    *logging.Logger
	config    *Config
	workQueue chan ImportTask
	db        *runner.DB
	cronner   *cron.Cron
)

// ##### Methods #############################################################

// Application entry point
func main() {

	fmt.Printf("\n%s (%s) %s\n\n", APP_TITLE, APP_NAME, APP_VERSION)

	initialiseLogging()

	opt := struct {
		ConfigFile string        `goptions:"-c, --config, description='Config file path'"`
		Help       goptions.Help `goptions:"-h, --help, description='Show this help'"`
	}{ // Default values
		ConfigFile: "./" + APP_NAME + ".config",
	}

	goptions.ParseAndFail(&opt)

	// Load the applications configuration such as database credentials
	config = loadConfig(opt.ConfigFile)

	initialiseDatabase()
	createProcessors()

	performHourlyTasks()
	return

	cronner = cron.New()
	//cronner.AddFunc("1 * * * * *", performHourlyTasks)
	cronner.AddFunc("@daily", performDataPurge)
	cronner.AddFunc("@hourly", performHourlyTasks)
	cronner.Start()

	var r *gin.Engine
	if config.Debug == true {
		// DEBUG
		r = gin.Default()
	} else {
		gin.SetMode(gin.ReleaseMode)
		r = gin.New()
		r.Use(gin.Recovery())
	}

	r.GET("/", index)
	r.GET("/:domain/:host", receive)
	r.POST("/:domain/:host", receiveData)
	r.RunTLS(config.HttpIp+":"+fmt.Sprintf("%d", config.HttpPort), config.ServerPem, config.ServerKey)
}

//
func initialiseDatabase() {

	// create a normal database connection through database/sql
	tempDb, err := sql.Open("postgres",
		fmt.Sprintf("host=%s dbname=%s user=%s password=%s sslmode=disable",
			config.DatabaseServer, config.DatabaseName, config.DatabaseUser, config.DatabasePassword))

	if err != nil {
		logger.Fatalf("Unable to open database connection: %v", err)
		return
	}

	// ensures the database can be pinged with an exponential backoff (15 min)
	runner.MustPing(tempDb)

	// set to reasonable values for production
	tempDb.SetMaxIdleConns(4)
	tempDb.SetMaxOpenConns(16)

	// set this to enable interpolation
	dat.EnableInterpolation = true

	// set to check things like sessions closing.
	// Should be disabled in production/release builds.
	dat.Strict = false

	// Log any query over 10ms as warnings. (optional)
	runner.LogQueriesThreshold = 50 * time.Millisecond

	db = runner.NewDB(tempDb, "postgres")
}

// Initialise the channels for the cross process comms and then start the workers
func createProcessors() {

	processorCount := runtime.NumCPU()
	if config.ProcessorThreads > 0 {
		processorCount = config.ProcessorThreads
	}

	workQueue = make(chan ImportTask, 100)

	// Create the workers that perform the actual processing
	for i := 0; i < processorCount; i++ {
		logger.Infof("Initialising processor: %d", i+1)
		p := NewProcessor(i, config, db)
		go func(p *Processor) {
			for j := range workQueue {
				p.Process(j)
			}
		}(p)
	}
}

// Sets up the logging infrastructure e.g. Stdout and /var/log
func initialiseLogging() {

	// Setup the actual loggers
	logger = logging.MustGetLogger(APP_NAME)

	// Check that we have a "nca" sub directory in /var/log
	if _, err := os.Stat("/var/log/" + APP_NAME); os.IsNotExist(err) {
		logger.Fatal("The /var/log/" + APP_NAME + " directory does not exist")
	}

	// Check that we have permission to write to the /var/log/APP_NAME directory
	f, err := os.Create("/var/log/" + APP_NAME + "/test.txt")
	if err != nil {
		logger.Fatal("Unable to write to /var/log/" + APP_NAME)
	}

	// Clear up our tests
	os.Remove("/var/log/" + APP_NAME + "/test.txt")
	f.Close()

	// Define the /var/log file
	logFile, err := os.OpenFile("/var/log/"+APP_NAME+"/log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logger.Fatalf("Error opening the log file: %v", err)
	}

	// Define the StdOut loggingDatabaser
	backendStdOut := logging.NewLogBackend(os.Stdout, "", 0)
	formatStdOut := logging.MustStringFormatter(
		"%{color}%{time:15:04:05.000000} %{color:reset} %{message}")
	formatterStdOut := logging.NewBackendFormatter(backendStdOut, formatStdOut)

	// Define the /var/log logging
	backendFile := logging.NewLogBackend(logFile, "", 0)
	formatFile := logging.MustStringFormatter(
		"%{time:15:04:05.000000} %{level:.4s} %{message}")
	formatterFile := logging.NewBackendFormatter(backendFile, formatFile)

	logging.SetBackend(formatterStdOut, formatterFile)
}

// Loads the applications config file contents (yaml) and marshals to a struct
func loadConfig(configPath string) *Config {

	c := new(Config)
	data, err := util.ReadTextFromFile(configPath)
	if err != nil {
		logger.Fatalf("Error reading the config file: %v", err)
	}

	err = yaml.Unmarshal([]byte(data), &c)
	if err != nil {
		logger.Fatalf("Error unmarshalling the config file: %v", err)
	}

	if len(c.DatabaseServer) == 0 {
		logger.Fatal("Database server not set in config file")
	}

	if len(c.DatabaseName) == 0 {
		logger.Fatal("Database name not set in config file")
	}

	if len(c.DatabaseUser) == 0 {
		logger.Fatal("Database user not set in config file")
	}

	if len(c.DatabasePassword) == 0 {
		logger.Fatal("Database password not set in config file")
	}

	if len(c.HttpIp) == 0 {
		logger.Fatal("HTTP IP not set in config file")
	}

	if len(c.ServerPem) == 0 {
		logger.Fatal("Server PEM file not set in config file")
	}

	if len(c.ServerKey) == 0 {
		logger.Fatal("Server key file not set in config file")
	}

	if len(c.TempDir) == 0 {
		logger.Fatal("Temp directory not set in config file")
	}

	if len(c.ExportDir) == 0 {
		logger.Fatal("Export directory not set in config file")
	}

	return c
}

//
func performHourlyTasks() {

	exportDataForStringStringTotal(
		SQL_PROCESS_CREATE_PATH_SHA256_GROUP_PATH_ORDER_PATH, "Process Create (Path, SHA256) By Path",
		EXPORT_TYPE_PROCESS_CREATE_PATH_SHA256_GROUP_PATH_ORDER_PATH, PREFIX_EXPORT_PROCESS_CREATE_PATH_SHA256_GROUP_PATH_ORDER_PATH)

	exportDataForStringStringTotal(
		SQL_PROCESS_CREATE_PATH_SHA256_GROUP_PATH_ORDER_HASH, "Process Create (Path, SHA256) By Hash",
		EXPORT_TYPE_PROCESS_CREATE_PATH_SHA256_GROUP_PATH_ORDER_HASH, PREFIX_EXPORT_PROCESS_CREATE_PATH_SHA256_GROUP_PATH_ORDER_HASH)

	exportDataForStringStringTotal(
		SQL_PROCESS_CREATE_PATH_MD5_GROUP_PATH_ORDER_PATH, "Process Create (Path, MD5) By Path",
		EXPORT_TYPE_PROCESS_CREATE_PATH_MD5_GROUP_PATH_ORDER_PATH, PREFIX_EXPORT_PROCESS_CREATE_PATH_MD5_GROUP_PATH_ORDER_PATH)

	exportDataForStringStringTotal(
		SQL_PROCESS_CREATE_PATH_MD5_GROUP_PATH_ORDER_HASH, "Process Create (Path, MD5) By Hash",
		EXPORT_TYPE_PROCESS_CREATE_PATH_MD5_GROUP_PATH_ORDER_HASH, PREFIX_EXPORT_PROCESS_CREATE_PATH_MD5_GROUP_PATH_ORDER_HASH)

	exportDataForStringStringTotal(
		SQL_DRIVER_LOADED_PATH_SHA256_GROUP_PATH_ORDER_PATH, "Driver Loaded (Path, SHA256) By Path",
		EXPORT_TYPE_DRIVER_LOADED_PATH_SHA256_GROUP_PATH_ORDER_PATH, PREFIX_EXPORT_DRIVER_LOADED_PATH_SHA256_GROUP_PATH_ORDER_PATH)

	exportDataForStringStringTotal(
		SQL_DRIVER_LOADED_PATH_SHA256_GROUP_PATH_ORDER_HASH, "Driver Loaded (Path, SHA256) By Hash",
		EXPORT_TYPE_DRIVER_LOADED_PATH_SHA256_GROUP_PATH_ORDER_HASH, PREFIX_EXPORT_DRIVER_LOADED_PATH_SHA256_GROUP_PATH_ORDER_HASH)

	exportDataForStringStringTotal(
		SQL_DRIVER_LOADED_PATH_MD5_GROUP_PATH_ORDER_PATH, "Driver Loaded (Path, MD5) By Path",
		EXPORT_TYPE_DRIVER_LOADED_PATH_MD5_GROUP_PATH_ORDER_PATH, PREFIX_EXPORT_DRIVER_LOADED_PATH_MD5_GROUP_PATH_ORDER_PATH)

	exportDataForStringStringTotal(
		SQL_DRIVER_LOADED_PATH_MD5_GROUP_PATH_ORDER_HASH, "Driver Loaded (Path, MD5) By Hash",
		EXPORT_TYPE_DRIVER_LOADED_PATH_MD5_GROUP_PATH_ORDER_HASH, PREFIX_EXPORT_DRIVER_LOADED_PATH_MD5_GROUP_PATH_ORDER_HASH)

	exportDataForStringStringTotal(
		SQL_IMAGE_LOADED_PATH_SHA256_GROUP_PATH_ORDER_PATH, "Image Loaded (Path, SHA256) By Path",
		EXPORT_TYPE_IMAGE_LOADED_PATH_SHA256_GROUP_PATH_ORDER_PATH, PREFIX_EXPORT_IMAGE_LOADED_PATH_SHA256_GROUP_PATH_ORDER_PATH)

	exportDataForStringStringTotal(
		SQL_IMAGE_LOADED_PATH_SHA256_GROUP_PATH_ORDER_HASH, "Image Loaded (Path, SHA256) By Hash",
		EXPORT_TYPE_IMAGE_LOADED_PATH_SHA256_GROUP_PATH_ORDER_HASH, PREFIX_EXPORT_IMAGE_LOADED_PATH_SHA256_GROUP_PATH_ORDER_HASH)

	exportDataForStringStringTotal(
		SQL_IMAGE_LOADED_PATH_MD5_GROUP_PATH_ORDER_PATH, "Image Loaded (Path, MD5) By Path",
		EXPORT_TYPE_IMAGE_LOADED_PATH_MD5_GROUP_PATH_ORDER_PATH, PREFIX_EXPORT_IMAGE_LOADED_PATH_MD5_GROUP_PATH_ORDER_PATH)

	exportDataForStringStringTotal(
		SQL_IMAGE_LOADED_PATH_MD5_GROUP_PATH_ORDER_HASH, "Image Loaded (Path, MD5) By Hash",
		EXPORT_TYPE_IMAGE_LOADED_PATH_MD5_GROUP_PATH_ORDER_HASH, PREFIX_EXPORT_IMAGE_LOADED_PATH_MD5_GROUP_PATH_ORDER_HASH)

	exportDataForStringStringTotal(
		SQL_FILE_STREAM_PATH_SHA256_GROUP_PATH_ORDER_PATH, "File Stream (Path, SHA256) By Path",
		EXPORT_TYPE_FILE_STREAM_PATH_SHA256_GROUP_PATH_ORDER_PATH, PREFIX_EXPORT_FILE_STREAM_PATH_SHA256_GROUP_PATH_ORDER_PATH)

	exportDataForStringStringTotal(
		SQL_FILE_STREAM_PATH_SHA256_GROUP_PATH_ORDER_HASH, "File Stream (Path, SHA256) By Hash",
		EXPORT_TYPE_FILE_STREAM_PATH_SHA256_GROUP_PATH_ORDER_HASH, PREFIX_EXPORT_FILE_STREAM_PATH_SHA256_GROUP_PATH_ORDER_HASH)

	exportDataForStringStringTotal(
		SQL_FILE_STREAM_PATH_MD5_GROUP_PATH_ORDER_PATH, "File Stream (Path, MD5) By Path",
		EXPORT_TYPE_FILE_STREAM_PATH_MD5_GROUP_PATH_ORDER_PATH, PREFIX_EXPORT_FILE_STREAM_PATH_MD5_GROUP_PATH_ORDER_PATH)

	exportDataForStringStringTotal(
		SQL_FILE_STREAM_PATH_MD5_GROUP_PATH_ORDER_HASH, "File Stream (Path, MD5) By Hash",
		EXPORT_TYPE_FILE_STREAM_PATH_MD5_GROUP_PATH_ORDER_HASH, PREFIX_EXPORT_FILE_STREAM_PATH_MD5_GROUP_PATH_ORDER_HASH)

	exportDataForString(
		SQL_NETWORK_CONNECTION_DISTINCT_DEST_IP, "Network Connection (Destination IP)",
		EXPORT_TYPE_NETWORK_CONNECTION_DISTINCT_DEST_IP, PREFIX_EXPORT_NETWORK_CONNECTION_DISTINCT_DEST_IP)

	exportDataForStringTotal(
		SQL_NETWORK_CONNECTION_COUNT_DEST_IP, "Network Connection (Destination IP Count)",
		EXPORT_TYPE_NETWORK_CONNECTION_DISTINCT_DEST_IP_COUNT, PREFIX_EXPORT_NETWORK_CONNECTION_DISTINCT_DEST_IP_COUNT)

	exportDataForString(
		SQL_NETWORK_CONNECTION_DISTINCT_DEST_HOST, "Network Connection (Destination Host)",
		EXPORT_TYPE_NETWORK_CONNECTION_DISTINCT_DEST_HOST, PREFIX_EXPORT_NETWORK_CONNECTION_DISTINCT_DEST_HOST)

	exportDataForStringTotal(
		SQL_NETWORK_CONNECTION_COUNT_DEST_HOST, "Network Connection (Destination Host Count)",
		EXPORT_TYPE_NETWORK_CONNECTION_DISTINCT_DEST_HOST_COUNT, PREFIX_EXPORT_NETWORK_CONNECTION_DISTINCT_DEST_HOST_COUNT)

	exportDataForString(
		SQL_SHA256_ALL, "SHA256 (All)",
		EXPORT_TYPE_SHA256_ALL, PREFIX_EXPORT_SHA256_ALL)

	exportDataForString(
		SQL_SHA256_PROCESS_CREATE, "SHA256 (Process Create)",
		EXPORT_TYPE_SHA256_PROCESS_CREATE, PREFIX_EXPORT_SHA256_PROCESS_CREATE)

	exportDataForString(
		SQL_SHA256_DRIVER_LOADED, "SHA256 (Driver Loaded)",
		EXPORT_TYPE_SHA256_DRIVER_LOADED, PREFIX_EXPORT_SHA256_DRIVER_LOADED)

	exportDataForString(
		SQL_SHA256_IMAGE_LOADED, "SHA256 (Image Loaded)",
		EXPORT_TYPE_SHA256_IMAGE_LOADED, PREFIX_EXPORT_SHA256_IMAGE_LOADED)

	exportDataForString(
		SQL_SHA256_FILE_STREAM, "SHA256 (File Stream)",
		EXPORT_TYPE_SHA256_FILE_STREAM, PREFIX_EXPORT_SHA256_FILE_STREAM)

	exportDataForString(
		SQL_MD5_ALL, "MD5 (All)",
		EXPORT_TYPE_MD5_ALL, PREFIX_EXPORT_MD5_ALL)

	exportDataForString(
		SQL_MD5_PROCESS_CREATE, "MD5 (Process Create)",
		EXPORT_TYPE_MD5_PROCESS_CREATE, PREFIX_EXPORT_MD5_PROCESS_CREATE)

	exportDataForString(
		SQL_MD5_DRIVER_LOADED, "MD5 (Driver Loaded)",
		EXPORT_TYPE_MD5_DRIVER_LOADED, PREFIX_EXPORT_MD5_DRIVER_LOADED)

	exportDataForString(
		SQL_MD5_IMAGE_LOADED, "MD5 (Image Loaded)",
		EXPORT_TYPE_MD5_IMAGE_LOADED, PREFIX_EXPORT_MD5_IMAGE_LOADED)

	exportDataForString(
		SQL_MD5_FILE_STREAM, "MD5 (File Stream)",
		EXPORT_TYPE_MD5_FILE_STREAM, PREFIX_EXPORT_MD5_FILE_STREAM)
}

//
func performDataPurge() {

	if config.MaxDataAgeDays == 0 || config.MaxDataAgeDays == -1 {
		return
	}

	// Use the config file value to determine what is classed as an old job
	staleTimestamp := time.Now().UTC().Add(-time.Duration(24*config.MaxDataAgeDays) * time.Hour)

	for _, table := range databaseTables {
		logger.Info("Purging stale data from table: %s", table)

		_, err := db.
		DeleteFrom(table).
			Where("utc_time < $1", staleTimestamp).
			Exec()

		if err != nil {
			logger.Errorf("Error deleting stale data (%s): %v", err, table)
		}
	}
}
