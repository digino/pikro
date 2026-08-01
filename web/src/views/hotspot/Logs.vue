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
        <span class="font-semibold text-text-primary">Logs</span>
        <div class="flex items-center gap-3">
          <input
            v-model="search"
            class="input w-64"
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
        v-else-if="filteredLogs.length === 0"
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
              v-for="(entry, i) in pagedLogs"
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
        <div
          v-if="filteredLogs.length > PAGE_SIZE"
          class="flex items-center justify-between pt-3 mt-1 border-t border-border"
        >
          <span class="text-xs text-text-muted">
            {{ (page - 1) * PAGE_SIZE + 1 }}–{{
              Math.min(page * PAGE_SIZE, filteredLogs.length)
            }}
            of {{ filteredLogs.length }}
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
import { ref, computed, watch } from "vue";
import {
  ArrowPathIcon,
  ChevronLeftIcon,
  ChevronRightIcon,
} from "@heroicons/vue/24/outline";
import NoRouterSelected from "@/components/NoRouterSelected.vue";
import EmptyState from "@/components/EmptyState.vue";
import { useRoutersStore } from "@/stores/routers";
import { getSystemLogs } from "@/api";
import PageLayout from "@/components/PageLayout.vue";

const store = useRoutersStore();

interface HotspotLogEntry {
  time: string;
  message: string;
  user: string;
}

const logs = ref<HotspotLogEntry[]>([]);
const loading = ref(false);
const search = ref("");

const filteredLogs = computed(() => {
  const q = search.value.trim().toLowerCase();
  if (!q) return logs.value;
  return logs.value.filter(
    (e) =>
      e.user.toLowerCase().includes(q) ||
      e.message.toLowerCase().includes(q) ||
      e.time.toLowerCase().includes(q),
  );
});

const PAGE_SIZE = 20;
const page = ref(1);
const pageCount = computed(() =>
  Math.max(1, Math.ceil(filteredLogs.value.length / PAGE_SIZE)),
);
const pagedLogs = computed(() => {
  const start = (page.value - 1) * PAGE_SIZE;
  return filteredLogs.value.slice(start, start + PAGE_SIZE);
});

watch(search, () => {
  page.value = 1;
});

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
