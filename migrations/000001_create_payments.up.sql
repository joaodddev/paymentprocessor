CREATE TABLE IF NOT EXISTS payments (

    id UUID PRIMARY KEY,

    amount BIGINT NOT NULL,

    currency VARCHAR(3) NOT NULL,

    customer_id UUID NOT NULL,

    status VARCHAR(20) NOT NULL,

    idempotency_key VARCHAR(255) UNIQUE,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMP NOT NULL DEFAULT NOW()

);