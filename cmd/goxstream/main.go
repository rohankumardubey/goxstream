package main

import (
    "fmt"
    "goxstream/internal/api"
)

func main() {
    fmt.Println("GoXStream REST API running on :8080")
    if err := api.StartAPIServer(":8080"); err != nil {
        panic(err)
    }
}
