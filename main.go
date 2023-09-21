package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Failed to load .env file:", err)
		return
	}

	// dbUser := os.Getenv("DBUSER")
	// dbPass := os.Getenv("DBPASS")
	// dbName := os.Getenv("DBNAME")
	// db, err := sqlx.Connect("postgres", fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", dbUser, dbPass, dbName))
	// if err != nil {
	// 	fmt.Println("Failed to connect to database:", err)
	// 	return
	// }

	mainRouter := mux.NewRouter()
	// test route
	mainRouter.HandleFunc("/", func(writer http.ResponseWriter, reader *http.Request) {
		// send hello world as json
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(map[string]string{"message": "Hello World!"})
	})

	port := os.Getenv("PORT")
	log.Fatal(http.ListenAndServe(port, mainRouter))
}
