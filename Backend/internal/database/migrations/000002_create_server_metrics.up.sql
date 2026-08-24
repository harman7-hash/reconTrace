CREATE EXTENSION IF NOT EXISTS timescaledb;
CREATE TABLE IF NOT EXISTS server_metrics (
    time TIMESTAMPTZ NOT NULL,
    user_id UUID REFERENCES users(id) NOT NULL,
    cpu_util DOUBLE PRECISION,
    mem_util DOUBLE PRECISION,
    disk_r   DOUBLE PRECISION, 
    disk_w   DOUBLE PRECISION
);



SELECT create_hypertable(
    'server_metrics',
    'time',
    chunk_time_interval => INTERVAL '1 day',
    if_not_exists => TRUE
);

SELECT add_retention_policy(
    'server_metrics',
    INTERVAL '30 days'
);