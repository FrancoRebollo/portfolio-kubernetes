-- 0) App role (owner of the schema/objects)
DO $$
BEGIN
  IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'auth_security') THEN
    CREATE ROLE auth_security
      LOGIN PASSWORD 'auth_security'
      NOSUPERUSER NOCREATEDB NOCREATEROLE NOINHERIT;
  ELSE
    ALTER ROLE auth_security WITH PASSWORD 'auth_security';
  END IF;
END$$;

-- CREATE DATABASE auth_security_db;

-- 1) Extensions you need
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- 2) Create a dedicated schema OWNED by your app role
CREATE SCHEMA IF NOT EXISTS sec AUTHORIZATION auth_security;

-- (Optional) keep public clean so only superuser can create there
REVOKE CREATE ON SCHEMA public FROM PUBLIC;
GRANT USAGE ON SCHEMA public TO PUBLIC;

-- 3) Make your app role default to that schema
ALTER ROLE auth_security IN DATABASE auth_security_db SET search_path = sec, public;

-- 4) Build everything as the app role so it OWNS the objects
SET ROLE auth_security;

-- Default privileges for NEW objects the *app role* creates later
ALTER DEFAULT PRIVILEGES IN SCHEMA sec
  GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO auth_security;
ALTER DEFAULT PRIVILEGES IN SCHEMA sec
  GRANT USAGE, SELECT, UPDATE ON SEQUENCES TO auth_security;

-- 5) Your model (moved to schema sec + small fixes)
CREATE TABLE sec.version_modelo (
  service_name       varchar(60)  NOT NULL,
  version_modelo     varchar(60)  NOT NULL,
  fecha_last_update  timestamp DEFAULT now() NOT NULL
);

CREATE TABLE sec.api (
  uuid_api           uuid NOT NULL,
  api                varchar(60) NOT NULL,
  "version"          varchar(12) NOT NULL,
  fecha_last_update  timestamp DEFAULT now() NOT NULL,
  actualizado_por    varchar(30) DEFAULT current_user NOT NULL,
  CONSTRAINT pk_api PRIMARY KEY (uuid_api)
);

CREATE TABLE sec.api_key (
  api_key                varchar(60) DEFAULT gen_random_uuid() NOT NULL,
  app_origen             varchar(60) NOT NULL,
  estado                 varchar(15) DEFAULT 'ACTIVO' NOT NULL,
  req_2fa                char(1) DEFAULT 'N' NOT NULL,
  ctd_hs_access_token_valido integer DEFAULT 1 NOT NULL,
  req_usuario_db         char(1) DEFAULT 'S' NOT NULL,
  fecha_vigencia         date DEFAULT current_date NOT NULL,
  fecha_fin_vigencia     date,
  ctrl_limite_acceso_tiempo char(1) DEFAULT 'N',
  ctd_accesos_unidad_tiempo integer,
  unidad_tiempo_acceso   varchar(15) DEFAULT 'MINUTO',
  fecha_last_update      timestamp DEFAULT now() NOT NULL,
  actualizado_por        varchar(30) DEFAULT current_user NOT NULL,
  is_super_user          char(1) DEFAULT 'N' NOT NULL,
  CONSTRAINT pk_api_key PRIMARY KEY (api_key)
);

CREATE TABLE sec.tipo_canal_digital_df (
  tipo_canal_digital   varchar(25) NOT NULL,
  acceso_revocado      char(1) DEFAULT 'N' NOT NULL,
  fecha_last_update    date DEFAULT current_date NOT NULL,
  actualizado_por      varchar(30) DEFAULT current_user NOT NULL,
  CONSTRAINT pk_canal_digital_df PRIMARY KEY (tipo_canal_digital)
);

CREATE TABLE sec.location (
  id_location        varchar(15) NOT NULL,
  fecha_last_update  date DEFAULT current_date NOT NULL,
  actualizado_por    varchar(30) DEFAULT current_user NOT NULL,
  CONSTRAINT pk_locations PRIMARY KEY (id_location)
);

