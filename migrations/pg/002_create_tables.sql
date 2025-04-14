CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE pvzs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    city TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);
CREATE TABLE recepcions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    pvz_id UUID NOT NULL REFERENCES pvz(id),
    status TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);
CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    reception_id UUID NOT NULL REFERENCES intakes(id),
    product_type TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

