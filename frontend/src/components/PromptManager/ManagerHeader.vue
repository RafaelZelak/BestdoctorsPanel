<template>
  <header class="flex-shrink-0 p-4 bg-neutral-800 border-b border-neutral-700 flex items-center justify-between">
    <div class="flex items-center gap-3">
      <!-- Hamburger: visible only on mobile -->
      <button
        @click="$emit('toggle-mobile-sidebar')"
        class="md:hidden p-2 rounded-lg text-neutral-400 hover:text-white hover:bg-neutral-700 transition active:scale-95"
        title="Arquivos"
        aria-label="Abrir lista de arquivos"
      >
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
        </svg>
      </button>

      <div>
        <h1 class="text-xl font-bold tracking-wider text-blue-500">{{ title }}</h1>
        <p class="text-xs text-neutral-400">Gerenciador de Prompts</p>
      </div>
    </div>

    <div class="flex items-center gap-2">
      <button v-if="showSync" @click="$emit('sync-all')" :disabled="syncingAll" 
        class="p-2 transition-all duration-300 rounded-lg flex items-center gap-0 hover:gap-2 text-neutral-400 hover:text-yellow-500 hover:bg-yellow-500/10 group h-10"
        title="Sincronizar Tudo">
        <svg v-if="syncingAll" class="animate-spin h-5 w-5 text-yellow-500" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
        <svg v-else class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
        </svg>
        <span class="max-w-0 overflow-hidden whitespace-nowrap group-hover:max-w-xs transition-all duration-300 text-sm font-medium">
          {{ syncingAll ? 'Sincronizando...' : 'Sincronizar Tudo' }}
        </span>
      </button>

      <button @click="$emit('logout')" 
        class="p-2 transition-all duration-300 rounded-lg flex items-center gap-0 hover:gap-2 text-neutral-400 hover:text-red-500 hover:bg-red-500/10 group h-10"
        title="Sair do Sistema">
        <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
        </svg>
        <span class="max-w-0 overflow-hidden whitespace-nowrap group-hover:max-w-xs transition-all duration-300 text-sm font-medium">
          Sair do Sistema
        </span>
      </button>
    </div>
  </header>
</template>

<script setup>
defineProps({
  title: {
    type: String,
    required: true
  },
  syncingAll: {
    type: Boolean,
    default: false
  },
  showSync: {
    type: Boolean,
    default: true
  }
})

defineEmits(['sync-all', 'logout', 'toggle-mobile-sidebar'])
</script>
