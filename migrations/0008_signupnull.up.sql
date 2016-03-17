update signups set referrer = '' where referrer is null;
alter table signups alter column referrer set default '';
alter table signups alter column referrer set not null;
