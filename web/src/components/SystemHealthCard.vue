<template>
  <div class="card px-5 py-4.5 flex items-center gap-4">
    <img :src="inkIcon('dashboard')" alt="" width="64" height="64" class="shrink-0" />
    <div class="flex-1 min-w-0 flex flex-col gap-2">
      <div class="flex justify-between items-baseline gap-2">
        <span class="eyebrow">{{ t('dashboard.systemHealth') }}</span>
        <span class="type-display text-3xl">
          {{ healthScore ?? '—' }}<span v-if="healthScore !== null" class="text-sm text-text-muted">/100</span>
        </span>
      </div>
      <div
        role="meter"
        aria-valuemin="0"
        aria-valuemax="100"
        :aria-valuenow="healthScore ?? 0"
        :aria-label="t('dashboard.systemHealth')"
        class="h-3 rounded-full border-2 border-ink bg-paper overflow-hidden"
      >
        <div
          class="h-full bg-ink transition-[width] duration-500"
          :style="{ width: `${healthScore ?? 0}%` }"
        />
      </div>
      <div v-if="error" class="flex items-center gap-1.5 text-xs text-red">
        <ExclamationTriangleIcon class="size-3.5 shrink-0" />{{ error }}
      </div>
      <div v-else class="font-mono text-xs text-text-muted truncate">
        <span v-if="loading && !resource['uptime']" class="spinner spinner--sm" />
        <template v-else>
          {{ t('dashboard.healthDetail', {
            cpu: cpuLoad,
            ram: formatBytes(freeMemory),
            uptime: resource['uptime'] ?? '—',
          }) }}
        </template>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { ExclamationTriangleIcon } from '@heroicons/vue/24/outline'
import { inkIcon } from '@/utils/inkIcons'

const props = defineProps<{
  loading: boolean
  error: string
  resource: Record<string, string>
  hasActiveId: boolean
}>()

const { t } = useI18n()

const cpuLoad = computed(() => parseInt(props.resource['cpu-load'] ?? '0') || 0)
const freeMemory = computed(() => parseInt(props.resource['free-memory'] ?? '0') || 0)
const totalMemory = computed(() => parseInt(props.resource['total-memory'] ?? '1') || 1)

const healthScore = computed(() => {
  if (!props.hasActiveId || props.error || !props.resource['cpu-load']) return null
  const cpuScore = Math.round((1 - cpuLoad.value / 100) * 40)
  const ramScore = Math.round(Math.min(freeMemory.value / totalMemory.value / 0.5, 1) * 30)
  return Math.min(100, cpuScore + ramScore + 30)
})

function formatBytes(n: number): string {
  if (!n || isNaN(n)) return '—'
  if (n >= 1_073_741_824) return (n / 1_073_741_824).toFixed(1) + ' GB'
  if (n >= 1_048_576) return (n / 1_048_576).toFixed(1) + ' MB'
  if (n >= 1_024) return (n / 1_024).toFixed(0) + ' KB'
  return n + ' B'
}
</script>
