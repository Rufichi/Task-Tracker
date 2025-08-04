// main.go 

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

// Task rappresenta la struttura di un singolo compito
type Task struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Status      string    `json:"status"` // "todo", "in-progress", "done"
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// Global variable per il percorso del file JSON
const tasksFilePath = "tasks.json"

// Funzione per caricare i compiti dal file JSON
func loadTasks() ([]Task, error) {
	_, err := os.Stat(tasksFilePath)
	if os.IsNotExist(err) {
		// Se il file non esiste, ritorna una lista vuota e nessun errore
		return []Task{}, nil
	}

	data, err := ioutil.ReadFile(tasksFilePath)
	if err != nil {
		return nil, fmt.Errorf("errore durante la lettura del file: %w", err)
	}

	if len(data) == 0 {
		return []Task{}, nil // File vuoto, ritorna lista vuota
	}

	var tasks []Task
	err = json.Unmarshal(data, &tasks)
	if err != nil {
		return nil, fmt.Errorf("errore durante il parsing del JSON: %w", err)
	}
	return tasks, nil
}

// Funzione per salvare i compiti nel file JSON
func saveTasks(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ") // Indent per leggibilità
	if err != nil {
		return fmt.Errorf("errore durante la serializzazione del JSON: %w", err)
	}
	err = ioutil.WriteFile(tasksFilePath, data, 0644) // Permessi di lettura/scrittura
	if err != nil {
		return fmt.Errorf("errore durante la scrittura del file: %w", err)
	}
	return nil
}

// Funzione per trovare il prossimo ID disponibile
func getNextID(tasks []Task) int {
	if len(tasks) == 0 {
		return 1
	}
	maxID := 0
	for _, task := range tasks {
		if task.ID > maxID {
			maxID = task.ID
		}
	}
	return maxID + 1
}

// Funzione per aggiungere un nuovo compito
func addTask(description string) {
	if description == "" {
		fmt.Println("Errore: La descrizione del compito non può essere vuota.")
		return
	}

	tasks, err := loadTasks()
	if err != nil {
		fmt.Println("Errore durante il caricamento dei compiti:", err)
		return
	}

	newID := getNextID(tasks)
	newTask := Task{
		ID:          newID,
		Description: description,
		Status:      "todo", // Stato iniziale predefinito
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	tasks = append(tasks, newTask) // Aggiungi il nuovo compito alla lista

	if err := saveTasks(tasks); err != nil {
		fmt.Println("Errore durante il salvataggio dei compiti:", err)
		return
	}

	fmt.Printf("Compito aggiunto con successo (ID: %d)\n", newID)
}

// Funzione per aggiornare un compito esistente
func updateTask(id int, newDescription string) {
	tasks, err := loadTasks()
	if err != nil {
		fmt.Println("Errore durante il caricamento dei compiti:", err)
		return
	}

	found := false
	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Description = newDescription
			tasks[i].UpdatedAt = time.Now()
			found = true
			break
		}
	}

	if !found {
		fmt.Printf("Errore: Compito con ID %d non trovato.\n", id)
		return
	}

	if err := saveTasks(tasks); err != nil {
		fmt.Println("Errore durante il salvataggio dei compiti:", err)
		return
	}

	fmt.Printf("Compito con ID %d aggiornato con successo.\n", id)
}

// Funzione per eliminare un compito
func deleteTask(id int) {
	tasks, err := loadTasks()
	if err != nil {
		fmt.Println("Errore durante il caricamento dei compiti:", err)
		return
	}

	newTasks := []Task{}
	found := false
	for _, task := range tasks {
		if task.ID == id {
			found = true
			continue // Salta il compito da eliminare
		}
		newTasks = append(newTasks, task)
	}

	if !found {
		fmt.Printf("Errore: Compito con ID %d non trovato.\n", id)
		return
	}

	if err := saveTasks(newTasks); err != nil {
		fmt.Println("Errore durante il salvataggio dei compiti:", err)
		return
	}

	fmt.Printf("Compito con ID %d eliminato con successo.\n", id)
}



