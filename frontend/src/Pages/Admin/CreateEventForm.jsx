import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom'; 

export default function CreateEventForm() {
  const navigate = useNavigate(); 
  const [event, setEvent] = useState({
    event_name: '',
    start_date: '',
    start_time: '',
    end_date: '',
    end_time: '',
    event_description: '',
  });

  const [role, setRole] = useState('');
  const [points, setPoints] = useState('');
  const [headcount, setHeadcount] = useState('');
  const [roles, setRoles] = useState([]);
  const [createdEvent, setCreatedEvent] = useState(null);
  const [error, setError] = useState(null);

  const handleChange = (e) => {
    setEvent({ ...event, [e.target.name]: e.target.value });
  };

  const handleAddRole = () => {
    if (!role || !points || !headcount) return;
    setRoles([...roles, { role, points, headcount }]);
    setRole('');
    setPoints('');
    setHeadcount('');
  };

  const handleDeleteRole = (index) => {
    const updatedRoles = roles.filter((_, i) => i !== index);
    setRoles(updatedRoles);
  };

  const handleSubmitEvent = async (e) => {
    e.preventDefault();
    try {
      const eventRes = await fetch('http://localhost:8085/events', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(event),
      });

      const eventData = await eventRes.json();
      if (!eventRes.ok) {
        setError('Error creating event: ' + eventData.error);
        return;
      }

      setCreatedEvent(eventData);
      setError(null);
    } catch (err) {
      console.error('Network error:', err);
      setError('Network error');
    }
  };

  const handleSubmitRoles = async (e) => {
    e.preventDefault();
    if (!createdEvent || !createdEvent.event_id) {
      setError('Please create the event first');
      return;
    }

    try {
      const eventId = createdEvent.event_id;

      for (let r of roles) {
        const rolePayload = {
          role_name: r.role,
          role_point: parseInt(r.points, 10),
          role_headcount: parseInt(r.headcount, 10),
          event_id: eventId,
        };

        console.log('Sending role data:', rolePayload); 

        const roleRes = await fetch('http://localhost:8085/roles', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(rolePayload),
        });

        const roleData = await roleRes.json();
        if (!roleRes.ok) {
          setError('Error creating role: ' + roleData.error);
          return;
        }
      }

      setRoles([]);
      alert('All roles submitted!');
      navigate('/events'); // Fixed navigation
    } catch (err) {
      console.error('Network error:', err);
      setError('Network error');
    }
  };

  return (
    <div className="max-w-3xl mx-auto p-6 bg-white rounded-lg shadow-md mt-10">
      <h2 className="text-2xl font-bold mb-6 text-gray-800">Create Event</h2>

      <form onSubmit={handleSubmitEvent} className="space-y-4">
        <div>
          <label htmlFor="event_name" className="block text-sm font-medium text-gray-700 mb-1">
            Event Name
          </label>
          <input
            id="event_name"
            type="text"
            name="event_name"
            placeholder="Event Name"
            value={event.event_name}
            onChange={handleChange}
            required
            className="w-full p-3 border rounded-md focus:ring-2 focus:ring-blue-500"
          />
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label htmlFor="start_date" className="block text-sm font-medium text-gray-700 mb-1">
              Start Date
            </label>
            <input
              id="start_date"
              type="date"
              name="start_date"
              value={event.start_date}
              onChange={handleChange}
              required
              className="w-full p-3 border rounded-md focus:ring-2 focus:ring-blue-500"
            />
          </div>
          <div>
            <label htmlFor="start_time" className="block text-sm font-medium text-gray-700 mb-1">
              Start Time
            </label>
            <input
              id="start_time"
              type="time"
              name="start_time"
              value={event.start_time}
              onChange={handleChange}
              required
              className="w-full p-3 border rounded-md focus:ring-2 focus:ring-blue-500"
            />
          </div>
          <div>
            <label htmlFor="end_date" className="block text-sm font-medium text-gray-700 mb-1">
              End Date
            </label>
            <input
              id="end_date"
              type="date"
              name="end_date"
              value={event.end_date}
              onChange={handleChange}
              required
              className="w-full p-3 border rounded-md focus:ring-2 focus:ring-blue-500"
            />
          </div>
          <div>
            <label htmlFor="end_time" className="block text-sm font-medium text-gray-700 mb-1">
              End Time
            </label>
            <input
              id="end_time"
              type="time"
              name="end_time"
              value={event.end_time}
              onChange={handleChange}
              required
              className="w-full p-3 border rounded-md focus:ring-2 focus:ring-blue-500"
            />
          </div>
        </div>

        <div>
          <label htmlFor="event_description" className="block text-sm font-medium text-gray-700 mb-1">
            Event Description
          </label>
          <textarea
            id="event_description"
            name="event_description"
            placeholder="Event Description"
            value={event.event_description}
            onChange={handleChange}
            required
            rows={4}
            className="w-full p-3 border rounded-md focus:ring-2 focus:ring-blue-500"
          />
        </div>

        <button
          type="submit"
          className="bg-blue-600 hover:bg-blue-700 text-white py-2 px-6 rounded-md transition duration-150"
        >
          Create Event
        </button>
      </form>

      {createdEvent && (
        <div className="mt-8 p-5 border rounded-md bg-green-50">
          <h3 className="text-lg font-semibold text-green-800 mb-2">Event Created Successfully</h3>
          <div className="text-sm text-gray-700 space-y-1">
            <p><strong>Name:</strong> {createdEvent.event_name}</p>
            <p><strong>Start:</strong> {createdEvent.start_date} {createdEvent.start_time}</p>
            <p><strong>End:</strong> {createdEvent.end_date} {createdEvent.end_time}</p>
            <p><strong>Description:</strong> {createdEvent.event_description}</p>
            <p><strong>ID:</strong> {createdEvent.event_id}</p>
          </div>
        </div>
      )}

      <div className="border-t mt-10 pt-6">
        <h2 className="text-xl font-semibold text-gray-800 mb-4">Add Volunteer Roles</h2>
        <form className="space-y-4">
          <div className="flex flex-col sm:flex-row gap-3">
            <input
              className="flex-1 p-3 border rounded-md focus:ring-2 focus:ring-blue-500 placeholder-gray-500"
              type="text"
              placeholder="Role title"
              value={role}
              onChange={(e) => setRole(e.target.value)}
              required
            />
            <input
              className="w-24 p-3 border rounded-md focus:ring-2 focus:ring-blue-500 placeholder-gray-500"
              type="number"
              placeholder="Points"
              value={points}
              onChange={(e) => setPoints(e.target.value)}
              required
            />
            <input
              className="w-32 p-3 border rounded-md focus:ring-2 focus:ring-blue-500 placeholder-gray-500"
              type="number"
              placeholder="Headcount"
              value={headcount}
              onChange={(e) => setHeadcount(e.target.value)}
              required
            />
          </div>

          <div className="flex gap-4">
            <button
              type="button"
              className="bg-green-600 hover:bg-green-700 text-white px-6 py-2 rounded-md"
              onClick={handleAddRole}
            >
              Add
            </button>
            <button
              type="button" 
              onClick={handleSubmitRoles}
              className="bg-indigo-600 hover:bg-indigo-700 text-white px-6 py-2 rounded-md"
            >
              Submit Roles
            </button>
          </div>
        </form>

        {roles.length > 0 && (
          <ul className="space-y-3 mt-6">
            {roles.map((r, index) => (
              <li key={index} className="flex justify-between items-center p-3 border rounded-md bg-white shadow-sm">
                <span className="text-gray-800 font-medium">
                  {r.role} - {r.points} pts | {r.headcount} needed
                </span>
                <button
                  type="button"
                  className="bg-red-600 hover:bg-red-700 text-white px-3 py-1 rounded-md"
                  onClick={() => handleDeleteRole(index)}
                >
                  Remove
                </button>
              </li>
            ))}
          </ul>
        )}
      </div>

      {error && (
        <div className="mt-6 text-red-600 font-medium bg-red-50 p-3 rounded-md">
          <p>{error}</p>
        </div>
      )}
    </div>
  );
}