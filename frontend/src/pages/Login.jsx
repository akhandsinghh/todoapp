import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

export default function Login() {
  const { login } = useAuth();
  const navigate = useNavigate();
  const [form, setForm] = useState({ email: '', password: '' });
  const [error, setError] = useState('');

  const submit = async (e) => {
    e.preventDefault();
    setError('');
    try {
      await login(form);
      navigate('/');
    } catch (err) {
      setError(err.response?.data?.message || err.response?.data?.error || 'Login failed');
    }
  };

  return (
    <main className="auth-screen">
      <form className="auth-panel" onSubmit={submit}>
        <h1>Todo App</h1>
        <p>Sign in to manage tasks, groups, and reminders.</p>
        {error && <div className="alert">{error}</div>}
        
        <label>
          Email
          <input
            type="email"
            value={form.email}
            onChange={(e) => setForm({ ...form, email: e.target.value })}
            required
          />
        </label>
        
        <label>
          Password
          <input
            type="password"
            value={form.password}
            onChange={(e) => setForm({ ...form, password: e.target.value })}
            required
          />
        </label>
        
        {/* HERE IS THE FORGOT PASSWORD LINK */}
        <div style={{ textAlign: 'right', marginTop: '-10px', marginBottom: '15px' }}>
          <Link to="/forgot-password" style={{ fontSize: '0.9rem', color: '#007bff', textDecoration: 'none' }}>
            Forgot Password?
          </Link>
        </div>

        <button type="submit">Sign in</button>
        
        <span style={{ marginTop: '15px', display: 'block', textAlign: 'center' }}>
          New here? <Link to="/register">Create account</Link>
        </span>
      </form>
    </main>
  );
}