// src/utils/formatters.js
export function toTime(iso) {
  if (!iso || iso.startsWith('0001-01-01')) return 0
  const t = Date.parse(iso)
  return Number.isNaN(t) ? 0 : t
}

export function formatTime(iso) {
  const t = toTime(iso)
  if (!t) return '—'
  return new Date(t).toLocaleString()
}

export function relativeTime(iso) {
  const t = toTime(iso)
  if (!t) return '—'
  const rtf = new Intl.RelativeTimeFormat(undefined, { numeric: 'auto' })
  const diff = Date.now() - t
  const mins = Math.round(diff / 60000)
  if (Math.abs(mins) < 60) return rtf.format(-mins, 'minute')
  const hours = Math.round(mins / 60)
  if (Math.abs(hours) < 24) return rtf.format(-hours, 'hour')
  const days = Math.round(hours / 24)
  return rtf.format(-days, 'day')
}

export function prettyPhone(p) {
  if (!p) return '-'

  const raw = String(p).trim()
  const hasPlus = raw.startsWith('+')
  const digits = raw.replace(/\D/g, '')

  // Detect Brazil (+55 or 55)
  let isBR = false
  let cc = ''
  let rest = digits

  if ((hasPlus && digits.startsWith('55')) || (!hasPlus && digits.startsWith('55') && digits.length >= 12)) {
    isBR = true
    cc = '+55'
    rest = digits.slice(2)
  } else if (digits.length === 10 || digits.length === 11) {
    // Assume BR when local sized 10/11 without explicit country code
    isBR = true
    rest = digits
  }

  if (isBR) {
    // Remove possíveis zeros à esquerda no local (DDD nunca começa com 0)
    if (rest.length >= 2 && rest[0] === '0') rest = rest.slice(1)

    const ddd = rest.slice(0, 2)
    const local = rest.slice(2)

    if (local.length === 9) {
      // (AA) 9XXXX-XXXX
      return `${cc ? cc + ' ' : ''}(${ddd}) ${local.slice(0, 5)}-${local.slice(5)}`
    }
    if (local.length === 8) {
      // (AA) XXXX-XXXX
      return `${cc ? cc + ' ' : ''}(${ddd}) ${local.slice(0, 4)}-${local.slice(4)}`
    }

    // Se tamanho inesperado, volta para fallback genérico abaixo
  }

  // Fallback genérico (ex.: outros países): agrupa em blocos legíveis
  const grouped = digits.replace(/(\d{3,4})(?=\d)/g, '$1 ').trim()
  return (hasPlus ? '+' : '') + grouped
}