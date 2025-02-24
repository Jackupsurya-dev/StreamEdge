import React, { useEffect, useState } from "react";
import "bootstrap/dist/css/bootstrap.min.css"; // Import Bootstrap

const UserTable = () => {
  const [users, setUsers] = useState([]);

  useEffect(() => {
    const eventSource = new EventSource("http://localhost:8081/users/stream"); // Update with your SSE endpoint

    eventSource.onmessage = (event) => {
      const newUser = JSON.parse(event.data);
      setUsers((prevUsers) => [...prevUsers, newUser]);
    };

    return () => {
      eventSource.close();
    };
  }, []);

  return (
    <div className="container-fluid p-0">
      <h2 className="text-center py-3 bg-light border-bottom shadow-sm">
        📌 Real-time User Data
      </h2>
      <div className="table-responsive">
        <table className="table table-bordered table-striped table-hover text-center">
          <thead className="table-dark">
            <tr>
              <th>ID</th>
              <th>First Name</th>
              <th>Last Name</th>
              <th>Email</th>
              <th>Created At</th>
              <th>Deleted At</th>
              <th>Merged At</th>
              <th>Parent User ID</th>
            </tr>
          </thead>
          <tbody>
            {users.length > 0 ? (
              users.map((user) => (
                <tr key={user.id}>
                  <td>{user.id}</td>
                  <td>{user.first_name}</td>
                  <td>{user.last_name}</td>
                  <td>{user.email_address}</td>
                  <td>{new Date(Number(user.created_at)).toLocaleString()}</td>
                  <td>{user.deleted_at === "-1" ? "N/A" : user.deleted_at}</td>
                  <td>{user.merged_at === "-1" ? "N/A" : user.merged_at}</td>
                  <td>{user.parent_user_id ? user.parent_user_id : "N/A"}</td>
                </tr>
              ))
            ) : (
              <tr>
                <td colSpan="8" className="text-center">
                  No data available
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
};

export default UserTable;
