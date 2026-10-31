<template>
  <div v-if="!store.activeId" class="p-6">
    <NoRouterSelected />
  </div>

  <PageLayout v-else :title="t('nav.tools')" :subtitle="t('nav.maintenance')" icon="maintenance">
    <p class="-mt-3 text-[15px] text-text-secondary">{{ t('maintenance.intro') }}</p>

    <section
      v-if="todos"
      aria-labelledby="todo-title"
      class="card flex flex-wrap items-center gap-4.5 px-5.5 py-4.5"
      :class="todos.length > 1 ? 'bg-primary text-primary-fg' : ''"
    >
      <img :src="inkIcon('maintenance')" alt="" width="56" height="56" class="flex-none" />
      <div class="flex-[1_1_260px] min-w-0 flex flex-col gap-1">
        <h2 id="todo-title" class="type-display m-0 text-[26px]">{{ t('maintenance.todo', todos.length) }}</h2>
        <p class="m-0 text-sm">{{ todos.length ? t('maintenance.todoText') : t('maintenance.allGood') }}</p>
      </div>
      <div class="flex flex-wrap gap-2.5">
        <a v-for="item in todos" :key="item.href" class="btn btn-sm" :href="item.href">{{ item.label }}</a>
      </div>
    </section>

    <div class="grid grid-cols-[repeat(auto-fit,minmax(380px,1fr))] gap-6 items-stretch">
      <section
        id="clock"
        class="card p-5.5 flex flex-col gap-3.5"
        :class="{ 'shadow-[5px_5px_0_var(--color-danger)]': clockWrong }"
        aria-labelledby="clock-title"
      >
        <div class="flex flex-wrap justify-between items-center gap-3">
          <h2 id="clock-title" class="type-display m-0 text-2xl">{{ t('maintenance.clockTitle') }}</h2>
          <span v-if="clock" class="chip" :class="clockWrong ? 'bg-danger-wash text-danger-text border-danger-text' : 'bg-paper'">
            <span aria-hidden="true" class="size-2 rounded-full" :class="clockWrong ? 'bg-danger' : 'bg-online'" />
            {{ clockWrong ? t('maintenance.clockWrong') : t('maintenance.clockOk') }}
          </span>
        </div>
        <template v-if="clock">
          <div class="flex flex-col gap-1">
            <span class="type-display text-[52px] leading-none">{{ routerTime.time }}</span>
            <span class="text-sm text-text-secondary">{{ routerTime.date }}</span>
          </div>
          <p class="m-0 text-sm leading-normal text-text-secondary">
            {{ clockWrong ? t('maintenance.clockWrongText') : t('maintenance.clockOkText') }}
          </p>
          <div class="flex flex-wrap items-center gap-3 mt-auto">
            <button
              v-if="clockWrong || !clock.ntpEnabled"
              type="button"
              class="btn"
              :class="clockWrong ? 'btn-primary' : ''"
              :disabled="settingTime"
              @click="fixClock"
            >
              <span v-if="settingTime" class="spinner spinner--sm" />
              {{ t('maintenance.turnOnAutoTime') }}
            </button>
            <span class="text-[13px] text-text-muted">
              {{ clock.ntpEnabled ? t('maintenance.ntpOn', { servers: clock.ntpServers || 'NTP' }) : t('maintenance.ntpOff') }}
            </span>
          </div>
        </template>
        <span v-else-if="clockError" role="alert" class="text-sm text-danger-text">{{ clockError }}</span>
        <span v-else class="spinner spinner--sm" />
      </section>

      <section id="updates" class="card p-5.5 flex flex-col gap-3.5" aria-labelledby="updates-title">
        <div class="flex flex-wrap justify-between items-center gap-3">
          <h2 id="updates-title" class="type-display m-0 text-2xl">{{ t('maintenance.updatesTitle') }}</h2>
          <span v-if="updates" class="chip" :class="updates.available ? 'bg-primary text-primary-fg' : 'bg-paper'">
            {{ updates.available ? t('maintenance.updateAvailable') : t('maintenance.upToDate') }}
          </span>
        </div>
        <dl v-if="updates" class="m-0 grid grid-cols-[auto_1fr] gap-x-4.5 gap-y-1.5 text-[15px]">
          <dt class="eyebrow self-center">{{ t('maintenance.installed') }}</dt>
          <dd class="m-0 font-mono font-semibold">{{ updates.installed }}</dd>
          <dt class="eyebrow self-center">{{ t('maintenance.available') }}</dt>
          <dd class="m-0 font-mono font-semibold">{{ updates.latest }} · {{ updates.channel }}</dd>
        </dl>
        <span v-else-if="updatesError" role="alert" class="text-sm text-danger-text">{{ updatesError }}</span>
        <p class="m-0 text-sm leading-normal text-text-secondary">{{ t('maintenance.updateNote') }}</p>
        <button type="button" class="btn self-start mt-auto" :disabled="checkingUpdates" @click="loadUpdates">
          <span v-if="checkingUpdates" class="spinner spinner--sm" />
          {{ checkingUpdates ? t('maintenance.checking') : t('maintenance.checkAgain') }}
        </button>
      </section>

      <section
        class="p-5.5 flex flex-col gap-3.5 border-2 border-dashed border-danger rounded-card bg-paper"
        aria-labelledby="restart-title"
      >
        <h2 id="restart-title" class="type-display m-0 text-2xl">{{ t('maintenance.restartTitle') }}</h2>
        <p class="m-0 text-sm leading-normal text-text-secondary">{{ t('maintenance.restartText') }}</p>
        <span v-if="restarting" role="status" class="text-[13px] font-semibold">{{ t('maintenance.restarting') }}</span>
        <span v-else-if="uptime" class="text-[13px] text-text-muted">{{ t('maintenance.runningFor', { uptime }) }}</span>
        <button
          type="button"
          class="btn btn-danger self-start mt-auto"
          :disabled="restarting"
          @click="confirming = true"
        >
          <span v-if="restarting" class="spinner spinner--sm" />
          {{ t('maintenance.restart') }}
        </button>
      </section>
    </div>

    <AppDialog v-model:open="confirming" :title="t('maintenance.confirmTitle')">
      <p class="m-0 text-sm leading-normal text-text-secondary">{{ t('maintenance.confirmBody') }}</p>
      <div class="flex justify-end gap-2.5 mt-6">
        <button type="button" class="btn" @click="confirming = false">{{ t('maintenance.cancel') }}</button>
        <button type="button" class="btn btn-danger" @click="restart">{{ t('maintenance.confirmRestart') }}</button>
      </div>
    </AppDialog>
  </PageLayout>
