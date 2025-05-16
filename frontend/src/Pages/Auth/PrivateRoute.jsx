import { Navigate } from "react-router-dom";
import { useAuth } from "./AuthContext";

const PrivateRoute = ({ children, allowedRoles }) => {
  const { auth } = useAuth();

  if (!auth) {
    // Not logged in, redirect to login
    return <Navigate to="/login" replace />;
  }

  if (allowedRoles && !allowedRoles.includes(auth.user.role)) {
    // Logged in but role not allowed, redirect to unauthorized page
    return <Navigate to="/unauthorized" replace />;
  }

  // Logged in and role allowed
  return children;
};

export default PrivateRoute;
