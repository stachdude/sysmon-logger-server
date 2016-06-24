--
-- PostgreSQL database dump
--

-- Dumped from database version 9.5.3
-- Dumped by pg_dump version 9.5.3

-- Started on 2016-06-24 15:04:58 BST

SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

DROP DATABASE sml;
--
-- TOC entry 2231 (class 1262 OID 18031)
-- Name: sml; Type: DATABASE; Schema: -; Owner: -
--

CREATE DATABASE sml WITH TEMPLATE = template0 ENCODING = 'UTF8' LC_COLLATE = 'en_GB.UTF-8' LC_CTYPE = 'en_GB.UTF-8';


SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

--
-- TOC entry 6 (class 2615 OID 2200)
-- Name: public; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA public;


--
-- TOC entry 2232 (class 0 OID 0)
-- Dependencies: 6
-- Name: SCHEMA public; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON SCHEMA public IS 'standard public schema';


--
-- TOC entry 1 (class 3079 OID 12393)
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- TOC entry 2233 (class 0 OID 0)
-- Dependencies: 1
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET search_path = public, pg_catalog;

SET default_with_oids = false;

--
-- TOC entry 185 (class 1259 OID 18161)
-- Name: create_remote_thread; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE create_remote_thread (
    id bigint NOT NULL,
    domain text,
    host text,
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
-- TOC entry 186 (class 1259 OID 18167)
-- Name: create_remote_thread_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE create_remote_thread_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 2234 (class 0 OID 0)
-- Dependencies: 186
-- Name: create_remote_thread_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE create_remote_thread_id_seq OWNED BY create_remote_thread.id;


--
-- TOC entry 187 (class 1259 OID 18169)
-- Name: driver_loaded; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE driver_loaded (
    id bigint NOT NULL,
    domain text,
    host text,
    utc_time timestamp without time zone,
    image_loaded text,
    md5 text,
    sha256 text,
    signed boolean,
    signature text
);


--
-- TOC entry 188 (class 1259 OID 18175)
-- Name: driver_loaded_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE driver_loaded_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 2235 (class 0 OID 0)
-- Dependencies: 188
-- Name: driver_loaded_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE driver_loaded_id_seq OWNED BY driver_loaded.id;


--
-- TOC entry 181 (class 1259 OID 18096)
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
-- TOC entry 182 (class 1259 OID 18102)
-- Name: event_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE event_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 2236 (class 0 OID 0)
-- Dependencies: 182
-- Name: event_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE event_id_seq OWNED BY event.id;


--
-- TOC entry 184 (class 1259 OID 18149)
-- Name: export; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE export (
    id bigint NOT NULL,
    data_type integer,
    file_name text,
    updated timestamp without time zone
);


--
-- TOC entry 189 (class 1259 OID 18285)
-- Name: file_creation_time; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE file_creation_time (
    id bigint NOT NULL,
    domain text,
    host text,
    utc_time timestamp without time zone,
    process_id bigint,
    image text,
    target_file_name text,
    creation_utc_time timestamp without time zone,
    previous_creation_utc_time timestamp without time zone
);


--
-- TOC entry 190 (class 1259 OID 18291)
-- Name: file_creation_time_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE file_creation_time_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 2237 (class 0 OID 0)
-- Dependencies: 190
-- Name: file_creation_time_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE file_creation_time_id_seq OWNED BY file_creation_time.id;


--
-- TOC entry 191 (class 1259 OID 18293)
-- Name: image_loaded; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE image_loaded (
    id bigint NOT NULL,
    domain text,
    host text,
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
-- TOC entry 192 (class 1259 OID 18299)
-- Name: image_loaded_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE image_loaded_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 2238 (class 0 OID 0)
-- Dependencies: 192
-- Name: image_loaded_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE image_loaded_id_seq OWNED BY image_loaded.id;


--
-- TOC entry 193 (class 1259 OID 18301)
-- Name: network_connection; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE network_connection (
    id bigint NOT NULL,
    domain text,
    host text,
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
-- TOC entry 194 (class 1259 OID 18307)
-- Name: network_connection_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE network_connection_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 2239 (class 0 OID 0)
-- Dependencies: 194
-- Name: network_connection_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE network_connection_id_seq OWNED BY network_connection.id;


--
-- TOC entry 195 (class 1259 OID 18309)
-- Name: process_create; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE process_create (
    id bigint NOT NULL,
    domain text,
    host text,
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
-- TOC entry 196 (class 1259 OID 18315)
-- Name: process_create_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE process_create_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 2240 (class 0 OID 0)
-- Dependencies: 196
-- Name: process_create_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE process_create_id_seq OWNED BY process_create.id;


--
-- TOC entry 197 (class 1259 OID 18317)
-- Name: process_terminate; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE process_terminate (
    id bigint NOT NULL,
    domain text,
    host text,
    utc_time timestamp without time zone,
    process_id bigint,
    image text
);


--
-- TOC entry 198 (class 1259 OID 18323)
-- Name: process_terminate_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE process_terminate_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 2241 (class 0 OID 0)
-- Dependencies: 198
-- Name: process_terminate_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE process_terminate_id_seq OWNED BY process_terminate.id;


--
-- TOC entry 199 (class 1259 OID 18325)
-- Name: raw_access_read; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE raw_access_read (
    id bigint NOT NULL,
    domain text,
    host text,
    utc_time timestamp without time zone,
    process_id bigint,
    image text,
    device text
);


--
-- TOC entry 200 (class 1259 OID 18331)
-- Name: raw_access_read_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE raw_access_read_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 2242 (class 0 OID 0)
-- Dependencies: 200
-- Name: raw_access_read_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE raw_access_read_id_seq OWNED BY raw_access_read.id;


--
-- TOC entry 183 (class 1259 OID 18147)
-- Name: summary_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE summary_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 2243 (class 0 OID 0)
-- Dependencies: 183
-- Name: summary_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE summary_id_seq OWNED BY export.id;


--
-- TOC entry 2083 (class 2604 OID 18358)
-- Name: id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY create_remote_thread ALTER COLUMN id SET DEFAULT nextval('create_remote_thread_id_seq'::regclass);


--
-- TOC entry 2084 (class 2604 OID 18359)
-- Name: id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY driver_loaded ALTER COLUMN id SET DEFAULT nextval('driver_loaded_id_seq'::regclass);


--
-- TOC entry 2081 (class 2604 OID 18737)
-- Name: id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY event ALTER COLUMN id SET DEFAULT nextval('event_id_seq'::regclass);


--
-- TOC entry 2082 (class 2604 OID 18152)
-- Name: id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY export ALTER COLUMN id SET DEFAULT nextval('summary_id_seq'::regclass);


--
-- TOC entry 2085 (class 2604 OID 18360)
-- Name: id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY file_creation_time ALTER COLUMN id SET DEFAULT nextval('file_creation_time_id_seq'::regclass);


--
-- TOC entry 2086 (class 2604 OID 18361)
-- Name: id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY image_loaded ALTER COLUMN id SET DEFAULT nextval('image_loaded_id_seq'::regclass);


--
-- TOC entry 2087 (class 2604 OID 18362)
-- Name: id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY network_connection ALTER COLUMN id SET DEFAULT nextval('network_connection_id_seq'::regclass);


--
-- TOC entry 2088 (class 2604 OID 18363)
-- Name: id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY process_create ALTER COLUMN id SET DEFAULT nextval('process_create_id_seq'::regclass);


--
-- TOC entry 2089 (class 2604 OID 18364)
-- Name: id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY process_terminate ALTER COLUMN id SET DEFAULT nextval('process_terminate_id_seq'::regclass);


--
-- TOC entry 2090 (class 2604 OID 18365)
-- Name: id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY raw_access_read ALTER COLUMN id SET DEFAULT nextval('raw_access_read_id_seq'::regclass);


--
-- TOC entry 2098 (class 2606 OID 18357)
-- Name: create_remote_thread_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY create_remote_thread
    ADD CONSTRAINT create_remote_thread_pkey PRIMARY KEY (id);


--
-- TOC entry 2100 (class 2606 OID 18355)
-- Name: driver_loaded_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY driver_loaded
    ADD CONSTRAINT driver_loaded_pkey PRIMARY KEY (id);


--
-- TOC entry 2092 (class 2606 OID 18130)
-- Name: event_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY event
    ADD CONSTRAINT event_pkey PRIMARY KEY (id);


--
-- TOC entry 2094 (class 2606 OID 18159)
-- Name: export_file_name_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY export
    ADD CONSTRAINT export_file_name_key UNIQUE (file_name);


--
-- TOC entry 2096 (class 2606 OID 18157)
-- Name: export_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY export
    ADD CONSTRAINT export_pkey PRIMARY KEY (id);


--
-- TOC entry 2102 (class 2606 OID 18353)
-- Name: file_creation_time_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY file_creation_time
    ADD CONSTRAINT file_creation_time_pkey PRIMARY KEY (id);


--
-- TOC entry 2104 (class 2606 OID 18351)
-- Name: image_loaded_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY image_loaded
    ADD CONSTRAINT image_loaded_pkey PRIMARY KEY (id);


--
-- TOC entry 2106 (class 2606 OID 18349)
-- Name: network_connection_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY network_connection
    ADD CONSTRAINT network_connection_pkey PRIMARY KEY (id);


--
-- TOC entry 2108 (class 2606 OID 18347)
-- Name: process_create_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY process_create
    ADD CONSTRAINT process_create_pkey PRIMARY KEY (id);


--
-- TOC entry 2110 (class 2606 OID 18345)
-- Name: process_terminate_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY process_terminate
    ADD CONSTRAINT process_terminate_pkey PRIMARY KEY (id);


--
-- TOC entry 2112 (class 2606 OID 18343)
-- Name: raw_access_read_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY raw_access_read
    ADD CONSTRAINT raw_access_read_pkey PRIMARY KEY (id);


-- Completed on 2016-06-24 15:04:58 BST

--
-- PostgreSQL database dump complete
--

