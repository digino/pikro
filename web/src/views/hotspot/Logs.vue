<template>
  <div v-if="!store.activeId" class="p-6">
    <NoRouterSelected />
  </div>

  <main v-else class="flex flex-col gap-6 px-4 sm:px-6 lg:px-10 pt-7 pb-14">
    <header class="flex flex-wrap items-end justify-between gap-4">
      <div class="flex items-center gap-4">
        <img :src="inkIcon('logs')" alt="" width="50" height="50" />
        <div>
          <nav :aria-label="t('logs.breadcrumb')" class="eyebrow">
            {{ t('logs.breadcrumb') }} / {{ t('logs.title') }}
          </nav>
          <h1 class="type-display text-3xl">{{ t('logs.title') }}</h1>
        </div>
      </div>
      <div class="flex flex-wrap gap-3">
        <button
          type="button"
          class="btn"
          :class="{ 'btn-ink shadow-[3px_3px_0_var(--color-primary)]': live }"
          :aria-pressed="live"
          @click="live = !live"
        >
          <span
            class="size-2.5 rounded-full"
            :class="live ? 'bg-primary motion-safe:animate-pulse' : 'bg-graphite'"
          />
          {{ live ? t('logs.live') : t('logs.paused') }}
        </button>
        <button type="button" class="btn" :disabled="filtered.length === 0" @click="exportCsv">
          <ArrowDownTrayIcon class="size-4" />
          {{ t('logs.exportCsv') }}
        </button>
      </div>
    </header>

    <section :aria-label="t('logs.statsCaption')" class="flex flex-col gap-2.5">
      <span class="eyebrow">{{ t('logs.statsCaption') }}</span>
      <div class="grid grid-cols-2 xl:grid-cols-4 gap-4.5">
        <div
          v-for="s in stats"
          :key="s.event"
          class="card px-4.5 py-4 flex flex-col gap-1.5"
          :class="s.card"
        >
          <span class="flex items-center gap-2 text-sm font-semibold">
            <span class="shrink-0" :class="EVENT_STYLE[s.event].marker" />
            {{ s.label }}
          </span>
          <span class="type-display text-4xl leading-none">{{ counts[s.event] }}</span>
        </div>
      </div>
    </section>

    <div class="flex flex-wrap items-center justify-between gap-3">
      <div role="group" :aria-label="t('logs.filterLabel')" class="flex flex-wrap gap-2">
        <button
          v-for="f in filters"
          :key="f.value"
          type="button"
          :aria-pressed="eventFilter === f.value"
          class="inline-flex items-center gap-2 min-h-11 px-4 rounded-full border-2 border-ink text-sm font-semibold cursor-pointer transition-colors"
          :class="
            eventFilter === f.value
              ? 'bg-ink text-paper shadow-[3px_3px_0_var(--color-primary)]'
              : 'bg-paper text-ink hover:bg-primary-wash'
          "
          @click="eventFilter = f.value"
        >
          {{ f.label }}
          <span class="font-mono text-xs opacity-75">{{ f.count }}</span>
        </button>
      </div>
      <label
        class="flex items-center gap-2 min-h-11 min-w-65 px-3.5 rounded-control border-2 border-ink bg-paper focus-within:outline-3 focus-within:outline-primary focus-within:outline-offset-1"
      >
        <MagnifyingGlassIcon class="size-4.5 shrink-0" />
        <span class="sr-only">{{ t('logs.search') }}</span>
        <input
          v-model="search"
          type="search"
          :placeholder="t('logs.searchPlaceholder')"
          class="flex-1 min-w-0 border-0 bg-transparent text-sm outline-none"
        />
      </label>
    </div>

    <div
      v-if="error"
      role="alert"
      class="flex items-center gap-2.5 p-4 rounded-control border-2 border-danger-text bg-danger-wash text-danger-text text-sm"
    >
      <ExclamationTriangleIcon class="size-5 shrink-0" />
      {{ error }}
      <button type="button" class="btn btn-sm ml-auto" @click="load">{{ t('logs.retry') }}</button>
    </div>

    <section class="card overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full min-w-225 text-sm">
          <thead>
            <tr class="border-b-2 border-ink">
              <th class="eyebrow text-left px-4 py-3.5 w-24">{{ t('logs.colTime') }}</th>
              <th class="eyebrow text-left px-4 py-3.5 w-33">{{ t('logs.colEvent') }}</th>
              <th class="eyebrow text-left px-4 py-3.5">{{ t('logs.colUser') }}</th>
              <th class="eyebrow text-left px-4 py-3.5">{{ t('logs.colAddress') }}</th>
              <th class="eyebrow text-left px-4 py-3.5">{{ t('logs.colMessage') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="loading && logs.length === 0">
              <td colspan="5" class="py-10 text-center"><span class="spinner" /></td>
            </tr>
            <tr v-else-if="visible.length === 0">
              <td colspan="5" class="px-4 py-10 text-center text-text-muted">
                {{ logs.length === 0 ? t('logs.empty') : t('logs.noMatch') }}
              </td>
            </tr>
            <template v-for="row in rows" :key="row.key">
              <tr v-if="row.day !== undefined">
                <td
                  colspan="5"
                  class="px-4 py-2 bg-ink text-paper text-[11px] font-bold uppercase tracking-widest"
                >
                  {{ row.day || t('logs.today') }}
                </td>
              </tr>
              <tr v-else class="border-b border-hairline last:border-0 hover:bg-primary-wash">
                <td class="px-4 py-3 font-mono text-[13px] text-text-muted">{{ row.entry.clock }}</td>
                <td class="px-4 py-3">
                  <span class="chip" :class="EVENT_STYLE[row.entry.event].chip">
                    <span class="shrink-0" :class="EVENT_STYLE[row.entry.event].marker" />
                    {{ t(`logs.event.${row.entry.event}`) }}
                  </span>
                </td>
                <td class="px-4 py-3 font-mono font-semibold">{{ row.entry.user || '—' }}</td>
                <td class="px-4 py-3 font-mono text-[13px] text-text-muted">{{ row.entry.address || '—' }}</td>
                <td class="px-4 py-3" :class="EVENT_STYLE[row.entry.event].message">{{ row.entry.message }}</td>
              </tr>
            </template>
          </tbody>
        </table>
      </div>
      <div v-if="filtered.length > visibleCount" class="flex justify-center p-3.5 border-t-2 border-ink">
        <button type="button" class="btn btn-sm" @click="visibleCount += PAGE">
          {{ t('logs.loadOlder') }}
        </button>
      </div>
    </section>
  </main>
</template>

<script setup lang="ts">
import { ref, computed, watch, onUnmounted } from "vue";
import { useI18n } from "vue-i18n";
import {
  ArrowDownTrayIcon,
  ExclamationTriangleIcon,
  MagnifyingGlassIcon,
} from "@heroicons/vue/24/outline";
import NoRouterSelected from "@/components/NoRouterSelected.vue";
import { useRoutersStore } from "@/stores/routers";
import { getSystemLogs } from "@/api";
import { friendlyError } from "@/utils/errors";
import { downloadCsv } from "@/utils/csv";
import { inkIcon } from "@/utils/inkIcons";

const store = useRoutersStore();
const { t } = useI18n();

type LogEvent = "login" | "logout" | "failed" | "cleanup" | "system";

interface HotspotLogEntry {
  day: string;
  clock: string;
  event: LogEvent;
  user: string;
  address: string;
  message: string;
}

const EVENTS: LogEvent[] = ["login", "logout", "failed", "cleanup", "system"];

// Events differ by marker shape and fill, not colour alone.
const EVENT_STYLE: Record<LogEvent, { chip: string; marker: string; message: string }> = {
  login: { chip: "", marker: "size-2 rounded-full bg-ink", message: "" },
  logout: { chip: "", marker: "size-2 rounded-full border-2 border-ink", message: "" },
  failed: {
    chip: "bg-danger-wash text-danger-text border-danger-text",
    marker: "size-2 rotate-45 bg-danger",
    message: "text-danger-text font-semibold",
  },
  cleanup: {
    chip: "bg-primary text-primary-fg",
    marker: "size-2 rounded-[2px] bg-primary-fg",
    message: "",
  },
  system: {
    chip: "bg-hairline text-text-secondary border-dashed",
    marker: "w-2 h-0.5 bg-graphite-deep",
    message: "text-text-muted",
  },
};

const PAGE = 100;
const LIVE_INTERVAL_MS = 10_000;

const logs = ref<HotspotLogEntry[]>([]);
const loading = ref(false);
const error = ref("");
const search = ref("");
const eventFilter = ref<LogEvent | "">("");
const visibleCount = ref(PAGE);
const live = ref(true);

// Hotspot log messages look like "<user> (<ip>): <action>", e.g.
// "gbiu (192.168.88.238): logged in" or "a3fm (192.168.88.238): login failed: invalid username or password".
const USER_MESSAGE_RE = /^(\S+)\s+\(([^)]+)\):\s*(.*)$/;
// Pikro's own scripts: "pikro-cleanup: removing <user> (<reason>)", "pikro-on-login: stamping <user> (...)".
const PIKRO_SCRIPT_RE = /^pikro-(cleanup|on-login):\s*(?:removing|stamping)\s+(\S+)\s*(.*)$/;
// RouterOS prefixes the clock with a date for entries not from today ("oct/24 21:48:12", "2026-10-24 21:48:12").
const TIME_RE = /^(.*?)\s*(\d{1,2}:\d{2}:\d{2})$/;

function classify(action: string): LogEvent {
  if (action.startsWith("logged in")) return "login";
  if (action.startsWith("logged out")) return "logout";
  if (action.startsWith("login failed")) return "failed";
  return "system";
}

function parseEntry(e: Record<string, string>): HotspotLogEntry {
  const message = e.message ?? "";
  const time = (e.time ?? "").match(TIME_RE);
  const base = { day: time?.[1] ?? "", clock: time?.[2] ?? e.time ?? "" };

  const script = message.match(PIKRO_SCRIPT_RE);
  if (script) {
    return {
      ...base,
      event: script[1] === "cleanup" ? "cleanup" : "system",
      user: script[2],
      address: "",
      message: script[3].replace(/^\((.*)\)$/, "$1"),
    };
  }
  const match = message.match(USER_MESSAGE_RE);
  return {
    ...base,
    event: match ? classify(match[3]) : "system",
    user: match?.[1] ?? "",
    address: match?.[2] ?? "",
    message: match ? match[3] : message,
  };
}

const counts = computed(() => {
  const c: Record<LogEvent, number> = { login: 0, logout: 0, failed: 0, cleanup: 0, system: 0 };
  for (const e of logs.value) c[e.event]++;
  return c;
});

const stats = computed(() => [
  { event: "login" as const, label: t("logs.statLogins"), card: "" },
  { event: "logout" as const, label: t("logs.statLogouts"), card: "" },
  { event: "failed" as const, label: t("logs.statFailed"), card: "bg-danger-wash text-danger-text" },
  { event: "cleanup" as const, label: t("logs.statCleanup"), card: "bg-primary text-primary-fg" },
]);

const filters = computed(() => [
  { value: "" as const, label: t("logs.all"), count: logs.value.length },
  ...EVENTS.map((e) => ({ value: e, label: t(`logs.filter.${e}`), count: counts.value[e] })),
]);

const filtered = computed(() => {
  const q = search.value.trim().toLowerCase();
  return logs.value.filter(
    (e) =>
      (!eventFilter.value || e.event === eventFilter.value) &&
      (!q || `${e.user} ${e.address} ${e.message}`.toLowerCase().includes(q)),
  );
});

const visible = computed(() => filtered.value.slice(0, visibleCount.value));

// Entries are newest first; insert an ink separator row whenever the day changes.
const rows = computed(() => {
  const out: ({ key: string; day: string } | { key: string; day?: undefined; entry: HotspotLogEntry })[] = [];
  let lastDay: string | null = null;
  visible.value.forEach((entry, i) => {
    if (entry.day !== lastDay) {
      out.push({ key: `day-${i}`, day: entry.day });
      lastDay = entry.day;
    }
    out.push({ key: `row-${i}`, entry });
  });
  return out;
});

watch([search, eventFilter], () => {
  visibleCount.value = PAGE;
});

async function load() {
  if (!store.activeId || loading.value) return;
  loading.value = true;
  error.value = "";
  try {
    const all = await getSystemLogs(store.activeId);
    logs.value = all.map(parseEntry);
  } catch (e: any) {
    error.value = friendlyError(e, t("logs.failedToLoad"));
  } finally {
    loading.value = false;
  }
}

function exportCsv() {
  downloadCsv(
    "pikro-hotspot-logs.csv",
    [t("logs.colTime"), t("logs.colEvent"), t("logs.colUser"), t("logs.colAddress"), t("logs.colMessage")],
    filtered.value.map((e) => [
      [e.day, e.clock].filter(Boolean).join(" "),
      t(`logs.event.${e.event}`),
      e.user,
      e.address,
      e.message,
    ]),
  );
}

let liveTimer: ReturnType<typeof setInterval> | undefined;

watch(
  [() => store.activeId, live],
  ([id, isLive], [prevId]) => {
    clearInterval(liveTimer);
    if (!id) return;
    if (id !== prevId) {
      logs.value = [];
      visibleCount.value = PAGE;
    }
    load();
    if (isLive) {
      liveTimer = setInterval(() => {
        if (!document.hidden) load();
      }, LIVE_INTERVAL_MS);
    }
  },
  { immediate: true },
);

onUnmounted(() => clearInterval(liveTimer));
</script>
