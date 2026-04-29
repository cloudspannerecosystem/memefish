ALTER TABLE foo ALTER COLUMN bar TIMESTAMP NOT NULL DEFAULT (pending_commit_timestamp()) ON UPDATE (pending_commit_timestamp())
