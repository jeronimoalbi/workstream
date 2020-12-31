BEGIN;

CREATE SCHEMA "account";

CREATE TABLE account.accounts
(
    id character varying(36) NOT NULL,
    name character varying(200) NOT NULL,
    active boolean NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone,
    CONSTRAINT accounts_pk PRIMARY KEY (id),
    CONSTRAINT name_idx UNIQUE (name)
);

COMMIT;
