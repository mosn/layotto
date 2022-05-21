-- Table: public.layotto_alloc

-- DROP TABLE IF EXISTS public.layotto_alloc;

CREATE TABLE IF NOT EXISTS public.layotto_alloc
(
    id integer NOT NULL,
    biz_tag character(128) COLLATE pg_catalog."default" NOT NULL,
    max_id bigint NOT NULL,
    step integer NOT NULL,
    description character(256) COLLATE pg_catalog."default" NOT NULL,
    update_time bigint NOT NULL,
    CONSTRAINT layotto_alloc_pkey PRIMARY KEY (id)
    )
    WITH (
        OIDS = FALSE
        )
    TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.layotto_alloc OWNER to postgres;

INSERT INTO public.layotto_alloc(
    id, biz_tag, max_id, step, description, update_time)
VALUES (1, "test", 30, 1, "description", 1653035287);

INSERT INTO public.layotto_alloc(
    id, biz_tag, max_id, step, description, update_time)
VALUES (2, "azh", 5000, 1, "test", 1652932196);