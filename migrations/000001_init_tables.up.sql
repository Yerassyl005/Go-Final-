CREATE TABLE service_points (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE queues (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    service_point_id INT REFERENCES service_points(id),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE operators (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    service_point_id INT REFERENCES service_points(id),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE tickets (
    id SERIAL PRIMARY KEY,
    queue_id INT REFERENCES queues(id),
    number INT NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    called_at TIMESTAMP,
    completed_at TIMESTAMP
);

CREATE TABLE event_logs (
    id SERIAL PRIMARY KEY,
    ticket_id INT REFERENCES tickets(id),
    event_type TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);