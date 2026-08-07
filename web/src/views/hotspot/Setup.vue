<template>
  <PageLayout title="Hotspot" subtitle="Setup Wizard">
    <div
      v-if="phase === 'loading'"
      class="flex flex-col items-center py-16 gap-3"
    >
      <span class="spinner" />
      <p class="text-sm text-text-secondary">Checking router configuration…</p>
    </div>

    <!-- ── Phase: blocked ─────────────────────────────────────────────────── -->
    <div v-else-if="phase === 'blocked'" class="space-y-4">
      <div
        class="flex text-sm items-start gap-3 p-4 border rounded-xl bg-amber/8 border-amber/20"
      >
        <ExclamationTriangleIcon class="size-5 shrink-0 mt-0.5 text-amber" />
        <div class="flex-1 min-w-0">
          <p class="font-medium text-amber">Hotspot already configured</p>
          <p class="text-text-secondary mt-1">
            A hotspot is active on this router. Reset it below to reconfigure,
            or go to Users to manage it.
          </p>
        </div>
      </div>

      <div
        class="border border-border rounded-xl p-4 text-sm space-y-2 bg-surface"
      >
        <p class="font-semibold text-text-secondary uppercase mb-3">
          Detected hotspot
        </p>
        <div class="grid grid-cols-2 gap-x-4 gap-y-2">
          <span class="text-text-secondary">Server name</span>
          <span class="font-mono text-text-primary">{{
            preflight?.hotspotName || "—"
          }}</span>
          <span class="text-text-secondary">Interface</span>
          <span class="font-mono text-text-primary">{{
            preflight?.hotspotOnIface || "—"
          }}</span>
          <span class="text-text-secondary">Profile</span>
          <span class="font-mono text-text-primary">{{
            preflight?.hotspotProfile || "—"
          }}</span>
          <span class="text-text-secondary">Address pool</span>
          <span class="font-mono text-text-primary">{{
            preflight?.hotspotAddressPool || "—"
          }}</span>
          <span class="text-text-secondary">DNS name</span>
          <span class="font-mono text-text-primary">{{
            preflight?.hotspotDnsName || "(none)"
          }}</span>
        </div>
      </div>

      <div class="flex items-center gap-3">
        <RouterLink to="/hotspot/settings" class="btn btn-primary mr-auto">
          <Cog6ToothIcon class="size-3.5" />
          Configure branding &amp; vouchers
        </RouterLink>
        <button class="btn btn-ghost" @click="runPreflight">
          <ArrowPathIcon class="size-3.5" />
          Re-check
        </button>
        <button
          class="btn border-red/30 text-red hover:bg-red/8"
          :disabled="tearing"
          @click="runTeardown"
        >
          <TrashIcon class="size-3.5" />
          {{ tearing ? "Resetting…" : "Reset config" }}
        </button>
      </div>

      <div
        v-if="teardownResult"
        class="border border-border rounded-xl p-4 bg-surface"
      >
        <p
          class="text-xs font-semibold text-text-muted uppercase tracking-wide mb-2"
        >
          Reset result
        </p>
        <ul class="space-y-1.5">
          <li
            v-for="step in teardownResult.steps"
            :key="step.name"
            class="flex items-center gap-2 text-sm"
          >
            <CheckCircleIcon
              v-if="step.ok && !step.skipped"
              class="size-3.5 shrink-0 text-green"
            />
            <MinusCircleIcon
              v-else-if="step.skipped"
              class="size-3.5 text-text-muted shrink-0"
            />
            <XCircleIcon v-else class="size-3.5 shrink-0 text-red" />
            <span
              :class="step.skipped ? 'text-text-muted' : 'text-text-secondary'"
              >{{ step.name }}</span
            >
            <span v-if="step.error" class="font-mono text-red"
              >— {{ step.error }}</span
            >
          </li>
        </ul>
      </div>
    </div>

    <!-- ── Phase: form ────────────────────────────────────────────────────── -->
    <div v-else-if="phase === 'form'" class="space-y-5">
      <div
        class="flex items-start gap-2 p-3 border rounded-lg text-xs bg-amber/8 border-amber/20 text-amber"
      >
        <ExclamationTriangleIcon class="size-4 shrink-0 mt-0.5" />
        <span
          >Ensure no hotspot is already configured on this router. Running setup
          on an existing configuration may cause conflicts.</span
        >
      </div>

      <div
        v-if="preflightError"
        class="flex items-start gap-2 p-3 border rounded-lg text-xs bg-red/8 border-red/20 text-red"
      >
        <ExclamationTriangleIcon class="size-4 shrink-0 mt-0.5" />
        <span
          >Could not read interfaces: {{ preflightError }}. Enter interface
          names manually.</span
        >
      </div>

      <div class="border border-border rounded-xl p-5 space-y-4 bg-surface">
        <h2 class="text-sm font-semibold text-text-primary">
          Network configuration
        </h2>

        <label class="flex flex-col gap-1">
          <span class="text-sm font-medium text-red"
            >LAN / Hotspot interface <span>*</span></span
          >
          <select
            v-if="interfaces.length"
            v-model="form.lanIface"
            class="input"
          >
            <option value="" disabled>Select interface…</option>
            <optgroup
              v-if="bridgeAndWlan.length"
              label="Bridge / WiFi (recommended)"
            >
              <option v-for="i in bridgeAndWlan" :key="i.name" :value="i.name">
                {{ i.name }}{{ !i.running ? " — down" : ""
                }}{{ i.comment ? ` (${i.comment})` : "" }}
              </option>
            </optgroup>
            <optgroup label="All interfaces">
              <option v-for="i in interfaces" :key="i.name" :value="i.name">
                {{ i.name }} [{{ i.type }}]{{ !i.running ? " — down" : "" }}
              </option>
            </optgroup>
          </select>
          <input
            v-else
            v-model="form.lanIface"
            class="input"
            placeholder="e.g. bridge1 or ether2"
          />
          <p class="text-xs text-text-muted">
            The interface clients connect to (usually a bridge or wlan)
          </p>
        </label>

        <label class="flex flex-col gap-1">
          <span class="text-sm font-medium text-red"
            >WAN interface <span>*</span></span
          >
          <select
            v-if="interfaces.length"
            v-model="form.wanIface"
            class="input"
          >
            <option value="" disabled>Select interface…</option>
            <option v-for="i in interfaces" :key="i.name" :value="i.name">
              {{ i.name }} [{{ i.type }}]{{ !i.running ? " — down" : "" }}
            </option>
          </select>
          <input
            v-else
            v-model="form.wanIface"
            class="input"
            placeholder="e.g. ether1"
          />
          <p class="text-xs text-text-muted">
            The interface connected to the internet (usually ether1)
          </p>
        </label>

        <label class="flex flex-col gap-1">
          <span class="text-sm font-medium text-red"
            >Subnet (CIDR) <span>*</span></span
          >
          <input
            v-model="form.subnet"
            class="input font-mono"
            placeholder="192.168.88.0/24"
          />
        </label>

        <label class="flex flex-col gap-1">
          <span class="text-sm font-medium text-text-secondary"
            >Hotspot name</span
          >
          <div class="flex items-stretch">
            <input
              v-model="form.hotspotName"
              class="input rounded-r-none border-r-0"
              placeholder="myhotspot"
              pattern="[a-zA-Z0-9\-]+"
            />
            <span
              class="flex items-center px-3 border border-border rounded-r-lg text-xs font-mono text-text-muted shrink-0 bg-surface"
              >.spot</span
            >
          </div>
          <p class="text-xs text-text-muted font-mono">
            Clients are redirected to
            <span>{{ form.hotspotName || "myhotspot" }}.spot</span> — leave
            blank to use the router IP instead
          </p>
        </label>

        <div
          v-if="derived"
          class="border border-border rounded-lg p-3 text-xs bg-base"
        >
          <div class="grid grid-cols-2 gap-x-4 gap-y-1.5">
            <span class="text-text-muted">Router IP (gateway)</span>
            <span class="font-mono text-text-secondary">{{
              derived.gateway
            }}</span>
            <span class="text-text-muted">DHCP pool</span>
            <span class="font-mono text-text-secondary"
              >{{ derived.poolStart }} – {{ derived.poolEnd }}</span
            >
            <span class="text-text-muted">Hotspot profile name</span>
            <span class="font-mono text-text-secondary">pikro-profile</span>
          </div>
        </div>
      </div>

      <details class="text-xs text-text-muted">
        <summary
          class="cursor-pointer hover:text-text-secondary transition-colors font-medium"
        >
          What will be configured
        </summary>
        <ul class="mt-2 space-y-1 pl-2 text-text-muted">
          <li>
            • Assign {{ derived?.gateway ?? "…" }} to
            {{ form.lanIface || "(LAN interface)" }}
          </li>
          <li>• Create IP pool <span class="font-mono">hotspot-pool</span></li>
          <li>• Create DHCP server + network</li>
          <li>• Create hotspot server profile with Pikro login page</li>
          <li>• Enable hotspot on {{ form.lanIface || "(LAN interface)" }}</li>
          <li>
            • Add NAT masquerade on {{ form.wanIface || "(WAN interface)" }}
          </li>
          <li>• Enable DNS remote requests</li>
        </ul>
      </details>

      <!-- Branding (form phase) -->
      <div class="border border-border rounded-xl p-5 space-y-4 bg-surface">
        <h2 class="text-sm font-semibold text-text-primary">Vouchers</h2>
        <label class="flex flex-col gap-1">
          <span class="text-sm font-medium text-text-secondary">Currency</span>
          <select v-model="branding.currency" class="input">
            <option value="">None</option>
            <option v-for="c in CURRENCIES" :key="c.value" :value="c.value">
              {{ c.label }}
            </option>
          </select>
        </label>
        <p class="text-xs text-text-muted">
          Saved alongside the network setup. You can update this anytime by
          editing the router.
        </p>
      </div>

      <div class="flex items-center gap-3">
        <button
          class="btn btn-primary"
          :disabled="!canSubmit"
          @click="runSetup"
        >
          Run Setup
        </button>
        <button
          class="btn btn-ghost btn-sm border-transparent"
          @click="runPreflight"
        >
          <ArrowPathIcon class="size-3.5" />
          Re-check router
        </button>
        <button
          class="btn btn-sm ml-auto border-red/30 text-red hover:bg-red/8"
          :disabled="tearing"
          @click="runTeardown"
        >
          <TrashIcon class="size-3.5" />
          {{ tearing ? "Resetting…" : "Reset config" }}
        </button>
      </div>

      <div
        v-if="teardownResult"
        class="border border-border rounded-xl p-4 bg-surface"
      >
        <p
          class="text-xs font-semibold text-text-muted uppercase tracking-wide mb-2"
        >
          Reset result
        </p>
        <ul class="space-y-1.5">
          <li
            v-for="step in teardownResult.steps"
            :key="step.name"
            class="flex items-center gap-2 text-sm"
          >
            <CheckCircleIcon
              v-if="step.ok && !step.skipped"
              class="size-3.5 shrink-0 text-green"
            />
            <MinusCircleIcon
              v-else-if="step.skipped"
              class="size-3.5 text-text-muted shrink-0"
            />
            <XCircleIcon v-else class="size-3.5 shrink-0 text-red" />
            <span
              :class="step.skipped ? 'text-text-muted' : 'text-text-secondary'"
              >{{ step.name }}</span
            >
            <span v-if="step.error" class="font-mono text-red"
              >— {{ step.error }}</span
            >
          </li>
        </ul>
      </div>
    </div>

    <!-- ── Phase: running ─────────────────────────────────────────────────── -->
    <div
      v-else-if="phase === 'running'"
      class="flex flex-col items-center py-16 gap-3"
    >
      <span class="spinner" />
      <p class="text-sm font-medium text-text-primary">Configuring hotspot…</p>
      <p class="text-xs text-text-muted">
        This takes a few seconds. Please wait.
      </p>
    </div>

    <!-- ── Phase: result ──────────────────────────────────────────────────── -->
    <div v-else-if="phase === 'result' && result" class="space-y-4">
      <div
        class="flex items-center gap-2 p-3 rounded-lg text-sm font-medium border"
        :class="
          result.success
            ? 'bg-green/8 border-green/20 text-green'
            : 'bg-red/8 border-red/20 text-red'
        "
      >
        <CheckCircleIcon v-if="result.success" class="size-4 shrink-0" />
        <ExclamationTriangleIcon v-else class="size-4 shrink-0" />
        {{
          result.success
            ? "Hotspot configured successfully!"
            : "Setup completed with errors"
        }}
      </div>

      <div class="border border-border rounded-xl p-4 bg-surface">
        <ul class="space-y-2">
          <li
            v-for="step in result.steps"
            :key="step.name"
            class="flex items-start gap-2 text-sm"
          >
            <CheckCircleIcon
              v-if="step.ok && !step.skipped"
              class="size-3.5 shrink-0 mt-0.5 text-green"
            />
            <MinusCircleIcon
              v-else-if="step.skipped"
              class="size-3.5 text-text-muted shrink-0 mt-0.5"
            />
            <XCircleIcon v-else class="size-3.5 shrink-0 mt-0.5 text-red" />
            <div class="flex-1 min-w-0">
              <span
                :class="
                  step.ok
                    ? 'text-text-secondary'
                    : step.skipped
                      ? 'text-text-muted'
                      : 'font-medium text-red'
                "
              >
                {{ step.name }}
                <span v-if="step.skipped" class="text-text-muted font-normal">
                  — already configured</span
                >
              </span>
              <p v-if="step.error" class="font-mono mt-0.5 break-all text-red">
                {{ step.error }}
              </p>
            </div>
          </li>
        </ul>
      </div>

      <div
        v-if="result.success"
        class="border border-border rounded-xl p-4 space-y-3 bg-surface"
      >
        <p
          class="text-xs font-semibold text-text-muted uppercase tracking-wide"
        >
          What's next
        </p>
        <ol class="space-y-2.5 text-sm text-text-secondary">
          <li class="flex items-start gap-2">
            <span
              class="size-5 shrink-0 rounded-full flex items-center justify-center text-xs font-bold mt-0.5 bg-muted text-text-primary"
              >1</span
            >
            <div>
              <span
                >Create your first hotspot user with a time or data limit.</span
              >
              <span class="block mt-0.5 text-amber">
                RouterOS added a default
                <span class="font-mono">admin</span> user with no password —
                delete or secure it from the Users page.
              </span>
            </div>
          </li>
          <li class="flex items-start gap-2">
            <span
              class="size-5 shrink-0 rounded-full flex items-center justify-center text-xs font-bold mt-0.5 bg-muted text-text-primary"
              >2</span
            >
            <span
              >Connect a client device to the router's WiFi or LAN — it will be
              redirected to
              <span class="font-mono">{{
                form.hotspotName ? form.hotspotName + ".spot" : "the login page"
              }}</span
              >.</span
            >
          </li>
          <li class="flex items-start gap-2">
            <span
              class="size-5 shrink-0 rounded-full flex items-center justify-center text-xs font-bold mt-0.5 bg-muted text-text-primary"
              >3</span
            >
            <span
              >Share the credentials — or create multiple users in bulk for
              voucher-style access.</span
            >
          </li>
          <li class="flex items-start gap-2">
            <span
              class="size-5 shrink-0 rounded-full flex items-center justify-center text-xs font-bold mt-0.5 bg-muted text-text-primary"
              >4</span
            >
            <span>
              Configure your business name, currency, voucher layout, and login
              page in
              <RouterLink
                to="/hotspot/settings"
                class="underline text-text-primary hover:text-text-secondary transition-colors"
                >Settings</RouterLink
              >.
            </span>
          </li>
        </ol>
      </div>

      <div
        v-if="!result.success"
        class="text-sm border border-border rounded-lg p-3 space-y-1 bg-surface"
      >
        <p class="font-medium text-text-primary">Some steps failed</p>
        <p class="text-text-muted">
          Steps marked as "already configured" are fine — they were skipped
          because the resource exists. Fix any red steps manually in RouterOS,
          then run setup again.
        </p>
      </div>

      <div class="flex items-center gap-3">
        <RouterLink
          v-if="result.success"
          to="/hotspot/users"
          class="btn btn-primary"
        >
          <UsersIcon class="size-4" />
          Add users
        </RouterLink>
        <button class="btn btn-ghost btn-sm border-transparent" @click="reset">
          <ArrowPathIcon class="size-3.5" />
          Run again
        </button>
      </div>
    </div>

    <!-- ── Teardown confirmation dialog ─────────────────────────────────── -->
    <AppDialog
      :open="showTeardownConfirm"
      title="Reset hotspot config?"
      @update:open="showTeardownConfirm = $event"
    >
      <div class="space-y-4">
        <div
          class="border border-border rounded-lg p-3 text-sm space-y-1.5 bg-base"
        >
          <div class="flex justify-between font-medium">
            <span class="text-text-secondary">Router</span>
            <span class="text-text-primary">{{
              store.routers.find((r) => r.id === store.activeId)?.name
            }}</span>
          </div>
          <div class="flex justify-between font-medium">
            <span class="text-text-secondary">Host</span>
            <span class="font-mono text-text-primary">{{
              store.routers.find((r) => r.id === store.activeId)?.host
            }}</span>
          </div>
          <div class="flex justify-between font-medium">
            <span class="text-text-secondary">Hotspot interface</span>
            <span class="font-mono text-text-primary">{{
              preflight?.hotspotOnIface || form.lanIface || "—"
            }}</span>
          </div>
        </div>

        <div>
          <p class="font-medium text-text-primary text-sm mb-1.5">
            This will remove:
          </p>
          <ul class="text-sm text-text-secondary space-y-1">
            <li class="flex items-center gap-1.5 shrink-0 text-red">
              <TrashIcon class="size-4" /> Hotspot server
            </li>
            <li class="flex items-center gap-1.5 shrink-0 text-red">
              <TrashIcon class="size-4" /> Hotspot server profile
            </li>
            <li class="flex items-center gap-1.5 shrink-0 text-red">
              <TrashIcon class="size-4" /> NAT masquerade rule
            </li>
            <li class="flex items-center gap-1.5 shrink-0 text-red">
              <TrashIcon class="size-4" /> Walled garden rules
            </li>
          </ul>
          <p class="text-sm text-text-secondary font-medium mt-2">
            IP address, DHCP server, IP pool and DNS settings are kept.
          </p>
        </div>

        <div class="flex justify-end gap-2 pt-1">
          <button class="btn btn-ghost" @click="showTeardownConfirm = false">
            Cancel
          </button>
          <button class="btn btn-danger" @click="confirmTeardown">
            Reset config
          </button>
        </div>
      </div>
    </AppDialog>
  </PageLayout>
