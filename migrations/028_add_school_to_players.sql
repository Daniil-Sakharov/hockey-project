-- +goose Up
-- +goose StatementBegin
ALTER TABLE players ADD COLUMN school VARCHAR(255);

COMMENT ON COLUMN players.school IS 'Школа/место воспитания игрока (для FHSPB)';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE players DROP COLUMN IF EXISTS school;
-- +goose StatementEnd
