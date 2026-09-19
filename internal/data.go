package internal

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

// DB is the global SQLite database connection
var DB *sql.DB

// InitDB initializes the SQLite database and creates the table if it doesn't exist
func InitDB(dbPath string) error {
	var err error
	DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS computers (
		id   INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT    NOT NULL UNIQUE,
		mac  TEXT    NOT NULL UNIQUE,
		ip   TEXT    NOT NULL UNIQUE
	);`

	if _, err = DB.Exec(createTableSQL); err != nil {
		return err
	}

	log.Printf("SQLite database initialized at %s", dbPath)
	return nil
}

// LoadComputerList loads all computers from the database
func LoadComputerList() ([]Computer, error) {
	rows, err := DB.Query("SELECT id, name, mac, ip FROM computers ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var computers []Computer
	for rows.Next() {
		var c Computer
		if err := rows.Scan(&c.ID, &c.Name, &c.Mac, &c.BroadcastIPAddress); err != nil {
			return nil, err
		}
		computers = append(computers, c)
	}
	return computers, rows.Err()
}

// AddComputer inserts a new computer into the database
func AddComputer(c Computer) (Computer, error) {
	result, err := DB.Exec(
		"INSERT INTO computers (name, mac, ip) VALUES (?, ?, ?)",
		c.Name, c.Mac, c.BroadcastIPAddress,
	)
	if err != nil {
		return c, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return c, err
	}
	c.ID = id
	return c, nil
}

// UpdateComputer updates an existing computer identified by oldName
func UpdateComputer(oldName string, c Computer) error {
	_, err := DB.Exec(
		"UPDATE computers SET name=?, mac=?, ip=? WHERE name=?",
		c.Name, c.Mac, c.BroadcastIPAddress, oldName,
	)
	return err
}

// DeleteComputer removes a computer by name from the database
func DeleteComputer(name string) error {
	_, err := DB.Exec("DELETE FROM computers WHERE name=?", name)
	return err
}

// FileExists kept for interface compatibility; SQLite file is created automatically
func FileExists(name string) bool {
	return false
}
