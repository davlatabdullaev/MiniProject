alter table if exists users
    add column if not exists login varchar(40) unique not null default uuid_generate_v4();
