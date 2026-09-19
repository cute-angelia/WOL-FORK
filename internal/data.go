package internal

import (
	"database/sql"
	"log"
	"net"
	"strings"
	"time"

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

// CheckOnline 探测设备是否在线。
// 从广播 IP（如 192.168.3.180:9）中提取主机部分，
// 依次尝试常用端口，遇到"连接成功"或"连接被拒绝"（主机在线但端口关闭）即视为在线。
func CheckOnline(broadcastIP string) bool {
	if broadcastIP == "" {
		return false
	}
	// 提取纯 IP，去掉 :port
	host := broadcastIP
	if idx := strings.LastIndex(broadcastIP, ":"); idx != -1 {
		host = broadcastIP[:idx]
	}

	// Windows: 135(RPC), 445(SMB), 3389(RDP)
	// macOS:   22(SSH), 548(AFP), 5900(VNC)
	// 通用:    80(HTTP)
	ports := []string{"135", "445", "22", "80", "3389", "548"}
	timeout := 2 * time.Second

	for _, port := range ports {
		conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
		if err == nil {
			conn.Close()
			return true
		}
		if strings.Contains(err.Error(), "refused") ||
			strings.Contains(err.Error(), "reset") {
			return true // 主机在线，只是端口关闭
		}
	}
	return false
}

// FileExists kept for interface compatibility
func FileExists(name string) bool {
	return false
}
