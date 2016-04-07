create unique index idx_signups_email on signups (lower(email) varchar_pattern_ops);
create unique index idx_users_email_active on users (lower(email) varchar_pattern_ops, active);
