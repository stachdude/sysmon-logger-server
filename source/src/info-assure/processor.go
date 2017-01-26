package main

import (
	_ "github.com/lib/pq"
	"gopkg.in/mgutz/dat.v1/sqlx-runner"
	"strings"
	"github.com/woanware/goutil"
	"time"
	"regexp"
	"fmt"
)

// ##### Types ###############################################################

// Encapsulates a Processor object and its properties
type Processor struct {
	id 						int
	config					*Config
	db 						*runner.DB
	regexEventName			*regexp.Regexp
	regexData				*regexp.Regexp
	regexUtcTime			*regexp.Regexp
}

// Encapsulates the data for an import task
type ImportTask struct {
	Domain	string
	Host 	string
	Data	string
}

// ##### Methods #############################################################

// Constructor/Initialiser for the Processor struct
func NewProcessor(id int, config *Config, db *runner.DB) *Processor {

	p := Processor{
		id:     id,
		config: config,
		db: db,
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

		// Extract the UTC Time from the event
		regexRes = p.regexUtcTime.FindAllStringSubmatch(v, -1)
		if regexRes != nil {
			parsedTimestamp, err := time.Parse(time.RFC3339Nano, strings.TrimSpace(regexRes[0][1]))
			if err != nil {
				logger.Error("Unable to parse event UTC Time: %v (%s)", err, regexRes[0][1])
				continue
			} else {
				e.UtcTime = parsedTimestamp
			}
		}

		if config.Debug == true {
			logger.Infof("Domain: %s, Host: %s, Event: %s", it.Domain, it.Host, eventName)
		}

		switch eventName {
		case "processcreate":
            messageHtml, message = p.parseProcessCreate(it, v)
			e.Type = "Process Create"

		case "filecreatetime":
			messageHtml, message = p.parseFileCreationTime(it, v)
			e.Type = "File Creation Time"

		case "networkconnect":
			messageHtml, message = p.parseNetworkConnection(it, v)
			e.Type = "Network Connection"

		case "processterminate":
            messageHtml, message = p.parseProcessTerminated(it, v)
			e.Type = "Process Terminated"

		case "driverload":
			messageHtml, message = p.parseDriverLoaded(it, v)
			e.Type = "Driver Loaded"

		case "imageload":
			messageHtml, message = p.parseImageLoaded(it, v)
			e.Type = "Image Loaded"

		case "createremotethread":
			messageHtml, message= p.parseCreateRemoteThread(it, v)
			e.Type = "Create Remote Thread"

		case "rawaccessread":
            messageHtml, message = p.parseRawAccessRead(it, v)
			e.Type = "Raw Access Read"

		//case "filecreate":
		//	messageHtml, message = p.parseFileCreationTime(it)
		//	e.Type = "File Create"

		default:
			continue
		}

        if len(message) == 0 {
            continue
        }

		// Add a generic/unified Event record.
		err = p.db.
			InsertInto("event").
			Columns("domain", "host", "utc_time", "type", "message", "message_html").
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
func (p *Processor) parseProcessCreate(it ImportTask, data string) (string, string) {

	regexRes := p.regexData.FindAllStringSubmatch(data, -1)
	if regexRes == nil {
		logger.Errorf(`Cannot locate Data elements for Process Create: %s`, data)
		return "", ""
	}

	pc := new(ProcessCreate)
	pc.Domain = it.Domain
	pc.Host = it.Host

	indexOf := 0
	for _, dataRes := range regexRes {

		// We are only interested if the array has 3 items
		// e.g. main match, data name and then data value
		if len(dataRes) != 3 {
			continue
		}

		// Lower for better matching
		dataRes[DATA_NAME] = strings.ToLower(dataRes[DATA_NAME])

		switch dataRes[DATA_NAME]{
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
			pc.Image = dataRes[DATA_VALUE]

		case "commandline":
			pc.CommandLine = strings.Replace(dataRes[DATA_VALUE], "&quot;", "\"", -1)

		case "currentdirectory":
			pc.CurrentDirectory = dataRes[DATA_VALUE]

		case "hashes":
			indexOf = strings.Index(dataRes[DATA_VALUE], "MD5=")
			if indexOf != -1 {
				pc.Md5 = dataRes[DATA_VALUE][indexOf+4:indexOf+4+32]
			}

			indexOf = strings.Index(dataRes[DATA_VALUE], "SHA256=")
			if indexOf != -1 {
				pc.Sha256 = dataRes[DATA_VALUE][indexOf+7:indexOf+7+64]
			}

		case "parentprocessid":
			dataRes[DATA_VALUE] = strings.Map(RemoveNonNumericChars, dataRes[DATA_VALUE])
			pc.ParentProcessId = goutil.ConvertStringToInt64(dataRes[DATA_VALUE])

		case "parentimage":
			pc.ParentImage = dataRes[DATA_VALUE]

		case "parentcommandline":
			pc.ParentCommandLine = strings.Replace(dataRes[DATA_VALUE], "&quot;", "\"", -1)
		case "user":
			pc.ProcessUser = dataRes[DATA_VALUE]
		}
	}

    err := p.db.
        InsertInto("process_create").
        Columns("domain", "host", "utc_time", "process_id", "image", "command_line", "current_directory",
            "md5", "sha256", "parent_process_id", "parent_image", "parent_command_line", "process_user").
        Values(pc.Domain, pc.Host, pc.UtcTime, pc.ProcessId, pc.Image,
            pc.CommandLine, pc.CurrentDirectory, pc.Md5, pc.Sha256, pc.ParentProcessId,
            pc.ParentImage, pc.ParentCommandLine, pc.ProcessUser).
        QueryStruct(&pc)

    if err != nil {
        if strings.Contains(err.Error(), "no rows in result set") == false {
            logger.Errorf("Error inserting Process Create record: %v", err)
            return "", ""
        }
    }

    return fmt.Sprintf(`<strong>Process ID:</strong> %d<br><strong>Image:</strong> %s<br><strong>Command Line:</strong> %s<br><strong>Current Directory:</strong> %s<br><strong>MD5:</strong> %s<br><strong>SHA256:</strong> %s<br><strong>Parent Process ID: </strong>%d<br><strong>Parent Image:</strong> %s<br><strong>Parent Command Line:</strong> %s<br><strong>Process User:</strong> %s`,
        pc.ProcessId, pc.Image, pc.CommandLine, pc.CurrentDirectory, pc.Md5, pc.Sha256, pc.ParentProcessId,
        pc.ParentImage, pc.ParentCommandLine, pc.ProcessUser),
	fmt.Sprintf(`Process ID: %d Image: %s Command Line: %s Current Directory: %s MD5: %s SHA256: %s Parent Process ID: %d Parent Image: %s Parent Command Line: %s Process User: %s`,
		pc.ProcessId, pc.Image, pc.CommandLine, pc.CurrentDirectory, pc.Md5, pc.Sha256, pc.ParentProcessId,
		pc.ParentImage, pc.ParentCommandLine, pc.ProcessUser)
}

//
func (p *Processor) parseFileCreationTime(it ImportTask, data string) (string, string) {

	regexRes := p.regexData.FindAllStringSubmatch(data, -1)
	if regexRes == nil {
		logger.Errorf(`Cannot locate Data elements for File Creation Time: %s`, data)
		return "", ""
	}

	fct := new(FileCreationTime)
	fct.Domain = it.Domain
	fct.Host = it.Host

	for _, dataRes := range regexRes{

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
			fct.Image = dataRes[DATA_VALUE]

		case "targetfilename":
			fct.TargetFileName = dataRes[DATA_VALUE]

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

	err := p.db.
		InsertInto("file_creation_time").
		Columns("domain", "host", "utc_time", "process_id", "image", "target_file_name", "creation_utc_time",
		"previous_creation_utc_time").
		Values(fct.Domain, fct.Host, fct.UtcTime, fct.ProcessId, fct.Image, fct.TargetFileName,
		fct.CreationUtcTime, fct.PreviousCreationUtcTime).
		QueryStruct(&fct)

	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") == false {
			logger.Errorf("Error inserting File Creation Time record: %v", err)
			return "", ""
		}
	}

	return fmt.Sprintf(`<strong>Process ID:</strong> %d<br><strong>Image:</strong> %s<br><strong>Target File Name:</strong> %s<br><strong>Creation Time (UTC):</strong> %s<br><strong>Previous Creation Time (UTC):</strong> %s`,
		fct.ProcessId, fct.Image, fct.TargetFileName, fct.CreationUtcTime.Format("15:04:05 02/01/2006"),
		fct.PreviousCreationUtcTime.Format("15:04:05 02/01/2006")),
		fmt.Sprintf(`Process ID: %d Image: %s Target File Name: %s Creation Time (UTC): %s Previous Creation Time (UTC): %s`,
			fct.ProcessId, fct.Image, fct.TargetFileName, fct.CreationUtcTime.Format("15:04:05 02/01/2006"),
			fct.PreviousCreationUtcTime.Format("15:04:05 02/01/2006"))
}

//
func (p *Processor) parseNetworkConnection(it ImportTask, data string) (string, string) {

	regexRes := p.regexData.FindAllStringSubmatch(data, -1)
	if regexRes == nil {
		logger.Errorf(`Cannot locate Data elements for Network Connection: %s`, data)
		return "", ""
	}

	nc := new(NetworkConnection)
	nc.Domain = it.Domain
	nc.Host = it.Host

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
			nc.Image = dataRes[DATA_VALUE]

		case "user":
			nc.ProcessUser = dataRes[DATA_VALUE]

		case "protocol":
			nc.Protocol = dataRes[DATA_VALUE]

		case "initiated":
			nc.Initiated = goutil.ParseBool(dataRes[DATA_VALUE])

		case "sourceip":
			nc.SourceIp.Scan(dataRes[DATA_VALUE])

		case "sourcehostname":
			nc.SourceHostName = dataRes[DATA_VALUE]

		case "sourceport":
			dataRes[DATA_VALUE] = strings.Map(RemoveNonNumericChars, dataRes[DATA_VALUE])
			nc.SourcePort = goutil.ConvertStringToInt32(dataRes[DATA_VALUE])

		case "sourceportname":
			nc.SourcePortName = dataRes[DATA_VALUE]

		case "destinationip":
			nc.DestinationIp.Scan(dataRes[DATA_VALUE])

		case "destinationhostname":
			nc.DestinationHostName = dataRes[DATA_VALUE]

		case "destinationport":
			dataRes[DATA_VALUE] = strings.Map(RemoveNonNumericChars, dataRes[DATA_VALUE])
			nc.DestinationPort = goutil.ConvertStringToInt32(dataRes[DATA_VALUE])

		case "destinationportname":
			nc.DestinationPortName = dataRes[DATA_VALUE]
		}
	}

	err := p.db.
		InsertInto("network_connection").
		Columns("domain", "host", "utc_time", "process_id", "image", "process_user", "protocol",
			"initiated", "source_ip", "source_host_name", "source_port", "source_port_name", "destination_ip",
			"destination_host_name", "destination_port", "destination_port_name").
		Values(nc.Domain, nc.Host, nc.UtcTime, nc.ProcessId, nc.Image,
			nc.ProcessUser, nc.Protocol, nc.Initiated, nc.SourceIp, nc.SourceHostName,
			nc.SourcePort, nc.SourcePortName, nc.DestinationIp, nc.DestinationHostName,
			nc.DestinationPort, nc.DestinationPortName).
		QueryStruct(&nc)

	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") == false {
			logger.Errorf("Error inserting Network Connection record: %v", err)
			return "", ""
		}
	}

	return fmt.Sprintf(`<strong>Process ID:</strong> %d<br><strong>Image:</strong> %s<br><strong>Process User:</strong> %s<br><strong>Protocol:</strong> %s<br><strong>Initiated:</strong> %t<br><strong>Source IP:</strong> %s<br><strong>Source Host Name: </strong>%s<br><strong>Source Port:</strong> %d<br><strong>Destination IP:</strong> %s<br><strong>Destination Host Name:</strong> %s<br><strong>Destination Port:</strong> %d<br><strong>Destination Port Name:</strong> %s`,
		nc.ProcessId, nc.Image, nc.ProcessUser, nc.Protocol, nc.Initiated, nc.SourceIp, nc.SourceHostName,
		nc.SourcePort, nc.SourcePortName, nc.DestinationIp, nc.DestinationHostName, nc.DestinationPort, nc.DestinationPortName),
		fmt.Sprintf( `Process ID: %d Image: %s Process User: %s Protocol: %s Initiated: %t Source IP: %s Source Host Name: %s Source Port: %d Destination IP: %s Destination Host Name: %s Destination Port: %d Destination Port Name: %s`,
			nc.ProcessId, nc.Image, nc.ProcessUser, nc.Protocol, nc.Initiated, nc.SourceIp, nc.SourceHostName,
			nc.SourcePort, nc.SourcePortName, nc.DestinationIp, nc.DestinationHostName, nc.DestinationPort, nc.DestinationPortName)
}


