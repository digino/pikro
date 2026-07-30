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
      <table v-else class="w-full text-xs">
        <thead>
          <tr class="text-text-muted border-b border-border">
            <th class="text-left pb-2 font-medium w-32">Time</th>
            <th class="text-left pb-2 font-medium w-32">Topic</th>
            <th class="text-left pb-2 font-medium">Message</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(entry, i) in logs" :key="i" class="border-b border-border/40 last:border-0">
            <td class="py-2 font-mono text-text-muted whitespace-nowrap">{{ entry.time }}</td>
            <td class="py-2">
              <span class="inline-flex items-center gap-1">
                <span class="size-1.5 rounded-full shrink-0" :style="`background:${topicColor(entry.topics)}`" />
                <span class="font-mono text-text-secondary">{{ entry.topics }}</span>
              </span>
            </td>
            <td class="py-2 text-text-primary">{{ entry.message }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </PageLayout>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { ArrowPathIcon } from '@heroicons/vue/24/outline'
import NoRouterSelected from '@/components/NoRouterSelected.vue'
import EmptyState from '@/components/EmptyState.vue'
import { useRoutersStore } from '@/stores/routers'
import { getSystemLogs } from '@/api'
import PageLayout from '@/components/PageLayout.vue'

const store = useRoutersStore()
const logs = ref<Record<string, string>[]>([])
const loading = ref(false)

async function load() {
  if (!store.activeId) return
  loading.value = true
  try {
    const all = await getSystemLogs(store.activeId)
    logs.value = all.filter(e => (e.topics ?? '').toLowerCase().includes('hotspot'))
  } catch {
    // non-critical
  } finally {
    loading.value = false
  }
}

watch(() => store.activeId, (id) => { if (id) load() }, { immediate: true })

function topicColor(topics: string | undefined): string {
  if (!topics) return 'var(--color-text-muted)'
  const t = topics.toLowerCase()
  if (t.includes('error') || t.includes('critical')) return 'var(--color-red)'
  if (t.includes('warning')) return 'var(--color-amber)'
  if (t.includes('info')) return 'var(--color-green)'
  return 'var(--color-text-muted)'
}
</script>
