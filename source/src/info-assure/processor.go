package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"github.com/woanware/goutil"
	"gopkg.in/mgutz/dat.v1/sqlx-runner"
	"regexp"
	"strings"
	"time"
)

// ##### Types ###############################################################

// Encapsulates a Processor object and its properties
type Processor struct {
	id             int
	config         *Config
	db             *runner.DB
	regexEventName *regexp.Regexp
	regexData      *regexp.Regexp
	regexUtcTime   *regexp.Regexp
}

// Encapsulates the data for an import task
type ImportTask struct {
	Domain string
	Host   string
	Data   string
}

// ##### Methods #############################################################

// Constructor/Initialiser for the Processor struct
func NewProcessor(id int, config *Config, db *runner.DB) *Processor {

	p := Processor{
		id:     id,
		config: config,
		db:     db,
	}

	p.regexEventName, _ = regexp.Compile(`\(rule:\s(.*?)\)`)
	p.regexData, _ = regexp.Compile(`<Data\sName='(.*?)'>(.*?)</Data>`)
	p.regexUtcTime, _ = regexp.Compile(`<TimeCreated\sSystemTime='(.*?)'/>`)

	return &p
}

// Process an individual set of host data
func (p Processor) Process(it ImportTask) {

	// Split the data into each separate message using the special delimiter :-)
	parts := strings.Split(it.Data, MESSAGE_DELIMITER)

	eventName := ""
	message := ""
	messageHtml := ""
	var e Event
	var err error
	var parsedTimestamp time.Time

	for _, v := range parts {

		// Extract the Event Name
		regexRes := p.regexEventName.FindAllStringSubmatch(v, -1)
		if regexRes == nil {
			logger.Errorf(`Cannot locate event name: %s`, v)
			continue
		}

		eventName = strings.TrimSpace(regexRes[0][1])
		// Lowercase for better matching
		eventName = strings.ToLower(eventName)

		// Set the various Event structure's properties
		e = Event{}
		e.Domain = it.Domain
		e.Host = it.Host

		// Extract the Event Log timestamp
		regexRes = p.regexUtcTime.FindAllStringSubmatch(v, -1)
		if regexRes != nil {
			parsedTimestamp, err = time.Parse(time.RFC3339Nano, strings.TrimSpace(regexRes[0][1]))
			if err != nil {
				logger.Error("Unable to parse Event Log Time: %v (%s)", err, regexRes[0][1])
				continue
			} else {
				e.EventLogTime = parsedTimestamp
			}
		}

		if config.Debug == true {
			logger.Infof("Domain: %s, Host: %s, Event: %s", it.Domain, it.Host, eventName)
		}

		switch eventName {
		case "processcreate":
			messageHtml, message = p.parseProcessCreate(it, parsedTimestamp, v)
			e.Type = "Process Create"

		case "filecreatetime":
			messageHtml, message = p.parseFileCreationTime(it, parsedTimestamp, v)
			e.Type = "File Creation Time"

		case "networkconnect":
			messageHtml, message = p.parseNetworkConnection(it, parsedTimestamp, v)
			e.Type = "Network Connection"

		case "processterminate":
			messageHtml, message = p.parseProcessTerminate(it, parsedTimestamp, v)
			e.Type = "Process Terminated"

		case "driverload":
			messageHtml, message = p.parseDriverLoaded(it, parsedTimestamp, v)
			e.Type = "Driver Loaded"

		case "imageload":
			messageHtml, message = p.parseImageLoaded(it, parsedTimestamp, v)
			e.Type = "Image Loaded"

		case "createremotethread":
			messageHtml, message = p.parseCreateRemoteThread(it, parsedTimestamp, v)
			e.Type = "Create Remote Thread"

		case "rawaccessread":
			messageHtml, message = p.parseRawAccessRead(it, parsedTimestamp, v)
			e.Type = "Raw Access Read"

		case "processaccess":
			messageHtml, message = p.parseProcessAccess(it, parsedTimestamp, v)
			e.Type = "Process Access"

		case "filecreate":
			messageHtml, message = p.parseFileCreate(it, parsedTimestamp, v)
			e.Type = "File Create"

		case "registryevent":
			messageHtml, message = p.parseRegistryEvent(it, parsedTimestamp, v)
			e.Type = "Registry"

		case "filecreatestreamhash":
			messageHtml, message = p.parseFileStream(it, parsedTimestamp, v)
			e.Type = "File Stream"

		default:
			logger.Errorf("Unsupported SysMon event: %s", eventName)
			logger.Errorf(`Event Data: %s`, v)
			continue
		}

		if len(message) == 0 {
			continue
		}

		// Add a generic/unified Event record.
		err = p.db.
			InsertInto("event").
			Columns("domain", "host", "utc_time", "type", "plain_text", "html").
			Values(e.Domain, e.Host, e.UtcTime, e.Type, message, messageHtml).
			QueryStruct(&e)

		if err != nil {
			if strings.Contains(err.Error(), "no rows in result set") == false {
				logger.Errorf("Error inserting Event record: %v", err)
				continue
			}
		}
	}
}

