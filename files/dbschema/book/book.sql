CREATE SCHEMA bookschema;
--end

SET search_path TO bookschema;
--end

CREATE TABLE book (
    "id" SERIAL,
    "title" varchar(80),
    "author" varchar(120),
    PRIMARY KEY("id")
);
--end