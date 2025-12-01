// src/utils/extractText.js
export function extractMessageText(message) {
  if (!message) return ''
  const c = message.content
  if (typeof c === 'string') return c
  if (c && typeof c === 'object') {
    if (c.output && typeof c.output.message === 'string') {
      return c.output.message
    }
    try {
      return JSON.stringify(c, null, 2)
    } catch {
      return String(c)
    }
  }
  return ''
}
