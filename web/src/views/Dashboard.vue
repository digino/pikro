<template>
  <PageLayout title="Dashboard" subtitle="Network Overview">
    <template #actions>
      <div
        class="flex items-center gap-1 rounded-lg border border-border"
        style="background: rgba(255, 255, 255, 0.02)"
      >
        <button
          v-for="r in RANGES"
          :key="r"
          class="btn btn-sm"
          :class="range === r ? 'btn-primary' : 'btn-ghost border-transparent'"
          @click="range = r"
        >
          {{ r.toUpperCase() }}
        </button>
      </div>
      <button class="btn btn-ghost btn-sm" @click="refreshBw">
        <ArrowPathIcon class="size-3.5" :class="{ 'animate-spin': spinning }" />
        Refresh
      </button>
    </template>

    <!-- No router selected -->
    <div
      v-if="!store.activeId"
      class="flex flex-col items-center justify-center flex-1 gap-3"
    >
      <ServerIcon class="size-10 text-text-muted" />
      <p class="text-sm text-text-secondary">
        Select a router to view the dashboard.
      </p>
    </div>

    <div v-else class="flex flex-col gap-4">
      <!-- ── Hero row ── -->
      <div class="flex gap-4">
        <!-- Health score -->
        <div
          class="flex items-center gap-6 flex-[1.5] rounded-xl border border-border p-6"
          style="background: linear-gradient(135deg, #0e0f11, #0b0c0d)"
        >
          <!-- Ring -->
          <div class="relative shrink-0">
            <svg width="112" height="112" viewBox="0 0 120 120">
              <circle
                cx="60"
                cy="60"
                r="50"
                fill="none"
                stroke="rgba(255,255,255,0.07)"
                stroke-width="9"
              />
              <circle
                cx="60"
                cy="60"
                r="50"
                fill="none"
                :stroke="healthColor"
                stroke-width="9"
                stroke-linecap="round"
                :stroke-dasharray="ringCirc"
                :stroke-dashoffset="ringOffset"
                transform="rotate(-90 60 60)"
                style="transition: stroke-dashoffset 0.6s ease"
              />
            </svg>
            <div
              class="absolute inset-0 flex flex-col items-center justify-center"
            >
              <span
                class="font-mono text-2xl font-bold tracking-tight text-text-primary"
                >{{ healthScore }}</span
              >
              <span class="text-xs uppercase tracking-widest text-text-muted"
                >Health</span
              >
            </div>
          </div>

          <!-- Status text -->
          <div class="flex-1 min-w-0">
            <div class="flex items-center gap-2">
              <span
                class="size-2 rounded-full shrink-0"
                :style="`background: ${healthColor}; box-shadow: 0 0 8px ${healthColor}`"
              ></span>
              <span class="font-semibold tracking-tight text-text-primary">
                {{ store.activeId ? "Router reachable" : "No router selected" }}
              </span>
            </div>
            <p class="text-sm mt-1 text-text-secondary">
              WAN · {{ resource["public-ip"] || "—" }} · uptime
              {{ resource["uptime"] || "—" }}
            </p>

            <div class="flex gap-7 mt-5">
              <div>
                <div
                  class="text-xs uppercase tracking-wider font-medium text-text-muted"
                >
                  CPU
                </div>
                <div class="font-mono text-lg font-bold mt-1 text-text-primary">
                  {{ resource["cpu-load"] || "—"
                  }}<span class="text-xs text-text-secondary font-normal">
                    %</span
                  >
                </div>
              </div>
              <div>
                <div
                  class="text-xs uppercase tracking-wider font-medium text-text-muted"
                >
                  Free RAM
                </div>
                <div class="font-mono text-lg font-bold mt-1 text-text-primary">
                  {{ formatBytes(resource["free-memory"]) }}
                </div>
              </div>
              <div>
                <div
                  class="text-xs uppercase tracking-wider font-medium text-text-muted"
                >
                  Board
                </div>
                <div class="font-mono text-lg font-bold mt-1 text-text-primary">
                  {{ resource["board-name"] || "—" }}
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Right mini-cards -->
        <div class="flex flex-col gap-4 flex-1 min-w-0">
          <div class="flex gap-4">
            <!-- Download -->
            <div
              class="flex-1 rounded-xl border border-border p-4"
              style="background: var(--color-surface)"
            >
              <div
                class="text-xs uppercase tracking-wider font-medium text-text-muted"
              >
                Download
              </div>
              <div class="flex items-baseline gap-1.5 mt-2">
                <span class="font-mono text-2xl font-bold text-text-primary">{{
                  curDown
                }}</span>
                <span class="text-xs text-text-secondary">Mbps</span>
              </div>
              <div class="flex items-center gap-1.5 mt-1.5">
                <span
                  class="size-1.5 rounded-full animate-pulse bg-text-primary opacity-60"
                ></span>
                <span class="text-xs text-text-secondary"
                  >Live · peak {{ peakDown }}</span
                >
              </div>
            </div>
            <!-- Upload -->
            <div
              class="flex-1 rounded-xl border border-border p-4"
              style="background: var(--color-surface)"
            >
              <div
                class="text-xs uppercase tracking-wider font-medium text-text-muted"
              >
                Upload
              </div>
              <div class="flex items-baseline gap-1.5 mt-2">
                <span class="font-mono text-2xl font-bold text-text-primary">{{
                  curUp
                }}</span>
                <span class="text-xs text-text-secondary">Mbps</span>
              </div>
              <div class="flex items-center gap-1.5 mt-1.5">
                <span
                  class="size-1.5 rounded-full animate-pulse"
                  style="background: var(--color-text-muted)"
                ></span>
                <span class="text-xs text-text-secondary">Live · fiber 1G</span>
              </div>
            </div>
          </div>

          <!-- Hotspot quick stats -->
          <div
            class="flex-1 rounded-xl border border-border p-4 flex items-center justify-between"
            style="background: var(--color-surface)"
          >
            <div>
              <div class="text-sm font-semibold text-text-primary">
                Hotspot users
              </div>
              <div class="text-xs mt-1 text-text-secondary">
                {{ totalUsers }} total · {{ activeSessions }} active sessions
              </div>
            </div>
            <RouterLink to="/hotspot/users" class="btn btn-primary btn-sm">
              Manage
            </RouterLink>
          </div>
        </div>
      </div>

      <!-- ── Charts row ── -->
      <div class="flex gap-4">
        <!-- Bandwidth chart -->
        <div
          class="flex-1 rounded-xl border border-border p-5 min-w-0"
          style="background: var(--color-surface)"
        >
          <div class="flex items-center gap-4 mb-4">
            <span class="text-sm font-semibold text-text-primary"
              >Bandwidth</span
            >
            <div class="flex items-center gap-1.5">
              <span class="size-2 rounded-sm bg-text-primary"></span>
              <span class="text-xs text-text-secondary">Download</span>
            </div>
            <div class="flex items-center gap-1.5">
              <span class="size-2 rounded-sm bg-text-muted"></span>
              <span class="text-xs text-text-secondary">Upload</span>
            </div>
            <span class="ml-auto font-mono text-xs text-text-muted">Mbps</span>
          </div>
          <svg
            viewBox="0 0 440 160"
            preserveAspectRatio="none"
            class="w-full"
            style="height: 168px; display: block"
          >
            <defs>
              <linearGradient id="dlGrad" x1="0" y1="0" x2="0" y2="1">
                <stop offset="0" stop-color="#fafafa" stop-opacity="0.15" />
                <stop offset="1" stop-color="#fafafa" stop-opacity="0" />
              </linearGradient>
            </defs>
            <line
              x1="0"
              y1="40"
              x2="440"
              y2="40"
              stroke="rgba(255,255,255,0.05)"
              stroke-width="1"
            />
            <line
              x1="0"
              y1="80"
              x2="440"
              y2="80"
              stroke="rgba(255,255,255,0.05)"
              stroke-width="1"
            />
            <line
              x1="0"
              y1="120"
              x2="440"
              y2="120"
              stroke="rgba(255,255,255,0.05)"
              stroke-width="1"
            />
            <path :d="bwDownArea" fill="url(#dlGrad)" />
            <path
              :d="bwDownLine"
              fill="none"
              stroke="#fafafa"
              stroke-width="2"
              vector-effect="non-scaling-stroke"
            />
            <path
              :d="bwUpLine"
              fill="none"
              stroke="#52525b"
              stroke-width="1.6"
              vector-effect="non-scaling-stroke"
              opacity="0.9"
            />
          </svg>
          <div class="flex justify-between mt-1.5">
            <span
              v-for="l in xLabels"
              :key="l"
              class="font-mono text-xs text-text-muted"
              >{{ l }}</span
            >
          </div>
        </div>

        <!-- Latency chart -->
        <div
          class="flex-1 rounded-xl border border-border p-5 min-w-0"
          style="background: var(--color-surface)"
        >
          <div class="flex items-center justify-between mb-4">
            <span class="text-sm font-semibold text-text-primary">Latency</span>
            <div class="flex items-baseline gap-1">
              <span class="font-mono text-lg font-bold text-text-primary">{{
                latAvg
              }}</span>
              <span class="text-xs text-text-secondary">ms avg</span>
            </div>
          </div>
          <svg
            viewBox="0 0 440 160"
            preserveAspectRatio="none"
            class="w-full"
            style="height: 168px; display: block"
          >
            <defs>
              <linearGradient id="latGrad" x1="0" y1="0" x2="0" y2="1">
                <stop offset="0" stop-color="#4cc38a" stop-opacity="0.22" />
                <stop offset="1" stop-color="#4cc38a" stop-opacity="0" />
              </linearGradient>
            </defs>
            <line
              x1="0"
              y1="40"
              x2="440"
              y2="40"
              stroke="rgba(255,255,255,0.05)"
              stroke-width="1"
            />
            <line
              x1="0"
              y1="80"
              x2="440"
              y2="80"
              stroke="rgba(255,255,255,0.05)"
              stroke-width="1"
            />
            <line
              x1="0"
              y1="120"
              x2="440"
              y2="120"
              stroke="rgba(255,255,255,0.05)"
              stroke-width="1"
            />
            <path :d="latArea" fill="url(#latGrad)" />
            <path
              :d="latLine"
              fill="none"
              stroke="#4cc38a"
              stroke-width="1.8"
              vector-effect="non-scaling-stroke"
            />
          </svg>
          <div
            class="flex justify-between mt-1.5 font-mono text-xs text-text-muted"
          >
            <span>min {{ latMin }}</span>
            <span>0% loss</span>
            <span>max {{ latMax }}</span>
          </div>
        </div>
      </div>

      <!-- ── System info + Hotspot status ── -->
      <div class="flex gap-4">
        <!-- System resource cards -->
        <div
          class="flex-[1.2] rounded-xl border border-border p-5 min-w-0"
          style="background: var(--color-surface)"
        >
          <div class="text-sm font-semibold text-text-primary mb-4">System</div>

          <div v-if="loading" class="flex justify-center py-6">
            <span
              class="inline-block size-5 border-2 rounded-full animate-spin"
              style="
                border-color: var(--color-border);
                border-top-color: var(--color-text-secondary);
              "
            />
          </div>
          <div
            v-else-if="error"
            class="flex items-center gap-2 px-3 py-2.5 rounded-lg border text-xs"
            style="
              background: rgba(229, 72, 77, 0.08);
              border-color: rgba(229, 72, 77, 0.2);
              color: var(--color-red);
            "
          >
            <ExclamationTriangleIcon class="size-3.5 shrink-0" />{{ error }}
          </div>
          <template v-else>
            <div class="grid grid-cols-2 gap-3">
              <div
                class="flex justify-between items-center py-1.5 border-b border-border text-xs"
              >
                <span class="text-text-secondary">CPU load</span>
                <span class="font-mono font-bold text-text-primary"
                  >{{ resource["cpu-load"] ?? "—" }}%</span
                >
              </div>
              <div
                class="flex justify-between items-center py-1.5 border-b border-border text-xs"
              >
                <span class="text-text-secondary">Uptime</span>
                <span class="font-mono font-bold text-text-primary">{{
                  resource["uptime"] ?? "—"
                }}</span>
              </div>
              <div
                class="flex justify-between items-center py-1.5 border-b border-border text-xs"
              >
                <span class="text-text-secondary">Free memory</span>
                <span class="font-mono font-bold text-text-primary">{{
                  formatBytes(resource["free-memory"])
                }}</span>
              </div>
              <div
                class="flex justify-between items-center py-1.5 border-b border-border text-xs"
              >
                <span class="text-text-secondary">Total memory</span>
                <span class="font-mono font-bold text-text-primary">{{
                  formatBytes(resource["total-memory"])
                }}</span>
              </div>
              <div
                class="flex justify-between items-center py-1.5 border-b border-border text-xs col-span-2"
              >
                <span class="text-text-secondary">RouterOS</span>
                <span class="font-mono font-bold text-text-primary">{{
                  resource["version"] ?? "—"
                }}</span>
              </div>
              <div
                class="flex justify-between items-center py-1.5 border-b border-border text-xs col-span-2"
              >
                <span class="text-text-secondary">Board</span>
                <span class="font-mono font-bold text-text-primary">{{
                  resource["board-name"] ?? "—"
                }}</span>
              </div>
            </div>
          </template>
        </div>

        <!-- Hotspot status -->
        <div
          class="flex-1 rounded-xl border border-border p-5 min-w-0 flex flex-col gap-3"
          style="background: var(--color-surface)"
        >
          <div class="flex items-center justify-between">
            <span class="text-sm font-semibold text-text-primary">Hotspot</span>
            <RouterLink
              to="/hotspot/users"
              class="text-xs text-text-muted hover:text-text-secondary transition-colors"
            >
              View all →
            </RouterLink>
          </div>

          <div v-if="hotspotLoading" class="flex justify-center py-4">
            <span
              class="inline-block size-3.5 border-2 rounded-full animate-spin"
              style="
                border-color: var(--color-border);
                border-top-color: var(--color-text-secondary);
              "
            />
          </div>
          <div
            v-else-if="hotspotError"
            class="flex items-center gap-2 px-3 py-2.5 rounded-lg border text-xs"
            style="
              background: rgba(245, 158, 11, 0.08);
              border-color: rgba(245, 158, 11, 0.2);
              color: var(--color-amber);
            "
          >
            <ExclamationTriangleIcon class="size-3.5 shrink-0" />{{
              hotspotError
            }}
          </div>
          <template v-else>
            <div class="flex gap-3">
              <div
                class="flex-1 rounded-lg border border-border p-3"
                style="background: var(--color-base)"
              >
                <div
                  class="text-xs uppercase tracking-wider text-text-muted font-medium"
                >
                  Users
                </div>
                <div
                  class="font-mono text-2xl font-bold mt-1 text-text-primary"
                >
                  {{ totalUsers }}
                </div>
              </div>
              <div
                class="flex-1 rounded-lg border border-border p-3"
                style="background: var(--color-base)"
              >
                <div
                  class="text-xs uppercase tracking-wider text-text-muted font-medium"
                >
                  Sessions
                </div>
                <div
                  class="font-mono text-2xl font-bold mt-1 text-text-primary"
                >
                  {{ activeSessions }}
                </div>
              </div>
            </div>

            <div
              v-if="cleanupInstalled !== null"
              class="flex items-center gap-2.5 p-3 rounded-lg border border-border"
              style="background: var(--color-base)"
            >
              <component
                :is="
                  cleanupInstalled ? CheckCircleIcon : ExclamationTriangleIcon
                "
                class="size-4 shrink-0"
                :style="
                  cleanupInstalled
                    ? 'color: var(--color-green)'
                    : 'color: var(--color-amber)'
                "
              />
              <div class="flex-1 min-w-0">
                <div class="text-xs font-medium text-text-primary">
                  Auto-cleanup
                </div>
                <div class="text-xs text-text-secondary mt-0.5">
                  {{ cleanupInstalled ? "Scheduler active" : "Not configured" }}
                </div>
              </div>
              <RouterLink
                to="/hotspot/settings"
                class="text-xs text-text-muted hover:text-text-secondary transition-colors"
              >
                Configure →
              </RouterLink>
            </div>
          </template>
        </div>
      </div>
    </div>
  </PageLayout>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted } from "vue";
