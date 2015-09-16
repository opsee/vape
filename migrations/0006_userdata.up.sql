CREATE OR REPLACE FUNCTION json_merge(data jsonb, merge_data jsonb)
RETURNS jsonb
IMMUTABLE
LANGUAGE sql
AS $$
    SELECT ('{'||string_agg(to_json(key)||':'||value, ',')||'}')::jsonb
    FROM (
        WITH to_merge AS (
            SELECT * FROM jsonb_each(merge_data)
        )
        SELECT *
        FROM jsonb_each(data)
        WHERE key NOT IN (SELECT key FROM to_merge)
        UNION ALL
        SELECT * FROM to_merge
    ) t;
$$;

CREATE FUNCTION insert_userdata()
RETURNS trigger
LANGUAGE plpgsql
AS $$
	BEGIN
	insert into userdata (user_id) values(NEW.id);
	RETURN NEW;
  END;
$$;

create table userdata (
	user_id integer not null references users(id) on delete cascade,
	data jsonb default '{}' not null,
  created_at timestamp with time zone DEFAULT now() NOT NULL,
  updated_at timestamp with time zone DEFAULT now() NOT NULL
);

CREATE TRIGGER update_userdata BEFORE UPDATE ON userdata FOR EACH ROW EXECUTE PROCEDURE update_time();
CREATE TRIGGER insert_userdata AFTER INSERT ON users FOR EACH ROW EXECUTE PROCEDURE insert_userdata();

-- add userdata to existing users
insert into userdata (user_id) select id from users;
