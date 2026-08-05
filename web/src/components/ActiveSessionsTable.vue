<template>
  <div v-if="loading && sessions.length === 0" class="flex justify-center py-6">
    <span class="spinner spinner--sm" />
  </div>
  <EmptyState v-else-if="sessions.length === 0" message="No active sessions" />
  <div v-else class="overflow-x-auto">
    <table class="w-full text-xs">
      <thead>
        <tr class="text-text-muted border-b border-border">
          <th class="text-left pb-2 font-medium">User</th>
          <th class="text-left pb-2 font-medium">IP</th>
          <th class="text-right pb-2 font-medium">Uptime</th>
          <th class="text-right pb-2 font-medium">Down</th>
          <th class="text-right pb-2 font-medium">Up</th>
          <th class="text-right pb-2 font-medium">Left</th>
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="s in sessions.slice(0, 8)"
          :key="s['.id']"
          class="border-b border-border/40 last:border-0"
        >
          <td class="py-2 font-mono font-medium text-text-primary">{{ s.user }}</td>
          <td class="py-2 font-mono text-text-secondary">{{ s.address }}</td>
          <td class="py-2 text-right font-mono text-text-secondary">{{ s.uptime || "—" }}</td>
          <td class="py-2 text-right font-mono text-text-primary">
            {{ formatBytes(parseInt(s["bytes-in"] ?? "0")) }}
          </td>
          <td class="py-2 text-right font-mono text-text-secondary">
            {{ formatBytes(parseInt(s["bytes-out"] ?? "0")) }}
          </td>
          <td class="py-2 text-right font-mono text-text-secondary">
            {{ s["session-time-left"] || "—" }}
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
import EmptyState from '@/components/EmptyState.vue'

defineProps<{
  sessions: Record<string, string>[]
  loading: boolean
}>()

function formatBytes(n: number): string {
  if (!n || isNaN(n)) return '—'
  if (n >= 1_073_741_824) return (n / 1_073_741_824).toFixed(1) + ' GB'
  if (n >= 1_048_576) return (n / 1_048_576).toFixed(1) + ' MB'
  if (n >= 1_024) return (n / 1_024).toFixed(0) + ' KB'
  return n + ' B'
}
</script>
