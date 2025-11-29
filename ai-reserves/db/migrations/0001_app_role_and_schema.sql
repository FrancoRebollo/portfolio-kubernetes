-- 0) App role (owner of the schema/objects)
DO $$
BEGIN
  IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'ai_reserves') THEN
    CREATE ROLE ai_reserves
      LOGIN PASSWORD 'ai_reserves'
      NOSUPERUSER NOCREATEDB NOCREATEROLE NOINHERIT;
  ELSE
    ALTER ROLE ai_reserves WITH PASSWORD 'ai_reserves';
  END IF;
END$$;

-- CREATE DATABASE ai_reserves_db;

-- 1) Extensions you need
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- 2) Create a dedicated schema OWNED by your app role
CREATE SCHEMA IF NOT EXISTS ai_res AUTHORIZATION ai_reserves;

-- (Optional) keep public clean so only superuser can create there
REVOKE CREATE ON SCHEMA public FROM PUBLIC;
GRANT USAGE ON SCHEMA public TO PUBLIC;

-- 3) Make your app role default to that schema
ALTER ROLE ai_reserves IN DATABASE ai_reserves_db SET search_path = ai_res, public;

-- 4) Build everything as the app role so it OWNS the objects
SET ROLE ai_reserves;

-- Default privileges for NEW objects the *app role* creates later
ALTER DEFAULT PRIVILEGES IN SCHEMA ai_res
  GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO ai_reserves;
ALTER DEFAULT PRIVILEGES IN SCHEMA ai_res
  GRANT USAGE, SELECT, UPDATE ON SEQUENCES TO ai_reserves;

-- Personas que usan la aplicación
CREATE TABLE ai_res.personas (
    id INT PRIMARY KEY,
    nombre VARCHAR(150) NOT NULL default 'PROVISORIO',
    apellido_razon_social varchar(150) not null DEFAULT 'PROVISORIO',
    persona_juridica  varchar(1) not null default 'N',
    tipo_doc_persona VARCHAR(50),
    nro_doc_persona varchar(70),
    email VARCHAR(100),
    telefono VARCHAR(30),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(100),
    updated_by VARCHAR(100)
);

-- Unidad lógica general (ej: Consultorio, Laboratorio)
CREATE TABLE ai_res.unidad_reserva (
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL,
    descripcion TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(100),
    updated_by VARCHAR(100)
);

-- Tipo dentro de una unidad (ej: Consultorio médico, Consultorio odontológico)
CREATE TABLE ai_res.tipo_unidad_reserva (
    id SERIAL PRIMARY KEY,
    id_unidad_reserva INT NOT NULL REFERENCES ai_res.unidad_reserva(id) ON DELETE CASCADE,
    nombre VARCHAR(100) NOT NULL,
    descripcion TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(100),
    updated_by VARCHAR(100)
);

-- Subtipo con configuración detallada (ej: duración)
CREATE TABLE ai_res.sub_tipo_unidad_reserva (
    id SERIAL PRIMARY KEY,
    id_tipo_unidad_reserva INT NOT NULL REFERENCES ai_res.tipo_unidad_reserva(id) ON DELETE CASCADE,
    nombre VARCHAR(100) NOT NULL,
    descripcion TEXT,
    duracion_reserva_minutos INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(100),
    updated_by VARCHAR(100)
);

-- Agendas por persona (una persona puede tener una agenda asociada a un subtipo de unidad)
CREATE TABLE ai_res.agendas (
    id SERIAL PRIMARY KEY,

    -- 1 de estos dos será NOT NULL
    id_conf_personal INT REFERENCES ai_res.conf_personal(id) ON DELETE CASCADE,
    id_conf_establecimiento INT REFERENCES ai_res.conf_establecimiento(id) ON DELETE CASCADE,

    fecha DATE NOT NULL,
    activa BOOLEAN DEFAULT TRUE,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(100),
    updated_by VARCHAR(100),

    -- Evita generar dos agendas para la misma entidad el mismo día
    UNIQUE(id_conf_personal, fecha),
    UNIQUE(id_conf_establecimiento, fecha)
);


