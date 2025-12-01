import { checkAdminAuth } from '@/api/admin'
import { checkAuth } from '@/api/auth'
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
    meta: { requiresAuth: true }
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
    const isAuthenticated = await checkAuth()
    
    if (!isAuthenticated) {
      next('/login')
    } else {
      next()
    }
  } catch (error) {
    console.error('Auth check failed:', error)
    next('/login')
  }
})

export default router
