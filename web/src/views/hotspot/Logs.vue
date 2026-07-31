<template>
  <PageLayout title="Hotspot" subtitle="Logs">
    <template #actions>
      <button class="btn btn-ghost btn-sm" :disabled="loading" @click="load">
        <ArrowPathIcon class="size-3.5" :class="{ 'animate-spin': loading }" />
        Refresh
      </button>
    </template>

    <NoRouterSelected v-if="!store.activeId" />

    <div v-else>
      <div v-if="loading && logs.length === 0" class="flex justify-center py-12">
        <span class="spinner" />
      </div>
      <EmptyState v-else-if="logs.length === 0" size="lg" message="No hotspot log entries found." />
      <template v-else>
        <table class="w-full text-xs">
          <thead>
            <tr class="text-text-muted border-b border-border">
              <th class="text-left pb-2 font-medium w-32">Time</th>
              <th class="text-left pb-2 font-medium w-28">User</th>
              <th class="text-left pb-2 font-medium">Message</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(entry, i) in pagedLogs" :key="i" class="border-b border-border/40 last:border-0">
              <td class="py-2 font-mono text-text-muted whitespace-nowrap">{{ entry.time }}</td>
              <td class="py-2 font-mono text-text-primary">{{ entry.user || '—' }}</td>
              <td class="py-2 text-text-secondary">{{ entry.message }}</td>
            </tr>
          </tbody>
        </table>
        <div v-if="logs.length > PAGE_SIZE" class="flex items-center justify-between pt-3 mt-1 border-t border-border">
          <span class="text-xs text-text-muted">
            {{ (page - 1) * PAGE_SIZE + 1 }}–{{ Math.min(page * PAGE_SIZE, logs.length) }} of {{ logs.length }}
          </span>
          <div class="flex items-center gap-1">
            <button
              class="p-1 rounded hover:bg-surface disabled:opacity-30 transition-colors"
              :disabled="page === 1"
              @click="page--"
            >
              <ChevronLeftIcon class="size-3.5" />
            </button>
            <button
              class="p-1 rounded hover:bg-surface disabled:opacity-30 transition-colors"
              :disabled="page >= pageCount"
              @click="page++"
            >
              <ChevronRightIcon class="size-3.5" />
            </button>
          </div>
        </div>
      </template>
    </div>
  </PageLayout>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { ArrowPathIcon, ChevronLeftIcon, ChevronRightIcon } from '@heroicons/vue/24/outline'
import NoRouterSelected from '@/components/NoRouterSelected.vue'
import EmptyState from '@/components/EmptyState.vue'
import { useRoutersStore } from '@/stores/routers'
import { getSystemLogs } from '@/api'
import PageLayout from '@/components/PageLayout.vue'

const store = useRoutersStore()

interface HotspotLogEntry {
  time: string
  message: string
  user: string
}

const logs = ref<HotspotLogEntry[]>([])
const loading = ref(false)

const PAGE_SIZE = 20
const page = ref(1)
const pageCount = computed(() => Math.max(1, Math.ceil(logs.value.length / PAGE_SIZE)))
const pagedLogs = computed(() => {
  const start = (page.value - 1) * PAGE_SIZE
  return logs.value.slice(start, start + PAGE_SIZE)
})

// Hotspot log messages look like "<user> (<ip>): <action>", e.g.
// "gbiu (192.168.88.238): logged in" or "a3fm (192.168.88.238): login failed: invalid username or password".
const USER_MESSAGE_RE = /^(\S+)\s+\([^)]+\):\s*(.*)$/

function parseEntry(e: Record<string, string>): HotspotLogEntry {
  const message = e.message ?? ''
  const match = message.match(USER_MESSAGE_RE)
  return {
    time: e.time ?? '',
    user: match?.[1] ?? '',
    message: match ? match[2] : message,
  }
}

async function load() {
  if (!store.activeId) return
  loading.value = true
  try {
    const all = await getSystemLogs(store.activeId)
    logs.value = all
      .filter(e => (e.topics ?? '').toLowerCase().includes('hotspot'))
      .map(parseEntry)
    page.value = 1
  } catch {
    // non-critical
  } finally {
    loading.value = false
  }
}

watch(() => store.activeId, (id) => { if (id) load() }, { immediate: true })
</script>
