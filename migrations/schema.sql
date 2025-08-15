--
-- PostgreSQL database dump
--

-- Dumped from database version 12.22 (Ubuntu 12.22-0ubuntu0.20.04.1)
-- Dumped by pg_dump version 14.18 (Ubuntu 14.18-0ubuntu0.22.04.1)

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
-- Name: active_tokens; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.active_tokens (
    id integer NOT NULL,
    user_id bigint NOT NULL,
    refresh_token text NOT NULL,
    created_at timestamp without time zone,
    expires_at timestamp without time zone,
    is_active boolean,
    ip text DEFAULT '0.0.0.0'::text,
    device text DEFAULT 'none'::text
);


--
-- Name: active_tokens_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.active_tokens_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: active_tokens_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.active_tokens_id_seq OWNED BY public.active_tokens.id;


--
-- Name: confirm_email_tokens; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.confirm_email_tokens (
    code text NOT NULL,
    created_at timestamp without time zone DEFAULT now(),
    expires_at timestamp without time zone DEFAULT (now() + '00:15:00'::interval),
    email character varying(255)
);


--
-- Name: forgot_tokens; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.forgot_tokens (
    token text NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    user_id bigint NOT NULL,
    expires_at timestamp without time zone DEFAULT (now() + '00:15:00'::interval)
);


--
-- Name: permission; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.permission (
    user_id integer,
    service_id integer,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    active boolean DEFAULT false NOT NULL,
    expires_at timestamp without time zone DEFAULT now() NOT NULL
);


--
-- Name: services; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.services (
    id integer NOT NULL,
    name character varying(255),
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


--
-- Name: services_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.services_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: services_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.services_id_seq OWNED BY public.services.id;


--
-- Name: tokens; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.tokens (
    token character varying(80),
    is_used integer,
    services character varying(255),
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


--
-- Name: users; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.users (
    id bigint NOT NULL,
    email character varying(255),
    login character varying(255),
    password character varying(255),
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    confirmed boolean DEFAULT false
);


--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: users_profile; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.users_profile (
    user_id integer NOT NULL,
    name text DEFAULT 'none'::text NOT NULL,
    last_name text DEFAULT 'none'::text NOT NULL,
    phone text,
    photo_url text DEFAULT '/assets/standart-profile.svg'::text NOT NULL
);


--
-- Name: active_tokens id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.active_tokens ALTER COLUMN id SET DEFAULT nextval('public.active_tokens_id_seq'::regclass);


--
-- Name: services id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.services ALTER COLUMN id SET DEFAULT nextval('public.services_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Name: active_tokens active_tokens_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.active_tokens
    ADD CONSTRAINT active_tokens_pkey PRIMARY KEY (id);


--
-- Name: services services_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.services
    ADD CONSTRAINT services_pkey PRIMARY KEY (id);


--
-- Name: users users_email_unique; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_unique UNIQUE (email);


--
-- Name: users users_login_unique; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_login_unique UNIQUE (login);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: idx_active_tokens_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_active_tokens_user_id ON public.active_tokens USING btree (user_id);


--
-- Name: idx_permission_service_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_permission_service_id ON public.permission USING btree (service_id);


--
-- Name: idx_permission_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_permission_user_id ON public.permission USING btree (user_id);


--
-- Name: idx_services_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_services_id ON public.services USING btree (id);


--
-- Name: idx_users_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_users_id ON public.users USING btree (id);


--
-- Name: idx_users_profile_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_users_profile_user_id ON public.users_profile USING btree (user_id);


--
-- Name: confirm_email_tokens confirm_email_tokens_email_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.confirm_email_tokens
    ADD CONSTRAINT confirm_email_tokens_email_fkey FOREIGN KEY (email) REFERENCES public.users(email) ON DELETE CASCADE;


--
-- Name: permission fk_services; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.permission
    ADD CONSTRAINT fk_services FOREIGN KEY (service_id) REFERENCES public.services(id) ON DELETE CASCADE;


--
-- Name: permission fk_users; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.permission
    ADD CONSTRAINT fk_users FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: forgot_tokens fk_users; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.forgot_tokens
    ADD CONSTRAINT fk_users FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: users_profile users_profile_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users_profile
    ADD CONSTRAINT users_profile_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

