import React, {useEffect, useState} from "react";
import axios from "axios";

const Meetings = () => {
    const [meetings, setMeetings] = useState([]);
    const [isEditing, setIsEditing] = useState(false);
    const [editIndex, setEditIndex] = useState(null);
    const [newParticipantUUID, setNewParticipantUUID] = useState('');
    const [newAttachmentUUID, setNewAttachmentUUID] = useState('');

    const [editingElement, setEditingElement] = useState({
        id: '',
        title: '',
        details: '',
        dateTime: '',
        location: '',
        participantIds: '',
    });

    const fetchMeetings = async () => {
        try {
            const response = await axios.get(`/api/meetings`);
            setMeetings(response.data);
        } catch (error) {
            console.error('Error fetching data: ', error);
        }
    };

    useEffect(() => {
        fetchMeetings();
    }, []);

    const handleCreate = async () => {
        await axios.post(`/api/meetings`, {
            ...editingElement,
            id: editingElement.id ? editingElement.id : undefined,
            participantIds: [editingElement.participantIds]
        });
        fetchMeetings();
        resetForm();
    }

    const handleEdit = (index, item) => {
        setIsEditing(true);
        setEditIndex(index);
        setEditingElement(item);
    };

    const handleDelete = async (index, id) => {
        await axios.delete(`/api/meetings/${id}`);
        fetchMeetings();
    };

    const handleInputChange = (e) => {
        setEditingElement({...editingElement, [e.target.name]: e.target.value});
    };

    const handleParticipantUUIDChange = (e) => {
        setNewParticipantUUID(e.target.value);
    };

    const handleAddParticipant = async () => {
        if (newParticipantUUID) {
            try {
                await axios.post(`/api/meetings/${editingElement.id}/participants`,
                    [newParticipantUUID]
                );
                fetchMeetings();
                setNewParticipantUUID('');
            } catch (error) {
                console.error('Error adding participant: ', error);
            }
        }
    };

    const handleAttachmentUUIDChange = (e) => {
        setNewAttachmentUUID(e.target.value);
    };

    const handleAddAttachment = async () => {
        if (newAttachmentUUID) {
            try {
                await axios.post(`/api/meetings/${editingElement.id}/attachments`,
                    [newAttachmentUUID]
                );
                fetchMeetings();
                setNewAttachmentUUID('');
            } catch (error) {
                console.error('Error adding participant: ', error);
            }
        }
    };

    const handleUpdate = async () => {
        const updatedItem = {
            id: editingElement.id,
            title: editingElement.title,
            details: editingElement.details,
            dateTime: editingElement.dateTime,
            location: editingElement.location,
        };
        const response = await axios.put(`/api/meetings/${updatedItem.id}`, updatedItem);
        setMeetings(meetings.map(meeting => meeting.id === response.data.id ? response.data : meeting));
        resetForm();
    };

    const resetForm = () => {
        setEditingElement({
            id: '',
            title: '',
            details: '',
            dateTime: '',
            location: '',
        });
        setIsEditing(false);
        setEditIndex(null);
    };

    return (
        <div>
            <div className="form">
                <input name="uuid" value={editingElement.id} onChange={handleInputChange}
                       placeholder="Meeting UUID"/>
                <input name="title" value={editingElement.title} onChange={handleInputChange}
                       placeholder="Title"/>
                <input name="dateTime" value={editingElement.dateTime} onChange={handleInputChange}
                       placeholder="Date and Time (YYYY-MM-DD HH:MM AM/PM)"/>
                <input name="location" value={editingElement.location} onChange={handleInputChange}
                       placeholder="Location"/>
                <textarea name="details" value={editingElement.details} onChange={handleInputChange}
                          placeholder="Details"/>
                <input name="participantIds" value={editingElement.participantIds} onChange={handleInputChange}
                       placeholder="Participant ID"/>
            </div>

            <button onClick={isEditing ? handleUpdate : handleCreate}>
                {isEditing ? 'Save Changes' : 'Create'}
            </button>

            <div className="participant-form">
                <h4>Add Participant by UUID</h4>
                <input name="participantUUID" value={newParticipantUUID} onChange={handleParticipantUUIDChange}
                       placeholder="Participant UUID"/>
                <button onClick={handleAddParticipant}>Add Participant</button>
            </div>

            <div className="attachment-form" style={{marginBottom: '50px'}}>
                <h4>Add Attachment by UUID</h4>
                <input name="attachmentUUID" value={newAttachmentUUID} onChange={handleAttachmentUUIDChange}
                       placeholder="Attachment UUID"/>
                <button onClick={handleAddAttachment}>Add Attachment</button>
            </div>

            {meetings.map((meeting, index) => (
                <div style={{marginBottom: '30px'}}>
                    <hr/>
                    <span>
                        <strong>UUID:</strong> {meeting.id} | <strong>Title:</strong> {meeting.title} | <strong>Date & Time:</strong> {meeting.dateTime} | <strong>Location:</strong> {meeting.location} | <strong>Details:</strong> {meeting.details}
                    </span>
                    <br></br>
                    {meeting.calendars.length > 0 &&
                        <span style={{display: 'block'}}>
                        Calendars:
                        </span>
                    }
                    {meeting.calendars.map((calendar, index) => (
                        <span style={{display: 'block', marginLeft: '20px'}}>
                            <strong>UUID:</strong> {calendar.id} |
                            <strong> Title:</strong> {calendar.title} |
                            <strong> Details:</strong> {calendar.details} |
                        </span>
                    ))}

                    {meeting.participants.length > 0 &&
                        <span style={{display: 'block'}}>
                            Participants:
                        </span>
                    }

                    {meeting.participants.map((participant, index) => (
                        <span style={{display: 'block', marginLeft: '20px'}}>
                            <strong>UUID:</strong> {participant.id} |
                            <strong> Name:</strong> {participant.name} |
                            <strong> Email:</strong> {participant.email} |
                        </span>
                    ))}
                    {meeting.attachments.length > 0 &&
                        <span style={{display: 'block'}}>
                            Attachments:
                        </span>
                    }

                    {meeting.attachments.map((attachment, index) => (
                        <span style={{display: 'block', marginLeft: '20px'}}>
                            <strong>UUID:</strong> {attachment.id} |
                            <strong> URL:</strong> {attachment.url} |
                        </span>
                    ))}
                    <button onClick={() => handleEdit(index, meeting)}>Edit</button>
                    <button onClick={() => handleDelete(index, meeting.id)}>Delete</button>
                    <hr/>
                </div>
            ))
            }
        </div>
    )
}

export default Meetings;