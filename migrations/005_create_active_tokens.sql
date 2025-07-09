CREATE TABLE active_tokens (
    id int PRIMARY KEY AUTO_INCREMENT,
    user_id bigint NOT NULL,
    refresh_token text NOT NULL,
    created_at TIMESTAMP NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    is_active boolean NOT NULL DEFAULT TRUE,
);