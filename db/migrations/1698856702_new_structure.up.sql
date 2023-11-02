ALTER TABLE spaces ALTER COLUMN coords DROP NOT NULL;
TRUNCATE TABLE spaces;
TRUNCATE TABLE park_values;
CREATE SEQUENCE space_number_sq;
ALTER TABLE spaces ALTER COLUMN number SET DEFAULT nextval('space_number_sq');
ALTER SEQUENCE space_number_sq OWNED BY spaces.number;
