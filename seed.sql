CREATE DATABASE url_shortener;

\c url_shortener

--
-- Name: url; Type: TABLE; Schema: public;
--

CREATE TABLE public.url (
    id uuid NOT NULL,
    full_url character varying(255) NOT NULL,
    tiny_url character varying(255) NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);

COMMENT ON COLUMN public.url.full_url IS 'Full URL including scheme (e.g. https://)';
COMMENT ON COLUMN public.url.tiny_url IS 'Tiny URL from Base62 encoding of UUID';
ALTER TABLE ONLY public.url
    ADD CONSTRAINT url_pkey PRIMARY KEY (id);
