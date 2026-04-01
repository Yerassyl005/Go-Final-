CREATE TABLE service_points (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT DEFAULT '',
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE queues (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    service_point_id INT REFERENCES service_points(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE operators (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    service_point_id INT REFERENCES service_points(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    phone VARCHAR(20) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE tickets (
    id SERIAL PRIMARY KEY,
    queue_id INT REFERENCES queues(id) ON DELETE CASCADE,
    user_id INT REFERENCES users(id) ON DELETE SET NULL,
    number INT NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    called_at TIMESTAMP,
    completed_at TIMESTAMP
);

CREATE TABLE event_logs (
    id SERIAL PRIMARY KEY,
    ticket_id INT REFERENCES tickets(id) ON DELETE CASCADE,
    event_type TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);