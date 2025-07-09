CREATE TABLE tokens (
    id int PRIMARY KEY AUTO_INCREMENT,
    user_id bigint NOT NULL,
    service_id int NOT NULL,
    token varchar(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
);