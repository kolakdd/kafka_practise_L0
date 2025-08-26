CREATE TYPE payment_currency AS ENUM ('RUB', 'USD', 'EUR');

CREATE TYPE bank_type AS ENUM ('alpha', 'sber', 'tbank');

CREATE TYPE locale_type AS ENUM ('eng', 'ru', 'kz');

CREATE table if not exists "Delivery" (
    uid SERIAL PRIMARY KEY,
    name TEXT,
    phone TEXT,
    zip TEXT,
    city TEXT,
    address TEXT,
    region TEXT,
    email TEXT,
    created_at timestamp without time zone NOT NULL DEFAULT now()
);

CREATE table if not exists "Payment" (
    transaction TEXT PRIMARY KEY,
    request_id TEXT,
    currency payment_currency,
    provider TEXT,
    amount BIGINT,
    payment_dt TIMESTAMP without time ZONE,
    bank bank_type,
    delivery_cost BIGINT,
    goods_total BIGINT,
    custom_fee BIGINT,
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    CHECK (
        goods_total + custom_fee + delivery_cost = amount
    )
);

CREATE table if not EXISTS "Items" (
    chrt_id SERIAL PRIMARY KEY,
    track_number TEXT,
    price bigint,
    rid TEXT,
    name TEXT,
    sale INTEGER,
    size SMALLINT,
    currency payment_currency,
    total_price BIGINT,
    nm_id bigint,
    brand TEXT,
    status SMALLINT,
    created_at timestamp without time zone NOT NULL DEFAULT now()
);

CREATE table if not EXISTS "Customer" (
    uid SERIAL PRIMARY KEY,
    name TEXT,
    created_at timestamp without time zone NOT NULL DEFAULT now()
);

CREATE TABLE if not exists "OrderXItems" (
    order_uid SERIAL NOT NULL,
    item_id SERIAL NOT NULL REFERENCES "Items" (chrt_id),
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    PRIMARY KEY (order_uid, item_id)
);

CREATE table "Order" (
    uid SERIAL PRIMARY KEY,
    track_number TEXT NOT NULL,
    entry TEXT NOT NULL,
    delivery SERIAL REFERENCES "Delivery" (uid) NOT NULL,
    payment TEXT REFERENCES "Payment" (transaction) NOT NULL,
    locale locale_type NOT NULL,
    internal_signature TEXT,
    customer_id SERIAL REFERENCES "Customer" (uid) NOT NULL,
    delivery_service TEXT,
    shardkey SMALLINT,
    sm_id SMALLINT,
    oof_shard SMALLINT,
    created_at timestamp without time zone NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX order_uidx ON "Order" (uid);

alter table "OrderXItems"
add foreign key (order_uid) REFERENCES "Order" (uid);