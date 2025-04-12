CREATE UNIQUE INDEX idx_users_email ON users(email);

CREATE INDEX idx_users_role_id ON users(role_id);

CREATE INDEX idx_pvz_created_at ON pvz(created_at);
CREATE INDEX idx_pvz_city_id ON pvz(city_id);

CREATE INDEX idx_intakes_pvz_id ON intakes(pvz_id);
CREATE INDEX idx_intakes_status_id ON intakes(status_id);
CREATE INDEX idx_intakes_pvz_status_created_at ON intakes(pvz_id, status_id, created_at DESC);

CREATE INDEX idx_products_intake_id ON products(intake_id);
CREATE INDEX idx_products_created_at ON products(created_at);
CREATE INDEX idx_products_intake_created_at ON products(intake_id, created_at DESC);

CREATE UNIQUE INDEX idx_roles_name ON roles(name);
CREATE UNIQUE INDEX idx_cities_name ON cities(name);
CREATE UNIQUE INDEX idx_product_types_name ON product_types(name);
CREATE UNIQUE INDEX idx_intake_statuses_name ON intake_statuses(name);
