CREATE TABLE forgot_tokens (
    token text NOT NULL, 
    created_at TIMESTAMP NOT NULL DEFAULT now(), 
    user_id bigint NOT NULL, 
    CONSTRAINT fk_users FOREIGN KEY (user_id) REFERENCES users(id)
);