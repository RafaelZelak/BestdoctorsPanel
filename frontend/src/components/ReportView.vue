<template>
  <Teleport to="body">
    <!-- Overlay -->
    <Transition
      appear
      enter-active-class="duration-200 ease-out"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
      leave-active-class="duration-150 ease-in"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div
        v-if="isVisible"
        class="fixed inset-0 bg-neutral-900/40 backdrop-blur-[2px] z-40"
        @click="closeModal"
        aria-hidden="true"
      />
    </Transition>

    <!-- Wrapper -->
    <div v-if="isVisible" class="fixed inset-0 z-50 flex items-center justify-center p-4 text-gray-200"
         role="dialog" aria-modal="true" @click.self="closeModal">

      <!-- Panel -->
      <Transition
        appear
        enter-active-class="bounce-in"
        leave-active-class="shrink-fade-out"
      >
        <div v-if="isVisible"
             class="rounded-2xl bg-neutral-900 shadow-xl border border-neutral-700 w-full max-w-7xl h-[90vh] flex flex-col">
          <header class="flex justify-between items-center p-6 border-b border-neutral-700">
            <h2 class="text-2xl font-semibold text-neutral-100">Gerador de Relatórios</h2>
            <button @click="closeModal" class="text-neutral-400 hover:text-neutral-200 text-3xl leading-none" aria-label="Close">&times;</button>
          </header>

          <!-- Content -->
         <div 
            v-if="isMobile" 
            class="flex items-center justify-center flex-grow text-center p-10"
          >
            <p class="text-lg text-neutral-300">
              Acesse por um Computador para gerar relatórios
            </p>
          </div>

          <!-- Desktop Content -->
          <div 
            v-else 
            class="p-6 flex-grow overflow-auto space-y-6 overflow-hidden"
          >
            <!-- Step 1: Report Type -->
            <section>
              <label for="reportType" class="block text-sm mb-2">Tipo de Relatório</label>
              <select
                id="reportType"
                v-model="form.report"
                class="w-full md:w-80 bg-neutral-800 border border-neutral-700 rounded-xl px-3 py-2 outline-none"
                aria-describedby="report-help"
              >
                <option disabled value="">Seleione…</option>
                <option value="session">Conversas Iniciadas</option>
                <option value="abandonment">Taxa de Abandono</option>
                <option value="flowDepth">Etapa Máxima do Fluxo</option>
                <option value="reengagement">Sessões Reengajadas</option>
              </select>
              <p id="report-help" class="text-xs text-neutral-400 mt-1">Selecione o "Tipo de Relatório" para Gerar</p>
            </section>

            <!-- Step 2: Period (unlocked after report type) -->
            <section :aria-disabled="!form.report" :class="!form.report ? 'opacity-50 pointer-events-none' : ''">
              <div class="grid grid-cols-1 md:grid-cols-3 gap-4 items-end">
                <div>
                  <label for="from" class="block text-sm mb-2">De: (Opicional)</label>
                  <input
                    id="from"
                    type="date"
                    v-model="dates.from"
                    class="w-full bg-neutral-800 border border-neutral-700 rounded-xl px-3 py-2 outline-none"
                    :disabled="!form.report"
                  />
                </div>
                <div>
                  <label for="to" class="block text-sm mb-2">Até (Opicional)</label>
                  <input
                    id="to"
                    type="date"
                    v-model="dates.to"
                    class="w-full bg-neutral-800 border border-neutral-700 rounded-xl px-3 py-2 outline-none"
                    :disabled="!form.report"
                  />
                </div>

                <div v-if="form.report === 'reengagement'">
                  <label class="block text-sm mb-2">Conteúdo</label>
                  <div class="flex items-center gap-3">
                    <input id="full" type="checkbox" v-model="form.full" class="h-4 w-4" />
                    <label for="full" class="text-sm">Todos os Detalhes</label>
                  </div>
                </div>
              </div>
              <p class="text-xs text-neutral-400 mt-1">
                Deixe as datas em branco para buscar tudo.
              </p>
            </section>

            <!-- Step 3: Auto Preview (JSON) -->
            <section :aria-disabled="!form.report" :class="!form.report ? 'opacity-50 pointer-events-none' : ''">
              <div class="flex items-center justify-between mb-2">
                <h3 class="text-lg font-medium">Preview</h3>
                <div class="text-xs text-neutral-400" role="status" aria-live="polite">
                  <span v-if="loading">Loading preview…</span>
                  <span v-else-if="error" class="text-rose-400">{{ error }}</span>
                  <span v-else-if="previewReady" class="text-emerald-400">Pré Vizualização</span>
                </div>
              </div>

              <div v-if="loading" class="border border-neutral-800 rounded-xl p-4 animate-pulse">
                <div class="h-4 bg-neutral-800 rounded w-2/3 mb-3"></div>
                <div class="h-4 bg-neutral-800 rounded w-1/2 mb-3"></div>
                <div class="h-4 bg-neutral-800 rounded w-3/4"></div>
              </div>

              <ReportPreviewTable
                v-else-if="previewReady"
                :report-type="form.report"
                :data="preview"
              />

              <p v-else class="text-sm text-neutral-400">Select a report type to see a preview.</p>
            </section>

            <!-- Step 4: Export (unlocked after preview) -->
            <section :aria-disabled="!previewReady" :class="!previewReady ? 'opacity-50 pointer-events-none' : ''">
              <h3 class="text-lg font-medium mb-2">Exportar</h3>
              <div class="flex flex-wrap gap-3 items-center">
                
                <!-- CSV -->
                <label class="flex items-center gap-2 cursor-pointer">
                  <input 
                    type="radio" 
                    name="export-type" 
                    value="csv" 
                    v-model="exportType" 
                    class="peer hidden"
                  />
                  <span
                    class="px-3 py-1 rounded-full border border-gray-500 text-sm 
                          peer-checked:bg-purple-600 peer-checked:text-white 
                          transition-colors"
                  >
                    CSV
                  </span>
                </label>

                <!-- PDF -->
                <label class="flex items-center gap-2 cursor-pointer">
                  <input 
                    type="radio" 
                    name="export-type" 
                    value="pdf" 
                    v-model="exportType" 
                    class="peer hidden"
                  />
                  <span
                    class="px-3 py-1 rounded-full border border-gray-500 text-sm 
                          peer-checked:bg-purple-600 peer-checked:text-white 
                          transition-colors"
                  >
                    PDF
                  </span>
                </label>

                <!-- XLSX -->
                <label class="flex items-center gap-2 cursor-pointer">
                  <input 
                    type="radio" 
                    name="export-type" 
                    value="xlsx" 
                    v-model="exportType" 
                    class="peer hidden"
                  />
                  <span
                    class="px-3 py-1 rounded-full border border-gray-500 text-sm 
                          peer-checked:bg-purple-600 peer-checked:text-white 
                          transition-colors"
                  >
                    XLSX
                  </span>
                </label>

                <button
                  class="ml-auto px-4 py-2 bg-emerald-600 hover:bg-emerald-500 rounded-xl disabled:opacity-50"
                  :disabled="!exportType || downloading"
                  @click="onDownload"
                >
                  <span v-if="!downloading">Download</span>
                  <span v-else>Downloading…</span>
                </button>
              </div>
            </section>
          </div>

          <footer class="p-4 border-t border-neutral-700 text-gray-500 flex justify-between">
            <p>Relatórios ChatBot "Bia" BestDoctors</p>
            <p>Dúvidas ou Erros entrar em contato com o Administrador</p>
          </footer>
        </div>
      </Transition>
    </div>
  </Teleport>
