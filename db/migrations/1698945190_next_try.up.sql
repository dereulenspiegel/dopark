TRUNCATE TABLE spaces RESTART IDENTITY;;
TRUNCATE TABLE park_values;
SELECT UpdateGeometrySRID('spaces', 'coords', 4326);
ALTER TABLE park_values ADD COLUMN updated_at TIMESTAMPTZ NOT NULL;
