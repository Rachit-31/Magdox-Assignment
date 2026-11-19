ALTER TABLE products
    ADD COLUMN status text NOT NULL DEFAULT 'active'
    CHECK (status IN ('active', 'retired'));
