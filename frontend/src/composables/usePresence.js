import { ref, onMounted, onUnmounted } from 'vue'

const WS_PATH = '/bestdoctors/presence/ws'
const RECONNECT_BASE_MS = 1000
const RECONNECT_MAX_MS = 30000

export function usePresence(scope) {
  const allPresence = ref({})
  const currentPromptName = ref(null)
  const currentUsername = ref(null)

  let socket = null
  let reconnectAttempts = 0
  let reconnectTimeoutId = null
  let intentionallyClosed = false

  function buildWebSocketURL() {
    const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws'
    return `${protocol}://${window.location.host}${WS_PATH}`
  }

  function send(payload) {
    if (socket && socket.readyState === WebSocket.OPEN) {
      socket.send(JSON.stringify(payload))
    }
  }

  function publishPresence(promptName, status) {
    currentPromptName.value = promptName
    send({ scope, prompt: promptName, status })
  }

  function clearPresence() {
    if (currentPromptName.value) {
      send({ scope, prompt: currentPromptName.value, status: 'gone' })
    }
    currentPromptName.value = null
  }

  function resolveCurrentUsername() {
    try {
      const storedUser = localStorage.getItem('user')
      if (storedUser) {
        const parsed = JSON.parse(storedUser)
        return parsed.username || null
      }
    } catch {
      return null
    }
    return null
  }

  function handleIncomingSnapshot(rawSnapshot) {
    const username = currentUsername.value
    const filtered = {}

    for (const [promptKey, users] of Object.entries(rawSnapshot)) {
      const [keyScope, ...promptParts] = promptKey.split(':')
      if (keyScope !== scope) continue

      const promptName = promptParts.join(':')
      const otherUsers = Object.values(users).filter(
        (userPresence) => userPresence.username !== username
      )

      if (otherUsers.length > 0) {
        filtered[promptName] = otherUsers
      }
    }

    allPresence.value = filtered
  }

  function connect() {
    if (intentionallyClosed) return

    socket = new WebSocket(buildWebSocketURL())

    socket.onopen = () => {
      reconnectAttempts = 0

      if (currentPromptName.value) {
        send({ scope, prompt: currentPromptName.value, status: 'viewing' })
      }
    }

    socket.onmessage = (event) => {
      try {
        const snapshot = JSON.parse(event.data)
        handleIncomingSnapshot(snapshot)
      } catch {
        // Ignore malformed frames silently
      }
    }

    socket.onclose = () => {
      if (intentionallyClosed) return

      const backoff = Math.min(
        RECONNECT_BASE_MS * 2 ** reconnectAttempts,
        RECONNECT_MAX_MS
      )
      reconnectAttempts++
      reconnectTimeoutId = setTimeout(connect, backoff)
    }

    socket.onerror = () => {
      socket.close()
    }
  }

  onMounted(() => {
    currentUsername.value = resolveCurrentUsername()
    connect()
  })

  onUnmounted(() => {
    intentionallyClosed = true
    clearTimeout(reconnectTimeoutId)
    clearPresence()
    if (socket) {
      socket.close()
      socket = null
    }
    allPresence.value = {}
  })

  return { allPresence, publishPresence, clearPresence }
}
