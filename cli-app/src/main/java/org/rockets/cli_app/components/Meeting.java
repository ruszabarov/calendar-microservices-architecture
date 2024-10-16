package org.rockets.cli_app.components;

import java.util.ArrayList;
import java.util.List;
import java.util.Objects;

public class Meeting {
    private final String meetingId;
    private String title;
    private String dateTime;
    private String location;
    private String details;
    private List<Participant> participants = new ArrayList<>();
    private List<Attachment> attachments = new ArrayList<>();
    private List<Calendar> calendars = new ArrayList<>();

    public Meeting(String meetingId) {
        this.meetingId = meetingId;
    }

    public Meeting(String meetingId, String title, String dateTime, String location, String details) {
        this(meetingId);
        this.title = title;
        this.dateTime = dateTime;
        this.location = location;
        this.details = details;
    }

    public String getMeetingId() {
        return meetingId;
    }

    public String getTitle() {
        return title;
    }

    public void setTitle(String title) {
        this.title = title;
    }

    public String getDateTime() {
        return dateTime;
    }

    public void setDateTime(String dateTime) {
        this.dateTime = dateTime;
    }

    public String getLocation() {
        return location;
    }

    public void setLocation(String location) {
        this.location = location;
    }

    public String getDetails() {
        return details;
    }

    public void setDetails(String details) {
        this.details = details;
    }

    public List<Participant> getParticipants() {
        return participants;
    }

    public void setParticipants(List<Participant> participants) {
        this.participants = participants;
    }

    public List<Attachment> getAttachments() {
        return attachments;
    }

    public void setAttachments(List<Attachment> attachments) {
        this.attachments = attachments;
    }

    public List<Calendar> getCalendars() {
        return calendars;
    }

    public void setCalendars(List<Calendar> calendars) {
        this.calendars = calendars;
    }

    public void addParticipant(Participant participant) {
        if (participant != null && participant.getParticipantId() != null && !participants.contains(participant)) {
            participants.add(participant);
        }
    }

    public void removeParticipant(Participant participant) {
        participants.remove(participant);
    }

    public void addAttachment(Attachment attachment) {
        if (attachment != null && attachment.getAttachmentId() != null && !attachments.contains(attachment)) {
            attachments.add(attachment);
        }
    }

    public void removeAttachment(Attachment attachment) {
        attachments.remove(attachment);
    }

    public void addCalendar(Calendar calendar) {
        if (calendar != null && calendar.getCalendarId() != null && !calendars.contains(calendar)) {
            calendars.add(calendar);
        }
    }

    public void removeCalendar(Calendar calendar) {
        calendars.remove(calendar);
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        Meeting meeting = (Meeting) o;
        return Objects.equals(meetingId, meeting.meetingId);
    }

    @Override
    public int hashCode() {
        return Objects.hash(meetingId);
    }

    @Override
    public String toString() {
        return "(Meeting) id: " + getMeetingId() + " | title: " + getTitle() + " | date: "
                + getDateTime() + " | location: " + getLocation() + " | details: " + getDetails();
    }
}
