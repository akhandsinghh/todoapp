import client from './axios';
import { normalizeListResponse, unwrapApiResponse } from './responseUtils';

export const listTasks = (params = {}) =>
  client.get('/tasks', { params }).then((r) => normalizeListResponse(r.data));

export const createTask = (payload) =>
  client.post('/tasks', payload).then((r) => unwrapApiResponse(r.data));

export const updateTask = (id, payload) =>
  client.put(`/tasks/${id}`, payload).then((r) => unwrapApiResponse(r.data));

export const deleteTask = (id) =>
  client.delete(`/tasks/${id}`).then((r) => unwrapApiResponse(r.data));

export const createReminder = (payload) =>
  client.post('/reminders', payload).then((r) => unwrapApiResponse(r.data));

export const listReminders = () =>
  client.get('/reminders').then((r) => normalizeListResponse(r.data));
