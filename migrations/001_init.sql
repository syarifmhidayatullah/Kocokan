CREATE TABLE IF NOT EXISTS koco_users (
  id          BIGSERIAL PRIMARY KEY,
  name        VARCHAR(100) NOT NULL,
  email       VARCHAR(200) NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS koco_groups (
  id               BIGSERIAL PRIMARY KEY,
  owner_id         BIGINT NOT NULL REFERENCES koco_users(id),
  name             VARCHAR(100) NOT NULL,
  emoji            VARCHAR(10) DEFAULT '🎰',
  description      TEXT DEFAULT '',
  num_participants INT NOT NULL DEFAULT 10,
  period_type      VARCHAR(20) NOT NULL DEFAULT 'monthly',
  prize_amount     BIGINT NOT NULL DEFAULT 0,
  total_rounds     INT NOT NULL DEFAULT 10,
  is_active        BOOLEAN NOT NULL DEFAULT TRUE,
  created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS koco_participants (
  id         BIGSERIAL PRIMARY KEY,
  group_id   BIGINT NOT NULL REFERENCES koco_groups(id) ON DELETE CASCADE,
  name       VARCHAR(100) NOT NULL,
  phone      VARCHAR(20) DEFAULT '',
  notes      TEXT DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS koco_rounds (
  id           BIGSERIAL PRIMARY KEY,
  group_id     BIGINT NOT NULL REFERENCES koco_groups(id) ON DELETE CASCADE,
  round_number INT NOT NULL,
  winner_id    BIGINT REFERENCES koco_participants(id),
  drawn_at     TIMESTAMPTZ,
  notes        TEXT DEFAULT '',
  created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_koco_groups_owner ON koco_groups(owner_id);
CREATE INDEX IF NOT EXISTS idx_koco_participants_group ON koco_participants(group_id);
CREATE INDEX IF NOT EXISTS idx_koco_rounds_group ON koco_rounds(group_id);
