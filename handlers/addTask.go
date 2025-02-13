package handlers

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
)

// Defining the type of task
type Task struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Status      bool   `json:"status"`
	Description string `json:"description"`
}

func AddTask(w http.ResponseWriter, r *http.Request) {
	// Ensure the request is a POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode the JSON body into a Person struct
	var task Task
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&task); err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	tasks := []Task{task}
	if err := writeTasksToCSV("./tasks.csv", tasks); err != nil {
		log.Println("Failed write to the CSV file; task not added")
		return
	}

	// Process the data (for example, log it)
	log.Printf("Added task: %+v\n", task)

	// Respond back to the client
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Data received successfully"))
}

// WriteTasksToCSV writes a slice of Task objects to a CSV file with the given filename.
func writeTasksToCSV(filename string, tasks []Task) error {
	// Create or truncate the CSV file.
	file, err := os.Create(filename)
	os.
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a new CSV writer.
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write the header row.
	// header := []string{"ID", "Name", "Status", "Description"}
	// if err := writer.Write(header); err != nil {
	// 	return err
	// }

	// Write each task as a CSV record.
	for _, task := range tasks {
		record := []string{
			strconv.Itoa(task.ID),           // Convert int to string.
			task.Name,                       // Already a string.
			strconv.FormatBool(task.Status), // Convert bool to string.
			task.Description,                // Already a string.
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	// Check if any error occurred during writing.
	return writer.Error()
}
