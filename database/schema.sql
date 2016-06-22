--
-- PostgreSQL database dump
--

-- Dumped from database version 9.5.3
-- Dumped by pg_dump version 9.5.3

-- Started on 2016-06-22 10:38:53 BST

SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

DROP DATABASE sml;
--
-- TOC entry 2219 (class 1262 OID 18031)
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
-- TOC entry 2220 (class 0 OID 0)
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
-- TOC entry 2222 (class 0 OID 0)
-- Dependencies: 1
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET search_path = public, pg_catalog;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- TOC entry 181 (class 1259 OID 18032)
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
-- TOC entry 182 (class 1259 OID 18038)
-- Name: create_remote_thread_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE create_remote_thread_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 2223 (class 0 OID 0)
-- Dependencies: 182
-- Name: create_remote_thread_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE create_remote_thread_id_seq OWNED BY create_remote_thread.id;


--
-- TOC entry 183 (class 1259 OID 18040)
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
-- TOC entry 184 (class 1259 OID 18046)
-- Name: driver_loaded_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE driver_loaded_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 2224 (class 0 OID 0)
-- Dependencies: 184
-- Name: driver_loaded_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE driver_loaded_id_seq OWNED BY driver_loaded.id;


--
-- TOC entry 185 (class 1259 OID 18048)
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
-- TOC entry 186 (class 1259 OID 18054)
-- Name: file_creation_time_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE file_creation_time_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 2225 (class 0 OID 0)
-- Dependencies: 186
-- Name: file_creation_time_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE file_creation_time_id_seq OWNED BY file_creation_time.id;


--
-- TOC entry 187 (class 1259 OID 18056)
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
-- TOC entry 188 (class 1259 OID 18062)
-- Name: image_loaded_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE image_loaded_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 2226 (class 0 OID 0)
-- Dependencies: 188
-- Name: image_loaded_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE image_loaded_id_seq OWNED BY image_loaded.id;


--
-- TOC entry 189 (class 1259 OID 18064)
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
-- TOC entry 190 (class 1259 OID 18070)
-- Name: network_connection_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE network_connection_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 2227 (class 0 OID 0)
-- Dependencies: 190
-- Name: network_connection_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE network_connection_id_seq OWNED BY network_connection.id;


--
-- TOC entry 191 (class 1259 OID 18072)
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
-- TOC entry 192 (class 1259 OID 18078)
-- Name: process_create_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE process_create_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 2228 (class 0 OID 0)
-- Dependencies: 192
-- Name: process_create_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE process_create_id_seq OWNED BY process_create.id;


--
-- TOC entry 193 (class 1259 OID 18080)
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
-- TOC entry 194 (class 1259 OID 18086)
-- Name: process_terminated_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE process_terminated_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 2229 (class 0 OID 0)
-- Dependencies: 194
-- Name: process_terminated_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE process_terminated_id_seq OWNED BY process_terminate.id;


--
-- TOC entry 195 (class 1259 OID 18088)
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
-- TOC entry 196 (class 1259 OID 18094)
-- Name: raw_access_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE raw_access_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 2230 (class 0 OID 0)
-- Dependencies: 196
-- Name: raw_access_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE raw_access_id_seq OWNED BY raw_access_read.id;


--
-- TOC entry 197 (class 1259 OID 18096)
-- Name: unified; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE unified (
    id bigint NOT NULL,
    domain text,
    host text,
    utc_time timestamp without time zone,
    type text,
    message text
);


--
-- TOC entry 198 (class 1259 OID 18102)
-- Name: unified_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE unified_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 2231 (class 0 OID 0)
-- Dependencies: 198
-- Name: unified_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE unified_id_seq OWNED BY unified.id;


--
-- TOC entry 2074 (class 2604 OID 18104)
-- Name: id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY create_remote_thread ALTER COLUMN id SET DEFAULT nextval('create_remote_thread_id_seq'::regclass);


