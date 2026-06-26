import React, { createContext, useContext, useEffect, useMemo, useState } from 'react';
import * as authApi from '../api/authApi';

const AuthContext = createContext(null);
export const useAuth = () => useContext(AuthContext);

export function AuthProvider({ children }) {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const token = localStorage.getItem('token');
    if (!token) {
      setLoading(false);
      return;
    }
    authApi
      .me()
      .then((payload) => {
        const normalizedUser = payload?.user || payload;
        setUser(normalizedUser);
      })
      .catch(() => localStorage.removeItem('token'))
      .finally(() => setLoading(false));
  }, []);

  const value = useMemo(
    () => ({
      user,
      loading,
      login: async (payload) => {
        const res = await authApi.login(payload);
        const token = res?.token || res?.model?.token;
        const normalizedUser = res?.user || res?.model?.user || res?.model;
        if (token) localStorage.setItem('token', token);
        setUser(normalizedUser || null);
        return res;
      },
      register: async (payload) => {
        const res = await authApi.register(payload);
        const token = res?.token || res?.model?.token;
        const normalizedUser = res?.user || res?.model?.user || res?.model;
        if (token) localStorage.setItem('token', token);
        setUser(normalizedUser || null);
        return res;
      },
      logout: () => {
        localStorage.removeItem('token');
        setUser(null);
      },
    }),
    [user, loading]
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}