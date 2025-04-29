DROP TABLE IF EXISTS station_entities;

CREATE TABLE station_entities(
    id int NOT NULL AUTO_INCREMENT,
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME,
    name VARCHAR(255) NOT NULL,
    x int NOT NULL,
    y int NOT NULL,
    PRIMARY KEY(id)
);
