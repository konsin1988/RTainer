import { Navigate } from "react-router-dom";
import { useAuth } from "react-oidc-context";

export function ProtectedRoute({
  children,
}: {
  children: React.ReactNode;
}) {

  const auth = useAuth();

  if (auth.isLoading) {
    return <div>Loading...</div>;
  }

  if (!auth.isAuthenticated) {
    return <Navigate to="/login" />;
  }

  return children;
}
