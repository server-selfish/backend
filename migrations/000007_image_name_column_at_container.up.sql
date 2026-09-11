ALTER TABLE container
    ADD COLUMN image_name VARCHAR;
    UPDATE container SET image_name = '' WHERE image_name IS NULL;
    ALTER TABLE container ALTER COLUMN image_name SET NOT NULL;