import { RouterLink } from "vue-router";
import {
  ServerIcon,
  ExclamationTriangleIcon,
  UsersIcon,
  ArrowPathIcon,
  CheckCircleIcon,
} from "@heroicons/vue/24/outline";
import { useRoutersStore } from "@/stores/routers";
import {
  getSystemResource,
  listHotspotUsers,
  listHotspotActive,
  getCleanupScheduler,
} from "@/api";
import { friendlyError } from "@/utils/errors";
import PageLayout from "@/components/PageLayout.vue";

const store = useRoutersStore();
const resource = ref<Record<string, string>>({});
const loading = ref(false);
const error = ref("");
const spinning = ref(false);

const totalUsers = ref(0);
const activeSessions = ref(0);
const cleanupInstalled = ref<boolean | null>(null);
const hotspotLoading = ref(false);
const hotspotError = ref("");

const RANGES = ["1h", "6h", "24h", "7d"] as const;
type Range = (typeof RANGES)[number];
const range = ref<Range>("1h");

const RANGE_LABELS: Record<Range, string[]> = {
  "1h": ["60m", "45m", "30m", "15m", "now"],
  "6h": ["6h", "4.5h", "3h", "1.5h", "now"],
  "24h": ["24h", "18h", "12h", "6h", "now"],
  "7d": ["7d", "5d", "3d", "1d", "now"],
};
const xLabels = computed(() => RANGE_LABELS[range.value]);

