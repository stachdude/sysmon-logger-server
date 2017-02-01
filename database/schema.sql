--
-- PostgreSQL database dump
--

-- Dumped from database version 9.5.5
-- Dumped by pg_dump version 9.5.5

-- Started on 2017-02-01 07:55:11 GMT

SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

--
-- TOC entry 1 (class 3079 OID 12393)
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- TOC entry 2283 (class 0 OID 0)
-- Dependencies: 1
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET search_path = public, pg_catalog;

SET default_with_oids = false;

--
-- TOC entry 186 (class 1259 OID 32226)
-- Name: create_remote_thread; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE create_remote_thread (
    id bigint NOT NULL,
    domain text,
    host text,
    event_log_time timestamp without time zone,
    utc_time timestamp without time zone,
    source_process_id bigint,
    source_image text,
    target_process_id bigint,
    target_image text,
    new_thread_id bigint,
    start_address text,
    start_module text,
    start_function text
);


--
-- TOC entry 185 (class 1259 OID 32224)
-- Name: create_remote_thread_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE create_remote_thread_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 2284 (class 0 OID 0)
-- Dependencies: 185
-- Name: create_remote_thread_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE create_remote_thread_id_seq OWNED BY create_remote_thread.id;


--
-- TOC entry 188 (class 1259 OID 32237)
-- Name: driver_loaded; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE driver_loaded (
    id bigint NOT NULL,
    domain text,
    host text,
    event_log_time timestamp without time zone,
    utc_time timestamp without time zone,
    image_loaded text,
    md5 text,
    sha256 text,
    signed boolean,
    signature text
);


--
-- TOC entry 187 (class 1259 OID 32235)
-- Name: driver_loaded_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE driver_loaded_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 2285 (class 0 OID 0)
-- Dependencies: 187
-- Name: driver_loaded_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE driver_loaded_id_seq OWNED BY driver_loaded.id;


--
-- TOC entry 181 (class 1259 OID 22070)
-- Name: event; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE event (
    id bigint NOT NULL,
    domain text,
    host text,
    utc_time timestamp without time zone,
    type text,
    message text,
    message_html text
);


--
-- TOC entry 182 (class 1259 OID 22076)
-- Name: event_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE event_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 2286 (class 0 OID 0)
-- Dependencies: 182
-- Name: event_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE event_id_seq OWNED BY event.id;


--
-- TOC entry 183 (class 1259 OID 22078)
-- Name: export; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE export (
    id bigint NOT NULL,
    data_type integer,
    file_name text,
    updated timestamp without time zone
);


--
-- TOC entry 184 (class 1259 OID 22132)
-- Name: export_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE export_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 2287 (class 0 OID 0)
-- Dependencies: 184
-- Name: export_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE export_id_seq OWNED BY export.id;


--
-- TOC entry 210 (class 1259 OID 33327)
-- Name: file_create; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE file_create (
    id bigint NOT NULL,
    domain text,
    host text,
    event_log_time timestamp without time zone,
    utc_time timestamp without time zone,
    process_id bigint,
    image text,
    target_file_name text,
    creation_utc_time timestamp without time zone
);


--
-- TOC entry 209 (class 1259 OID 33325)
-- Name: file_create_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE file_create_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 2288 (class 0 OID 0)
-- Dependencies: 209
-- Name: file_create_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE file_create_id_seq OWNED BY file_create.id;


--
-- TOC entry 190 (class 1259 OID 32250)
-- Name: file_creation_time; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE file_creation_time (
    id bigint NOT NULL,
    domain text,
    host text,
    event_log_time timestamp without time zone,
    utc_time timestamp without time zone,
    process_id bigint,
    image text,
    target_file_name text,
    creation_utc_time timestamp without time zone,
    previous_creation_utc_time timestamp without time zone
);


--
-- TOC entry 189 (class 1259 OID 32248)
-- Name: file_creation_time_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE file_creation_time_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 2289 (class 0 OID 0)
-- Dependencies: 189
-- Name: file_creation_time_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE file_creation_time_id_seq OWNED BY file_creation_time.id;


