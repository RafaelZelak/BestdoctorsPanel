// src/api/bia.js
// Use backend-prefixed path to avoid clashing with existing Bia WebChat frontend routes
const API_BASE = '/bestdoctors/bia/homol';

export async function fetchPromptsList() {
    const res = await fetch(`${API_BASE}/prompts?content=false`, {
        credentials: 'include'
    });
    if (!res.ok) {
        throw new Error(`Failed to fetch prompts list: ${res.statusText}`);
    }
    return res.json();
}

export async function fetchPromptDetails(name) {
    const res = await fetch(`${API_BASE}/prompts/${name}`, {
        credentials: 'include'
    });
    if (!res.ok) {
        throw new Error(`Failed to fetch prompt ${name}: ${res.statusText}`);
    }
    return res.json();
}

export async function updatePrompt(name, promptContent) {
    const res = await fetch(`${API_BASE}/prompts/${name}`, {
        method: 'PATCH',
        headers: {
            'Content-Type': 'application/json'
        },
        credentials: 'include',
        body: JSON.stringify({ prompt: promptContent })
    });
    if (!res.ok) {
        throw new Error(`Failed to update prompt ${name}: ${res.statusText}`);
    }
    return res.json();
}

export async function syncAllPrompts() {
    const res = await fetch(`${API_BASE}/prompts/sync`, {
        method: 'PUT',
        credentials: 'include'
    });
    if (!res.ok) {
        throw new Error(`Failed to sync all prompts: ${res.statusText}`);
    }
    return res.json();
}

export async function syncSinglePrompt(name) {
    const res = await fetch(`${API_BASE}/prompts/sync/${name}`, {
        method: 'PUT',
        credentials: 'include'
    });
    if (!res.ok) {
        throw new Error(`Failed to sync prompt ${name}: ${res.statusText}`);
    }
    return res.json();
}
