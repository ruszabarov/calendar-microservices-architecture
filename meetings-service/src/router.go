package main

import (
	"github.com/gorilla/mux"
)

func initializeRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/meetings", GetMeetingsByIds).Methods("GET").Queries("ids", "{ids}")
	router.HandleFunc("/meetings", GetMeetings).Methods("GET")
	router.HandleFunc("/meetings", CreateMeeting).Methods("POST")
	router.HandleFunc("/meetings/{id}", UpdateMeetingById).Methods("PUT")
	router.HandleFunc("/meetings/{id}", DeleteMeetingById).Methods("DELETE")
	router.HandleFunc("/meetings/{id}/addParticipants", AddParticipantsToMeeting).Methods("POST")
	router.HandleFunc("/meetings/{id}/removeParticipants", RemoveParticipantsFromMeeting).Methods("POST")
	router.HandleFunc("/meetings/{id}/addAttachments", AddAttachmentsToMeeting).Methods("POST")
	router.HandleFunc("/meetings/{id}/removeAttachments", RemoveAttachmentsFromMeeting).Methods("POST")

	return router
}
