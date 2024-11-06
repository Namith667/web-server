package controller

import (
	"encoding/json"
	"net/http"
	"web-server/logger"
	"web-server/model"
	"web-server/service"

	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

var logg *logger.Logger = logger.Init()

func HandleGetUsers(db *service.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := db.Read("users")
		if err != nil {
			logg.Error("Error reading users", zap.Error(err))
			http.Error(w, "Error reading users", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}

func HandleCreateUser(db *service.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user model.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			logg.Error("Error decoding JSON", zap.Error(err))
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
}

func HandleUpdateUser(db *service.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var update model.UpdateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
			logg.Error("Error decoding JSON", zap.Error(err))
			http.Error(w, "Error decoding JSON", http.StatusBadRequest)
			return
		}

		filter := bson.D{}
		for key, value := range update.Filter {
			filter = append(filter, bson.E{Key: key, Value: value})
		}

		updateDoc := bson.D{{"$set", bson.D{}}}
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
		json.NewEncoder(w).Encode(result)
		logg.Info("User updated successfully", zap.Int64("matchedCount", result.MatchedCount), zap.Int64("modifiedCount", result.ModifiedCount))
	}
}

func HandleDeleteUser(db *service.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var filter model.FilterRequest
		if err := json.NewDecoder(r.Body).Decode(&filter); err != nil {
			logg.Error("Error decoding JSON", zap.Error(err))
			http.Error(w, "Error decoding JSON", http.StatusBadRequest)
			return
		}

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
}
