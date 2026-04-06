<template>
  <div class="iza-chat-layout">
    <div class="iza-chat-root">
    <!-- Header -->
    <header class="iza-header">
      <div class="iza-header-left">
        <div class="iza-avatar-badge">
          <svg width="20" height="20" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M8 10h.01M12 10h.01M16 10h.01M9 16H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-5l-5 5v-5z" />
          </svg>
        </div>
        <div>
          <h1 class="iza-title">IZA Chat</h1>
        </div>
      </div>
      <div class="iza-header-right">
        <span class="iza-session-id" :title="conversationId">
          Session: {{ conversationId.slice(0, 8) }}...
        </span>
        <button class="iza-btn-logout" @click="handleLogout">
          <svg width="16" height="16" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
          </svg>
          Sair
        </button>
      </div>
    </header>

    <!-- Messages area -->
    <main class="iza-messages" ref="messagesContainer">
      <!-- Empty state removed to avoid mixing with typing indicator -->

      <!-- Messages -->
      <template v-for="(msg, index) in messages" :key="index">
        <!-- User message -->
        <div v-if="msg.role === 'user'" class="iza-bubble-row iza-bubble-row--user">
          <div class="iza-bubble iza-bubble--user">
            <p>{{ msg.text }}</p>
          </div>
          <div class="iza-bubble-avatar iza-bubble-avatar--user">
            <svg width="14" height="14" fill="currentColor" viewBox="0 0 24 24">
              <path d="M12 12c2.7 0 4.8-2.1 4.8-4.8S14.7 2.4 12 2.4 7.2 4.5 7.2 7.2 9.3 12 12 12zm0 2.4c-3.2 0-9.6 1.6-9.6 4.8v2.4h19.2v-2.4c0-3.2-6.4-4.8-9.6-4.8z"/>
            </svg>
          </div>
        </div>

        <!-- IZA message -->
        <div v-else class="iza-bubble-row iza-bubble-row--iza">
          <div class="iza-bubble-avatar iza-bubble-avatar--iza">
            <svg width="14" height="14" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z" />
            </svg>
          </div>
          <div class="iza-bubble iza-bubble--iza">
            <p>{{ msg.text }}</p>
          </div>
        </div>
      </template>

      <!-- Loading indicator -->
      <div v-if="loading" class="iza-bubble-row iza-bubble-row--iza">
        <div class="iza-bubble-avatar iza-bubble-avatar--iza">
          <svg width="14" height="14" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z" />
          </svg>
        </div>
        <div class="iza-bubble iza-bubble--iza iza-bubble--typing">
          <span></span>
          <span></span>
          <span></span>
        </div>
      </div>

      <!-- Error banner -->
      <div v-if="errorMessage" class="iza-error-banner">
        <svg width="16" height="16" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        {{ errorMessage }}
      </div>
    </main>

    <!-- Input area -->
    <footer class="iza-input-area">
      <div class="iza-input-wrapper" @click="inputRef?.focus()">
        <textarea
          id="iza-chat-input"
          ref="inputRef"
          v-model="currentMessage"
          class="iza-input"
          placeholder="Digite sua mensagem..."
          rows="1"
          :disabled="loading"
          @keydown.enter.exact.prevent="handleSend"
          @input="autoResizeTextarea"
        ></textarea>
        <button
          id="iza-chat-send"
          class="iza-send-btn"
          :disabled="loading || !currentMessage.trim()"
          @click="handleSend"
        >
          <svg v-if="!loading" width="20" height="20" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8" />
          </svg>
          <svg v-else class="iza-spin" width="20" height="20" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor"
              d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z">
            </path>
          </svg>
        </button>
      </div>
      <p class="iza-footer-hint">Enter para enviar · Shift+Enter para nova linha</p>
    </footer>
    </div>

    <!-- Extractor Sidebar -->
    <aside class="iza-extractor-sidebar" :class="{ 'is-collapsed': isSidebarCollapsed }">
      <div class="extractor-header" @click="isSidebarCollapsed = !isSidebarCollapsed">
        <h3 v-if="!isSidebarCollapsed">Webhooks Recebidos</h3>
        <svg v-if="isSidebarCollapsed" width="20" height="20" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" /></svg>
        <svg v-else width="20" height="20" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" /></svg>
      </div>
      <div class="extractor-content" v-show="!isSidebarCollapsed">
        <div v-if="extractorPayloads.length === 0" class="extractor-empty">Nenhum webhook recebido</div>
        <div v-for="(w, i) in extractorPayloads" :key="i" class="extractor-card">
          <div class="extractor-time">{{ new Date(w.timestamp).toLocaleTimeString() }}</div>
          <pre class="extractor-json">{{ JSON.stringify(w.data, null, 2) }}</pre>
        </div>
      </div>
    </aside>
  </div>
