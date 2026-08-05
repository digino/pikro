<template>
  <PageLayout title="Network" subtitle="DHCP leases">
    <div class="rounded-xl border border-border p-5 bg-surface">
      <div class="flex items-center justify-between gap-3 mb-4">
        <div class="flex items-center gap-2">
          <span class="font-semibold text-text-primary">DHCP leases</span>
          <span
            v-if="leases.length > 0"
            class="inline-flex items-center whitespace-nowrap shrink-0 text-xs px-2 py-0.5 rounded-full text-text-secondary bg-muted"
          >
            {{ filtered.length }} total
          </span>
        </div>
        <div class="w-64 shrink-0">
          <input
            v-model="search"
            class="input"
            placeholder="Search leases…"
          />
        </div>
      </div>
      <div
        v-if="leasesLoading && leases.length === 0"
        class="flex justify-center py-8"
      >
        <span class="spinner spinner--sm" />
      </div>
      <EmptyState
        v-else-if="leases.length === 0"
        message="No DHCP leases found"
      />
      <EmptyState
        v-else-if="filtered.length === 0"
        message="No leases match your search."
      />
      <template v-else>
        <table class="w-full text-sm">
          <thead>
            <tr class="border-b border-border">
              <th class="text-left pb-2 font-medium">Host</th>
              <th class="text-left pb-2 font-medium">IP</th>
              <th class="text-left pb-2 font-medium">MAC</th>
              <th class="text-left pb-2 font-medium">Active server</th>
              <th class="text-left pb-2 font-medium">Status</th>
              <th class="text-right pb-2 font-medium">Expires</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="lease in paged"
              :key="lease['.id']"
              class="border-b border-border/40 last:border-0"
            >
              <td class="py-1.5 text-text-primary">
                {{ lease["host-name"] || "—" }}
              </td>
              <td class="py-1.5 font-mono text-text-primary">
                {{ lease.address || "—" }}
              </td>
              <td class="py-1.5 font-mono text-text-muted">
                {{ lease["mac-address"] || "—" }}
              </td>
              <td class="py-1.5 font-mono text-text-secondary">
                {{ lease["active-server"] || "—" }}
              </td>
              <td class="py-1.5">
                <StatusBadge
                  :color="
                    lease.status === 'bound'
                      ? 'green'
                      : lease.status === 'waiting'
                        ? 'amber'
                        : 'muted'
                  "
                  :label="leaseStatusLabel(lease.status)"
                />
              </td>
              <td class="py-1.5 text-right font-mono text-text-secondary">
                {{ lease["expires-after"] || "—" }}
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
import { getDHCPLeases } from "@/api";
import PageLayout from "@/components/PageLayout.vue";
import { useSearchPagination } from "@/composables/useSearchPagination";

const store = useRoutersStore();

const leases = ref<Record<string, string>[]>([]);
const leasesLoading = ref(false);

const { search, filtered, page, pageCount, paged, pageSize } = useSearchPagination(
  leases,
  (l) => `${l["host-name"] ?? ""} ${l.address ?? ""} ${l["mac-address"] ?? ""} ${l["active-server"] ?? ""}`,
);

function leaseStatusLabel(status: string | undefined): string {
  if (!status) return "—";
  return status.charAt(0).toUpperCase() + status.slice(1);
}

async function loadLeases() {
  if (!store.activeId) return;
  leasesLoading.value = true;
  try {
    leases.value = await getDHCPLeases(store.activeId);
  } catch {
    // non-critical
  } finally {
    leasesLoading.value = false;
  }
}

let leasesTimer: ReturnType<typeof setInterval>;

watch(
  () => store.activeId,
  async (id) => {
    clearInterval(leasesTimer);
    leases.value = [];
    page.value = 1;
    if (!id) return;
    await loadLeases();
    leasesTimer = setInterval(() => {
      if (!document.hidden) loadLeases();
    }, 5 * 60_000);
  },
  { immediate: true },
);

onUnmounted(() => clearInterval(leasesTimer));
</script>
