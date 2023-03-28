CREATE SCHEMA sh_ilo;

CREATE TABLE sh_ilo.events_messages (
    id bigserial not null,
    created_at timestamp,
    created_id bytea,
    device_id int8,
    object_id int4,
    mes_id int8,
    mes_time timestamp,
    mes_code int4,
    mes_status jsonb,
    mes_data jsonb,
    event_value varchar,
    event_data jsonb
)