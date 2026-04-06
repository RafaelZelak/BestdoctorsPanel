<template>
  <div v-if="isOpen" class="fixed inset-0 bg-black/60 z-[60] flex items-center justify-center p-4 backdrop-blur-sm">
    <div class="bg-neutral-800 rounded-xl p-6 max-w-md w-full border border-neutral-700 shadow-2xl transform transition-all">
      <div class="flex items-center gap-3 mb-4">
        <div class="w-10 h-10 rounded-full bg-yellow-500/20 flex items-center justify-center text-yellow-500">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        </div>
        <h3 class="text-xl font-bold text-white">{{ info?.title }}</h3>
      </div>
      <p class="text-neutral-300 mb-6 leading-relaxed">{{ info?.message }}</p>
      <div class="flex justify-end">
        <button @click="$emit('close')" class="px-5 py-2.5 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition font-medium">
          Entendido
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  isOpen: {
    type: Boolean,
    default: false
  },
  info: {
    type: Object,
    default: () => ({})
  }
})

defineEmits(['close'])
</script>
