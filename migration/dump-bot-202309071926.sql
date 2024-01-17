--
-- PostgreSQL database dump
--

-- Dumped from database version 15.4 (Debian 15.4-1.pgdg120+1)
-- Dumped by pg_dump version 15.3

-- Started on 2023-09-07 19:26:38

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

--
-- TOC entry 7 (class 2615 OID 24585)
-- Name: questions; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA questions;


--
-- TOC entry 6 (class 2615 OID 24584)
-- Name: users; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA users;


--
-- TOC entry 216 (class 1259 OID 24586)
-- Name: questions_seq; Type: SEQUENCE; Schema: questions; Owner: -
--

CREATE SEQUENCE questions.questions_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- TOC entry 218 (class 1259 OID 24596)
-- Name: list; Type: TABLE; Schema: questions; Owner: -
--

CREATE TABLE questions.list (
    question_id integer DEFAULT nextval('questions.questions_seq'::regclass) NOT NULL,
    question_slug character varying(100) NOT NULL,
    question_text character varying NOT NULL
);


--
-- TOC entry 219 (class 1259 OID 24604)
-- Name: list; Type: TABLE; Schema: users; Owner: -
--

CREATE TABLE users.list (
    user_id integer NOT NULL,
    dialog_id integer NOT NULL
);


--
-- TOC entry 217 (class 1259 OID 24587)
-- Name: users_seq; Type: SEQUENCE; Schema: users; Owner: -
--

CREATE SEQUENCE users.users_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 3357 (class 0 OID 24596)
-- Dependencies: 218
-- Data for Name: list; Type: TABLE DATA; Schema: questions; Owner: -
--

INSERT INTO questions.list VALUES (1, 'first', 'AgreeWithStatement');


--
-- TOC entry 3358 (class 0 OID 24604)
-- Dependencies: 219
-- Data for Name: list; Type: TABLE DATA; Schema: users; Owner: -
--

INSERT INTO users.list VALUES (1, 2);
INSERT INTO users.list VALUES (3, 7);
INSERT INTO users.list VALUES (4, 8);
INSERT INTO users.list VALUES (6, 9);


--
-- TOC entry 3364 (class 0 OID 0)
-- Dependencies: 216
-- Name: questions_seq; Type: SEQUENCE SET; Schema: questions; Owner: -
--

SELECT pg_catalog.setval('questions.questions_seq', 1, false);


--
-- TOC entry 3365 (class 0 OID 0)
-- Dependencies: 217
-- Name: users_seq; Type: SEQUENCE SET; Schema: users; Owner: -
--

SELECT pg_catalog.setval('users.users_seq', 1, false);


--
-- TOC entry 3209 (class 2606 OID 24603)
-- Name: list pk_question; Type: CONSTRAINT; Schema: questions; Owner: -
--

ALTER TABLE ONLY questions.list
    ADD CONSTRAINT pk_question PRIMARY KEY (question_id);


--
-- TOC entry 3212 (class 2606 OID 24608)
-- Name: list pk_user; Type: CONSTRAINT; Schema: users; Owner: -
--

ALTER TABLE ONLY users.list
    ADD CONSTRAINT pk_user PRIMARY KEY (user_id);


--
-- TOC entry 3207 (class 1259 OID 24609)
-- Name: idx1_question; Type: INDEX; Schema: questions; Owner: -
--

CREATE UNIQUE INDEX idx1_question ON questions.list USING btree (question_slug);


--
-- TOC entry 3210 (class 1259 OID 24610)
-- Name: idx1_user; Type: INDEX; Schema: users; Owner: -
--

CREATE UNIQUE INDEX idx1_user ON users.list USING btree (dialog_id);


-- Completed on 2023-09-07 19:26:38

--
-- PostgreSQL database dump complete
--