CREATE TABLE sec.acceso_api (
  api_key              varchar(60) NOT NULL,
  uuid_api             uuid NOT NULL,
  fecha_last_update    timestamp DEFAULT now() NOT NULL,
  actualizado_por      varchar(30) NOT NULL,
  CONSTRAINT pk_api_key_api PRIMARY KEY (uuid_api, api_key),
  CONSTRAINT unq_acceso_api_api_key UNIQUE (api_key, uuid_api)
);

-- Prefer identity to serial
CREATE TABLE sec.persona (
  id_persona            integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  last_location         varchar(15) NOT NULL,
  acceso_revocado       char(1) DEFAULT 'N' NOT NULL,
  fecha_last_update     date DEFAULT current_date NOT NULL,
  actualizado_por       varchar(30) DEFAULT current_user
);

CREATE TABLE sec.canal_digital_persona (
  id_canal_digital_persona integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  id_persona           integer NOT NULL,
  tipo_canal_digital   varchar(25) NOT NULL,
  password_acceso_hash varchar(256) NOT NULL,
  id_mail_persona      integer,
  id_te_persona        integer,
  login_name           varchar(100),
  canal_validado       char(1) DEFAULT 'N',
  fecha_validacion_canal date,
  acceso_revocado      char(1) DEFAULT 'N' NOT NULL,
  "req_2fa"            char(1) DEFAULT 'N' NOT NULL,
  fecha_last_update    date DEFAULT current_date NOT NULL,
  actualizado_por      varchar(30) DEFAULT current_user NOT NULL,
  CONSTRAINT unique_login_name UNIQUE (login_name)
);

CREATE TABLE sec.token (
  id_token                 integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  api_key                  varchar(60) NOT NULL,
  id_canal_digital_persona integer NOT NULL,
  access_token             varchar(500) NOT NULL,
  fecha_creacion_token     timestamp,
  fecha_exp_access_token   timestamp,
  refresh_token            varchar(500) NOT NULL,
  fecha_exp_refresh_token  timestamp,
  acceso_revocado          char(1) DEFAULT 'N' NOT NULL,
  fecha_last_update        date DEFAULT current_date NOT NULL,
  actualizado_por          varchar(30) DEFAULT current_user,
  "2fa_seed"               varchar(100),
  last_code_2fa            numeric
);

CREATE TABLE sec.hist_token (
  id_hist_token           integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  id_token                integer NOT NULL,
  api_key                 varchar(60) NOT NULL,
  id_canal_digital_persona integer NOT NULL,
  access_token            varchar(500) NOT NULL,
  fecha_creacion_token    timestamp,
  fecha_exp_access_token  timestamp,
  refresh_token           varchar(500) NOT NULL,
  fecha_exp_refresh_token timestamp,
  acceso_revocado         char(1) DEFAULT 'N' NOT NULL,
  fecha_last_update       date DEFAULT current_date NOT NULL,
  actualizado_por         varchar(30) DEFAULT current_user
);

CREATE TABLE sec.error_log (
  id_error_log      integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  message_error     varchar(5000) NOT NULL,
  endpoint          varchar(400),
  id_tipo_error     integer DEFAULT 0,
  ip_address        varchar(50),
  id_persona        integer NOT NULL,
  canal_digital     varchar(60) NOT NULL,
  api_key           varchar(60) NOT NULL,
  id_token          integer NOT NULL,
  access_token      varchar(500) NOT NULL,
  fecha_last_update date DEFAULT current_date NOT NULL,
  actualizado_por   varchar(30) DEFAULT current_user
);

-- FKs
ALTER TABLE sec.acceso_api
  ADD CONSTRAINT fk_acceso_api_api      FOREIGN KEY (uuid_api) REFERENCES sec.api(uuid_api),
  ADD CONSTRAINT fk_acceso_api_api_key  FOREIGN KEY (api_key)  REFERENCES sec.api_key(api_key);

ALTER TABLE sec.canal_digital_persona
  ADD CONSTRAINT fk_cdp_persona  FOREIGN KEY (id_persona)        REFERENCES sec.persona(id_persona),
  ADD CONSTRAINT fk_cdp_tipo     FOREIGN KEY (tipo_canal_digital) REFERENCES sec.tipo_canal_digital_df(tipo_canal_digital);

