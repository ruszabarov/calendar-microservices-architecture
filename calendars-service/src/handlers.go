package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo/options"

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

	var calendars []Calendar
	for _, calendarSummary := range calendarSummaries {
		calendars = append(calendars, ConvertSummaryToFull(calendarSummary))
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

func CreateCalendar(w http.ResponseWriter, r *http.Request) {
	var calendarSummary CalendarSummary
	err := json.NewDecoder(r.Body).Decode(&calendarSummary)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if calendarSummary.ID == "" {
		calendarSummary.ID = uuid.New().String()
	}

	collection := database.Collection("calendars")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = collection.InsertOne(ctx, calendarSummary)
	if err != nil {
		http.Error(w, "Error inserting calendarSummary", http.StatusInternalServerError)
		return
	}

	responseWrapper := map[string]interface{}{
		"data": ConvertSummaryToFull(calendarSummary),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseWrapper)
}

func UpdateCalendarById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var updatedCalendar Calendar
	err := json.NewDecoder(r.Body).Decode(&updatedCalendar)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	updatedCalendar.ID = ""

	collection := database.Collection("calendars")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"title":   updatedCalendar.Title,
			"details": updatedCalendar.Details,
		},
	}

	var calendarSummary CalendarSummary

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err = collection.FindOneAndUpdate(ctx, bson.M{"_id": id}, update, opts).Decode(&calendarSummary)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error fetching updated calendar")
		return
	}

	responseWrapper := map[string]interface{}{
		"data": ConvertSummaryToFull(calendarSummary),
	}

	respondWithJSON(w, http.StatusOK, responseWrapper)
}

func DeleteCalendarById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	collection := database.Collection("calendars")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error deleting calendar")
		return
	}

	if result.DeletedCount == 0 {
		http.Error(w, "Calendar not found", http.StatusNotFound)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "calendar successfully deleted"})
}

func AddMeetingsToCalendar(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Missing 'id' parameter")
		return
	}

	var request CalendarSummary
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if len(request.Meetings) == 0 {
		respondWithError(w, http.StatusBadRequest, "No meetings provided to add")
		return
	}

	for _, meetingID := range request.Meetings {
		_, err := uuid.Parse(meetingID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid meeting ID: "+meetingID)
			return
		}
	}

	collection := database.Collection("calendars")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{
		"$addToSet": bson.M{
			"meetings": bson.M{
				"$each": request.Meetings,
			},
		},
	}

	var calendarSummary CalendarSummary

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err = collection.FindOneAndUpdate(ctx, bson.M{"_id": id}, update, opts).Decode(&calendarSummary)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error fetching updated calendar")
		return
	}

	responseWrapper := map[string]interface{}{
		"data": ConvertSummaryToFull(calendarSummary),
	}

	respondWithJSON(w, http.StatusOK, responseWrapper)
}

func RemoveMeetingsFromCalendar(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Missing 'id' parameter")
		return
	}

	var meetingsToRemove CalendarSummary
	err := json.NewDecoder(r.Body).Decode(&meetingsToRemove)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if len(meetingsToRemove.Meetings) == 0 {
		respondWithError(w, http.StatusBadRequest, "No meetings provided to add")
		return
	}

	for _, meetingID := range meetingsToRemove.Meetings {
		_, err := uuid.Parse(meetingID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid meeting ID: "+meetingID)
			return
		}
	}

	collection := database.Collection("calendars")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{
		"$pull": bson.M{
			"meetings": bson.M{
				"$in": meetingsToRemove,
			},
		},
	}

	var calendarSummary CalendarSummary

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err = collection.FindOneAndUpdate(ctx, bson.M{"_id": id}, update, opts).Decode(&calendarSummary)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error fetching updated calendar")
		return
	}

	responseWrapper := map[string]interface{}{
		"data": ConvertSummaryToFull(calendarSummary),
	}

	respondWithJSON(w, http.StatusOK, responseWrapper)
}

func GetCalendarsByIds(w http.ResponseWriter, r *http.Request) {
	idsParam := r.URL.Query().Get("ids")
	if idsParam == "" {
		respondWithError(w, http.StatusBadRequest, "Missing 'ids' query parameter")
		return
	}

	ids := strings.Split(idsParam, ",")
	if len(ids) == 0 {
		respondWithError(w, http.StatusBadRequest, "No IDs provided")
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
		respondWithError(w, http.StatusInternalServerError, "Error fetching calendars")
		return
	}
	defer cursor.Close(ctx)

	var calendars []Calendar
	if err = cursor.All(ctx, &calendars); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding calendars")
		return
	}

	responseWrapper := map[string]interface{}{
		"data": calendars,
	}

	respondWithJSON(w, http.StatusOK, responseWrapper)
}

func ConvertSummaryToFull(calendarSummary CalendarSummary) Calendar {
	meetings := GetMeetingsByIds(calendarSummary.Meetings)
	var calendar = Calendar{
		ID:       calendarSummary.ID,
		Title:    calendarSummary.Title,
		Details:  calendarSummary.Details,
		Meetings: meetings,
	}

	return calendar
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}
