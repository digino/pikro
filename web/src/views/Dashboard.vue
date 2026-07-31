<template>
  <PageLayout title="Dashboard" subtitle="Router Overview">
    <NoRouterSelected v-if="!store.activeId" />

    <div v-else class="grid gap-4">
      <!-- ── Row 1: hero ── -->
      <div class="grid gap-4" style="grid-template-columns: 42% 1fr">
        <!-- System card: router image | live stats | health ring -->
        <div
          class="grid rounded-xl border border-border p-5 bg-surface gap-5"
          style="grid-template-columns: auto 1fr auto"
        >

          <!-- LEFT: router illustration + device identity -->
          <div class="flex flex-col items-center gap-2">
            <RouterArt
              :board-name="resource['board-name'] ?? ''"
              :size="120"
              :power-led="healthColor"
              :wifi-led="activeSessions > 0 ? 'var(--color-green)' : 'var(--color-border)'"
              wan-led="var(--color-amber)"
            />
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
            <div v-else class="grid grid-cols-1 text-xs">
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
              <div class="flex justify-between items-center py-1.5 border-b border-border/60 ">
                <span class="text-text-muted">Date</span>
                <span class="font-mono font-semibold text-text-primary">{{ routerDate || "—" }}</span>
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
              <svg width="84" height="84" viewBox="0 0 120 120">
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
                <span class="font-mono text-lg font-bold tracking-tight text-text-primary">{{ healthScore }}</span>
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

      <!-- ── Row 2: sales report + active sessions ── -->
      <div class="grid grid-cols-2 gap-4">
        <!-- Sales report -->
        <div
          class="rounded-xl border border-border p-5 grid gap-3 bg-surface"
          style="grid-template-rows: auto 1fr"
        >
          <div class="flex items-center justify-between">
            <span class="font-semibold text-text-primary">Sales this month</span>
            <RouterLink to="/hotspot/reports" class="btn btn-ghost btn-sm"
              >View reports</RouterLink
            >
          </div>
          <div
            v-if="salesLoading && salesLedger.length === 0"
            class="flex items-center justify-center"
          >
            <span class="spinner spinner--sm" />
          </div>
          <EmptyState v-else-if="salesLedger.length === 0" message="No sales recorded yet" />
          <div v-else class="flex flex-col gap-3">
            <div class="grid grid-cols-2 gap-3">
              <div class="rounded-lg border border-border p-3 bg-base">
                <div class="text-sm font-medium text-text-secondary">Vouchers generated</div>
                <div class="font-mono text-2xl font-bold mt-1 text-text-primary">{{ monthSales.generated }}</div>
              </div>
              <div class="rounded-lg border border-border p-3 bg-base">
                <div class="text-sm font-medium text-text-secondary">Revenue</div>
                <div class="font-mono text-2xl font-bold mt-1">{{ fmtAmount(monthSales.revenue) }}</div>
              </div>
            </div>
            <div class="h-40">
              <SalesBarChart :points="dailySalesPoints" :currency="salesCurrency" />
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
          <EmptyState v-else-if="activeList.length === 0" message="No active sessions" />
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
import EmptyState from "@/components/EmptyState.vue";
import { useRoutersStore } from "@/stores/routers";
import {
  getSystemResource,
  getSystemClock,
  getPollSnapshot,
  listHotspotUsers,
  listHotspotActive,
  getCleanupScheduler,
  getSalesLedger,
  type SaleEntry,
} from "@/api";
import { friendlyError } from "@/utils/errors";
import { formatCompactAmount } from "@/utils/currencies";
import PageLayout from "@/components/PageLayout.vue";
import RouterArt from "@/components/router-art/RouterArt.vue";
import SalesBarChart from "@/components/SalesBarChart.vue";

const store = useRoutersStore();

// ── State ─────────────────────────────────────────────────────
const resource = ref<Record<string, string>>({});
const loading = ref(false);
const error = ref("");
const routerTime = ref("");
const routerDate = ref("");

