// All Global Constants and Program Default Values
package internal

// DefaultHTTPPort defines the Default HTTP Port to listen
const DefaultHTTPPort int = 8080

// DefaultComputerFilePath defines the Default File Path for the SQLite database
const DefaultComputerFilePath string = "computer.db"

// DefaultComputerFilePathEnvironmentName defines the Name of the Environment Variable where we should look for a DB path
const DefaultComputerFilePathEnvironmentName string = "WOLDB"

// DefaultHTTPPortEnvironmentVariableName defines the Name of the Environment Variable what TCP Port we should use for the Webserver
const DefaultHTTPPortEnvironmentVariableName string = "WOLHTTPPORT"
