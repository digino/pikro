<template>
  <div v-if="!store.activeId" class="p-6">
    <NoRouterSelected />
  </div>

  <PageLayout v-else :title="t('nav.tools')" :subtitle="t('nav.troubleshoot')" icon="troubleshoot">
    <p class="-mt-3 text-text-secondary">{{ t('troubleshoot.intro') }}</p>

    <section
      class="card p-6 flex flex-wrap gap-7"
      :class="{ 'shadow-[5px_5px_0_var(--color-danger)]': problem }"
      aria-labelledby="check-title"
    >
      <div class="flex-[1_1_300px] min-w-0 flex flex-col gap-3.5">
        <h2 id="check-title" class="eyebrow m-0">{{ t('troubleshoot.checkTitle') }}</h2>
        <div class="flex items-start gap-3" aria-live="polite">
          <span
            v-if="verdict"
            aria-hidden="true"
            class="flex-none size-3.5 mt-2.5 rounded-full"
            :class="problem
              ? 'bg-danger shadow-[0_0_0_5px_color-mix(in_srgb,var(--color-danger)_18%,transparent)]'
              : 'bg-online shadow-[0_0_0_5px_color-mix(in_srgb,var(--color-online)_18%,transparent)]'"
          />
          <p class="type-display m-0 text-3xl leading-tight">
            {{ verdict ? verdict.title : t('troubleshoot.checking') }}
          </p>
        </div>
        <p v-if="verdict" class="m-0 leading-normal text-text-secondary">{{ verdict.advice }}</p>
        <div class="flex flex-wrap gap-2.5 mt-auto">
          <button type="button" class="btn btn-primary" :disabled="checking" @click="runCheck">
            <span v-if="checking" class="spinner spinner--sm" />
            {{ t('troubleshoot.checkAgain') }}
          </button>
          <button
            v-if="failedStep?.key === 'dns'"
            type="button"
            class="btn btn-ink"
            :disabled="fixing"
            @click="fixDNS"
          >
            <span v-if="fixing" class="spinner spinner--sm" />
            {{ t('troubleshoot.usePublicDns') }}
          </button>
        </div>
        <span v-if="check" class="text-xs text-text-muted">
          {{ t('troubleshoot.checkedAt', { time: clockTime(new Date(check.checkedAt)) }) }}
        </span>
        <RouterLink
          v-if="!checkError"
          to="/tools/maintenance"
          class="inline-flex items-center gap-2 font-semibold text-ink"
        >
          <img :src="inkIcon('maintenance')" alt="" width="24" height="24" />
          {{ t('troubleshoot.stillStuck') }}
        </RouterLink>
      </div>

      <ol
        :aria-label="t('troubleshoot.stepsLabel')"
        class="flex-[1_1_380px] m-0 p-0 list-none flex flex-col border-2 border-ink rounded-control overflow-hidden"
      >
        <li
          v-for="(s, i) in rows"
          :key="s.key"
          class="flex items-center gap-3.5 px-4 py-3.5"
          :class="[i > 0 && 'border-t-[1.5px] border-hairline', s.state === 'fail' ? 'bg-danger-wash' : 'bg-paper']"
        >
          <span
            aria-hidden="true"
            class="flex-none size-7.5 rounded-full border-2 border-ink inline-flex items-center justify-center font-extrabold"
            :class="MARK[s.state].class"
          >
            <span v-if="s.state === 'pending'" class="spinner spinner--sm" />
            <template v-else>{{ MARK[s.state].glyph }}</template>
          </span>
          <span class="flex flex-col gap-0.5 flex-1 min-w-0">
            <span class="font-bold">{{ t(`troubleshoot.steps.${s.key}.title`) }}</span>
            <span class="text-sm text-text-secondary">{{ s.detail }}</span>
          </span>
          <span class="chip" :class="MARK[s.state].chip">{{ t(`troubleshoot.state.${s.state}`) }}</span>
        </li>
      </ol>
    </section>

    <section class="card p-5.5 flex flex-col gap-4 max-w-xl" aria-labelledby="speed-title">
      <h2 id="speed-title" class="type-display m-0 text-2xl">{{ t('troubleshoot.speedTitle') }}</h2>
      <div class="flex flex-col gap-1">
        <span class="eyebrow">{{ t('troubleshoot.lastTest') }}</span>
        <span class="type-display text-[52px] leading-none">
          {{ speed ? speed.mbps : '—' }} <span class="text-xl">Mbps</span>
        </span>
        <span class="text-text-muted">
          {{ speed ? t('troubleshoot.ranAt', { time: speed.at, size: speed.size }) : t('troubleshoot.notRun') }}
        </span>
      </div>
      <div class="flex flex-col gap-2">
        <span :id="sizeLabelId" class="eyebrow">{{ t('troubleshoot.testSize') }}</span>
        <div
          role="group"
          :aria-labelledby="sizeLabelId"
          class="flex self-start border-2 border-ink rounded-control overflow-hidden bg-paper"
        >
          <button
            v-for="(n, i) in SIZES"
            :key="n"
            type="button"
            :aria-pressed="size === n"
            class="min-h-10 px-4 text-sm font-semibold cursor-pointer"
            :class="[i > 0 && 'border-l-2 border-ink', size === n ? 'bg-ink text-paper' : 'text-ink hover:bg-primary-wash']"
            @click="size = n"
          >
            {{ n }} MB
          </button>
        </div>
        <span class="text-text-secondary">{{ t(`troubleshoot.sizeHint${size}`) }}</span>
      </div>
      <button type="button" class="btn btn-ink self-start" :disabled="testing" @click="runSpeed">
        <span v-if="testing" class="spinner spinner--sm" />
        {{ testing ? t('troubleshoot.testing') : t('troubleshoot.runTest') }}
      </button>
      <p v-if="speedError" role="alert" class="m-0 text-sm text-danger-text">{{ speedError }}</p>
      <p class="leading-normal text-text-muted">{{ t('troubleshoot.speedNote') }}</p>
    </section>
  </PageLayout>
