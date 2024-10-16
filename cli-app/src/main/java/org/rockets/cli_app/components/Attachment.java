package org.rockets.cli_app.components;

import java.util.Objects;

public class Attachment {
    private final String attachmentId;
    private String attachmentUrl;

    public Attachment(String attachmentId) {
        this.attachmentId = attachmentId;
    }

    public Attachment(String attachmentId, String attachmentUrl) {
        this(attachmentId);
        this.attachmentUrl = attachmentUrl;
    }

    public String getAttachmentId() {
        return attachmentId;
    }

    public String getAttachmentUrl() {
        return attachmentUrl;
    }

    public void setAttachmentUrl(String attachmentUrl) {
        this.attachmentUrl = attachmentUrl;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        Attachment attachment = (Attachment) o;
        return Objects.equals(attachmentId, attachment.attachmentId);
    }

    @Override
    public int hashCode() {
        return Objects.hash(attachmentId);
    }

    @Override
    public String toString() {
        return "(Attachment) id: " + getAttachmentId() + " | url: " + getAttachmentUrl();
    }
}