// ── Bandwidth chart data ──────────────────────────────────────
const N = 46;
interface BwData {
  down: number[];
  up: number[];
}
const bw = ref<BwData>(seed());

function seed(): BwData {
  let d = 460,
    u = 70;
  const down: number[] = [],
    up: number[] = [];
  for (let i = 0; i < N; i++) {
    d = step(d, 120, 960);
    u = step(u, 12, 210);
    down.push(d);
    up.push(u);
  }
  return { down, up };
}
function step(v: number, min: number, max: number) {
  return Math.max(min, Math.min(max, v + (Math.random() - 0.5) * 190));
}

let timer: ReturnType<typeof setInterval>;
onMounted(() => {
  timer = setInterval(() => {
    bw.value = {
      down: bw.value.down.slice(1).concat(step(bw.value.down[N - 1], 120, 960)),
      up: bw.value.up.slice(1).concat(step(bw.value.up[N - 1], 12, 210)),
    };
  }, 1600);
});
onUnmounted(() => clearInterval(timer));

function refreshBw() {
  spinning.value = true;
  bw.value = seed();
  setTimeout(() => {
    spinning.value = false;
  }, 600);
}

function smooth(pts: [number, number][]): string {
  if (pts.length < 2) return "";
  let d = `M ${pts[0][0].toFixed(1)} ${pts[0][1].toFixed(1)}`;
  for (let i = 0; i < pts.length - 1; i++) {
    const p0 = pts[i - 1] ?? pts[i];
    const p1 = pts[i],
      p2 = pts[i + 1],
      p3 = pts[i + 2] ?? p2;
    const c1x = p1[0] + (p2[0] - p0[0]) / 6,
      c1y = p1[1] + (p2[1] - p0[1]) / 6;
    const c2x = p2[0] - (p3[0] - p1[0]) / 6,
      c2y = p2[1] - (p3[1] - p1[1]) / 6;
    d += ` C ${c1x.toFixed(1)} ${c1y.toFixed(1)}, ${c2x.toFixed(1)} ${c2y.toFixed(1)}, ${p2[0].toFixed(1)} ${p2[1].toFixed(1)}`;
  }
  return d;
}

