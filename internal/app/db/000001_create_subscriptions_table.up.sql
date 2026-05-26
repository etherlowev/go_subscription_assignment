begin transaction;

create table if not exists public.subscriptions (
    id uuid primary key not null,
    name text not null,
    price integer not null,
    user_id uuid not null,
    start_date date not null,
    end_date date
);

commit;