// Funzione per elencare i compiti
func listTasks(filter string) {
	tasks, err := loadTasks()
	if err != nil {
		fmt.Println("Errore durante il caricamento dei compiti:", err)
		return
	}

	if len(tasks) == 0 {
		fmt.Println("Nessun compito trovato.")
		return
	}

	fmt.Println("--- Elenco Compiti ---")
	for _, task := range tasks {
		display := true
		switch filter {
		case "done":
			if task.Status != "done" {
				display = false
			}
		case "todo":
			if task.Status != "todo" {
				display = false
			}
		case "in-progress":
			if task.Status != "in-progress" {
				display = false
			}
		case "": // Nessun filtro, mostra tutto
			display = true
		default:
			fmt.Printf("Filtro '%s' non valido. Usa 'done', 'todo', 'in-progress' o nessun filtro.\n", filter)
			return
		}

		if display {
			fmt.Printf("ID: %d | Stato: %-12s | Descrizione: %s (Creato: %s, Aggiornato: %s)\n",
				task.ID,
				strings.ToUpper(task.Status), // Rende lo stato in maiuscolo per leggibilità
				task.Description,
				task.CreatedAt.Format("02 Jan 15:04"),
				task.UpdatedAt.Format("02 Jan 15:04"),
			)
		}
	}
	fmt.Println("----------------------")
}

// Funzione per cambiare lo stato di un compito
func markTask(id int, status string) {
	if status != "in-progress" && status != "done" && status != "todo" { // Aggiunto "todo" per completezza, anche se non richiesto esplicitamente dai comandi
		fmt.Println("Errore: Stato non valido. Usa 'in-progress' o 'done'.")
		return
	}

	tasks, err := loadTasks()
	if err != nil {
		fmt.Println("Errore durante il caricamento dei compiti:", err)
		return
	}

	found := false
	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Status = status
			tasks[i].UpdatedAt = time.Now()
			found = true
			break
		}
	}

	if !found {
		fmt.Printf("Errore: Compito con ID %d non trovato.\n", id)
		return
	}

	if err := saveTasks(tasks); err != nil {
		fmt.Println("Errore durante il salvataggio dei compiti:", err)
		return
	}

	fmt.Printf("Compito con ID %d marcato come '%s' con successo.\n", id, status)
}

// Funzione principale che analizza gli argomenti della riga di comando
func main() {
	args := os.Args[1:] // os.Args[0] è il nome del programma stesso, quindi iniziamo da 1

	if len(args) == 0 {
		fmt.Println("Utilizzo: task-cli [comando] [argomenti]")
		fmt.Println("Comandi disponibili:")
		fmt.Println("  add \"descrizione\"              - Aggiunge un nuovo compito")
		fmt.Println("  update <ID> \"nuova descrizione\" - Aggiorna la descrizione di un compito")
		fmt.Println("  delete <ID>                    - Elimina un compito")
		fmt.Println("  mark-in-progress <ID>          - Marca un compito come 'in corso'")
		fmt.Println("  mark-done <ID>                 - Marca un compito come 'completato'")
		fmt.Println("  list [done|todo|in-progress]   - Elenca tutti i compiti o quelli filtrati")
		return
	}

	command := args[0]

	switch command {
	case "add":
		if len(args) < 2 {
			fmt.Println("Errore: 'add' richiede una descrizione. Esempio: task-cli add \"Comprare pane\"")
			return
		}
		addTask(strings.Join(args[1:], " ")) // Unisci tutti gli argomenti rimanenti come descrizione
	case "update":
		if len(args) < 3 {
			fmt.Println("Errore: 'update' richiede un ID e una nuova descrizione. Esempio: task-cli update 1 \"Finire il report\"")
			return
		}
		id, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Errore: L'ID deve essere un numero intero valido.")
			return
		}
		updateTask(id, strings.Join(args[2:], " "))
	case "delete":
		if len(args) < 2 {
			fmt.Println("Errore: 'delete' richiede un ID. Esempio: task-cli delete 1")
			return
		}
		id, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Errore: L'ID deve essere un numero intero valido.")
			return
		}
		deleteTask(id)
	case "mark-in-progress":
		if len(args) < 2 {
			fmt.Println("Errore: 'mark-in-progress' richiede un ID. Esempio: task-cli mark-in-progress 1")
			return
		}
		id, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Errore: L'ID deve essere un numero intero valido.")
			return
		}
		markTask(id, "in-progress")
	case "mark-done":
		if len(args) < 2 {
			fmt.Println("Errore: 'mark-done' richiede un ID. Esempio: task-cli mark-done 1")
			return
		}
		id, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Errore: L'ID deve essere un numero intero valido.")
			return
		}
		markTask(id, "done")
	case "list":
		filter := ""
		if len(args) >= 2 {
			filter = args[1]
		}
		listTasks(filter)
	default:
		fmt.Printf("Comando '%s' non riconosciuto.\n", command)
		fmt.Println("Digita 'task-cli' senza argomenti per vedere l'utilizzo.")
	}
}
