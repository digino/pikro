<template>
  <PageLayout title="Dashboard" subtitle="Router Overview">
    <NoRouterSelected v-if="!store.activeId" />

    <div v-else class="grid gap-4">
      <!-- ── Row 1: hero ── -->
      <div class="grid gap-4" style="grid-template-columns: 42% 1fr">
        <!-- System card: router image | live stats | health ring -->
        <div class="grid grid-cols-3 rounded-xl border border-border p-5 bg-surface gap-5">

          <!-- LEFT: router illustration + device identity -->
          <div class="flex flex-col items-center gap-2">
            <RouterArt
              :board-name="resource['board-name'] ?? ''"
              :size="120"
              :power-led="healthColor"
              :wifi-led="activeSessions > 0 ? 'var(--color-green)' : 'var(--color-border)'"
              wan-led="var(--color-amber)"
            />
            <!-- Status badge -->
            <div class="flex items-center gap-1.5 px-2 py-0.5 rounded-full border border-border text-xs font-medium text-text-secondary">
              <span class="size-1.5 rounded-full shrink-0" :style="`background:${healthColor}; box-shadow:0 0 5px ${healthColor}`"/>
              {{ healthLabel }}
            </div>
            <!-- Device identity -->
            <div class="text-center space-y-0.5">
              <div class="text-sm font-semibold text-text-primary font-mono leading-tight">
                {{ resource['board-name'] || store.active()?.name || '—' }}
              </div>
              <div class="text-xs text-text-muted font-mono leading-tight">
                RouterOS {{ resource['version']?.split(' ')[0] || '—' }}
              </div>
            </div>
          </div>

          <!-- MIDDLE: live stats -->
          <div class="flex flex-col justify-center min-w-0">
            <div v-if="loading" class="flex items-center gap-2 text-sm text-text-muted">
              <span class="spinner spinner--sm" /> Loading…
            </div>
            <div v-else-if="error" class="flex items-center gap-1.5 text-xs text-red">
              <ExclamationTriangleIcon class="size-3.5 shrink-0" />{{ error }}
            </div>
            <div v-else class="grid grid-cols-1 text-sm">
              <div class="flex justify-between items-center py-1.5 border-b border-border/60 ">
                <span class="text-text-muted">Uptime</span>
                <span class="font-mono font-semibold text-text-primary">{{ resource["uptime"] ?? "—" }}</span>
              </div>
              <div class="flex justify-between items-center py-1.5 border-b border-border/60 ">
                <span class="text-text-muted">CPU</span>
                <span class="font-mono font-semibold text-text-primary">{{ resource["cpu-load"] ?? "—" }}%</span>
              </div>
              <div class="flex justify-between items-center py-1.5 border-b border-border/60 ">
                <span class="text-text-muted">Free RAM</span>
                <span class="font-mono font-semibold text-text-primary">{{ formatBytes(freeMemory) }}</span>
              </div>
              <div class="flex justify-between items-center py-1.5 border-b border-border/60 ">
                <span class="text-text-muted">Free disk</span>
                <span class="font-mono font-semibold text-text-primary">{{ formatBytes(parseInt(resource["free-hdd-space"] ?? "0") || 0) }}</span>
              </div>
              <div class="flex justify-between items-center py-1.5 ">
                <span class="text-text-muted">Time</span>
                <span class="font-mono font-semibold text-text-primary">{{ routerTime || "—" }}</span>
              </div>
            </div>
          </div>

          <!-- RIGHT: health ring -->
          <div class="flex flex-col items-center justify-center gap-2 shrink-0">
            <div class="relative">
              <svg width="110" height="110" viewBox="0 0 120 120">
                <circle cx="60" cy="60" r="50" fill="none" stroke="var(--color-border)" stroke-width="8"/>
                <circle
                  cx="60" cy="60" r="50" fill="none"
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
                <span class="font-mono text-2xl font-bold tracking-tight text-text-primary">{{ healthScore }}</span>
                <span class="text-xs font-bold text-text-secondary">Health</span>
              </div>
            </div>
          </div>

        </div>

        <!-- Right column: hotspot card -->
        <div>
          <!-- Hotspot card -->
          <div
            class="rounded-xl border border-border p-4 grid gap-3 bg-surface"
            style="grid-template-rows: auto 1fr auto"
          >
            <div class="flex items-center justify-between">
              <span class="font-semibold text-text-primary">Hotspot</span>
              <RouterLink to="/hotspot/users" class="btn btn-primary"
                >Manage</RouterLink
              >
            </div>
            <div
              v-if="hotspotLoading && totalUsers === 0"
              class="flex items-center justify-center"
            >
              <span class="spinner spinner--sm" />
            </div>
            <template v-else>
              <div
                class="grid grid-cols-4 gap-3 tracking-tight text-sm font-semibold text-secondary"
              >
                <div class="rounded-lg border border-border p-3 bg-base">
                  <div class="text-text-secondary">Users</div>
                  <div class="font-mono text-2xl font-bold mt-1 text-text-primary">{{ totalUsers }}</div>
                </div>
                <div class="rounded-lg border border-border p-3 bg-base">
                  <div class="text-text-secondary">Active</div>
                  <div class="font-mono text-2xl font-bold mt-1 text-green">{{ activeSessions }}</div>
                </div>
                <div class="rounded-lg border border-border p-3 bg-base">
                  <div class="text-text-secondary">Expired</div>
                  <div class="font-mono text-2xl font-bold mt-1 text-amber">{{ expiredUsers }}</div>
                </div>
                <div class="rounded-lg border border-border p-3 bg-base">
                  <div class="text-text-secondary">Disabled</div>
                  <div class="font-mono text-2xl font-bold mt-1 text-text-secondary">{{ disabledUsers }}</div>
                </div>
              </div>
              <div
                v-if="cleanupInstalled !== null"
                class="flex items-center gap-2.5 p-3 rounded-lg border border-border bg-base"
              >
                <component
                  :is="cleanupInstalled ? CheckCircleIcon : ExclamationTriangleIcon"
                  class="size-4 shrink-0"
                  :class="cleanupInstalled ? 'text-green' : 'text-amber'"
                />
                <div class="flex-1 min-w-0">
                  <div class="text-sm tracking-tight text-text-primary">
                    Auto-cleanup
                  </div>
                  <div class="text-xs text-text-secondary mt-1">
                    {{
                      cleanupInstalled
                        ? "Scheduler active"
                        : "Not configured — expired users accumulate"
                    }}
                  </div>
                </div>
                <RouterLink
                  v-if="!cleanupInstalled"
                  to="/hotspot/settings"
                  class="text-amber btn btn-ghost"
                >
                  Configure
                </RouterLink>
              </div>
            </template>
          </div>
        </div>
      </div>

      <!-- ── Row 2: bandwidth chart + interface traffic ── -->
      <div class="grid grid-cols-2 gap-4">
        <!-- Bandwidth chart -->
        <div class="rounded-xl border border-border p-5 bg-surface h-full">
          <div class="flex items-center gap-4 mb-4">
            <span class="font-semibold text-text-primary">Bandwidth</span>
            <div class="flex items-center gap-1.5">
              <span class="size-2 rounded-sm shrink-0" style="background:#22d3ee"/>
              <span class="text-xs text-text-secondary">Download</span>
              <span class="font-mono text-xs font-semibold text-text-primary">{{ curDown }} Mbps</span>
            </div>
            <div class="flex items-center gap-1.5">
              <span class="size-2 rounded-sm shrink-0" style="background:#f59e0b"/>
              <span class="text-xs text-text-secondary">Upload</span>
              <span class="font-mono text-xs font-semibold text-text-primary">{{ curUp }} Mbps</span>
            </div>
          </div>
          <div class="h-60">
            <BandwidthChart :history="bwHistory" />
          </div>
        </div>

        <!-- Interface traffic -->
        <div class="rounded-xl border border-border p-5 bg-surface">
          <div class="flex items-center justify-between mb-4">
            <span class="font-semibold text-text-primary">Interfaces</span>
            <span v-if="trafficUpdatedAt" class="text-xs text-text-muted">{{
              trafficUpdatedAt
            }}</span>
          </div>
          <div
            v-if="trafficLoading && traffic.length === 0"
            class="flex justify-center py-8"
          >
            <span class="spinner spinner--sm" />
          </div>
          <div
            v-else-if="traffic.length === 0"
            class="text-xs text-text-muted py-8 text-center"
          >
            No interface data
          </div>
          <table v-else class="w-full text-xs">
            <thead>
              <tr class="text-text-muted border-b border-border">
                <th class="text-left pb-2 font-medium">Interface</th>
                <th class="text-left pb-2 font-medium">IP</th>
                <th class="text-right pb-2 font-medium">RX</th>
                <th class="text-right pb-2 font-medium">TX</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="iface in traffic"
                :key="iface.name"
                class="border-b border-border/40 last:border-0"
              >
                <td class="py-1.5 font-mono text-text-primary">{{ iface.name }}</td>
                <td class="py-1.5 font-mono text-text-muted">{{ ifaceIP(iface.name) }}</td>
                <td class="py-1.5 text-right font-mono text-text-primary">{{ formatBps(iface["rx-bits-per-second"]) }}</td>
                <td class="py-1.5 text-right font-mono text-text-secondary">{{ formatBps(iface["tx-bits-per-second"]) }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- ── Row 3: recent logs + active sessions ── -->
      <div class="grid grid-cols-2 gap-4">
        <!-- Recent logs -->
        <div
          class="rounded-xl border border-border p-5 grid gap-3 bg-surface"
          style="grid-template-rows: auto 1fr"
        >
          <div class="flex items-center justify-between">
            <span class="font-semibold text-text-primary">Recent logs</span>
            <RouterLink to="/hotspot/logs" class="btn btn-ghost btn-sm"
              >View all</RouterLink
            >
          </div>
          <div
            v-if="logsLoading && logs.length === 0"
            class="flex justify-center py-4"
          >
            <span class="spinner spinner--sm" />
          </div>
          <div
            v-else-if="logs.length === 0"
            class="text-xs text-text-muted py-2 text-center"
          >
            No log entries
          </div>
          <div v-else class="grid" style="align-content: start">
            <div
              v-for="(entry, i) in logs"
              :key="i"
              class="flex items-start gap-2 py-1.5 border-b border-border/40 last:border-0"
            >
              <span
                class="shrink-0 mt-1 size-1.5 rounded-full"
                :style="`background:${logColor(entry.topics)}`"
              />
              <div class="min-w-0 flex-1">
                <p class="text-xs text-text-primary truncate">
                  {{ entry.message }}
                </p>
                <p class="text-xs text-text-muted font-mono">
                  {{ entry.time }}
                </p>
              </div>
            </div>
          </div>
        </div>

        <!-- Active sessions -->
        <div
          class="rounded-xl border border-border p-5 grid gap-3 bg-surface"
          style="grid-template-rows: auto 1fr"
        >
          <div class="flex items-center justify-between">
            <span class="font-semibold text-text-primary">Active sessions</span>
            <RouterLink to="/hotspot/users" class="btn btn-ghost btn-sm"
              >View all</RouterLink
            >
          </div>
          <div
            v-if="hotspotLoading && activeList.length === 0"
            class="flex justify-center py-6"
          >
            <span class="spinner spinner--sm" />
          </div>
          <div
            v-else-if="activeList.length === 0"
            class="text-xs text-text-muted py-4 text-center"
          >
            No active sessions
          </div>
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
                  v-for="s in activeList.slice(0, 8)"
                  :key="s['.id']"
                  class="border-b border-border/40 last:border-0"
                >
                  <td class="py-2 font-mono font-medium text-text-primary">
                    {{ s.user }}
                  </td>
                  <td class="py-2 font-mono text-text-secondary">
                    {{ s.address }}
                  </td>
                  <td class="py-2 text-right font-mono text-text-secondary">
                    {{ s.uptime || "—" }}
                  </td>
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
        </div>
      </div>
    </div>
  </PageLayout>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted } from "vue";