</template>

<script setup>
import { downloadReport, requestReport } from '@/api'
import ReportPreviewTable from '@/components/report/ReportPreviewTable.vue'
import { useResponsive } from '@/composables/useResponsive'
import { computed, reactive, ref, watch } from 'vue'

const props = defineProps({ isVisible: { type: Boolean, default: false } })
const emit = defineEmits(['close'])
const closeModal = () => emit('close')

const { isMobile } = useResponsive()

const form = reactive({
  report: '',   // 'session' | 'abandonment' | 'flowDepth' | 'reengagement'
  full: false,  // only used when report === 'reengagement'
})

const dates = reactive({
  from: '', // HTML date 'YYYY-MM-DD'
  to: '',   // HTML date 'YYYY-MM-DD'
})

const preview = ref(null)
const loading = ref(false)
const error = ref('')
const exportType = ref('')       // 'csv' | 'pdf' | 'xlsx'
const previewReady = computed(() => !loading.value && !error.value && preview.value !== null)
const downloading = ref(false)

let controller = null

function toIsoEdge(d, edge='start') {
  if (!d) return null
  const iso = edge === 'start' ? `${d}T00:00:00` : `${d}T23:59:59`
  return new Date(iso).toISOString() // server expects Z
}

async function runPreview() {
  if (!form.report) {
    preview.value = null
    return
  }
  controller?.abort()
  controller = new AbortController()
  loading.value = true
  error.value = ''
  preview.value = null
  try {
    const filters = {
      from: toIsoEdge(dates.from, 'start'),
      to: toIsoEdge(dates.to, 'end'),
      full: form.report === 'reengagement' ? form.full : false,
    }
    const data = await requestReport({ report: form.report, type: 'json', filters }, controller.signal)
    preview.value = data ?? (form.report === 'session' ? [] : {})
  } catch (e) {
    if (e.name !== 'AbortError') error.value = e.message || 'Failed to load preview'
  } finally {
    loading.value = false
  }
}

async function onDownload() {
  if (!exportType.value || !form.report) return
  downloading.value = true
  try {
    const filters = {
      from: toIsoEdge(dates.from, 'start'),
      to: toIsoEdge(dates.to, 'end'),
      full: form.report === 'reengagement' ? form.full : false,
    }
    const { blob, filename } = await downloadReport(
      { report: form.report, type: exportType.value, filters },
      undefined
    )
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = filename || `report-${form.report}.${exportType.value}`
    document.body.appendChild(a)
    a.click()
    a.remove()
    URL.revokeObjectURL(url)
  } catch (e) {
    alert(e?.message || 'Failed to download')
  } finally {
    downloading.value = false
  }
}

// Auto-run preview when inputs change (debounced)
let debounce
watch([() => form.report, () => dates.from, () => dates.to, () => form.full], () => {
  clearTimeout(debounce)
  debounce = setTimeout(runPreview, 300)
})

// Reset state when modal opens/closes
watch(() => props.isVisible, (v) => {
  if (!v) {
    // reset everything on close
    form.report = ''
    form.full = false
    dates.from = ''
    dates.to = ''
    preview.value = null
    error.value = ''
    exportType.value = ''
    controller?.abort()
  }
})
</script>

<style scoped>
@keyframes bounceIn {
  0%   { transform: scale(0.1); opacity: 0; }
  80%  { transform: scale(1.00); opacity: 1; }
  100% { transform: scale(1.00); }
}
@keyframes shrinkFadeOut {
  0%   { transform: scale(1.00); opacity: 1; }
  100% { transform: scale(0.10); opacity: 1; }
}
.bounce-in { animation: bounceIn 520ms cubic-bezier(.2,.8,.2,1) both; }
.shrink-fade-out { animation: shrinkFadeOut 520ms ease-out both; }
</style>
