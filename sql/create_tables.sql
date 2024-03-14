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
    lat character varying(512),
    lon character varying(512),
    image character varying(255),
    created_at timestamp without time zone,
    updated_at timestamp without time zone
);

ALTER TABLE public.memes OWNER TO esusu;
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

ALTER TABLE public.users OWNER TO esusu;

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

COPY public.memes (id, lat, lon, image, created_at, updated_at) FROM stdin;
1	40.730610   -73.935242    /8Z8dptJEypuLoOQro1WugD855YE.jpg	2024-03-13 00:00:00	2024-03-13 00:00:00
2	40.730610   -73.935242    /ceG9VzoRAVGwivFU403Wc3AHRys.jpg	2024-03-13 00:00:00	2024-03-13 00:00:00
3	40.730610   -73.935242    /3bhkrj58Vtu7enYsRolD1fZdja1.jpg	2024-03-13 00:00:00	2024-03-13 00:00:00
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

