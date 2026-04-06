<template>
  <div class="h-screen flex flex-col bg-neutral-900 text-neutral-100 overflow-hidden relative">
    
    <!-- Modals -->
    <InfoModal 
      :is-open="!!modalInfo" 
      :info="modalInfo || {}" 
      @close="modalInfo = null" 
    />
    
    <CreatePromptModal 
      :is-open="showCreateModal" 
      :creating="creating" 
      :error="createError"
      @close="closeCreateModal" 
      @create="handleCreatePrompt" 
    />

    <CreateFolderModal
      :is-open="showCreateFolderModal"
      @close="showCreateFolderModal = false"
      @create="handleCreateFolder"
    />

    <!-- Header -->
    <ManagerHeader 
      :title="title" 
      :syncingAll="syncingAll" 
      :showSync="showSync"
      @sync-all="handleSyncAll" 
      @logout="handleLogout" 
      @toggle-mobile-sidebar="isMobileSidebarOpen = !isMobileSidebarOpen"
    />

    <!-- Main Content -->
    <div class="flex-1 flex overflow-hidden relative">
      <!-- Sidebar (drawer on mobile, static on desktop) -->
      <PromptSidebar 
        :prompts="prompts"
        :loading="loadingList"
        :error="listError"
        :folders="folders"
        :loadingFolders="loadingFolders"
        :folderError="folderError"
        :selectedPromptName="selectedPromptName"
        :allPresence="allPresence"
        :isMobileOpen="isMobileSidebarOpen"
        @select-prompt="selectPrompt"
        @create-new="showCreateModal = true"
        @create-folder="showCreateFolderModal = true"
        @update-folder="handleUpdateFolder"
        @delete-folder="handleDeleteFolder"
        @assign-prompt="handleAssignPromptToFolder"
        @close="isMobileSidebarOpen = false"
      />

      <!-- Viewer -->
      <PromptViewer 
        :selectedPromptName="selectedPromptName"
        :selectedPromptData="selectedPromptData"
        :loading="loadingDetail"
        :error="detailError"
        :isEditing="isEditing"
        v-model="editContent"
        :syncingSingle="syncingSingle"
        :deleting="deleting"
        :saving="saving"
        :showHistoryPanel="showHistoryPanel"
        :selectedVersionDiff="selectedVersionDiff"
        :loadingDiff="loadingDiff"
        :diffError="diffError"
        :reverting="reverting"
        :showSync="showSync"
        @sync-single="handleSyncSingle"
        @start-editing="startEditing"
        @cancel-editing="cancelEditing"
        @delete-prompt="handleDeletePrompt"
        @save-prompt="savePrompt"
        @toggle-history="toggleHistory"
        @revert-to-version="revertToVersion(selectedHistoryVersion)"
      />

      <!-- History Panel (bottom sheet on mobile, right aside on desktop) -->
      <PromptHistoryPanel
        v-if="showHistoryPanel"
        :historyList="historyList"
        :loading="loadingHistory"
        :error="historyError"
        :selectedVersion="selectedHistoryVersion"
        :isOpen="showHistoryPanel"
        @close="toggleHistory"
        @select-version="viewVersionDiff"
        @clear-version="clearHistoryVersion"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import InfoModal from '@/components/PromptManager/InfoModal.vue'
import CreatePromptModal from '@/components/PromptManager/CreatePromptModal.vue'
import CreateFolderModal from '@/components/PromptManager/CreateFolderModal.vue'
import ManagerHeader from '@/components/PromptManager/ManagerHeader.vue'
import PromptSidebar from '@/components/PromptManager/PromptSidebar.vue'
import PromptViewer from '@/components/PromptManager/PromptViewer.vue'
import PromptHistoryPanel from '@/components/PromptManager/PromptHistoryPanel.vue'
import { usePromptManager } from '@/composables/usePromptManager'

const props = defineProps({
  title: {
    type: String,
    required: true
  },
  systemId: {
    type: String,
    required: true
  },
  apiBase: {
    type: String,
    required: true
  }
})

const showCreateFolderModal = ref(false)
const isMobileSidebarOpen = ref(false)

const showSync = computed(() => !props.apiBase.includes('/prod'))

const {
  prompts,
  loadingList,
  listError,
  folders,
  loadingFolders,
  folderError,
  selectedPromptName,
  selectedPromptData,
  loadingDetail,
  detailError,
  isEditing,
  editContent,
  saving,
  syncingAll,
  syncingSingle,
  deleting,
  modalInfo,
  showCreateModal,
  creating,
  createError,
  
  showHistoryPanel,
  historyList,
  loadingHistory,
  historyError,
  selectedHistoryVersion,
  selectedVersionDiff,
  loadingDiff,
  diffError,
  reverting,
  
  allPresence,
  handleLogout,
  selectPrompt,
  startEditing,
  cancelEditing,
  savePrompt,
  handleSyncAll,
  handleSyncSingle,
  closeCreateModal,
  handleCreatePrompt,
  handleDeletePrompt,
  handleCreateFolder,
  handleUpdateFolder,
  handleDeleteFolder,
  handleAssignPromptToFolder,
  toggleHistory,
  clearHistoryVersion,
  viewVersionDiff,
  revertToVersion
} = usePromptManager(props.systemId, props.apiBase)
</script>