--
-- TOC entry 192 (class 1259 OID 32295)
-- Name: image_loaded; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE image_loaded (
    id bigint NOT NULL,
    domain text,
    host text,
    event_log_time timestamp without time zone,
    utc_time timestamp without time zone,
    process_id bigint,
    image text,
    image_loaded text,
    md5 text,
    sha256 text,
    signed boolean,
    signature text
);


--
-- TOC entry 191 (class 1259 OID 32293)
-- Name: image_loaded_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE image_loaded_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 2290 (class 0 OID 0)
-- Dependencies: 191
-- Name: image_loaded_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE image_loaded_id_seq OWNED BY image_loaded.id;


--
-- TOC entry 194 (class 1259 OID 32362)
-- Name: network_connection; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE network_connection (
    id bigint NOT NULL,
    domain text,
    host text,
    event_log_time timestamp without time zone,
    utc_time timestamp without time zone,
    process_id bigint,
    image text,
    process_user text,
    protocol text,
    initiated boolean,
    source_ip inet,
    source_host_name text,
    source_port integer,
    source_port_name text,
    destination_ip inet,
    destination_host_name text,
    destination_port integer,
    destination_port_name text
);


--
-- TOC entry 193 (class 1259 OID 32360)
-- Name: network_connection_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE network_connection_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 2291 (class 0 OID 0)
-- Dependencies: 193
-- Name: network_connection_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE network_connection_id_seq OWNED BY network_connection.id;


--
-- TOC entry 208 (class 1259 OID 33316)
-- Name: process_access; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE process_access (
    id bigint NOT NULL,
    domain text,
    host text,
    event_log_time timestamp without time zone,
    utc_time timestamp without time zone,
    source_process_id bigint,
    source_image text,
    target_process_id bigint,
    target_image text,
    call_trace text,
    granted_access text
);


--
-- TOC entry 207 (class 1259 OID 33314)
-- Name: process_access_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE process_access_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 2292 (class 0 OID 0)
-- Dependencies: 207
-- Name: process_access_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE process_access_id_seq OWNED BY process_access.id;


--
-- TOC entry 196 (class 1259 OID 32373)
-- Name: process_create; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE process_create (
    id bigint NOT NULL,
    domain text,
    host text,
    event_log_time timestamp without time zone,
    utc_time timestamp without time zone,
    process_id bigint,
    image text,
    command_line text,
    current_directory text,
    md5 text,
    sha256 text,
    parent_process_id bigint,
    parent_image text,
    parent_command_line text,
    process_user text
);


--
-- TOC entry 195 (class 1259 OID 32371)
-- Name: process_create_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE process_create_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 2293 (class 0 OID 0)
-- Dependencies: 195
-- Name: process_create_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE process_create_id_seq OWNED BY process_create.id;


--
-- TOC entry 198 (class 1259 OID 32387)
-- Name: process_terminate; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE process_terminate (
    id bigint NOT NULL,
    domain text,
    host text,
    event_log_time timestamp without time zone,
    utc_time timestamp without time zone,
    process_id bigint,
    image text
);


--
-- TOC entry 197 (class 1259 OID 32385)
-- Name: process_terminate_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE process_terminate_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 2294 (class 0 OID 0)
-- Dependencies: 197
-- Name: process_terminate_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE process_terminate_id_seq OWNED BY process_terminate.id;


--
-- TOC entry 200 (class 1259 OID 32400)
-- Name: raw_access_read; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE raw_access_read (
    id bigint NOT NULL,
    domain text,
    host text,
    event_log_time timestamp without time zone,
    utc_time timestamp without time zone,
    process_id bigint,
    image text,
    device text
);


--
-- TOC entry 199 (class 1259 OID 32398)
-- Name: raw_access_read_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE raw_access_read_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 2295 (class 0 OID 0)
-- Dependencies: 199
-- Name: raw_access_read_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE raw_access_read_id_seq OWNED BY raw_access_read.id;


