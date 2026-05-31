import { Navigate, Outlet } from "react-router-dom";
import { useAuth } from "react-oidc-context";

export default function ProtectedRoute(){
  const auth = useAuth();

  if (auth.isLoading) {
    return <div>Loading...</div>;
  }

  if (!auth.isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  return <Outlet />;
}

