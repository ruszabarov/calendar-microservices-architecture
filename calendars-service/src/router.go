package main

import (
	"github.com/gorilla/mux"
)

func initializeRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/calendars", GetCalendarsByIds).Methods("GET").Queries("ids", "{ids}")
	router.HandleFunc("/calendars", GetCalendars).Methods("GET")
	router.HandleFunc("/calendars", CreateCalendar).Methods("POST")
	router.HandleFunc("/calendars/{id}", UpdateCalendarById).Methods("PUT")
	router.HandleFunc("/calendars/{id}", DeleteCalendarById).Methods("DELETE")
	router.HandleFunc("/calendars/{id}/meetings", AddMeetingsToCalendar).Methods("POST")

	return router
}