import { RouterLink } from "vue-router";
import {
  ExclamationTriangleIcon,
  CheckCircleIcon,
} from "@heroicons/vue/24/outline";
import NoRouterSelected from "@/components/NoRouterSelected.vue";
import { useRoutersStore } from "@/stores/routers";
import {
  getSystemResource,
  getSystemClock,
  getSystemLogs,
  getPollSnapshot,
  listHotspotUsers,
  listHotspotActive,
  getCleanupScheduler,
} from "@/api";
import { friendlyError } from "@/utils/errors";
import PageLayout from "@/components/PageLayout.vue";
import BandwidthChart from "@/components/BandwidthChart.vue";
import RouterArt from "@/components/router-art/RouterArt.vue";

const store = useRoutersStore();

// ── State ─────────────────────────────────────────────────────
const resource = ref<Record<string, string>>({});
const loading = ref(false);
const error = ref("");
const routerTime = ref("");

const allUsers = ref<Record<string, string>[]>([]);
const activeList = ref<Record<string, string>[]>([]);
const cleanupInstalled = ref<boolean | null>(null);
const hotspotLoading = ref(false);
const hotspotError = ref("");

const traffic = ref<Record<string, string>[]>([]);
const addresses = ref<Record<string, string>[]>([]);
const trafficLoading = ref(false);
const trafficUpdatedAt = ref("");

