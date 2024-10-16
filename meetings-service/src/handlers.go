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

func GetMeetings(w http.ResponseWriter, r *http.Request) {
	collection := database.Collection("meetings")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, "Error fetching meetingSummaries", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var meetingSummaries []MeetingSummary
	if err = cursor.All(ctx, &meetingSummaries); err != nil {
		http.Error(w, "Error decoding meetingSummaries", http.StatusInternalServerError)
		return
	}

	var meetings []Meeting
	for _, meetingSummary := range meetingSummaries {
		meetings = append(meetings, ConvertSummaryToFull(meetingSummary))
	}

	responseData, err := json.Marshal(meetings)
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

func CreateMeeting(w http.ResponseWriter, r *http.Request) {
	var meetingSummary MeetingSummary
	err := json.NewDecoder(r.Body).Decode(&meetingSummary)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if meetingSummary.ID == "" {
		meetingSummary.ID = uuid.New().String()
	}

	collection := database.Collection("meetings")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = collection.InsertOne(ctx, meetingSummary)
	if err != nil {
		http.Error(w, "Error inserting meetingSummary", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(meetingSummary)
}

func UpdateMeetingById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var updatedMeeting Meeting
	err := json.NewDecoder(r.Body).Decode(&updatedMeeting)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	updatedMeeting.ID = ""

	collection := database.Collection("meetings")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"title":   updatedMeeting.Title,
			"details": updatedMeeting.Details,
		},
	}

	var meetingSummary MeetingSummary

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err = collection.FindOneAndUpdate(ctx, bson.M{"_id": id}, update, opts).Decode(&meetingSummary)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error fetching updated meeting")
		return
	}

	respondWithJSON(w, http.StatusOK, meetingSummary)
}

func DeleteMeetingById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	collection := database.Collection("meetings")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error deleting meeting")
		return
	}

	if result.DeletedCount == 0 {
		http.Error(w, "Meeting not found", http.StatusNotFound)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "meeting successfully deleted"})
}

func AddCalendarToMeeting(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	meetingId := vars["meetingId"]
	calendarId := vars["calendarId"]
	if meetingId == "" || calendarId == "" {
		respondWithError(w, http.StatusBadRequest, "Missing 'id' parameter")
		return
	}

	collection := database.Collection("meetings")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{
		"$addToSet": bson.M{
			"calendars": calendarId,
		},
	}

	var meetingSummary MeetingSummary

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err := collection.FindOneAndUpdate(ctx, bson.M{"_id": meetingId}, update, opts).Decode(&meetingSummary)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error fetching updated meeting")
		return
	}

	respondWithJSON(w, http.StatusOK, ConvertSummaryToFull(meetingSummary))
}

func RemoveCalendarFromMeeting(w http.ResponseWriter, r *http.Request) {
	log.Println("hello world!")
	vars := mux.Vars(r)
	meetingId := vars["meetingId"]
	calendarId := vars["calendarId"]
	if meetingId == "" || calendarId == "" {
		respondWithError(w, http.StatusBadRequest, "Missing 'id' parameter")
		return
	}

	collection := database.Collection("meetings")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{
		"$pull": bson.M{
			"calendars": calendarId,
		},
	}

	var meetingSummary MeetingSummary

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err := collection.FindOneAndUpdate(ctx, bson.M{"_id": meetingId}, update, opts).Decode(&meetingSummary)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error fetching updated meeting")
		return
	}

	respondWithJSON(w, http.StatusOK, ConvertSummaryToFull(meetingSummary))
}

func AddParticipantsToMeeting(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Missing 'id' parameter")
		return
	}

	var request MeetingSummary
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if len(request.Participants) == 0 {
		respondWithError(w, http.StatusBadRequest, "No participants provided to add")
		return
	}

	for _, participantID := range request.Participants {
		_, err := uuid.Parse(participantID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid participant ID: "+participantID)
			return
		}
	}

	collection := database.Collection("meetings")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{
		"$addToSet": bson.M{
			"participants": bson.M{
				"$each": request.Participants,
			},
		},
	}

	var meetingSummary MeetingSummary

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err = collection.FindOneAndUpdate(ctx, bson.M{"_id": id}, update, opts).Decode(&meetingSummary)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error fetching updated meeting")
		return
	}

	respondWithJSON(w, http.StatusOK, ConvertSummaryToFull(meetingSummary))
}

