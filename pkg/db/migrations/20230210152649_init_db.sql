-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS records (
   uuid uuid UNIQUE NOT NULL,
   amount bigint NOT NULL,
   datetime timestamp NOT NULL,
   PRIMARY KEY (uuid)
);

CREATE INDEX idx_datetime ON records (datetime);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS records;
-- +goose StatementEnd