const logs = ref<Record<string, string>[]>([]);
const logsLoading = ref(false);

// ── Bandwidth history ─────────────────────────────────────────
const N = 46;
interface BwPoint {
  down: number;
  up: number;
}
const bwHistory = ref<BwPoint[]>(
  Array.from({ length: N }, () => ({ down: 0, up: 0 })),
);

// ── Hotspot computed counts ───────────────────────────────────
const totalUsers = computed(() => allUsers.value.length);
const activeSessions = computed(() => activeList.value.length);
const expiredUsers = computed(() => {
  const n = Math.floor(Date.now() / 1000);
  return allUsers.value.filter((u) => {
    const c = u.comment ?? "";
    if (!c.startsWith("exp:")) return false;
    const exp = c.slice(4);
    const epoch = parseInt(exp);
    if (!isNaN(epoch) && epoch > 1_000_000_000) return epoch < n;
    if (exp.includes("-")) return new Date(exp.replace(" ", "T")) < new Date();
    return false;
  }).length;
});
const disabledUsers = computed(
  () => allUsers.value.filter((u) => u.disabled === "true").length,
);

// ── Health computeds ──────────────────────────────────────────
const cpuLoad = computed(
  () => parseInt(resource.value["cpu-load"] ?? "0") || 0,
);
const freeMemory = computed(
  () => parseInt(resource.value["free-memory"] ?? "0") || 0,
);
const totalMemory = computed(
  () => parseInt(resource.value["total-memory"] ?? "1") || 1,
);

