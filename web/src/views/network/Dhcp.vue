<template>
  <div v-if="!store.activeId" class="p-6">
    <NoRouterSelected />
  </div>

  <PageLayout v-else :title="t('nav.network')" :subtitle="t('dhcp.title')" icon="dhcp">
    <p class="-mt-3 text-[15px] text-text-secondary">{{ t('dhcp.intro') }}</p>

    <section
      v-if="pools.length"
      :aria-label="t('dhcp.pools')"
      class="grid grid-cols-[repeat(auto-fit,minmax(320px,1fr))] gap-5"
    >
      <div
        v-for="p in pools"
        :key="p.name"
        class="card px-5 py-4.5 flex flex-col gap-2.5"
        :class="{ 'border-dashed shadow-none': p.disabled }"
      >
        <div class="flex flex-wrap justify-between items-baseline gap-2.5">
          <span class="font-bold text-[16px]">{{ p.name }} · {{ p.interface }}</span>
          <span class="font-mono text-xs text-text-muted">{{ p.range }}</span>
        </div>
        <div class="flex items-baseline gap-2">
          <span class="type-display text-[40px] leading-none">{{ p.used }}</span>
          <span class="text-sm text-text-secondary">
            {{ p.size ? t('dhcp.inUse', { size: p.size }) : t('dhcp.leases', p.used) }}
          </span>
        </div>
        <div
          v-if="p.size"
          class="h-3.5 border-2 border-ink rounded-full overflow-hidden bg-paper"
          role="meter"
          :aria-valuenow="p.used"
          aria-valuemin="0"
          :aria-valuemax="p.size"
          :aria-label="t('dhcp.inUse', { size: p.size })"
        >
          <div
            class="h-full"
            :class="p.full ? 'bg-danger' : 'bg-ink'"
            :style="{ width: `${Math.min(100, (p.used / p.size) * 100)}%` }"
          />
        </div>
        <span class="text-[13px] text-text-muted">{{ p.note }}</span>
      </div>
    </section>

    <div class="flex flex-wrap items-center justify-between gap-3">
      <div role="group" :aria-label="t('dhcp.filterLabel')" class="flex flex-wrap gap-2">
        <button
          v-for="f in FILTERS"
          :key="f"
          type="button"
          :aria-pressed="filter === f"
          class="inline-flex items-center gap-2 min-h-11 px-4 rounded-full border-2 border-ink text-sm font-semibold cursor-pointer transition-colors"
          :class="
            filter === f
              ? 'bg-ink text-paper shadow-[3px_3px_0_var(--color-primary)]'
              : 'bg-paper text-ink hover:bg-primary-wash'
          "
          @click="filter = f"
        >
          {{ t(`dhcp.filters.${f}`) }}
          <span class="font-mono text-xs opacity-75">{{ counts[f] }}</span>
        </button>
      </div>

      <label
        class="flex items-center gap-2 min-h-11 min-w-65 px-3.5 rounded-control border-2 border-ink bg-paper focus-within:outline-3 focus-within:outline-primary focus-within:outline-offset-1"
      >
        <MagnifyingGlassIcon class="size-4.5 shrink-0" />
        <span class="sr-only">{{ t('dhcp.search') }}</span>
        <input
          v-model="search"
          type="search"
          :placeholder="t('dhcp.searchPlaceholder')"
          class="flex-1 min-w-0 border-0 bg-transparent text-sm outline-none"
        />
      </label>
    </div>

    <section class="card overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full min-w-215 text-sm">
          <thead>
            <tr class="border-b-2 border-ink">
              <th class="eyebrow text-left px-4 py-3.5">{{ t('dhcp.colDevice') }}</th>
              <th class="eyebrow text-left px-4 py-3.5">{{ t('dhcp.colIp') }}</th>
              <th class="eyebrow text-left px-4 py-3.5">{{ t('dhcp.colMac') }}</th>
              <th class="eyebrow text-left px-4 py-3.5">{{ t('dhcp.colServer') }}</th>
              <th class="eyebrow text-left px-4 py-3.5">{{ t('dhcp.colStatus') }}</th>
              <th class="eyebrow text-right px-4 py-3.5">{{ t('dhcp.colExpires') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="loading && leases.length === 0">
              <td colspan="6" class="px-4 py-10 text-center">
                <span class="spinner spinner--sm" />
              </td>
            </tr>
            <tr v-else-if="filtered.length === 0">
              <td colspan="6" class="px-4 py-10 text-center text-text-muted">
                {{ leases.length ? t('dhcp.noMatch') : t('dhcp.noLeases') }}
                <button
                  v-if="leases.length"
                  type="button"
                  class="btn btn-sm ml-3"
                  @click="clearFilters"
                >
                  {{ t('dhcp.clearFilters') }}
                </button>
              </td>
            </tr>
            <tr
              v-for="l in paged"
              :key="l['.id']"
              class="border-b border-hairline last:border-0 hover:bg-primary-wash"
            >
              <td class="px-4 py-3" :class="deviceName(l) ? 'font-semibold' : 'italic text-text-muted'">
                {{ deviceName(l) || t('dhcp.noName') }}
              </td>
              <td class="px-4 py-3 font-mono">{{ l.address || '—' }}</td>
              <td class="px-4 py-3 font-mono text-text-muted">{{ l['mac-address'] || '—' }}</td>
              <td class="px-4 py-3">{{ l.server || l['active-server'] || '—' }}</td>
              <td class="px-4 py-3">
                <span v-if="l.status" class="chip" :class="l.status === 'bound' ? 'bg-primary text-primary-fg' : 'bg-paper'">
                  <span
                    class="size-2 shrink-0 rounded-full"
                    :class="l.status === 'bound' ? 'bg-ink' : 'border-2 border-ink'"
                  />
                  {{ te(`dhcp.status.${l.status}`) ? t(`dhcp.status.${l.status}`) : l.status }}
                </span>
                <span v-else>—</span>
              </td>
              <td class="px-4 py-3 text-right">{{ expires(l) }}</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div
        v-if="filtered.length"
        class="flex flex-wrap items-center justify-between gap-3 px-5 py-3.5 border-t-2 border-ink text-[13px]"
      >
        <span>
          {{ t('dhcp.showing', {
            from: (page - 1) * pageSize + 1,
            to: Math.min(page * pageSize, filtered.length),
            total: filtered.length,
          }) }}
        </span>
        <PaginationRoot
          v-if="pageCount > 1"
          :page="page"
          :items-per-page="pageSize"
          :total="filtered.length"
          class="flex gap-2"
          @update:page="page = $event"
        >
          <PaginationPrev class="btn btn-sm" :aria-label="t('dhcp.previous')">
            <ChevronLeftIcon class="size-4" />
          </PaginationPrev>
          <PaginationNext class="btn btn-sm" :aria-label="t('dhcp.next')">
            <ChevronRightIcon class="size-4" />
          </PaginationNext>
        </PaginationRoot>
      </div>
    </section>
  </PageLayout>
</template>

<script setup lang="ts">
import { ref, computed, watch, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { PaginationRoot, PaginationPrev, PaginationNext } from 'reka-ui'
import { MagnifyingGlassIcon, ChevronLeftIcon, ChevronRightIcon } from '@heroicons/vue/24/outline'
import { useRoutersStore } from '@/stores/routers'
import { getDHCPLeases, getDHCPServers, type DHCPServer } from '@/api'
import PageLayout from '@/components/PageLayout.vue'
import NoRouterSelected from '@/components/NoRouterSelected.vue'
import { useSearchPagination } from '@/composables/useSearchPagination'
import { spacedDuration } from '@/utils/duration'

type Filter = 'all' | 'bound' | 'waiting'
const FILTERS: Filter[] = ['all', 'bound', 'waiting']
// Above this share of the pool the bar turns red.
const FULL_RATIO = 0.85

const { t, te } = useI18n()
const store = useRoutersStore()

const leases = ref<Record<string, string>[]>([])
const servers = ref<DHCPServer[]>([])
const loading = ref(false)
const filter = ref<Filter>('all')

const isStatic = (l: Record<string, string>) => l.dynamic !== 'true'
const serverOf = (l: Record<string, string>) => l.server || l['active-server']
// Static leases are usually named by the operator's comment, not the device.
const deviceName = (l: Record<string, string>) => l['host-name'] || l.comment

const pools = computed(() =>
  servers.value.map((s) => {
    const own = leases.value.filter(l => serverOf(l) === s.name)
    const statics = own.filter(isStatic).length
    // Dynamic leases are the ones taken from the pool; static ones are
    // reserved by hand and may sit outside it.
    const used = s.size ? own.length - statics : own.length
    const ratio = s.size ? used / s.size : 0
    const note = s.disabled
      ? t('dhcp.notes.disabled')
      : s.size
        ? [
            t(ratio > FULL_RATIO ? 'dhcp.notes.almostFull' : ratio > 0.6 ? 'dhcp.notes.filling' : 'dhcp.notes.plenty'),
            s.leaseTime ? t('dhcp.notes.leaseTime', { time: spacedDuration(s.leaseTime) }) : '',
            statics ? t('dhcp.notes.static', { n: statics }, statics) : '',
          ].filter(Boolean).join(' ')
        : t('dhcp.notes.staticOnly')
    return {
      name: s.name,
      interface: s.interface,
      disabled: s.disabled,
      range: s.ranges.replace(/-/g, ' – ').replace(/,\s*/g, ', '),
      size: s.size,
      used,
      full: ratio > FULL_RATIO,
      note,
    }
  }),
)

const counts = computed(() => ({
  all: leases.value.length,
  bound: leases.value.filter(l => l.status === 'bound').length,
  waiting: leases.value.filter(l => l.status === 'waiting').length,
}))

const byStatus = computed(() =>
  filter.value === 'all' ? leases.value : leases.value.filter(l => l.status === filter.value),
)

const { search, filtered, page, pageCount, paged, pageSize, reset } = useSearchPagination(
  byStatus,
  (l) => `${deviceName(l) ?? ''} ${l.address ?? ''} ${l['mac-address'] ?? ''} ${serverOf(l) ?? ''}`,
)

watch(filter, () => { page.value = 1 })

function clearFilters() {
  reset()
  filter.value = 'all'
}

function expires(l: Record<string, string>) {
  if (isStatic(l)) return t('dhcp.static')
  const after = l['expires-after']
  return after && after !== 'never' ? t('dhcp.expiresIn', { time: spacedDuration(after) }) : '—'
}

async function load() {
  const id = store.activeId
  if (!id) return
  loading.value = true
  try {
    const [l, s] = await Promise.all([getDHCPLeases(id), getDHCPServers(id)])
    if (id !== store.activeId) return
    leases.value = l
    servers.value = s
  } catch {
    // non-critical
  } finally {
    loading.value = false
  }
}

let timer: ReturnType<typeof setInterval>

watch(
  () => store.activeId,
  async (id) => {
    clearInterval(timer)
    leases.value = []
    servers.value = []
    clearFilters()
    if (!id) return
    await load()
    timer = setInterval(() => {
      if (!document.hidden) load()
    }, 5 * 60_000)
  },
  { immediate: true },
)

onUnmounted(() => clearInterval(timer))
</script>
