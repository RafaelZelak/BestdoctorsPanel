<template>
  <div
    class="h-screen grid grid-cols-12 bg-neutral-900 text-neutral-100 overflow-hidden"
  >
    <SessionList
      class="col-span-12 md:col-span-4 lg:col-span-3 h-screen"
      :selected-id="selected?.session_id || ''"
      @select="onSelect"
    />
    <div
      class="col-span-12 md:col-span-8 lg:col-span-9 h-screen flex flex-col overflow-hidden"
    >
      <!-- Header Component -->
      <Header class="flex-shrink-0" />

      <!-- Chat View -->
      <ChatView class="flex-1 overflow-hidden" :session="selected" />
    </div>
  </div>
</template>

<script setup>
import { ref } from "vue";
import SessionList from "@/components/SessionList.vue";
import ChatView from "@/components/ChatView.vue";
import Header from "@/components/Header.vue";

const selected = ref(null);

function onSelect(s) {
  selected.value = s;
}
</script>
