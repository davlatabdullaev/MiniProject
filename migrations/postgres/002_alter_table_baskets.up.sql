alter table if exists baskets
    add column if not exists quantity int default 0;