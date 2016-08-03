package main

import (
    util "github.com/woanware/goutil"
    "io/ioutil"
    "path"
    "os"
    "time"
    "strings"
)

// ##### Types ##############################################################

type ExportStringStringTotal struct {
    String1     string  `db:"string1"`
    String2     string  `db:"string2"`
    Total       int64   `db:"total"`
}

type ExportStringTotal struct {
    String  string  `db:"string"`
    Total   int64   `db:"total"`
}

type ExportString struct {
    String  string  `db:"string"`
}

// ##### Constants ##########################################################

const (
    SQL_PROCESS_CREATE_PATH_SHA256_GROUP_PATH_ORDER_PATH string = "SELECT image AS string1, sha256 AS string2 , COUNT(image) AS total FROM process_create GROUP BY string1, string2 ORDER BY string1"
    SQL_PROCESS_CREATE_PATH_SHA256_GROUP_PATH_ORDER_HASH string = "SELECT image AS string1, sha256 AS string2 , COUNT(image) AS total FROM process_create GROUP BY string1, string2 ORDER BY string2"
    SQL_PROCESS_CREATE_PATH_MD5_GROUP_PATH_ORDER_PATH string = "SELECT image AS string1, md5 AS string2 , COUNT(image) AS total FROM process_create GROUP BY string1, string2 ORDER BY string1"
    SQL_PROCESS_CREATE_PATH_MD5_GROUP_PATH_ORDER_HASH string = "SELECT image AS string1, md5 AS string2 , COUNT(image) AS total FROM process_create GROUP BY string1, string2 ORDER BY string2"

    SQL_DRIVER_LOADED_PATH_SHA256_GROUP_PATH_ORDER_PATH string = "SELECT image_loaded AS string1, sha256 AS string2, COUNT(image_loaded) AS total FROM driver_loaded GROUP BY string1, string2 ORDER BY string1"
    SQL_DRIVER_LOADED_PATH_SHA256_GROUP_PATH_ORDER_HASH string = "SELECT image_loaded AS string1, sha256 AS string2, COUNT(image_loaded) AS total FROM driver_loaded GROUP BY string1, string2 ORDER BY string2"
    SQL_DRIVER_LOADED_PATH_MD5_GROUP_PATH_ORDER_PATH string = "SELECT image_loaded AS string1, md5 AS string2, COUNT(image_loaded) AS total FROM driver_loaded GROUP BY string1, string2 ORDER BY string1"
    SQL_DRIVER_LOADED_PATH_MD5_GROUP_PATH_ORDER_HASH string = "SELECT image_loaded AS string1, md5 AS string2, COUNT(image_loaded) AS total FROM driver_loaded GROUP BY string1, string2 ORDER BY string2"

    SQL_IMAGE_LOADED_PATH_SHA256_GROUP_PATH_ORDER_PATH string = "SELECT image_loaded AS string1, sha256 AS string2, COUNT(image_loaded) AS total FROM image_loaded GROUP BY string1, string2 ORDER BY string1"
    SQL_IMAGE_LOADED_PATH_SHA256_GROUP_PATH_ORDER_HASH string = "SELECT image_loaded AS string1, sha256 AS string2, COUNT(image_loaded) AS total FROM image_loaded GROUP BY string1, string2 ORDER BY string2"
    SQL_IMAGE_LOADED_PATH_MD5_GROUP_PATH_ORDER_PATH string = "SELECT image_loaded AS string1, md5 AS string2, COUNT(image_loaded) AS total FROM image_loaded GROUP BY string1, string2 ORDER BY string1"
    SQL_IMAGE_LOADED_PATH_MD5_GROUP_PATH_ORDER_HASH string = "SELECT image_loaded AS string1, md5 AS string2, COUNT(image_loaded) AS total FROM image_loaded GROUP BY string1, string2 ORDER BY string2"

    SQL_NETWORK_CONNECTION_DISTINCT_DEST_IP string = "SELECT DISTINCT destination_ip as string FROM network_connection ORDER BY string"
    SQL_NETWORK_CONNECTION_COUNT_DEST_IP string = "SELECT DISTINCT destination_ip as string, COUNT(destination_ip) as total FROM network_connection GROUP BY string ORDER BY string"

    SQL_NETWORK_CONNECTION_DISTINCT_DEST_HOST string = "SELECT DISTINCT destination_host_name as string FROM network_connection ORDER BY string"
    SQL_NETWORK_CONNECTION_COUNT_DEST_HOST string = "SELECT DISTINCT destination_host_name as string, COUNT(destination_host_name) as total FROM network_connection GROUP BY string ORDER BY string"

    SQL_SHA256_ALL string = `SELECT DISTINCT sha256 FROM process_create UNION SELECT DISTINCT sha256 FROM driver_loaded UNION SELECT DISTINCT sha256 FROM image_loaded`
    SQL_SHA256_PROCESS_CREATE string = `SELECT DISTINCT sha256 FROM process_create`
    SQL_SHA256_DRIVER_LOADED string = `SELECT DISTINCT sha256 FROM driver_loaded`
    SQL_SHA256_IMAGE_LOADED string = `SELECT DISTINCT sha256 FROM image_loaded`

    SQL_MD5_ALL string = `SELECT DISTINCT md5 FROM process_create UNION SELECT DISTINCT md5 FROM driver_loaded UNION SELECT DISTINCT md5 FROM image_loaded`
    SQL_MD5_PROCESS_CREATE string = `SELECT DISTINCT md5 FROM process_create`
    SQL_MD5_DRIVER_LOADED string = `SELECT DISTINCT md5 FROM driver_loaded`
    SQL_MD5_IMAGE_LOADED string = `SELECT DISTINCT md5 FROM image_loaded`
)

