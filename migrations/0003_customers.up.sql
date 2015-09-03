drop table orgs cascade;
CREATE TABLE customers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v1mc(),
    name character varying(255),
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);

-- fix users to use customers
alter table users drop column org_id cascade;
alter table users add column customer_id UUID;
alter table users add constraint fk_users_customers foreign key (customer_id) references customers(id);
create index idx_users_customers on users (customer_id);

-- fix bastions to use customers
alter table bastions drop column org_id cascade;
alter table bastions add column customer_id UUID;
alter table bastions add constraint fk_bastions_customers foreign key (customer_id) references customers(id);
create index idx_bastions_customers on bastions (customer_id);
