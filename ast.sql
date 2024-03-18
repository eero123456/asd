-- Table: public.todos

-- DROP TABLE IF EXISTS public.todos;

CREATE TABLE IF NOT EXISTS public.todos
(
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    user_id bigint NOT NULL,
    text text COLLATE pg_catalog."default",
    completed boolean,
    CONSTRAINT todos_pkey PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.todos
    OWNER to postgres;

CREATE INDEX IF NOT EXISTS idx_user_id
    ON public.todos USING btree
    (user_id ASC NULLS LAST)
    WITH (deduplicate_items=True)
    TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.todos
    CLUSTER ON idx_user_id;





    -- Table: public.user

-- DROP TABLE IF EXISTS public."user";

CREATE TABLE IF NOT EXISTS public."user"
(
    id bigint NOT NULL,
    username uuid DEFAULT gen_random_uuid(),
    CONSTRAINT user_pkey PRIMARY KEY (id),
    CONSTRAINT user_idx_name UNIQUE (username)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public."user"
    OWNER to postgres;