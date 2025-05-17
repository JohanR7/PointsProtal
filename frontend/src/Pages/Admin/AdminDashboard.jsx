import React, { useState, useEffect } from "react";
import axios from "axios";
import { Link, useNavigate } from "react-router-dom";

// Main Admin Dashboard Component (No longer needs Routes)
const AdminDashboard = () => {
  const [events, setEvents] = useState([]);
  const navigate = useNavigate();
  const [showCreateForm, setShowCreateForm] = useState(false);
  const [formData, setFormData] = useState({
    eventName: "",
    eventDescription: "",
    startDate: "",
    endDate: "",
  });

  useEffect(() => {
    fetchEvents();
  }, []);

  const fetchEvents = () => {
    axios
      .get("https://09da6b2c-7088-466c-b94a-4662e3e1bd28.mock.pstmn.io/events")
      .then((res) => setEvents(res.data))
      .catch((err) => console.error("Error fetching events:", err));
  };

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
  };

  const handleFormSubmit = (e) => {
    e.preventDefault();
    
    // In a real application, you'd make a POST request to your backend
    axios.post("https://09da6b2c-7088-466c-b94a-4662e3e1bd28.mock.pstmn.io/events", {
      event_name: formData.eventName,
      event_description: formData.eventDescription,
      start_date: formData.startDate,
      end_date: formData.endDate
    })
    .then(response => {
      console.log("Event created:", response.data);
      // Reset form
      setFormData({
        eventName: "",
        eventDescription: "",
        startDate: "",
        endDate: "",
      });
      // Hide form
      setShowCreateForm(false);
      // Refresh event list
      fetchEvents();
    })
    .catch(error => {
      console.error("Error creating event:", error);
    });
  };

  return (
    <div>
      <h2 className="text-xl font-semibold mb-4">Admin Dashboard</h2>
      
      <button
        className="mb-6 px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
        onClick={() => setShowCreateForm(!showCreateForm)}
      >
        {showCreateForm ? "Cancel" : "Create Event"}
      </button>

      {showCreateForm && (
        <div className="bg-white shadow rounded p-6 mb-6">
          <h3 className="text-xl font-semibold mb-4">Create New Event</h3>
          <form onSubmit={handleFormSubmit} className="space-y-4">
            <div>
              <label className="block mb-1 font-medium">Event Name</label>
              <input
                type="text"
                name="eventName"
                value={formData.eventName}
                onChange={handleInputChange}
                className="w-full px-3 py-2 border rounded"
                required
              />
            </div>
            <div>
              <label className="block mb-1 font-medium">Event Description</label>
              <textarea
                name="eventDescription"
                value={formData.eventDescription}
                onChange={handleInputChange}
                className="w-full px-3 py-2 border rounded"
                rows="3"
              />
            </div>
            <div>
              <label className="block mb-1 font-medium">Start Date</label>
              <input
                type="date"
                name="startDate"
                value={formData.startDate}
                onChange={handleInputChange}
                className="w-full px-3 py-2 border rounded"
                required
              />
            </div>
            <div>
              <label className="block mb-1 font-medium">End Date</label>
              <input
                type="date"
                name="endDate"
                value={formData.endDate}
                onChange={handleInputChange}
                className="w-full px-3 py-2 border rounded"
                required
              />
            </div>
            <div className="flex items-center gap-4">
              <button
                type="submit"
                className="px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700"
              >
                Create Event
              </button>
              <button
                type="button"
                className="px-4 py-2 bg-gray-400 text-white rounded hover:bg-gray-500"
                onClick={() => setShowCreateForm(false)}
              >
                Cancel
              </button>
            </div>
          </form>
        </div>
      )}

      <div>
        <h1 className="text-2xl font-bold mb-4">Current Events</h1>
        <ul className="space-y-2">
          {events.map((event) => (
            <li
              key={event.event_id}
              className="p-4 border rounded bg-white shadow hover:bg-blue-50 cursor-pointer transition"
              onClick={() => navigate(`/admin/event/${event.event_id}`)}
            >
              <h2 className="text-lg font-semibold">{event.event_name}</h2>
              <p className="text-sm text-gray-600">{event.event_description}</p>
              <div className="flex justify-between mt-2">
                <span className="text-xs text-gray-500">
                  {new Date(event.start_date).toLocaleDateString()} - {new Date(event.end_date).toLocaleDateString()}
                </span>
                <Link
                  to={`/admin/event/${event.event_id}`}
                  className="text-blue-600 hover:text-blue-800 text-sm"
                >
                  Manage Event →
                </Link>
              </div>
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
};

export default AdminDashboard;