</template>

<script setup lang="ts">
import { ref, computed, watch } from "vue";
import {
  ExclamationTriangleIcon,
  ArrowPathIcon,
  CheckCircleIcon,
  XCircleIcon,
  MinusCircleIcon,
  UsersIcon,
  TrashIcon,
  Cog6ToothIcon,
} from "@heroicons/vue/24/outline";
import { RouterLink } from "vue-router";
import {
  hotspotPreflight,
  setupHotspot,
  teardownHotspot,
  getHotspotSettings,
  putHotspotSettings,
  type PreflightResult,
  type SetupResult,
  type TeardownResult,
  type InterfaceInfo,
} from "@/api";
import { useRoutersStore } from "@/stores/routers";
import { friendlyError } from "@/utils/errors";
import { CURRENCIES } from "@/utils/currencies";
import AppDialog from "@/components/AppDialog.vue";
import PageLayout from "@/components/PageLayout.vue";

const store = useRoutersStore();

type Phase = "loading" | "form" | "blocked" | "running" | "result";
const phase = ref<Phase>("loading");
const preflight = ref<PreflightResult | null>(null);
const preflightError = ref("");
const result = ref<SetupResult | null>(null);
const interfaces = ref<InterfaceInfo[]>([]);
const form = ref({
  lanIface: "",
  wanIface: "",
  subnet: "192.168.88.0/24",
  hotspotName: "",
});
const tearing = ref(false);
const teardownResult = ref<TeardownResult | null>(null);
const showTeardownConfirm = ref(false);

