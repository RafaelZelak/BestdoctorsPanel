// src/composables/useSessions.js
import { ref, onMounted, onUnmounted } from 'vue';
import { fetchSessions } from '@/api';
import { toTime } from '@/utils/formatters';

export function useSessions() {
  const sessions = ref([]);
  const loading = ref(false);
  const error = ref('');

  const POLL_MS_DEFAULT = 5000;
  const isPolling = ref(false);
  const pollMs = ref(POLL_MS_DEFAULT);
  let timer = null;
  let controller = null;

  async function loadSessions() {
    if (loading.value) return;
    error.value = '';
    controller?.abort();
    controller = new AbortController();
    loading.value = true;
    try {
      const data = await fetchSessions(controller.signal);
      sessions.value = (data || []).slice().sort((a, b) => {
        const ta = toTime(a.last_message_at);
        const tb = toTime(b.last_message_at);
        return tb - ta;
      });
      pollMs.value = POLL_MS_DEFAULT;
    } catch (e) {
      if (e.name !== 'AbortError') {
        error.value = `Failed to load sessions: ${e.message}`;
        pollMs.value = Math.min(pollMs.value * 2, 60000);
      }
    } finally {
      loading.value = false;
    }
  }

  function scheduleNext() {
    clearTimeout(timer);
    if (!isPolling.value) return;
    timer = setTimeout(async () => {
      if (!isPolling.value) return;
      if (document.hidden) {
        scheduleNext();
        return;
      }
      await loadSessions();
      scheduleNext();
    }, pollMs.value);
  }

  function startPolling() {
    if (isPolling.value) return;
    isPolling.value = true;
    pollMs.value = POLL_MS_DEFAULT;
    loadSessions().finally(scheduleNext);
  }

  function stopPolling() {
    isPolling.value = false;
    clearTimeout(timer);
    controller?.abort();
  }

  function onVisibilityChange() {
    if (!isPolling.value) return;
    if (!document.hidden) {
      clearTimeout(timer);
      loadSessions().finally(scheduleNext);
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
    sessions,
    loading,
    error,
    loadSessions,
    isPolling,
    startPolling,
    stopPolling,
  };
}