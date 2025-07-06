CREATE TABLE user_roles(
  user_id VARCHAR(26),
  role_id VARCHAR(26),

  PRIMARY KEY(user_id, role_id),
  FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY(role_id) REFERENCES roles(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;