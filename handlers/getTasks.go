package handlers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

// GetTasks is an HTTP handler that reads tasks from a CSV file and returns them as JSON.
func GetTasks(w http.ResponseWriter, r *http.Request) {
	// Read tasks from the CSV file.
	tasks, err := ReadTasksFromCSV("./tasks.csv")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read tasks: %v", err), http.StatusInternalServerError)
		return
	}

	// Set the response header to application/json.
	w.Header().Set("Content-Type", "application/json")

	// Encode tasks as JSON and write to the response.
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode tasks: %v", err), http.StatusInternalServerError)
		return
	}
}

// ReadTasksFromCSV reads a CSV file with the provided filename and returns a slice of Task.
func ReadTasksFromCSV(filename string) ([]Task, error) {
	// Open the CSV file.
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("unable to open CSV file: %w", err)
	}
	defer file.Close()

	// Create a new CSV reader.
	reader := csv.NewReader(file)

	// Read all records from the CSV file.
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("unable to read CSV file: %w", err)
	}

	// Check that there is at least a header row.
	if len(records) < 1 {
		return nil, fmt.Errorf("CSV file is empty")
	}

	var tasks []Task

	// Start at index 1 to skip the header row.
	for i, record := range records[1:] {
		// Each record should have exactly 4 fields.
		if len(record) < 4 {
			return nil, fmt.Errorf("record %d is malformed: %v", i+1, record)
		}

		// Convert ID from string to int.
		id, err := strconv.Atoi(record[0])
		if err != nil {
			return nil, fmt.Errorf("error parsing ID on row %d: %w", i+2, err)
		}

		// Convert Status from string to bool.
		status, err := strconv.ParseBool(record[2])
		if err != nil {
			return nil, fmt.Errorf("error parsing Status on row %d: %w", i+2, err)
		}

		task := Task{
			ID:          id,
			Name:        record[1],
			Status:      status,
			Description: record[3],
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}