//
func (p *Processor) parseProcessTerminated(it ImportTask, data string) (string, string) {

	regexRes := p.regexData.FindAllStringSubmatch(data, -1)
	if regexRes == nil {
		logger.Errorf("Cannot locate Data elements for Process Terminated: %s", data)
		return "", ""
	}

	pt := new(ProcessTerminate)
	pt.Domain = it.Domain
	pt.Host = it.Host

	for _, dataRes := range regexRes {

		// We are only interested if the array has 3 items
		// e.g. main match, data name and then data value
		if len(dataRes) != 3 {
			continue
		}

		// Lower for better matching
		dataRes[DATA_NAME] = strings.ToLower(dataRes[DATA_NAME])

		switch dataRes[DATA_NAME]{
		case "utctime":
			parsedTimestamp, err := time.Parse(LAYOUT_PROCESS_UTC_TIME, strings.TrimSpace(dataRes[DATA_VALUE]))
			if err != nil {
				logger.Error("Unable to parse Process Terminated UTC Time: %v (%s)", err, dataRes[DATA_VALUE])
				continue
			}

			pt.UtcTime = parsedTimestamp

		case "processid":
			dataRes[DATA_VALUE] = strings.Map(RemoveNonNumericChars, dataRes[DATA_VALUE])
			pt.ProcessId = goutil.ConvertStringToInt64(dataRes[DATA_VALUE])

		case "image":
			pt.Image = dataRes[DATA_VALUE]
		}
	}

    err := p.db.
        InsertInto("process_terminated").
        Columns("domain", "host", "utc_time", "process_id", "image").
        Values(pt.Domain, pt.Host, pt.UtcTime, pt.ProcessId, pt.Image).
        QueryStruct(&pt)

    if err != nil {
        if strings.Contains(err.Error(), "no rows in result set") == false {
            logger.Errorf("Error inserting Process Terminated record: %v", err)
            return "", ""
        }
    }

    return fmt.Sprintf(`<strong>Process ID:</strong> %d<br><strong>Image:</strong> %s`, pt.ProcessId, pt.Image),
        fmt.Sprintf(`Process ID: %d Image: %s`, pt.ProcessId, pt.Image)
}

