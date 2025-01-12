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

type User struct {
	Username       string `json:"username"`
	Password       string `json:"password"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request, ctx context.Context) {

	if r.Method != http.MethodPost {
		log.Println("Invalid request method")
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var user User
	// Decode the JSON body into the User struct
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println("Invalid JSON data")
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}


	// Send register request to RabbitMQ
	registerResponse:= producers.SendRegisterRequest(user.Username, user.Password, ctx)

	var RegisterResponse spec.RegisterUserResponse
	// Unmarshal the Protobuf message
	err = proto.Unmarshal(registerResponse, &RegisterResponse)
	if err != nil {
		log.Fatalf("Failed to unmarshal Protobuf message: %v", err)
	}

	log.Printf("token: %s" , RegisterResponse.Token)

	//Send the response back to the client
	if RegisterResponse.Token != "" {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Registration successful!", "token": RegisterResponse.Token})
		
	} else {
		http.Error(w, "Registration failed", http.StatusInternalServerError)
		w.WriteHeader(http.StatusFailedDependency)
		json.NewEncoder(w).Encode(map[string]string{"message": "Registration failed", "token": RegisterResponse.Token})
	}
}
