<template>
  <div v-if="isOpen" class="fixed inset-0 bg-black/60 z-50 flex items-center justify-center p-4 backdrop-blur-sm">
    <div class="bg-neutral-800 rounded-xl p-6 max-w-lg w-full border border-neutral-700 shadow-2xl">
      <div class="flex items-center gap-3 mb-5">
        <div class="w-10 h-10 rounded-full bg-green-500/20 flex items-center justify-center text-green-400">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
        </div>
        <h3 class="text-xl font-bold text-white">Novo Prompt</h3>
      </div>
      <div class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-neutral-300 mb-1.5">Nome do arquivo</label>
          <input
            v-model="internalName"
            type="text"
            placeholder="ex: meu_prompt.md"
            class="w-full bg-neutral-900 text-neutral-100 px-3 py-2 rounded-lg border border-neutral-700 focus:outline-none focus:border-green-500 font-mono text-sm transition"
          />
        </div>
        <div>
          <label class="block text-sm font-medium text-neutral-300 mb-1.5">Conteúdo</label>
          <textarea
            v-model="internalContent"
            rows="8"
            placeholder="Digite o conteúdo markdown do prompt..."
            class="w-full bg-neutral-900 text-neutral-100 px-3 py-2 rounded-lg border border-neutral-700 focus:outline-none focus:border-green-500 font-mono text-sm resize-none transition"
          ></textarea>
        </div>
        <p v-if="error" class="text-sm text-red-400">{{ error }}</p>
      </div>
      <div class="flex justify-end gap-3 mt-6">
        <button @click="$emit('close')" class="px-4 py-2 bg-neutral-700 hover:bg-neutral-600 text-white rounded-lg text-sm font-medium transition">
          Cancelar
        </button>
        <button @click="handleCreate" :disabled="creating" class="px-4 py-2 bg-green-600 hover:bg-green-700 disabled:opacity-50 text-white rounded-lg text-sm font-medium transition flex items-center gap-2">
          <svg v-if="creating" class="animate-spin h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          <span>{{ creating ? 'Criando...' : 'Criar Prompt' }}</span>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'

const props = defineProps({
  isOpen: {
    type: Boolean,
    default: false
  },
  creating: {
    type: Boolean,
    default: false
  },
  error: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['close', 'create'])

const internalName = ref('')
const internalContent = ref('')

watch(() => props.isOpen, (newVal) => {
  if (newVal) {
    internalName.value = ''
    internalContent.value = ''
  }
})

function handleCreate() {
  emit('create', {
    name: internalName.value,
    content: internalContent.value
  })
}
</script>
