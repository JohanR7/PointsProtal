import React from "react";
import { useState } from "react";
import EventList from "./EventList";
import Sidebar from "../Layout/Sidebar";
const AdminDashboard = () => {
  const [selectedEventId, setSelectedEventId] = useState(null);

  const handleEventClick = (eventId) => {
    setSelectedEventId(eventId);
    console.log("Clicked event ID:", eventId);
  };

  return (
    <div className="flex">
      <main className="flex-1 p-6 bg-gray-100 min-h-screen">
        <h2 className="text-xl font-semibold mb-4">Welcome to the Dashboard</h2>

        <EventList onEventClick={handleEventClick} />

        {selectedEventId && (
          <div className="mt-6 p-4 border rounded bg-white shadow">
            <h3 className="text-lg font-semibold">Selected Event ID:</h3>
            <p>{selectedEventId}</p>
          </div>
        )}
      </main>
    </div>
  );
};

export default AdminDashboard;