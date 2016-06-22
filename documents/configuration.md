# Configuration

There are various configuration options for the server application. This document details the configuration options and the required values:

- database_server: The host name or IP address of the PostgreSQL database server
- database_name: Name of the database (sml)
- database_user:  Database server user name (postgres)
- database_password: Password for the database user
- http_ip: IP address of the interface for the HTTPS API. Use 0.0.0.0 to access on all interfaces
- http_port: Port for the HTTPS server. Use the **sml-setbind.sh** file to allow lower port access such as port 80
- debug: Show each HTTPS request in the logs (true/false)
- processor_threads. The number of processor threads. Use the value 0 to auto configure
- server_pem: Full path to the server PEM file (server.pem)
- server_key: Full path to the server key file (server.key)
