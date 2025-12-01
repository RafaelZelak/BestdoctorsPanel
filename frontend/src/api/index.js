// src/api/index.js
const API_BASE = import.meta.env.VITE_API_BASE || '';

async function safeFetch(url, opts = {}) {
  const res = await fetch(url, {
    ...opts,
    credentials: 'include' // CRITICAL: Include session cookies
  });

  // Handle 401 Unauthorized - redirect to login
  if (res.status === 401) {
    window.location.href = '/login';
    throw new Error('Unauthorized');
  }

  if (!res.ok) {
    let detail = '';
    const ct = res.headers.get('content-type') || '';
    try {
      if (ct.includes('application/json')) {
        const data = await res.json();
        detail = data?.message || data?.error || JSON.stringify(data);
      } else {
        detail = await res.text();
      }
    } catch {
      // ignore parse errors
    }
    const msg = `${res.status} ${res.statusText}${detail ? ` — ${detail}` : ''}`;
    throw new Error(msg);
  }
  // try json first
  const ct = res.headers.get('content-type') || '';
  if (ct.includes('application/json')) return res.json();
  return res.text();
}

export async function fetchSessions(signal) {
  return safeFetch(`${API_BASE}/bestdoctors/sessionphone`, { signal });
}

export async function fetchHistory(sessionId, signal) {
  return safeFetch(
    `${API_BASE}/bestdoctors/chathistory?session_id=${encodeURIComponent(sessionId)}`,
    { signal }
  );
}

export async function toggleSessionActive(sessionId, signal) {
  return safeFetch(
    `${API_BASE}/bestdoctors/sessionphone/active?session_id=${encodeURIComponent(sessionId)}`,
    { method: 'PATCH', signal }
  );
}

export async function sendMessage({ to, message, sessionId, vars }, signal) {
  return safeFetch(`${API_BASE}/bestdoctors/sendmessage`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    signal,
    body: JSON.stringify({
      to,
      message,
      session_id: String(sessionId),
      vars: {
        nome_do_lead: vars?.nome_do_lead ?? '',
        numero_de_usuarios: vars?.numero_de_usuarios ?? '',
        especialidade: vars?.especialidade ?? '',
        funcionalidades_desejadas: vars?.funcionalidades_desejadas ?? '',
        plano: vars?.plano ?? '',
        desinteresse: vars?.desinteresse ?? '',
        saudacao_enviada: vars?.saudacao_enviada ?? '',
        finalizar: vars?.finalizar ?? ''
      }
    })
  })
}

export async function requestReport({ report, type, filters }, signal) {
  return safeFetch(`${API_BASE}/bestdoctors/report`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    signal,
    body: JSON.stringify({
      report,            // "session" | "abandonment" | "flowDepth" | "reengagement"
      type,              // "json"
      filters: {
        from: filters?.from ?? null,   // ISO string "YYYY-MM-DDTHH:mm:ssZ" or null
        to: filters?.to ?? null,       // ISO string "YYYY-MM-DDTHH:mm:ssZ" or null
        full: Boolean(filters?.full ?? false),
      },
    }),
  });
}

/**
 * Use this for file downloads (csv/pdf/xlsx). It returns { blob, filename }.
 * We do not use safeFetch here because we need response.blob().
 */
export async function downloadReport({ report, type, filters }, signal) {
  const res = await fetch(`${API_BASE}/bestdoctors/report`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    signal,
    body: JSON.stringify({
      report,
      type, // "csv" | "pdf" | "xlsx"
      filters: {
        from: filters?.from ?? null,
        to: filters?.to ?? null,
        full: Boolean(filters?.full ?? false),
      },
    }),
  });

  if (!res.ok) {
    let detail = '';
    try {
      const ct = res.headers.get('content-type') || '';
      if (ct.includes('application/json')) {
        const data = await res.json();
        detail = data?.message || data?.error || JSON.stringify(data);
      } else {
        detail = await res.text();
      }
    } catch { /* ignore */ }
    throw new Error(`${res.status} ${res.statusText}${detail ? ` — ${detail}` : ''}`);
  }

  const blob = await res.blob();
  const cd = res.headers.get('content-disposition') || '';
  const filenameMatch = cd.match(/filename\*=UTF-8''([^;]+)|filename="?([^"]+)"?/i);
  const filename = decodeURIComponent(filenameMatch?.[1] || filenameMatch?.[2] || `report-${report}.${type}`);
  return { blob, filename };
}

export { API_BASE };

