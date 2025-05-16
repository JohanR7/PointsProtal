import React from 'react';
import './App.css';

import { BrowserRouter as Router, Routes, Route , Navigate } from 'react-router-dom';

import CreateEventForm from './Pages/Admin/CreateEventForm';
import AddTeacher from './Pages/Admin/AddTeacher';
import EventList from './Pages/Admin/EventList';
import Signup from './Pages/Auth/Signup';
import Login from './Pages/Auth/Login';
import AdminDashboard from './Pages/Admin/AdminDashboard';
import DashboardLayout from './Pages/Layout/DashboardLayout';
import EventHistory from './Pages/Admin/EventHistory';
import FacultyDashboard  from './Pages/Faculty/FacultyDashboard';

import PrivateRoute from './Pages/Auth/PrivateRoute';
import PublicRoute from './Pages/Auth/PublicRoute';

import { AuthProvider } from './Pages/Auth/AuthContext'; 

function App() {
  return (
<AuthProvider>
  <Router>
    <Routes>
      {/* Public routes */}
      <Route path="/" element={<Signup />} />
      <Route path="/signup" element={<Signup />} />
      <Route path="/login" element={<Login />} />

      {/* Admin routes - only accessible if role is 'admin' */}
      <Route
        path="/admin"
        element={
          <PrivateRoute allowedRoles={['admin']}>
            <DashboardLayout />
          </PrivateRoute>
        }
      >
        <Route index element={<AdminDashboard />} />
        <Route path="create-event" element={<CreateEventForm />} />
        <Route path="add-teacher" element={<AddTeacher />} />
        <Route path="event-history" element={<EventHistory />} />
        <Route path="events" element={<EventList />} />
      </Route>

      {/* Faculty route - only accessible if role is 'faculty' */}
      <Route
        path="/faculty"
        element={
          <PrivateRoute allowedRoles={['faculty']}>
            <FacultyDashboard />
          </PrivateRoute>
        }
      />

      {/* Unauthorized page route (you should create this page) */}
      <Route path="/unauthorized" element={<div>Unauthorized Access</div>} />

      {/* Catch-all route */}
      <Route path="*" element={<Navigate to="/" />} />
    </Routes>
  </Router>
</AuthProvider>

  );
}

export default App;
