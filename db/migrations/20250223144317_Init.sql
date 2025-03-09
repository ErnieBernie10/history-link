-- migrate:up
create table record (
    id uuid primary key default gen_random_uuid (),
    title varchar(255) not null,
    description varchar(255) not null,
    location varchar(255) default '',
    significance varchar(255) default '',
    url varchar(255) not null,
    start_date timestamp,
    end_date timestamp,
    type smallint not null,
    status smallint not null
);

create table impact (
    id uuid primary key default gen_random_uuid (),
    record_id uuid not null references record (id) on delete cascade,
    description varchar(255) not null,
    value smallint not null,
    category smallint not null
);

create table link (
    id uuid primary key default gen_random_uuid (),
    record_id uuid not null references record (id) on delete cascade,
    record_id2 uuid not null references record (id) on delete cascade,
    strength smallint not null
);

create table source (
    id uuid primary key default gen_random_uuid (),
    record_id uuid not null references record (id) on delete cascade,
    title varchar(255) not null,
    type smallint not null,
    url varchar(255) not null,
    description varchar(255)
);

-- migrate:down
drop table source;

drop table link;

drop table impact;

drop table record;
