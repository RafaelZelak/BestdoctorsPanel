// src/api/prompts.js

/**
 * Creates an API client for prompt management targeting the specified base URL.
 * @param {string} apiBase - The base URL prefix (e.g., '/digesac/homol' or '/bestdoctors/bia/homol')
 */
export function createPromptApi(apiBase) {
    const domainBase = apiBase.replace(/\/homol$/, '').replace(/\/prod$/, '');

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
        
        async updatePrompt(name, promptContent, edited_by) {
            const res = await fetch(`${apiBase}/prompts/${name}`, {
                method: 'PATCH',
                headers: {
                    'Content-Type': 'application/json'
                },
                credentials: 'include',
                body: JSON.stringify({ prompt: promptContent, edited_by })
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
        
        async createPrompt(name, promptContent, edited_by) {
            const res = await fetch(`${apiBase}/prompts`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                credentials: 'include',
                body: JSON.stringify({ name, prompt: promptContent, edited_by })
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
        },
        
        async fetchPromptHistory(name) {
            const res = await fetch(`${apiBase}/prompts/${name}/history`, {
                credentials: 'include'
            });
            if (!res.ok) {
                throw new Error(`Failed to fetch history for ${name}: ${res.statusText}`);
            }
            return res.json();
        },
        
        async fetchPromptDiff(name, version) {
            const res = await fetch(`${apiBase}/prompts/${name}/history/${version}`, {
                credentials: 'include'
            });
            if (!res.ok) {
                throw new Error(`Failed to fetch diff for ${name} version ${version}: ${res.statusText}`);
            }
            return res.json();
        },

        // Folder Methods
        async fetchFolders() {
            const res = await fetch(`${domainBase}/folders`, {
                credentials: 'include'
            });
            if (!res.ok) {
                throw new Error(`Failed to fetch folders: ${res.statusText}`);
            }
            return res.json();
        },
        
        async createFolder(folderName, folderIdx = null) {
            const payload = { folder_name: folderName };
            if (folderIdx !== null) {
                payload.folder_idx = folderIdx;
            }
            const res = await fetch(`${domainBase}/folders`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                credentials: 'include',
                body: JSON.stringify(payload)
            });
            if (!res.ok) {
                throw new Error(`Failed to create folder: ${res.statusText}`);
            }
            return res.json();
        },
        
        async updateFolder(id, params) {
            const res = await fetch(`${domainBase}/folders/${id}`, {
                method: 'PATCH',
                headers: { 'Content-Type': 'application/json' },
                credentials: 'include',
                body: JSON.stringify(params)
            });
            if (!res.ok) {
                throw new Error(`Failed to update folder: ${res.statusText}`);
            }
            return res.json();
        },
        
        async deleteFolder(id) {
            const res = await fetch(`${domainBase}/folders/${id}`, {
                method: 'DELETE',
                credentials: 'include'
            });
            if (!res.ok) {
                throw new Error(`Failed to delete folder: ${res.statusText}`);
            }
            return res.json();
        },
        
        async assignPromptToFolder(name, folderId) {
            const res = await fetch(`${apiBase}/prompts/${name}/folder`, {
                method: 'PUT',
                headers: { 'Content-Type': 'application/json' },
                credentials: 'include',
                body: JSON.stringify({ folder_id: folderId })
            });
            if (!res.ok) {
                throw new Error(`Failed to assign prompt: ${res.statusText}`);
            }
            return res.json();
        }
    }
}
