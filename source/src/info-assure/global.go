package main

const MESSAGE_DELIMITER string = "###SML###"
const LAYOUT_PROCESS_UTC_TIME = "2006-01-02 15:04:05.000"

//const LAYOUT_EVENT_UTC_TIME = "2006-01-02T15:04:05.000"
const LAYOUT_DAILY_EXPORT = "2006-01-02"

// These constants are the positions in the regex array
const DATA_NAME = 1
const DATA_VALUE = 2

// String values for the database tables e.g. easier to
// loop through them for functionality such as data purge
var databaseTables = []string{
	"create_remote_thread",
	"driver_loaded",
	"event",
	"file_create",
	"file_creation_time",
	"file_stream",
	"image_loaded",
	"network_connection",
	"process_access",
	"process_create",
	"process_terminate",
	"raw_access_read",
	"registry_add_delete",
	"registry_rename",
	"registry_set",
}