const (
    PREFIX_EXPORT_PROCESS_CREATE_PATH_SHA256_GROUP_PATH_ORDER_PATH string = "process-create-path-sha256-count-by-path-"
    PREFIX_EXPORT_PROCESS_CREATE_PATH_SHA256_GROUP_PATH_ORDER_HASH string = "process-create-path-sha256-count-by-hash-"
    PREFIX_EXPORT_PROCESS_CREATE_PATH_MD5_GROUP_PATH_ORDER_PATH string = "process-create-path-md5-count-by-path-"
    PREFIX_EXPORT_PROCESS_CREATE_PATH_MD5_GROUP_PATH_ORDER_HASH string = "process-create-path-md5-count-by-hash-"

    PREFIX_EXPORT_DRIVER_LOADED_PATH_SHA256_GROUP_PATH_ORDER_PATH string = "driver-loaded-path-sha256-count-by-path-"
    PREFIX_EXPORT_DRIVER_LOADED_PATH_SHA256_GROUP_PATH_ORDER_HASH string = "driver-loaded-path-sha256-count-by-hash-"
    PREFIX_EXPORT_DRIVER_LOADED_PATH_MD5_GROUP_PATH_ORDER_PATH string = "driver-loaded-path-md5-count-by-path-"
    PREFIX_EXPORT_DRIVER_LOADED_PATH_MD5_GROUP_PATH_ORDER_HASH string = "driver-loaded-path-md5-count-by-hash-"

    PREFIX_EXPORT_IMAGE_LOADED_PATH_SHA256_GROUP_PATH_ORDER_PATH string = "image-loaded-path-sha256-count-by-path-"
    PREFIX_EXPORT_IMAGE_LOADED_PATH_SHA256_GROUP_PATH_ORDER_HASH string = "image-loaded-path-sha256-count-by-hash-"
    PREFIX_EXPORT_IMAGE_LOADED_PATH_MD5_GROUP_PATH_ORDER_PATH string = "image-loaded-path-md5-count-by-path-"
    PREFIX_EXPORT_IMAGE_LOADED_PATH_MD5_GROUP_PATH_ORDER_HASH string = "image-loaded-path-md5-count-by-hash-"

    PREFIX_EXPORT_NETWORK_CONNECTION_DISTINCT_DEST_IP string = "network-connection-destination-ip-"
    PREFIX_EXPORT_NETWORK_CONNECTION_DISTINCT_DEST_IP_COUNT string = "network-connection-destination-ip-count-"
    PREFIX_EXPORT_NETWORK_CONNECTION_DISTINCT_DEST_HOST string = "network-connection-destination-host-"
    PREFIX_EXPORT_NETWORK_CONNECTION_DISTINCT_DEST_HOST_COUNT string = "network-connection-destination-host-count-"

    PREFIX_EXPORT_SHA256_ALL string = "sha256-all-"
    PREFIX_EXPORT_SHA256_PROCESS_CREATE string = "sha256-process-create-"
    PREFIX_EXPORT_SHA256_DRIVER_LOADED string = "sha256-driver-imaged-"
    PREFIX_EXPORT_SHA256_IMAGE_LOADED string = "sha256-image-loaded-"

    PREFIX_EXPORT_MD5_ALL string = "md5-all-"
    PREFIX_EXPORT_MD5_PROCESS_CREATE string = "md5-process-create-"
    PREFIX_EXPORT_MD5_DRIVER_LOADED string = "md5-driver-imaged-"
    PREFIX_EXPORT_MD5_IMAGE_LOADED string = "md5-image-loaded-"
)

