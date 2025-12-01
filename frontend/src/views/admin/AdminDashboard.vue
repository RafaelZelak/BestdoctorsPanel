<!-- Admin Dashboard - User Management -->
<template>
  <div class="min-h-screen bg-gray-900 text-white">
    <!-- Admin Header -->
    <header class="bg-gray-800 border-b border-gray-700 p-4 flex items-center justify-between">
      <div>
        <h1 class="text-xl font-bold">SuperAdmin Panel</h1>
        <p class="text-sm text-gray-400">User Management</p>
      </div>
      <button
        @click="handleLogout"
        class="px-4 py-2 bg-red-600 hover:bg-red-700 rounded-lg text-sm font-medium transition"
      >
        Logout
      </button>
    </header>

    <div class="p-6">
      <!-- Actions Bar -->
      <div class="flex items-center justify-between mb-6">
        <h2 class="text-2xl font-bold">Users</h2>
        <button
          @click="openCreateModal"
          class="px-4 py-2 bg-blue-600 hover:bg-blue-700 rounded-lg font-medium transition flex items-center gap-2"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
          </svg>
          Create User
        </button>
      </div>

      <!-- Users Table -->
      <div class="bg-gray-800 rounded-lg border border-gray-700 overflow-hidden">
        <table class="w-full">
          <thead class="bg-gray-900 border-b border-gray-700">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-400 uppercase">ID</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-400 uppercase">Username</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-400 uppercase">Email</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-400 uppercase">Full Name</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-400 uppercase">Status</th>
              <th class="px-6 py-3 text-right text-xs font-medium text-gray-400 uppercase">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-700">
            <tr v-for="user in users" :key="user.id" class="hover:bg-gray-750 transition">
              <td class="px-6 py-4 text-sm">{{ user.id }}</td>
              <td class="px-6 py-4 text-sm font-medium">{{ user.username }}</td>
              <td class="px-6 py-4 text-sm text-gray-300">{{ user.email }}</td>
              <td class="px-6 py-4 text-sm text-gray-300">{{ user.full_name }}</td>
              <td class="px-6 py-4 text-sm">
                <span :class="user.is_active ? 'bg-green-900 text-green-200' : 'bg-red-900 text-red-200'" class="px-2 py-1 rounded text-xs">
                  {{ user.is_active ? 'Active' : 'Inactive' }}
                </span>
              </td>
              <td class="px-6 py-4 text-sm text-right space-x-2">
                <button @click="openEditModal(user)" class="text-blue-400 hover:text-blue-300">Edit</button>
                <button @click="confirmDelete(user)" class="text-red-400 hover:text-red-300">Delete</button>
              </td>
            </tr>
          </tbody>
        </table>

        <div v-if="loading" class="p-8 text-center text-gray-400">Loading users...</div>
        <div v-if="!loading && users.length === 0" class="p-8 text-center text-gray-400">No users found</div>
      </div>
    </div>

    <!-- Create/Edit User Modal -->
    <div v-if="showModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50" @click.self="closeModal">
      <div class="bg-gray-800 rounded-lg p-6 max-w-md w-full border border-gray-700">
        <h3 class="text-xl font-bold mb-4">{{ isEditing ? 'Edit User' : 'Create User' }}</h3>
        
        <form @submit.prevent="submitForm" class="space-y-4">
          <div v-if="!isEditing">
            <label class="block text-sm font-medium text-gray-300 mb-1">Username</label>
            <input v-model="formData.username" type="text" required class="w-full px-3 py-2 bg-gray-900 border border-gray-700 rounded text-white" />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-300 mb-1">Email</label>
            <input v-model="formData.email" type="email" required class="w-full px-3 py-2 bg-gray-900 border border-gray-700 rounded text-white" />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-300 mb-1">Full Name</label>
            <input v-model="formData.full_name" type="text" required class="w-full px-3 py-2 bg-gray-900 border border-gray-700 rounded text-white" />
          </div>

          <div v-if="!isEditing">
            <label class="block text-sm font-medium text-gray-300 mb-1">Password</label>
            <input v-model="formData.password" type="password" required class="w-full px-3 py-2 bg-gray-900 border border-gray-700 rounded text-white" />
          </div>

          <div class="flex items-center gap-2">
            <input v-model="formData.is_active" type="checkbox" id="is_active" class="w-4 h-4" />
            <label for="is_active" class="text-sm text-gray-300">Active</label>
          </div>

          <div v-if="modalError" class="p-3 bg-red-900/50 border border-red-700 rounded text-red-200 text-sm">
            {{ modalError }}
          </div>

          <div class="flex gap-3 pt-2">
            <button type="button" @click="closeModal" class="flex-1 px-4 py-2 bg-gray-700 hover:bg-gray-600 rounded font-medium transition">
              Cancel
            </button>
            <button type="submit" :disabled="modalLoading" class="flex-1 px-4 py-2 bg-blue-600 hover:bg-blue-700 disabled:bg-gray-700 rounded font-medium transition">
              {{ modalLoading ? 'Saving...' : (isEditing ? 'Update' : 'Create') }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { adminLogout, createUser, deleteUser, getUsers, updateUser } from '@/api/admin'
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const users = ref([])
const loading = ref(false)
const showModal = ref(false)
const isEditing = ref(false)
const modalLoading = ref(false)
const modalError = ref('')

const formData = ref({
  username: '',
  email: '',
  full_name: '',
  password: '',
  role: 'user',
  is_active: true
})

onMounted(async () => {
  await loadUsers()
})

async function loadUsers() {
  loading.value = true
  try {
    const result = await getUsers()
    users.value = result.users || []
  } catch (err) {
    console.error('Failed to load users:', err)
  } finally {
    loading.value = false
  }
}

function openCreateModal() {
  isEditing.value = false
  formData.value = {
    username: '',
    email: '',
    full_name: '',
    password: '',
    is_active: true
  }
  modalError.value = ''
  showModal.value = true
}

function openEditModal(user) {
  isEditing.value = true
  formData.value = {
    id: user.id,
    email: user.email,
    full_name: user.full_name,
    is_active: user.is_active
  }
  modalError.value = ''
  showModal.value = true
}

function closeModal() {
  showModal.value = false
  modalError.value = ''
}

async function submitForm() {
  modalError.value = ''
  modalLoading.value = true

  try {
    if (isEditing.value) {
      await updateUser(formData.value.id, {
        email: formData.value.email,
        full_name: formData.value.full_name,
        is_active: formData.value.is_active
      })
    } else {
      await createUser(formData.value)
    }
    
    closeModal()
    await loadUsers()
  } catch (err) {
    modalError.value = err.message
  } finally {
    modalLoading.value = false
  }
}

async function confirmDelete(user) {
  if (!confirm(`Delete user "${user.username}"?`)) return

  try {
    await deleteUser(user.id)
    await loadUsers()
  } catch (err) {
    alert(`Failed to delete user: ${err.message}`)
  }
}

async function handleLogout() {
  await adminLogout()
  router.push('/admin/login')
}
</script>
