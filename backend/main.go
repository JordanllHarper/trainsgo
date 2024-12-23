package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func handleHttpError(w http.ResponseWriter, err HttpError) {
	code, msg := err.Status()
	http.Error(w, msg, code)
}

// Crud operations on trains
func handleTrains(w http.ResponseWriter, r *http.Request, db *gorm.DB) error {
	method := r.Method
	encoder := json.NewEncoder(w)

	switch method {
	case "GET":
		response, httpError := onTrainGet(db, r)
		if httpError != nil {
			handleHttpError(w, httpError)
			return nil
		}
		if err := encoder.Encode(response); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return err
		} else {
			fmt.Println("Responded to GET request with", response)
		}

	case "POST":
		response, err := onTrainPost(db, r)
		if err != nil {
			code, msg := err.Status()
			http.Error(w, msg, code)
			return nil
		}

		if err := encoder.Encode(response); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return err
		} else {
			fmt.Println("Responded to POST request with", response)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	return nil
}

func main() {
	// db
	dsn := "root:@tcp(127.0.0.1:3306)/trainsgo?charset=utf8mb4&parseTime=True"
	db, err := gorm.Open(
		mysql.Open(dsn),
		&gorm.Config{},
	)

	if err != nil {
		log.Fatalf("Database error %s", err)
		return // for explicitness
	}

	{
		http.HandleFunc("/", handleRoot)
		http.HandleFunc("/trains", func(w http.ResponseWriter, r *http.Request) {
			err := handleTrains(
				w,
				r,
				db,
			)
			if err != nil {
				log.Println("[ERROR] handle trains failed:", err)
			}
		})
	}

	{
		address := "127.0.0.1"
		port := "3333"
		listening := fmt.Sprintf("%v:%v", address, port)

		fmt.Printf("Starting server. Listening on port: %v\n", listening)
		err = http.ListenAndServe(listening, nil)
	}

	fmt.Printf("%v\n", err)

}