</template>

<script setup>
import { ref, nextTick, onMounted, onBeforeUnmount } from 'vue'
import { useRouter } from 'vue-router'
import { logout } from '@/api/auth'
import { sendIzaChatMessage, pollIzaChatResponse, pollExtractorWebhook } from './api.js'

const router = useRouter()

const isSidebarCollapsed = ref(true)
const extractorPayloads = ref([])
let pollInterval = null

const conversationId = ref(generateUUID())
const messages = ref([])
const currentMessage = ref('')
const loading = ref(true)
const errorMessage = ref('')
const messagesContainer = ref(null)
const inputRef = ref(null)

function generateUUID() {
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, (character) => {
    const randomValue = (Math.random() * 16) | 0
    const resolvedValue = character === 'x' ? randomValue : (randomValue & 0x3) | 0x8
    return resolvedValue.toString(16)
  })
}

async function scrollToBottom() {
  await nextTick()
  if (messagesContainer.value) {
    messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
  }
}

function autoResizeTextarea() {
  const textarea = inputRef.value
  if (!textarea) return
  // Reset height to correctly calculate scrollHeight
  textarea.style.height = '40px'
  textarea.style.height = Math.min(Math.max(textarea.scrollHeight, 40), 160) + 'px'
}

async function handleSend() {
  const trimmedMessage = currentMessage.value.trim()
  if (!trimmedMessage || loading.value) return

  errorMessage.value = ''
  messages.value.push({ role: 'user', text: trimmedMessage })
  currentMessage.value = ''
  loading.value = true

  if (inputRef.value) {
    inputRef.value.style.height = 'auto'
  }

  await scrollToBottom()

  try {
    // Send message to backend (returns immediately with status=processing)
    await sendIzaChatMessage({
      message: trimmedMessage,
      conversationId: conversationId.value,
    })

    // Wait for the webhook response by long polling
    let responseMessage = ''
    while(true) {
      const pollRes = await pollIzaChatResponse(conversationId.value)
      
      if (pollRes.error === 'timeout') {
         // Long poll timeout (120s), we could try again, or just show error.
         errorMessage.value = 'Tempo de resposta excedido.'
         break
      }

      if (pollRes.message !== undefined) {
         responseMessage = pollRes.message
         break
      }
      
      if (pollRes.done === false) {
         // Keep waiting
         await new Promise((r) => setTimeout(r, 1000))
         continue
      }
    }
    
    if (responseMessage) {
      messages.value.push({ role: 'iza', text: responseMessage })
    }
  } catch (fetchError) {
    errorMessage.value = 'Falha ao obter resposta da IZA. Tente novamente.'
    console.error(fetchError)
  } finally {
    loading.value = false
    await scrollToBottom()
  }
}

async function handleLogout() {
  try {
    await logout()
    localStorage.removeItem('user')
  } catch (logoutError) {
    console.error(logoutError)
  } finally {
    router.push('/login')
  }
}

onMounted(() => {
  inputRef.value?.focus()

  // Simulate initial typing delay
  const delay = Math.floor(Math.random() * 2000) + 2000 // Between 2 to 4 seconds
  setTimeout(async () => {
    messages.value.push({
      role: 'iza',
      text: 'Oi, User! 😊 Eu sou a Iza, sua assistente jurídica. Estou aqui pra te ouvir com atenção e entender direitinho sua situação antes de acionar um advogado(a). Pode me contar o que tá te preocupando ou o que você quer resolver?'
    })
    loading.value = false
    await scrollToBottom()
  }, delay)

  pollInterval = setInterval(async () => {
    try {
      const payloads = await pollExtractorWebhook()
      if (payloads && payloads.length > 0) {
        extractorPayloads.value.push(...payloads)
      }
    } catch (e) {
      console.error('Failed to poll extractor', e)
    }
  }, 2000)
})

onBeforeUnmount(() => {
  if (pollInterval) clearInterval(pollInterval)
})
</script>

