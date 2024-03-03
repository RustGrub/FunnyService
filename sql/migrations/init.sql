CREATE TABLE IF NOT EXISTS public.projects
(
    id bigserial PRIMARY KEY,
    name text NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    CONSTRAINT not_empty_name CHECK (name <> '')
);

ALTER TABLE IF EXISTS public.projects
    OWNER TO postgres;

CREATE TABLE IF NOT EXISTS public.goods
(
    id bigserial NOT NULL,
    project_id bigint NOT NULL,
    name text NOT NULL,
    description text,
    priority bigint NOT NULL,
    removed boolean NOT NULL DEFAULT false,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    CONSTRAINT goods_pkey PRIMARY KEY (id, project_id),
    CONSTRAINT project_id FOREIGN KEY (project_id)
        REFERENCES public.projects (id)
);

ALTER TABLE IF EXISTS public.goods
    OWNER TO postgres;

CREATE OR REPLACE FUNCTION set_priority()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.priority = COALESCE((SELECT MAX(priority) + 1 FROM goods), 1);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_priority_trigger
    BEFORE INSERT ON goods
    FOR EACH ROW
EXECUTE FUNCTION set_priority();

CREATE INDEX goods_index ON GOODS (id, project_id, name);


INSERT INTO projects(name) VALUES ('Первая запись');

/*
CREATE TABLE Goods (
                       GoodId      Int,
                       ProjectId   Int,
                       Name        String,
                       Description Nullable(String),
                       Priority    Int,
                       Removed     bool,
                       EventTime    timestamp,
        INDEX idx_goods (GoodId, ProjectId, Name) TYPE  minmax GRANULARITY 1
) ENGINE = MergeTree()
      ORDER BY (GoodId, ProjectId)
      PRIMARY KEY (GoodId, ProjectId)

 */