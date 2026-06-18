import React, { useState } from 'react';
import { changePassword } from '../api/authApi';
import { useNavigate } from 'react-router-dom';

export default function ChangePassword({ onSuccess, onCancel }) {
  const navigate = useNavigate();
  const [formData, setFormData] = useState({
    old_password: '',
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
      return setError('New passwords do not match');
    }

    setLoading(true);
    try {
      await changePassword(formData);
      setSuccess('Password updated successfully!');
      setFormData({ old_password: '', new_password: '', confirm_password: '' });
      setTimeout(() => navigate('/'), 2000); 
      if (onSuccess) onSuccess();
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to update password');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="change-password-container" style={{ padding: '20px', maxWidth: '400px', margin: '0 auto', backgroundColor: 'white', borderRadius: '8px', boxShadow: '0 2px 4px rgba(0,0,0,0.1)' }}>
      <h2>Change Password</h2>
      {error && <div className="alert" style={{ color: 'red', marginBottom: '10px' }}>{error}</div>}
      {success && <div className="alert" style={{ color: 'green', marginBottom: '10px' }}>{success}</div>}
      
      <form onSubmit={handleSubmit} style={{ display: 'flex', flexDirection: 'column', gap: '15px' }}>
        <input
          type="password"
          placeholder="Current Password"
          value={formData.old_password}
          onChange={(e) => setFormData({ ...formData, old_password: e.target.value })}
          required
          style={{ padding: '8px', borderRadius: '4px', border: '1px solid #ccc' }}
        />
        <input
          type="password"
          placeholder="New Password"
          value={formData.new_password}
          onChange={(e) => setFormData({ ...formData, new_password: e.target.value })}
          required
          minLength="6"
          style={{ padding: '8px', borderRadius: '4px', border: '1px solid #ccc' }}
        />
        <input
          type="password"
          placeholder="Confirm New Password"
          value={formData.confirm_password}
          onChange={(e) => setFormData({ ...formData, confirm_password: e.target.value })}
          required
          minLength="6"
          style={{ padding: '8px', borderRadius: '4px', border: '1px solid #ccc' }}
        />
        <div style={{ display: 'flex', gap: '10px', marginTop: '10px' }}>
          <button type="submit" disabled={loading} style={{ padding: '10px', cursor: 'pointer', backgroundColor: '#007bff', color: 'white', border: 'none', borderRadius: '4px' }}>
            {loading ? 'Updating...' : 'Update Password'}
          </button>
          <button type="button" onClick={() => navigate('/')} style={{ padding: '10px', cursor: 'pointer', backgroundColor: '#050505', border: '1px solid #ccc', borderRadius: '4px' }}>
            Cancel
          </button>
        </div>
      </form>
    </div>
  );
}