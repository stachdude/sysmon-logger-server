package main

// ##### Structs #############################################################

// Stores the YAML config file data
type Config struct {
	DatabaseServer   string `yaml:"database_server"`
	DatabaseName     string `yaml:"database_name"`
	DatabaseUser     string `yaml:"database_user"`
	DatabasePassword string `yaml:"database_password"`
	HttpIp           string `yaml:"http_ip"`
	HttpPort         int16  `yaml:"http_port"`
	ProcessorThreads int    `yaml:"processor_threads"`
	Debug            bool   `yaml:"debug"`
	ServerPem        string `yaml:"server_pem"`
	ServerKey        string `yaml:"server_key"`
	TempDir          string `yaml:"temp_dir"`
	ExportDir        string `yaml:"export_dir"`
	MaxDataAgeDays   int    `yaml:"max_data_age_days"`
}