//
func (p *Processor) parseProcessCreate(it ImportTask, eventLogTime time.Time, data string) (string, string) {

	regexRes := p.regexData.FindAllStringSubmatch(data, -1)
	if regexRes == nil {
		logger.Errorf(`Cannot locate Data elements for Process Create: %s`, data)
		return "", ""
	}

	pc := new(ProcessCreate)
	pc.Domain = it.Domain
	pc.Host = it.Host
	pc.EventLogTime = eventLogTime

	indexOf := 0
	for _, dataRes := range regexRes {

		// We are only interested if the array has 3 items
		// e.g. main match, data name and then data value
		if len(dataRes) != 3 {
			continue
		}

		// Lower for better matching
		dataRes[DATA_NAME] = strings.ToLower(dataRes[DATA_NAME])

		switch dataRes[DATA_NAME] {
		case "utctime":
			parsedTimestamp, err := time.Parse(LAYOUT_PROCESS_UTC_TIME, strings.TrimSpace(dataRes[DATA_VALUE]))
			if err != nil {
				logger.Error("Unable to parse Process UTC Time: %v (%s)", err, dataRes[DATA_VALUE])
				continue
			}

			pc.UtcTime = parsedTimestamp

		case "processid":
			dataRes[DATA_VALUE] = strings.Map(RemoveNonNumericChars, dataRes[DATA_VALUE])
			pc.ProcessId = goutil.ConvertStringToInt64(dataRes[DATA_VALUE])

		case "image":
			pc.Image = strings.ToLower(dataRes[DATA_VALUE])

		case "commandline":
			pc.CommandLine = strings.ToLower(strings.Replace(dataRes[DATA_VALUE], "&quot;", "\"", -1))

		case "currentdirectory":
			pc.CurrentDirectory = strings.ToLower(dataRes[DATA_VALUE])

		case "hashes":
			indexOf = strings.Index(dataRes[DATA_VALUE], "MD5=")
			if indexOf != -1 {
				pc.Md5 = strings.ToLower(dataRes[DATA_VALUE][indexOf+4 : indexOf+4+32])
			}

			indexOf = strings.Index(dataRes[DATA_VALUE], "SHA256=")
			if indexOf != -1 {
				pc.Sha256 = strings.ToLower(dataRes[DATA_VALUE][indexOf+7 : indexOf+7+64])
			}

		case "parentprocessid":
			dataRes[DATA_VALUE] = strings.Map(RemoveNonNumericChars, dataRes[DATA_VALUE])
			pc.ParentProcessId = goutil.ConvertStringToInt64(dataRes[DATA_VALUE])

		case "parentimage":
			pc.ParentImage = strings.ToLower(dataRes[DATA_VALUE])

		case "parentcommandline":
			pc.ParentCommandLine = strings.ToLower(strings.Replace(dataRes[DATA_VALUE], "&quot;", "\"", -1))
		case "user":
			pc.ProcessUser = strings.ToLower(dataRes[DATA_VALUE])
		}
	}

	html := fmt.Sprintf(`<strong>Process ID:</strong> %d<br><strong>Image:</strong> %s<br><strong>Command Line:</strong> %s<br><strong>Current Directory:</strong> %s<br><strong>MD5:</strong> %s<br><strong>SHA256:</strong> %s<br><strong>Parent Process ID: </strong>%d<br><strong>Parent Image:</strong> %s<br><strong>Parent Command Line:</strong> %s<br><strong>Process User:</strong> %s`,
		pc.ProcessId, pc.Image, pc.CommandLine, pc.CurrentDirectory, pc.Md5, pc.Sha256, pc.ParentProcessId,
		pc.ParentImage, pc.ParentCommandLine, pc.ProcessUser)

	err := p.db.
		InsertInto("process_create").
		Columns("domain", "host", "event_log_time", "utc_time", "process_id", "image", "command_line", "current_directory",
			"md5", "sha256", "parent_process_id", "parent_image", "parent_command_line", "process_user", "html").
		Values(pc.Domain, pc.Host, pc.EventLogTime, pc.UtcTime, pc.ProcessId, pc.Image,
			pc.CommandLine, pc.CurrentDirectory, pc.Md5, pc.Sha256, pc.ParentProcessId,
			pc.ParentImage, pc.ParentCommandLine, pc.ProcessUser, html).
		QueryStruct(&pc)

	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") == false {
			logger.Errorf("Error inserting Process Create record: %v", err)
			return "", ""
		}
	}

	return html,
		fmt.Sprintf(`Process ID: %d Image: %s Command Line: %s Current Directory: %s MD5: %s SHA256: %s Parent Process ID: %d Parent Image: %s Parent Command Line: %s Process User: %s`,
			pc.ProcessId, pc.Image, pc.CommandLine,
			pc.CurrentDirectory, pc.Md5, pc.Sha256, pc.ParentProcessId,
			pc.ParentImage, pc.ParentCommandLine, pc.ProcessUser)
}

//
func (p *Processor) parseFileCreationTime(it ImportTask, eventLogTime time.Time, data string) (string, string) {

	regexRes := p.regexData.FindAllStringSubmatch(data, -1)
	if regexRes == nil {
		logger.Errorf(`Cannot locate Data elements for File Creation Time: %s`, data)
		return "", ""
	}

	fct := new(FileCreationTime)
	fct.Domain = it.Domain
	fct.Host = it.Host
	fct.EventLogTime = eventLogTime

	for _, dataRes := range regexRes {

		// We are only interested if the array has 3 items
		// e.g. main match, data name and then data value
		if len(dataRes) != 3 {
			continue
		}

		// Lower for better matching
		dataRes[DATA_NAME] = strings.ToLower(dataRes[DATA_NAME])

		switch dataRes[DATA_NAME] {
		case "utctime":
			parsedTimestamp, err := time.Parse(LAYOUT_PROCESS_UTC_TIME, strings.TrimSpace(dataRes[DATA_VALUE]))
			if err != nil {
				logger.Error("Unable to parse File Creation Time UTC Time: %v (%s)", err, dataRes[DATA_VALUE])
				continue
			}

			fct.UtcTime = parsedTimestamp

		case "processid":
			dataRes[DATA_VALUE] = strings.Map(RemoveNonNumericChars, dataRes[DATA_VALUE])
			fct.ProcessId = goutil.ConvertStringToInt64(dataRes[DATA_VALUE])

		case "image":
			fct.Image = strings.ToLower(dataRes[DATA_VALUE])

		case "targetfilename":
			fct.TargetFileName = strings.ToLower(dataRes[DATA_VALUE])

		case "creationutctime":
			parsedTimestamp, err := time.Parse(LAYOUT_PROCESS_UTC_TIME, strings.TrimSpace(dataRes[DATA_VALUE]))
			if err != nil {
				logger.Error("Unable to parse File Creation Time Creation UTC Time: %v (%s)", err, dataRes[DATA_VALUE])
				continue
			}

			fct.CreationUtcTime = parsedTimestamp

		case "previouscreationutctime":
			parsedTimestamp, err := time.Parse(LAYOUT_PROCESS_UTC_TIME, strings.TrimSpace(dataRes[DATA_VALUE]))
			if err != nil {
				logger.Errorf("Unable to parse File Creation Time Previous Creation UTC Time: %v (%s)", err, dataRes[DATA_VALUE])
				continue
			}

			fct.PreviousCreationUtcTime = parsedTimestamp
		}
	}

	html := fmt.Sprintf(`<strong>Process ID:</strong> %d<br><strong>Image:</strong> %s<br>
		<strong>Target File Name:</strong> %s<br><strong>Creation Time (UTC):</strong> %s<br>
		<strong>Previous Creation Time (UTC):</strong> %s`,
		fct.ProcessId, fct.Image, fct.TargetFileName, fct.CreationUtcTime.Format("15:04:05 02/01/2006"),
		fct.PreviousCreationUtcTime.Format("15:04:05 02/01/2006"))

	err := p.db.
		InsertInto("file_creation_time").
		Columns("domain", "host", "event_log_time", "utc_time", "process_id", "image", "target_file_name", "creation_utc_time",
			"previous_creation_utc_time", "html").
		Values(fct.Domain, fct.Host, fct.EventLogTime, fct.UtcTime, fct.ProcessId, fct.Image, fct.TargetFileName,
			fct.CreationUtcTime, fct.PreviousCreationUtcTime, html).
		QueryStruct(&fct)

	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") == false {
			logger.Errorf("Error inserting File Creation Time record: %v", err)
			return "", ""
		}
	}

	return html,
		fmt.Sprintf(`Process ID: %d Image: %s Target File Name: %s Creation Time (UTC): %s Previous Creation Time (UTC): %s`,
			fct.ProcessId, fct.Image, fct.TargetFileName,
			fct.CreationUtcTime.Format("15:04:05 02/01/2006"),
			fct.PreviousCreationUtcTime.Format("15:04:05 02/01/2006"))
}

