-- migrate:up
CREATE OR REPLACE FUNCTION update_record_history() RETURNS TRIGGER AS $$
DECLARE
  prev_created_at TIMESTAMP;
  prev_updated_at TIMESTAMP;
BEGIN
  -- Initialize the variables to NULL
  prev_created_at := NULL;
  prev_updated_at := NULL;

  -- Get created_at and updated_at values from the history table if any exist
  IF (TG_OP = 'UPDATE' OR TG_OP = 'DELETE') THEN
    SELECT created_at, updated_at
    INTO prev_created_at, prev_updated_at
    FROM record_history
    WHERE record_id = OLD.id
    ORDER BY updated_at DESC
    LIMIT 1;
  END IF;

  IF (TG_OP = 'DELETE') THEN
    INSERT INTO record_history (record_id, title, description, location, significance, url, start_date, end_date, type, status, created_at, updated_at)
    VALUES (OLD.id, OLD.title, OLD.description, OLD.location, OLD.significance, OLD.url, OLD.start_date, OLD.end_date, OLD.type, OLD.status, prev_created_at, NOW());
    RETURN OLD;
  ELSIF (TG_OP = 'UPDATE') THEN
    INSERT INTO record_history (record_id, title, description, location, significance, url, start_date, end_date, type, status, created_at, updated_at)
    VALUES (NEW.id, NEW.title, NEW.description, NEW.location, NEW.significance, NEW.url, NEW.start_date, NEW.end_date, NEW.type, NEW.status, prev_created_at, NOW());
    RETURN NEW;
  ELSIF (TG_OP = 'INSERT') THEN
    INSERT INTO record_history (record_id, title, description, location, significance, url, start_date, end_date, type, status, created_at, updated_at)
    VALUES (NEW.id, NEW.title, NEW.description, NEW.location, NEW.significance, NEW.url, NEW.start_date, NEW.end_date, NEW.type, NEW.status, NOW(), NOW());
    RETURN NEW;
  END IF;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tr_record_history
AFTER INSERT OR UPDATE OR DELETE ON record
FOR EACH ROW EXECUTE FUNCTION update_record_history();


CREATE OR REPLACE FUNCTION update_impact_history() RETURNS TRIGGER AS $$
DECLARE
  prev_created_at TIMESTAMP;
  prev_updated_at TIMESTAMP;
BEGIN
    -- Initialize the variables to NULL
    prev_created_at := NULL;
    prev_updated_at := NULL;

    -- Get created_at and updated_at values from the history table if any exist
    IF (TG_OP = 'UPDATE' OR TG_OP = 'DELETE') THEN
        SELECT created_at, updated_at
        INTO prev_created_at, prev_updated_at
        FROM impact_history
        WHERE impact_id = OLD.id
        ORDER BY updated_at DESC
        LIMIT 1;
    END IF;

    IF (TG_OP = 'DELETE') THEN
        INSERT INTO impact_history (impact_id, record_id, description, value, category, created_at, updated_at)
        VALUES (OLD.id, OLD.record_id, OLD.description, OLD.value, OLD.category, prev_created_at, NOW());
        RETURN OLD;
    ELSIF (TG_OP = 'UPDATE') THEN
        INSERT INTO impact_history (impact_id, record_id, description, value, category, created_at, updated_at)
        VALUES (NEW.id, NEW.record_id, NEW.description, NEW.value, NEW.category, prev_created_at, NOW());
        RETURN NEW;
    ELSIF (TG_OP = 'INSERT') THEN
        INSERT INTO impact_history (impact_id, record_id, description, value, category, created_at, updated_at)
        VALUES (NEW.id, NEW.record_id, NEW.description, NEW.value, NEW.category, NOW(), NOW());
        RETURN NEW;
    END IF;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tr_impact_history
AFTER INSERT OR UPDATE OR DELETE ON impact
FOR EACH ROW EXECUTE FUNCTION update_impact_history();

-- migrate:down
DROP TRIGGER tr_record_history ON record;
DROP FUNCTION update_record_history();

DROP TRIGGER tr_impact_history ON impact;
DROP FUNCTION update_impact_history();
