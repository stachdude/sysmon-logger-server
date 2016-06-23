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

type StringStringTotal struct {
    String1     string  `db:"string1"`
    String2     string  `db:"string2"`
    Total       int64   `db:"total"`
}

type StringTotal struct {
    String  string  `db:"string"`
    Total   int64   `db:"total"`
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
)

const (
    PREFIX_SUMMARY_PROCESS_CREATE_PATH_SHA256_GROUP_PATH_ORDER_PATH string = "process-create-path-sha256-count-by-path-"
    PREFIX_SUMMARY_PROCESS_CREATE_PATH_SHA256_GROUP_PATH_ORDER_HASH string = "process-create-path-sha256-count-by-hash-"
    PREFIX_SUMMARY_PROCESS_CREATE_PATH_MD5_GROUP_PATH_ORDER_PATH string = "process-create-path-md5-count-by-path-"
    PREFIX_SUMMARY_PROCESS_CREATE_PATH_MD5_GROUP_PATH_ORDER_HASH string = "process-create-path-md5-count-by-hash-"

    PREFIX_SUMMARY_DRIVER_LOADED_PATH_SHA256_GROUP_PATH_ORDER_PATH string = "driver-loaded-path-sha256-count-by-path-"
    PREFIX_SUMMARY_DRIVER_LOADED_PATH_SHA256_GROUP_PATH_ORDER_HASH string = "driver-loaded-path-sha256-count-by-hash-"
    PREFIX_SUMMARY_DRIVER_LOADED_PATH_MD5_GROUP_PATH_ORDER_PATH string = "driver-loaded-path-md5-count-by-path-"
    PREFIX_SUMMARY_DRIVER_LOADED_PATH_MD5_GROUP_PATH_ORDER_HASH string = "driver-loaded-path-md5-count-by-hash-"

    PREFIX_SUMMARY_IMAGE_LOADED_PATH_SHA256_GROUP_PATH_ORDER_PATH string = "image-loaded-path-sha256-count-by-path-"
    PREFIX_SUMMARY_IMAGE_LOADED_PATH_SHA256_GROUP_PATH_ORDER_HASH string = "image-loaded-path-sha256-count-by-hash-"
    PREFIX_SUMMARY_IMAGE_LOADED_PATH_MD5_GROUP_PATH_ORDER_PATH string = "image-loaded-path-md5-count-by-path-"
    PREFIX_SUMMARY_IMAGE_LOADED_PATH_MD5_GROUP_PATH_ORDER_HASH string = "image-loaded-path-md5-count-by-hash-"

    PREFIX_SUMMARY_NETWORK_CONNECTION_DISTINCT_DEST_IP string = "network-connection-destination-ip-"
    PREFIX_SUMMARY_NETWORK_CONNECTION_DISTINCT_DEST_IP_COUNT string = "network-connection-destination-ip-count-"

    PREFIX_SUMMARY_NETWORK_CONNECTION_DISTINCT_DEST_HOST string = "network-connection-destination-host-"
    PREFIX_SUMMARY_NETWORK_CONNECTION_DISTINCT_DEST_HOST_COUNT string = "network-connection-destination-host-count-"
)

const (
    SUMMARY_TYPE_PROCESS_CREATE_PATH_SHA256_GROUP_PATH_ORDER_PATH = 1
    SUMMARY_TYPE_PROCESS_CREATE_PATH_SHA256_GROUP_PATH_ORDER_HASH = 2
    SUMMARY_TYPE_PROCESS_CREATE_PATH_MD5_GROUP_PATH_ORDER_PATH = 3
    SUMMARY_TYPE_PROCESS_CREATE_PATH_MD5_GROUP_PATH_ORDER_HASH = 4
    SUMMARY_TYPE_DRIVER_LOADED_PATH_SHA256_GROUP_PATH_ORDER_PATH = 5
    SUMMARY_TYPE_DRIVER_LOADED_PATH_SHA256_GROUP_PATH_ORDER_HASH = 6
    SUMMARY_TYPE_DRIVER_LOADED_PATH_MD5_GROUP_PATH_ORDER_PATH = 7
    SUMMARY_TYPE_DRIVER_LOADED_PATH_MD5_GROUP_PATH_ORDER_HASH = 8

    SUMMARY_TYPE_IMAGE_LOADED_PATH_SHA256_GROUP_PATH_ORDER_PATH = 9
    SUMMARY_TYPE_IMAGE_LOADED_PATH_SHA256_GROUP_PATH_ORDER_HASH = 10
    SUMMARY_TYPE_IMAGE_LOADED_PATH_MD5_GROUP_PATH_ORDER_PATH = 11
    SUMMARY_TYPE_IMAGE_LOADED_PATH_MD5_GROUP_PATH_ORDER_HASH = 12

    SUMMARY_TYPE_NETWORK_CONNECTION_DISTINCT_DEST_IP = 13
    SUMMARY_TYPE_NETWORK_CONNECTION_DISTINCT_DEST_IP_COUNT = 14

    SUMMARY_TYPE_NETWORK_CONNECTION_DISTINCT_DEST_HOST = 15
    SUMMARY_TYPE_NETWORK_CONNECTION_DISTINCT_DEST_HOST_COUNT = 16
)

