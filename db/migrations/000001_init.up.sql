--
-- PostgreSQL database dump
--

-- Dumped from database version 14.11 (Ubuntu 14.11-0ubuntu0.22.04.1)
-- Dumped by pg_dump version 14.11 (Ubuntu 14.11-0ubuntu0.22.04.1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: expense; Type: TABLE; Schema: public; Owner: expense_user
--

CREATE TABLE public.expense (
    id integer NOT NULL,
    expense_type_id integer NOT NULL,
    reated_at timestamp without time zone DEFAULT now() NOT NULL,
    spent_money numeric(16,2) NOT NULL
);


ALTER TABLE public.expense OWNER TO expense_user;

--
-- Name: expense_id_seq; Type: SEQUENCE; Schema: public; Owner: expense_user
--

CREATE SEQUENCE public.expense_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.expense_id_seq OWNER TO expense_user;

--
-- Name: expense_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: expense_user
--

ALTER SEQUENCE public.expense_id_seq OWNED BY public.expense.id;


--
-- Name: expense_type; Type: TABLE; Schema: public; Owner: expense_user
--

CREATE TABLE public.expense_type (
    id integer NOT NULL,
    users_id integer NOT NULL,
    title character varying(50) NOT NULL
);


ALTER TABLE public.expense_type OWNER TO expense_user;

--
-- Name: expense_type_id_seq; Type: SEQUENCE; Schema: public; Owner: expense_user
--

CREATE SEQUENCE public.expense_type_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.expense_type_id_seq OWNER TO expense_user;

--
-- Name: expense_type_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: expense_user
--

ALTER SEQUENCE public.expense_type_id_seq OWNED BY public.expense_type.id;


--
-- Name: files; Type: TABLE; Schema: public; Owner: expense_user
--

CREATE TABLE public.files (
    id integer NOT NULL,
    expense_id integer NOT NULL,
    path_file character varying(300) NOT NULL,
    type_file character varying(50) NOT NULL
);


ALTER TABLE public.files OWNER TO expense_user;

--
-- Name: files_id_seq; Type: SEQUENCE; Schema: public; Owner: expense_user
--

CREATE SEQUENCE public.files_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.files_id_seq OWNER TO expense_user;

--
-- Name: files_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: expense_user
--

ALTER SEQUENCE public.files_id_seq OWNED BY public.files.id;


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: expense_user
--

CREATE TABLE IF NOT EXISTS public.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);


ALTER TABLE public.schema_migrations OWNER TO expense_user;

--
-- Name: users; Type: TABLE; Schema: public; Owner: expense_user
--

CREATE TABLE public.users (
    id integer NOT NULL,
    name character varying(50) NOT NULL,
    surname character varying(25) NOT NULL,
    login character varying(25) NOT NULL,
    pass character varying(30) NOT NULL,
    email character varying(30) NOT NULL
);


ALTER TABLE public.users OWNER TO expense_user;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: expense_user
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.users_id_seq OWNER TO expense_user;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: expense_user
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: expense id; Type: DEFAULT; Schema: public; Owner: expense_user
--

ALTER TABLE ONLY public.expense ALTER COLUMN id SET DEFAULT nextval('public.expense_id_seq'::regclass);


--
-- Name: expense_type id; Type: DEFAULT; Schema: public; Owner: expense_user
--

ALTER TABLE ONLY public.expense_type ALTER COLUMN id SET DEFAULT nextval('public.expense_type_id_seq'::regclass);


--
-- Name: files id; Type: DEFAULT; Schema: public; Owner: expense_user
--

ALTER TABLE ONLY public.files ALTER COLUMN id SET DEFAULT nextval('public.files_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: expense_user
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Name: expense expense_pkey; Type: CONSTRAINT; Schema: public; Owner: expense_user
--

ALTER TABLE ONLY public.expense
    ADD CONSTRAINT expense_pkey PRIMARY KEY (id);


--
-- Name: expense_type expense_type_pkey; Type: CONSTRAINT; Schema: public; Owner: expense_user
--

ALTER TABLE ONLY public.expense_type
    ADD CONSTRAINT expense_type_pkey PRIMARY KEY (id);


--
-- Name: files files_path_file_key; Type: CONSTRAINT; Schema: public; Owner: expense_user
--

ALTER TABLE ONLY public.files
    ADD CONSTRAINT files_path_file_key UNIQUE (path_file);


--
-- Name: files files_pkey; Type: CONSTRAINT; Schema: public; Owner: expense_user
--

ALTER TABLE ONLY public.files
    ADD CONSTRAINT files_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: expense_user
--

--ALTER TABLE ONLY public.schema_migrations
--    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: expense_type uniq_exp_type; Type: CONSTRAINT; Schema: public; Owner: expense_user
--

ALTER TABLE ONLY public.expense_type
    ADD CONSTRAINT uniq_exp_type UNIQUE (users_id, title);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: expense_user
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: exp_type; Type: INDEX; Schema: public; Owner: expense_user
--

CREATE INDEX exp_type ON public.expense_type USING btree (title);


--
-- Name: expense_type expense_type_users_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: expense_user
--

ALTER TABLE ONLY public.expense_type
    ADD CONSTRAINT expense_type_users_id_fkey FOREIGN KEY (users_id) REFERENCES public.users(id);


--
-- Name: expense expfk; Type: FK CONSTRAINT; Schema: public; Owner: expense_user
--

ALTER TABLE ONLY public.expense
    ADD CONSTRAINT expfk FOREIGN KEY (expense_type_id) REFERENCES public.expense_type(id) ON DELETE CASCADE;


--
-- Name: files files_expense_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: expense_user
--

ALTER TABLE ONLY public.files
    ADD CONSTRAINT files_expense_id_fkey FOREIGN KEY (expense_id) REFERENCES public.expense(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

