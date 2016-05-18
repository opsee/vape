CREATE TYPE status_types AS ENUM ('invited', 'active', 'inactive');
ALTER TABLE users ADD COLUMN status status_types not null default 'invited';
ALTER TABLE users ADD COLUMN perms BIGINT not null default 0;
ALTER TABLE signups ADD COLUMN customer_id varchar(36) not null default '';
ALTER TABLE signups ADD COLUMN perms BIGINT not null default 0;
