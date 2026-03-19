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
    name: 'Dashboard',
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
        // Redireciona o usuário para algum sistema que ele TENHA acesso
        if (systems.includes('bestdoctors_chat')) {
          next('/')
        } else if (systems.includes('bia_homol')) {
          next('/bia-homol')
        } else if (systems.includes('digesac_homol')) {
          next('/digesac-homol')
        } else {
          // Último caso (sem sistemas), volta pro login
          next('/login')
        }
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
