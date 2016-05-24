CREATE TYPE subscription_types AS ENUM ('free', 'basic', 'advanced');
ALTER TABLE customers ADD COLUMN subscription subscription_types not null default 'free';
