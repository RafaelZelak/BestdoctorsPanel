<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-gray-900 via-slate-900 to-gray-900">
    <div class="max-w-md w-full px-6">
      <!-- Login Card -->
      <div class="bg-gray-800 backdrop-blur-sm rounded-2xl shadow-2xl p-8">
        <!-- Header -->
        <div class="text-center mb-8">
          <div class="w-16 h-16 bg-blue-600 rounded-2xl mx-auto mb-4 flex items-center justify-center">
            <svg class="w-10 h-10 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
            </svg>
          </div>
          <h1 class="text-3xl font-bold text-white mb-2">Painel de Sistemas</h1>
          <p class="text-neutral-400">Autenticação</p>
        </div>

        <!-- Form -->
        <!-- Form Step 1: Login Credentials -->
        <form v-if="step === 1" @submit.prevent="handleLogin" class="space-y-6">
          <!-- Username -->
          <div>
            <label class="block text-sm font-medium text-neutral-300 mb-2">
              Usuário
            </label>
            <input
              v-model="username"
              type="text"
              required
              autocomplete="username"
              class="w-full px-4 py-3 bg-gray-600/50 text-white rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition"
              placeholder="Enter your username"
            />
          </div>

          <!-- Password -->
          <div>
            <label class="block text-sm font-medium text-neutral-300 mb-2">
              Senha
            </label>
            <input
              v-model="password"
              type="password"
              required
              autocomplete="current-password"
              class="w-full px-4 py-3 bg-gray-600/50 text-white rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition"
              placeholder="Enter your password"
            />
          </div>

          <!-- Error Message -->
          <div v-if="error" class="bg-red-500/20 border border-red-500 text-red-200 px-4 py-3 rounded-lg text-sm">
            {{ error }}
          </div>

          <!-- Submit Button -->
          <button
            type="submit"
            :disabled="loading"
            class="w-full bg-blue-600 hover:bg-blue-700 disabled:bg-neutral-600 disabled:cursor-not-allowed text-white font-semibold py-3 rounded-lg transition duration-200 flex items-center justify-center gap-2"
          >
            <svg v-if="loading" class="animate-spin h-5 w-5 text-white" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            <span>{{ loading ? 'Logging in...' : 'Login' }}</span>
          </button>
        </form>

        <!-- Form Step 2: System Selection -->
        <div v-else-if="step === 2" class="space-y-4">
          <h2 class="text-xl font-semibold text-white mb-4 text-center">Selecione o Sistema</h2>
          <div class="grid gap-3">
            <button
              v-for="sys in availableSystems"
              :key="sys"
              @click="handleSystemSelection(sys)"
              class="w-full bg-slate-700 hover:bg-slate-600 text-white font-medium py-3 rounded-lg transition border border-slate-600"
            >
              {{ formatSystemName(sys) }}
            </button>
          </div>
          
          <button @click="step = 1" class="w-full mt-4 text-sm text-neutral-400 hover:text-white transition">
            Voltar
          </button>
        </div>

        <!-- Footer -->
        <div class="mt-6 text-center text-sm text-neutral-500">
          Secure authentication with session management
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { login } from '@/api/auth'
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const username = ref('')
const password = ref('')
const loading = ref(false)
const error = ref('')
const step = ref(1)
const availableSystems = ref([])

function formatSystemName(sys) {
  const names = {
    'bestdoctors_chat': 'Chat BestDoctors',
    'digesac_homol': 'Digesac (Homologação)'
  }
  return names[sys] || sys
}

async function handleLogin() {
  error.value = ''
  loading.value = true

  try {
    const response = await login(username.value, password.value)
    if (response.success) {
      if (response.user) {
        localStorage.setItem('user', JSON.stringify(response.user))
        
        const systems = response.user.system || [];
        
        if (systems.length > 1) {
          availableSystems.value = systems;
          step.value = 2; // Move to select screen
          return;
        } else if (systems.length === 1) {
          handleSystemSelection(systems[0]);
          return;
        } else {
          error.value = 'Usuário sem sistemas vinculados.';
          return;
        }
      }
      router.push('/')
    } else {
      error.value = response.message || 'Login failed'
    }
  } catch (err) {
    error.value = err.message || 'Network error. Please try again.'
  } finally {
    loading.value = false
  }
}

function handleSystemSelection(sys) {
  if (sys === 'bestdoctors_chat') {
    router.push('/')
  } else if (sys === 'digesac_homol') {
    router.push('/digesac-homol')
  } else {
    // Default fallback
    router.push('/')
  }
}
</script>