--
-- TOC entry 202 (class 1259 OID 32411)
-- Name: registry_add_delete; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE registry_add_delete (
    id bigint NOT NULL,
    domain text,
    host text,
    event_log_time timestamp without time zone,
    utc_time timestamp without time zone,
    process_id bigint,
    image text,
    event_type text,
    target_object text
);


--
-- TOC entry 201 (class 1259 OID 32409)
-- Name: registry_add_delete_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE registry_add_delete_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 2296 (class 0 OID 0)
-- Dependencies: 201
-- Name: registry_add_delete_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE registry_add_delete_id_seq OWNED BY registry_add_delete.id;


--
-- TOC entry 204 (class 1259 OID 32422)
-- Name: registry_rename; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE registry_rename (
    id bigint NOT NULL,
    domain text,
    host text,
    event_log_time timestamp without time zone,
    utc_time timestamp without time zone,
    process_id bigint,
    image text,
    event_type text,
    target_object text,
    new_name text
);


--
-- TOC entry 203 (class 1259 OID 32420)
-- Name: registry_rename_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE registry_rename_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 2297 (class 0 OID 0)
-- Dependencies: 203
-- Name: registry_rename_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE registry_rename_id_seq OWNED BY registry_rename.id;


--
-- TOC entry 206 (class 1259 OID 32433)
-- Name: registry_set; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE registry_set (
    id bigint NOT NULL,
    domain text,
    host text,
    event_log_time timestamp without time zone,
    utc_time timestamp without time zone,
    process_id bigint,
    image text,
    event_type text,
    target_object text,
    details text
);


--
-- TOC entry 205 (class 1259 OID 32431)
-- Name: registry_set_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE registry_set_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 2298 (class 0 OID 0)
-- Dependencies: 205
-- Name: registry_set_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE registry_set_id_seq OWNED BY registry_set.id;


--
-- TOC entry 2118 (class 2604 OID 32229)
-- Name: id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY create_remote_thread ALTER COLUMN id SET DEFAULT nextval('create_remote_thread_id_seq'::regclass);


--
-- TOC entry 2119 (class 2604 OID 32240)
-- Name: id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY driver_loaded ALTER COLUMN id SET DEFAULT nextval('driver_loaded_id_seq'::regclass);


--
-- TOC entry 2116 (class 2604 OID 22136)
-- Name: id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY event ALTER COLUMN id SET DEFAULT nextval('event_id_seq'::regclass);


--
-- TOC entry 2117 (class 2604 OID 22137)
-- Name: id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY export ALTER COLUMN id SET DEFAULT nextval('export_id_seq'::regclass);


--
-- TOC entry 2130 (class 2604 OID 33330)
-- Name: id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY file_create ALTER COLUMN id SET DEFAULT nextval('file_create_id_seq'::regclass);


--
-- TOC entry 2120 (class 2604 OID 32253)
-- Name: id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY file_creation_time ALTER COLUMN id SET DEFAULT nextval('file_creation_time_id_seq'::regclass);


--
-- TOC entry 2121 (class 2604 OID 32298)
-- Name: id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY image_loaded ALTER COLUMN id SET DEFAULT nextval('image_loaded_id_seq'::regclass);


--
-- TOC entry 2122 (class 2604 OID 32365)
-- Name: id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY network_connection ALTER COLUMN id SET DEFAULT nextval('network_connection_id_seq'::regclass);


--
-- TOC entry 2129 (class 2604 OID 33319)
-- Name: id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY process_access ALTER COLUMN id SET DEFAULT nextval('process_access_id_seq'::regclass);


--
-- TOC entry 2123 (class 2604 OID 32376)
-- Name: id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY process_create ALTER COLUMN id SET DEFAULT nextval('process_create_id_seq'::regclass);


--
-- TOC entry 2124 (class 2604 OID 32390)
-- Name: id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY process_terminate ALTER COLUMN id SET DEFAULT nextval('process_terminate_id_seq'::regclass);


