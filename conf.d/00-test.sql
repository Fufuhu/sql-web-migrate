-- Code generated by yamltmpl from academic_background.yaml. DO NOT EDIT.

-- +migrate Up
CREATE TABLE test
(
    id serial NOT NULL,
    name varchar(100) NOT NULL,
    PRIMARY KEY(id)
) WITHOUT OIDS;

ALTER SEQUENCE test_id_SEQ INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 1 RESTART 1 CYCLE;

-- +migrate Down
DROP TABLE test;