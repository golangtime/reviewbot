package db

import "database/sql"

func MigrateDB(db *sql.DB) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS repositories (
			owner  text NOT NULL,
			name   text NOT NULL,
			min_approvals integer NOT NULL DEFAULT 1,
			active bool NOT NULL DEFAULT true
		)
	`)
	if err != nil {
		panic(err)
	}

	// _, err = db.Exec(`
	// 	CREATE TABLE IF NOT EXISTS peers (
	// 		peer_id integer NOT NULL,
	// 		active  bool NOT NULL DEFAULT true
	// 	)
	// `)
	// if err != nil {
	// 	panic(err)
	// }

	// _, err = db.Exec(`
	// 	CREATE TABLE IF NOT EXISTS merge_request_rules (
	// 		repo_id integer NOT NULL,
	// 		rule    text NOT NULL,
	// 		active  bool NOT NULL DEFAULT true
	// 	)
	// `)
	// if err != nil {
	// 	panic(err)
	// }
}
