CREATE TABLE IF NOT EXISTS group_shares (
  group_id BIGINT NOT NULL,
  user_id BIGINT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (group_id, user_id),
  CONSTRAINT fk_group_shares_group FOREIGN KEY (group_id) REFERENCES task_groups(id) ON DELETE CASCADE,
  CONSTRAINT fk_group_shares_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
