BEGIN;

CREATE TYPE result_status_enum AS ENUM (
    'queued',
    'in_progress',
    'success',
    'failure'
);

CREATE TABLE source_repositories
(
    id         serial                 NOT NULL
        constraint source_repositories_pkey primary key,
    name       character varying(255) NOT NULL,
    link       character varying(255) NOT NULL,
    is_active  boolean DEFAULT true,
    created_at timestamp without time zone not null,
    updated_at timestamp without time zone
);

CREATE TABLE results
(
    id                   serial                 NOT NULL
        constraint results_pkey primary key,
    source_repository_id integer                NOT NULL,
    status               result_status_enum     NOT NULL,
    name                 character varying(255) NOT NULL,
    link                 character varying(255) NOT NULL,
    findings             jsonb,
    queued_at            timestamp without time zone not null,
    scanning_at          timestamp without time zone,
    finished_at          timestamp without time zone
);

CREATE INDEX results_status_idx on results (status);

CREATE UNIQUE INDEX results_source_repository_id_uidx on results (source_repository_id);

ALTER TABLE results
    ADD CONSTRAINT results_source_repository_id_fkey FOREIGN KEY (source_repository_id) REFERENCES source_repositories (id);

COMMIT;