const healthScore = computed(() => {
  if (!store.activeId || error.value) return "—";
  const cpuScore = Math.round((1 - cpuLoad.value / 100) * 40);
  const ramScore = Math.round(
    Math.min(freeMemory.value / totalMemory.value / 0.5, 1) * 30,
  );
  return String(Math.min(100, cpuScore + ramScore + 30));
});
const healthColor = computed(() => {
  if (error.value) return "var(--color-red)";
  const s = parseInt(healthScore.value) || 0;
  if (s >= 75) return "var(--color-green)";
  if (s >= 50) return "var(--color-amber)";
  return "var(--color-red)";
});
const healthLabel = computed(() => {
  if (error.value) return "Unreachable";
  const s = parseInt(healthScore.value) || 0;
  if (s >= 75) return "Router reachable";
  if (s >= 50) return "Router degraded";
  return "Router under stress";
});
const ringCirc = (2 * Math.PI * 50).toFixed(1);
const ringOffset = computed(() => {
  const s = parseInt(healthScore.value) || 0;
  return (2 * Math.PI * 50 * (1 - s / 100)).toFixed(1);
});

// ── Chart ─────────────────────────────────────────────────────
const curDown = computed(() => bwHistory.value[N - 1].down.toFixed(2));
const curUp = computed(() => bwHistory.value[N - 1].up.toFixed(2));

// ── Data loading ──────────────────────────────────────────────
async function loadStatic() {
  if (!store.activeId) return;
  loading.value = true;
  error.value = "";
  try {
    const [res, clock] = await Promise.all([
      getSystemResource(store.activeId),
      getSystemClock(store.activeId).catch(
        () => ({}) as Record<string, string>,
      ),
    ]);
    resource.value = res;
    routerTime.value = clock["time"] ?? "";
  } catch (e: any) {
    error.value = friendlyError(e, "Could not reach router");
  } finally {
    loading.value = false;
  }
}

async function loadHotspot() {
  if (!store.activeId) return;
  hotspotLoading.value = true;
  hotspotError.value = "";
  try {
    const [users, active, cleanup] = await Promise.all([
      listHotspotUsers(store.activeId),
      listHotspotActive(store.activeId),
      getCleanupScheduler(store.activeId).catch(() => null),
    ]);
    allUsers.value = users;
    activeList.value = active;
    cleanupInstalled.value = cleanup?.installed ?? null;
  } catch (e: any) {
    hotspotError.value = friendlyError(e, "Could not load hotspot data");
  } finally {
    hotspotLoading.value = false;
  }
}

