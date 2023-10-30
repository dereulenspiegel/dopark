CREATE TABLE spaces(
    name VARCHAR(256) UNIQUE NOT NULL,
    coords geometry(Point, 925832) NOT NULL,
    number INT PRIMARY KEY
);

CREATE TABLE park_values(
    spaces_id INT NOT NULL,
    free INT NOT NULL,
    total INT NOT NULL,
    time TIMESTAMPTZ NOT NULL
);

SELECT create_hypertable('park_values', 'time', if_not_exists => TRUE);
