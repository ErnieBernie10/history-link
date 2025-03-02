-- migrate:up
create table record_history (
    id uuid primary key default gen_random_uuid (),
    record_id uuid references record (id) on delete cascade,
    title varchar(255) not null,
    description varchar(255) not null,
    location varchar(255) default '',
    significance varchar(255) default '',
    url varchar(255) not null,
    start_date timestamp,
    end_date timestamp,
    type smallint not null,
    status smallint not null,
    created_at timestamp not null default now (),
    updated_at timestamp not null default now ()
);

create table impact_history (
    id uuid primary key default gen_random_uuid (),
    impact_id uuid references impact (id) on delete cascade,
    record_id uuid references record (id) on delete cascade,
    description varchar(255) not null,
    value smallint not null,
    category smallint not null,
    created_at timestamp not null default now (),
    updated_at timestamp not null default now ()
);

-- migrate:down
drop table impact_history;

drop table record_history;