//
func (p *Processor) parseNetworkConnection(it ImportTask, eventLogTime time.Time, data string) (string, string) {

	regexRes := p.regexData.FindAllStringSubmatch(data, -1)
	if regexRes == nil {
		logger.Errorf(`Cannot locate Data elements for Network Connection: %s`, data)
		return "", ""
	}

	nc := new(NetworkConnection)
	nc.Domain = it.Domain
	nc.Host = it.Host
	nc.EventLogTime = eventLogTime

	for _, dataRes := range regexRes {

		// We are only interested if the array has 3 items
		// e.g. main match, data name and then data value
		if len(dataRes) != 3 {
			continue
		}

		// Lower for better matching
		dataRes[DATA_NAME] = strings.ToLower(dataRes[DATA_NAME])

		switch dataRes[DATA_NAME] {
		case "utctime":
			parsedTimestamp, err := time.Parse(LAYOUT_PROCESS_UTC_TIME, strings.TrimSpace(dataRes[DATA_VALUE]))
			if err != nil {
				logger.Errorf("Unable to parse Network Connection UTC Time: %v (%s)", err, dataRes[DATA_VALUE])
				continue
			}

			nc.UtcTime = parsedTimestamp

		case "processid":
			dataRes[DATA_VALUE] = strings.Map(RemoveNonNumericChars, dataRes[DATA_VALUE])
			nc.ProcessId = goutil.ConvertStringToInt64(dataRes[DATA_VALUE])

		case "image":
			nc.Image = strings.ToLower(dataRes[DATA_VALUE])

		case "user":
			nc.ProcessUser = strings.ToLower(dataRes[DATA_VALUE])

		case "protocol":
			nc.Protocol = strings.ToLower(dataRes[DATA_VALUE])

		case "initiated":
			nc.Initiated = goutil.ParseBool(dataRes[DATA_VALUE])

		case "sourceip":
			nc.SourceIp.Scan(dataRes[DATA_VALUE])

		case "sourcehostname":
			nc.SourceHostName = strings.ToLower(dataRes[DATA_VALUE])

		case "sourceport":
			dataRes[DATA_VALUE] = strings.Map(RemoveNonNumericChars, dataRes[DATA_VALUE])
			nc.SourcePort = goutil.ConvertStringToInt32(dataRes[DATA_VALUE])

		case "sourceportname":
			nc.SourcePortName = strings.ToLower(dataRes[DATA_VALUE])

		case "destinationip":
			nc.DestinationIp.Scan(dataRes[DATA_VALUE])

		case "destinationhostname":
			nc.DestinationHostName = strings.ToLower(dataRes[DATA_VALUE])

		case "destinationport":
			dataRes[DATA_VALUE] = strings.Map(RemoveNonNumericChars, dataRes[DATA_VALUE])
			nc.DestinationPort = goutil.ConvertStringToInt32(dataRes[DATA_VALUE])

		case "destinationportname":
			nc.DestinationPortName = strings.ToLower(dataRes[DATA_VALUE])
		}
	}

	html := fmt.Sprintf(`<strong>Process ID:</strong> %d<br><strong>Image:</strong> %s<br>
	<strong>Process User:</strong> %s<br><strong>Protocol:</strong> %s<br><strong>Initiated:</strong> %t<br>
	<strong>Source IP:</strong> %s<br><strong>Source Host Name: </strong>%s<br><strong>Source Port:</strong> %d<br>
	<strong>Source Port Name: </strong>%s<br><strong>Destination IP:</strong> %s<br>
	<strong>Destination Host Name:</strong> %s<br><strong>Destination Port:</strong> %d
	<br><strong>Destination Port Name:</strong> %s`,
		nc.ProcessId, nc.Image, nc.ProcessUser, nc.Protocol, nc.Initiated, nc.SourceIp.String, nc.SourceHostName,
		nc.SourcePort, nc.SourcePortName, nc.DestinationIp.String, nc.DestinationHostName, nc.DestinationPort, nc.DestinationPortName)

	err := p.db.
		InsertInto("network_connection").
		Columns("domain", "host", "event_log_time", "utc_time", "process_id", "image", "process_user", "protocol",
			"initiated", "source_ip", "source_host_name", "source_port", "source_port_name", "destination_ip",
			"destination_host_name", "destination_port", "destination_port_name", "html").
		Values(nc.Domain, nc.Host, nc.EventLogTime, nc.UtcTime, nc.ProcessId, nc.Image,
			nc.ProcessUser, nc.Protocol, nc.Initiated, nc.SourceIp, nc.SourceHostName,
			nc.SourcePort, nc.SourcePortName, nc.DestinationIp, nc.DestinationHostName,
			nc.DestinationPort, nc.DestinationPortName, html).
		QueryStruct(&nc)

	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") == false {
			logger.Errorf("Error inserting Network Connection record: %v", err)
			return "", ""
		}
	}

	return html,
		fmt.Sprintf(`Process ID: %d Image: %s Process User: %s Protocol: %s Initiated: %t Source IP: %s
		Source Host Name: %s Source Port: %d Source Port Name: %s Destination IP: %s Destination Host Name: %s
		Destination Port: %d Destination Port Name: %s`,
			nc.ProcessId, nc.Image, nc.ProcessUser, nc.Protocol,
			nc.Initiated, nc.SourceIp.String, nc.SourceHostName,
			nc.SourcePort, nc.SourcePortName, nc.DestinationIp.String, nc.DestinationHostName,
			nc.DestinationPort, nc.DestinationPortName)
}

