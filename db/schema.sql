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
-- Name: impact_revisions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.impact_revisions (
    id text NOT NULL,
    impact_id text NOT NULL
);


--
-- Name: impacts; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.impacts (
    id text NOT NULL,
    subject_id text NOT NULL,
    reasoning character varying NOT NULL,
    category character varying(50) NOT NULL,
    value integer NOT NULL
);


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.schema_migrations (
    version character varying(128) NOT NULL
);


--
-- Name: subject_relations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.subject_relations (
    subject_1 text NOT NULL,
    subject_2 text NOT NULL
);


--
-- Name: subjects; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.subjects (
    id text NOT NULL,
    title character varying NOT NULL,
    summary character varying NOT NULL,
    subject_type character varying(50),
    url character varying NOT NULL,
    weight integer,
    from_date date NOT NULL,
    until_date date NOT NULL
);


--
-- Name: impact_revisions impact_revisions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.impact_revisions
    ADD CONSTRAINT impact_revisions_pkey PRIMARY KEY (id);


--
-- Name: impacts impacts_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.impacts
    ADD CONSTRAINT impacts_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: subject_relations subject_relations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.subject_relations
    ADD CONSTRAINT subject_relations_pkey PRIMARY KEY (subject_1, subject_2);


--
-- Name: subjects subjects_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.subjects
    ADD CONSTRAINT subjects_pkey PRIMARY KEY (id);


--
-- Name: impact_revisions impact_revisions_impact_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.impact_revisions
    ADD CONSTRAINT impact_revisions_impact_id_fkey FOREIGN KEY (impact_id) REFERENCES public.impacts(id);


--
-- Name: impacts impacts_subject_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.impacts
    ADD CONSTRAINT impacts_subject_id_fkey FOREIGN KEY (subject_id) REFERENCES public.subjects(id);


--
-- Name: subject_relations subject_relations_subject_1_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.subject_relations
    ADD CONSTRAINT subject_relations_subject_1_fkey FOREIGN KEY (subject_1) REFERENCES public.subjects(id);


--
-- Name: subject_relations subject_relations_subject_2_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.subject_relations
    ADD CONSTRAINT subject_relations_subject_2_fkey FOREIGN KEY (subject_2) REFERENCES public.subjects(id);


--
-- PostgreSQL database dump complete
--


--
-- Dbmate schema migrations
--

INSERT INTO public.schema_migrations (version) VALUES
    ('20250218182804'),
    ('20250218183047');
