import { apiPost } from './client.js';

export const authApi = {
  login: (email, password) => apiPost('/api/auth/login', { email, password }),
  register: (username, email, password) => apiPost('/api/auth/register', { username, email, password }),
  logout: () => apiPost('/api/auth/logout', {}),
};

