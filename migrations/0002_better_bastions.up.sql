-- add orgs and "active" to bastions. we have to give org a default of 0 for existing bastions
-- we'll remove that later
alter table only bastions add column org_id integer default 0 not null;
alter table only bastions add column active boolean default false not null;

-- insert our test user and associate existing bastions to their org
do $$
declare
	orgid orgs.id%TYPE;
begin
	insert into orgs (name) values (NULL) returning id into orgid;
	insert into users (email, name, password_hash, verified, active, admin, org_id)
		values ('cliff@leaninto.it', 'cliff', '$2a$10$gTv2kAfkIjy0GX67zxbM.ucgP37Na7qDrj42uBnBBrpU3Q/AvWl.G', 
			true, true, true, orgid);
	update bastions set org_id = orgid, active = true where org_id = 0;
end $$;

-- add indexes and constraints
alter table only bastions add constraint fk_bastions_orgs foreign key (org_id) references orgs(id);
create index idx_bastions_orgs on bastions (org_id);
create index idx_users_orgs on users (org_id);

-- we don't ever want an org id of 0
alter table bastions alter column org_id drop default;
