CREATE TABLE
    IF NOT EXISTS t_blog(
        id UUID UNIQUE DEFAULT gen_random_uuid(),
        title character varying(128) UNIQUE NOT NULL,
        content text NOT NULL,
        author email_address references t_user(email) NOT NULL,
        create_time timestamp NOT NULL DEFAULT (now()) :: timestamp,
        update_time timestamp NOT NULL DEFAULT (now()) :: timestamp,
        PRIMARY KEY(id)
    );