<template>
  <div
    class="h-screen flex flex-col md:grid md:grid-cols-12 bg-neutral-900 text-neutral-100 overflow-hidden"
  >
    <!-- Session List: visible on desktop always, on mobile only when list view is active -->
    <SessionList
      v-show="!isMobile || isListVisible"
      class="md:col-span-4 lg:col-span-3 h-screen"
      :selected-id="selected?.session_id || ''"
      @select="onSelect"
    />
    
    <!-- Chat Area: visible on desktop always, on mobile only when chat view is active -->
    <div
      v-show="!isMobile || isChatVisible"
      class="md:col-span-8 lg:col-span-9 h-screen flex flex-col overflow-hidden"
    >
      <!-- Header Component -->
      <Header class="flex-shrink-0" />

      <!-- Chat View with mobile back button support -->
      <ChatView 
        class="flex-1 overflow-hidden" 
        :session="selected"
        :is-mobile="isMobile"
        @back="handleBack"
      />
    </div>
  </div>
</template>

<script setup>
import ChatView from "@/components/ChatView.vue";
import Header from "@/components/Header.vue";
import SessionList from "@/components/SessionList.vue";
import { useMobileNavigation } from "@/composables/useMobileNavigation";
import { useResponsive } from "@/composables/useResponsive";
import { ref } from "vue";

const selected = ref(null);

// Responsive detection
const { isMobile } = useResponsive();

// Mobile navigation management
const { isListVisible, isChatVisible, showList, showChat } = useMobileNavigation();

function onSelect(s) {
  selected.value = s;
  
  // On mobile, switch to chat view when session is selected
  if (isMobile.value) {
    showChat();
  }
}

function handleBack() {
  // On mobile, return to list view
  if (isMobile.value) {
    showList();
  }
}
</script>

