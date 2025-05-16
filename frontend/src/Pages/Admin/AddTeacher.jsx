import React, { useState } from 'react';

function AddTeacher() {
  const [teacher, setTeacher] = useState({
    teacher_name: '',
    email: '',
    department: '',
    position: '',
    profile_photo: ''
  });

  const [responseMessage, setResponseMessage] = useState('');
  const [errorMessage, setErrorMessage] = useState('');

  const handleChange = (e) => {
    setTeacher({ ...teacher, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    const res = await fetch('http://localhost:8085/teachers', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(teacher)
    });

    if (res.ok) {
      setResponseMessage('Teacher created successfully!');
      setTeacher({
        teacher_name: '',
        email: '',
        department: '',
        position: '',
        profile_photo: ''
      });
    } else {
      const errorData = await res.json();
      setResponseMessage(`Error: ${errorData.error}`);
    }
};


  return (
    <div className="p-4 max-w-md mx-auto">
      <h2 className="text-2xl font-bold mb-4">Create Teacher</h2>

      {/* Show success message */}
      {responseMessage && <div className="mb-4 text-green-600">{responseMessage}</div>}

      {/* Show error message */}
      {errorMessage && <div className="mb-4 text-red-600">{errorMessage}</div>}

      <form onSubmit={handleSubmit} className="space-y-4">
        <input
          type="text"
          name="teacher_name"
          placeholder="Name"
          value={teacher.name}
          onChange={handleChange}
          className="w-full border p-2 rounded"
          required
        />
        <input
          type="email"
          name="email"
          placeholder="Email"
          value={teacher.email}
          onChange={handleChange}
          className="w-full border p-2 rounded"
          required
        />
        <input
          type="text"
          name="department"
          placeholder="Department"
          value={teacher.department}
          onChange={handleChange}
          className="w-full border p-2 rounded"
        />
        <button type="submit" className="bg-blue-500 text-white px-4 py-2 rounded">
          Create
        </button>
      </form>
    </div>
  );
}

export default AddTeacher;