// ##### Methods #############################################################

// Exports data that needs to be grouped by two text fields, and has
// a total/count field e.g. GROUPED BY PATH and SHA256, with a COUNT
func exportProcessCreateSummaryDataForStringStringTotal(sql string, typeName string, dataType int, prefix string) {

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
                logger.Errorf("Error deleting temporary %s summary file: %v (%s)", typeName, err, tf.Name)
            }
        }
    }()

    var sst StringStringTotal
    for rows.Next() {
        err = rows.Scan(&sst.String1, &sst.String2, &sst.Total)
        if err != nil {
            logger.Errorf("Error scanning struct for %s export: %v", typeName, err)
            return
        }

        tf.WriteString(sst.String1)
        tf.WriteString(",")
        tf.WriteString(sst.String2)
        tf.WriteString(",")
        tf.WriteString(util.ConvertInt64ToString(sst.Total))
        tf.WriteString("\n")
    }

    timestamp := time.Now().UTC()
    fileName := prefix + timestamp.Format(LAYOUT_DAILY_SUMMARY) + ".csv"

    // Move the file
    err = os.Rename(tf.Name(), path.Join(config.SummaryDir, fileName))
    if err != nil {
        logger.Errorf("Error moving file to summary directory: %v (%s)", err, fileName)
        return
    }

    // Insert the summary record
    setSummaryRecord(dataType, fileName)
}

// Exports data that needs to be grouped by one text field, and
// has a total/count field e.g. GROUPED BY PATH, with a COUNT
func exportProcessCreateSummaryDataForStringTotal(sql string, typeName string, dataType int, prefix string) {

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
                logger.Errorf("Error deleting temporary %s summary file: %v (%s)", typeName, err, tf.Name)
            }
        }
    }()

    var st StringTotal
    for rows.Next() {
        err = rows.Scan(&st.String, &st.Total)
        if err != nil {
            logger.Errorf("Error scanning struct for %s export: %v", typeName, err)
            return
        }

        tf.WriteString(st.String)
        tf.WriteString(",")
        tf.WriteString(",")
        tf.WriteString(util.ConvertInt64ToString(st.Total))
        tf.WriteString("\n")
    }

    timestamp := time.Now().UTC()
    fileName := prefix + timestamp.Format(LAYOUT_DAILY_SUMMARY) + ".csv"

    // Move the file
    err = os.Rename(tf.Name(), path.Join(config.SummaryDir, fileName))
    if err != nil {
        logger.Errorf("Error moving file to summary directory: %v (%s)", err, fileName)
        return
    }

    // Insert the summary record
    setSummaryRecord(dataType, fileName)
}

// Exports data that needs a DISTINCT string field summarised
func exportProcessCreateSummaryDataForString(sql string, typeName string, dataType int, prefix string) {

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
                logger.Errorf("Error deleting temporary %s summary file: %v (%s)", typeName, err, tf.Name)
            }
        }
    }()

    var st StringTotal
    for rows.Next() {
        err = rows.Scan(&st.String, &st.Total)
        if err != nil {
            logger.Errorf("Error scanning struct for %s export: %v", typeName, err)
            return
        }

        tf.WriteString(st.String)
        tf.WriteString("\n")
    }

    timestamp := time.Now().UTC()
    fileName := prefix + timestamp.Format(LAYOUT_DAILY_SUMMARY) + ".csv"

    // Move the file
    err = os.Rename(tf.Name(), path.Join(config.SummaryDir, fileName))
    if err != nil {
        logger.Errorf("Error moving file to summary directory: %v (%s)", err, fileName)
        return
    }

    // Insert the summary record
    setSummaryRecord(dataType, fileName)
}

//
func setSummaryRecord(dataType int, fileName string) {

    var s Summary

    err := db.
        Select("*").
        From("summary").
        Where("data_type = $1 and file_name = $2", dataType, fileName).
        OrderBy("id ASC").
        QueryStruct(&s)

    if s.Id > 0 {
        return
    }

    err = db.
        InsertInto("summary").
        Columns("data_type", "file_name").
        Values(dataType, fileName).
        QueryStruct(&s)

    if err != nil {
        if strings.Contains(err.Error(), "no rows in result set") == false {
            logger.Errorf("Error inserting summary record: %v (%s, %s, %s)", err, dataType, fileName)
        }
    }
}