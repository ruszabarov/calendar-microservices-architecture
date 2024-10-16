package org.rockets.cli_app.components;

import java.util.Objects;

public class Participant {
    private final String participantId;
    private String name;
    private String email;

    public Participant(String participantId) {
        this.participantId = participantId;
    }

    public Participant(String participantId, String name, String email) {
        this(participantId);
        this.name = name;
        this.email = email;
    }

    public String getParticipantId() {
        return participantId;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public String getEmail() {
        return email;
    }

    public void setEmail(String email) {
        this.email = email;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        Participant participant = (Participant) o;
        return Objects.equals(participantId, participant.participantId);
    }

    @Override
    public int hashCode() {
        return Objects.hash(participantId);
    }

    @Override
    public String toString() {
        return "(Participant) id: " + getParticipantId() + " | name: " + getName() + " | email: " + getEmail();
    }
}
