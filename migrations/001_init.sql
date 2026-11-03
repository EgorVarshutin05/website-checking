CREATE TABLE sites (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    url TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE checks (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    site_id BIGINT NOT NULL REFERENCES sites(id),
    status_code INTEGER,
    status VARCHAR(10) NOT NULL CHECK (status IN ('up', 'down')),
    response_time_ms BIGINT NOT NULL CHECK (response_time_ms >= 0),
    error_message TEXT,
    checked_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_checks_site_time
ON checks (site_id, checked_at DESC); -- Быстрое получение истории конкретного сайта, начиная с последних проверок