//
func (p *Processor) parseProcessTerminate(it ImportTask, eventLogTime time.Time, data string) (string, string) {

	regexRes := p.regexData.FindAllStringSubmatch(data, -1)
	if regexRes == nil {
		logger.Errorf("Cannot locate Data elements for Process Terminate: %s", data)
		return "", ""
	}

	pt := new(ProcessTerminate)
	pt.Domain = it.Domain
	pt.Host = it.Host
	pt.EventLogTime = eventLogTime

	for _, dataRes := range regexRes {

		// We are only interested if the array has 3 items
		// e.g. main match, data name and then data value
		if len(dataRes) != 3 {
			continue
		}

		// Lower for better matching
		dataRes[DATA_NAME] = strings.ToLower(dataRes[DATA_NAME])

		switch dataRes[DATA_NAME] {
		case "utctime":
			parsedTimestamp, err := time.Parse(LAYOUT_PROCESS_UTC_TIME, strings.TrimSpace(dataRes[DATA_VALUE]))
			if err != nil {
				logger.Error("Unable to parse Process Terminate UTC Time: %v (%s)", err, dataRes[DATA_VALUE])
				continue
			}

			pt.UtcTime = parsedTimestamp

		case "processid":
			dataRes[DATA_VALUE] = strings.Map(RemoveNonNumericChars, dataRes[DATA_VALUE])
			pt.ProcessId = goutil.ConvertStringToInt64(dataRes[DATA_VALUE])

		case "image":
			pt.Image = strings.ToLower(dataRes[DATA_VALUE])
		}
	}

	html := fmt.Sprintf(`<strong>Process ID:</strong> %d<br><strong>Image:</strong> %s`, pt.ProcessId, pt.Image)

	err := p.db.
		InsertInto("process_terminate").
		Columns("domain", "host", "event_log_time", "utc_time", "process_id", "image", "html").
		Values(pt.Domain, pt.Host, pt.EventLogTime, pt.UtcTime, pt.ProcessId, pt.Image, html).
		QueryStruct(&pt)

	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") == false {
			logger.Errorf("Error inserting Process Terminate record: %v", err)
			return "", ""
		}
	}

	return html,
		fmt.Sprintf(`Process ID: %d Image: %s`, pt.ProcessId, pt.Image)
}

//
func (p *Processor) parseDriverLoaded(it ImportTask, eventLogTime time.Time, data string) (string, string) {

	regexRes := p.regexData.FindAllStringSubmatch(data, -1)
	if regexRes == nil {
		logger.Errorf("Cannot locate Data elements for Driver Loaded: %s", data)
		return "", ""
	}

	dl := new(DriverLoaded)
	dl.Domain = it.Domain
	dl.Host = it.Host
	dl.EventLogTime = eventLogTime

	indexOf := 0
	for _, dataRes := range regexRes {

		// We are only interested if the array has 3 items
		// e.g. main match, data name and then data value
		if len(dataRes) != 3 {
			continue
		}

		// Lower for better matching
		dataRes[DATA_NAME] = strings.ToLower(dataRes[DATA_NAME])

		switch dataRes[DATA_NAME] {
		case "utctime":
			parsedTimestamp, err := time.Parse(LAYOUT_PROCESS_UTC_TIME, strings.TrimSpace(dataRes[DATA_VALUE]))
			if err != nil {
				logger.Error("Unable to parse Driver Loaded UTC Time: %v (%s)", err, dataRes[DATA_VALUE])
				continue
			}

			dl.UtcTime = parsedTimestamp

		case "imageloaded":
			dl.ImageLoaded = strings.ToLower(dataRes[DATA_VALUE])

		case "hashes":
			indexOf = strings.Index(dataRes[DATA_VALUE], "MD5=")
			if indexOf != -1 {
				dl.Md5 = strings.ToLower(dataRes[DATA_VALUE][indexOf+4 : indexOf+4+32])
			}

			indexOf = strings.Index(dataRes[DATA_VALUE], "SHA256=")
			if indexOf != -1 {
				dl.Sha256 = strings.ToLower(dataRes[DATA_VALUE][indexOf+7 : indexOf+7+64])
			}

		case "signed":
			dl.Signed = goutil.ParseBool(dataRes[DATA_VALUE])

		case "signature":
			dl.Signature = strings.ToLower(dataRes[DATA_VALUE])
		}
	}

	html := fmt.Sprintf(`<strong>Image Loaded:</strong> %s<br><strong>MD5:</strong> %s<br>
		<strong>SHA256:</strong> %s<br><strong>Signed:</strong> %t<br><strong>Signature:</strong> %s`,
		dl.ImageLoaded, dl.Md5, dl.Sha256, dl.Signed, dl.Signature)

	err := p.db.
		InsertInto("driver_loaded").
		Columns("domain", "host", "event_log_time", "utc_time", "image_loaded", "md5", "sha256", "signed", "signature", "html").
		Values(dl.Domain, dl.Host, dl.EventLogTime, dl.UtcTime, dl.ImageLoaded, dl.Md5, dl.Sha256, dl.Signed, dl.Signature, html).
		QueryStruct(&dl)

	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") == false {
			logger.Errorf("Error inserting Driver Loaded record: %v", err)
			return "", ""
		}
	}

	return html,
		fmt.Sprintf(`Image Loaded: %s MD5: %s SHA256: %s Signed: %t Signature: %s`,
			dl.ImageLoaded, dl.Md5, dl.Sha256, dl.Signed, dl.Signature)
}