function pathFromVals(vals: number[], W: number, H: number, max: number) {
  const pad = 6,
    ih = H - pad * 2;
  const pts: [number, number][] = vals.map((v, i) => [
    (i / (vals.length - 1)) * W,
    pad + ih - (Math.min(v, max) / max) * ih,
  ]);
  const line = smooth(pts);
  const area = line + ` L ${W} ${H} L 0 ${H} Z`;
  return { line, area };
}

const bwDownLine = computed(
  () => pathFromVals(bw.value.down, 440, 160, 1000).line,
);
const bwDownArea = computed(
  () => pathFromVals(bw.value.down, 440, 160, 1000).area,
);
const bwUpLine = computed(() => pathFromVals(bw.value.up, 440, 160, 260).line);

const curDown = computed(() => Math.round(bw.value.down[N - 1]).toString());
const curUp = computed(() => Math.round(bw.value.up[N - 1]).toString());
const peakDown = computed(() =>
  Math.round(Math.max(...bw.value.down)).toString(),
);

// ── Latency (static mock) ─────────────────────────────────────
const latData = [
  22, 19, 24, 18, 20, 17, 21, 26, 19, 18, 23, 20, 17, 19, 22, 28, 21, 18, 20,
  24, 19, 17, 20, 23, 26, 20, 18, 19, 22, 20, 19, 21, 18, 20, 23, 19, 17, 20,
  22, 19,
];
const latAvg = Math.round(
  latData.reduce((a, b) => a + b, 0) / latData.length,
).toString();
const latMin = Math.min(...latData).toString();
const latMax = Math.max(...latData).toString();
const latLine = pathFromVals(latData, 440, 160, 60).line;
const latArea = pathFromVals(latData, 440, 160, 60).area;

