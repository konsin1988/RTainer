import {
  BrowserRouter,
  Routes,
  Route,
  Navigate,
} from "react-router-dom";

import LoginPage from "../pages/LoginPage";
import Dashboard from "../pages/Dashboard";
import AuthCallback from "../pages/AuthCallback";

import { ProtectedRoute } from "../auth/ProtectedRoute";

export default function AppRouter(){
  return (
  <BrowserRouter>
    <Routes>
      <Route 
        path="*" 
        element={<Navigate to="/dashboard" replace />} 
      />
    
      <Route
        path="/login"
        element={<LoginPage />}
      />
      
      <Route
        path="/auth/callback"
        element={<AuthCallback/>}
      />

      <Route
        path="/dashboard"
        element={
          <ProtectedRoute>
            <Dashboard />
          </ProtectedRoute>
        }
      />
  
    </Routes>
  </BrowserRouter>
  );
}
