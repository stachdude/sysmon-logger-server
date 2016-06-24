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
	regRes 					[][]string
	parts					[]string
	lines 					[]string
	db 						*runner.DB
	regexFormattedMessage	*regexp.Regexp
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

	p.regexFormattedMessage, _ = regexp.Compile(`FormattedMessage="(?s)(.*?)"`)
	p.regexUtcTime, _ = regexp.Compile(`UtcTime="(.*?)"`)

	return &p
}

// Process an individual set of host data
func (p Processor) Process(it ImportTask) {
	var err error

	// Split the data into each separate message using the special delimiter :-)
	p.parts = strings.Split(it.Data, MESSAGE_DELIMITER)

	eventName := ""
	temp := ""
    message := ""
    messageHtml := ""
	var e Event

	for _, v := range p.parts {
		p.regRes = p.regexFormattedMessage.FindAllStringSubmatch(v, -1)
		if p.regRes == nil {
			// Cannot locate the formatted message bit
			logger.Errorf(`Cannot locate FormattedMessage XML: %s`, v)
			continue
		}

		// Hold just the formatted message
        temp = p.regRes[0][1]
		// Split the message into lines
		p.lines = strings.Split(temp, "\n")

		// Set the various Event structure's properties
		e = Event{}
		e.Domain = it.Domain
		e.Host = it.Host

		// Extract the UTC Time from the event
		p.regRes = p.regexUtcTime.FindAllStringSubmatch(v, -1)
		if p.regRes != nil {
			parsedTimestamp, err := time.Parse(LAYOUT_PROCESS_UTC_TIME, strings.TrimSpace(p.regRes[0][1]))
			if err != nil {
				logger.Error("Unable to parse Event UTC Time: %v (%s)", err, p.regRes[0][1])
			} else {
				e.UtcTime = parsedTimestamp
			}
		}

		// Get the first line
		eventName = strings.TrimSpace(p.lines[0])
		// Lowercase for better matching
		eventName = strings.ToLower(eventName)
		// Remove any non-alphanumeric chars
		eventName = strings.Map(RemoveNonAlphaNumericChars, eventName)

		logger.Errorf("Event: %s", eventName)
		switch (eventName) {
		case "processcreate":
            messageHtml, message = p.parseProcessCreate(it)
			e.Type = "Process Create"
		case "processterminated":
            messageHtml, message = p.parseProcessTerminated(it)
			e.Type = "Process Terminated"
		case "networkconnectiondetected":
            messageHtml, message = p.parseNetworkConnection(it)
			e.Type = "Network Connection Detected"
		case "rawaccessreaddetected":
            messageHtml, message = p.parseRawAccessRead(it)
			e.Type = "Raw Access Read"
		case "filecreationtimechanged":
            messageHtml, message = p.parseFileCreationTime(it)
			e.Type = "File Creation Time Changed"
		case "driverloaded":
            messageHtml, message = p.parseDriverLoaded(it)
			e.Type = "Driver Loaded"
		case "imageloaded":
            messageHtml, message = p.parseImageLoaded(it)
			e.Type = "Image Loaded"
		case "createremotethreaddetected":
            messageHtml, message= p.parseCreateRemoteThread(it)
			e.Type = "Create Remote Thread"
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

	// This is probably unnecessary
	p.parts = []string{}
	p.lines = []string{}
	p.regRes = [][]string{}
}

//
func (p *Processor) parseProcessCreate(it ImportTask) (string, string) {
	fieldName := ""
	fieldValue := ""
	indexOf := 0

	pc := new(ProcessCreate)
    pc.Domain = it.Domain
    pc.Host = it.Host

	for k, v := range p.lines {
		// Ignore the first line as it contains the event type that we have already parsed
		if k == 0 {
			continue
		}

		indexOf = strings.Index(v, ":")
		if indexOf == 0 {
			continue
		}
		fieldName = strings.ToLower(v[0:indexOf])
		fieldValue = v[indexOf +1:]
		fieldValue = strings.TrimSpace(fieldValue)

		switch (fieldName){
		case "utctime":
			parsedTimestamp, err := time.Parse(LAYOUT_PROCESS_UTC_TIME, strings.TrimSpace(fieldValue))
			if err != nil {
				logger.Error("Unable to parse Process UTC Time: %v (%s)", err, fieldValue)
				continue
			}

			pc.UtcTime = parsedTimestamp
		case "processid":
			fieldValue = strings.Map(RemoveNonNumericChars, fieldValue)
			pc.ProcessId = goutil.ConvertStringToInt64(fieldValue)
		case "image":
			pc.Image = fieldValue
		case "commandline":
			pc.CommandLine = strings.Replace(fieldValue, "&quot;", "\"", -1)
		case "currentdirectory":
			pc.CurrentDirectory = fieldValue
		case "hashes":
			indexOf = strings.Index(fieldValue, "MD5=")
			if indexOf != -1 {
				pc.Md5 = fieldValue[indexOf+4:indexOf+4+32]
			}

			indexOf = strings.Index(fieldValue, "SHA256=")
			if indexOf != -1 {
				pc.Sha256 = fieldValue[indexOf+7:indexOf+7+64]
			}
		case "parentprocessid":
			fieldValue = strings.Map(RemoveNonNumericChars, fieldValue)
			pc.ParentProcessId = goutil.ConvertStringToInt64(fieldValue)
		case "parentimage":
			pc.ParentImage = fieldValue
		case "parentcommandline":
			pc.ParentCommandLine = strings.Replace(fieldValue, "&quot;", "\"", -1)
		case "user":
			pc.ProcessUser = fieldValue
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
func (p *Processor) parseProcessTerminated(it ImportTask) (string, string) {
	fieldName := ""
	fieldValue := ""
	indexOf := 0

	pt := new(ProcessTerminate)

	for k, v := range p.lines {
		// Ignore the first line as it contains the event type that we have already parsed
		if k == 0 {
			continue
		}

		indexOf = strings.Index(v, ":")
		if indexOf == 0 {
			continue
		}
		fieldName = strings.ToLower(v[0:indexOf])
		fieldValue = v[indexOf +1:]
		fieldValue = strings.TrimSpace(fieldValue)

		switch (fieldName){
		case "utctime":
			parsedTimestamp, err := time.Parse(LAYOUT_PROCESS_UTC_TIME, strings.TrimSpace(fieldValue))
			if err != nil {
				logger.Error("Unable to parse Process Terminated UTC Time: %v (%s)", err, fieldValue)
				continue
			}

			pt.UtcTime = parsedTimestamp
		case "processid":
			fieldValue = strings.Map(RemoveNonNumericChars, fieldValue)
			pt.ProcessId = goutil.ConvertStringToInt64(fieldValue)
		case "image":
			pt.Image = fieldValue
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
func (p *Processor) parseNetworkConnection(it ImportTask) (string, string) {
	fieldName := ""
	fieldValue := ""
	indexOf := 0

	nc := new(NetworkConnection)

	for k, v := range p.lines {
		// Ignore the first line as it contains the event type that we have already parsed
		if k == 0 {
			continue
		}

		indexOf = strings.Index(v, ":")
		if indexOf == 0 {
			continue
		}
		fieldName = strings.ToLower(v[0:indexOf])
		fieldValue = v[indexOf +1:]
		fieldValue = strings.TrimSpace(fieldValue)

		switch (fieldName){
		case "utctime":
			parsedTimestamp, err := time.Parse(LAYOUT_PROCESS_UTC_TIME, strings.TrimSpace(fieldValue))
			if err != nil {
				logger.Error("Unable to parse Network Connection UTC Time: %v (%s)", err, fieldValue)
				continue
			}

			nc.UtcTime = parsedTimestamp
		case "processid":
			fieldValue = strings.Map(RemoveNonNumericChars, fieldValue)
			nc.ProcessId = goutil.ConvertStringToInt64(fieldValue)
		case "image":
			nc.Image = fieldValue
		case "user":
			nc.ProcessUser = fieldValue
		case "protocol":
			nc.Protocol = fieldValue
		case "initiated":
			nc.Initiated = goutil.ParseBool(fieldValue)
		case "source_ip":
			nc.SourceIp = fieldValue
		case "source_host_name":
			nc.SourceHostName =fieldValue
		case "source_port":
			fieldValue = strings.Map(RemoveNonNumericChars, fieldValue)
			nc.SourcePort = goutil.ConvertStringToInt32(fieldValue)
		case "source_port_name":
			nc.SourcePortName = fieldValue
		case "destination_ip":
			nc.DestinationIp = fieldValue
		case "destination_host_name":
			nc.DestinationHostName = fieldValue
		case "destination_port":
			fieldValue = strings.Map(RemoveNonNumericChars, fieldValue)
			nc.DestinationPort = goutil.ConvertStringToInt32(fieldValue)
		case "destination_port_name":
			nc.DestinationPortName = fieldValue
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
func (p *Processor) parseRawAccessRead(it ImportTask) (string, string) {
	fieldName := ""
	fieldValue := ""
	indexOf := 0

	rawAccess := new(RawAccess)

	for k, v := range p.lines {
		// Ignore the first line as it contains the event type that we have already parsed
		if k == 0 {
			continue
		}

		indexOf = strings.Index(v, ":")
		if indexOf == 0 {
			continue
		}
		fieldName = strings.ToLower(v[0:indexOf])
		fieldValue = v[indexOf +1:]
		fieldValue = strings.TrimSpace(fieldValue)

		switch (fieldName){
		case "utctime":
			parsedTimestamp, err := time.Parse(LAYOUT_PROCESS_UTC_TIME, strings.TrimSpace(fieldValue))
			if err != nil {
				logger.Error("Unable to parse Raw Access UTC Time: %v (%s)", err, fieldValue)
				continue
			}

			rawAccess.UtcTime = parsedTimestamp
		case "processid":
			fieldValue = strings.Map(RemoveNonNumericChars, fieldValue)
			rawAccess.ProcessId = goutil.ConvertStringToInt64(fieldValue)
		case "image":
			rawAccess.Image = fieldValue
		case "device":
			rawAccess.Device = fieldValue
		}
	}

    err := p.db.
        InsertInto("raw_access").
        Columns("domain", "host", "utc_time", "process_id", "image", "device").
        Values(rawAccess.Domain, rawAccess.Host, rawAccess.UtcTime,
            rawAccess.ProcessId, rawAccess.Image, rawAccess.Device).
        QueryStruct(&rawAccess)

    if err != nil {
        if strings.Contains(err.Error(), "no rows in result set") == false {
            logger.Errorf("Error inserting Raw Access record: %v", err)
            return "", ""
        }
    }

    return fmt.Sprintf(`<strong>Process ID:</strong> %d<br><strong>Image:</strong> %s<br><strong>Device:</strong> %s`,
        rawAccess.ProcessId, rawAccess.Image, rawAccess.Device),
    fmt.Sprintf(`Process ID: %d Image: %s Device: %s`,
        rawAccess.ProcessId, rawAccess.Image, rawAccess.Device)
}

//
func (p *Processor) parseFileCreationTime(it ImportTask) (string, string) {
	fieldName := ""
	fieldValue := ""
	indexOf := 0

	fct := new(FileCreationTime)

	for k, v := range p.lines {
		// Ignore the first line as it contains the event type that we have already parsed
		if k == 0 {
			continue
		}

		indexOf = strings.Index(v, ":")
		if indexOf == 0 {
			continue
		}
		fieldName = strings.ToLower(v[0:indexOf])
		fieldValue = v[indexOf +1:]
		fieldValue = strings.TrimSpace(fieldValue)

		switch (fieldName){
		case "utctime":
			parsedTimestamp, err := time.Parse(LAYOUT_PROCESS_UTC_TIME, strings.TrimSpace(fieldValue))
			if err != nil {
				logger.Error("Unable to parse File Creation Time UTC Time: %v (%s)", err, fieldValue)
				continue
			}

			fct.UtcTime = parsedTimestamp
		case "processid":
			fieldValue = strings.Map(RemoveNonNumericChars, fieldValue)
			fct.ProcessId = goutil.ConvertStringToInt64(fieldValue)
		case "image":
			fct.Image = fieldValue
		case "targetfilename":
			fct.TargetFileName = fieldValue
		case "creationutctime":
			parsedTimestamp, err := time.Parse(LAYOUT_PROCESS_UTC_TIME, strings.TrimSpace(fieldValue))
			if err != nil {
				logger.Error("Unable to parse File Creation Time Creation UTC Time: %v (%s)", err, fieldValue)
				continue
			}

			fct.CreationUtcTime = parsedTimestamp
		case "previouscreationutctime":
			parsedTimestamp, err := time.Parse(LAYOUT_PROCESS_UTC_TIME, strings.TrimSpace(fieldValue))
			if err != nil {
				logger.Errorf("Unable to parse File Creation Time Previous Creation UTC Time: %v (%s)", err, fieldValue)
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
func (p *Processor) parseDriverLoaded(it ImportTask) (string, string) {
	fieldName := ""
	fieldValue := ""
	indexOf := 0

	dl := new(DriverLoaded)

	for k, v := range p.lines {
		// Ignore the first line as it contains the event type that we have already parsed
		if k == 0 {
			continue
		}

		indexOf = strings.Index(v, ":")
		if indexOf == 0 {
			continue
		}
		fieldName = strings.ToLower(v[0:indexOf])
		fieldValue = v[indexOf +1:]
		fieldValue = strings.TrimSpace(fieldValue)

		switch (fieldName){
		case "utctime":
			parsedTimestamp, err := time.Parse(LAYOUT_PROCESS_UTC_TIME, strings.TrimSpace(fieldValue))
			if err != nil {
				logger.Error("Unable to parse Driver Loaded UTC Time: %v (%s)", err, fieldValue)
				continue
			}

			dl.UtcTime = parsedTimestamp
		case "imageloaded":
			dl.ImageLoaded = fieldValue
		case "hashes":
			indexOf = strings.Index(fieldValue, "MD5=")
			if indexOf != -1 {
				dl.Md5 = fieldValue[indexOf+4:indexOf+4+32]
			}

			indexOf = strings.Index(fieldValue, "SHA256=")
			if indexOf != -1 {
				dl.Sha256 = fieldValue[indexOf+7:indexOf+7+64]
			}
		case "signed":
			dl.Signed = goutil.ParseBool(fieldValue)
		case "signature":
			dl.Signature = fieldValue
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
func (p *Processor) parseImageLoaded(it ImportTask) (string, string) {
	fieldName := ""
	fieldValue := ""
	indexOf := 0

	il := new(ImageLoaded)

	for k, v := range p.lines {
		// Ignore the first line as it contains the event type that we have already parsed
		if k == 0 {
			continue
		}

		indexOf = strings.Index(v, ":")
		if indexOf == 0 {
			continue
		}
		fieldName = strings.ToLower(v[0:indexOf])
		fieldValue = v[indexOf +1:]
		fieldValue = strings.TrimSpace(fieldValue)

		switch (fieldName){
		case "utctime":
			parsedTimestamp, err := time.Parse(LAYOUT_PROCESS_UTC_TIME, strings.TrimSpace(fieldValue))
			if err != nil {
				logger.Error("Unable to parse Image Loaded UTC Time: %v (%s)", err, fieldValue)
				continue
			}

			il.UtcTime = parsedTimestamp
		case "processid":
			fieldValue = strings.Map(RemoveNonNumericChars, fieldValue)
			il.ProcessId = goutil.ConvertStringToInt64(fieldValue)
		case "image":
			il.Image = fieldValue
		case "imageloaded":
			il.ImageLoaded = fieldValue
		case "hashes":
			indexOf = strings.Index(fieldValue, "MD5=")
			if indexOf != -1 {
				il.Md5 = fieldValue[indexOf+4:indexOf+4+32]
			}

			indexOf = strings.Index(fieldValue, "SHA256=")
			if indexOf != -1 {
				il.Sha256 = fieldValue[indexOf+7:indexOf+7+64]
			}
		case "signed":
			il.Signed = goutil.ParseBool(fieldValue)
		case "signature":
			il.Signature = fieldValue
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
func (p *Processor) parseCreateRemoteThread(it ImportTask) (string, string) {
	fieldName := ""
	fieldValue := ""
	indexOf := 0

	crt := new(CreateRemoteThread)

	for k, v := range p.lines {
		// Ignore the first line as it contains the event type that we have already parsed
		if k == 0 {
			continue
		}

		indexOf = strings.Index(v, ":")
		if indexOf == 0 {
			continue
		}
		fieldName = strings.ToLower(v[0:indexOf])
		fieldValue = v[indexOf +1:]
		fieldValue = strings.TrimSpace(fieldValue)

		switch (fieldName){
		case "utctime":
			parsedTimestamp, err := time.Parse(LAYOUT_PROCESS_UTC_TIME, strings.TrimSpace(fieldValue))
			if err != nil {
				logger.Error("Unable to parse Create Remote Thread UTC Time: %v (%s)", err, fieldValue)
				continue
			}

			crt.UtcTime = parsedTimestamp
		case "sourceprocessid":
			fieldValue = strings.Map(RemoveNonNumericChars, fieldValue)
			crt.SourceProcessId = goutil.ConvertStringToInt64(fieldValue)
		case "targetprocessid":
			fieldValue = strings.Map(RemoveNonNumericChars, fieldValue)
			crt.SourceProcessId = goutil.ConvertStringToInt64(fieldValue)
		case "sourceimage":
			crt.SourceImage = fieldValue
		case "targetimage":
			crt.TargetImage = fieldValue
		case "newthreadid":
			fieldValue = strings.Map(RemoveNonNumericChars, fieldValue)
			crt.NewThreadId = goutil.ConvertStringToInt64(fieldValue)
		case "startaddress":
			crt.StartAddress = fieldValue
		case "startmodule":
			crt.StartModule = fieldValue
		case "startfunction":
			crt.StartFunction = fieldValue
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