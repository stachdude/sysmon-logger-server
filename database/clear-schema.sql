TRUNCATE TABLE create_remote_thread;
ALTER SEQUENCE create_remote_thread_id_seq RESTART WITH 1;

TRUNCATE TABLE driver_loaded;
ALTER SEQUENCE driver_loaded_id_seq RESTART WITH 1;

TRUNCATE TABLE file_creation_time;
ALTER SEQUENCE file_creation_time_id_seq RESTART WITH 1;

TRUNCATE TABLE image_loaded;
ALTER SEQUENCE image_loaded_id_seq RESTART WITH 1;

TRUNCATE TABLE network_connection;
ALTER SEQUENCE network_connection_id_seq RESTART WITH 1;

TRUNCATE TABLE process_create;
ALTER SEQUENCE process_create_id_seq RESTART WITH 1;

TRUNCATE TABLE process_terminate;
ALTER SEQUENCE process_terminate_id_seq RESTART WITH 1;

TRUNCATE TABLE raw_access_read;
ALTER SEQUENCE raw_access_read_id_seq RESTART WITH 1;

TRUNCATE TABLE event;
ALTER SEQUENCE event_id_seq RESTART WITH 1;
