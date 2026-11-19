CREATE TABLE orgs (
    id   uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name text NOT NULL UNIQUE
);

CREATE TABLE users (
    id     uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id uuid NOT NULL REFERENCES orgs (id),
    email  text NOT NULL UNIQUE,
    name   text NOT NULL,
    role   text NOT NULL CHECK (role IN ('viewer', 'analyst', 'admin'))
);

CREATE TABLE products (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id      uuid NOT NULL REFERENCES orgs (id),
    name        text NOT NULL,
    sku         text NOT NULL,
    device_type text NOT NULL CHECK (device_type IN ('camera', 'router', 'gateway', 'other')),
    created_at  timestamptz NOT NULL DEFAULT now(),
    UNIQUE (org_id, sku)
);