--
-- TOC entry 2075 (class 2604 OID 18105)
-- Name: id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY driver_loaded ALTER COLUMN id SET DEFAULT nextval('driver_loaded_id_seq'::regclass);


--
-- TOC entry 2076 (class 2604 OID 18106)
-- Name: id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY file_creation_time ALTER COLUMN id SET DEFAULT nextval('file_creation_time_id_seq'::regclass);


--
-- TOC entry 2077 (class 2604 OID 18107)
-- Name: id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY image_loaded ALTER COLUMN id SET DEFAULT nextval('image_loaded_id_seq'::regclass);


--
-- TOC entry 2078 (class 2604 OID 18108)
-- Name: id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY network_connection ALTER COLUMN id SET DEFAULT nextval('network_connection_id_seq'::regclass);


--
-- TOC entry 2079 (class 2604 OID 18109)
-- Name: id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY process_create ALTER COLUMN id SET DEFAULT nextval('process_create_id_seq'::regclass);


--
-- TOC entry 2080 (class 2604 OID 18110)
-- Name: id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY process_terminate ALTER COLUMN id SET DEFAULT nextval('process_terminated_id_seq'::regclass);


--
-- TOC entry 2081 (class 2604 OID 18111)
-- Name: id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY raw_access_read ALTER COLUMN id SET DEFAULT nextval('raw_access_id_seq'::regclass);


--
-- TOC entry 2082 (class 2604 OID 18112)
-- Name: id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY unified ALTER COLUMN id SET DEFAULT nextval('unified_id_seq'::regclass);


--
-- TOC entry 2084 (class 2606 OID 18114)
-- Name: create_remote_thread_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY create_remote_thread
    ADD CONSTRAINT create_remote_thread_pkey PRIMARY KEY (id);


--
-- TOC entry 2086 (class 2606 OID 18116)
-- Name: driver_loaded_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY driver_loaded
    ADD CONSTRAINT driver_loaded_pkey PRIMARY KEY (id);


--
-- TOC entry 2088 (class 2606 OID 18118)
-- Name: file_creation_time_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY file_creation_time
    ADD CONSTRAINT file_creation_time_pkey PRIMARY KEY (id);


--
-- TOC entry 2090 (class 2606 OID 18120)
-- Name: image_loaded_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY image_loaded
    ADD CONSTRAINT image_loaded_pkey PRIMARY KEY (id);


--
-- TOC entry 2092 (class 2606 OID 18122)
-- Name: network_connection_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY network_connection
    ADD CONSTRAINT network_connection_pkey PRIMARY KEY (id);


--
-- TOC entry 2094 (class 2606 OID 18124)
-- Name: process_create_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY process_create
    ADD CONSTRAINT process_create_pkey PRIMARY KEY (id);


--
-- TOC entry 2096 (class 2606 OID 18126)
-- Name: process_terminate_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY process_terminate
    ADD CONSTRAINT process_terminate_pkey PRIMARY KEY (id);


--
-- TOC entry 2098 (class 2606 OID 18128)
-- Name: raw_access_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY raw_access_read
    ADD CONSTRAINT raw_access_pkey PRIMARY KEY (id);


--
-- TOC entry 2100 (class 2606 OID 18130)
-- Name: unified_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY unified
    ADD CONSTRAINT unified_pkey PRIMARY KEY (id);


--
-- TOC entry 2221 (class 0 OID 0)
-- Dependencies: 6
-- Name: public; Type: ACL; Schema: -; Owner: -
--

REVOKE ALL ON SCHEMA public FROM PUBLIC;
REVOKE ALL ON SCHEMA public FROM postgres;
GRANT ALL ON SCHEMA public TO postgres;
GRANT ALL ON SCHEMA public TO PUBLIC;


-- Completed on 2016-06-22 10:38:53 BST

--
-- PostgreSQL database dump complete
--

