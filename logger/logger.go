package logger

import (
	"fmt"
	"log"
	"os"
)

// LogFile holds basic information about the programs logfile.
type LogFile struct {
	ErrorCount int
	FilePath   string
}

// NewLogFile creates a reference to the program's logfile.
func NewLogFile(filePath string) *LogFile {
	return &LogFile{0, filePath}
}

// CreateLog creates a basic CSV logfile with a header row.
func (logger *LogFile) CreateLog() {
	// TODO: Too verbose?
	// Create logfile in current directory.
	file, err := os.Create(logger.FilePath)
	if err != nil {
		log.Fatalf("Could not create error file in current directory: %s", err)
	}
	// Ensure file is closed at end.
	defer file.Close()

	_, err = file.WriteString("error,row,id,sku,handle,image_url,reason\n")
	if err != nil {
		log.Printf("Error writing error file header: %s", err)
	}
}

// WriteEntry takes a RowError struct and writes it to the CSV logfile.
func (logger *LogFile) WriteEntry(entry RowError) {
	logger.ErrorCount++

	// Open existing log file.
	file, err := os.OpenFile(logger.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Ironic error in writing error to error file: %s", err)
	}

	// Write RowError to CSV file.
	_, err = file.WriteString(fmt.Sprintf("%s,%d,%s,%s,%s,%s,%s\n", entry.Error,
		entry.Row, entry.ID, entry.SKU, entry.Handle, entry.ImageURL, entry.Reason))
	if err != nil {
		log.Printf("Error writing entry to CSV error file: %s", err)
	}
}
