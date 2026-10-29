<template>
  <div v-if="!store.activeId" class="p-6">
    <NoRouterSelected />
  </div>

  <PageLayout v-else :title="t('nav.network')" :subtitle="t('nav.hosts')" icon="hosts">
    <p class="-mt-3 text-[15px] text-text-secondary">{{ t('hosts.intro') }}</p>

    <section
      v-if="counts.notLoggedIn > 0"
      class="card bg-primary text-primary-fg px-5.5 py-4.5 flex flex-wrap items-center gap-4.5"
    >
      <img :src="inkIcon('vouchers')" alt="" width="56" height="56" class="flex-none" />
      <p class="m-0 flex-[1_1_320px] text-[16px] leading-normal">
        <b class="type-display text-[22px]">{{ t('hosts.bannerCount', { n: counts.notLoggedIn }, counts.notLoggedIn) }}</b>
        {{ t('hosts.bannerRest', counts.notLoggedIn) }} {{ t('hosts.sellHint', counts.notLoggedIn) }}
      </p>
      <RouterLink to="/hotspot/vouchers" class="btn btn-ink">{{ t('hosts.printVouchers') }}</RouterLink>
    </section>

    <div class="flex flex-wrap items-center justify-between gap-3">
      <div role="group" :aria-label="t('hosts.filterLabel')" class="flex flex-wrap gap-2">
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
          {{ t(`hosts.filters.${f}`) }}
          <span class="font-mono text-xs opacity-75">{{ counts[f] }}</span>
        </button>
      </div>

      <label
        class="flex items-center gap-2 min-h-11 min-w-65 px-3.5 rounded-control border-2 border-ink bg-paper focus-within:outline-3 focus-within:outline-primary focus-within:outline-offset-1"
      >
        <MagnifyingGlassIcon class="size-4.5 shrink-0" />
        <span class="sr-only">{{ t('hosts.search') }}</span>
        <input
          v-model="search"
          type="search"
          :placeholder="t('hosts.searchPlaceholder')"
          class="flex-1 min-w-0 border-0 bg-transparent text-sm outline-none"
        />
      </label>
    </div>

    <section class="card overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full min-w-225 text-sm">
          <thead>
            <tr class="border-b-2 border-ink">
              <th class="eyebrow text-left px-4 py-3.5">{{ t('hosts.colMac') }}</th>
              <th class="eyebrow text-left px-4 py-3.5">{{ t('hosts.colIp') }}</th>
              <th class="eyebrow text-left px-4 py-3.5">{{ t('hosts.colTo') }}</th>
              <th class="eyebrow text-left px-4 py-3.5">{{ t('hosts.colServer') }}</th>
              <th class="eyebrow text-left px-4 py-3.5">{{ t('hosts.colStatus') }}</th>
              <th class="eyebrow text-right px-4 py-3.5">{{ t('hosts.colUptime') }}</th>
              <th class="eyebrow text-right px-4 py-3.5">{{ t('hosts.colIdle') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="loading && hosts.length === 0">
              <td colspan="7" class="px-4 py-10 text-center">
                <span class="spinner spinner--sm" />
              </td>
            </tr>
            <tr v-else-if="filtered.length === 0">
              <td colspan="7" class="px-4 py-10 text-center text-text-muted">
                {{ hosts.length ? t('hosts.noMatch') : t('hosts.noHosts') }}
                <button
                  v-if="hosts.length"
                  type="button"
                  class="btn btn-sm ml-3"
                  @click="clearFilters"
                >
                  {{ t('hosts.clearFilters') }}
                </button>
              </td>
            </tr>
            <tr
              v-for="h in paged"
              :key="h['.id']"
              class="border-b border-hairline last:border-0 hover:bg-primary-wash"
            >
              <td class="px-4 py-3 font-mono font-semibold">{{ h['mac-address'] || '—' }}</td>
              <td class="px-4 py-3 font-mono">{{ h.address || '—' }}</td>
              <td class="px-4 py-3 font-mono text-text-muted">{{ h['to-address'] || '—' }}</td>
              <td class="px-4 py-3">{{ h.server || '—' }}</td>
              <td class="px-4 py-3">
                <span class="chip" :class="STATUS_CHIP[statusOf(h)].chip">
                  <span class="size-2 shrink-0 rounded-full" :class="STATUS_CHIP[statusOf(h)].dot" />
                  {{ t(`hosts.status.${statusOf(h)}`) }}
                </span>
              </td>
              <td class="px-4 py-3 text-right font-mono">{{ spaced(h.uptime) }}</td>
              <td
                class="px-4 py-3 text-right font-mono"
                :class="{ 'text-text-muted': idleLong(h['idle-time']) }"
              >
                {{ spaced(h['idle-time']) }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div
        v-if="filtered.length"
        class="flex flex-wrap items-center justify-between gap-3 px-5 py-3.5 border-t-2 border-ink text-[13px]"
      >
        <span>
          {{ t('hosts.showing', {
            from: (page - 1) * pageSize + 1,
            to: Math.min(page * pageSize, filtered.length),
            total: filtered.length,
          }) }}
        </span>
        <div class="flex flex-wrap items-center gap-4">
          <span class="text-text-muted">{{ t('hosts.idleNote') }}</span>
          <PaginationRoot
            v-if="pageCount > 1"
            :page="page"
            :items-per-page="pageSize"
            :total="filtered.length"
            class="flex gap-2"
            @update:page="page = $event"
          >
            <PaginationPrev class="btn btn-sm" :aria-label="t('hosts.previous')">
              <ChevronLeftIcon class="size-4" />
            </PaginationPrev>
            <PaginationNext class="btn btn-sm" :aria-label="t('hosts.next')">
              <ChevronRightIcon class="size-4" />
            </PaginationNext>
          </PaginationRoot>
        </div>
      </div>
    </section>
  </PageLayout>
</template>

<script setup lang="ts">
import { ref, computed, watch, onUnmounted } from 'vue'
import { RouterLink } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { PaginationRoot, PaginationPrev, PaginationNext } from 'reka-ui'
import { MagnifyingGlassIcon, ChevronLeftIcon, ChevronRightIcon } from '@heroicons/vue/24/outline'
import { useRoutersStore } from '@/stores/routers'
import { listHotspotHosts } from '@/api'
import PageLayout from '@/components/PageLayout.vue'
import NoRouterSelected from '@/components/NoRouterSelected.vue'
import { useSearchPagination } from '@/composables/useSearchPagination'
import { inkIcon } from '@/utils/inkIcons'
import { spacedDuration as spaced } from '@/utils/duration'

type Filter = 'all' | 'loggedIn' | 'notLoggedIn'
// Bypassed devices (IP bindings) skip the login page, so they are neither
// guests online nor guests to sell to: they only show under "All".
type Status = 'loggedIn' | 'notLoggedIn' | 'bypassed'

const FILTERS: Filter[] = ['all', 'loggedIn', 'notLoggedIn']
const STATUS_CHIP: Record<Status, { chip: string; dot: string }> = {
  loggedIn: { chip: 'bg-primary text-primary-fg', dot: 'bg-ink' },
  notLoggedIn: { chip: 'bg-paper', dot: 'border-2 border-ink' },
  bypassed: { chip: 'bg-paper border-dashed text-text-secondary', dot: 'border-[1.5px] border-dashed border-text-secondary' },
}

const { t } = useI18n()
const store = useRoutersStore()

const hosts = ref<Record<string, string>[]>([])
const loading = ref(false)
const filter = ref<Filter>('all')

function statusOf(h: Record<string, string>): Status {
  if (h.authorized === 'true') return 'loggedIn'
  return h.bypassed === 'true' ? 'bypassed' : 'notLoggedIn'
}

const counts = computed(() => {
  const c: Record<Filter, number> = { all: hosts.value.length, loggedIn: 0, notLoggedIn: 0 }
  for (const h of hosts.value) {
    const s = statusOf(h)
    if (s !== 'bypassed') c[s]++
  }
  return c
})

const byStatus = computed(() =>
  filter.value === 'all' ? hosts.value : hosts.value.filter(h => statusOf(h) === filter.value),
)

const { search, filtered, page, pageCount, paged, pageSize, reset } = useSearchPagination(
  byStatus,
  (h) => `${h['mac-address'] ?? ''} ${h.address ?? ''} ${h['to-address'] ?? ''} ${h.server ?? ''}`,
)

watch(filter, () => { page.value = 1 })

function clearFilters() {
  reset()
  filter.value = 'all'
}

function idleLong(d?: string) {
  if (!d) return false
  if (/\d+[wdh]/.test(d)) return true
  return parseInt(/(\d+)m/.exec(d)?.[1] ?? '0') >= 10
}

async function load() {
  if (!store.activeId) return
  loading.value = true
  try {
    hosts.value = await listHotspotHosts(store.activeId)
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
    hosts.value = []
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
