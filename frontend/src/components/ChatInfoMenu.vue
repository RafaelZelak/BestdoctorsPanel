<!-- src/components/ChatInfoMenu.vue -->
<template>
  <!-- Fade no overlay -->
  <Transition
    appear
    enter-active-class="transition-opacity duration-150 ease-out"
    enter-from-class="opacity-0"
    enter-to-class="opacity-100"
    leave-active-class="transition-opacity duration-100 ease-in"
    leave-from-class="opacity-100"
    leave-to-class="opacity-0"
  >
    <div
      class="fixed inset-0 z-40"
      @click="$emit('close')"
      aria-hidden="true"
    ></div>
  </Transition>

  <!-- Slide+scale no painel -->
  <Transition
    appear
    enter-active-class="transition ease-out duration-200"
    enter-from-class="opacity-0 -translate-y-2 scale-95"
    enter-to-class="opacity-100 translate-y-0 scale-100"
    leave-active-class="transition ease-in duration-150"
    leave-from-class="opacity-100 translate-y-0 scale-100"
    leave-to-class="opacity-0 -translate-y-2 scale-95"
  >
    <div
      class="absolute right-4 top-18 z-50 w-[420px] rounded-2xl bg-neutral-900/80 backdrop-blur shadow-xl border border-neutral-700"
      role="dialog"
      aria-label="Conversation info"
      @click.stop
    >
      <header class="px-4 py-3 border-b border-neutral-800 flex items-center gap-2">
        <span class="text-sm font-semibold">Conversation info</span>
        <span
          class="ml-auto text-xs px-2 py-0.5 rounded-full border"
          :class="aiActive ? 'border-emerald-500 text-emerald-400' : 'border-neutral-600 text-neutral-400'"
        >
          {{ aiActive ? 'AI active' : 'AI off' }}
        </span>
      </header>

      <section class="px-4 py-3 space-y-3">
        <!-- Card de status da IA -->
        <div class="text-sm border border-neutral-800 rounded-xl p-3 flex items-center gap-3">
          <div class="flex-1">
            <div class="text-neutral-400 text-xs mb-0.5">AI status</div>
            <div class="font-medium" :class="aiActive ? 'text-emerald-400' : 'text-neutral-300'">
              {{ aiActive ? 'AI active' : 'AI off' }}
            </div>
            <div v-if="toggleError" class="text-xs text-red-400 mt-1">{{ toggleError }}</div>
          </div>

          <!-- Toggle switch -->
          <button
            type="button"
            class="relative inline-flex h-6 w-11 items-center rounded-full transition-colors duration-200 border border-neutral-700 disabled:opacity-60 disabled:cursor-not-allowed"
            :class="aiActive ? 'bg-emerald-600/70' : 'bg-neutral-700'"
            :disabled="toggling"
            @click="onToggleAI"
            :aria-pressed="aiActive ? 'true' : 'false'"
            aria-label="Toggle AI"
          >
            <span
              class="inline-block h-5 w-5 transform rounded-full bg-white transition-transform duration-200"
              :class="aiActive ? 'translate-x-5' : 'translate-x-1'"
            />
            <!-- spinner fino sobreposto quando em progresso -->
            <span
              v-if="toggling"
              class="absolute inset-0 grid place-items-center"
            >
              <svg class="animate-spin w-3.5 h-3.5" viewBox="0 0 24 24" fill="none">
                <circle cx="12" cy="12" r="9" stroke="currentColor" stroke-opacity=".25" stroke-width="3" />
                <path d="M21 12a9 9 0 0 1-9 9" stroke="currentColor" stroke-width="3" stroke-linecap="round"/>
              </svg>
            </span>
          </button>
        </div>

        <!-- Lead info -->
        <div class="text-sm border border-neutral-800 rounded-xl p-3">
          <div class="text-neutral-400 text-xs mb-1">Lead:</div>
          <div class="font-medium">
            <span>
                {{ session?.lead_name || formatValue(vars?.nome_do_lead) || '—' }}
            </span>
            <span class="text-xs text-neutral-400 ml-3">
                {{ phone }}
            </span>
          </div>
        </div>

        <ul class="text-sm space-y-2 border border-neutral-800 rounded-xl p-3">
          <li class="flex items-start gap-2">
            <span class="text-neutral-400 min-w-[130px]">Especialidade</span>
            <span class="flex-1">{{ formatValue(vars?.especialidade) }}</span>
          </li>
          <li class="flex items-start gap-2">
            <span class="text-neutral-400 min-w-[130px]">Nº de usuários</span>
            <span class="flex-1">{{ formatValue(vars?.numero_de_usuarios) }}</span>
          </li>
          <li class="flex items-start gap-2">
            <span class="text-neutral-400 min-w-[130px]">Funcionalidades</span>
            <span class="flex-1">
              {{ formatArray(vars?.funcionalidades_desejadas) }}
            </span>
          </li>
          <li class="flex items-start gap-2">
            <span class="text-neutral-400 min-w-[130px]">Plano</span>
            <span class="flex-1">{{ formatValue(vars?.plano) }}</span>
          </li>
          <li class="flex items-start gap-2">
            <span class="text-neutral-400 min-w-[130px]">Posição hierárq.</span>
            <span class="flex-1">{{ formatValue(vars?.posicao_hierarquica) }}</span>
          </li>
          <li class="flex items-start gap-2">
            <span class="text-neutral-400 min-w-[130px]">Desinteresse</span>
            <span class="flex-1">{{ formatValue(vars?.desinteresse) }}</span>
          </li>
          <li class="flex items-start gap-2">
            <span class="text-neutral-400 min-w-[130px]">Saudação enviada</span>
            <span class="flex-1">{{ formatBoolean(vars?.saudacao_enviada) }}</span>
          </li>
          <li class="flex items-start gap-2">
            <span class="text-neutral-400 min-w-[130px]">Finalizar</span>
            <span class="flex-1">{{ formatBoolean(vars?.finalizar) }}</span>
          </li>
        </ul>
      </section>
    </div>
  </Transition>
</template>

<script setup>
import { computed, ref, watch } from 'vue'
import { prettyPhone } from '@/utils/formatters'
import { toggleSessionActive } from '@/api'

const props = defineProps({
  session: { type: Object, default: null },
  vars: { type: Object, default: null },
})

const emit = defineEmits(['close', 'updated'])

const phone = computed(() => prettyPhone(props.session?.phone))

// Estado local do AI active refletindo props, mas sem "otimismo".
const aiActive = ref(!!props.session?.ai_active)
watch(
  () => props.session?.ai_active,
  (v) => { aiActive.value = !!v },
  { immediate: true }
)

const toggling = ref(false)
const toggleError = ref('')
let toggleController = null

async function onToggleAI() {
  if (!props.session?.session_id || toggling.value) return
  toggleError.value = ''
  toggling.value = true

  // cancela requisição anterior, se houver
  try { toggleController?.abort?.() } catch {}
  toggleController = new AbortController()

  try {
    const res = await toggleSessionActive(props.session.session_id, toggleController.signal)
    aiActive.value = !!res.ai_active
    emit('updated', res) // informa o pai para refletir no header/listas
  } catch (e) {
    toggleError.value = e?.message || 'Failed to toggle AI'
  } finally {
    toggling.value = false
  }
}

function formatBoolean(v) {
  if (v === true) return 'Yes'
  if (v === false) return 'No'
  return '—'
}
function formatArray(a) {
  if (!Array.isArray(a) || a.length === 0) return '—'
  return a.join(', ')
}
function formatValue(v) {
  if (v === null || v === undefined || v === '') return '—'
  return String(v)
}
</script>
