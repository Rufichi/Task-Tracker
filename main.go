// main.go
package main

import (
    "fmt"
    "os"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: task-cli <command> [args...]")
        return
    }
    
    command := os.Args[1]
    fmt.Printf("Command: %s\n", command)
}
