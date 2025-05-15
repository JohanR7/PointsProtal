import { Navigate } from "react-router-dom";
import { useAuth } from "./AuthContext";

const PublicRoute = ({ children }) => {
  const { auth } = useAuth();
  return !auth ? children : <Navigate to="/event" replace />;
};

export default PublicRoute;
