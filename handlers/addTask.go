package handlers

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Task struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Status      bool   `json:"status"`
	Description string `json:"description"`
}

func AddTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var task Task
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&task); err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	tasks := []Task{task}
	if err := writeTasksToCSV("./tasks.csv", tasks); err != nil {
		log.Println("Failed to write to CSV:", err)
		http.Error(w, "Failed to save task", http.StatusInternalServerError)
		return
	}

	log.Printf("Added task: %+v\n", task)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Data received successfully"))
}

func writeTasksToCSV(filename string, tasks []Task) error {
	// Open file in append mode, create if not exists
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Check if file is empty to write header
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header only if file is new/empty
	if fileInfo.Size() == 0 {
		header := []string{"ID", "Name", "Status", "Description"}
		if err := writer.Write(header); err != nil {
			return err
		}
	}

	// Write task records
	for _, task := range tasks {
		record := []string{
			strconv.Itoa(task.ID),
			task.Name,
			strconv.FormatBool(task.Status),
			task.Description,
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return writer.Error()
}