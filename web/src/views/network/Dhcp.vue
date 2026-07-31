<template>
  <PageLayout title="Network" subtitle="DHCP leases">
    <div class="rounded-xl border border-border p-5 bg-surface">
      <div class="flex items-center justify-between mb-4">
        <span class="font-semibold text-text-primary">DHCP leases</span>
        <span v-if="leases.length > 0" class="text-xs text-text-muted">{{ leases.length }} total</span>
      </div>
      <div v-if="leasesLoading && leases.length === 0" class="flex justify-center py-8">
        <span class="spinner spinner--sm" />
      </div>
      <EmptyState v-else-if="leases.length === 0" message="No DHCP leases found" />
      <template v-else>
        <table class="w-full text-xs">
          <thead>
            <tr class="text-text-muted border-b border-border">
              <th class="text-left pb-2 font-medium">Host</th>
              <th class="text-left pb-2 font-medium">IP</th>
              <th class="text-left pb-2 font-medium">MAC</th>
              <th class="text-left pb-2 font-medium">Status</th>
              <th class="text-right pb-2 font-medium">Expires</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="lease in pagedLeases"
              :key="lease['.id']"
              class="border-b border-border/40 last:border-0"
            >
              <td class="py-1.5 text-text-primary">{{ lease['host-name'] || '—' }}</td>
              <td class="py-1.5 font-mono text-text-primary">{{ lease.address || '—' }}</td>
              <td class="py-1.5 font-mono text-text-muted">{{ lease['mac-address'] || '—' }}</td>
              <td class="py-1.5">
                <StatusBadge
                  :color="lease.status === 'bound' ? 'green' : lease.status === 'waiting' ? 'amber' : 'muted'"
                  :label="leaseStatusLabel(lease.status)"
                />
              </td>
              <td class="py-1.5 text-right font-mono text-text-secondary">{{ lease['expires-after'] || '—' }}</td>
            </tr>
          </tbody>
        </table>
        <div v-if="leases.length > LEASES_PAGE_SIZE" class="flex items-center justify-between pt-3 mt-1 border-t border-border">
          <span class="text-xs text-text-muted">
            {{ (leasesPage - 1) * LEASES_PAGE_SIZE + 1 }}–{{ Math.min(leasesPage * LEASES_PAGE_SIZE, leases.length) }} of {{ leases.length }}
          </span>
          <div class="flex items-center gap-1">
            <button
              class="p-1 rounded hover:bg-surface disabled:opacity-30 transition-colors"
              :disabled="leasesPage === 1"
              @click="leasesPage--"
            >
              <ChevronLeftIcon class="size-3.5" />
            </button>
            <button
              class="p-1 rounded hover:bg-surface disabled:opacity-30 transition-colors"
              :disabled="leasesPage >= leasesPageCount"
              @click="leasesPage++"
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
import { ref, computed, watch, onUnmounted } from 'vue'
import { ChevronLeftIcon, ChevronRightIcon } from '@heroicons/vue/24/outline'
import EmptyState from '@/components/EmptyState.vue'
import StatusBadge from '@/components/StatusBadge.vue'
import { useRoutersStore } from '@/stores/routers'
import { getDHCPLeases } from '@/api'
import PageLayout from '@/components/PageLayout.vue'

const store = useRoutersStore()

const leases = ref<Record<string, string>[]>([])
const leasesLoading = ref(false)
const LEASES_PAGE_SIZE = 10
const leasesPage = ref(1)
const leasesPageCount = computed(() =>
  Math.max(1, Math.ceil(leases.value.length / LEASES_PAGE_SIZE)),
)
const pagedLeases = computed(() => {
  const start = (leasesPage.value - 1) * LEASES_PAGE_SIZE
  return leases.value.slice(start, start + LEASES_PAGE_SIZE)
})

function leaseStatusLabel(status: string | undefined): string {
  if (!status) return '—'
  return status.charAt(0).toUpperCase() + status.slice(1)
}

async function loadLeases() {
  if (!store.activeId) return
  leasesLoading.value = true
  try {
    leases.value = await getDHCPLeases(store.activeId)
  } catch {
    // non-critical
  } finally {
    leasesLoading.value = false
  }
}

let leasesTimer: ReturnType<typeof setInterval>

watch(() => store.activeId, async (id) => {
  clearInterval(leasesTimer)
  leases.value = []
  leasesPage.value = 1
  if (!id) return
  await loadLeases()
  leasesTimer = setInterval(() => { if (!document.hidden) loadLeases() }, 5 * 60_000)
}, { immediate: true })

onUnmounted(() => clearInterval(leasesTimer))
</script>
