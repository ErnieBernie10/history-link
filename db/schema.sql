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
-- Name: update_impact_history(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.update_impact_history() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
DECLARE
  prev_created_at TIMESTAMP;
  prev_updated_at TIMESTAMP;
BEGIN
    -- Initialize the variables to NULL
    prev_created_at := NULL;
    prev_updated_at := NULL;

    -- Get created_at and updated_at values from the history table if any exist
    IF (TG_OP = 'UPDATE' OR TG_OP = 'DELETE') THEN
        SELECT created_at, updated_at
        INTO prev_created_at, prev_updated_at
        FROM impact_history
        WHERE impact_id = OLD.id
        ORDER BY updated_at DESC
        LIMIT 1;
    END IF;

    IF (TG_OP = 'DELETE') THEN
        INSERT INTO impact_history (impact_id, record_id, description, value, category, created_at, updated_at)
        VALUES (OLD.id, OLD.record_id, OLD.description, OLD.value, OLD.category, prev_created_at, NOW());
        RETURN OLD;
    ELSIF (TG_OP = 'UPDATE') THEN
        INSERT INTO impact_history (impact_id, record_id, description, value, category, created_at, updated_at)
        VALUES (NEW.id, NEW.record_id, NEW.description, NEW.value, NEW.category, prev_created_at, NOW());
        RETURN NEW;
    ELSIF (TG_OP = 'INSERT') THEN
        INSERT INTO impact_history (impact_id, record_id, description, value, category, created_at, updated_at)
        VALUES (NEW.id, NEW.record_id, NEW.description, NEW.value, NEW.category, NOW(), NOW());
        RETURN NEW;
    END IF;

    IF (TG_OP = 'INSERT' OR TG_OP = 'UPDATE' OR TG_OP = 'DELETE') THEN
        INSERT INTO record_history (record_id, title, description, location, significance, url, start_date, end_date, type, status, created_at, updated_at)
        SELECT r.id, r.title, r.description, r.location, r.significance, r.url, r.start_date, r.end_date, r.type, r.status, rh.created_at, NOW()
        FROM record r
        LEFT JOIN record_history rh ON r.id = rh.record_id
        WHERE r.id = NEW.record_id
        ORDER BY rh.updated_at DESC
        LIMIT 1;
    END IF;
    END;
$$;


--
-- Name: update_record_history(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.update_record_history() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
DECLARE
  prev_created_at TIMESTAMP;
  prev_updated_at TIMESTAMP;
BEGIN
  -- Initialize the variables to NULL
  prev_created_at := NULL;
  prev_updated_at := NULL;

  -- Get created_at and updated_at values from the history table if any exist
  IF (TG_OP = 'UPDATE' OR TG_OP = 'DELETE') THEN
    SELECT created_at, updated_at
    INTO prev_created_at, prev_updated_at
    FROM record_history
    WHERE record_id = OLD.id
    ORDER BY updated_at DESC
    LIMIT 1;
  END IF;

  IF (TG_OP = 'DELETE') THEN
    INSERT INTO record_history (record_id, title, description, location, significance, url, start_date, end_date, type, status, created_at, updated_at)
    VALUES (OLD.id, OLD.title, OLD.description, OLD.location, OLD.significance, OLD.url, OLD.start_date, OLD.end_date, OLD.type, OLD.status, prev_created_at, NOW());
    RETURN OLD;
  ELSIF (TG_OP = 'UPDATE') THEN
    INSERT INTO record_history (record_id, title, description, location, significance, url, start_date, end_date, type, status, created_at, updated_at)
    VALUES (NEW.id, NEW.title, NEW.description, NEW.location, NEW.significance, NEW.url, NEW.start_date, NEW.end_date, NEW.type, NEW.status, prev_created_at, NOW());
    RETURN NEW;
  ELSIF (TG_OP = 'INSERT') THEN
    INSERT INTO record_history (record_id, title, description, location, significance, url, start_date, end_date, type, status, created_at, updated_at)
    VALUES (NEW.id, NEW.title, NEW.description, NEW.location, NEW.significance, NEW.url, NEW.start_date, NEW.end_date, NEW.type, NEW.status, NOW(), NOW());
    RETURN NEW;
  END IF;
END;
$$;


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: impact; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.impact (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    record_id uuid NOT NULL,
    description character varying(255) NOT NULL,
    value smallint NOT NULL,
    category smallint NOT NULL
);


--
-- Name: impact_history; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.impact_history (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    impact_id uuid,
    record_id uuid,
    description character varying(255) NOT NULL,
    value smallint NOT NULL,
    category smallint NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);


--
-- Name: link; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.link (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    record_id uuid NOT NULL,
    record_id2 uuid NOT NULL,
    strength smallint NOT NULL
);


--
-- Name: record; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.record (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    title character varying(255) NOT NULL,
    description character varying(255) NOT NULL,
    location character varying(255) DEFAULT ''::character varying,
    significance character varying(255) DEFAULT ''::character varying,
    url character varying(255) NOT NULL,
    start_date timestamp without time zone,
    end_date timestamp without time zone,
    type smallint NOT NULL,
    status smallint NOT NULL
);


--
-- Name: record_history; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.record_history (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    record_id uuid,
    title character varying(255) NOT NULL,
    description character varying(255) NOT NULL,
    location character varying(255) DEFAULT ''::character varying,
    significance character varying(255) DEFAULT ''::character varying,
    url character varying(255) NOT NULL,
    start_date timestamp without time zone,
    end_date timestamp without time zone,
    type smallint NOT NULL,
    status smallint NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.schema_migrations (
    version character varying(128) NOT NULL
);


--
-- Name: source; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.source (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    record_id uuid NOT NULL,
    title character varying(255) NOT NULL,
    type smallint NOT NULL,
    url character varying(255) NOT NULL,
    description character varying(255)
);


--
-- Name: impact_history impact_history_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.impact_history
    ADD CONSTRAINT impact_history_pkey PRIMARY KEY (id);


--
-- Name: impact impact_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.impact
    ADD CONSTRAINT impact_pkey PRIMARY KEY (id);


--
-- Name: link link_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.link
    ADD CONSTRAINT link_pkey PRIMARY KEY (id);


--
-- Name: record_history record_history_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.record_history
    ADD CONSTRAINT record_history_pkey PRIMARY KEY (id);


--
-- Name: record record_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.record
    ADD CONSTRAINT record_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: source source_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.source
    ADD CONSTRAINT source_pkey PRIMARY KEY (id);


--
-- Name: idx_impact_history_impact_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_impact_history_impact_id ON public.impact_history USING btree (impact_id);


--
-- Name: idx_record_history_record_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_record_history_record_id ON public.record_history USING btree (record_id);


--
-- Name: idx_record_impacts; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_record_impacts ON public.impact USING btree (record_id);


--
-- Name: impact tr_impact_history; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER tr_impact_history AFTER INSERT OR UPDATE ON public.impact FOR EACH ROW EXECUTE FUNCTION public.update_impact_history();


--
-- Name: record tr_record_history; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER tr_record_history AFTER INSERT OR UPDATE ON public.record FOR EACH ROW EXECUTE FUNCTION public.update_record_history();


--
-- Name: impact_history impact_history_impact_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.impact_history
    ADD CONSTRAINT impact_history_impact_id_fkey FOREIGN KEY (impact_id) REFERENCES public.impact(id) ON DELETE CASCADE;


--
-- Name: impact_history impact_history_record_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.impact_history
    ADD CONSTRAINT impact_history_record_id_fkey FOREIGN KEY (record_id) REFERENCES public.record(id) ON DELETE CASCADE;


--
-- Name: impact impact_record_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.impact
    ADD CONSTRAINT impact_record_id_fkey FOREIGN KEY (record_id) REFERENCES public.record(id) ON DELETE CASCADE;


--
-- Name: link link_record_id2_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.link
    ADD CONSTRAINT link_record_id2_fkey FOREIGN KEY (record_id2) REFERENCES public.record(id) ON DELETE CASCADE;


--
-- Name: link link_record_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.link
    ADD CONSTRAINT link_record_id_fkey FOREIGN KEY (record_id) REFERENCES public.record(id) ON DELETE CASCADE;


--
-- Name: record_history record_history_record_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.record_history
    ADD CONSTRAINT record_history_record_id_fkey FOREIGN KEY (record_id) REFERENCES public.record(id) ON DELETE CASCADE;


--
-- Name: source source_record_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.source
    ADD CONSTRAINT source_record_id_fkey FOREIGN KEY (record_id) REFERENCES public.record(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--


--
-- Dbmate schema migrations
--

INSERT INTO public.schema_migrations (version) VALUES
    ('20250223144317'),
    ('20250302122704'),
    ('20250302131546'),
    ('20250303074713');
