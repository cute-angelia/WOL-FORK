/* Processing and handling Shell Arguments */

package internal

import (
	"flag"
)

// ProcessShellArgs processes command-line arguments
func ProcessShellArgs() (int, string) {
	port := flag.Int("port", DefaultHTTPPort,
		"Define the port on which the webserver will listen to")
	dbPath := flag.String("db", DefaultComputerFilePath,
		"Path to the SQLite database file")

	flag.Parse()

	return *port, *dbPath
}