</template>

<script setup lang="ts">
import { ref, computed, watch, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoutersStore } from '@/stores/routers'
import { useToastStore } from '@/stores/toast'
import {
  getRouterClock, enableAutoTime, checkUpdates, rebootRouter, getSystemResource,
  type RouterClock, type UpdateInfo,
} from '@/api'
import PageLayout from '@/components/PageLayout.vue'
import NoRouterSelected from '@/components/NoRouterSelected.vue'
import AppDialog from '@/components/AppDialog.vue'
import { inkIcon } from '@/utils/inkIcons'
import { friendlyError } from '@/utils/errors'
import { spacedDuration } from '@/utils/duration'

// A router clock further off than this from this computer is wrong.
const CLOCK_TOLERANCE_MS = 5 * 60_000
// After a restart: wait before the first check, so the router has gone down,
// then keep checking until it answers or this long has passed.
const RESTART_FIRST_CHECK_MS = 15_000
const RESTART_CHECK_EVERY_MS = 5_000
const RESTART_GIVE_UP_MS = 5 * 60_000

const { t, locale } = useI18n()
const store = useRoutersStore()
const toast = useToastStore()

// ── Clock ───────────────────────────────────────────────────────────────────
const clock = ref<RouterClock | null>(null)
const clockOffBy = ref(0)
const clockError = ref('')
const settingTime = ref(false)

