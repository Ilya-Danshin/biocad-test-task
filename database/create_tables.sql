CREATE TABLE IF NOT EXISTS files(
    file text NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS data (
    n int NOT NULL,
    mqtt bytea,
    invid text NOT NULL,
    unit_guid uuid NOT NULL,
    msg_id text NOT NULL,
    text text NOT NULL,
    context bytea,
    class text NOT NULL,
    level int NOT NULL,
    area text NOT NULL,
    addr text NOT NULL,
    block text,
    type text,
    bit int,
    invert_bit int
);