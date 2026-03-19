// src/api/prompts.js

/**
 * Creates an API client for prompt management targeting the specified base URL.
 * @param {string} apiBase - The base URL prefix (e.g., '/digesac/homol' or '/bestdoctors/bia/homol')
 */
export function createPromptApi(apiBase) {
    return {
        async fetchPromptsList() {
            const res = await fetch(`${apiBase}/prompts?content=false`, {
                credentials: 'include'
            });
            if (!res.ok) {
                throw new Error(`Failed to fetch prompts list: ${res.statusText}`);
            }
            return res.json();
        },
        
        async fetchPromptDetails(name) {
            const res = await fetch(`${apiBase}/prompts/${name}`, {
                credentials: 'include'
            });
            if (!res.ok) {
                throw new Error(`Failed to fetch prompt ${name}: ${res.statusText}`);
            }
            return res.json();
        },
        
        async updatePrompt(name, promptContent) {
            const res = await fetch(`${apiBase}/prompts/${name}`, {
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
        },
        
        async syncAllPrompts() {
            const res = await fetch(`${apiBase}/prompts/sync`, {
                method: 'PUT',
                credentials: 'include'
            });
            if (!res.ok) {
                throw new Error(`Failed to sync all prompts: ${res.statusText}`);
            }
            return res.json();
        },
        
        async syncSinglePrompt(name) {
            const res = await fetch(`${apiBase}/prompts/sync/${name}`, {
                method: 'PUT',
                credentials: 'include'
            });
            if (!res.ok) {
                throw new Error(`Failed to sync prompt ${name}: ${res.statusText}`);
            }
            return res.json();
        },
        
        async createPrompt(name, promptContent) {
            const res = await fetch(`${apiBase}/prompts`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                credentials: 'include',
                body: JSON.stringify({ name, prompt: promptContent })
            });
            if (!res.ok) {
                throw new Error(`Failed to create prompt: ${res.statusText}`);
            }
            return res.json();
        },
        
        async deletePrompt(name) {
            const res = await fetch(`${apiBase}/prompts/${name}`, {
                method: 'DELETE',
                credentials: 'include'
            });
            if (!res.ok) {
                throw new Error(`Failed to delete prompt ${name}: ${res.statusText}`);
            }
            return res.json();
        }
    }
}
