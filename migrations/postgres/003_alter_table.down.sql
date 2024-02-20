alter table customers
    drop column if exists created_at,
    drop column if exists updated_at,
    drop column if exists deleted_at;

alter table products
    drop column if exists created_at,
    drop column if exists updated_at,
    drop column if exists deleted_at;

alter table baskets
    drop column if exists created_at,
    drop column if exists updated_at,
    drop column if exists deleted_at;

