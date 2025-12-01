<!-- src/components/Header.vue -->
<template>
  <header
    class="p-3 bg-neutral-800 border-b border-neutral-700 flex items-center justify-between"
  >
    <!-- User Info -->
    <div class="flex items-center gap-2">
      <p class="text-sm font-medium text-white">
        {{ user?.full_name || user?.username || "Usuário" }}
      </p>
      <span class="text-neutral-600">•</span>
      <p class="text-xs text-neutral-400">{{ user?.role || "user" }}</p>
    </div>

    <!-- Logout Button -->
    <button
      @click="handleLogout"
      class="px-3 py-1.5 bg-red-600 hover:bg-red-700 text-white rounded-lg text-sm font-medium transition flex items-center gap-2"
      aria-label="Logout"
    >
      <svg
        class="w-4 h-4"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"
        />
      </svg>
      Sair
    </button>
  </header>
</template>

<script setup>
import { ref, onMounted } from "vue";
import { useRouter } from "vue-router";
import { logout, getCurrentUser } from "@/api/auth";

const router = useRouter();
const user = ref(null);

onMounted(async () => {
  try {
    user.value = await getCurrentUser();
  } catch (err) {
    console.error("Failed to load user:", err);
  }
});

async function handleLogout() {
  try {
    await logout();
    localStorage.removeItem("user");
    router.push("/login");
  } catch (err) {
    console.error("Logout failed:", err);
    // Force redirect anyway
    router.push("/login");
  }
}
</script>
