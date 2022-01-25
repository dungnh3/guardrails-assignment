BEGIN;

DROP TYPE IF EXISTS result_status_enum CASCADE;

ALTER TABLE results
DROP CONSTRAINT results_source_repository_id_fkey;

DROP TABLE IF EXISTS source_repositories CASCADE;

DROP TABLE IF EXISTS results CASCADE;

DROP INDEX IF EXISTS results_status_idx;

DROP INDEX IF EXISTS results_source_repository_id_uidx;


COMMIT;