SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;

CREATE EXTENSION "uuid-ossp";

--
-- Name: update_time(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION update_time() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
      BEGIN
      NEW.updated_at := CURRENT_TIMESTAMP;
      RETURN NEW;
      END;
      $$;


SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: users; Type: TABLE; Schema: public; Owner: -; Tablespace: 
--

CREATE TABLE users (
    id integer NOT NULL,
    email character varying(255) NOT NULL,
    password_hash character varying(60) NOT NULL,
    admin boolean DEFAULT false NOT NULL,
    active boolean DEFAULT false NOT NULL,
    verified boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    org_id integer NOT NULL,
    name character varying(255) NOT NULL
);


--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE users_id_seq OWNED BY users.id;


--
-- Name: bastions; Type: TABLE; Schema: public; Owner: -; Tablespace: 
--

CREATE TABLE bastions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v1mc(),
    password_hash character varying(60) NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: orgs; Type: TABLE; Schema: public; Owner: -; Tablespace: 
--

CREATE TABLE orgs (
    id integer NOT NULL,
    name character varying(255),
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: orgs_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE orgs_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: orgs_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE orgs_id_seq OWNED BY orgs.id;


--
-- Name: signups; Type: TABLE; Schema: public; Owner: -; Tablespace: 
--

CREATE TABLE signups (
    id integer NOT NULL,
    email character varying(255),
    name character varying(255),
    claimed boolean default false,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: signups_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE signups_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: signups_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE signups_id_seq OWNED BY signups.id;


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY users ALTER COLUMN id SET DEFAULT nextval('users_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY orgs ALTER COLUMN id SET DEFAULT nextval('orgs_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY signups ALTER COLUMN id SET DEFAULT nextval('signups_id_seq'::regclass);


--
-- Name: users_email_key; Type: CONSTRAINT; Schema: public; Owner: -; Tablespace: 
--

ALTER TABLE ONLY users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: pk_users; Type: CONSTRAINT; Schema: public; Owner: -; Tablespace: 
--

ALTER TABLE ONLY users
    ADD CONSTRAINT pk_users PRIMARY KEY (id);


--
-- Name: pk_orgs; Type: CONSTRAINT; Schema: public; Owner: -; Tablespace: 
--

ALTER TABLE ONLY orgs
    ADD CONSTRAINT pk_orgs PRIMARY KEY (id);


--
-- Name: pk_signups; Type: CONSTRAINT; Schema: public; Owner: -; Tablespace: 
--

ALTER TABLE ONLY signups
    ADD CONSTRAINT pk_signups PRIMARY KEY (id);


--
-- Name: signups_email_key; Type: CONSTRAINT; Schema: public; Owner: -; Tablespace: 
--

ALTER TABLE ONLY signups
    ADD CONSTRAINT signups_email_key UNIQUE (email);


--
-- Name: update_users; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER update_users BEFORE UPDATE ON users FOR EACH ROW EXECUTE PROCEDURE update_time();



--
-- Name: update_bastions; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER update_bastions BEFORE UPDATE ON bastions FOR EACH ROW EXECUTE PROCEDURE update_time();

  

--
-- Name: update_orgs; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER update_orgs BEFORE UPDATE ON orgs FOR EACH ROW EXECUTE PROCEDURE update_time();


--
-- Name: update_signups; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER update_signups BEFORE UPDATE ON signups FOR EACH ROW EXECUTE PROCEDURE update_time();


--
-- Name: fk_users_orgs; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY users
    ADD CONSTRAINT fk_users_orgs FOREIGN KEY (org_id) REFERENCES orgs(id);

