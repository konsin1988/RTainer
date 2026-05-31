import {
  BrowserRouter,
  Routes,
  Route,
  Navigate,
} from "react-router-dom";

import LoginPage from "../pages/LoginPage";
import Dashboard from "../pages/Dashboard";
import AuthCallback from "../pages/AuthCallback";
import ProtectedRoute from "../auth/ProtectedRoute";
import Layout from "../layouts/Layout.tsx";

export default function AppRouter(){
  return (
  <BrowserRouter>
    <Routes>
      <Route
        path="/auth/callback"
        element={<AuthCallback/>}
      />
      <Route element={<Layout/>}>
        <Route
          path="/login"
          element={<LoginPage />}
        />

        <Route
          path="/user"
          element={<ProtectedRoute/>}
        >
          <Route 
            path="dashboard"
            element={<Dashboard/>}
          />
        </Route>
      </Route>
      <Route 
        path="*" 
        element={<Navigate to="/user/dashboard" replace />} 
      />
    </Routes>
  </BrowserRouter>
  );
}
