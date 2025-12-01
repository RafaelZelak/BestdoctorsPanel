<template>
  <main class="flex flex-col h-full">
    <!-- Header -->
    <header class="p-4 border-b border-neutral-800 relative">
      <div class="flex items-center justify-between">
        <!-- Session info -->
        <div class="flex flex-col text-left min-w-0">
          <span class="text-lg font-semibold truncate">
            <template v-if="session">
              {{ session.lead_name ? session.lead_name : prettyPhone(session.phone) }}
              <span
                class="ml-2 text-xs px-2 py-0.5 rounded-full border align-middle"
                :class="session.ai_active ? 'border-emerald-500 text-emerald-400' : 'border-neutral-600 text-neutral-400'"
              >
                {{ session.ai_active ? 'AI active' : 'AI off' }}
              </span>
            </template>
            <template v-else>
              Selecione uma conversa
            </template>
          </span>

          <!-- Phone under name (when lead_name exists) -->
          <span
            v-if="session && session.lead_name"
            class="text-sm text-neutral-400 truncate"
          >
            {{ prettyPhone(session.phone) }}
          </span>
        </div>

        <!-- Info menu trigger (hamburgers.css – spring) -->
        <button
          v-if="session"
          class="ml-4 inline-flex items-center justify-center w-9 h-9 rounded-lg hover:bg-neutral-800/60 focus:outline-none focus:ring-2 focus:ring-blue-500"
          @click="showInfo = !showInfo"
          aria-label="Show conversation info"
          title="Info"
        >
          <div class="hamburger hamburger--spring m-0" :class="{ 'is-active': showInfo }">
            <span class="hamburger-box">
              <span class="hamburger-inner"></span>
            </span>
          </div>
        </button>
      </div>

      <!-- Info menu -->
      <ChatInfoMenu
        v-if="showInfo"
        :session="session"
        :vars="latestVars"
        @updated="onSessionUpdated"
        @close="showInfo = false"
      />
    </header>

    <!-- Messages scroll area -->
    <section
      class="flex-1 overflow-y-auto p-4"
      ref="scroller"
      aria-live="polite"
      aria-busy="loading"
    >
      <!-- Empty state: no session -->
      <div v-if="!session" class="h-full grid place-items-center text-neutral-400">
        <p>Nenhuma conversa selecionada ;-;</p>
      </div>

      <template v-else>
        <!-- Error -->
        <div v-if="error" class="text-sm text-red-400 p-2">{{ error }}</div>

        <!-- Loading skeleton -->
        <div v-if="loading && !messages.length" class="space-y-3">
          <div class="w-2/3 max-w-[560px] bg-neutral-800 h-5 rounded animate-pulse"></div>
          <div class="w-1/2 max-w-[420px] bg-neutral-800 h-5 rounded animate-pulse"></div>
          <div class="w-3/4 max-w-[640px] bg-neutral-800 h-5 rounded animate-pulse"></div>
        </div>

        <!-- Messages -->
        <TransitionGroup
          name="fade-move"
          tag="div"
        >
          <MessageBubble
            v-for="h in messages"
            :key="h.id || h.created_at"
            :role="h.message?.type || 'ai'"
            :text="extractMessageText(h.message)"
            :timestamp="formatTime(h.created_at)"
          />
        </TransitionGroup>

        <!-- No content -->
        <div v-if="!messages.length && !loading" class="text-sm text-neutral-400 p-2">
          No messages in this session.
        </div>
      </template>
    </section>
    <SendMessage
      v-if="canSendMessage"
      :session-id="session.session_id"
      :to="session.phone"
    />
  </main>
</template>

<script setup>
import { ref, watch, nextTick, onMounted, onUnmounted, computed } from 'vue'
import { useChatHistory } from '@/composables/useChatHistory'
import { extractMessageText } from '@/utils/extractText'
import { formatTime, prettyPhone } from '@/utils/formatters'
import MessageBubble from './MessageBubble.vue'
import ChatInfoMenu from './ChatInfoMenu.vue'
import SendMessage from './SendMessage.vue'

