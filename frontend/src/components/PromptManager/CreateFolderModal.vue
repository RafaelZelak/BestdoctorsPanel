<template>
  <div v-if="isOpen" class="fixed inset-0 bg-black/60 z-50 flex items-center justify-center p-4 backdrop-blur-sm">
    <div class="bg-neutral-800 rounded-xl p-6 max-w-lg w-full border border-neutral-700 shadow-2xl">
      <div class="flex items-center gap-3 mb-5">
        <div class="w-10 h-10 rounded-full bg-blue-500/20 flex items-center justify-center text-blue-400">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10.5v6m3-3H9m4.06-7.19-2.12-2.12a1.5 1.5 0 0 0-1.061-.44H4.5A2.25 2.25 0 0 0 2.25 6v12a2.25 2.25 0 0 0 2.25 2.25h15A2.25 2.25 0 0 0 21.75 18V9a2.25 2.25 0 0 0-2.25-2.25h-5.379a1.5 1.5 0 0 1-1.06-.44Z" />
          </svg>
        </div>
        <h3 class="text-xl font-bold text-white">Nova Pasta</h3>
      </div>
      <div class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-neutral-300 mb-1.5">Nome da pasta</label>
          <input
            ref="nameInput"
            v-model="internalName"
            @keyup.enter="handleCreate"
            type="text"
            placeholder="ex: Prompts de Vendas"
            class="w-full bg-neutral-900 text-neutral-100 px-3 py-2 rounded-lg border border-neutral-700 focus:outline-none focus:border-blue-500 font-mono text-sm transition"
          />
        </div>
      </div>
      <div class="flex justify-end gap-3 mt-6">
        <button @click="$emit('close')" class="px-4 py-2 bg-neutral-700 hover:bg-neutral-600 text-white rounded-lg text-sm font-medium transition">
          Cancelar
        </button>
        <button @click="handleCreate" :disabled="!internalName.trim()" class="px-4 py-2 bg-blue-600 hover:bg-blue-700 disabled:opacity-50 text-white rounded-lg text-sm font-medium transition flex items-center gap-2">
          <span>Criar Pasta</span>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, nextTick } from 'vue'

const props = defineProps({
  isOpen: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['close', 'create'])

const internalName = ref('')
const nameInput = ref(null)

watch(() => props.isOpen, async (newVal) => {
  if (newVal) {
    internalName.value = ''
    await nextTick()
    nameInput.value?.focus()
  }
})

function handleCreate() {
  const trimmed = internalName.value.trim()
  if (trimmed) {
    emit('create', trimmed)
    emit('close')
  }
}
</script>
