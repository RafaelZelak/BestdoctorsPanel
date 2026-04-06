<template>
  <div class="system-selection-container">
    <div class="selection-box">
      <!-- HEADER -->
      <div class="header">
        <div class="header-icon">
          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" class="icon-lock">
            <path stroke-linecap="round" stroke-linejoin="round" d="M16.5 10.5V6.75a4.5 4.5 0 10-9 0v3.75m-.75 11.25h10.5a2.25 2.25 0 002.25-2.25v-6.75a2.25 2.25 0 00-2.25-2.25H6.75a2.25 2.25 0 00-2.25 2.25v6.75a2.25 2.25 0 002.25 2.25z" />
          </svg>
        </div>
        <div class="header-text">
          <h1>Painel de Sistemas</h1>
          <span class="subtitle">AUTENTICAÇÃO</span>
        </div>
      </div>

      <!-- STATES -->
      <div v-if="loading" class="loading">Carregando sistemas...</div>
      <div v-else-if="systems.length === 0" class="no-systems">
        Você não possui acesso a nenhum sistema associado à sua conta.
      </div>
      
      <!-- SYSTEM LIST -->
      <div v-else class="system-section">
        <div class="section-title">SELECIONE O SISTEMA</div>
        <div class="system-list">
          <router-link
            v-for="sys in systems"
            :key="sys.id"
            :to="sys.path"
            class="system-item"
            :class="{ 'has-env': sys.env }"
          >
            <div class="system-item-icon-box">
              <span class="system-icon">{{ sys.icon }}</span>
            </div>
            <div class="system-item-name">{{ sys.name }}</div>
            <div v-if="sys.env" class="system-env" :class="sys.env.toLowerCase()">
              {{ sys.env }}
            </div>
          </router-link>
        </div>
      </div>

      <hr class="divider"/>

      <!-- LOGOUT -->
      <div class="footer-actions">
        <button @click="handleLogout" class="logout-btn">
          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" class="icon-logout">
            <path stroke-linecap="round" stroke-linejoin="round" d="M15.75 9V5.25A2.25 2.25 0 0013.5 3h-6a2.25 2.25 0 00-2.25 2.25v13.5A2.25 2.25 0 007.5 21h6a2.25 2.25 0 002.25-2.25V15M12 9l-3 3m0 0l3 3m-3-3h12.75" />
          </svg>
          Sair
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getCurrentUser, logout } from '@/api/auth'

const router = useRouter()
const loading = ref(true)
const systems = ref([])

const SYSTEM_MAP = [
  { id: 'bestdoctors_chat', path: '/bestdoctors-chat', name: 'Chat BestDoctors', icon: '💬', env: null },
  { id: 'iza_chat', path: '/iza-chat', name: 'IZA Chat', icon: '💬', env: null },
  { id: 'digesac_homol', path: '/digesac-homol', name: 'DIGESAC', icon: '⚙️', env: 'HOMOL' },
  { id: 'digesac_prod', path: '/digesac-prod', name: 'DIGESAC', icon: '⚙️', env: 'PROD' },
  { id: 'bia_homol', path: '/bia-homol', name: 'BIA Prompt', icon: '🤖', env: 'HOMOL' },
  { id: 'bia_prod', path: '/bia-prod', name: 'BIA Prompt', icon: '🤖', env: 'PROD' },
  { id: 'iza_homol', path: '/iza-homol', name: 'IZA Prompt', icon: '✨', env: 'HOMOL' },
  { id: 'iza_prod', path: '/iza-prod', name: 'IZA Prompt', icon: '✨', env: 'PROD' },
  { id: 'iza_extractor_homol', path: '/iza-extractor-homol', name: 'IZA Extractor', icon: '📄', env: 'HOMOL' },
  { id: 'iza_extractor_prod', path: '/iza-extractor-prod', name: 'IZA Extractor', icon: '📄', env: 'PROD' },
  { id: 'iza_classifier_homol', path: '/iza-classifier-homol', name: 'IZA Classifier', icon: '🏷️', env: 'HOMOL' },
  { id: 'iza_classifier_prod', path: '/iza-classifier-prod', name: 'IZA Classifier', icon: '🏷️', env: 'PROD' },
]

onMounted(async () => {
  try {
    const user = await getCurrentUser()
    const userSystems = user?.system || []

    if (userSystems.length === 1) {
      const singleSys = SYSTEM_MAP.find(s => s.id === userSystems[0])
      if (singleSys) {
        router.replace(singleSys.path)
        return
      }
    }

    // Preserve the order defined in SYSTEM_MAP
    systems.value = SYSTEM_MAP.filter(sys => userSystems.includes(sys.id))
  } catch (err) {
    console.error(err)
    router.replace('/login')
  } finally {
    loading.value = false
  }
})

