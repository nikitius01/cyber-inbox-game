import { createContext, useContext, useEffect, useMemo, useState } from 'react';
import { apiGet, apiPost } from '../api/client.js';

const AuthContext = createContext(null);

export function AuthProvider({ children }) {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    apiGet('/api/auth/me')
      .then(setUser)
      .catch(() => setUser(null))
      .finally(() => setLoading(false));
  }, []);

  async function login(email, password) {
    const data = await apiPost('/api/auth/login', { email, password });
    setUser(data.user);
  }

  async function register(username, email, password) {
    const data = await apiPost('/api/auth/register', { username, email, password });
    setUser(data.user);
  }

  async function logout() {
    await apiPost('/api/auth/logout', {});
    setUser(null);
  }

  const value = useMemo(() => ({ user, loading, login, register, logout, setUser }), [user, loading]);
  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  return useContext(AuthContext);
}

