<template>
  <button
    :class="[
      'w-full text-left px-4 py-3 transition-colors rounded-lg m-0.5',
      selected
        ? 'bg-neutral-600/60'
        : 'bg-neutral-700/10'
    ]"
    @click="$emit('select', session)"
    :aria-current="selected ? 'true' : 'false'"
    :data-selected="selected ? 'true' : 'false'"
  >
    <div class="flex items-start gap-3 p-1">
      <span
        class="mt-2 inline-block h-2 w-2 rounded-full"
        :class="session.ai_active ? 'bg-emerald-500' : 'bg-neutral-600'"
      />
      <div class="min-w-0 flex-1">
        <div class="flex items-center justify-between gap-3">
          <span class="truncate font-medium">
            {{ session.lead_name ? session.lead_name : prettyPhone(session.phone) }}
          </span>
        </div>
      </div>
    </div>
  </button>
</template>

<script setup>
import { prettyPhone } from '@/utils/formatters'

defineProps({
  session: { type: Object, required: true },
  selected: { type: Boolean, default: false },
})

defineEmits(['select'])
</script>