const clockWrong = computed(() => !!clock.value && Math.abs(clockOffBy.value) > CLOCK_TOLERANCE_MS)

// `local` is the router's own wall time: shown as is, in Pikro's language.
const routerTime = computed(() => {
  const d = new Date(clock.value?.local ?? 0)
  return {
    time: d.toLocaleTimeString(locale.value, { hour: '2-digit', minute: '2-digit' }),
    date: d.toLocaleDateString(locale.value, { weekday: 'long', day: 'numeric', month: 'long', year: 'numeric' }),
  }
})

async function loadClock(id: string) {
  clockError.value = ''
  try {
    const c = await getRouterClock(id)
    if (id !== store.activeId) return
    clock.value = c
    clockOffBy.value = new Date(c.utc).getTime() - Date.now()
  } catch (e) {
    clockError.value = friendlyError(e, t('maintenance.loadFailed'))
  }
}

async function fixClock() {
  const id = store.activeId
  if (!id) return
  settingTime.value = true
  try {
    await enableAutoTime(id)
    toast.success(t('maintenance.autoTimeDone'))
    await loadClock(id)
  } catch (e) {
    toast.error(t('maintenance.autoTimeFailed'), friendlyError(e))
  } finally {
    settingTime.value = false
  }
}

// ── Updates ─────────────────────────────────────────────────────────────────
const updates = ref<UpdateInfo | null>(null)
const updatesError = ref('')
const checkingUpdates = ref(false)

async function loadUpdates() {
  const id = store.activeId
  if (!id) return
  checkingUpdates.value = true
  updatesError.value = ''
  try {
    const u = await checkUpdates(id)
    if (id === store.activeId) updates.value = u
  } catch (e) {
    updatesError.value = friendlyError(e, t('maintenance.loadFailed'))
  } finally {
    checkingUpdates.value = false
  }
}

// ── To do ───────────────────────────────────────────────────────────────────
// Shown once both checks have answered, so it doesn't change under the reader.
const todos = computed(() => {
  if (!clock.value || checkingUpdates.value) return null
  const items: { label: string; href: string }[] = []
  if (clockWrong.value) items.push({ label: t('maintenance.todoClock'), href: '#clock' })
  if (updates.value?.available) {
    items.push({ label: t('maintenance.todoUpdate', { version: updates.value.latest }), href: '#updates' })
  }
  return items
})

// ── Restart ─────────────────────────────────────────────────────────────────
const uptime = ref('')
const confirming = ref(false)
const restarting = ref(false)
let restartTimer: ReturnType<typeof setTimeout> | undefined

async function loadUptime(id: string) {
  try {
    const r = await getSystemResource(id)
    if (id === store.activeId) uptime.value = spacedDuration(r.uptime)
    return true
  } catch {
    return false
  }
}

async function restart() {
  confirming.value = false
  const id = store.activeId
  if (!id) return
  try {
    await rebootRouter(id)
  } catch (e) {
    toast.error(t('maintenance.restartFailed'), friendlyError(e))
    return
  }
  restarting.value = true
  const giveUpAt = Date.now() + RESTART_GIVE_UP_MS
  const waitForRouter = async () => {
    if (id !== store.activeId) return
    if (await loadUptime(id)) {
      restarting.value = false
      toast.success(t('maintenance.backOnline'))
      loadClock(id)
      return
    }
    if (Date.now() < giveUpAt) restartTimer = setTimeout(waitForRouter, RESTART_CHECK_EVERY_MS)
    else restarting.value = false
  }
  restartTimer = setTimeout(waitForRouter, RESTART_FIRST_CHECK_MS)
}

watch(() => store.activeId, (id) => {
  clearTimeout(restartTimer)
  restarting.value = false
  clock.value = null
  updates.value = null
  uptime.value = ''
  if (!id) return
  loadClock(id)
  loadUpdates()
  loadUptime(id)
}, { immediate: true })

onUnmounted(() => clearTimeout(restartTimer))
</script>
