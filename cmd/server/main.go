package main

import (
	"fmt"
	"go-graphql-sample/internal/database"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})
	if err := database.InitDB(); err != nil {
		log.Fatal(err)
	}
	defer database.CloseDB()

	fmt.Println("Database connection successful!")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