ALTER TABLE sec.error_log
  ADD CONSTRAINT fk_err_persona FOREIGN KEY (id_persona) REFERENCES sec.persona(id_persona);

ALTER TABLE sec.persona
  ADD CONSTRAINT fk_persona_location FOREIGN KEY (last_location) REFERENCES sec.location(id_location);

ALTER TABLE sec.token
  ADD CONSTRAINT fk_token_api_key  FOREIGN KEY (api_key) REFERENCES sec.api_key(api_key),
  ADD CONSTRAINT fk_token_cdp      FOREIGN KEY (id_canal_digital_persona) REFERENCES sec.canal_digital_persona(id_canal_digital_persona);

ALTER TABLE sec.canal_digital_persona DROP COLUMN id_mail_persona;
ALTER TABLE sec.canal_digital_persona ADD COLUMN mail_persona varchar(100);
ALTER TABLE sec.canal_digital_persona DROP COLUMN id_te_persona;
ALTER TABLE sec.canal_digital_persona ADD COLUMN telefono_persona varchar(50);

-- Indexes (schema-qualified)
CREATE INDEX idx_persona_2 ON sec.persona (id_persona, acceso_revocado);

CREATE INDEX idx_cdp_1 ON sec.canal_digital_persona (id_persona, tipo_canal_digital);
CREATE INDEX idx_cdp_2 ON sec.canal_digital_persona (tipo_canal_digital, id_persona, canal_validado);
CREATE INDEX idx_cdp_3 ON sec.canal_digital_persona (tipo_canal_digital, login_name);
CREATE INDEX idx_cdp_4 ON sec.canal_digital_persona (login_name);
CREATE INDEX idx_cdp_5 ON sec.canal_digital_persona (id_persona, tipo_canal_digital, acceso_revocado);

CREATE INDEX idx_tipo_canal_1 ON sec.tipo_canal_digital_df (tipo_canal_digital);
CREATE INDEX idx_tipo_canal_2 ON sec.tipo_canal_digital_df (tipo_canal_digital, acceso_revocado);

CREATE INDEX idx_token_2 ON sec.token (id_canal_digital_persona, api_key);
CREATE INDEX idx_token_3 ON sec.token (api_key, id_canal_digital_persona);
CREATE INDEX idx_token_4 ON sec.token (api_key, id_canal_digital_persona, refresh_token);
CREATE INDEX idx_token_5 ON sec.token (api_key, id_canal_digital_persona, access_token);

CREATE INDEX idx_api_key_1 ON sec.api_key (api_key);
CREATE INDEX idx_api_key_2 ON sec.api_key (api_key, fecha_fin_vigencia);

INSERT INTO sec.api_key (
    api_key, app_origen, estado, req_2fa, ctd_hs_access_token_valido,
    req_usuario_db, fecha_vigencia, fecha_fin_vigencia, ctrl_limite_acceso_tiempo,
    ctd_accesos_unidad_tiempo, unidad_tiempo_acceso, fecha_last_update, actualizado_por
) VALUES
('3b2a4a2d-5f91-4a83-9f21-9a8c2adcc123', 'APP_SECURITY', 'ACTIVO', 'N', 1, 'N', CURRENT_DATE, NULL, 'N', NULL, 'MINUTO', CURRENT_TIMESTAMP, 'FRANCO'),
('5e9b31f7-c423-4a6c-9a5f-88929a5a7f0c', 'APP_MESSAGING', 'ACTIVO', 'N', 2, 'N', CURRENT_DATE, NULL, 'S', 10, 'MINUTO', CURRENT_TIMESTAMP, 'FRANCO'),
('7c1fd418-2140-463a-88ef-fdc821b4f571', 'APP_ADMIN', 'ACTIVO', 'N', 1, 'N', CURRENT_DATE, NULL, 'N', NULL, 'MINUTO', CURRENT_TIMESTAMP, 'FRANCO'),
('9a2f563b-9b3a-4025-89ee-8b7cefb4d482', 'APP_API_INTEGRATION', 'ACTIVO', 'N', 3, 'N', CURRENT_DATE, NULL, 'S', 5, 'MINUTO', CURRENT_TIMESTAMP, 'FRANCO'),
('b3c52b87-6d34-4b3d-b5a5-0b2a6a7c813f', 'APP_ADITIONAL', 'ACTIVO', 'N', 1, 'N', CURRENT_DATE, NULL, 'N', NULL, 'MINUTO', CURRENT_TIMESTAMP, 'FRANCO');


