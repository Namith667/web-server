package main

import (
	"context"
	"encoding/json"

	//"fmt"
	//"log"
	"net/http"
	"os"

	"web-server/logger"
	"web-server/model"
	"web-server/service"

	"github.com/joho/godotenv"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

var logg *logger.Logger

func init() {
	logg = logger.Init()
}

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		logg.Error("Error loading .env file", zap.Error(err))
		return
	}

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		logg.Error("Set the 'MONGODB_URI' environment variable.")
		return
	}

	// Create database connection
	db, err := NewDatabase(uri)
	if err != nil {
		logg.Warn("Error Connecting to Database", zap.Error(err))
		return
	}
	defer db.Disconnect()

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetUsers(w, db)
		case http.MethodPost:
			handleCreateUser(w, r, db)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/users/update", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			handleUpdateUser(w, r, db)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/users/delete", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			handleDeleteUser(w, r, db)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	logg.Info("Server started on port 8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		logg.Error("Failed to start server", zap.Error(err))
	}
}

func NewDatabase(uri string) (*service.Database, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	return &service.Database{Client: client}, nil
}

// Handlers
func handleGetUsers(w http.ResponseWriter, db *service.Database) {
	users, err := db.Read("users")
	if err != nil {
		logg.Error("Error reading users", zap.Error(err))
		http.Error(w, "Error reading users", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func handleCreateUser(w http.ResponseWriter, r *http.Request, db *service.Database) {
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		logg.Warn("Error decoding JSON", zap.Error(err))
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}
	result, err := db.Create("users", user)
	if err != nil {
		logg.Error("Error creating user", zap.Error(err))
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func handleUpdateUser(w http.ResponseWriter, r *http.Request, db *service.Database) {
	var update model.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		logg.Warn("Error decoding JSON", zap.Error(err))
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	// Convert filter to primitive.D
	filter := bson.D{}
	for key, value := range update.Filter {
		filter = append(filter, bson.E{Key: key, Value: value})
	}

	// Create the update document with $set
	updateDoc := bson.D{{"$set", bson.D{}}}

	// Append fields to the $set operator
	for key, value := range update.Update {
		updateDoc[0].Value = append(updateDoc[0].Value.(bson.D), bson.E{Key: key, Value: value})
	}

	result, err := db.Update("users", filter, updateDoc)
	if err != nil {
		logg.Error("Error updating user", zap.Error(err))
		http.Error(w, "Error updating user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result) // Send back the result of the update
	logg.Info("User updated successfully", zap.Int64("matchedCount", result.MatchedCount), zap.Int64("modifiedCount", result.ModifiedCount))
}
func handleDeleteUser(w http.ResponseWriter, r *http.Request, db *service.Database) {
	var filter model.FilterRequest
	if err := json.NewDecoder(r.Body).Decode(&filter); err != nil {
		logg.Warn("Error decoding JSON", zap.Error(err))
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	// Convert filter to primitive.D
	deleteFilter := bson.D{}
	for key, value := range filter.Filter {
		deleteFilter = append(deleteFilter, bson.E{Key: key, Value: value})
	}

	result, err := db.Delete("users", deleteFilter)
	if err != nil {
		logg.Error("Error deleting user", zap.Error(err))
		http.Error(w, "Error deleting user", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
