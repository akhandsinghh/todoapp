import client from './axios';
import { unwrapApiResponse } from './responseUtils';

const unwrapAuthResponse = (response) => {
  const payload = unwrapApiResponse(response);
  if (payload && typeof payload === 'object' && !Array.isArray(payload)) {
    return payload;
  }
  return payload;
};

export const register = (payload) =>
  client.post('/auth/register', payload).then((r) => unwrapAuthResponse(r.data));

export const login = (payload) =>
  client.post('/auth/login', payload).then((r) => unwrapAuthResponse(r.data));

export const me = () => client.get('/auth/me').then((r) => unwrapAuthResponse(r.data));

export const forgotPassword = (payload) =>
  client.post('/auth/forgot-password', payload).then((r) => unwrapAuthResponse(r.data));

export const changePassword = (payload) =>
  client.post('/auth/change-password', payload).then((r) => unwrapAuthResponse(r.data));