UPDATE sec.API_KEY SET IS_SUPER_USER = 'S' WHERE APP_ORIGEN = 'APP_ADMIN';

INSERT INTO sec.location (id_location, fecha_last_update, actualizado_por) VALUES (1, CURRENT_DATE, 'FRANCO');
INSERT INTO sec.location (id_location, fecha_last_update, actualizado_por) VALUES (2, CURRENT_DATE, 'FRANCO');
INSERT INTO sec.location (id_location, fecha_last_update, actualizado_por) VALUES (3, CURRENT_DATE, 'FRANCO');
INSERT INTO sec.location (id_location, fecha_last_update, actualizado_por) VALUES (4, CURRENT_DATE, 'FRANCO');
INSERT INTO sec.location (id_location, fecha_last_update, actualizado_por) VALUES (5, CURRENT_DATE, 'FRANCO');
INSERT INTO sec.location (id_location, fecha_last_update, actualizado_por) VALUES (6, CURRENT_DATE, 'FRANCO');
INSERT INTO sec.location (id_location, fecha_last_update, actualizado_por) VALUES (7, CURRENT_DATE, 'FRANCO');
INSERT INTO sec.location (id_location, fecha_last_update, actualizado_por) VALUES (8, CURRENT_DATE, 'FRANCO');
INSERT INTO sec.location (id_location, fecha_last_update, actualizado_por) VALUES (9, CURRENT_DATE, 'FRANCO');
INSERT INTO sec.location (id_location, fecha_last_update, actualizado_por) VALUES (10, CURRENT_DATE, 'FRANCO');
INSERT INTO sec.location (id_location, fecha_last_update, actualizado_por) VALUES (11, CURRENT_DATE, 'FRANCO');
INSERT INTO sec.location (id_location, fecha_last_update, actualizado_por) VALUES (12, CURRENT_DATE, 'FRANCO');
INSERT INTO sec.location (id_location, fecha_last_update, actualizado_por) VALUES (13, CURRENT_DATE, 'FRANCO');
INSERT INTO sec.location (id_location, fecha_last_update, actualizado_por) VALUES (14, CURRENT_DATE, 'FRANCO');
INSERT INTO sec.location (id_location, fecha_last_update, actualizado_por) VALUES (15, CURRENT_DATE, 'FRANCO');
INSERT INTO sec.location (id_location, fecha_last_update, actualizado_por) VALUES (16, CURRENT_DATE, 'FRANCO');
INSERT INTO sec.location (id_location, fecha_last_update, actualizado_por) VALUES (17, CURRENT_DATE, 'FRANCO');
INSERT INTO sec.location (id_location, fecha_last_update, actualizado_por) VALUES (18, CURRENT_DATE, 'FRANCO');
INSERT INTO sec.location (id_location, fecha_last_update, actualizado_por) VALUES (19, CURRENT_DATE, 'FRANCO');
INSERT INTO sec.location (id_location, fecha_last_update, actualizado_por) VALUES (20, CURRENT_DATE, 'FRANCO');
INSERT INTO sec.location (id_location, fecha_last_update, actualizado_por) VALUES (21, CURRENT_DATE, 'FRANCO');
INSERT INTO sec.location (id_location, fecha_last_update, actualizado_por) VALUES (22, CURRENT_DATE, 'FRANCO');
INSERT INTO sec.location (id_location, fecha_last_update, actualizado_por) VALUES (23, CURRENT_DATE, 'FRANCO');


-- Return to superuser at the end (optional)
RESET ROLE;
