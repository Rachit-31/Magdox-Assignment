CREATE TABLE findings (
    id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id     uuid NOT NULL REFERENCES orgs (id),
    product_id uuid NOT NULL REFERENCES products (id),
    title      text NOT NULL,
    severity   text NOT NULL CHECK (severity IN ('critical', 'high', 'medium', 'low')),
    component  text
);
