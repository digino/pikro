<template>
  <div v-if="loading && sessions.length === 0" class="flex justify-center py-6">
    <span class="spinner spinner--sm" />
  </div>
  <EmptyState v-else-if="sessions.length === 0" :message="t('dashboard.noActiveSessions')" />
  <div v-else class="overflow-x-auto">
    <table class="w-full min-w-180 text-sm">
      <thead>
        <tr class="border-b-2 border-ink">
          <th class="eyebrow text-left px-4 py-3">{{ t('dashboard.colUser') }}</th>
          <th class="eyebrow text-left px-4 py-3">{{ t('dashboard.colProfile') }}</th>
          <th class="eyebrow text-left px-4 py-3">{{ t('dashboard.colIp') }}</th>
          <th class="eyebrow text-left px-4 py-3">{{ t('dashboard.colMac') }}</th>
          <th class="eyebrow text-left px-4 py-3">{{ t('dashboard.colOnlineFor') }}</th>
          <th class="eyebrow text-left px-4 py-3">{{ t('dashboard.colDataUsed') }}</th>
          <th class="px-4 py-3"><span class="sr-only">{{ t('dashboard.colActions') }}</span></th>
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="s in sessions.slice(0, 8)"
          :key="s['.id']"
          class="border-b border-hairline last:border-0"
        >
          <td class="px-4 py-3.5 font-mono font-semibold">{{ s.user }}</td>
          <td class="px-4 py-3.5">
            <span v-if="userByName.get(s.user)?.profile" class="chip">
              {{ userByName.get(s.user)?.profile }}
            </span>
          </td>
          <td class="px-4 py-3.5 font-mono">{{ s.address }}</td>
          <td class="px-4 py-3.5 font-mono text-text-muted">{{ s['mac-address'] || '—' }}</td>
          <td class="px-4 py-3.5">{{ s.uptime || '—' }}</td>
          <td class="px-4 py-3.5">
            <div class="flex items-center gap-2.5">
              <div
                v-if="quotaPct(s) !== null"
                class="w-22 h-2 shrink-0 rounded-full border-[1.5px] border-ink overflow-hidden"
              >
                <div
                  class="h-full"
                  :class="quotaPct(s)! > 70 ? 'bg-ink' : 'bg-primary'"
                  :style="{ width: `${quotaPct(s)}%` }"
                />
              </div>
              <span class="font-mono text-xs">{{ formatBytes(bytesUsed(s)) }}</span>
            </div>
          </td>
          <td class="px-4 py-3.5 text-right">
            <button
              type="button"
              class="btn btn-sm"
              :aria-label="t('dashboard.disconnectUser', { user: s.user })"
              :disabled="disconnecting === s['.id']"
              @click="$emit('disconnect', s)"
            >
              {{ t('dashboard.disconnect') }}
            </button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import EmptyState from '@/components/EmptyState.vue'

const props = defineProps<{
  sessions: Record<string, string>[]
  users: Record<string, string>[]
  loading: boolean
  disconnecting: string | null
}>()

defineEmits<{ disconnect: [session: Record<string, string>] }>()

const { t } = useI18n()

const userByName = computed(() => new Map(props.users.map(u => [u.name, u])))

function bytesUsed(s: Record<string, string>): number {
  return (parseInt(s['bytes-in'] ?? '0') || 0) + (parseInt(s['bytes-out'] ?? '0') || 0)
}

// Share of the user's data quota used, or null when the user has no quota.
function quotaPct(s: Record<string, string>): number | null {
  const limit = parseInt(userByName.value.get(s.user)?.['limit-bytes-total'] ?? '0') || 0
  if (!limit) return null
  return Math.min(100, Math.round((bytesUsed(s) / limit) * 100))
}

function formatBytes(n: number): string {
  if (!n || isNaN(n)) return '—'
  if (n >= 1_073_741_824) return (n / 1_073_741_824).toFixed(1) + ' GB'
  if (n >= 1_048_576) return (n / 1_048_576).toFixed(1) + ' MB'
  if (n >= 1_024) return (n / 1_024).toFixed(0) + ' KB'
  return n + ' B'
}
</script>