//
func (p *Processor) parseImageLoaded(it ImportTask, eventLogTime time.Time, data string) (string, string) {

	regexRes := p.regexData.FindAllStringSubmatch(data, -1)
	if regexRes == nil {
		logger.Errorf("Cannot locate Data elements for Image Loaded: %s", data)
		return "", ""
	}

	il := new(ImageLoaded)
	il.Domain = it.Domain
	il.Host = it.Host
	il.EventLogTime = eventLogTime

	indexOf := 0
	for _, dataRes := range regexRes {

		// We are only interested if the array has 3 items
		// e.g. main match, data name and then data value
		if len(dataRes) != 3 {
			continue
		}

		// Lower for better matching
		dataRes[DATA_NAME] = strings.ToLower(dataRes[DATA_NAME])

		switch dataRes[DATA_NAME] {
		case "utctime":
			parsedTimestamp, err := time.Parse(LAYOUT_PROCESS_UTC_TIME, strings.TrimSpace(dataRes[DATA_VALUE]))
			if err != nil {
				logger.Error("Unable to parse Image Loaded UTC Time: %v (%s)", err, dataRes[DATA_VALUE])
				continue
			}

			il.UtcTime = parsedTimestamp
		case "processid":
			dataRes[DATA_VALUE] = strings.Map(RemoveNonNumericChars, dataRes[DATA_VALUE])
			il.ProcessId = goutil.ConvertStringToInt64(dataRes[DATA_VALUE])

		case "image":
			il.Image = strings.ToLower(dataRes[DATA_VALUE])

		case "imageloaded":
			il.ImageLoaded = strings.ToLower(dataRes[DATA_VALUE])

		case "hashes":
			indexOf = strings.Index(dataRes[DATA_VALUE], "MD5=")
			if indexOf != -1 {
				il.Md5 = strings.ToLower(dataRes[DATA_VALUE][indexOf+4 : indexOf+4+32])
			}

			indexOf = strings.Index(dataRes[DATA_VALUE], "SHA256=")
			if indexOf != -1 {
				il.Sha256 = strings.ToLower(dataRes[DATA_VALUE][indexOf+7 : indexOf+7+64])
			}

		case "signed":
			il.Signed = goutil.ParseBool(dataRes[DATA_VALUE])

		case "signature":
			il.Signature = strings.ToLower(dataRes[DATA_VALUE])
		}
	}

	html := fmt.Sprintf(`<strong>Process ID:</strong> %d<br><strong>Image:</strong> %s<br>
		<strong>Image Loaded:</strong> %s<br><strong>MD5:</strong> %s<br><strong>SHA256:</strong> %s<br>
		<strong>Signed:</strong> %t<br><strong>Signature:</strong> %s`,
		il.ProcessId, il.Image, il.ImageLoaded, il.Md5, il.Sha256, il.Signed, il.Signature)

	err := p.db.
		InsertInto("image_loaded").
		Columns("domain", "host", "event_log_time", "utc_time", "process_id", "image",
		"image_loaded", "md5", "sha256", "signed", "signature", "html").
		Values(il.Domain, il.Host, il.EventLogTime, il.UtcTime, il.ProcessId, il.Image,
		il.ImageLoaded, il.Md5, il.Sha256, il.Signed, il.Signature, html).
		QueryStruct(&il)

	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") == false {
			logger.Errorf("Error inserting Image Loaded record: %v", err)
			return "", ""
		}
	}

	return html,
		fmt.Sprintf(`Process ID: %d Image: %s Image Loaded: %s MD5: %s SHA256: %s Signed: %t Signature: %s`,
			il.ProcessId, il.Image, il.ImageLoaded, il.Md5,
			il.Sha256, il.Signed, il.Signature)
}

//
func (p *Processor) parseCreateRemoteThread(it ImportTask, eventLogTime time.Time, data string) (string, string) {

	regexRes := p.regexData.FindAllStringSubmatch(data, -1)
	if regexRes == nil {
		logger.Errorf("Cannot locate Data elements for Create Remote Thread: %s", data)
		return "", ""
	}

	crt := new(CreateRemoteThread)
	crt.Domain = it.Domain
	crt.Host = it.Host
	crt.EventLogTime = eventLogTime

	for _, dataRes := range regexRes {

		// We are only interested if the array has 3 items
		// e.g. main match, data name and then data value
		if len(dataRes) != 3 {
			continue
		}

		// Lower for better matching
		dataRes[DATA_NAME] = strings.ToLower(dataRes[DATA_NAME])

		switch dataRes[DATA_NAME] {
		case "utctime":
			parsedTimestamp, err := time.Parse(LAYOUT_PROCESS_UTC_TIME, strings.TrimSpace(dataRes[DATA_VALUE]))
			if err != nil {
				logger.Error("Unable to parse Create Remote Thread UTC Time: %v (%s)", err, dataRes[DATA_VALUE])
				continue
			}

			crt.UtcTime = parsedTimestamp
		case "sourceprocessid":
			dataRes[DATA_VALUE] = strings.Map(RemoveNonNumericChars, dataRes[DATA_VALUE])
			crt.SourceProcessId = goutil.ConvertStringToInt64(dataRes[DATA_VALUE])

		case "targetprocessid":
			dataRes[DATA_VALUE] = strings.Map(RemoveNonNumericChars, dataRes[DATA_VALUE])
			crt.SourceProcessId = goutil.ConvertStringToInt64(dataRes[DATA_VALUE])

		case "sourceimage":
			crt.SourceImage = strings.ToLower(dataRes[DATA_VALUE])

		case "targetimage":
			crt.TargetImage = strings.ToLower(dataRes[DATA_VALUE])

		case "newthreadid":
			dataRes[DATA_VALUE] = strings.Map(RemoveNonNumericChars, dataRes[DATA_VALUE])
			crt.NewThreadId = goutil.ConvertStringToInt64(dataRes[DATA_VALUE])

		case "startaddress":
			crt.StartAddress = dataRes[DATA_VALUE]

		case "startmodule":
			crt.StartModule = strings.ToLower(dataRes[DATA_VALUE])

		case "startfunction":
			crt.StartFunction = strings.ToLower(dataRes[DATA_VALUE])
		}
	}

	html := fmt.Sprintf(`<strong>Source Process ID:</strong> %d<br><strong>Source Image:</strong> %s<br>
		<strong>Target Process ID:</strong> %d<br><strong>Target Image:</strong> %s<br>
		<strong>New Thread ID:</strong> %d<br><strong>Start Address:</strong> %s<br>
		<strong>Start Module:</strong> %s<br><strong>Start Function:</strong> %s`,
		crt.SourceProcessId, crt.SourceImage, crt.TargetProcessId, crt.TargetImage, crt.NewThreadId,
		crt.StartAddress, crt.StartModule, crt.StartFunction)

	err := p.db.
		InsertInto("create_remote_thread").
		Columns("domain", "host", "event_log_time", "utc_time", "source_process_id", "source_image", "target_process_id",
			"target_image", "new_thread_id", "start_address", "start_module", "start_function", "html").
		Values(crt.Domain, crt.Host, crt.EventLogTime, crt.UtcTime, crt.SourceProcessId, crt.SourceImage, crt.TargetProcessId,
			crt.TargetImage, crt.NewThreadId, crt.StartAddress, crt.StartModule, crt.StartFunction, "html").
		QueryStruct(&crt)

	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") == false {
			logger.Errorf("Error inserting Create Remote Thread record: %v", err)
			return "", ""
		}
	}

	return html,
		fmt.Sprintf(`Source Process ID: %d Source Image: %s Target Process ID: %d Target Image: %s
			New Thread ID: %d Start Address: %s Start Module: %s Start Function: %s`,
			crt.SourceProcessId, crt.SourceImage, crt.TargetProcessId, crt.TargetImage,
			crt.NewThreadId, crt.StartAddress, crt.StartModule, crt.StartFunction)
}

