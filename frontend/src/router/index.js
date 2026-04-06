import { checkAdminAuth } from '@/api/admin'
import { getCurrentUser } from '@/api/auth'
import Dashboard from '@/components/Dashboard.vue'
import Login from '@/views/Login.vue'
import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: Login,
    meta: { public: true }
  },
  {
    path: '/',
    name: 'Root',
    component: () => import('@/views/SystemSelection.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/bestdoctors-chat',
    name: 'BestdoctorsChat',
    component: Dashboard,
    meta: { requiresAuth: true, requiresSystem: 'bestdoctors_chat' }
  },
  {
    path: '/digesac-homol',
    name: 'DigesacHomol',
    component: () => import('@/views/PromptManager.vue'),
    props: { title: 'DIGESAC HOMOL', systemId: 'digesac_homol', apiBase: '/api/proxy/digesac/homol' },
    meta: { requiresAuth: true, requiresSystem: 'digesac_homol' }
  },
  {
    path: '/bia-homol',
    name: 'BiaHomol',
    component: () => import('@/views/PromptManager.vue'),
    props: { title: 'BIA HOMOL', systemId: 'bia_homol', apiBase: '/api/proxy/bia/homol' },
    meta: { requiresAuth: true, requiresSystem: 'bia_homol' }
  },
  {
    path: '/bia-prod',
    name: 'BiaProd',
    component: () => import('@/views/PromptManager.vue'),
    props: { title: 'BIA PROD', systemId: 'bia_prod', apiBase: '/api/proxy/bia/prod' },
    meta: { requiresAuth: true, requiresSystem: 'bia_prod' }
  },
  {
    path: '/digesac-prod',
    name: 'DigesacProd',
    component: () => import('@/views/PromptManager.vue'),
    props: { title: 'DIGESAC PROD', systemId: 'digesac_prod', apiBase: '/api/proxy/digesac/prod' },
    meta: { requiresAuth: true, requiresSystem: 'digesac_prod' }
  },
  {
    path: '/iza-prod',
    name: 'IzaProd',
    component: () => import('@/views/PromptManager.vue'),
    props: { title: 'IZA PROD', systemId: 'iza_prod', apiBase: '/api/proxy/iza/prod' },
    meta: { requiresAuth: true, requiresSystem: 'iza_prod' }
  },
  {
    path: '/iza-homol',
    name: 'IzaHomol',
    component: () => import('@/views/PromptManager.vue'),
    props: { title: 'IZA HOMOL', systemId: 'iza_homol', apiBase: '/api/proxy/iza/homol' },
    meta: { requiresAuth: true, requiresSystem: 'iza_homol' }
  },
  {
    path: '/iza-chat',
    name: 'IzaChat',
    component: () => import('@/views/iza_chat/IzaChat.vue'),
    meta: { requiresAuth: true, requiresSystem: 'iza_chat' }
  },
  {
    path: '/iza-extractor-homol',
    name: 'IzaExtractorHomol',
    component: () => import('@/views/PromptManager.vue'),
    props: { title: 'IZA EXTRACTOR HOMOL', systemId: 'iza_extractor_homol', apiBase: '/api/proxy/iza-extractor/homol' },
    meta: { requiresAuth: true, requiresSystem: 'iza_extractor_homol' }
  },
  {
    path: '/iza-extractor-prod',
    name: 'IzaExtractorProd',
    component: () => import('@/views/PromptManager.vue'),
    props: { title: 'IZA EXTRACTOR PROD', systemId: 'iza_extractor_prod', apiBase: '/api/proxy/iza-extractor/prod' },
    meta: { requiresAuth: true, requiresSystem: 'iza_extractor_prod' }
  },
  {
    path: '/iza-classifier-homol',
    name: 'IzaClassifierHomol',
    component: () => import('@/views/PromptManager.vue'),
    props: { title: 'IZA CLASSIFIER HOMOL', systemId: 'iza_classifier_homol', apiBase: '/api/proxy/iza-classifier/homol' },
    meta: { requiresAuth: true, requiresSystem: 'iza_classifier_homol' }
  },
  {
    path: '/iza-classifier-prod',
    name: 'IzaClassifierProd',
    component: () => import('@/views/PromptManager.vue'),
    props: { title: 'IZA CLASSIFIER PROD', systemId: 'iza_classifier_prod', apiBase: '/api/proxy/iza-classifier/prod' },
    meta: { requiresAuth: true, requiresSystem: 'iza_classifier_prod' }
  },
  // Admin routes
  {
    path: '/admin/login',
    name: 'AdminLogin',
    component: () => import('@/views/admin/SuperAdminLogin.vue'),
    meta: { public: true }
  },
  {
    path: '/admin',
    name: 'AdminDashboard',
    component: () => import('@/views/admin/AdminDashboard.vue'),
    meta: { requiresSuperAdmin: true }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

function getDefaultRouteForUser(user) {
  const systems = user?.system || []
  if (systems.includes('bestdoctors_chat')) return '/bestdoctors-chat'
  if (systems.includes('bia_homol')) return '/bia-homol'
  if (systems.includes('bia_prod')) return '/bia-prod'
  if (systems.includes('digesac_homol')) return '/digesac-homol'
  if (systems.includes('digesac_prod')) return '/digesac-prod'
  if (systems.includes('iza_homol')) return '/iza-homol'
  if (systems.includes('iza_prod')) return '/iza-prod'
  if (systems.includes('iza_chat')) return '/iza-chat'
  if (systems.includes('iza_extractor_homol')) return '/iza-extractor-homol'
  if (systems.includes('iza_extractor_prod')) return '/iza-extractor-prod'
  if (systems.includes('iza_classifier_homol')) return '/iza-classifier-homol'
  if (systems.includes('iza_classifier_prod')) return '/iza-classifier-prod'
  return '/login'
}

// Authentication guard
router.beforeEach(async (to, from, next) => {
  const isPublic = to.meta.public

  if (isPublic) {
    next()
    return
  }

  // For admin routes, check admin auth
  if (to.meta.requiresSuperAdmin) {
    const isAdminAuthenticated = await checkAdminAuth()
    if (!isAdminAuthenticated) {
      next('/admin/login')
    } else {
      next()
    }
    return
  }

  // For regular routes, check regular auth
  try {
    const user = await getCurrentUser()

    if (to.meta.requiresSystem) {
      const systems = user?.system || []
      const hasAccess = systems.includes(to.meta.requiresSystem)

      if (!hasAccess) {
        next(getDefaultRouteForUser(user))
        return
      }
    }

    next()
  } catch (error) {
    console.error('Auth check failed:', error)
    next('/login')
  }
})

export default router
