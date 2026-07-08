export async function apiGet(path) {
  const response = await fetch(path, { credentials: 'include' });
  return parseResponse(response);
}

export async function apiPost(path, payload) {
  const response = await fetch(path, {
    method: 'POST',
    credentials: 'include',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload),
  });
  return parseResponse(response);
}

async function parseResponse(response) {
  if (!response.ok) {
    const message = await response.text();
    throw new Error(message || 'Ошибка запроса');
  }
  if (response.status === 204) {
    return null;
  }
  return response.json();
}

