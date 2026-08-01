<template>
  <PageLayout title="Network" subtitle="Hosts">
    <div class="rounded-xl border border-border p-5 bg-surface">
      <div class="flex items-center justify-between mb-4">
        <span class="font-semibold text-text-primary">Connected hosts</span>
        <span v-if="hosts.length > 0" class="text-sm text-text-secon"
          >{{ hosts.length }} total</span
        >
      </div>
      <div
        v-if="loading && hosts.length === 0"
        class="flex justify-center py-8"
      >
        <span class="spinner spinner--sm" />
      </div>
      <EmptyState
        v-else-if="hosts.length === 0"
        message="No connected hosts found"
      />
      <template v-else>
        <table class="w-full text-sm">
          <thead>
            <tr class="border-b border-border">
              <th class="text-left pb-2 font-medium">MAC address</th>
              <th class="text-left pb-2 font-medium">IP address</th>
              <th class="text-left pb-2 font-medium">To address</th>
              <th class="text-left pb-2 font-medium">Server</th>
              <th class="text-left pb-2 font-medium">Status</th>
              <th class="text-right pb-2 font-medium">Uptime</th>
              <th class="text-right pb-2 font-medium">Idle</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="host in pagedHosts"
              :key="host['.id']"
              class="border-b border-border/40 last:border-0"
            >
              <td class="py-1.5 font-mono text-text-primary">
                {{ host["mac-address"] || "—" }}
              </td>
              <td class="py-1.5 font-mono text-text-primary">
                {{ host.address || "—" }}
              </td>
              <td class="py-1.5 font-mono text-text-secondary">
                {{ host["to-address"] || "—" }}
              </td>
              <td class="py-1.5 text-text-secondary">
                {{ host.server || "—" }}
              </td>
              <td class="py-1.5">
                <StatusBadge
                  :color="host.authorized === 'true' ? 'green' : 'muted'"
                  :label="
                    host.authorized === 'true' ? 'Authorized' : 'Unauthorized'
                  "
                />
              </td>
              <td class="py-1.5 text-right font-mono text-text-secondary">
                {{ host.uptime || "—" }}
              </td>
              <td class="py-1.5 text-right font-mono text-text-muted">
                {{ host["idle-time"] || "—" }}
              </td>
            </tr>
          </tbody>
        </table>
        <div
          v-if="hosts.length > PAGE_SIZE"
          class="flex items-center justify-between pt-3 mt-1 border-t border-border"
        >
          <span class="text-xs text-text-muted">
            {{ (page - 1) * PAGE_SIZE + 1 }}–{{
              Math.min(page * PAGE_SIZE, hosts.length)
            }}
            of {{ hosts.length }}
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
import { ref, computed, watch, onUnmounted } from "vue";
import { ChevronLeftIcon, ChevronRightIcon } from "@heroicons/vue/24/outline";
import EmptyState from "@/components/EmptyState.vue";
import StatusBadge from "@/components/StatusBadge.vue";
import { useRoutersStore } from "@/stores/routers";
import { listHotspotHosts } from "@/api";
import PageLayout from "@/components/PageLayout.vue";

const store = useRoutersStore();

const hosts = ref<Record<string, string>[]>([]);
const loading = ref(false);
const PAGE_SIZE = 10;
const page = ref(1);
const pageCount = computed(() =>
  Math.max(1, Math.ceil(hosts.value.length / PAGE_SIZE)),
);
const pagedHosts = computed(() => {
  const start = (page.value - 1) * PAGE_SIZE;
  return hosts.value.slice(start, start + PAGE_SIZE);
});

async function load() {
  if (!store.activeId) return;
  loading.value = true;
  try {
    hosts.value = await listHotspotHosts(store.activeId);
  } catch {
    // non-critical
  } finally {
    loading.value = false;
  }
}

let timer: ReturnType<typeof setInterval>;

watch(
  () => store.activeId,
  async (id) => {
    clearInterval(timer);
    hosts.value = [];
    page.value = 1;
    if (!id) return;
    await load();
    timer = setInterval(() => {
      if (!document.hidden) load();
    }, 5 * 60_000);
  },
  { immediate: true },
);

onUnmounted(() => clearInterval(timer));
</script>
