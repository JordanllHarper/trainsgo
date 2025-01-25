package api

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

func handleResponse(
	db *gorm.DB,
	r *http.Request,
	w http.ResponseWriter,
	handler func(db *gorm.DB, r *http.Request) (ResponseBody, HttpError),
) error {
	response, httpError := handler(db, r)
	if httpError != nil {
		handleHttpError(w, httpError)
		return nil
	}

	json, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Operation was successful but response couldn't be serialized", http.StatusInternalServerError)
		return err
	}

	w.WriteHeader(response.StatusCode())
	w.Write(json)

	fmt.Printf("Responded to %v request with %v\n", r.Method, response)

	return nil
}

func handleTrains(w http.ResponseWriter, r *http.Request, db *gorm.DB) error {
	method := r.Method

	switch method {
	case "GET":
		return handleResponse(db, r, w, onTrainGet)

	case "POST":
		return handleResponse(db, r, w, onTrainPost)

	case "PUT":
		return handleResponse(db, r, w, onTrainPut)

	case "DELETE":
		return handleResponse(db, r, w, onTrainDelete)

	default:
		http.Error(w, "Supported methods are: GET POST PUT DELETE", http.StatusMethodNotAllowed)
	}

	return nil
}

func handleStations(w http.ResponseWriter, r *http.Request, db *gorm.DB) error {

	method := r.Method

	switch method {
	case "GET":
		return handleResponse(db, r, w, onStationGet)

	case "POST":
		return handleResponse(db, r, w, onStationPost)

	case "PUT":
		return handleResponse(db, r, w, onStationPut)

	case "DELETE":
		return handleResponse(db, r, w, onStationDelete)

	default:
		http.Error(w, "Supported methods are: GET POST PUT DELETE", http.StatusMethodNotAllowed)
	}

	return nil
}

func Run() {
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
		// http.HandleFunc("/", )
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

		http.HandleFunc("/stations", func(w http.ResponseWriter, r *http.Request) {
			err := handleStations(
				w,
				r,
				db,
			)
			if err != nil {
				log.Println("[ERROR] handle station failed:", err)
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
