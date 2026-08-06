<template>
  <AppDialog :open="open" title="Hotspot Setup" @update:open="onClose">
    <div
      v-if="phase === 'loading'"
      class="flex flex-col items-center py-8 gap-3"
    >
      <span
        class="inline-block size-6 border-2 rounded-full animate-spin border-border border-t-text-secondary"
      />
      <p class="text-sm text-text-secondary">Checking router configuration…</p>
    </div>

    <!-- Phase: form -->
    <template v-else-if="phase === 'form'">
      <!-- Safety notice -->
      <div
        class="flex items-start gap-2 p-3 mb-4 rounded-lg border text-xs bg-amber/8 border-amber/20 text-amber"
      >
        <ExclamationTriangleIcon class="size-4 shrink-0 mt-0.5" />
        <span
          >Make sure no hotspot is already configured on this router. Running
          setup on an existing hotspot configuration may cause conflicts.</span
        >
      </div>

      <!-- Preflight connection error -->
      <div
        v-if="preflightError"
        class="flex items-start gap-2 p-3 mb-4 rounded-lg border text-xs bg-red/8 border-red/20 text-red"
      >
        <ExclamationTriangleIcon class="size-4 shrink-0 mt-0.5" />
        <span
          >Could not read interfaces: {{ preflightError }}. Enter interface
          names manually.</span
        >
      </div>

      <div class="space-y-4">
        <!-- LAN interface -->
        <label class="flex flex-col gap-1">
          <span class="text-sm font-medium text-red"
            >LAN / Hotspot interface <span>*</span></span
          >
          <select
            v-if="interfaces.length"
            v-model="form.lanIface"
            class="w-full px-3 py-2 text-sm border border-border rounded-lg focus:outline-2 focus:outline-accent transition-colors bg-base text-text-primary"
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
            class="w-full px-3 py-2 text-sm border border-border rounded-lg focus:outline-2 focus:outline-accent transition-colors bg-base text-text-primary"
            placeholder="e.g. bridge1 or ether2"
            required
          />
          <p class="text-xs text-text-muted">
            The interface clients connect to (usually a bridge or wlan)
          </p>
        </label>

        <!-- WAN interface -->
        <label class="flex flex-col gap-1">
          <span class="text-sm font-medium text-red"
            >WAN interface <span>*</span></span
          >
          <select
            v-if="interfaces.length"
            v-model="form.wanIface"
            class="w-full px-3 py-2 text-sm border border-border rounded-lg focus:outline-2 focus:outline-accent transition-colors bg-base text-text-primary"
          >
            <option value="" disabled>Select interface…</option>
            <option v-for="i in interfaces" :key="i.name" :value="i.name">
              {{ i.name }} [{{ i.type }}]{{ !i.running ? " — down" : "" }}
            </option>
          </select>
          <input
            v-else
            v-model="form.wanIface"
            class="w-full px-3 py-2 text-sm border border-border rounded-lg focus:outline-2 focus:outline-accent transition-colors bg-base text-text-primary"
            placeholder="e.g. ether1"
            required
          />
          <p class="text-xs text-text-muted">
            The interface connected to the internet (usually ether1)
          </p>
        </label>

        <!-- Subnet -->
        <label class="flex flex-col gap-1">
          <span class="text-sm font-medium text-red"
            >Subnet (CIDR) <span>*</span></span
          >
          <input
            v-model="form.subnet"
            class="w-full px-3 py-2 text-sm border border-border rounded-lg font-mono focus:outline-2 focus:outline-accent transition-colors bg-base text-text-primary"
            placeholder="192.168.88.0/24"
            required
          />
        </label>

        <!-- Derived values preview -->
        <div
          v-if="derived"
          class="border border-border rounded-lg p-3 text-xs bg-surface"
        >
          <div class="grid grid-cols-2 gap-x-4 gap-y-1">
            <span class="text-text-muted">Router IP (gateway)</span>
            <span class="font-mono text-text-secondary">{{
              derived.gateway
            }}</span>
            <span class="text-text-muted">DHCP pool</span>
            <span class="font-mono text-text-secondary"
              >{{ derived.poolStart }} – {{ derived.poolEnd }}</span
            >
            <span class="text-text-muted">Hotspot profile</span>
            <span class="font-mono text-text-secondary">pikro-profile</span>
          </div>
        </div>

        <!-- What will be configured -->
        <details class="text-xs text-text-muted cursor-pointer">
          <summary class="hover:text-text-secondary transition-colors">
            What will be configured
          </summary>
          <ul class="mt-2 space-y-1 pl-2 text-text-muted">
            <li>
              • Assign {{ derived?.gateway ?? "…" }} to
              {{ form.lanIface || "(LAN)" }}
            </li>
            <li>• Create IP pool hotspot-pool</li>
            <li>• Create DHCP server + network</li>
            <li>• Create hotspot server profile with Pikro login page</li>
            <li>• Enable hotspot on {{ form.lanIface || "(LAN)" }}</li>
            <li>• Add NAT masquerade on {{ form.wanIface || "(WAN)" }}</li>
            <li>• Enable DNS remote requests</li>
          </ul>
        </details>
      </div>

      <div class="flex justify-end gap-2 pt-4 mt-2 border-t border-border">
        <button type="button" class="btn btn-ghost" @click="onClose(false)">
          Cancel
        </button>
        <button
          type="button"
          class="btn btn-primary"
          :disabled="!canSubmit"
          @click="runSetup"
        >
          Run Setup
        </button>
      </div>
    </template>

    <!-- Phase: blocked -->
    <div v-else-if="phase === 'blocked'" class="space-y-4">
      <div
        class="flex items-start gap-3 p-4 border rounded-xl bg-red/8 border-red/20"
      >
        <ExclamationTriangleIcon class="size-5 shrink-0 mt-0.5 text-red" />
        <div>
          <p class="text-sm font-medium text-red">Hotspot already configured</p>
          <p class="text-sm mt-1 text-text-secondary">
            A hotspot is already enabled on interface
            <span class="font-mono font-semibold">{{
              preflight?.hotspotOnIface
            }}</span
            >. Remove the existing hotspot from RouterOS before running setup.
          </p>
        </div>
      </div>
      <div class="flex justify-end">
        <button type="button" class="btn btn-ghost" @click="onClose(false)">
          Close
        </button>
      </div>
    </div>

    <!-- Phase: running -->
    <div
      v-else-if="phase === 'running'"
      class="flex flex-col items-center py-10 gap-3"
    >
      <span
        class="inline-block size-6 border-2 rounded-full animate-spin border-border border-t-text-secondary"
      />
      <p class="text-sm font-medium text-text-secondary">
        Configuring hotspot…
      </p>
      <p class="text-xs text-text-muted">
        This takes a few seconds. Please wait.
      </p>
    </div>

    <!-- Phase: result -->
    <div v-else-if="phase === 'result' && result" class="space-y-4">
      <div
        class="flex items-center gap-2 p-3 rounded-lg border text-sm font-medium"
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

      <ul class="space-y-1.5">
        <li
          v-for="step in result.steps"
          :key="step.name"
          class="flex items-start gap-2 text-sm"
        >
          <CheckCircleIcon
            v-if="step.ok"
            class="size-3.5 shrink-0 mt-0.5 text-green"
          />
          <MinusCircleIcon
            v-else-if="step.skipped"
            class="size-3.5 shrink-0 mt-0.5 text-text-muted"
          />
          <XCircleIcon v-else class="size-3.5 shrink-0 mt-0.5 text-red" />
          <div class="flex-1 min-w-0">
            <span
              :class="
                step.ok
                  ? 'text-text-secondary'
                  : step.skipped
                    ? 'text-text-muted'
                    : 'text-red'
              "
            >
              {{ step.name }}
            </span>
            <p v-if="step.error" class="font-mono mt-0.5 break-all text-red">
              {{ step.error }}
            </p>
          </div>
        </li>
      </ul>

      <p
        v-if="!result.success"
        class="text-sm text-text-muted border border-border rounded-lg p-3 bg-surface"
      >
        Fix the failed steps manually in your router's configuration, then try
        again or continue using the app.
      </p>

      <div class="flex justify-end">
        <button type="button" class="btn btn-primary" @click="onClose(false)">
          Done
        </button>
      </div>
    </div>
  </AppDialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from "vue";
