// Sidebar.jsx
import React from "react";
import { Home, History, Users } from "lucide-react";
import { useNavigate } from "react-router-dom";

const Sidebar = () => {
  const navigate = useNavigate();

  const menuItems = [
    { icon: <Home size={20} />, label: "Home", path: "/admin" },
    { icon: <History size={20} />, label: "Past Events", path: "/admin/event-history" },
    { icon: <Users size={20} />, label: "Users", path: "/admin/add-teacher" },
  ];

  return (
    <div className="w-64 h-screen bg-gray-900 text-white flex flex-col p-4">
      <h1 className="text-2xl font-bold mb-6">Dashboard</h1>
      <nav className="space-y-2">
        {menuItems.map((item, index) => (
          <button
            key={index}
            onClick={() => navigate(item.path)}
            className="flex items-center space-x-3 p-2 rounded-xl hover:bg-gray-700 w-full text-left"
          >
            {item.icon}
            <span>{item.label}</span>
          </button>
        ))}
      </nav>
    </div>
  );
};

export default Sidebar;
