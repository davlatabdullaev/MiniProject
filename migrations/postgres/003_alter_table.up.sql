
alter table customers
    add column if not exists created_at timestamp default now(),
    add column if not exists updated_at timestamp,
    add column if not exists deleted_at integer default 0;

alter table products
    add column if not exists created_at timestamp default now(),
    add column if not exists updated_at timestamp,
    add column if not exists deleted_at integer default 0;

alter table baskets
    add column if not exists created_at timestamp default now(),
    add column if not exists updated_at timestamp,
    add column if not exists deleted_at integer default 0;