import {
  ExclamationTriangleIcon,
  CheckCircleIcon,
  XCircleIcon,
  MinusCircleIcon,
} from "@heroicons/vue/24/outline";
import {
  hotspotPreflight,
  setupHotspot,
  type PreflightResult,
  type SetupResult,
  type InterfaceInfo,
} from "@/api";
import { useRoutersStore } from "@/stores/routers";
import AppDialog from "@/components/AppDialog.vue";

const props = defineProps<{ open: boolean }>();
const emit = defineEmits<{ "update:open": [value: boolean]; done: [] }>();

const store = useRoutersStore();

type Phase = "loading" | "form" | "blocked" | "running" | "result";
const phase = ref<Phase>("loading");
const preflight = ref<PreflightResult | null>(null);
const preflightError = ref("");
const result = ref<SetupResult | null>(null);

const interfaces = ref<InterfaceInfo[]>([]);
const form = ref({ lanIface: "", wanIface: "", subnet: "192.168.88.0/24" });

const bridgeAndWlan = computed(() =>
  interfaces.value.filter(
    (i) => i.type === "bridge" || i.type === "wlan" || i.type === "vlan",
  ),
);

const derived = computed(() => {
  const cidr = form.value.subnet.trim();
  const match = cidr.match(/^(\d+\.\d+\.\d+)\.(\d+)\/(\d+)$/);
  if (!match) return null;
  const base3 = match[1];
  const prefix = parseInt(match[3]);
  if (prefix < 1 || prefix > 30) return null;
  const totalHosts = (1 << (32 - prefix)) - 2;
  if (totalHosts < 2) return null;
  const gateway = `${base3}.1`;
  const poolStart = `${base3}.2`;
  const broadcastLast =
    255 - (255 & ((1 << (32 - prefix)) - 1)) + ((1 << (32 - prefix)) - 1);
  const poolEnd = `${base3}.${Math.min(broadcastLast - 1, 254)}`;
  return { gateway, poolStart, poolEnd };
});

