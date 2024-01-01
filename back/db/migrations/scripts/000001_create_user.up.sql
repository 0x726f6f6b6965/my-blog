CREATE EXTENSION IF NOT EXISTS citext;

CREATE DOMAIN email_address AS citext
  CHECK ( value ~ '^[a-zA-Z0-9.!#$%&''*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$' );

CREATE TABLE
    IF NOT EXISTS t_user(
        id UUID UNIQUE DEFAULT gen_random_uuid(),
        username character varying(128) UNIQUE NOT NULL,
        email email_address UNIQUE NOT NULL,
        salt character varying(64) NOT NULL,
        password character varying(128) NOT NULL,
        create_time timestamp NOT NULL DEFAULT (now()) :: timestamp,
        update_time timestamp NOT NULL DEFAULT (now()) :: timestamp,
        PRIMARY KEY(id)
    )