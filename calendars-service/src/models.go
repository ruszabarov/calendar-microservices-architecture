package main

import (
	"net/url"
	"time"
)

type Calendar struct {
	ID       string    `json:"id" bson:"_id,omitempty"`
	Title    string    `json:"title" bson:"title"`
	Details  string    `json:"details" bson:"details"`
	Meetings []Meeting `json:"meetings" bson:"meetings"`
}

type CalendarSummary struct {
	ID       string   `json:"id" bson:"_id,omitempty"`
	Title    string   `json:"title" bson:"title"`
	Details  string   `json:"details" bson:"details"`
	Meetings []string `json:"meetings" bson:"meetings"`
}

type Meeting struct {
	ID           string        `json:"id" bson:"_id,omitempty"`
	Title        string        `json:"title" bson:"title"`
	Details      string        `json:"details" bson:"details"`
	DateTime     time.Time     `json:"datetime" bson:"datetime"`
	Location     string        `json:"location" bson:"location"`
	Participants []Participant `json:"participants" bson:"participants"`
	Attachments  []Attachment  `json:"attachments" bson:"attachments"`
}

type Participant struct {
	ID    string `json:"id" bson:"_id,omitempty"`
	Name  string `json:"name" bson:"name"`
	Email string `json:"email" bson:"email"`
}

type Attachment struct {
	ID  string   `json:"id" bson:"_id,omitempty"`
	URL *url.URL `json:"url" bson:"url"`
}