//
func (p *Processor) parseDriverLoaded(it ImportTask, data string) (string, string) {

	regexRes := p.regexData.FindAllStringSubmatch(data, -1)
	if regexRes == nil {
		logger.Errorf("Cannot locate Data elements for Driver Loaded: %s", data)
		return "", ""
	}

	dl := new(DriverLoaded)
	dl.Domain = it.Domain
	dl.Host = it.Host

	indexOf := 0
	for _, dataRes := range regexRes {

		// We are only interested if the array has 3 items
		// e.g. main match, data name and then data value
		if len(dataRes) != 3 {
			continue
		}

		// Lower for better matching
		dataRes[DATA_NAME] = strings.ToLower(dataRes[DATA_NAME])

		switch dataRes[DATA_NAME]{
		case "utctime":
			parsedTimestamp, err := time.Parse(LAYOUT_PROCESS_UTC_TIME, strings.TrimSpace(dataRes[DATA_VALUE]))
			if err != nil {
				logger.Error("Unable to parse Driver Loaded UTC Time: %v (%s)", err, dataRes[DATA_VALUE])
				continue
			}

			dl.UtcTime = parsedTimestamp

		case "imageloaded":
			dl.ImageLoaded = dataRes[DATA_VALUE]

		case "hashes":
			indexOf = strings.Index(dataRes[DATA_VALUE], "MD5=")
			if indexOf != -1 {
				dl.Md5 = dataRes[DATA_VALUE][indexOf+4:indexOf+4+32]
			}

			indexOf = strings.Index(dataRes[DATA_VALUE], "SHA256=")
			if indexOf != -1 {
				dl.Sha256 = dataRes[DATA_VALUE][indexOf+7:indexOf+7+64]
			}

		case "signed":
			dl.Signed = goutil.ParseBool(dataRes[DATA_VALUE])

		case "signature":
			dl.Signature = dataRes[DATA_VALUE]
		}
	}

	err := p.db.
		InsertInto("driver_loaded").
		Columns("domain", "host", "utc_time", "image_loaded", "md5", "sha256", "signed", "signature").
		Values(dl.Domain, dl.Host, dl.UtcTime, dl.ImageLoaded, dl.Md5, dl.Sha256, dl.Signed, dl.Signature).
		QueryStruct(&dl)

	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") == false {
			logger.Errorf("Error inserting Driver Loaded record: %v", err)
			return "", ""
		}
	}

	return fmt.Sprintf(`<strong>Image Loaded:</strong> %s<br><strong>MD5:</strong> %s<br><strong>SHA256:</strong> %s<br><strong>Signed:</strong> %t<br><strong>Signature:</strong> %s`,
		dl.ImageLoaded, dl.Md5, dl.Sha256, dl.Signed, dl.Signature),
		fmt.Sprintf(`Image Loaded: %s MD5: %s SHA256: %s Signed: %t Signature: %s`,
			dl.ImageLoaded, dl.Md5, dl.Sha256, dl.Signed, dl.Signature)
}

