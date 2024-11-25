CREATE TABLE IF NOT EXISTS users (
  id VARCHAR NOT NULL,
  role varchar NOT NULL DEFAULT 'guest',
  phone VARCHAR  NOT NULL CHECK (phone ~* '^[0-9]{11,11}$'), 
  name VARCHAR NOT NULL DEFAULT '',
  image_url VARCHAR NOT NULL DEFAULT '',
  status_text VARCHAR NOT NULL DEFAULT '',
  is_blocked BOOLEAN NOT NULL DEFAULT FALSE,
  blocked_reason VARCHAR NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ  NOT NULL DEFAULT (now()),
  PRIMARY KEY (id)
);

