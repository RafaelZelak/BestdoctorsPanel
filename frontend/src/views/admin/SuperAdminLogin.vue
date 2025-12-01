<!-- SuperAdmin Login Page -->
<template>
  <div class="min-h-screen bg-gradient-to-br from-gray-900 via-slate-900 to-gray-900 flex items-center justify-center p-4">
    <div class="w-full max-w-md">
      <div class="bg-gray-800 rounded-2xl shadow-2xl p-8 border border-gray-700">
        <!-- Header -->
        <div class="text-center mb-8">
          <div class="w-16 h-16 bg-blue-600 rounded-full mx-auto mb-4 flex items-center justify-center">
            <svg class="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
            </svg>
          </div>
          <h1 class="text-2xl font-bold text-white">SuperAdmin Panel</h1>
          <p class="text-gray-400 text-sm mt-1">BestDoctors Administration</p>
        </div>

        <!-- Error Message -->
        <div v-if="error" class="mb-4 p-3 bg-red-900/50 border border-red-700 rounded-lg text-red-200 text-sm">
          {{ error }}
        </div>

        <!-- Login Form -->
        <form @submit.prevent="handleLogin" class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-300 mb-1">Username</label>
            <input
              v-model="username"
              type="text"
              required
              class="w-full px-4 py-2 bg-gray-900 border border-gray-700 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="admin"
              :disabled="loading"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-300 mb-1">Password</label>
            <input
              v-model="password"
              type="password"
              required
              class="w-full px-4 py-2 bg-gray-900 border border-gray-700 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="••••••••"
              :disabled="loading"
            />
          </div>

          <button
            type="submit"
            :disabled="loading"
            class="w-full py-3 bg-blue-600 hover:bg-blue-700 disabled:bg-gray-700 disabled:cursor-not-allowed text-white font-medium rounded-lg transition flex items-center justify-center gap-2"
          >
            <svg v-if="loading" class="animate-spin h-5 w-5" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            <span>{{ loading ? 'Logging in...' : 'Login' }}</span>
          </button>
        </form>

        <!-- Back Link -->
        <div class="mt-6 text-center">
          <router-link to="/login" class="text-sm text-gray-400 hover:text-white transition">
            ← Back to regular login
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { adminLogin } from '@/api/admin'
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const username = ref('')
const password = ref('')
const loading = ref(false)
const error = ref('')

async function handleLogin() {
  error.value = ''
  loading.value = true

  try {
    const result = await adminLogin(username.value, password.value)
    
    if (result.success) {
      router.push('/admin')
    } else {
      error.value = result.message || 'Login failed'
    }
  } catch (err) {
    error.value = err.message || 'Failed to connect to server'
  } finally {
    loading.value = false
  }
}
</script>
