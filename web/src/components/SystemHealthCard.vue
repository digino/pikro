<template>
  <div v-if="loading" class="flex items-center gap-2 text-sm text-text-muted py-4 justify-center">
    <span class="spinner spinner--sm" /> Loading…
  </div>
  <div v-else-if="error" class="flex items-center gap-1.5 text-xs text-red py-4 justify-center">
    <ExclamationTriangleIcon class="size-3.5 shrink-0" />{{ error }}
  </div>
  <template v-else>
    <div class="flex flex-col items-center gap-2">
      <RouterArt
        :board-name="resource['board-name'] ?? ''"
        :size="88"
        :power-led="healthColor"
        :wifi-led="activeSessions > 0 ? 'var(--color-green)' : 'var(--color-border)'"
        wan-led="var(--color-amber)"
      />
      <div class="text-center space-y-0.5">
        <div class="text-sm font-semibold text-text-primary font-mono leading-tight">
          {{ resource["board-name"] || routerName || "—" }}
        </div>
        <div class="text-xs text-text-muted font-mono leading-tight">
          RouterOS {{ resource["version"]?.split(" ")[0] || "—" }}
        </div>
      </div>
      <div class="relative mt-1">
        <svg width="72" height="72" viewBox="0 0 120 120">
          <circle cx="60" cy="60" r="50" fill="none" stroke="var(--color-border)" stroke-width="8" />
          <circle
            cx="60"
            cy="60"
            r="50"
            fill="none"
            :stroke="healthColor"
            stroke-width="8"
            stroke-linecap="round"
            :stroke-dasharray="ringCirc"
            :stroke-dashoffset="ringOffset"
            transform="rotate(-90 60 60)"
            style="transition: stroke-dashoffset 0.6s ease, stroke 0.4s ease"
          />
        </svg>
        <div class="absolute inset-0 flex flex-col items-center justify-center">
          <span class="font-mono font-bold tracking-tight text-text-primary">{{ healthScore }}</span>
          <span class="text-xs font-bold text-text-secondary">Health</span>
        </div>
      </div>
    </div>

    <div class="grid grid-cols-1 text-xs">
      <div class="flex justify-between items-center py-1.5 border-b border-border/60">
        <span class="text-text-muted">Uptime</span>
        <span class="font-mono font-semibold text-text-primary">{{ resource["uptime"] ?? "—" }}</span>
      </div>
      <div class="flex justify-between items-center py-1.5 border-b border-border/60">
        <span class="text-text-muted">CPU</span>
        <span class="font-mono font-semibold text-text-primary">{{ resource["cpu-load"] ?? "—" }}%</span>
      </div>
      <div class="flex justify-between items-center py-1.5 border-b border-border/60">
        <span class="text-text-muted">Free RAM</span>
        <span class="font-mono font-semibold text-text-primary">{{ formatBytes(freeMemory) }}</span>
      </div>
      <div class="flex justify-between items-center py-1.5">
        <span class="text-text-muted">Free disk</span>
        <span class="font-mono font-semibold text-text-primary">{{
          formatBytes(parseInt(resource["free-hdd-space"] ?? "0") || 0)
        }}</span>
      </div>
    </div>
  </template>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { ExclamationTriangleIcon } from '@heroicons/vue/24/outline'
import RouterArt from '@/components/router-art/RouterArt.vue'

const props = defineProps<{
  loading: boolean
  error: string
  resource: Record<string, string>
  routerName: string
  activeSessions: number
  hasActiveId: boolean
}>()

const cpuLoad = computed(() => parseInt(props.resource['cpu-load'] ?? '0') || 0)
const freeMemory = computed(() => parseInt(props.resource['free-memory'] ?? '0') || 0)
const totalMemory = computed(() => parseInt(props.resource['total-memory'] ?? '1') || 1)

const healthScore = computed(() => {
  if (!props.hasActiveId || props.error) return '—'
  const cpuScore = Math.round((1 - cpuLoad.value / 100) * 40)
  const ramScore = Math.round(Math.min(freeMemory.value / totalMemory.value / 0.5, 1) * 30)
  return String(Math.min(100, cpuScore + ramScore + 30))
})
const healthColor = computed(() => {
  if (props.error) return 'var(--color-red)'
  const s = parseInt(healthScore.value) || 0
  if (s >= 75) return 'var(--color-green)'
  if (s >= 50) return 'var(--color-amber)'
  return 'var(--color-red)'
})
const ringCirc = (2 * Math.PI * 50).toFixed(1)
const ringOffset = computed(() => {
  const s = parseInt(healthScore.value) || 0
  return (2 * Math.PI * 50 * (1 - s / 100)).toFixed(1)
})

function formatBytes(n: number): string {
  if (!n || isNaN(n)) return '—'
  if (n >= 1_073_741_824) return (n / 1_073_741_824).toFixed(1) + ' GB'
  if (n >= 1_048_576) return (n / 1_048_576).toFixed(1) + ' MB'
  if (n >= 1_024) return (n / 1_024).toFixed(0) + ' KB'
  return n + ' B'
}
</script>