func RemoveParticipantsFromMeeting(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Missing 'id' parameter")
		return
	}

	var request MeetingSummary
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if len(request.Participants) == 0 {
		respondWithError(w, http.StatusBadRequest, "No meetings provided to add")
		return
	}

	for _, participantID := range request.Participants {
		_, err := uuid.Parse(participantID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid meeting ID: "+participantID)
			return
		}
	}

	collection := database.Collection("meetings")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{
		"$pull": bson.M{
			"participants": bson.M{
				"$in": request.Participants,
			},
		},
	}

	var meetingSummary MeetingSummary

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err = collection.FindOneAndUpdate(ctx, bson.M{"_id": id}, update, opts).Decode(&meetingSummary)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error fetching updated meeting")
		return
	}

	respondWithJSON(w, http.StatusOK, ConvertSummaryToFull(meetingSummary))
}

func AddAttachmentsToMeeting(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Missing 'id' parameter")
		return
	}

	var request MeetingSummary
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if len(request.Attachments) == 0 {
		respondWithError(w, http.StatusBadRequest, "No attachments provided to add")
		return
	}

	for _, attachmentID := range request.Attachments {
		_, err := uuid.Parse(attachmentID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid attachment ID: "+attachmentID)
			return
		}
	}

	collection := database.Collection("meetings")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{
		"$addToSet": bson.M{
			"attachments": bson.M{
				"$each": request.Attachments,
			},
		},
	}

	var meetingSummary MeetingSummary

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err = collection.FindOneAndUpdate(ctx, bson.M{"_id": id}, update, opts).Decode(&meetingSummary)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error fetching updated meeting")
		return
	}

	respondWithJSON(w, http.StatusOK, ConvertSummaryToFull(meetingSummary))
}

func RemoveAttachmentsFromMeeting(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Missing 'id' parameter")
		return
	}

	var request MeetingSummary
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if len(request.Attachments) == 0 {
		respondWithError(w, http.StatusBadRequest, "No attachments provided to remove")
		return
	}

	for _, attachmentID := range request.Attachments {
		_, err := uuid.Parse(attachmentID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid attachment ID: "+attachmentID)
			return
		}
	}

	collection := database.Collection("meetings")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{
		"$pull": bson.M{
			"attachments": bson.M{
				"$in": request.Attachments,
			},
		},
	}

	var meetingSummary MeetingSummary

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err = collection.FindOneAndUpdate(ctx, bson.M{"_id": id}, update, opts).Decode(&meetingSummary)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error fetching updated meeting")
		return
	}

	respondWithJSON(w, http.StatusOK, ConvertSummaryToFull(meetingSummary))
}

func GetMeetingsByIds(w http.ResponseWriter, r *http.Request) {
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

	collection := database.Collection("meetings")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"_id": bson.M{
			"$in": ids,
		},
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error fetching meetings")
		return
	}
	defer cursor.Close(ctx)

	var meetingSummaries []MeetingSummary
	if err = cursor.All(ctx, &meetingSummaries); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding meetings")
		return
	}

	var meetings []Meeting
	for _, meetingSummary := range meetingSummaries {
		meetings = append(meetings, ConvertSummaryToFull(meetingSummary))
	}

	respondWithJSON(w, http.StatusOK, meetingSummaries)
}

func ConvertSummaryToFull(meetingSummary MeetingSummary) Meeting {
	calendars := GetCalendarsByIds(meetingSummary.Calendars)
	// attachments := GetAttachmentsByIds(meetingSummary.Attachments)
	// participants := getParticipantsByIds(meetingSummary.Participants)

	var meeting = Meeting{
		ID:        meetingSummary.ID,
		Title:     meetingSummary.Title,
		Details:   meetingSummary.Details,
		Calendars: calendars,
		// Attachments:  attachments,
		// Participants: participants,
	}

	return meeting
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
