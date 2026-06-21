-- +goose Up
CREATE TABLE products (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name        TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    price       BIGINT NOT NULL,
    currency    TEXT NOT NULL DEFAULT 'USD',
    category_id UUID NOT NULL REFERENCES categories(id),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_products_category_id ON products(category_id);

-- +goose Down
DROP TABLE IF EXISTS products;