//
func (p *Processor) parseRawAccessRead(it ImportTask, eventLogTime time.Time, data string) (string, string) {

	regexRes := p.regexData.FindAllStringSubmatch(data, -1)
	if regexRes == nil {
		logger.Errorf("Cannot locate Data elements for Raw Access Read: %s", data)
		return "", ""
	}

	ra := new(RawAccess)
	ra.Domain = it.Domain
	ra.Host = it.Host
	ra.EventLogTime = eventLogTime

	for _, dataRes := range regexRes {

		// We are only interested if the array has 3 items
		// e.g. main match, data name and then data value
		if len(dataRes) != 3 {
			continue
		}

		// Lower for better matching
		dataRes[DATA_NAME] = strings.ToLower(dataRes[DATA_NAME])

		switch dataRes[DATA_NAME] {
		case "utctime":
			parsedTimestamp, err := time.Parse(LAYOUT_PROCESS_UTC_TIME, strings.TrimSpace(dataRes[DATA_VALUE]))
			if err != nil {
				logger.Error("Unable to parse Raw Access UTC Time: %v (%s)", err, dataRes[DATA_VALUE])
				continue
			}

			ra.UtcTime = parsedTimestamp

		case "processid":
			dataRes[DATA_VALUE] = strings.Map(RemoveNonNumericChars, dataRes[DATA_VALUE])
			ra.ProcessId = goutil.ConvertStringToInt64(dataRes[DATA_VALUE])

		case "image":
			ra.Image = strings.ToLower(dataRes[DATA_VALUE])

		case "device":
			ra.Device = strings.ToLower(dataRes[DATA_VALUE])

		}
	}

	html := fmt.Sprintf(`<strong>Process ID:</strong> %d<br><strong>Image:</strong> %s<br><strong>Device:</strong> %s`,
		ra.ProcessId, ra.Image, ra.Device)

	err := p.db.
		InsertInto("raw_access").
		Columns("domain", "host", "event_log_time", "utc_time", "process_id", "image", "device", "html").
		Values(ra.Domain, ra.Host, ra.EventLogTime, ra.UtcTime, ra.ProcessId, ra.Image, ra.Device, html).
		QueryStruct(&ra)

	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") == false {
			logger.Errorf("Error inserting Raw Access record: %v", err)
			return "", ""
		}
	}

	return html,
		fmt.Sprintf(`Process ID: %d Image: %s Device: %s`,
			ra.ProcessId, ra.Image, ra.Device)
}

//
func (p *Processor) parseProcessAccess(it ImportTask, eventLogTime time.Time, data string) (string, string) {

	regexRes := p.regexData.FindAllStringSubmatch(data, -1)
	if regexRes == nil {
		logger.Errorf("Cannot locate Data elements for Process Access: %s", data)
		return "", ""
	}

	pa := new(ProcessAccess)
	pa.Domain = it.Domain
	pa.Host = it.Host
	pa.EventLogTime = eventLogTime

	for _, dataRes := range regexRes {

		// We are only interested if the array has 3 items
		// e.g. main match, data name and then data value
		if len(dataRes) != 3 {
			continue
		}

		// Lower for better matching
		dataRes[DATA_NAME] = strings.ToLower(dataRes[DATA_NAME])

		switch dataRes[DATA_NAME] {
		case "utctime":
			parsedTimestamp, err := time.Parse(LAYOUT_PROCESS_UTC_TIME, strings.TrimSpace(dataRes[DATA_VALUE]))
			if err != nil {
				logger.Error("Unable to parse Process Accessed UTC Time: %v (%s)", err, dataRes[DATA_VALUE])
				continue
			}

			pa.UtcTime = parsedTimestamp

		case "sourceprocessid":
			dataRes[DATA_VALUE] = strings.Map(RemoveNonNumericChars, dataRes[DATA_VALUE])
			pa.SourceProcessId = goutil.ConvertStringToInt64(dataRes[DATA_VALUE])

		case "sourceimage":
			pa.SourceImage = strings.ToLower(dataRes[DATA_VALUE])

		case "targetprocessid":
			pa.TargetProcessId = goutil.ConvertStringToInt64(dataRes[DATA_VALUE])

		case "targetimage":
			pa.TargetImage = strings.ToLower(dataRes[DATA_VALUE])

		case "grantedaccess":
			pa.GrantedAccess = strings.ToLower(dataRes[DATA_VALUE])

		case "calltrace":
			pa.CallTrace = strings.ToLower(dataRes[DATA_VALUE])
		}
	}

	html := fmt.Sprintf(`<strong>Source Process ID:</strong> %d<br>
	<strong>Source Image:</strong> %s<strong>Target Process ID:</strong> %d<br><strong>Target Image:</strong> %s<br>
	<strong>Granted Access:</strong> %s<br><strong>Call Trace:</strong> %s<br>`,
		pa.SourceProcessId, pa.SourceImage, pa.TargetProcessId, pa.TargetImage, pa.GrantedAccess, pa.CallTrace)

	err := p.db.
		InsertInto("process_access").
		Columns("domain", "host", "event_log_time", "utc_time", "source_process_id", "source_image",
			"target_process_id", "target_image", "granted_access", "call_trace", "html").
		Values(pa.Domain, pa.Host, pa.EventLogTime, pa.UtcTime, pa.SourceProcessId, pa.SourceImage,
			pa.TargetProcessId, pa.TargetImage, pa.GrantedAccess, pa.CallTrace, html).
		QueryStruct(&pa)

	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") == false {
			logger.Errorf("Error inserting Process Access record: %v", err)
			return "", ""
		}
	}

	return html,
		fmt.Sprintf(`Source Process ID: %d Source Image: %s Target Process ID: %d Target Image: %s`,
			pa.SourceProcessId, pa.SourceImage, pa.TargetProcessId, pa.TargetImage)
}

