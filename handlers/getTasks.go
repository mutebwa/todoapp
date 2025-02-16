package handlers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

// GetTasks reads tasks from CSV and returns them as JSON
func GetTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := ReadTasksFromCSV("./tasks.csv")
	if err != nil {
		log.Printf("Error reading tasks: %v", err) // Log detailed error
		http.Error(w, "Failed to retrieve tasks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		log.Printf("JSON encoding error: %v", err)
		http.Error(w, "Failed to format response", http.StatusInternalServerError)
	}
}

// ReadTasksFromCSV reads and validates tasks from CSV
func ReadTasksFromCSV(filename string) ([]Task, error) {
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) { // Handle missing file gracefully
			return []Task{}, nil
		}
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading records: %w", err)
	}

	if len(records) == 0 {
		return []Task{}, nil // Empty file (only header case handled later)
	}

	// Validate CSV header
	expectedHeader := []string{"ID", "Name", "Status", "Description"}
	if !sliceEqual(records[0], expectedHeader) {
		return nil, fmt.Errorf("invalid CSV header. Got %v, expected %v", records[0], expectedHeader)
	}

	var tasks []Task
	for i, record := range records[1:] {
		if len(record) != 4 {
			return nil, fmt.Errorf("row %d: invalid field count (%d), expected 4", i+2, len(record))
		}

		id, err := strconv.Atoi(record[0])
		if err != nil {
			return nil, fmt.Errorf("row %d: invalid ID format: %w", i+2, err)
		}

		status, err := strconv.ParseBool(record[2])
		if err != nil {
			return nil, fmt.Errorf("row %d: invalid status format: %w", i+2, err)
		}

		tasks = append(tasks, Task{
			ID:          id,
			Name:        record[1],
			Status:      status,
			Description: record[3],
		})
	}
	return tasks, nil
}

// Helper function to compare string slices
func sliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}