import React, { useState, useEffect } from "react";
import axios from "axios";
import { useParams, useNavigate } from "react-router-dom";

// Separate EventDashboard Component
const EventDashboard = () => {
  const { eventId } = useParams();
  const navigate = useNavigate();
  const [eventDetails, setEventDetails] = useState(null);
  const [facultyList, setFacultyList] = useState([]);
  const [roles, setRoles] = useState([]);
  const [loading, setLoading] = useState(true);
  const [showAssignForm, setShowAssignForm] = useState(false);
  const [assignmentFormData, setAssignmentFormData] = useState({
    facultyId: "",
    roleId: "",
    points: 0
  });

  useEffect(() => {
    // Fetch event details
    axios.get(`https://09da6b2c-7088-466c-b94a-4662e3e1bd28.mock.pstmn.io/events/${eventId}`)
      .then(res => {
        setEventDetails(res.data);
      })
      .catch(err => console.error("Error fetching event details:", err));

    // Fetch faculty list
    axios.get(`https://09da6b2c-7088-466c-b94a-4662e3e1bd28.mock.pstmn.io/faculty`)
      .then(res => {
        setFacultyList(res.data);
      })
      .catch(err => console.error("Error fetching faculty list:", err));

    // Fetch roles
    axios.get(`https://09da6b2c-7088-466c-b94a-4662e3e1bd28.mock.pstmn.io/roles`)
      .then(res => {
        setRoles(res.data);
        setLoading(false);
      })
      .catch(err => {
        console.error("Error fetching roles:", err);
        setLoading(false);
      });
  }, [eventId]);

  const handleAssignmentInputChange = (e) => {
    const { name, value } = e.target;
    setAssignmentFormData(prev => ({
      ...prev,
      [name]: name === "points" ? parseInt(value) || 0 : value
    }));
  };

  const handleAssignRole = (e) => {
    e.preventDefault();
    
    // In a real application, you'd make a POST request to your backend
    axios.post(`https://09da6b2c-7088-466c-b94a-4662e3e1bd28.mock.pstmn.io/events/${eventId}/assignments`, {
      faculty_id: assignmentFormData.facultyId,
      role_id: assignmentFormData.roleId,
      points: assignmentFormData.points
    })
    .then(response => {
      console.log("Role assigned:", response.data);
      setShowAssignForm(false);
      // Refresh faculty list with updated assignments
      axios.get(`https://09da6b2c-7088-466c-b94a-4662e3e1bd28.mock.pstmn.io/faculty`)
        .then(res => {
          setFacultyList(res.data);
        })
        .catch(err => console.error("Error fetching faculty list:", err));
    })
    .catch(error => {
      console.error("Error assigning role:", error);
    });
  };

  if (loading) {
    return <div>Loading event details...</div>;
  }

  return (
    <div>
      <div className="mb-6 flex items-center">
        <button
          onClick={() => navigate("/admin")}
          className="mr-4 text-blue-600 hover:text-blue-800"
        >
          ← Back to Dashboard
        </button>
        <h2 className="text-xl font-semibold">
          Event Dashboard: {eventDetails?.event_name || eventId}
        </h2>
      </div>

      {eventDetails && (
        <div className="bg-white shadow rounded p-4 mb-6">
          <h3 className="text-lg font-medium mb-2">Event Details</h3>
          <p><strong>Name:</strong> {eventDetails.event_name}</p>
          <p><strong>Description:</strong> {eventDetails.event_description}</p>
          <p>
            <strong>Duration:</strong> {new Date(eventDetails.start_date).toLocaleDateString()} - {new Date(eventDetails.end_date).toLocaleDateString()}
          </p>
        </div>
      )}

      <div className="bg-white shadow rounded p-4 mb-6">
        <div className="flex justify-between items-center mb-4">
          <h3 className="text-lg font-medium">Faculty Assignments</h3>
          <button
            className="px-3 py-1 bg-blue-600 text-white rounded hover:bg-blue-700"
            onClick={() => setShowAssignForm(!showAssignForm)}
          >
            {showAssignForm ? "Cancel" : "Assign Role"}
          </button>
        </div>

        {showAssignForm && (
          <div className="border rounded p-4 mb-4 bg-gray-50">
            <h4 className="font-medium mb-3">Assign Role to Faculty</h4>
            <form onSubmit={handleAssignRole} className="space-y-4">
              <div>
                <label className="block mb-1">Faculty Member</label>
                <select
                  name="facultyId"
                  value={assignmentFormData.facultyId}
                  onChange={handleAssignmentInputChange}
                  className="w-full px-3 py-2 border rounded"
                  required
                >
                  <option value="">-- Select Faculty --</option>
                  {facultyList.map(faculty => (
                    <option key={faculty.id} value={faculty.id}>
                      {faculty.name}
                    </option>
                  ))}
                </select>
              </div>
              
              <div>
                <label className="block mb-1">Role</label>
                <select
                  name="roleId"
                  value={assignmentFormData.roleId}
                  onChange={handleAssignmentInputChange}
                  className="w-full px-3 py-2 border rounded"
                  required
                >
                  <option value="">-- Select Role --</option>
                  {roles.map(role => (
                    <option key={role.id} value={role.id}>
                      {role.name}
                    </option>
                  ))}
                </select>
              </div>
              
              <div>
                <label className="block mb-1">Points</label>
                <input
                  type="number"
                  name="points"
                  value={assignmentFormData.points}
                  onChange={handleAssignmentInputChange}
                  min="0"
                  className="w-full px-3 py-2 border rounded"
                  required
                />
              </div>
              
              <button
                type="submit"
                className="px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700"
              >
                Assign
              </button>
            </form>
          </div>
        )}

        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Faculty</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Role</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Points</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {facultyList.filter(faculty => faculty.assignments?.some(a => a.event_id === eventId)).map(faculty => (
                faculty.assignments
                  .filter(assignment => assignment.event_id === eventId)
                  .map((assignment, index) => (
                    <tr key={`${faculty.id}-${index}`}>
                      <td className="px-6 py-4 whitespace-nowrap">{faculty.name}</td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        {roles.find(r => r.id === assignment.role_id)?.name || assignment.role_id}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">{assignment.points}</td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <button 
                          className="text-red-600 hover:text-red-900 mr-2"
                          onClick={() => {
                            // Handle role removal
                            if (window.confirm("Are you sure you want to remove this assignment?")) {
                              console.log("Remove assignment:", assignment.id);
                              // API call would go here
                            }
                          }}
                        >
                          Remove
                        </button>
                        <button 
                          className="text-blue-600 hover:text-blue-900"
                          onClick={() => {
                            // Set up edit form
                            setAssignmentFormData({
                              facultyId: faculty.id,
                              roleId: assignment.role_id,
                              points: assignment.points
                            });
                            setShowAssignForm(true);
                          }}
                        >
                          Edit
                        </button>
                      </td>
                    </tr>
                  ))
              ))}
              {facultyList.filter(faculty => faculty.assignments?.some(a => a.event_id === eventId)).length === 0 && (
                <tr>
                  <td colSpan="4" className="px-6 py-4 text-center text-gray-500">
                    No faculty assignments for this event yet.
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
};

export default EventDashboard;