-- 0) App role (owner of the schema/objects)
DO $$
BEGIN
  IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'async_messaging') THEN
    CREATE ROLE async_messaging
      LOGIN PASSWORD 'async_messaging'
      NOSUPERUSER NOCREATEDB NOCREATEROLE NOINHERIT;
  ELSE
    ALTER ROLE async_messaging WITH PASSWORD 'async_messaging';
  END IF;
END$$;

-- CREATE DATABASE async_messaging_db;

-- 1) Extensions you need
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- 2) Create a dedicated schema OWNED by your app role
CREATE SCHEMA IF NOT EXISTS asyn_m AUTHORIZATION async_messaging;

-- (Optional) keep public clean so only superuser can create there
REVOKE CREATE ON SCHEMA public FROM PUBLIC;
GRANT USAGE ON SCHEMA public TO PUBLIC;

-- 3) Make your app role default to that schema
ALTER ROLE async_messaging IN DATABASE async_messaging_db SET search_path = asyn_m, public;

-- 4) Build everything as the app role so it OWNS the objects
SET ROLE async_messaging;

-- Default privileges for NEW objects the *app role* creates later
ALTER DEFAULT PRIVILEGES IN SCHEMA asyn_m
  GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO async_messaging;
ALTER DEFAULT PRIVILEGES IN SCHEMA asyn_m
  GRANT USAGE, SELECT, UPDATE ON SEQUENCES TO async_messaging;

-- 5) Your model (moved to schema asyn_m + small fixes)
CREATE TABLE asyn_m.version_modelo (
  service_name       varchar(60)  NOT NULL,
  version_modelo     varchar(60)  NOT NULL,
  fecha_last_update  timestamp DEFAULT now() NOT NULL
);


CREATE TABLE asyn_m.message_event (
    id_event        VARCHAR(50) PRIMARY KEY,
    source_system   VARCHAR(50) NOT NULL,
    destiny_system      VARCHAR(50) NOT NULL,
    payload         JSONB NOT NULL,
    status          VARCHAR(20) NOT NULL DEFAULT 'RECEIVED',
    error_msg       TEXT,
    fecha_recepcion TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    fecha_envio     TIMESTAMP,
    fecha_last_update TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    actualizado_por VARCHAR(30) NOT NULL DEFAULT 'SYSTEM'
);

ALTER TABLE asyn_m.message_event
ADD CONSTRAINT uk_event_source UNIQUE (id_event, source_system);


CREATE TABLE asyn_m.dead_letter (
    id_dead             SERIAL PRIMARY KEY,
    original_event_id   VARCHAR(50) NOT NULL,
    source_system       VARCHAR(50) NOT NULL,
    queue_name          VARCHAR(100) NOT NULL,
    payload             JSONB NOT NULL,
    error_msg           TEXT NOT NULL,
    retry_count         INTEGER NOT NULL DEFAULT 0,
    fecha_error         TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    fecha_last_update   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    actualizado_por     VARCHAR(50) NOT NULL DEFAULT 'SYSTEM',
    CONSTRAINT fk_dead_event
        FOREIGN KEY (original_event_id, source_system)
        REFERENCES asyn_m.message_event (id_event, source_system)
        ON DELETE CASCADE
);


CREATE INDEX idx_message_event_status ON asyn_m.message_event(status);
