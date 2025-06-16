CREATE TABLE IF NOT EXISTS short_urls (
    id SERIAL PRIMARY KEY,
    short_code VARCHAR(16) UNIQUE NOT NULL,
    long_url TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    expire_at TIMESTAMP,
    clicks BIGINT DEFAULT 0
);

CREATE TABLE IF NOT EXISTS click_logs (
  id SERIAL PRIMARY KEY,
  short_code VARCHAR(16) NOT NULL,
  clicked_at TIMESTAMP NOT NULL DEFAULT NOW(),
  referrer TEXT,
  user_agent TEXT,
  device_type VARCHAR(32),  -- e.g. mobile, desktop, tablet
  os VARCHAR(64),           -- e.g. Android, iOS, Windows
  browser VARCHAR(64),      -- e.g. Chrome, Safari, Firefox
  country VARCHAR(64),      -- e.g. Vietnam
  city VARCHAR(64),         -- e.g. Ho Chi Minh City
  ip_address INET,          -- e.g. 192.168.1.1
  FOREIGN KEY (short_code) REFERENCES short_urls(short_code)
);
