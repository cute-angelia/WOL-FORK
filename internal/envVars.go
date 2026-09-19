package internal

import (
	"log"
	"os"
	"strconv"
)

// ProcessEnvVars processes environment variables for port and DB path
func ProcessEnvVars(port int, dbFile string) (int, string) {
	if os.Getenv(DefaultComputerFilePathEnvironmentName) != "" {
		dbFile = os.Getenv(DefaultComputerFilePathEnvironmentName)
	}

	if os.Getenv(DefaultHTTPPortEnvironmentVariableName) != "" {
		var err error
		if port, err = strconv.Atoi(os.Getenv(DefaultHTTPPortEnvironmentVariableName)); err != nil {
			log.Fatalf("Environment variable \"%s\" should be an integer", DefaultHTTPPortEnvironmentVariableName)
		}
	}
	return port, dbFile
}