const allUsers = ref<Record<string, string>[]>([]);
const activeList = ref<Record<string, string>[]>([]);
const cleanupInstalled = ref<boolean | null>(null);
const hotspotLoading = ref(false);
const hotspotError = ref("");

const salesLedger = ref<SaleEntry[]>([]);
const salesLoading = ref(false);

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
    routerDate.value = clock["date"] ?? "";
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

async function loadSales() {
  if (!store.activeId) return;
  salesLoading.value = true;
  try {
    salesLedger.value = await getSalesLedger(
      store.activeId,
      new Date().getFullYear(),
    );
  } catch {
    // non-critical
  } finally {
    salesLoading.value = false;
  }
}

async function poll() {
  if (!store.activeId) return;
  try {
    const snap = await getPollSnapshot(store.activeId);
    resource.value = snap.resource;
    error.value = "";
    if (snap.clock["time"]) {
      routerTime.value = snap.clock["time"];
    }
    if (snap.clock["date"]) {
      routerDate.value = snap.clock["date"];
    }
  } catch {
    // non-critical
  }
}

let pollTimer: ReturnType<typeof setInterval>;
let hotspotTimer: ReturnType<typeof setInterval>;
let salesTimer: ReturnType<typeof setInterval>;

function startTimers() {
  pollTimer = setInterval(() => {
    if (!document.hidden) poll();
  }, 5000);
  hotspotTimer = setInterval(() => {
    if (!document.hidden) loadHotspot();
  }, 15_000);
  salesTimer = setInterval(() => {
    if (!document.hidden) loadSales();
  }, 60_000);
}

function stopTimers() {
  clearInterval(pollTimer);
  clearInterval(hotspotTimer);
  clearInterval(salesTimer);
}

watch(
  () => store.activeId,
  async (id) => {
    stopTimers();
    allUsers.value = [];
    activeList.value = [];
    salesLedger.value = [];
    if (!id) return;
    await Promise.all([loadStatic(), loadHotspot(), loadSales()]);
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

// ── Sales summary (current month) ──────────────────────────────
const monthSalesEntries = computed(() => {
  const now = new Date();
  return salesLedger.value.filter((e) => {
    const d = new Date(e.at);
    return d.getFullYear() === now.getFullYear() && d.getMonth() === now.getMonth();
  });
});

const monthSales = computed(() => {
  let generated = 0;
  let revenue = 0;
  for (const e of monthSalesEntries.value) {
    generated += e.count;
    revenue += (parseFloat(e.price) || 0) * e.count;
  }
  return { generated, revenue };
});

const salesCurrency = computed(
  () => salesLedger.value.find((e) => e.currency)?.currency ?? "",
);

// By-day breakdown for the current month, for the bar chart.
const dailySalesPoints = computed(() => {
  const bucket = new Map<string, { count: number; revenue: number }>();
  for (const e of monthSalesEntries.value) {
    const d = new Date(e.at);
    const key = String(d.getDate());
    const entry = bucket.get(key) ?? { count: 0, revenue: 0 };
    entry.count += e.count;
    entry.revenue += (parseFloat(e.price) || 0) * e.count;
    bucket.set(key, entry);
  }
  return [...bucket.entries()]
    .sort(([a], [b]) => parseInt(a) - parseInt(b))
    .map(([label, v]) => ({ label, count: v.count, revenue: v.revenue }));
});

function fmtAmount(n: number): string {
  return formatCompactAmount(n, salesCurrency.value);
}

// ── Formatters ────────────────────────────────────────────────
function formatBytes(n: number): string {
  if (!n || isNaN(n)) return "—";
  if (n >= 1_073_741_824) return (n / 1_073_741_824).toFixed(1) + " GB";
  if (n >= 1_048_576) return (n / 1_048_576).toFixed(1) + " MB";
  if (n >= 1_024) return (n / 1_024).toFixed(0) + " KB";
  return n + " B";
}
</script>
