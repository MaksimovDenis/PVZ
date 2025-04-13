CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role VARCHAR(255) NOT NULL CHECK (role IN ('moderator', 'employee')),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS pvz (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    city VARCHAR(255) NOT NULL CHECK (city IN ('Москва', 'Санкт-Петербург', 'Казань')),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_pvz_user FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS receptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    pvz_id UUID NOT NULL,
    status VARCHAR(255) NOT NULL CHECK (status IN ('in_progress', 'close')),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    close_at TIMESTAMP,

    CONSTRAINT fk_reception_user FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT fk_reception_pvz FOREIGN KEY (pvz_id) REFERENCES pvz(id)
);

CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    pvz_id UUID NOT NULL,
    reception_id UUID NOT NULL,
    product_type VARCHAR(255) NOT NULL CHECK (product_type IN ('электроника', 'одежда', 'обувь')),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_product_user FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT fk_product_reception FOREIGN KEY (reception_id) REFERENCES receptions(id),
    CONSTRAINT fk_product_pvz FOREIGN KEY (pvz_id) REFERENCES pvz(id)
);

CREATE INDEX idx_receptions_pvz_id ON receptions(pvz_id);
CREATE INDEX idx_products_reception_id ON products(reception_id);

INSERT INTO users (id, email, password_hash, role)
VALUES
    ('11111111-1111-1111-1111-111111111111', 'moderatorDummyLogin@example.com', crypt('moderatorpass', gen_salt('bf')), 'moderator'),
    ('22222222-2222-2222-2222-222222222222', 'employeeDummyLogin@example.com', crypt('employeepass', gen_salt('bf')), 'employee');