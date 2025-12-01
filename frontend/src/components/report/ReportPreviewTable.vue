<template>
  <div>
    <!-- Sessions: list the first 5 entries -->
    <table
      v-if="reportType === 'session' && Array.isArray(data)"
      class="w-full text-left border border-neutral-800 rounded-xl overflow-hidden"
      aria-label="Preview table for session report"
    >
      <thead class="bg-neutral-800">
        <tr>
          <th class="px-3 py-2">Phone</th>
          <th class="px-3 py-2">AI</th>
          <th class="px-3 py-2">Last Message</th>
          <th class="px-3 py-2">Lead</th>
          <th class="px-3 py-2">Session</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(row, idx) in data.slice(0, 5)" :key="idx" class="border-t border-neutral-800">
          <td class="px-3 py-2">{{ pretty(row.phone) }}</td>
          <td class="px-3 py-2">
            <span :class="row.ai_active ? 'text-emerald-400' : 'text-rose-400'">
              {{ row.ai_active ? 'active' : 'off' }}
            </span>
          </td>
          <td class="px-3 py-2">{{ human(row.last_message_at) }}</td>
          <td class="px-3 py-2">{{ row.lead_name || '—' }}</td>
          <td class="px-3 py-2"><code class="text-xs opacity-80">{{ shortId(row.session_id) }}</code></td>
        </tr>
      </tbody>
    </table>

    <!-- Key/Value: show up to 5 metrics -->
    <table
      v-else
      class="w-full text-left border border-neutral-800 rounded-xl overflow-hidden"
      aria-label="Preview table for aggregated report"
    >
      <thead class="bg-neutral-800">
        <tr>
          <th class="px-3 py-2">Metrica</th>
          <th class="px-3 py-2">Valor</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(item, idx) in kvList" :key="idx" class="border-t border-neutral-800">
          <td class="px-3 py-2 capitalize">{{ item[0].replaceAll('_',' ') }}</td>
          <td class="px-3 py-2">{{ formatValue(item[0], item[1]) }}</td>
        </tr>
      </tbody>
    </table>

    <p v-if="isEmpty" class="text-sm text-neutral-400 mt-2">No data available for the selected range.</p>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  reportType: { type: String, required: true },
  data: { type: [Array, Object], default: () => ([]) },
})

const isEmpty = computed(() => {
  if (props.reportType === 'session') return !Array.isArray(props.data) || props.data.length === 0
  return !props.data || (typeof props.data === 'object' && Object.keys(props.data).length === 0)
})

const kvList = computed(() => {
  if (props.reportType === 'session') return []
  if (!props.data || typeof props.data !== 'object') return []

  // Pick up to 5 relevant keys depending on report type
  const pick = (obj, keys) => keys.filter(k => k in obj).map(k => [k, obj[k]])

  if (props.reportType === 'abandonment') {
    return pick(props.data, [
      'total_sessions',
      'completed_sessions',
      'abandonment_rate',
      'total_engaged_sessions',
      'engaged_abandonment_rate',
    ])
  }

  if (props.reportType === 'flowDepth') {
    const out = pick(props.data, ['average_depth'])
    const dist = props.data?.distribution_percent || {}
    // Take up to 4 most representative states
    const topStates = Object.entries(dist)
      .sort((a,b) => b[1] - a[1])
      .slice(0, 4)
      .map(([k,v]) => [`state_${k}_percent`, v])
    return [...out, ...topStates]
  }

  if (props.reportType === 'reengagement') {
    const base = pick(props.data, ['total_recapture_sessions', 'reengaged_sessions', 'reengagement_rate'])
    // Derive remaining as convenience
    const non = (props.data?.total_recapture_sessions ?? 0) - (props.data?.reengaged_sessions ?? 0)
    return [...base, ['non_reengaged_sessions', non]]
  }

  return Object.entries(props.data).slice(0, 5)
})

function formatValue(key, val) {
  if (String(key).includes('rate')) return `${Number(val).toFixed(2)}%`
  if (typeof val === 'number') return String(val)
  if (typeof val === 'string') return val
  try { return JSON.stringify(val) } catch { return String(val) }
}

function shortId(id) {
  if (!id) return '—'
  return String(id).slice(0, 8) + '…'
}

function pretty(raw) {
  // Lightweight pretty: keeps masks like (+xx) and groups
  const s = String(raw || '').trim()
  if (!s) return '—'
  return s
}

function human(iso) {
  if (!iso) return '—'
  const t = Date.parse(iso)
  return Number.isNaN(t) ? '—' : new Date(t).toLocaleString()
}
</script>
