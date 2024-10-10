import React, {useEffect, useState} from "react";
import axios from "axios";
import {API_URL} from "./Home";

const Attachments = () => {
    const [attachments, setAttachments] = useState([]);
    const [isEditing, setIsEditing] = useState(false);
    const [editIndex, setEditIndex] = useState(null);

    const [editingElement, setEditingElement] = useState({
        id: '',
        url: '',
    });

    const resetForm = () => {
        setEditingElement({
            id: '',
            url: '',
        });
        setIsEditing(false);
        setEditIndex(null);
    };


    const fetchAttachments = async () => {
        try {
            // Replace 'your-api-endpoint' with the actual endpoint you want to call
            const response = await axios.get(`/api/attachments`);
            setAttachments(response.data);
        } catch (error) {
            console.error('Error fetching data: ', error);
        }
    };

    useEffect(() => {
        fetchAttachments();
    }, []);

    const handleCreate = async () => {
        await axios.post(`/api/attachments`, editingElement);
        fetchAttachments();
        resetForm();
    }

    const handleEdit = (index, item) => {
        setIsEditing(true);
        setEditIndex(index);
        setEditingElement(item);
    };

    const handleDelete = async (index, id) => {
        await axios.delete(`/api/attachments/${id}`);
        fetchAttachments();
    };


    const handleUpdate = async () => {
        const updatedItem = {...editingElement};
        const response = await axios.put(`/api/attachments/${updatedItem.id}`, updatedItem);
        setAttachments(attachments.map(attachment => attachment.id === response.data.id ? response.data : attachment));
        resetForm();
    };

    const handleInputChange = (e) => {
        setEditingElement({...editingElement, [e.target.name]: e.target.value});
    };

    return (
        <div>
            <div className="form">
                <input name="id" value={editingElement.id} onChange={handleInputChange}
                       placeholder="Attachment UUID"/>
                <input name="url" value={editingElement.url} onChange={handleInputChange}
                       placeholder="Attachment URL"/>
            </div>

            <button onClick={isEditing ? handleUpdate : handleCreate}>
                {isEditing ? 'Save Changes' : 'Create'}
            </button>

            {attachments.map((attachment, index) => (
                <div>
                    <hr/>
                    <span>
                        <strong>UUID:</strong> {attachment.id} | <strong>URL:</strong> {attachment.url}
                    </span>
                    <button onClick={() => handleEdit(index, attachment)}>Edit</button>
                    <button onClick={() => handleDelete(index, attachment.id)}>Delete</button>
                    <hr/>
                </div>
            ))}


        </div>
    )
}

export default Attachments;