import React from 'react';
import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom';
import { AuthProvider } from './context/AuthContext';
import PrivateRoute from './routes/PrivateRoute';
import Login from './pages/Login';
import Register from './pages/Register';
import Dashboard from './pages/Dashboard';
import ForgotPassword from './pages/ForgotPassword'; 
import ChangePassword from './components/ChangePassword'; // <-- NEW IMPORT
import './styles/global.css';
import './styles/dashboard.css';

export default function App() {
  return (
    <AuthProvider>
      <BrowserRouter>
        <Routes>
          {/* PUBLIC ROUTES */}
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />
          <Route path="/forgot-password" element={<ForgotPassword />} />
          
          {/* PROTECTED ROUTES (Requires valid JWT token) */}
          <Route path="/" element={<PrivateRoute><Dashboard /></PrivateRoute>} />
          <Route 
            path="/change-password" 
            element={
              <PrivateRoute>
                <div className="auth-container" style={{ marginTop: '50px' }}>
                  {/* Wrapping it in a container so it looks nice as a standalone page */}
                  <ChangePassword />
                </div>
              </PrivateRoute>
            } 
          />
          
          {/* FALLBACK */}
          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </BrowserRouter>
    </AuthProvider>
  );
}