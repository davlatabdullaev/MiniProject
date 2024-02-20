create table if not exists customers (
  id uuid primary key not null,
  full_name varchar(50),
  phone varchar(13) unique not null,
  password varchar(128) not null,
);
create table if not exists products (
  id uuid primary key not null,
  name varchar(50) not null,
  price int default 0,
  quantity int default 0
);
create table if not exists baskets (
  id uuid primary key not null,
  product_id uuid references products(id),
  customer_id uuid references customers(id),
  total_sum int default 0
);