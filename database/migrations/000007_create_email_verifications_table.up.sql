CREATE TABLE email_verifications (
  id VARCHAR(26) NOT NULL PRIMARY KEY,
  user_id VARCHAR(26) NOT NULL,
  token TEXT NOT NULL,
  expires_at DATETIME NOT NULL,
  is_used BOOLEAN DEFAULT FALSE,
  action_type ENUM('register', 'email_change', 'username_change', 'password_change') NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

  -- Indexing
  INDEX idx_token (token(255)),
  INDEX idx_user_id (user_id),
  INDEX idx_expires_used (expires_at, is_used),
  INDEX idx_action_type (action_type),

  -- Foreign key
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
