document.addEventListener("DOMContentLoaded", () => {
    // Load tasks when the page loads
    loadTasks();
  
    // Attach event listener for the Add Task form
    const addTaskForm = document.getElementById("add-task-form");
    addTaskForm.addEventListener("submit", (e) => {
      e.preventDefault();
      addTask();
    });
  });
  
  // Fetch and display tasks from the server
  function loadTasks() {
    fetch("/tasks")
      .then((response) => {
        if (!response.ok) {
          throw new Error("Network response was not ok");
        }
        return response.json();
      })
      .then((tasks) => displayTasks(tasks))
      .catch((error) => console.error("Error fetching tasks:", error));
  }
  
  // Render tasks in a table
  function displayTasks(tasks) {
    const container = document.getElementById("tasks-container");
    container.innerHTML = ""; // Clear previous content
  
    if (tasks.length === 0) {
      container.innerHTML = "<p>No tasks found.</p>";
      return;
    }
  
    // Create a table for tasks
    const table = document.createElement("table");
  
    // Create header row
    const headerRow = document.createElement("tr");
    const headers = ["ID", "Name", "Status", "Description"];
    headers.forEach((headerText) => {
      const th = document.createElement("th");
      th.textContent = headerText;
      headerRow.appendChild(th);
    });
    table.appendChild(headerRow);
  
    // Create a row for each task
    tasks.forEach((task) => {
      const row = document.createElement("tr");
  
      const idCell = document.createElement("td");
      idCell.textContent = task.id;
      row.appendChild(idCell);
  
      const nameCell = document.createElement("td");
      nameCell.textContent = task.name;
      row.appendChild(nameCell);
  
      const statusCell = document.createElement("td");
      // Display status as "Completed" or "Pending"
      statusCell.textContent = task.status ? "Completed" : "Pending";
      row.appendChild(statusCell);
  
      const descCell = document.createElement("td");
      descCell.textContent = task.description;
      row.appendChild(descCell);
  
      table.appendChild(row);
    });
  
    container.appendChild(table);
  }
  
  // Send a new task to the server
  function addTask() {
    const nameInput = document.getElementById("task-name");
    const statusInput = document.getElementById("task-status");
    const descInput = document.getElementById("task-desc");
  
    // Build the task object from form values
    const task = {
      name: nameInput.value,
      status: statusInput.checked,
      description: descInput.value,
    };
  
    fetch("/add", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(task),
    })
      .then((response) => {
        if (!response.ok) {
          throw new Error("Failed to add task");
        }
        return response.json(); // Assume server responds with JSON data
      })
      .then((data) => {
        // Clear the form fields
        nameInput.value = "";
        statusInput.checked = false;
        descInput.value = "";
  
        // Optionally log the added task
        console.log("Task added:", data);
        // Refresh the tasks list to show the new task
        loadTasks();
      })
      .catch((error) => console.error("Error adding task:", error));
  }
  