//
func (p *Processor) parseFileCreate(it ImportTask, eventLogTime time.Time, data string) (string, string) {

	regexRes := p.regexData.FindAllStringSubmatch(data, -1)
	if regexRes == nil {
		logger.Errorf(`Cannot locate Data elements for File Create: %s`, data)
		return "", ""
	}

	fct := new(FileCreate)
	fct.Domain = it.Domain
	fct.Host = it.Host
	fct.EventLogTime = eventLogTime

	for _, dataRes := range regexRes {

		// We are only interested if the array has 3 items
		// e.g. main match, data name and then data value
		if len(dataRes) != 3 {
			continue
		}

		// Lower for better matching
		dataRes[DATA_NAME] = strings.ToLower(dataRes[DATA_NAME])

		switch dataRes[DATA_NAME] {
		case "utctime":
			parsedTimestamp, err := time.Parse(LAYOUT_PROCESS_UTC_TIME, strings.TrimSpace(dataRes[DATA_VALUE]))
			if err != nil {
				logger.Error("Unable to parse File Created UTC Time: %v (%s)", err, dataRes[DATA_VALUE])
				continue
			}

			fct.UtcTime = parsedTimestamp

		case "processid":
			dataRes[DATA_VALUE] = strings.Map(RemoveNonNumericChars, dataRes[DATA_VALUE])
			fct.ProcessId = goutil.ConvertStringToInt64(dataRes[DATA_VALUE])

		case "image":
			fct.Image = strings.ToLower(dataRes[DATA_VALUE])

		case "targetfilename":
			fct.TargetFileName = strings.ToLower(dataRes[DATA_VALUE])

		case "creationutctime":
			parsedTimestamp, err := time.Parse(LAYOUT_PROCESS_UTC_TIME, strings.TrimSpace(dataRes[DATA_VALUE]))
			if err != nil {
				logger.Error("Unable to parse File Created UTC Time: %v (%s)", err, dataRes[DATA_VALUE])
				continue
			}

			fct.CreationUtcTime = parsedTimestamp
		}
	}

	html := fmt.Sprintf(`<strong>Process ID:</strong> %d<br><strong>Image:</strong> %s<br>
	<strong>Target File Name:</strong> %s<br><strong>Creation Time (UTC):</strong> %s<br><strong>`,
		fct.ProcessId, fct.Image, fct.TargetFileName, fct.CreationUtcTime.Format("15:04:05 02/01/2006"))

	err := p.db.
		InsertInto("file_create").
		Columns("domain", "host", "event_log_time", "utc_time", "process_id", "image",
			"target_file_name", "creation_utc_time", "html").
		Values(fct.Domain, fct.Host, fct.EventLogTime, fct.UtcTime, fct.ProcessId, fct.Image, fct.TargetFileName,
			fct.CreationUtcTime, html).
		QueryStruct(&fct)

	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") == false {
			logger.Errorf("Error inserting File Creation Time record: %v", err)
			return "", ""
		}
	}

	return html,
		fmt.Sprintf(`Process ID: %d Image: %s Target File Name: %s Creation Time (UTC): %s`,
			fct.ProcessId, fct.Image, fct.TargetFileName,
			fct.CreationUtcTime.Format("15:04:05 02/01/2006"))
}

//
func (p *Processor) parseRegistryEvent(it ImportTask, eventLogTime time.Time, data string) (string, string) {

	regexRes := p.regexData.FindAllStringSubmatch(data, -1)
	if regexRes == nil {
		logger.Errorf(`Cannot locate Data elements for Registry: %s`, data)
		return "", ""
	}

	var processId int64
	var utcTime time.Time
	var image string
	var eventType string
	var targetObject string
	var details string
	var newName string

	for _, dataRes := range regexRes {

		// We are only interested if the array has 3 items
		// e.g. main match, data name and then data value
		if len(dataRes) != 3 {
			continue
		}

		// Lower for better matching
		dataRes[DATA_NAME] = strings.ToLower(dataRes[DATA_NAME])

		switch dataRes[DATA_NAME] {
		case "utctime":
			parsedTimestamp, err := time.Parse(LAYOUT_PROCESS_UTC_TIME, strings.TrimSpace(dataRes[DATA_VALUE]))
			if err != nil {
				logger.Error("Unable to parse File Created UTC Time: %v (%s)", err, dataRes[DATA_VALUE])
				continue
			}

			utcTime = parsedTimestamp

		case "processid":
			dataRes[DATA_VALUE] = strings.Map(RemoveNonNumericChars, dataRes[DATA_VALUE])
			processId = goutil.ConvertStringToInt64(dataRes[DATA_VALUE])

		case "image":
			image = strings.ToLower(dataRes[DATA_VALUE])

		case "eventtype":
			eventType = strings.ToLower(dataRes[DATA_VALUE])

		case "targetobject":
			targetObject = strings.ToLower(dataRes[DATA_VALUE])

		case "details":
			details = strings.ToLower(dataRes[DATA_VALUE])

		case "newname":
			newName = strings.ToLower(dataRes[DATA_VALUE])
		}
	}

	eventType = strings.ToLower(eventType)

	switch eventType {
	case "createkey", "deletekey", "createvalue", "deletevalue":
		rad := new(RegistryAddDelete)
		rad.Domain = it.Domain
		rad.Host = it.Host
		rad.EventLogTime = eventLogTime
		rad.UtcTime = utcTime
		rad.ProcessId = processId
		rad.Image = image
		rad.EventType = eventType
		rad.TargetObject = targetObject

		return p.insertRegistryAddDeleteRecord(rad)

	case "renamekey", "renamevalue":
		rr := new(RegistryRename)
		rr.Domain = it.Domain
		rr.Host = it.Host
		rr.EventLogTime = eventLogTime
		rr.UtcTime = utcTime
		rr.ProcessId = processId
		rr.Image = image
		rr.EventType = eventType
		rr.TargetObject = targetObject
		rr.NewName = newName

		return p.insertRegistryRenameRecord(rr)

	case "setvalue":
		rs := new(RegistrySet)
		rs.Domain = it.Domain
		rs.Host = it.Host
		rs.EventLogTime = eventLogTime
		rs.UtcTime = utcTime
		rs.ProcessId = processId
		rs.Image = image
		rs.EventType = eventType
		rs.TargetObject = targetObject
		rs.Details = details

		return p.insertRegistrySetValueRecord(rs)

	default:
		logger.Errorf("Unknown Registry event type: %s", eventType)
		return "", ""
	}
}

//
func (p *Processor) insertRegistryAddDeleteRecord(rad *RegistryAddDelete) (string, string) {

	html := fmt.Sprintf(`<strong>Process ID:</strong> %d<br><strong>Image:</strong> %s<br>
			<strong>Event Type:</strong> %s<br><strong>Target Object:</strong> %s`,
		rad.ProcessId, rad.Image, rad.EventType, rad.TargetObject)

	err := p.db.
		InsertInto("registry_add_delete").
		Columns("domain", "host", "event_log_time", "utc_time", "process_id", "image", "event_type", "target_object", "html").
		Values(rad.Domain, rad.Host, rad.EventLogTime, rad.UtcTime, rad.ProcessId, rad.Image, rad.EventType, rad.TargetObject, html).
		QueryStruct(&rad)

	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") == false {
			logger.Errorf("Error inserting Registry Add Delete record: %v", err)
			return "", ""
		}
	}

	return html,
		fmt.Sprintf(`Process ID: %d Image: %s Event Type: %s Target Object: %s`,
			rad.ProcessId, strings.ToLower(rad.Image), rad.EventType, rad.TargetObject)
}

