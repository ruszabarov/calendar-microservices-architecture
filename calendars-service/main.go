package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

type Calendar struct {
	ID          string   `json:"id" bson:"_id,omitempty"`
	Title       string   `json:"title" bson:"title"`
	Description string   `json:"description" bson:"description"`
	Meetings    []string `json:"meetings" bson:"meetings"`
}

func GetCalendars(w http.ResponseWriter, r *http.Request) {
	collection := client.Database("store").Collection("calendars")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, "Error fetching calendars", http.StatusInternalServerError)
		return
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			http.Error(w, "Error closing cursor", http.StatusInternalServerError)
			return
		}
	}(cursor, ctx)

	var calendars []Calendar
	for cursor.Next(ctx) {
		var calendar Calendar
		if err := cursor.Decode(&calendar); err != nil {
			http.Error(w, "Error decoding calendar", http.StatusInternalServerError)
			return
		}
		calendars = append(calendars, calendar)
	}

	responseWrapper := map[string]interface{}{
		"data": calendars,
	}

	// Marshal the wrapped response to JSON
	responseData, err := json.Marshal(responseWrapper)
	if err != nil {
		http.Error(w, "Error marshalling response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(responseData)
	if err != nil {
		http.Error(w, "Error writing response", http.StatusInternalServerError)
		return
	}
}

func CreateCalendar(w http.ResponseWriter, r *http.Request) {
	var calendar Calendar
	err := json.NewDecoder(r.Body).Decode(&calendar)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	calendar.ID = uuid.New().String()

	collection := client.Database("store").Collection("calendars")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := collection.InsertOne(ctx, calendar)
	if err != nil {
		http.Error(w, "Error inserting calendar", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(result.InsertedID)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusBadRequest)
		return
	}
}

func connectToMongoDB() {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI environment variable not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB!")
}

func initializeRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/calendars", GetCalendars).Methods("GET")
	router.HandleFunc("/calendars", CreateCalendar).Methods("POST")
	return router
}

func main() {
	connectToMongoDB()

	router := initializeRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}
