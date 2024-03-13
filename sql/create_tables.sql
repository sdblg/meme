--
-- PostgreSQL database dump
--

-- Dumped from database version 14.5 (Debian 14.5-1.pgdg110+1)
-- Dumped by pg_dump version 14.5 (Homebrew)

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
-- Name: memes; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.memes (
    id integer NOT NULL,
    title character varying(512),
    release_date date,
    runtime integer,
    mpaa_rating character varying(10),
    description text,
    image character varying(255),
    created_at timestamp without time zone,
    updated_at timestamp without time zone
);

--
-- Name: memes_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

ALTER TABLE public.memes ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.memes_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: users; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.users (
    id integer NOT NULL,
    first_name character varying(255),
    last_name character varying(255),
    email character varying(255),
    password character varying(255),
    created_at timestamp without time zone,
    updated_at timestamp without time zone
);


--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

ALTER TABLE public.users ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);

--
-- Data for Name: memes; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.memes (id, title, release_date, runtime, mpaa_rating, description, image, created_at, updated_at) FROM stdin;
1	Highlander	1986-03-07	116	R	He fought his first battle on the Scottish Highlands in 1536. He will fight his greatest battle on the streets of New York City in 1986. His name is Connor MacLeod. He is immortal.	/8Z8dptJEypuLoOQro1WugD855YE.jpg	2022-09-23 00:00:00	2022-09-23 00:00:00
2	Raiders of the Lost Ark	1981-06-12	115	PG-13	Archaeology professor Indiana Jones ventures to seize a biblical artefact known as the Ark of the Covenant. While doing so, he puts up a fight against Renee and a troop of Nazis.	/ceG9VzoRAVGwivFU403Wc3AHRys.jpg	2022-09-23 00:00:00	2022-09-23 00:00:00
3	The Godfather	1972-03-24	175	18A	The aging patriarch of an organized crime dynasty in postwar New York City transfers control of his clandestine empire to his reluctant youngest son.	/3bhkrj58Vtu7enYsRolD1fZdja1.jpg	2022-09-23 00:00:00	2022-09-23 00:00:00
\.

--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.users (id, first_name, last_name, email, password, created_at, updated_at) FROM stdin;
1	Admin	User	admin@esusu.com	$2a$14$wVsaPvJnJJsomWArouWCtusem6S/.Gauq/GjOIEHpyh2DAMmso1wy	2022-09-23 00:00:00	2022-09-23 00:00:00
\.

--
-- Name: memes memes_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.memes
    ADD CONSTRAINT memes_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);

--
-- PostgreSQL database dump complete
--

