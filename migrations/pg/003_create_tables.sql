CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role_id INTEGER NOT NULL REFERENCES roles(id),
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE pvz (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    city_id INTEGER NOT NULL REFERENCES cities(id),
    created_at TIMESTAMP NOT NULL DEFAULT now()
);
CREATE TABLE intakes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    pvz_id UUID NOT NULL REFERENCES pvz(id),
    status_id INTEGER NOT NULL REFERENCES intake_statuses(id),
    created_at TIMESTAMP NOT NULL DEFAULT now()
);
CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    intake_id UUID NOT NULL REFERENCES intakes(id),
    product_type_id INTEGER NOT NULL REFERENCES product_types(id),
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

