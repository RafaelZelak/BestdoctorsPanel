const API_BASE = import.meta.env.VITE_API_BASE || '';

export async function sendIzaChatMessage({ message, conversationId }) {
  const res = await fetch(`${API_BASE}/api/iza-chat/message`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
    body: JSON.stringify({
      message,
      conversation_id: conversationId,
    }),
  });

  if (res.status === 401) {
    window.location.href = '/login';
    throw new Error('Unauthorized');
  }

  if (!res.ok) {
    const detail = await res.text().catch(() => '');
    throw new Error(`${res.status} ${res.statusText}${detail ? ` — ${detail}` : ''}`);
  }

  return res.json();
}

export async function pollIzaChatResponse(conversationId) {
  const res = await fetch(`${API_BASE}/api/iza-chat/response?conversation_id=${encodeURIComponent(conversationId)}`, {
    method: 'GET',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
  });

  if (res.status === 401) {
    window.location.href = '/login';
    throw new Error('Unauthorized');
  }

  if (!res.ok) {
    const detail = await res.text().catch(() => '');
    throw new Error(`${res.status} ${res.statusText}${detail ? ` — ${detail}` : ''}`);
  }

	return res.json();
}

export async function pollExtractorWebhook() {
  const res = await fetch(`${API_BASE}/api/iza-chat-extractor/poll?_t=${Date.now()}`, {
    method: 'GET',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
  });

  if (res.status === 401) {
    window.location.href = '/login';
    throw new Error('Unauthorized');
  }

  if (!res.ok) {
    const detail = await res.text().catch(() => '');
    throw new Error(`${res.status} ${res.statusText}${detail ? ` — ${detail}` : ''}`);
  }

  return res.json();
}
