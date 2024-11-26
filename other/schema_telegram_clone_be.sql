create table if not exists public.users
(
    id         serial
        primary key,
    created_at timestamp with time zone default CURRENT_TIMESTAMP,
    updated_at timestamp with time zone default CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone,
    username   text
        unique,
    password   text,
    email      text
        unique,
    is_online  boolean                  default false,
    last_seen  timestamp with time zone default CURRENT_TIMESTAMP
);

alter table public.users
    owner to admin;

create index if not exists idx_users_username
    on public.users (username);

create index if not exists idx_users_email
    on public.users (email);

create index if not exists idx_users_deleted_at
    on public.users (deleted_at);

create table if not exists public.rooms
(
    id          serial
        primary key,
    created_at  timestamp with time zone default CURRENT_TIMESTAMP,
    updated_at  timestamp with time zone default CURRENT_TIMESTAMP,
    deleted_at  timestamp with time zone,
    name        text,
    description text
);

alter table public.rooms
    owner to admin;

create index if not exists idx_rooms_deleted_at
    on public.rooms (deleted_at);

create table if not exists public.room_users
(
    room_id integer not null
        references public.rooms
        constraint fk_room_users_room
            references public.rooms,
    user_id integer not null
        references public.users
        constraint fk_room_users_user
            references public.users,
    primary key (room_id, user_id)
);

alter table public.room_users
    owner to admin;

create table if not exists public.messages
(
    id         serial
        primary key,
    created_at timestamp with time zone default CURRENT_TIMESTAMP,
    updated_at timestamp with time zone default CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone,
    content    text not null,
    user_id    bigint
        references public.users
        constraint fk_messages_user
            references public.users,
    room_id    bigint
        references public.rooms
        constraint fk_rooms_messages
            references public.rooms,
    is_private boolean,
    to_user_id bigint
        references public.users
        constraint fk_messages_to_user
            references public.users,
    file_url   text
);

alter table public.messages
    owner to admin;

create index if not exists idx_messages_user_id
    on public.messages (user_id);

create index if not exists idx_messages_room_id
    on public.messages (room_id);

create index if not exists idx_messages_to_user_id
    on public.messages (to_user_id);

create index if not exists idx_messages_deleted_at
    on public.messages (deleted_at);