const canSubmit = computed(
  () =>
    form.value.lanIface.trim() !== "" &&
    form.value.wanIface.trim() !== "" &&
    derived.value !== null,
);

watch(
  () => props.open,
  async (open) => {
    if (!open || !store.activeId) return;
    resetState();
    phase.value = "loading";
    try {
      preflight.value = await hotspotPreflight(store.activeId);
      interfaces.value = preflight.value.interfaces;
      if (preflight.value.hotspotExists) {
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
      preflightError.value =
        e?.response?.data?.error ?? e?.message ?? "Connection failed";
      phase.value = "form";
    }
  },
);

async function runSetup() {
  if (!store.activeId || !canSubmit.value) return;
  phase.value = "running";
  try {
    result.value = await setupHotspot(store.activeId, {
      lanIface: form.value.lanIface,
      wanIface: form.value.wanIface,
      subnet: form.value.subnet,
      hotspotName: "",
    });
    phase.value = "result";
    if (result.value.success) emit("done");
  } catch (e: any) {
    result.value = {
      success: false,
      steps: [
        {
          name: "Setup request",
          ok: false,
          error: e?.response?.data?.error ?? e?.message ?? "Request failed",
        },
      ],
    };
    phase.value = "result";
  }
}

function onClose(val: boolean) {
  if (phase.value === "running") return;
  emit("update:open", val);
}

function resetState() {
  preflight.value = null;
  preflightError.value = "";
  result.value = null;
  interfaces.value = [];
  form.value = { lanIface: "", wanIface: "", subnet: "192.168.88.0/24" };
  phase.value = "loading";
}
</script>
