// src/api/auth.js
const API_BASE = import.meta.env.VITE_API_BASE || '';

export async function login(username, password) {
    const res = await fetch(`${API_BASE}/auth/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include', // CRITICAL for cookies
        body: JSON.stringify({ username, password })
    });

    if (!res.ok) {
        const error = await res.text();
        throw new Error(error || 'Login failed');
    }

    return res.json();
}

export async function logout() {
    const res = await fetch(`${API_BASE}/auth/logout`, {
        method: 'POST',
        credentials: 'include'
    });

    return res.json();
}

export async function getCurrentUser() {
    const res = await fetch(`${API_BASE}/auth/me`, {
        credentials: 'include'
    });

    if (!res.ok) {
        throw new Error('Not authenticated');
    }

    return res.json();
}

export async function checkAuth() {
    try {
        await getCurrentUser();
        return true;
    } catch {
        return false;
    }
}
