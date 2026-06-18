import React from 'react';
import { useAuth } from '../context/AuthContext';

export default function Navbar() {
  const { user, logout } = useAuth();

  return (
    <header className="navbar">
      <div>
        <strong>Todo App</strong>
        <span>{user?.name}</span>
      </div>
      <button className="secondary" onClick={logout}>
        Logout
      </button>
    </header>
  );
}