async function loadLogs() {
  if (!store.activeId) return;
  logsLoading.value = true;
  try {
    const all = await getSystemLogs(store.activeId);
    logs.value = all
      .filter((e) => (e.topics ?? "").toLowerCase().includes("hotspot"))
      .slice(0, 8);
  } catch {
    // non-critical
  } finally {
    logsLoading.value = false;
  }
}

async function poll() {
  if (!store.activeId) return;
  trafficLoading.value = true;
  try {
    const snap = await getPollSnapshot(store.activeId);
    resource.value = snap.resource;
    error.value = "";
    if (snap.traffic.length > 0) {
      traffic.value = snap.traffic.filter(
        (i: Record<string, string>) => i.name,
      );
      trafficUpdatedAt.value = new Date().toLocaleTimeString();
      const totalRx = snap.traffic.reduce(
        (s: number, i: Record<string, string>) =>
          s + (parseInt(i["rx-bits-per-second"] ?? "0") || 0),
        0,
      );
      const totalTx = snap.traffic.reduce(
        (s: number, i: Record<string, string>) =>
          s + (parseInt(i["tx-bits-per-second"] ?? "0") || 0),
        0,
      );
      bwHistory.value = [
        ...bwHistory.value.slice(1),
        { down: totalRx / 1_000_000, up: totalTx / 1_000_000 },
      ];
    }
    if (snap.addresses) {
      addresses.value = snap.addresses;
    }
    if (snap.clock["time"]) {
      routerTime.value = snap.clock["time"];
    }
  } catch {
    // non-critical
  } finally {
    trafficLoading.value = false;
  }
}


let pollTimer: ReturnType<typeof setInterval>;
let hotspotTimer: ReturnType<typeof setInterval>;
let logTimer: ReturnType<typeof setInterval>;

function startTimers() {
  pollTimer = setInterval(() => {
    if (!document.hidden) poll();
  }, 3000);
  hotspotTimer = setInterval(() => {
    if (!document.hidden) loadHotspot();
  }, 15_000);
  logTimer = setInterval(() => {
    if (!document.hidden) loadLogs();
  }, 30_000);
}

function stopTimers() {
  clearInterval(pollTimer);
  clearInterval(hotspotTimer);
  clearInterval(logTimer);
}

watch(
  () => store.activeId,
  async (id) => {
    stopTimers();
    bwHistory.value = Array.from({ length: N }, () => ({ down: 0, up: 0 }));
    traffic.value = [];
    addresses.value = [];
    logs.value = [];
    allUsers.value = [];
    activeList.value = [];
    if (!id) return;
    await Promise.all([loadStatic(), loadHotspot(), loadLogs()]);
    await poll();
    startTimers();
  },
  { immediate: true },
);

onMounted(() => {
  document.addEventListener("visibilitychange", () => {
    if (!document.hidden && store.activeId) poll();
  });
});

onUnmounted(stopTimers);

// ── Formatters ────────────────────────────────────────────────
function ifaceIP(name: string): string {
  const match = addresses.value.find(a => a.interface === name)
  return match ? match.address : '—'
}

function formatBytes(n: number): string {
  if (!n || isNaN(n)) return "—";
  if (n >= 1_073_741_824) return (n / 1_073_741_824).toFixed(1) + " GB";
  if (n >= 1_048_576) return (n / 1_048_576).toFixed(1) + " MB";
  if (n >= 1_024) return (n / 1_024).toFixed(0) + " KB";
  return n + " B";
}

function formatBps(val: string | undefined): string {
  const n = parseInt(val ?? "0") || 0;
  if (n >= 1_000_000_000) return (n / 1_000_000_000).toFixed(2) + " Gbps";
  if (n >= 1_000_000) return (n / 1_000_000).toFixed(2) + " Mbps";
  if (n >= 1_000) return (n / 1_000).toFixed(1) + " Kbps";
  return n + " bps";
}

function logColor(topics: string | undefined): string {
  if (!topics) return "var(--color-text-muted)";
  const t = topics.toLowerCase();
  if (t.includes("error") || t.includes("critical")) return "var(--color-red)";
  if (t.includes("warning")) return "var(--color-amber)";
  if (t.includes("info")) return "var(--color-green)";
  return "var(--color-text-muted)";
}
</script>
