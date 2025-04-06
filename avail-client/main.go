package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/dhananjaypai08/availDAMCP/avail"
	"github.com/gorilla/mux"
)

type AvailData struct {
	AppId   uint32 `json:"AppId"`
	Message string `json:"Message"`
}

var txns = []AvailData{}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/check", getItems).Methods("GET")
	r.HandleFunc("/send-data", sendData).Methods("POST")

	fmt.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func getItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(txns)
}

func sendData(w http.ResponseWriter, r *http.Request) {
	var DataBlob AvailData
	err := json.NewDecoder(r.Body).Decode(&DataBlob)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	msg, err := avail.SendDataToDA(DataBlob.AppId, DataBlob.Message)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	txns = append(txns, DataBlob)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(msg)
}
