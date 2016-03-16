alter table signups add column referrer character varying(64);
create index idx_signups_referrer on signups (referrer);