</template>

<script setup lang="ts">
import { ref, computed, watch, useId } from 'vue'
import { RouterLink } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useRoutersStore } from '@/stores/routers'
import { useToastStore } from '@/stores/toast'
import {
  getInternetCheck, switchToPublicDNS, runSpeedTest,
  type InternetCheck, type CheckStep,
} from '@/api'
import PageLayout from '@/components/PageLayout.vue'
import NoRouterSelected from '@/components/NoRouterSelected.vue'
import { inkIcon } from '@/utils/inkIcons'
import { friendlyError } from '@/utils/errors'

type RowState = CheckStep['state'] | 'pending'

const MARK: Record<RowState, { glyph: string; class: string; chip: string }> = {
  ok: { glyph: '✓', class: 'bg-ink text-primary', chip: 'bg-primary text-primary-fg' },
  fail: { glyph: '!', class: 'bg-danger text-paper', chip: 'bg-danger-wash text-danger-text border-danger-text' },
  skip: { glyph: '–', class: 'bg-ground text-text-muted border-dashed', chip: 'bg-ground text-text-secondary border-dashed' },
  pending: { glyph: '', class: 'bg-paper', chip: 'bg-paper' },
}
const STEP_KEYS: CheckStep['key'][] = ['cable', 'gateway', 'internet', 'dns']
const SIZES = [10, 50] as const

const { t, locale } = useI18n()
const store = useRoutersStore()
const toast = useToastStore()
const sizeLabelId = useId()

const clockTime = (d: Date) => d.toLocaleTimeString(locale.value, { hour: '2-digit', minute: '2-digit' })

// ── Internet check ──────────────────────────────────────────────────────────
const check = ref<InternetCheck | null>(null)
const checking = ref(false)
const checkError = ref('')
const fixing = ref(false)

const failedStep = computed(() => check.value?.steps.find(s => s.state === 'fail'))
const problem = computed(() => !!failedStep.value || !!checkError.value)

const verdict = computed(() => {
  if (checking.value) return null
  // Pikro itself can't reach the router: that is the answer, not a side error.
  if (checkError.value) {
    return { title: t('troubleshoot.verdict.unreachable.title'), advice: checkError.value }
  }
  if (!check.value) return null
  const f = failedStep.value
  const key = !f ? 'ok' : f.key === 'cable' ? (f.reason ?? 'noLink') : f.key
  const named = { wan: f?.target, gw: f?.target, via: f?.via }
  return {
    title: t(`troubleshoot.verdict.${key}.title`),
    advice: t(`troubleshoot.verdict.${key}.advice`, named),
  }
})

function detail(s: CheckStep) {
  const named = { target: s.target, ms: s.rttMs ?? 0, via: s.via?.replace(/,/g, ', ') || '—' }
  if (s.state === 'skip') return t('troubleshoot.skipped')
  if (s.key === 'cable') {
    return t(`troubleshoot.steps.cable.${s.state === 'ok' ? 'ok' : (s.reason ?? 'noLink')}`, named)
  }
  // No round trip on the gateway step means a PPPoE-style link, not an address.
  if (s.key === 'gateway' && s.state === 'ok' && !s.rttMs) return t('troubleshoot.steps.gateway.okLink', named)
  return t(`troubleshoot.steps.${s.key}.${s.state}`, named)
}

const rows = computed(() =>
  checking.value || !check.value
    ? STEP_KEYS.map(key => ({
        key,
        // Nothing was checked when the router can't be reached.
        state: (checkError.value && !checking.value ? 'skip' : 'pending') as RowState,
        detail: checkError.value && !checking.value ? t('troubleshoot.notChecked') : '',
      }))
    : check.value.steps.map(s => ({ key: s.key, state: s.state as RowState, detail: detail(s) })),
)

async function runCheck() {
  const id = store.activeId
  if (!id) return
  checking.value = true
  checkError.value = ''
  try {
    const result = await getInternetCheck(id)
    if (id === store.activeId) check.value = result
  } catch (e) {
    checkError.value = friendlyError(e, t('troubleshoot.checkFailed'))
  } finally {
    checking.value = false
  }
}

async function fixDNS() {
  if (!store.activeId) return
  fixing.value = true
  try {
    await switchToPublicDNS(store.activeId)
    toast.success(t('troubleshoot.publicDnsDone'))
    await runCheck()
  } catch (e) {
    toast.error(t('troubleshoot.publicDnsFailed'), friendlyError(e))
  } finally {
    fixing.value = false
  }
}

// ── Speed test ──────────────────────────────────────────────────────────────
const size = ref<(typeof SIZES)[number]>(10)
const testing = ref(false)
const speedError = ref('')
const speed = ref<{ mbps: string; at: string; size: number } | null>(null)

async function runSpeed() {
  if (!store.activeId) return
  testing.value = true
  speedError.value = ''
  try {
    const r = await runSpeedTest(store.activeId, String(size.value))
    speed.value = {
      mbps: ((parseInt(r['rx-speed']) || 0) / 1_000_000).toFixed(1),
      at: clockTime(new Date()),
      size: size.value,
    }
  } catch (e) {
    speedError.value = friendlyError(e, t('troubleshoot.speedFailed'))
  } finally {
    testing.value = false
  }
}

watch(() => store.activeId, () => {
  check.value = null
  speed.value = null
  runCheck()
}, { immediate: true })
</script>
