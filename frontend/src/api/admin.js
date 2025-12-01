// Admin API client
const API_BASE = '/admin'

export async function adminLogin(username, password) {
  const response = await fetch(`${API_BASE}/auth`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
    body: JSON.stringify({ username, password })
  })
  
  if (!response.ok) {
    const data = await response.json()
    throw new Error(data.message || 'Login failed')
  }
  
  return response.json()
}

export async function adminLogout() {
  const response = await fetch(`${API_BASE}/logout`, {
    method: 'POST',
    credentials: 'include'
  })
  return response.json()
}

export async function checkAdminAuth() {
  try {
    // Try to fetch users (lightweight check) to validate session
    const response = await fetch(`${API_BASE}/users`, {
      method: 'GET',
      credentials: 'include'
    })
    return response.ok
  } catch (e) {
    return false
  }
}

export async function createUser(userData) {
  const response = await fetch(`${API_BASE}/users`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
    body: JSON.stringify(userData)
  })
  
  if (!response.ok) {
    const data = await response.json()
    throw new Error(data.message || 'Failed to create user')
  }
  
  return response.json()
}

export async function getUsers() {
  const response = await fetch(`${API_BASE}/users`, {
    credentials: 'include'
  })
  
  if (!response.ok) {
    throw new Error('Failed to fetch users')
  }
  
  return response.json()
}

export async function getUser(id) {
  const response = await fetch(`${API_BASE}/users/${id}`, {
    credentials: 'include'
  })
  
  if (!response.ok) {
    throw new Error('Failed to fetch user')
  }
  
  return response.json()
}

export async function updateUser(id, userData) {
  const response = await fetch(`${API_BASE}/users/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
    body: JSON.stringify(userData)
  })
  
  if (!response.ok) {
    const data = await response.json()
    throw new Error(data.message || 'Failed to update user')
  }
  
  return response.json()
}

export async function deleteUser(id) {
  const response = await fetch(`${API_BASE}/users/${id}`, {
    method: 'DELETE',
    credentials: 'include'
  })
  
  if (!response.ok) {
    throw new Error('Failed to delete user')
  }
  
  return response.json()
}

export async function resetPassword(id, newPassword) {
  const response = await fetch(`${API_BASE}/users/${id}/password`, {
    method: 'PATCH',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
    body: JSON.stringify({ new_password: newPassword })
  })
  
  if (!response.ok) {
    const data = await response.json()
    throw new Error(data.message || 'Failed to reset password')
  }
  
  return response.json()
}