const handleLogout = async () => {
  await logout()
  router.push('/login')
}
</script>

<style scoped>
.system-selection-container {
  display: flex;
  justify-content: center;
  align-items: flex-start;
  min-height: 100vh;
  background-color: #1a1a21; /* Very dark background */
  font-family: 'Inter', system-ui, sans-serif;
  padding: 40px 20px;
  box-sizing: border-box;
}

.selection-box {
  background-color: #1f2026; /* Darker card background */
  padding: 2.5rem;
  border-radius: 12px;
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.5), 0 10px 10px -5px rgba(0, 0, 0, 0.3);
  width: 100%;
  max-width: 580px;
}

/* HEADER */
.header {
  display: flex;
  align-items: center;
  gap: 1rem;
  margin-bottom: 2.5rem;
}

.header-icon {
  background-color: #3b82f6;
  color: white;
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.icon-lock {
  width: 24px;
  height: 24px;
}

.header-text {
  display: flex;
  flex-direction: column;
}

h1 {
  font-size: 1.5rem;
  color: #f3f4f6;
  margin: 0;
  font-weight: 700;
  line-height: 1.2;
}

.subtitle {
  font-size: 0.75rem;
  color: #6b7280;
  letter-spacing: 0.1em;
  margin-top: 0.25rem;
  font-weight: 600;
}

/* SECTION TITLE */
.section-title {
  font-size: 0.75rem;
  color: #6b7280;
  font-weight: 700;
  letter-spacing: 0.1em;
  margin-bottom: 1rem;
  text-transform: uppercase;
}

/* SYSTEM LIST */
.system-list {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
  max-height: 45vh;
  overflow-y: auto;
  padding-right: 8px;
}

.system-list::-webkit-scrollbar {
  width: 6px;
}

.system-list::-webkit-scrollbar-track {
  background: transparent;
}

.system-list::-webkit-scrollbar-thumb {
  background-color: #3f3f46;
  border-radius: 4px;
}

.system-list::-webkit-scrollbar-thumb:hover {
  background-color: #52525b;
}

.system-item {
  display: flex;
  align-items: center;
  padding: 0.75rem 1rem;
  border-radius: 8px;
  background-color: transparent;
  text-decoration: none;
  transition: all 0.2s ease;
  border: 1px solid transparent;
}

.system-item:hover, .system-item:focus {
  background-color: #273045;
  border-color: #3b82f6;
}

.system-item-icon-box {
  background-color: #2a2a35;
  width: 32px;
  height: 32px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 1rem;
  flex-shrink: 0;
}

.system-item:hover .system-item-icon-box {
  background-color: transparent;
}

.system-icon {
  font-size: 1.1rem;
}

.system-item-name {
  color: #d1d5db;
  font-weight: 500;
  font-size: 0.95rem;
  flex-grow: 1;
}

.system-item:hover .system-item-name {
  color: #f3f4f6;
}

.system-env {
  font-size: 0.65rem;
  font-weight: 800;
  padding: 0.2rem 0.5rem;
  border-radius: 4px;
  letter-spacing: 0.05em;
  margin-left: 1rem;
}

.system-env.homol {
  background-color: rgba(245, 158, 11, 0.15);
  color: #fbbf24;
  border: 1px solid rgba(245, 158, 11, 0.3);
}

.system-env.prod {
  background-color: rgba(16, 185, 129, 0.15);
  color: #34d399;
  border: 1px solid rgba(16, 185, 129, 0.3);
}

hr.divider {
  border: 0;
  height: 1px;
  background-color: #2d2d3a;
  margin: 2rem 0 1.5rem 0;
}

/* LOGOUT BUTTON */
.footer-actions {
  display: flex;
  justify-content: flex-start;
}

.logout-btn {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  background-color: transparent;
  color: #d1d5db;
  border: 1px solid #3f3f46;
  padding: 0.5rem 1rem;
  border-radius: 6px;
  font-size: 0.85rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.logout-btn:hover {
  background-color: rgba(255, 255, 255, 0.05);
  color: #f3f4f6;
  border-color: #52525b;
}

.icon-logout {
  width: 16px;
  height: 16px;
}

.loading, .no-systems {
  color: #9ca3af;
  margin-bottom: 1.5rem;
  font-size: 0.95rem;
}
</style>
