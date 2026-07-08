import { apiGet } from './client.js';

export const profileApi = {
  get: () => apiGet('/api/profile'),
};

