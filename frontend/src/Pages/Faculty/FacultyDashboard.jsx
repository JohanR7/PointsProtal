import React, { useEffect, useState } from 'react';
import axios from 'axios';

 const FacultyDashboard = () => {
  const [currentEvents, setCurrentEvents] = useState([]);
  const [pastEvents, setPastEvents] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // Replace with your actual backend URLs
    const fetchCurrent = axios.get('https://example.com/api/events/current');
    const fetchPast = axios.get('https://example.com/api/events/history');

    Promise.all([fetchCurrent, fetchPast])
      .then(([currentRes, pastRes]) => {
        setCurrentEvents(currentRes.data);
        setPastEvents(pastRes.data);
        setLoading(false);
      })
      .catch((err) => {
        console.error('Error fetching events:', err);
        setLoading(false);
      });
  }, []);

  if (loading) return <p>Loading events...</p>;

  return (
    <div className="p-6 space-y-10">
      <div>
        <h2 className="text-2xl font-bold mb-4">Current Events</h2>
        {currentEvents.length === 0 ? (
          <p>No current events.</p>
        ) : (
          <ul className="space-y-4">
            {currentEvents.map((event) => (
              <li
                key={event.id}
                className="p-4 border rounded bg-white shadow-sm hover:bg-gray-50"
              >
                <h3 className="text-lg font-semibold">{event.name}</h3>
                <p className="text-sm text-gray-600">{event.description}</p>
                <p className="text-xs text-gray-500">Date: {event.date}</p>
              </li>
            ))}
          </ul>
        )}
      </div>

      <div>
        <h2 className="text-2xl font-bold mb-4">Event History</h2>
        {pastEvents.length === 0 ? (
          <p>No past events found.</p>
        ) : (
          <ul className="space-y-4">
            {pastEvents.map((event) => (
              <li
                key={event.id}
                className="p-4 border rounded bg-white shadow-sm hover:bg-gray-50"
              >
                <h3 className="text-lg font-semibold">{event.name}</h3>
                <p className="text-sm text-gray-600">{event.description}</p>
                <p className="text-xs text-gray-500">Date: {event.date}</p>
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  );
};
export default FacultyDashboard;