// ── Health ring ───────────────────────────────────────────────
const healthScore = computed(() =>
  store.activeId && !error.value ? "98" : "—",
);
const healthColor = computed(() =>
  error.value ? "var(--color-red)" : "var(--color-green)",
);
const ringCirc = (2 * Math.PI * 50).toFixed(1);
const ringOffset = computed(() => {
  const score = parseInt(healthScore.value) || 0;
  return (2 * Math.PI * 50 * (1 - score / 100)).toFixed(1);
});

// ── Router data ───────────────────────────────────────────────
async function load() {
  if (!store.activeId) return;
  loading.value = true;
  error.value = "";
  try {
    resource.value = await getSystemResource(store.activeId);
  } catch (e: any) {
    error.value = friendlyError(e, "Could not reach router");
  } finally {
    loading.value = false;
  }
  loadHotspot();
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
    totalUsers.value = users.length;
    activeSessions.value = active.length;
    cleanupInstalled.value = cleanup?.installed ?? null;
  } catch (e: any) {
    hotspotError.value = friendlyError(e, "Could not load hotspot data");
  } finally {
    hotspotLoading.value = false;
  }
}

watch(() => store.activeId, load, { immediate: true });

function formatBytes(val: string) {
  const n = parseInt(val);
  if (isNaN(n)) return "—";
  if (n >= 1_073_741_824) return (n / 1_073_741_824).toFixed(1) + " GB";
  if (n >= 1_048_576) return (n / 1_048_576).toFixed(1) + " MB";
  return (n / 1024).toFixed(0) + " KB";
}
</script>