const branding = ref({ currency: "" });
const brandingSaving = ref(false);
const brandingError = ref("");
const brandingSaved = ref(false);

const effectiveCurrency = computed(() => branding.value.currency);

const bridgeAndWlan = computed(() =>
  interfaces.value.filter(
    (i) => i.type === "bridge" || i.type === "wlan" || i.type === "vlan",
  ),
);

const derived = computed(() => {
  const m = form.value.subnet.trim().match(/^(\d+\.\d+\.\d+)\.(\d+)\/(\d+)$/);
  if (!m) return null;
  const base3 = m[1];
  const prefix = parseInt(m[3]);
  if (prefix < 1 || prefix > 30) return null;
  return {
    gateway: `${base3}.1`,
    poolStart: `${base3}.2`,
    poolEnd: `${base3}.254`,
  };
});

const canSubmit = computed(
  () =>
    form.value.lanIface.trim() !== "" &&
    form.value.wanIface.trim() !== "" &&
    derived.value !== null,
);

async function loadBranding() {
  if (!store.activeId) return;
  try {
    const s = await getHotspotSettings(store.activeId);
    branding.value.currency = s.currency ?? "";
  } catch {
    // non-fatal — branding stays blank
  }
}

async function saveBranding() {
  if (!store.activeId) return;
  brandingSaving.value = true;
  brandingError.value = "";
  brandingSaved.value = false;
  try {
    const existing = await getHotspotSettings(store.activeId);
    await putHotspotSettings(store.activeId, {
      ...existing,
      currency: effectiveCurrency.value,
    });
    brandingSaved.value = true;
    setTimeout(() => {
      brandingSaved.value = false;
    }, 3000);
  } catch (e: any) {
    brandingError.value = friendlyError(e, "Failed to save");
  } finally {
    brandingSaving.value = false;
  }
}

