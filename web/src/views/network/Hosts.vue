<template>
  <PageLayout title="Network" subtitle="Hosts">
    <div class="rounded-xl border border-border p-5 bg-surface">
      <div class="flex items-center justify-between gap-3 mb-4">
        <div class="flex items-center gap-2">
          <span class="font-semibold text-text-primary">Connected hosts</span>
          <span
            v-if="hosts.length > 0"
            class="inline-flex items-center whitespace-nowrap shrink-0 text-xs px-2 py-0.5 rounded-full text-text-secondary bg-muted"
          >
            {{ filtered.length }} total
          </span>
        </div>
        <div class="w-64 shrink-0">
          <input
            v-model="search"
            class="input"
            placeholder="Search hosts…"
          />
        </div>
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
      <EmptyState
        v-else-if="filtered.length === 0"
        message="No hosts match your search."
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
              v-for="host in paged"
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
        <TablePager
          :page="page"
          :page-count="pageCount"
          :page-size="pageSize"
          :total="filtered.length"
          @update:page="page = $event"
        />
      </template>
    </div>
  </PageLayout>
</template>

<script setup lang="ts">
import { ref, watch, onUnmounted } from "vue";
import EmptyState from "@/components/EmptyState.vue";
import StatusBadge from "@/components/StatusBadge.vue";
import TablePager from "@/components/TablePager.vue";
import { useRoutersStore } from "@/stores/routers";
import { listHotspotHosts } from "@/api";
import PageLayout from "@/components/PageLayout.vue";
import { useSearchPagination } from "@/composables/useSearchPagination";

const store = useRoutersStore();

const hosts = ref<Record<string, string>[]>([]);
const loading = ref(false);

const { search, filtered, page, pageCount, paged, pageSize } = useSearchPagination(
  hosts,
  (h) => `${h["mac-address"] ?? ""} ${h.address ?? ""} ${h["to-address"] ?? ""} ${h.server ?? ""}`,
);

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
