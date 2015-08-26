--
-- PostgreSQL database dump
--

SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;
SET search_path = public, pg_catalog;
SET default_tablespace = '';
SET default_with_oids = false;

CREATE OR REPLACE FUNCTION update_time() RETURNS trigger
LANGUAGE plpgsql
AS $function$
BEGIN
NEW.updated_at := CURRENT_TIMESTAMP;
RETURN NEW;
END;
$function$;

--
-- Name: logins; Type: TABLE; Schema: public; Owner: -; Tablespace: 
--

CREATE TABLE logins (
    id integer NOT NULL,
    email character varying(255) NOT NULL,
    password_hash character varying(60) NOT NULL,
    admin boolean DEFAULT false NOT NULL,
    active boolean DEFAULT false NOT NULL,
    verified boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    onboard boolean DEFAULT true NOT NULL,
    customer_id character varying(64) DEFAULT 'cliff'::character varying,
    name character varying(255)
);


--
-- Name: logins_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE logins_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: logins_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE logins_id_seq OWNED BY logins.id;


--
-- Name: orgs; Type: TABLE; Schema: public; Owner: -; Tablespace: 
--

CREATE TABLE orgs (
    name character varying(255) NOT NULL,
    subdomain character varying(64) NOT NULL
);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY logins ALTER COLUMN id SET DEFAULT nextval('logins_id_seq'::regclass);


--
-- Name: logins_email_key; Type: CONSTRAINT; Schema: public; Owner: -; Tablespace: 
--

ALTER TABLE ONLY logins
    ADD CONSTRAINT logins_email_key UNIQUE (email);


--
-- Name: orgs_subdomain_key; Type: CONSTRAINT; Schema: public; Owner: -; Tablespace: 
--

ALTER TABLE ONLY orgs
    ADD CONSTRAINT orgs_subdomain_key UNIQUE (subdomain);


--
-- Name: pk_logins; Type: CONSTRAINT; Schema: public; Owner: -; Tablespace: 
--

ALTER TABLE ONLY logins
    ADD CONSTRAINT pk_logins PRIMARY KEY (id);


--
-- Name: update_logins; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER update_logins BEFORE UPDATE ON logins FOR EACH ROW EXECUTE PROCEDURE update_time();


--
-- Name: fk_logins_orgs; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY logins
    ADD CONSTRAINT fk_logins_orgs FOREIGN KEY (customer_id) REFERENCES orgs(subdomain);


--
-- PostgreSQL database dump complete
--

