CREATE TABLE IF NOT EXISTS roles (
    id serial primary key,
    uuid uuid not null unique,
    name varchar(32)
);

CREATE TABLE IF NOT EXISTS statuses (
    id serial primary key,
    uuid uuid not null unique,
    name varchar(64)
);

CREATE TABLE IF NOT EXISTS product_types (
    id serial primary key,
    uuid uuid not null unique,
    name varchar(64)
);

CREATE TABLE IF NOT EXISTS cities (
    id serial primary key,
    uuid uuid not null unique,
    name varchar(64)
);

CREATE TABLE IF NOT EXISTS users (
    id uuid primary key,
    role_id integer references roles(id),
    email varchar(320),
    password varchar(64)
);

CREATE TABLE IF NOT EXISTS sessions (
    id uuid primary key,
    user_id uuid references users(id),
    token varchar(64)
);

CREATE TABLE IF NOT EXISTS pvzs (
    id uuid primary key,
    city_id integer references cities(id),
    registration_date timestamp with time zone
);

CREATE TABLE IF NOT EXISTS receptions (
    id uuid primary key,
    pvz_id uuid references pvzs(id),
    status_id integer references statuses(id),
    date_time timestamp with time zone
);

CREATE TABLE IF NOT EXISTS products (
    id uuid primary key,
    type_id integer references product_types(id),
    reception_id uuid references receptions(id),
    date_time timestamp with time zone
);