//
func (p *Processor) insertRegistryRenameRecord(rr *RegistryRename) (string, string) {

	html := fmt.Sprintf(`<strong>Process ID:</strong> %d<br><strong>Image:</strong> %s<br>
			<strong>Event Type:</strong> %s<br><strong>Target Object:</strong> %s<br>
			<strong>New Name:</strong> %s`,
		rr.ProcessId, rr.Image, rr.EventType, rr.TargetObject, rr.NewName)

	err := p.db.
		InsertInto("registry_rename").
		Columns("domain", "host", "event_log_time", "utc_time", "process_id", "image", "event_type", "target_object", "new_name", "html").
		Values(rr.Domain, rr.Host, rr.EventLogTime, rr.UtcTime, rr.ProcessId, rr.Image, rr.EventType, rr.TargetObject, rr.NewName, html).
		QueryStruct(&rr)

	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") == false {
			logger.Errorf("Error inserting Registry Rename record: %v", err)
			return "", ""
		}
	}

	return html,
		fmt.Sprintf(`Process ID: %d Image: %s Event Type: %s Target Object: %s New Name: %s`,
			rr.ProcessId, rr.Image, rr.EventType, rr.TargetObject, rr.NewName)
}

//
func (p *Processor) insertRegistrySetValueRecord(rs *RegistrySet) (string, string) {

	html := fmt.Sprintf(`<strong>Process ID:</strong> %d<br><strong>Image:</strong> %s<br>
			<strong>Event Type:</strong> %s<br><strong>Target Object:</strong> %s<br>
			<strong>Details:</strong> %s`,
			rs.ProcessId, rs.Image, rs.EventType, rs.TargetObject, rs.Details)

	err := p.db.
		InsertInto("registry_set").
		Columns("domain", "host", "event_log_time", "utc_time", "process_id", "image", "event_type", "target_object", "details", "html").
		Values(rs.Domain, rs.Host, rs.EventLogTime, rs.UtcTime, rs.ProcessId, rs.Image, rs.EventType, rs.TargetObject, rs.Details, html).
		QueryStruct(&rs)

	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") == false {
			logger.Errorf("Error inserting Registry Set record: %v", err)
			return "", ""
		}
	}

	return html,
		fmt.Sprintf(`Process ID: %d Image: %s Event Type: %s Target Object: %s Details: %s`,
			rs.ProcessId, rs.Image, rs.EventType, rs.TargetObject, rs.Details)
}

//
func (p *Processor) parseFileStream(it ImportTask, eventLogTime time.Time, data string) (string, string) {

	regexRes := p.regexData.FindAllStringSubmatch(data, -1)
	if regexRes == nil {
		logger.Errorf(`Cannot locate Data elements for File Stream: %s`, data)
		return "", ""
	}

	fs := new(FileStream)
	fs.Domain = it.Domain
	fs.Host = it.Host
	fs.EventLogTime = eventLogTime

	indexOf := 0
	for _, dataRes := range regexRes {

		// We are only interested if the array has 3 items
		// e.g. main match, data name and then data value
		if len(dataRes) != 3 {
			continue
		}

		// Lower for better matching
		dataRes[DATA_NAME] = strings.ToLower(dataRes[DATA_NAME])

		switch dataRes[DATA_NAME] {
		case "utctime":
			parsedTimestamp, err := time.Parse(LAYOUT_PROCESS_UTC_TIME, strings.TrimSpace(dataRes[DATA_VALUE]))
			if err != nil {
				logger.Error("Unable to parse File Stream UTC Time: %v (%s)", err, dataRes[DATA_VALUE])
				continue
			}

			fs.UtcTime = parsedTimestamp

		case "processid":
			dataRes[DATA_VALUE] = strings.Map(RemoveNonNumericChars, dataRes[DATA_VALUE])
			fs.ProcessId = goutil.ConvertStringToInt64(dataRes[DATA_VALUE])

		case "image":
			fs.Image = strings.ToLower(dataRes[DATA_VALUE])

		case "targetfilename":
			fs.TargetFileName = strings.ToLower(dataRes[DATA_VALUE])

		case "creationutctime":
			parsedTimestamp, err := time.Parse(LAYOUT_PROCESS_UTC_TIME, strings.TrimSpace(dataRes[DATA_VALUE]))
			if err != nil {
				logger.Error("Unable to parse File Stream Creation UTC Time: %v (%s)", err, dataRes[DATA_VALUE])
				continue
			}

			fs.CreationUtcTime = parsedTimestamp

		case "hash":
			indexOf = strings.Index(dataRes[DATA_VALUE], "MD5=")
			if indexOf != -1 {
				fs.Md5 = strings.ToLower(dataRes[DATA_VALUE][indexOf+4 : indexOf+4+32])
			}

			indexOf = strings.Index(dataRes[DATA_VALUE], "SHA256=")
			if indexOf != -1 {
				fs.Sha256 = strings.ToLower(dataRes[DATA_VALUE][indexOf+7 : indexOf+7+64])
			}
		}
	}

	html := fmt.Sprintf(`<strong>Process ID:</strong> %d<br><strong>Image:</strong> %s<br>
		<strong>Target File Name:</strong> %s<br><strong>Creation UTC Time:</strong> %s<br>
		<strong>MD5:</strong> %s<br><strong>SHA256:</strong> %s`,
		fs.ProcessId, fs.Image, fs.TargetFileName, fs.CreationUtcTime.Format("15:04:05 02/01/2006"), fs.Md5, fs.Sha256)

	err := p.db.
		InsertInto("file_stream").
		Columns("domain", "host", "event_log_time", "utc_time", "process_id", "image",
			"target_file_name", "creation_utc_time", "md5", "sha256", "html").
		Values(fs.Domain, fs.Host, fs.EventLogTime, fs.UtcTime, fs.ProcessId, fs.Image,
			fs.TargetFileName, fs.CreationUtcTime, fs.Md5, fs.Sha256, html).
		QueryStruct(&fs)

	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") == false {
			logger.Errorf("Error inserting File Stream record: %v", err)
			return "", ""
		}
	}

	return html,
		fmt.Sprintf(`Process ID: %d Image: %s Target File Name: %s Creation UTC Time: %s MD5: %s SHA256: %s`,
			fs.ProcessId, fs.Image, fs.TargetFileName,
			fs.CreationUtcTime.Format("15:04:05 02/01/2006"), fs.Md5, fs.Sha256)
}
