DROP TABLE IF EXISTS train_entities;

CREATE TABLE train_entities(
    id int NOT NULL AUTO_INCREMENT,
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME,
    name VARCHAR(255) NOT NULL,
    top_speed int NOT NULL,
    x int NOT NULL,
    y int NOT NULL,
    status ENUM(
        'Travelling', 
        'Transferring', 
        'Unused',
        'Emergency'
    ) NOT NULL, 
    PRIMARY KEY(id)
);