//
func (p *Processor) parseImageLoaded(it ImportTask, data string) (string, string) {

	regexRes := p.regexData.FindAllStringSubmatch(data, -1)
	if regexRes == nil {
		logger.Errorf("Cannot locate data: %s", data)
		return "", ""
	}

	il := new(ImageLoaded)
	il.Domain = it.Domain
	il.Host = it.Host

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
			il.Image = dataRes[DATA_VALUE]

		case "imageloaded":
			il.ImageLoaded = dataRes[DATA_VALUE]

		case "hashes":
			indexOf = strings.Index(dataRes[DATA_VALUE], "MD5=")
			if indexOf != -1 {
				il.Md5 = dataRes[DATA_VALUE][indexOf+4:indexOf+4+32]
			}

			indexOf = strings.Index(dataRes[DATA_VALUE], "SHA256=")
			if indexOf != -1 {
				il.Sha256 = dataRes[DATA_VALUE][indexOf+7:indexOf+7+64]
			}

		case "signed":
			il.Signed = goutil.ParseBool(dataRes[DATA_VALUE])

		case "signature":
			il.Signature = dataRes[DATA_VALUE]
		}
	}

	err := p.db.
		InsertInto("image_loaded").
		Columns("domain", "host", "utc_time", "process_id", "image", "image_loaded", "md5", "sha256", "signed", "signature").
		Values(il.Domain, il.Host, il.UtcTime, il.ProcessId, il.Image, il.ImageLoaded, il.Md5, il.Sha256, il.Signed, il.Signature).
		QueryStruct(&il)

	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") == false {
			logger.Errorf("Error inserting Image Loaded record: %v", err)
			return "", ""
		}
	}

	return fmt.Sprintf(`<strong>Process ID:</strong> %d<br><strong>Image:</strong> %s<br><strong>Image Loaded:</strong> %s<br><strong>MD5:</strong> %s<br><strong>SHA256:</strong> %s<br><strong>Signed:</strong> %t<br><strong>Signature:</strong> %s`,
		il.ProcessId, il.Image, il.ImageLoaded, il.Md5, il.Sha256, il.Signed, il.Signature),
		fmt.Sprintf(`Process ID: %d Image: %s Image Loaded: %s MD5: %s SHA256: %s Signed: %t Signature: %s`,
			il.ProcessId, il.Image, il.ImageLoaded, il.Md5, il.Sha256, il.Signed, il.Signature)
}

