package db

import "database/sql"

func migrate(sqlDB *sql.DB) error {
	if _, err := sqlDB.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return err
	}
	if _, err := sqlDB.Exec(schema); err != nil {
		return err
	}
	return seed(sqlDB)
}

func seed(sqlDB *sql.DB) error {
	var n int
	if err := sqlDB.QueryRow("SELECT COUNT(1) FROM users").Scan(&n); err != nil {
		return err
	}
	if n > 0 {
		return nil
	}
	_, err := sqlDB.Exec("INSERT INTO users(name, token_hash, role, created_at) VALUES('admin','bootstrap','admin', datetime('now'))")
	if err != nil {
		return err
	}
	_, err = sqlDB.Exec("INSERT INTO settings(k,v) VALUES('app_name','GridPulse Relay')")
	return err
}
