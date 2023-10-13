-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "orders"(
    "id" SERIAL,
    "order_uuid" varchar(36) NOT NULL,
    "track_number" varchar(36) NOT NULL,
    "entry" varchar(16) NOT NULL,
    "locale" varchar(8) NOT NULL,
    "internal_signature" text NOT NULL,
    "customer_id" varchar(36) NOT NULL,
    "delivery_service" varchar(16) NOT NULL,
    "shardkey" varchar(8) NOT NULL,
    "sm_id" integer NOT NULL,
    "date_created" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "oof_shard" varchar(8) NOT NULL,
    "delivery" jsonb NOT NULL DEFAULT '{}'::jsonb,
    "payment" jsonb NOT NULL DEFAULT '{}'::jsonb,
    "items" jsonb NOT NULL DEFAULT '{}'::jsonb,
    CONSTRAINT "order_id" PRIMARY KEY ("id")
);



-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "orders";

-- +goose StatementEnd
