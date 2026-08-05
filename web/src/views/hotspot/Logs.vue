<template>
  <PageLayout title="Hotspot" subtitle="Logs">
    <template #actions>
      <button class="btn btn-ghost btn-sm" :disabled="loading" @click="load">
        <ArrowPathIcon class="size-3.5" :class="{ 'animate-spin': loading }" />
        Refresh
      </button>
    </template>

    <NoRouterSelected v-if="!store.activeId" />

    <div v-else class="rounded-xl border border-border p-5 bg-surface">
      <div class="flex items-center justify-between gap-3 mb-4">
        <div class="flex items-center gap-2">
          <span class="font-semibold text-text-primary">Logs</span>
          <span
            v-if="logs.length > 0"
            class="inline-flex items-center whitespace-nowrap shrink-0 text-xs px-2 py-0.5 rounded-full text-text-secondary bg-muted"
          >
            {{ filtered.length }} total
          </span>
        </div>
        <div class="w-64 shrink-0">
          <input
            v-model="search"
            class="input"
            placeholder="Search logs…"
          />
        </div>
      </div>

      <div v-if="loading && logs.length === 0" class="flex justify-center py-8">
        <span class="spinner spinner--sm" />
      </div>
      <EmptyState
        v-else-if="logs.length === 0"
        message="No hotspot log entries found."
      />
      <EmptyState
        v-else-if="filtered.length === 0"
        message="No log entries match your search."
      />
      <template v-else>
        <table class="w-full text-sm">
          <thead>
            <tr class="text-text-muted border-b border-border">
              <th class="text-left pb-2 font-medium w-32">Time</th>
              <th class="text-left pb-2 font-medium w-28">User</th>
              <th class="text-left pb-2 font-medium">Message</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="(entry, i) in paged"
              :key="i"
              class="border-b border-border/40 last:border-0"
            >
              <td class="py-1.5 font-mono text-text-muted whitespace-nowrap">
                {{ entry.time }}
              </td>
              <td class="py-1.5 font-mono text-text-primary">
                {{ entry.user || "—" }}
              </td>
              <td class="py-1.5 text-text-secondary">{{ entry.message }}</td>
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
import { ref, watch } from "vue";
import { ArrowPathIcon } from "@heroicons/vue/24/outline";
import NoRouterSelected from "@/components/NoRouterSelected.vue";
import EmptyState from "@/components/EmptyState.vue";
import TablePager from "@/components/TablePager.vue";
import { useRoutersStore } from "@/stores/routers";
import { getSystemLogs } from "@/api";
import PageLayout from "@/components/PageLayout.vue";
import { useSearchPagination } from "@/composables/useSearchPagination";

const store = useRoutersStore();

interface HotspotLogEntry {
  time: string;
  message: string;
  user: string;
}

const logs = ref<HotspotLogEntry[]>([]);
const loading = ref(false);

const { search, filtered, page, pageCount, paged, pageSize } = useSearchPagination(
  logs,
  (e) => `${e.user} ${e.message} ${e.time}`,
);

// Hotspot log messages look like "<user> (<ip>): <action>", e.g.
// "gbiu (192.168.88.238): logged in" or "a3fm (192.168.88.238): login failed: invalid username or password".
const USER_MESSAGE_RE = /^(\S+)\s+\([^)]+\):\s*(.*)$/;

function parseEntry(e: Record<string, string>): HotspotLogEntry {
  const message = e.message ?? "";
  const match = message.match(USER_MESSAGE_RE);
  return {
    time: e.time ?? "",
    user: match?.[1] ?? "",
    message: match ? match[2] : message,
  };
}

async function load() {
  if (!store.activeId) return;
  loading.value = true;
  try {
    const all = await getSystemLogs(store.activeId);
    logs.value = all
      .filter((e) => (e.topics ?? "").toLowerCase().includes("hotspot"))
      .map(parseEntry);
    page.value = 1;
  } catch {
    // non-critical
  } finally {
    loading.value = false;
  }
}

watch(
  () => store.activeId,
  (id) => {
    if (id) load();
  },
  { immediate: true },
);
</script>
