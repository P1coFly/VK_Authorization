DROP TABLE IF EXISTS public."USERS" CASCADE;

CREATE TABLE public."USERS"
(
    id bigserial NOT NULL,
	email text UNIQUE,
	password bytea,
    PRIMARY KEY (id)
);

ALTER TABLE public."USERS"
    OWNER to postgres;		

--------------------------------------------

CREATE ROLE auth_service WITH LOGIN PASSWORD '12345678';
GRANT CONNECT ON DATABASE "VK_Authorization" TO auth_service;
GRANT SELECT, INSERT ON ALL TABLES IN SCHEMA public TO auth_service;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO auth_service;


-- REVOKE ALL PRIVILEGES ON DATABASE "VK_MOVIES" FROM auth_service;
-- REVOKE ALL PRIVILEGES ON ALL TABLES IN SCHEMA public FROM auth_service;
-- REVOKE ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public FROM auth_service;
-- DROP ROLE auth_service