<style scoped>
.iza-chat-layout {
  display: flex;
  flex-direction: row;
  height: 100vh;
  width: 100vw;
  overflow: hidden;
  background: #0f1117;
  color: #e2e8f0;
  font-family: 'Inter', system-ui, sans-serif;
}

.iza-chat-root {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
}

/* ─── Header ─── */
.iza-header {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1rem 1.5rem;
  background: #161b27;
  border-bottom: 1px solid #1e2535;
  box-shadow: 0 1px 16px 0 rgba(0, 0, 0, 0.4);
}

.iza-header-left {
  display: flex;
  align-items: center;
  gap: 0.875rem;
}

.iza-avatar-badge {
  width: 42px;
  height: 42px;
  border-radius: 12px;
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  box-shadow: 0 4px 14px rgba(99, 102, 241, 0.4);
}

.iza-title {
  font-size: 1.1rem;
  font-weight: 700;
  color: #f1f5f9;
  letter-spacing: 0.02em;
  margin: 0;
}

.iza-subtitle {
  font-size: 0.72rem;
  color: #64748b;
  margin: 0;
  margin-top: 2px;
}

.iza-header-right {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.iza-session-id {
  font-size: 0.7rem;
  font-family: 'JetBrains Mono', monospace;
  color: #475569;
  background: #1e293b;
  padding: 0.25rem 0.6rem;
  border-radius: 6px;
  border: 1px solid #2d3748;
}

.iza-btn-logout {
  display: flex;
  align-items: center;
  gap: 0.4rem;
  padding: 0.4rem 0.9rem;
  background: transparent;
  border: 1px solid #374151;
  color: #94a3b8;
  border-radius: 8px;
  font-size: 0.82rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.18s ease;
}

.iza-btn-logout:hover {
  background: #ef444430;
  border-color: #ef4444;
  color: #fca5a5;
}

/* ─── Messages ─── */
.iza-messages {
  flex: 1;
  overflow-y: auto;
  padding: 2rem 1.5rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
  scroll-behavior: smooth;
}

.iza-messages::-webkit-scrollbar {
  width: 6px;
}

.iza-messages::-webkit-scrollbar-track {
  background: transparent;
}

.iza-messages::-webkit-scrollbar-thumb {
  background: #2d3748;
  border-radius: 99px;
}

/* Empty state */
.iza-empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  flex: 1;
  gap: 0.5rem;
  text-align: center;
  padding: 3rem 1rem;
}