const (
    EXPORT_TYPE_PROCESS_CREATE_PATH_SHA256_GROUP_PATH_ORDER_PATH = 1
    EXPORT_TYPE_PROCESS_CREATE_PATH_SHA256_GROUP_PATH_ORDER_HASH = 2
    EXPORT_TYPE_PROCESS_CREATE_PATH_MD5_GROUP_PATH_ORDER_PATH = 3
    EXPORT_TYPE_PROCESS_CREATE_PATH_MD5_GROUP_PATH_ORDER_HASH = 4

    EXPORT_TYPE_DRIVER_LOADED_PATH_SHA256_GROUP_PATH_ORDER_PATH = 5
    EXPORT_TYPE_DRIVER_LOADED_PATH_SHA256_GROUP_PATH_ORDER_HASH = 6
    EXPORT_TYPE_DRIVER_LOADED_PATH_MD5_GROUP_PATH_ORDER_PATH = 7
    EXPORT_TYPE_DRIVER_LOADED_PATH_MD5_GROUP_PATH_ORDER_HASH = 8

    EXPORT_TYPE_IMAGE_LOADED_PATH_SHA256_GROUP_PATH_ORDER_PATH = 9
    EXPORT_TYPE_IMAGE_LOADED_PATH_SHA256_GROUP_PATH_ORDER_HASH = 10
    EXPORT_TYPE_IMAGE_LOADED_PATH_MD5_GROUP_PATH_ORDER_PATH = 11
    EXPORT_TYPE_IMAGE_LOADED_PATH_MD5_GROUP_PATH_ORDER_HASH = 12

    EXPORT_TYPE_NETWORK_CONNECTION_DISTINCT_DEST_IP = 13
    EXPORT_TYPE_NETWORK_CONNECTION_DISTINCT_DEST_IP_COUNT = 14
    EXPORT_TYPE_NETWORK_CONNECTION_DISTINCT_DEST_HOST = 15
    EXPORT_TYPE_NETWORK_CONNECTION_DISTINCT_DEST_HOST_COUNT = 16

    EXPORT_TYPE_SHA256_ALL = 17
    EXPORT_TYPE_SHA256_PROCESS_CREATE = 18
    EXPORT_TYPE_SHA256_DRIVER_LOADED = 19
    EXPORT_TYPE_SHA256_IMAGE_LOADED = 20

    EXPORT_TYPE_MD5_ALL = 21
    EXPORT_TYPE_MD5_PROCESS_CREATE = 22
    EXPORT_TYPE_MD5_DRIVER_LOADED = 23
    EXPORT_TYPE_MD5_IMAGE_LOADED = 24
)

// ##### Methods #############################################################

// Exports data that needs to be grouped by two text fields, and has
// a total/count field e.g. GROUPED BY PATH and SHA256, with a COUNT
func exportDataForStringStringTotal(sql string, typeName string, dataType int, prefix string) {

    rows, err := db.DB.Query(sql)
    if err != nil {
        logger.Errorf("Error querying for %s export: %v", typeName, err)
        return
    }
    defer rows.Close()

    tf, err:= ioutil.TempFile(config.TempDir, "sml-summary-")
    if err != nil {
        logger.Errorf("Error creating temp file for %s export: %v", typeName, err)
        return
    }
    defer tf.Close()

    defer func() {
        if util.DoesFileExist(path.Join(config.TempDir, tf.Name())) == true {
            err := os.Remove(path.Join(config.TempDir, tf.Name()))
            if err != nil {
                logger.Errorf("Error deleting temporary %s export file: %v (%s)", typeName, err, tf.Name)
            }
        }
    }()

    var esst ExportStringStringTotal
    for rows.Next() {
        err = rows.Scan(&esst.String1, &esst.String2, &esst.Total)
        if err != nil {
            logger.Errorf("Error scanning struct for %s export: %v", typeName, err)
            return
        }

        tf.WriteString(esst.String1)
        tf.WriteString(",")
        tf.WriteString(esst.String2)
        tf.WriteString(",")
        tf.WriteString(util.ConvertInt64ToString(esst.Total))
        tf.WriteString("\n")
    }

    timestamp := time.Now().UTC()
    fileName := prefix + timestamp.Format(LAYOUT_DAILY_EXPORT) + ".csv"

    // Move the file
    err = os.Rename(tf.Name(), path.Join(config.ExportDir, fileName))
    if err != nil {
        logger.Errorf("Error moving file to export directory: %v (%s)", err, fileName)
        return
    }

    setExportRecord(dataType, fileName)
}

