-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    Username VARCHAR(255) UNIQUE NOT NULL,
    Email VARCHAR(255) UNIQUE NOT NULL,
    Password VARCHAR(255) NOT NULL,
    createdAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE url
ADD COLUMN assignedToId INT,
ADD CONSTRAINT fk_url FOREIGN KEY (assignedToId) REFERENCES users(id);

ALTER TABLE url
ALTER COLUMN assignedToId SET NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE url 
DROP CONSTRAINT IF EXISTS fk_url,
DROP COLUMN IF EXISTS assignedToId;

DROP TABLE IF EXISTS users;
-- +goose StatementEnd
