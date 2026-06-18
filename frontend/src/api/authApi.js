import client from './axios';

export const register = (payload) => client.post('/auth/register', payload).then((r) => r.data);
export const login = (payload) => client.post('/auth/login', payload).then((r) => r.data);
export const me = () => client.get('/auth/me').then((r) => r.data);

// NEW API CALLS
export const forgotPassword = (payload) => client.post('/auth/forgot-password', payload).then((r) => r.data);
export const changePassword = (payload) => client.post('/auth/change-password', payload).then((r) => r.data);