// src/composables/useChatHistory.js
import { ref, onMounted, onUnmounted } from 'vue';
import { fetchHistory } from '@/api';

export function useChatHistory() {
  const messages = ref([]);
  const loading = ref(false);
  const error = ref('');
  const currentSessionId = ref(null);

  const POLL_MS_BASE = 3000;
  let pollMs = POLL_MS_BASE;
  let timer = null;
  let controller = null;
  const isPolling = ref(false);

  async function loadHistory(id) {
    const sid = id || currentSessionId.value;
    if (!sid) return;
    error.value = '';
    controller?.abort();
    controller = new AbortController();
    loading.value = true;
    try {
      const data = await fetchHistory(sid, controller.signal);
      messages.value = Array.isArray(data) ? data : [];
      pollMs = POLL_MS_BASE; // reset backoff
    } catch (e) {
      if (e.name !== 'AbortError') {
        error.value = `Failed to load history: ${e.message}`;
        pollMs = Math.min(pollMs * 2, 30000); // light backoff
      }
    } finally {
      loading.value = false;
    }
  }

  function scheduleNext() {
    clearTimeout(timer);
    if (!isPolling.value || !currentSessionId.value) return;
    timer = setTimeout(async () => {
      if (!isPolling.value || !currentSessionId.value) return;
      if (document.hidden) {
        // keep timer cycling but do not fetch while hidden
        scheduleNext();
        return;
      }
      await loadHistory();
      scheduleNext();
    }, pollMs);
  }

  function startPolling(sessionId) {
    stopPolling();
    currentSessionId.value = sessionId || null;
    if (!currentSessionId.value) return;
    isPolling.value = true;
    // Optional: clear to avoid showing messages from previous session for a moment
    messages.value = [];
    loadHistory(currentSessionId.value).finally(scheduleNext);
  }

  function stopPolling() {
    isPolling.value = false;
    clearTimeout(timer);
    controller?.abort();
  }

  function onVisibilityChange() {
    if (!isPolling.value || !currentSessionId.value) return;
    if (!document.hidden) {
      clearTimeout(timer);
      loadHistory().finally(scheduleNext);
    }
  }

  onMounted(() => {
    document.addEventListener('visibilitychange', onVisibilityChange);
  });

  onUnmounted(() => {
    document.removeEventListener('visibilitychange', onVisibilityChange);
    stopPolling();
  });

  return {
    messages,
    loading,
    error,
    currentSessionId,
    loadHistory,
    startPolling,
    stopPolling,
    isPolling,
  };
}
