CREATE TABLE user_sessions (
  id VARCHAR(26) PRIMARY KEY,
  user_id VARCHAR(26) NOT NULL,
  refresh_token_hash TEXT NOT NULL,
  device_id TEXT,
  ip_address VARCHAR(45),
  user_agent TEXT,
  revoked BOOLEAN DEFAULT FALSE,
  expires_at DATETIME NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;