import React, {useEffect, useState} from "react";
import axios from "axios";

const Calendars = () => {
    // useStates
    const [calendars, setCalendars] = useState([]);
    const [isEditing, setIsEditing] = useState(false);
    const [editIndex, setEditIndex] = useState(null);
    const [newMeetingUUID, setNewMeetingUUID] = useState('');


    const [editingElement, setEditingElement] = useState({
        id: '',
        title: '',
        details: '',
        meetingId: ''
    });

    // Fetch calendar data
    const fetchCalendars = async () => {
        try {
            const response = await axios.get(`/api/calendars`);
            setCalendars(response.data);
        } catch (error) {
            console.error("Error fetching calendars data: ", error);
        }
    };

    // Use effect fetches calendar data from the server
    useEffect(() => {
        fetchCalendars();
    }, []);

    // Validations
    const validateInputs = () => {
        // Title is not longer than 2000 chars
        if (editingElement.title.length > 2000) {
            alert("Title should not exceed 2000 characters.");
            return false;
        }
        // Details is not longer than 10000 chars
        if (editingElement.details.length > 10000) {
            alert("Details should not exceed 10000 characters.");
            return false;
        }
        return true;
    };

    // Handle creation or update of a calendar
    const handleCreateOrUpdate = async () => {
        if (!validateInputs()) return;

        if (isEditing) {
            await axios.put(`/api/calendars/${editingElement.id}`, editingElement);
        } else {
            await axios.post(`/api/calendars`, {
                ...editingElement,
                id: editingElement.id ? editingElement : undefined,
                meetingIds: [editingElement.meetingId]
            });
        }
        fetchCalendars();
        resetForm();
    };

    // Edit mode
    const handleEdit = (index, item) => {
        setIsEditing(true);
        setEditIndex(index);
        setEditingElement(item);
    };

    // Handles delete
    const handleDelete = async (id) => {
        await axios.delete(`/api/calendars/${id}`);
        fetchCalendars();
    };

    // Handles the change
    const handleInputChange = (e) => {
        setEditingElement({...editingElement, [e.target.name]: e.target.value});
    };

    // Resets form
    const resetForm = () => {
        setEditingElement({
            id: '',
            title: '',
            details: '',
            meetingId: ''
        });
        setIsEditing(false);
        setEditIndex(null);
    };

    const handleMeetingUUIDChange = (e) => {
        setNewMeetingUUID(e.target.value);
    };

    const handleAddMeeting = async () => {
        if (newMeetingUUID) {
            try {
                const response = await axios.post(`/api/meetings/${editingElement.id}/participants`,
                    [newMeetingUUID]
                );
                setCalendars(calendars.map(calendar => calendar.id === response.data.id ? response.data : calendar));
                setNewMeetingUUID('');
            } catch (error) {
                console.error('Error adding participant: ', error);
            }
        }
    };

    return (
        <div>
            <div className="form">
                <input name="uuid" value={editingElement.id} onChange={handleInputChange} placeholder="Calendar UUID"/>
                <input name="title" value={editingElement.title} onChange={handleInputChange}
                       placeholder="Title (max 2000 characters)"/>
                <textarea name="details" value={editingElement.details} onChange={handleInputChange}
                          placeholder="Details (max 10000 characters)"/>
                <input name="meetingId" value={editingElement.meetingId} onChange={handleInputChange}
                       placeholder="Meeting ID"/>
            </div>

            <button onClick={handleCreateOrUpdate}>
                {isEditing ? 'Save Changes' : 'Create'}
            </button>

            <div className="meeting-form" style={{marginBottom: '50px'}}>
                <h4>Add Meeting by UUID</h4>
                <input name="meetingUUID" value={newMeetingUUID} onChange={handleMeetingUUIDChange}
                       placeholder="Participant UUID"/>
                <button onClick={handleAddMeeting}>Add Meeting</button>
            </div>

            {calendars.map((calendar, index) => (
                <div key={index} style={{marginBottom: '30px'}}>
                    <hr/>
                    <span>
                        <strong>UUID:</strong> {calendar.id} | 
                        <strong> Title:</strong> {calendar.title} | 
                        <strong> Details:</strong> {calendar.details} |
                    </span>

                    {calendar.meetings && calendar.meetings.length > 0 && (
                        <span style={{display: 'block', marginTop: '10px'}}>
                            <strong>Meetings:</strong>
                        </span>
                    )}
                    {calendar.meetings && calendar.meetings.map((meeting, index) => (
                        <span key={index} style={{display: 'block', marginLeft: '20px'}}>
                            <strong>UUID:</strong> {meeting.id} | 
                            <strong> Title:</strong> {meeting.title} | 
                            <strong> Date & Time:</strong> {meeting.dateTime} | 
                            <strong> Location:</strong> {meeting.location} | 
                            <strong> Details:</strong> {meeting.details}
                        </span>
                    ))}

                    <button onClick={() => handleEdit(index, calendar)}>Edit</button>
                    <button onClick={() => handleDelete(calendar.id)}>Delete</button>
                    <hr/>
                </div>
            ))}
        </div>
    );
};

export default Calendars;
