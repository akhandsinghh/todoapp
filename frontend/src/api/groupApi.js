import client from './axios';
import { normalizeArrayResponse, unwrapApiResponse } from './responseUtils';

export const listGroups = () =>
  client.get('/groups').then((r) => normalizeArrayResponse(r.data));

export const createGroup = (payload) =>
  client.post('/groups', payload).then((r) => unwrapApiResponse(r.data));

export const updateGroup = (id, payload) =>
  client.put(`/groups/${id}`, payload).then((r) => unwrapApiResponse(r.data));

export const deleteGroup = (id) =>
  client.delete(`/groups/${id}`).then((r) => unwrapApiResponse(r.data));

export const shareGroup = (id, payload) =>
  client.post(`/groups/${id}/share`, payload).then((r) => unwrapApiResponse(r.data));