async function runPreflight() {
  if (!store.activeId) return;
  phase.value = "loading";
  preflightError.value = "";
  try {
    const [pf] = await Promise.all([
      hotspotPreflight(store.activeId),
      loadBranding(),
    ]);
    preflight.value = pf;
    interfaces.value = pf.interfaces;
    if (pf.hotspotExists) {
      phase.value = "blocked";
      return;
    }
    const bridge = interfaces.value.find((i) => i.type === "bridge");
    const wlan = interfaces.value.find((i) => i.type === "wlan");
    const ether = interfaces.value.find((i) => i.type === "ether");
    form.value.lanIface = bridge?.name ?? wlan?.name ?? "";
    form.value.wanIface =
      interfaces.value.find(
        (i) => i.type === "ether" && i.name !== form.value.lanIface,
      )?.name ??
      ether?.name ??
      "";
    phase.value = "form";
  } catch (e: any) {
    preflightError.value = friendlyError(e, "Connection failed");
    phase.value = "form";
  }
}

async function runSetup() {
  if (!store.activeId || !canSubmit.value) return;
  phase.value = "running";
  try {
    result.value = await setupHotspot(store.activeId, {
      lanIface: form.value.lanIface,
      wanIface: form.value.wanIface,
      subnet: form.value.subnet,
      hotspotName: form.value.hotspotName,
    });
    if (effectiveCurrency.value) {
      await saveBranding();
    }
  } catch (e: any) {
    result.value = {
      success: false,
      steps: [
        {
          name: "Setup request",
          ok: false,
          error: friendlyError(e, "Request failed"),
        },
      ],
    };
  }
  phase.value = "result";
}

function reset() {
  result.value = null;
  runPreflight();
}

function runTeardown() {
  showTeardownConfirm.value = true;
}

async function confirmTeardown() {
  if (!store.activeId) return;
  showTeardownConfirm.value = false;
  tearing.value = true;
  teardownResult.value = null;
  try {
    teardownResult.value = await teardownHotspot(store.activeId);
    await runPreflight();
  } catch (e: any) {
    teardownResult.value = {
      steps: [
        {
          name: "Teardown",
          ok: false,
          error: friendlyError(e, "Reset failed"),
        },
      ],
    };
  } finally {
    tearing.value = false;
  }
}

watch(
  () => store.activeId,
  (id) => {
    if (id) runPreflight();
  },
  { immediate: true },
);
</script>