-- Reservas realizadas sobre una agenda
CREATE TABLE ai_res.reservas (
    id SERIAL PRIMARY KEY,
    id_agenda INT NOT NULL REFERENCES ai_res.agendas(id) ON DELETE CASCADE,
    fecha DATE NOT NULL,
    hora_inicio TIME NOT NULL,
    hora_fin TIME NOT NULL,
    id_sub_tipo_unidad_reserva INT NOT NULL REFERENCES ai_res.sub_tipo_unidad_reserva(id)
    id_paciente INT REFERENCES ai_res.personas(id),
    estado VARCHAR(50) DEFAULT 'PENDIENTE', -- pendiente / confirmada / cancelada / finalizada
    observaciones TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(100),
    updated_by VARCHAR(100)
);

-- Relación entre agendas y reservas (si necesitas tabla intermedia, por ahora la FK está en reservas)
-- Si querés agendar múltiples agendas para una misma reserva, usás esta tabla adicional:
-- Si no, ignorala.
CREATE TABLE ai_res.agendas_reservas (
    id SERIAL PRIMARY KEY,
    id_agenda INT NOT NULL REFERENCES ai_res.agendas(id) ON DELETE CASCADE,
    id_reserva INT NOT NULL REFERENCES ai_res.reservas(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(100),
    updated_by VARCHAR(100),
    UNIQUE(id_agenda, id_reserva)
);



CREATE TABLE ai_res.conf_personal (
    id SERIAL PRIMARY KEY,
    id_persona INT NOT NULL REFERENCES ai_res.personas(id) ON DELETE CASCADE,

    hora_inicio TIME NOT NULL,
    hora_fin TIME NOT NULL,

    lunes BOOLEAN DEFAULT TRUE,
    martes BOOLEAN DEFAULT TRUE,
    miercoles BOOLEAN DEFAULT TRUE,
    jueves BOOLEAN DEFAULT TRUE,
    viernes BOOLEAN DEFAULT TRUE,
    sabado BOOLEAN DEFAULT FALSE,
    domingo BOOLEAN DEFAULT FALSE,

    genera_feriados BOOLEAN DEFAULT FALSE,

    -- FLAG CLAVE:
    modo_agenda VARCHAR(20) NOT NULL DEFAULT 'PREGENERADA',
    -- valores: 'PREGENERADA' o 'MULTIAGENDA'

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(100),
    updated_by VARCHAR(100)
);

CREATE TABLE ai_res.conf_personal_sub_tipo_unidad_reserva (
    id SERIAL PRIMARY KEY,
    id_conf_personal INT NOT NULL REFERENCES ai_res.conf_personal(id) ON DELETE CASCADE,
    id_sub_tipo_unidad_reserva INT NOT NULL REFERENCES ai_res.sub_tipo_unidad_reserva(id) ON DELETE CASCADE,

    duracion_reserva_minutos INT,   -- opcional, override del estándar
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(100),
    updated_by VARCHAR(100),

    UNIQUE(id_conf_personal, id_sub_tipo_unidad_reserva)
);

CREATE TABLE ai_res.agenda_slots (
    id SERIAL PRIMARY KEY,
    id_agenda INT NOT NULL REFERENCES ai_res.agendas(id) ON DELETE CASCADE,
    hora_inicio TIME NOT NULL,
    hora_fin TIME NOT NULL,
    estado VARCHAR(20) DEFAULT 'LIBRE',
    id_reserva INT REFERENCES ai_res.reservas(id),

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    UNIQUE(id_agenda, hora_inicio)
);

CREATE TABLE ai_res.conf_establecimiento (
    id SERIAL PRIMARY KEY,
    id_persona INT NOT NULL REFERENCES ai_res.personas(id) ON DELETE CASCADE,

    -- Cada espacio físico concreto (cancha 1, cancha 2, mesa 1, etc.)
    nombre VARCHAR(100) NOT NULL,

    id_sub_tipo_unidad_reserva INT NOT NULL 
        REFERENCES ai_res.sub_tipo_unidad_reserva(id) ON DELETE CASCADE,

    hora_inicio TIME NOT NULL,
    hora_fin TIME NOT NULL,

    lunes BOOLEAN DEFAULT TRUE,
    martes BOOLEAN DEFAULT TRUE,
    miercoles BOOLEAN DEFAULT TRUE,
    jueves BOOLEAN DEFAULT TRUE,
    viernes BOOLEAN DEFAULT TRUE,
    sabado BOOLEAN DEFAULT TRUE,
    domingo BOOLEAN DEFAULT TRUE,

    genera_feriados BOOLEAN DEFAULT FALSE,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(100),
    updated_by VARCHAR(100)
);
