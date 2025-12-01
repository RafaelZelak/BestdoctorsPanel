<!-- src/components/MessageBubble.vue -->
<template>
  <div class="mb-4">
    <!-- Aviso do sistema (não é balão de chat) -->
    <div v-if="isNotice" class="flex justify-end">
      <div
        class="max-w-[85%] w-full rounded-xl border border-amber-500/30 bg-amber-400/10 text-amber-200 px-3 py-2 text-sm shadow-sm"
        role="note"
        aria-label="Aviso automático"
      >
        <div class="flex items-center gap-2 mb-1">
          <span class="text-[10px] uppercase tracking-wide px-2 py-0.5 rounded-full bg-amber-500/20 border border-amber-500/30">
            sistema
          </span>
          <span class="opacity-80 text-[12px]">Regra de recapture</span>
          <span class="ml-auto text-[10px] opacity-70">{{ timestamp }}</span>
        </div>
        <div class="opacity-90">
          {{ cleanNotice }}
        </div>
      </div>
    </div>

    <!-- Balão padrão -->
    <div v-else class="flex" :class="role === 'human' ? 'justify-end' : 'justify-start'">
      <div
        class="max-w-[80%] rounded-2xl px-4 py-3 shadow"
        :class="role === 'human'
          ? 'bg-blue-600 text-white'
          : 'bg-neutral-800 text-neutral-100 border border-neutral-700'"
      >
        <div class="whitespace-pre-wrap break-words">
          {{ text }}
        </div>
        <div class="text-[10px] opacity-70 mt-1 text-right">
          {{ timestamp }}
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  role: { type: String, required: true },      // 'human' | 'ai'
  text: { type: String, required: true },
  timestamp: { type: String, required: true }, // already formatted
})

// Detecta mensagens automáticas do seu sistema (padrão “Recapture - ...”)
const isNotice = computed(() => {
  const t = (props.text || '').trim()
  return /^\s*recapture\s*-\s*/i.test(t)
})

// Mostra o conteúdo sem o prefixo “Recapture - ”
const cleanNotice = computed(() => {
  const t = (props.text || '').trim()
  return t.replace(/^\s*recapture\s*-\s*/i, '')
})
</script>
