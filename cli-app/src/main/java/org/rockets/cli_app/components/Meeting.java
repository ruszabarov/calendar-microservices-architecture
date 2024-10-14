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
    ;

    private List<Calendar> calendars = new ArrayList<>();
    ;

    public Meeting(String uuid) {
        this.meetingId = uuid;
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

    // Adds ParticipantID
    public void addParticipant(Participant participant) {
        if (participant.getParticipantId() != null && !participants.contains(participant)) {
            participants.add(participant);
        }
    }

    // Removes ParticipantID
    public void removeParticipant(Participant participant) {
        participants.remove(participant);
    }

    // Adds AttachmentId
    public void addAttachment(Attachment attachment) {
        if (!attachments.contains(attachment)) {
            attachments.add(attachment);
        }
    }

    // Removes AttachmentId
    public void removeAttachment(Attachment attachment) {
        attachments.remove(attachment);
    }

    public List<Calendar> getCalendars() {
        return this.calendars;
    }

    public void setCalendars(List<Calendar> calendars) {
        this.calendars = calendars;
    }

    public void addCalendar(Calendar calendar) {
        if (!this.calendars.contains(calendar)) {
            this.calendars.add(calendar);
        }
    }

    public void removeCalendar(Calendar calendar) {
        if (calendar != null) {
            this.calendars.remove(calendar);
        }
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
                + getDateTime() + " | location: " + getLocation() + " | details: "
                + getDetails();
    }

    public String participantsToString() {
        StringBuilder result = new StringBuilder("Participants:\n");
        for (Participant p : getParticipants()) {
            result.append("\t").append(p.toString()).append("\n");
        }
        return result.toString();
    }

    public String attachmentsToString() {
        StringBuilder result = new StringBuilder("Attachments:\n");

        for (Attachment a : getAttachments()) {
            result.append("\t").append(a.toString()).append("\n");
        }

        return result.toString();
    }

    public String calendarsToString() {
        StringBuilder result = new StringBuilder("Calendars:\n");

        for (Calendar c : getCalendars()) {
            result.append("\t").append(c.toString()).append("\n");
        }

        return result.toString();
    }

}