//
func (p *Processor) parseCreateRemoteThread(it ImportTask, data string) (string, string) {

	regexRes := p.regexData.FindAllStringSubmatch(data, -1)
	if regexRes == nil {
		logger.Errorf("Cannot locate data: %s", data)
		return "", ""
	}

	crt := new(CreateRemoteThread)
	crt.Domain = it.Domain
	crt.Host = it.Host

	for _, dataRes := range regexRes{

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
			crt.SourceImage = dataRes[DATA_VALUE]

		case "targetimage":
			crt.TargetImage = dataRes[DATA_VALUE]

		case "newthreadid":
			dataRes[DATA_VALUE] = strings.Map(RemoveNonNumericChars, dataRes[DATA_VALUE])
			crt.NewThreadId = goutil.ConvertStringToInt64(dataRes[DATA_VALUE])

		case "startaddress":
			crt.StartAddress = dataRes[DATA_VALUE]

		case "startmodule":
			crt.StartModule = dataRes[DATA_VALUE]

		case "startfunction":
			crt.StartFunction = dataRes[DATA_VALUE]
		}
	}

	err := p.db.
		InsertInto("image_loaded").
		Columns("domain", "host", "utc_time", "source_process_id", "source_image", "target_process_id",
		"target_image", "new_thread_id", "start_address", "start_module", "start_function").
		Values(crt.Domain, crt.Host, crt.UtcTime, crt.SourceProcessId, crt.SourceImage, crt.TargetProcessId,
		crt.TargetImage, crt.NewThreadId, crt.StartAddress, crt.StartModule, crt.StartFunction).
		QueryStruct(&crt)

	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") == false {
			logger.Errorf("Error inserting Create Remote Thread record: %v", err)
			return "", ""
		}
	}

	return fmt.Sprintf(`<strong>Source Process ID:</strong> %d<br><strong>Source Image:</strong> %s<br><strong>Target Process ID:</strong> %d<br><strong>Target Image:</strong> %s<br><strong>New Thread ID:</strong> %d<br><strong>Start Address:</strong> %s<br><strong>Start Module:</strong> %s<br><strong>Start Function:</strong> %s`,
		crt.SourceProcessId, crt.SourceImage, crt.TargetProcessId, crt.TargetImage, crt.NewThreadId, crt.StartAddress, crt.StartModule, crt.StartFunction),
		fmt.Sprintf(`Source Process ID: %d Source Image: %s Target Process ID: %d Target Image: %s New Thread ID: %d Start Address: %s Start Module: %s Start Function: %s`,
			crt.SourceProcessId, crt.SourceImage, crt.TargetProcessId, crt.TargetImage, crt.NewThreadId, crt.StartAddress, crt.StartModule, crt.StartFunction)
}

