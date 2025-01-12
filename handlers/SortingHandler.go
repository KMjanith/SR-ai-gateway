package handlers

import (
	"api-gateway/producers"
	"api-gateway/spec"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"google.golang.org/protobuf/proto"
)

type SortingJob struct {
	Method string  `json:"method"`
	Arr    []int64 `json:"arr"`
}

func SortingHandler(w http.ResponseWriter, r *http.Request, ctx context.Context) {

	log.Println("received sorting request")

	if r.Method != http.MethodPost {
		log.Println("Invalid request method")
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var sortingJob SortingJob
	// Decode the JSON body into the User struct
	err := json.NewDecoder(r.Body).Decode(&sortingJob)
	if err != nil {
		log.Println("Invalid JSON data")
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	// Send register request to RabbitMQ
	sortingResponse := producers.SendSortingRequest(sortingJob.Method, sortingJob.Arr, ctx)

	var SortingResponse spec.SoritingResponse
	// Unmarshal the Protobuf message
	err = proto.Unmarshal(sortingResponse, &SortingResponse)
	if err != nil {
		log.Fatalf("Failed to unmarshal Protobuf message: %v", err)
	}
	
	//Send the response back to the client
	if sortingResponse != nil {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{"message": "Registration successful!", "Method": SortingResponse.Method, "array": SortingResponse.SortedArr, "time": SortingResponse.Time})

	} else {
		http.Error(w, "Registration failed", http.StatusInternalServerError)
		w.WriteHeader(http.StatusFailedDependency)
		json.NewEncoder(w).Encode(map[string]string{"message": "Registration failed"})
	}

}