--
-- TOC entry 2125 (class 2604 OID 32403)
-- Name: id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY raw_access_read ALTER COLUMN id SET DEFAULT nextval('raw_access_read_id_seq'::regclass);


--
-- TOC entry 2126 (class 2604 OID 32414)
-- Name: id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY registry_add_delete ALTER COLUMN id SET DEFAULT nextval('registry_add_delete_id_seq'::regclass);


--
-- TOC entry 2127 (class 2604 OID 32425)
-- Name: id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY registry_rename ALTER COLUMN id SET DEFAULT nextval('registry_rename_id_seq'::regclass);


--
-- TOC entry 2128 (class 2604 OID 32436)
-- Name: id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY registry_set ALTER COLUMN id SET DEFAULT nextval('registry_set_id_seq'::regclass);


--
-- TOC entry 2138 (class 2606 OID 32234)
-- Name: create_remote_thread_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY create_remote_thread
    ADD CONSTRAINT create_remote_thread_pkey PRIMARY KEY (id);


--
-- TOC entry 2140 (class 2606 OID 32245)
-- Name: driver_loaded_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY driver_loaded
    ADD CONSTRAINT driver_loaded_pkey PRIMARY KEY (id);


--
-- TOC entry 2132 (class 2606 OID 22149)
-- Name: event_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY event
    ADD CONSTRAINT event_pkey PRIMARY KEY (id);


--
-- TOC entry 2134 (class 2606 OID 22151)
-- Name: export_file_name_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY export
    ADD CONSTRAINT export_file_name_key UNIQUE (file_name);


--
-- TOC entry 2136 (class 2606 OID 22153)
-- Name: export_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY export
    ADD CONSTRAINT export_pkey PRIMARY KEY (id);


--
-- TOC entry 2162 (class 2606 OID 33335)
-- Name: file_create_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY file_create
    ADD CONSTRAINT file_create_pkey PRIMARY KEY (id);


--
-- TOC entry 2142 (class 2606 OID 32258)
-- Name: file_creation_time_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY file_creation_time
    ADD CONSTRAINT file_creation_time_pkey PRIMARY KEY (id);


--
-- TOC entry 2144 (class 2606 OID 32303)
-- Name: image_loaded_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY image_loaded
    ADD CONSTRAINT image_loaded_pkey PRIMARY KEY (id);


--
-- TOC entry 2146 (class 2606 OID 32370)
-- Name: network_connection_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY network_connection
    ADD CONSTRAINT network_connection_pkey PRIMARY KEY (id);


--
-- TOC entry 2160 (class 2606 OID 33324)
-- Name: process_access_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY process_access
    ADD CONSTRAINT process_access_pkey PRIMARY KEY (id);


--
-- TOC entry 2148 (class 2606 OID 32381)
-- Name: process_create_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY process_create
    ADD CONSTRAINT process_create_pkey PRIMARY KEY (id);


--
-- TOC entry 2150 (class 2606 OID 32395)
-- Name: process_terminate_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY process_terminate
    ADD CONSTRAINT process_terminate_pkey PRIMARY KEY (id);


--
-- TOC entry 2152 (class 2606 OID 32408)
-- Name: raw_access_read_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY raw_access_read
    ADD CONSTRAINT raw_access_read_pkey PRIMARY KEY (id);


--
-- TOC entry 2154 (class 2606 OID 32419)
-- Name: registry_add_delete_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY registry_add_delete
    ADD CONSTRAINT registry_add_delete_pkey PRIMARY KEY (id);


--
-- TOC entry 2156 (class 2606 OID 32430)
-- Name: registry_renamed_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY registry_rename
    ADD CONSTRAINT registry_renamed_pkey PRIMARY KEY (id);


--
-- TOC entry 2158 (class 2606 OID 32441)
-- Name: registry_set_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY registry_set
    ADD CONSTRAINT registry_set_pkey PRIMARY KEY (id);


-- Completed on 2017-02-01 07:55:11 GMT

--
-- PostgreSQL database dump complete
--

