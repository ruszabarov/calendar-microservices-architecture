package main

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

func GetCalendars(w http.ResponseWriter, r *http.Request) {
	collection := database.Collection("calendars")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, "Error fetching calendarSummaries", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var calendarSummaries []CalendarSummary
	if err = cursor.All(ctx, &calendarSummaries); err != nil {
		http.Error(w, "Error decoding calendarSummaries", http.StatusInternalServerError)
		return
	}

	log.Println(calendarSummaries)

	var calendars []Calendar
	for _, calendarSummary := range calendarSummaries {
		meetings := GetMeetingsByIds(calendarSummary.Meetings)
		calendar := Calendar{
			ID:       calendarSummary.ID,
			Title:    calendarSummary.Title,
			Details:  calendarSummary.Details,
			Meetings: meetings,
		}
		calendars = append(calendars, calendar)
	}

	responseWrapper := map[string]interface{}{
		"data": calendars,
	}

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

	log.Println("Success: GetCalendars")
}

func GetCalendarsByIds(w http.ResponseWriter, r *http.Request) {
	idsParam := r.URL.Query().Get("ids")
	if idsParam == "" {
		http.Error(w, "Missing 'ids' query parameter", http.StatusBadRequest)
		return
	}

	ids := strings.Split(idsParam, ",")
	if len(ids) == 0 {
		http.Error(w, "No IDs provided", http.StatusBadRequest)
		return
	}

	collection := database.Collection("calendars")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"_id": bson.M{
			"$in": ids,
		},
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		http.Error(w, "Error fetching calendars", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var calendars []Calendar
	if err = cursor.All(ctx, &calendars); err != nil {
		http.Error(w, "Error decoding calendars", http.StatusInternalServerError)
		return
	}

	responseWrapper := map[string]interface{}{
		"data": calendars,
	}

	responseData, err := json.Marshal(responseWrapper)
	if err != nil {
		http.Error(w, "Error marshalling response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseData)
}

func UpdateCalendarById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var updatedCalendar Calendar
	err := json.NewDecoder(r.Body).Decode(&updatedCalendar)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	updatedCalendar.ID = ""

	collection := database.Collection("calendars")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{
		"$set": updatedCalendar,
	}

	result, err := collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		http.Error(w, "Error updating calendar", http.StatusInternalServerError)
		return
	}

	if result.MatchedCount == 0 {
		http.Error(w, "Calendar not found", http.StatusNotFound)
		return
	}

	responseWrapper := map[string]interface{}{
		"data": update,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseWrapper)
}

func DeleteCalendarById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	collection := database.Collection("calendars")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		http.Error(w, "Error deleting calendar", http.StatusInternalServerError)
		return
	}

	if result.DeletedCount == 0 {
		http.Error(w, "Calendar not found", http.StatusNotFound)
		return
	}

	responseWrapper := map[string]interface{}{
		"message": "Calendar deleted successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseWrapper)
}

func CreateCalendar(w http.ResponseWriter, r *http.Request) {
	var calendar CalendarSummary
	err := json.NewDecoder(r.Body).Decode(&calendar)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if calendar.ID == "" {
		calendar.ID = uuid.New().String()
	}

	collection := database.Collection("calendars")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = collection.InsertOne(ctx, calendar)
	if err != nil {
		http.Error(w, "Error inserting calendar", http.StatusInternalServerError)
		return
	}

	responseWrapper := map[string]interface{}{
		"data": calendar,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseWrapper)
}

func AddMeetingsToCalendar(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "Missing 'id' parameter", http.StatusBadRequest)
		return
	}

	var meetingsToAdd []string
	err := json.NewDecoder(r.Body).Decode(&meetingsToAdd)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if len(meetingsToAdd) == 0 {
		http.Error(w, "No meetings provided to add", http.StatusBadRequest)
		return
	}

	for _, meetingID := range meetingsToAdd {
		_, err := uuid.Parse(meetingID)
		if err != nil {
			http.Error(w, "Invalid meeting ID: "+meetingID, http.StatusBadRequest)
			return
		}
	}

	collection := database.Collection("calendars")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{
		"$addToSet": bson.M{
			"meetings": bson.M{
				"$each": meetingsToAdd,
			},
		},
	}

	result, err := collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		http.Error(w, "Error updating calendar", http.StatusInternalServerError)
		return
	}

	if result.MatchedCount == 0 {
		http.Error(w, "Calendar not found", http.StatusNotFound)
		return
	}

	responseWrapper := map[string]interface{}{
		"data": update,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseWrapper)
}
