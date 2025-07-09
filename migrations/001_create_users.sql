CREATE TABLE users (
    id bigint PRIMARY KEY AUTO_INCREMENT,
    email varchar(255) NOT NULL UNIQUE,
    password varchar(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
);