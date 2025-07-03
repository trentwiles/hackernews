import { Navigate } from 'react-router-dom';
import Cookies from "js-cookie";

const ProtectedRoute = ({ children }) => {
  const hasAuth = Cookies.get("token") == undefined
  
  if (!hasAuth) {
    return <Navigate to="/login" replace />;
  }
  
  return children;
};

export default ProtectedRoute;