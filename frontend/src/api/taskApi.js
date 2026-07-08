import { apiGet, apiPost } from './client.js';

export const taskApi = {
  random: (category) => apiGet(`/api/tasks/random?category=${encodeURIComponent(category)}`),
  answer: (task, answer) => {
    if (task.category === 'AI') {
      return apiPost('/api/ai/tasks/answer', { taskId: task.id, answer });
    }
    return apiPost(`/api/tasks/${encodeURIComponent(task.id)}/answer`, { answer });
  },
};