const emit = defineEmits(['session-updated'])

const props = defineProps({
  session: { type: Object, default: null },
})

const {
  messages,
  loading,
  error,
  startPolling,
  stopPolling,
} = useChatHistory()

/* Keep a local copy so header badges update instantly after PATCH */
const localSession = ref(null)
watch(
  () => props.session,
  (s) => { localSession.value = s ? { ...s } : null },
  { immediate: true }
)
const session = computed(() => localSession.value)

/* Handle PATCH success from ChatInfoMenu */
function onSessionUpdated(updated) {
  if (!updated) return
  if (session.value && updated.session_id === session.value.session_id) {
    localSession.value = { ...session.value, ...updated }
    emit('session-updated', updated)
  }
}

/* Scroll behavior during polling (stable near bottom) */
const scroller = ref(null)

function isNearBottom(el) {
  if (!el) return true
  const threshold = 56
  return el.scrollTop + el.clientHeight >= el.scrollHeight - threshold
}

watch(
  () => props.session?.session_id,
  async (id) => {
    stopPolling()
    showInfo.value = false
    messages.value = []
    if (!id) return
    startPolling(id)
  },
  { immediate: true }
)

watch(messages, async () => {
  const el = scroller.value
  if (!el) return
  const wasNear = isNearBottom(el)
  const prevTop = el.scrollTop
  const prevHeight = el.scrollHeight
  await nextTick()
  const newHeight = el.scrollHeight
  if (wasNear) {
    el.scrollTop = el.scrollHeight
  } else {
    const delta = newHeight - prevHeight
    if (delta > 0) el.scrollTop = prevTop + delta
  }
})

onMounted(() => {
  const el = scroller.value
  if (el) el.scrollTop = el.scrollHeight
})

onUnmounted(() => {
  stopPolling()
})

/* Vars from the most recent history entry */
const latestVars = computed(() => {
  const arr = messages.value
  if (!arr || !arr.length) return null
  const last = arr[arr.length - 1]
  return last?.message?.content?.output?.vars ?? null
})

/* UI */
const showInfo = ref(false)

/* Footer input state (manual message when AI is OFF) */
const draft = ref('')
const isSending = ref(false)
const canSend = computed(() => draft.value.trim().length > 0 && !isSending.value)

function onSubmit() {
  if (!canSend.value) return
  isSending.value = true
  // No backend integration for now. Just clear input to simulate a successful send.
  setTimeout(() => {
    draft.value = ''
    isSending.value = false
  }, 0)
}

/* Regra final para exibir SendMessage */
const canSendMessage = computed(() => {
  if (!session.value) return false
  if (session.value.ai_active) return false
  if (!messages.value.length) return false

  const lastUserMsg = [...messages.value].reverse().find(
    m => ['user', 'client', 'human', 'outbound'].includes(m.message?.type)
  )
  if (!lastUserMsg) return false

  const msgTime = new Date(lastUserMsg.created_at).getTime()
  const diffHours = (Date.now() - msgTime) / (1000 * 60 * 60)

  return diffHours <= 24
})
</script>

<style scoped>
/* hamburgers spring – smaller and white lines */
.hamburger { padding: 0; }
.hamburger .hamburger-box { width: 22px; height: 16px; }
.hamburger .hamburger-inner,
.hamburger .hamburger-inner::before,
.hamburger .hamburger-inner::after {
  width: 22px;
  height: 2px;
  background-color: #fff;
}
.hamburger.is-active .hamburger-inner,
.hamburger.is-active .hamburger-inner::before,
.hamburger.is-active .hamburger-inner::after {
  background-color: #fff;
}

/* message list transitions */
.fade-move-enter-active,
.fade-move-leave-active {
  transition: opacity .18s ease, transform .18s ease;
}
.fade-move-enter-from,
.fade-move-leave-to {
  opacity: 0;
  transform: translateY(4px);
}
</style>
