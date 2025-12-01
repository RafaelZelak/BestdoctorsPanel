<template>
  <form @submit.prevent aria-labelledby="report-form-title" class="space-y-4">
    <h3 id="report-form-title" class="text-lg font-semibold text-neutral-100">Report options</h3>

    <!-- Step 1: Report Type -->
    <div class="flex flex-col gap-2">
      <label for="reportType" class="text-sm text-neutral-300">Report type</label>
      <select
        id="reportType"
        class="rounded-xl bg-neutral-800 border border-neutral-700 px-3 py-2 outline-none"
        :value="modelValue.report"
        @change="onChangeField('report', $event.target.value)"
        aria-describedby="reportTypeHelp"
      >
        <option value="" disabled>Select a report type…</option>
        <option value="session">session</option>
        <option value="abandonment">abandonment</option>
        <option value="flowDepth">flowDepth</option>
        <option value="reengagement">reengagement</option>
      </select>
      <p id="reportTypeHelp" class="text-xs text-neutral-400">Choose the report you want to generate.</p>
    </div>

    <!-- Step 2: Period (enabled only if report chosen) -->
    <div v-if="hasReport" class="grid grid-cols-1 sm:grid-cols-2 gap-3">
      <div class="flex flex-col gap-2">
        <label for="fromDate" class="text-sm text-neutral-300">From (UTC ISO)</label>
        <input
          id="fromDate"
          type="datetime-local"
          class="rounded-xl bg-neutral-800 border border-neutral-700 px-3 py-2 outline-none"
          :value="toLocalInput(modelValue.filters.from)"
          @change="onChangeField('from', toIsoUTC($event.target.value))"
          placeholder="e.g. 2025-01-01T00:00"
        />
      </div>
      <div class="flex flex-col gap-2">
        <label for="toDate" class="text-sm text-neutral-300">To (UTC ISO)</label>
        <input
          id="toDate"
          type="datetime-local"
          class="rounded-xl bg-neutral-800 border border-neutral-700 px-3 py-2 outline-none"
          :value="toLocalInput(modelValue.filters.to)"
          @change="onChangeField('to', toIsoUTC($event.target.value))"
          placeholder="e.g. 2025-01-31T23:59"
        />
      </div>
    </div>

    <!-- Step 3: "full" only for reengagement -->
    <div v-if="modelValue.report === 'reengagement'" class="flex items-center gap-2">
      <input
        id="fullFlag"
        type="checkbox"
        class="h-4 w-4 rounded border-neutral-600 bg-neutral-800"
        :checked="!!modelValue.filters.full"
        @change="onChangeField('full', $event.target.checked)"
      />
      <label for="fullFlag" class="text-sm text-neutral-300">Full reengagement dataset</label>
    </div>
  </form>
</template>

<script setup>
import { computed, toRaw } from 'vue'

const props = defineProps({
  modelValue: {
    type: Object,
    required: true,
    // shape: { report: string, filters: { from?: string|null, to?: string|null, full?: boolean } }
  }
})
const emit = defineEmits(['update:model-value', 'request-preview'])

const hasReport = computed(() => !!props.modelValue.report)

function onChangeField(key, value) {
  // Build a plain object (avoid cloning Vue proxies)
  const base = props.modelValue || {}
  const next = {
    report: base.report || '',
    filters: { ...(toRaw(base.filters) || {}) },
  }

  if (key === 'report') {
    next.report = value
    if (value !== 'reengagement') next.filters.full = false
  } else if (key === 'from') {
    next.filters.from = value || null
  } else if (key === 'to') {
    next.filters.to = value || null
  } else if (key === 'full') {
    next.filters.full = !!value
  }

  emit('update:model-value', next)
  if (next.report) emit('request-preview')
}

// Helpers to handle datetime-local <-> UTC ISO
function toLocalInput(iso) {
  if (!iso) return ''
  const d = new Date(iso)
  if (Number.isNaN(d.getTime())) return ''
  const pad = (n)=>String(n).padStart(2,'0')
  return `${d.getFullYear()}-${pad(d.getMonth()+1)}-${pad(d.getDate())}T${pad(d.getHours())}:${pad(d.getMinutes())}`
}
function toIsoUTC(localValue) {
  if (!localValue) return null
  const d = new Date(localValue)
  if (Number.isNaN(d.getTime())) return null
  return new Date(Date.UTC(
    d.getFullYear(), d.getMonth(), d.getDate(), d.getHours(), d.getMinutes(), 0, 0
  )).toISOString().replace('.000','')
}
</script>