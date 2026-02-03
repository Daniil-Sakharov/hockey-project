-- +goose Up
-- +goose StatementBegin
ALTER TABLE match_events ADD COLUMN IF NOT EXISTS penalty_reason_code VARCHAR(20);

COMMENT ON COLUMN match_events.penalty_reason_code IS 'Код нарушения (ПОДН, ГРУБ, НП-АТ и т.д.)';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE match_events DROP COLUMN IF EXISTS penalty_reason_code;
-- +goose StatementEnd
