package utils

import "database/sql"

// ---------------------- DB Initialization ----------------------
func InitDB(db *sql.DB) {
	_, _ = db.Exec(`CREATE TABLE IF NOT EXISTS debts(
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        customer TEXT,
        phone TEXT,
        amount REAL,
		balance REAL
    );`)

	_, _ = db.Exec(`CREATE TABLE IF NOT EXISTS payments(
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        customer TEXT,
        amount REAL,
        balance REAL,
		created_at REAL
    );`)
}