//
func (p *Processor) parseRawAccessRead(it ImportTask, data string) (string, string) {

	regexRes := p.regexData.FindAllStringSubmatch(data, -1)
	if regexRes == nil {
		logger.Errorf("Cannot locate data: %s", data)
		return "", ""
	}

	ra := new(RawAccess)
	ra.Domain = it.Domain
	ra.Host = it.Host

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
			ra.Image = dataRes[DATA_VALUE]

		case "device":
			ra.Device = dataRes[DATA_VALUE]

		}
	}

    err := p.db.
        InsertInto("raw_access").
        Columns("domain", "host", "utc_time", "process_id", "image", "device").
        Values(ra.Domain, ra.Host, ra.UtcTime, ra.ProcessId, ra.Image, ra.Device).
        QueryStruct(&ra)

    if err != nil {
        if strings.Contains(err.Error(), "no rows in result set") == false {
            logger.Errorf("Error inserting Raw Access record: %v", err)
            return "", ""
        }
    }

    return fmt.Sprintf(`<strong>Process ID:</strong> %d<br><strong>Image:</strong> %s<br><strong>Device:</strong> %s`,
		ra.ProcessId, ra.Image, ra.Device),
    fmt.Sprintf(`Process ID: %d Image: %s Device: %s`,
		ra.ProcessId, ra.Image, ra.Device)
}
