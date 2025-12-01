<template>
  <div class="send-bar border-t border-neutral-800 bg-neutral-900/70 backdrop-blur supports-[backdrop-filter]:bg-neutral-900/60">
    <form class="w-full px-3 py-3 flex items-center" @submit.prevent="onSubmit" :aria-busy="isSending ? 'true' : 'false'">
      <label for="message-input" class="sr-only">Mensagem</label>
      <div class="relative w-full">
        <input
          id="message-input"
          v-model="draft"
          type="text"
          name="message"
          autocomplete="off"
          placeholder="Type a message…"
          class="w-full h-12 rounded-full border border-neutral-700 bg-neutral-800 pl-5 pr-12 text-sm text-white placeholder-neutral-400 outline-none
                 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition duration-200 ease-in-out"
          aria-label="Type a message"
          :disabled="isSending"
        />
        <button
          type="submit"
          class="absolute right-2 top-1/2 -translate-y-1/2 h-8 w-8 flex items-center justify-center rounded-full 
                 text-blue-400 hover:text-blue-300 disabled:opacity-40 disabled:cursor-not-allowed
                 outline-none focus:ring-2 focus:ring-blue-500 transition duration-200 ease-in-out"
          :disabled="!canSend || isSending"
          aria-label="Send message"
          title="Send"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="currentColor" viewBox="0 0 24 24">
            <path d="M2.01 21 23 12 2.01 3 2 10l15 2-15 2z"/>
          </svg>
        </button>
      </div>
    </form>
    <p v-if="errorMsg" class="px-4 pb-3 text-sm text-red-400">{{ errorMsg }}</p>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { useChatHistory } from '@/composables/useChatHistory'
import { sendMessage } from '@/api'

const props = defineProps({
  sessionId: { type: [String, Number], required: true },
  to: { type: String, required: true }
})

const draft = ref('')
const errorMsg = ref('')
const isSending = ref(false)

const { messages, startPolling, stopPolling } = useChatHistory()

watch(
  () => props.sessionId,
  (id) => {
    stopPolling()
    messages.value = []
    if (id) startPolling(String(id))
  },
  { immediate: true }
)

onMounted(() => {
  if (props.sessionId) startPolling(String(props.sessionId))
})

onUnmounted(() => {
  stopPolling()
})

const latestVars = computed(() => {
  const arr = messages.value
  if (!arr || !arr.length) return {}
  const last = arr[arr.length - 1]
  return last?.message?.content?.output?.vars ?? {}
})

function normalizeBrazilPhone(raw) {
  let v = String(raw || '').trim()
  if (v.startsWith('whatsapp:')) v = v.slice('whatsapp:'.length)
  const digits = v.replace(/[^\d+]/g, '')
  let d = digits.startsWith('+') ? digits : `+${digits}`
  if (!d.startsWith('+55')) {
    if (d.startsWith('+')) d = `+55${d.slice(1)}`
    else d = `+55${d.replace(/^\+/, '')}`
  }
  const rest = d.replace(/^\+55/, '')
  if (!/^\d{10,11}$/.test(rest)) return null
  return `whatsapp:${d}`
}

const normalizedTo = computed(() => normalizeBrazilPhone(props.to))

const canSend = computed(() => {
  if (!draft.value.trim()) return false
  if (!normalizedTo.value) return false
  return true
})

async function onSubmit() {
  if (!canSend.value || isSending.value) return
  errorMsg.value = ''
  isSending.value = true
  try {
    const payload = {
      to: normalizedTo.value,
      message: draft.value.trim(),
      sessionId: String(props.sessionId),
      vars: {
        nome_do_lead: latestVars.value?.nome_do_lead ?? '',
        numero_de_usuarios: latestVars.value?.numero_de_usuarios ?? '',
        especialidade: latestVars.value?.especialidade ?? '',
        funcionalidades_desejadas: latestVars.value?.funcionalidades_desejadas ?? '',
        plano: latestVars.value?.plano ?? '',
        desinteresse: latestVars.value?.desinteresse ?? '',
        saudacao_enviada: latestVars.value?.saudacao_enviada ?? '',
        finalizar: latestVars.value?.finalizar ?? ''
      }
    }
    await sendMessage(payload)
    draft.value = ''
  } catch (e) {
    errorMsg.value = e?.message || 'Failed to send message'
  } finally {
    isSending.value = false
  }
}
</script>

<style scoped>
.send-bar {
  position: sticky;
  bottom: 0;
}
</style>
