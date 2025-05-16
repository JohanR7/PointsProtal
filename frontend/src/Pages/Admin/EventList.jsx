import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import axios from "axios";

export default function EventList({ onEventClick }) {
  const [events, setEvents] = useState([]);

  useEffect(() => {
    axios.get("https://09da6b2c-7088-466c-b94a-4662e3e1bd28.mock.pstmn.io/events")
      .then((res) => setEvents(res.data))
      .catch((err) => console.error(err));
  }, []);

  return (
    <div>
      <h1 className="text-2xl font-bold mb-4">Current Events</h1>
      <ul className="space-y-2">
        {events.map((event) => (
          <li
            key={event.ID}
            className="p-4 border rounded hover:bg-gray-100 cursor-pointer"
            onClick={() => onEventClick(event.event_id)} 
          >
            <h2 className="text-lg font-semibold">{event.event_name}</h2>
            <p className="text-sm">{event.event_description}</p>
          </li>
        ))}
      </ul>
    </div>
  );
}

