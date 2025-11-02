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