.iza-empty-icon {
  width: 72px;
  height: 72px;
  border-radius: 50%;
  background: linear-gradient(135deg, #6366f120, #8b5cf220);
  border: 1px solid #6366f130;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #818cf8;
  margin-bottom: 0.5rem;
}

.iza-empty-title {
  font-size: 1.25rem;
  font-weight: 600;
  color: #e2e8f0;
  margin: 0;
}

.iza-empty-desc {
  font-size: 0.875rem;
  color: #64748b;
  max-width: 320px;
  margin: 0;
  line-height: 1.6;
}

/* Bubbles */
.iza-bubble-row {
  display: flex;
  align-items: flex-end;
  gap: 0.5rem;
  max-width: 75%;
  animation: fadeSlideIn 0.22s ease;
}

@keyframes fadeSlideIn {
  from { opacity: 0; transform: translateY(8px); }
  to   { opacity: 1; transform: translateY(0); }
}

.iza-bubble-row--user {
  align-self: flex-end;
  flex-direction: row-reverse;
}

.iza-bubble-row--iza {
  align-self: flex-start;
}

.iza-bubble-avatar {
  width: 30px;
  height: 30px;
  border-radius: 50%;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.7rem;
  font-weight: 700;
}

.iza-bubble-avatar--user {
  background: #3b82f6;
  color: #fff;
}

.iza-bubble-avatar--iza {
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
  color: #fff;
}

.iza-bubble {
  padding: 0.75rem 1rem;
  border-radius: 16px;
  line-height: 1.6;
  font-size: 0.9rem;
}

.iza-bubble p {
  margin: 0;
  white-space: pre-wrap;
  word-break: break-word;
}

.iza-bubble--user {
  background: linear-gradient(135deg, #3b82f6, #2563eb);
  color: #fff;
  border-bottom-right-radius: 4px;
  box-shadow: 0 2px 12px rgba(59, 130, 246, 0.25);
}

.iza-bubble--iza {
  background: #1e2535;
  color: #cbd5e1;
  border-bottom-left-radius: 4px;
  border: 1px solid #2d3748;
}

/* Typing animation */
.iza-bubble--typing {
  display: flex;
  align-items: center;
  gap: 5px;
  padding: 0.85rem 1.1rem;
}

.iza-bubble--typing span {
  display: inline-block;
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: #6366f1;
  animation: typingBounce 1.2s infinite ease-in-out;
}

.iza-bubble--typing span:nth-child(2) { animation-delay: 0.2s; }
.iza-bubble--typing span:nth-child(3) { animation-delay: 0.4s; }

@keyframes typingBounce {
  0%, 60%, 100% { transform: translateY(0); opacity: 0.5; }
  30%           { transform: translateY(-6px); opacity: 1; }
}

/* Error */
.iza-error-banner {
  align-self: center;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  background: #7f1d1d30;
  border: 1px solid #ef444450;
  color: #fca5a5;
  padding: 0.5rem 1rem;
  border-radius: 10px;
  font-size: 0.82rem;
  animation: fadeSlideIn 0.2s ease;
}

/* ─── Input area ─── */
.iza-input-area {
  flex-shrink: 0;
  padding: 1rem 1.5rem 1.25rem;
  background: #161b27;
  border-top: 1px solid #1e2535;
}

.iza-input-wrapper {
  display: flex;
  align-items: flex-end;
  gap: 0.75rem;
  background: #1e2535;
  border: 1px solid #2d3748;
  border-radius: 14px;
  padding: 0.6rem 0.6rem 0.6rem 1rem;
  transition: border-color 0.18s;
}

.iza-input-wrapper:focus-within {
  border-color: #6366f1;
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.12);
}

.iza-input {
  flex: 1;
  background: transparent;
  border: none;
  outline: none;
  color: #e2e8f0;
  font-size: 0.9rem;
  line-height: 1.5;
  resize: none;
  min-height: 40px;
  padding: 9px 0;
  max-height: 160px;
  font-family: inherit;
  box-sizing: border-box;
}

.iza-input::placeholder {
  color: #475569;
}

.iza-input:disabled {
  opacity: 0.6;
}

.iza-send-btn {
  flex-shrink: 0;
  width: 40px;
  height: 40px;
  border-radius: 10px;
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
  border: none;
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.18s ease;
  box-shadow: 0 2px 8px rgba(99, 102, 241, 0.35);
}

.iza-send-btn:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 4px 14px rgba(99, 102, 241, 0.5);
}

.iza-send-btn:disabled {
  opacity: 0.45;
  cursor: not-allowed;
  transform: none;
  box-shadow: none;
}

.iza-footer-hint {
  font-size: 0.7rem;
  color: #374151;
  text-align: center;
  margin: 0.5rem 0 0 0;
}

/* Spin animation */
.iza-spin {
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to   { transform: rotate(360deg); }
}

/* ─── Sidebar ─── */
.iza-extractor-sidebar {
  width: 320px;
  flex-shrink: 0;
  background: #161b27;
  border-left: 1px solid #1e2535;
  display: flex;
  flex-direction: column;
  transition: width 0.3s ease;
  overflow: hidden;
}

.iza-extractor-sidebar.is-collapsed {
  width: 50px;
}

.extractor-header {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 1rem;
  border-bottom: 1px solid #1e2535;
  cursor: pointer;
  background: #1a202c;
}
.extractor-header:hover {
  background: #2d3748;
}
.is-collapsed .extractor-header {
  justify-content: center;
  padding: 0;
}
.extractor-header h3 {
  font-size: 0.9rem;
  margin: 0;
  color: #f1f5f9;
}

.extractor-content {
  flex: 1;
  overflow-y: auto;
  padding: 1rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.extractor-empty {
  text-align: center;
  font-size: 0.8rem;
  color: #64748b;
  margin-top: 2rem;
}

.extractor-card {
  background: #1e2535;
  border: 1px solid #2d3748;
  border-radius: 8px;
  padding: 0.75rem;
}

.extractor-time {
  font-size: 0.7rem;
  color: #94a3b8;
  margin-bottom: 0.5rem;
  text-align: right;
}

.extractor-json {
  font-family: 'JetBrains Mono', monospace;
  font-size: 0.75rem;
  margin: 0;
  color: #a5b4fc;
  background: #0f1117;
  padding: 0.5rem;
  border-radius: 4px;
  overflow-x: auto;
}
</style>
