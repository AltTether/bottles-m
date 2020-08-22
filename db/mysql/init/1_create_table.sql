CREATE DATABASE
IF NOT EXISTS
bottle;

USE bottle;

CREATE TABLE
IF NOT EXISTS
messages (
        id BIGINT AUTO_INCREMENT NOT NULL PRIMARY KEY,
        text TEXT NOT NULL,
        created_at BIGINT NOT NULL,
        updated_at BIGINT NOT NULL,
        deleted_at BIGINT
        )
ENGINE=InnoDB DEFAULT CHARSET=utf8;
