CREATE SCHEMA users;

CREATE SEQUENCE users.users_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE users.list
(
    user_id   integer DEFAULT nextval('users.users_seq'::regclass) NOT NULL,
    dialog_id bigint                                               NOT NULL
);

ALTER TABLE ONLY users.list
    ADD CONSTRAINT pk_user PRIMARY KEY (user_id);

CREATE UNIQUE INDEX idx1_user ON users.list USING btree (dialog_id);