// Exports data that needs to be grouped by one text field, and
// has a total/count field e.g. GROUPED BY PATH, with a COUNT
func exportDataForStringTotal(sql string, typeName string, dataType int, prefix string) {

    rows, err := db.DB.Query(sql)
    if err != nil {
        logger.Errorf("Error querying for %s export: %v", typeName, err)
        return
    }
    defer rows.Close()

    tf, err:= ioutil.TempFile(config.TempDir, "sml-export-")
    if err != nil {
        logger.Errorf("Error creating temp file for %s export: %v", typeName, err)
        return
    }
    defer tf.Close()

    defer func() {
        if util.DoesFileExist(path.Join(config.TempDir, tf.Name())) == true {
            err := os.Remove(path.Join(config.TempDir, tf.Name()))
            if err != nil {
                logger.Errorf("Error deleting temporary %s export file: %v (%s)", typeName, err, tf.Name)
            }
        }
    }()

    var est ExportStringTotal
    for rows.Next() {
        err = rows.Scan(&est.String, &est.Total)
        if err != nil {
            logger.Errorf("Error scanning struct for %s export: %v", typeName, err)
            return
        }

        tf.WriteString(est.String)
        tf.WriteString(",")
        tf.WriteString(",")
        tf.WriteString(util.ConvertInt64ToString(est.Total))
        tf.WriteString("\n")
    }

    timestamp := time.Now().UTC()
    fileName := prefix + timestamp.Format(LAYOUT_DAILY_EXPORT) + ".csv"

    // Move the file
    err = os.Rename(tf.Name(), path.Join(config.ExportDir, fileName))
    if err != nil {
        logger.Errorf("Error moving file to export directory: %v (%s)", err, fileName)
        return
    }

    setExportRecord(dataType, fileName)
}

// Exports data that needs a DISTINCT string field summarised
func exportDataForString(sql string, typeName string, dataType int, prefix string) {

    rows, err := db.DB.Query(sql)
    if err != nil {
        logger.Errorf("Error querying for %s export: %v", typeName, err)
        return
    }
    defer rows.Close()

    tf, err:= ioutil.TempFile(config.TempDir, "sml-export-")
    if err != nil {
        logger.Errorf("Error creating temp file for %s export: %v", typeName, err)
        return
    }
    defer tf.Close()

    defer func() {
        if util.DoesFileExist(path.Join(config.TempDir, tf.Name())) == true {
            err := os.Remove(path.Join(config.TempDir, tf.Name()))
            if err != nil {
                logger.Errorf("Error deleting temporary %s export file: %v (%s)", typeName, err, tf.Name)
            }
        }
    }()

    var es ExportString
    for rows.Next() {
        err = rows.Scan(&es.String)
        if err != nil {
            logger.Errorf("Error scanning struct for %s export: %v", typeName, err)
            return
        }

        tf.WriteString(es.String)
        tf.WriteString("\n")
    }

    timestamp := time.Now().UTC()
    fileName := prefix + timestamp.Format(LAYOUT_DAILY_EXPORT) + ".csv"

    // Move the file
    err = os.Rename(tf.Name(), path.Join(config.ExportDir, fileName))
    if err != nil {
        logger.Errorf("Error moving file to export directory: %v (%s)", err, fileName)
        return
    }

    setExportRecord(dataType, fileName)
}

//
func setExportRecord(dataType int, fileName string) {

    var e Export

    err := db.
        Select("id, data_type, file_name, updated").
        From("export").
        Where("data_type = $1 and file_name = $2", dataType, fileName).
        OrderBy("id ASC").
        QueryStruct(&e)

    if e.Id > 0 {
        err = db.
            Update("export").
            Set("updated", time.Now().UTC()).
            Where("id = $1", e.Id).
            QueryScalar(&e.Updated)
        return
    }

    err = db.
        InsertInto("export").
        Columns("data_type", "file_name", "updated").
        Values(dataType, fileName, time.Now().UTC()).
        QueryStruct(&e)

    if err != nil {
        if strings.Contains(err.Error(), "no rows in result set") == false {
            logger.Errorf("Error inserting export record: %v (%s, %s, %s)", err, dataType, fileName)
        }
    }
}