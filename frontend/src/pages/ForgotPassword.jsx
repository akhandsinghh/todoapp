import React, { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { forgotPassword } from '../api/authApi';

export default function ForgotPassword() {
  const navigate = useNavigate();
  const [formData, setFormData] = useState({
    email: '',
    new_password: '',
    confirm_password: '',
  });
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setSuccess('');

    if (formData.new_password !== formData.confirm_password) {
      return setError('Passwords do not match');
    }

    setLoading(true);
    try {
      await forgotPassword(formData);
      setSuccess('Password reset successfully! You can now log in.');
      setTimeout(() => navigate('/login'), 3000); // Auto-redirect after 3 seconds
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to reset password');
    } finally {
      setLoading(false);
    }
  };

  return (
    <main className="auth-screen">
      <div className="auth-panel">
        <h2>Forgot Password</h2>
        {error && <div className="alert" style={{ color: 'red' }}>{error}</div>}
        {success && <div className="alert" style={{ color: 'green' }}>{success}</div>}
        
        <form onSubmit={handleSubmit} style={{ display: 'flex', flexDirection: 'column', gap: '15px' }}>
          <label>
            Email Address
            <input
              type="email"
              value={formData.email}
              onChange={(e) => setFormData({ ...formData, email: e.target.value })}
              required
            />
          </label>
          <label>
            New Password
            <input
              type="password"
              value={formData.new_password}
              onChange={(e) => setFormData({ ...formData, new_password: e.target.value })}
              required
              minLength="6"
            />
          </label>
          <label>
            Confirm New Password
            <input
              type="password"
              value={formData.confirm_password}
              onChange={(e) => setFormData({ ...formData, confirm_password: e.target.value })}
              required
              minLength="6"
            />
          </label>
          <button type="submit" disabled={loading} style={{ marginTop: '10px' }}>
            {loading ? 'Resetting...' : 'Reset Password'}
          </button>
        </form>
        <div style={{ marginTop: '15px', textAlign: 'center' }}>
          <Link to="/login" style={{ color: '#007bff', textDecoration: 'none' }}>Back to Login</Link>
        </div>
      </div>
    </main>
  );
}