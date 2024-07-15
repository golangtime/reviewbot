package db

import "database/sql"

func MigrateDB(db *sql.DB) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS repositories (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			owner  text NOT NULL,
			name   text NOT NULL,
			min_approvals integer NOT NULL DEFAULT 1,
			active bool NOT NULL DEFAULT true,
			provider text NOT NULL DEFAULT 'github'
		)
	`)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS notification_queue (
			id           INTEGER PRIMARY KEY AUTOINCREMENT,
			rule_id      INTEGER NOT NULL,
			recepient    text NOT NULL,
			link         text NOT NULL,
			user_id      integer NOT NULL,
			created_at   timestamp NOT NULL,
			reserved_for timestamp NOT NULL,
			status       text NOT NULL DEFAULT '',
			source       text NOT NULL,

			UNIQUE (rule_id, recepient, link, user_id)
		)
	`)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS notification_rules (
			id           INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id      integer NOT NULL,
			notification_type text NOT NULL,
			provider_id  text NOT NULL,
			chat_id      text,
			priority     integer NOT NULL,

			UNIQUE (user_id, notification_type)
		)
	`)
	if err != nil {
		panic(err)
	}
}
