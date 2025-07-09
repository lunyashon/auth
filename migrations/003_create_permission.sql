CREATE TABLE user_tokens (
    user_id bigint NOT NULL,
    service_id int NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    CONSTRAINT fk_users FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT fk_services FOREIGN KEY (service_id) REFERENCES services(id),
); 