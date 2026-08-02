<template>
  <PageLayout title="Routers" subtitle="Manage your MikroTik connections">
    <template #actions>
      <button class="btn btn-ghost" :disabled="scanning" @click="scan">
        <MagnifyingGlassIcon
          class="size-3.5"
          :class="{ 'animate-pulse': scanning }"
        />
        {{ scanning ? "Scanning…" : "Scan network" }}
      </button>
      <button class="btn btn-primary" @click="showAdd = true">
        <PlusIcon class="size-3.5" />
        Add router
      </button>
    </template>

    <!-- Scan error -->
    <div
      v-if="scanError"
      class="flex items-center gap-2 p-3 rounded-lg text-xs border text-red bg-red/8 border-red/20"
    >
      <ExclamationTriangleIcon class="size-4 shrink-0" />
      {{ scanError }}
    </div>

    <!-- Scanning spinner -->
    <div
      v-if="scanning"
      class="border border-dashed border-border rounded-xl py-16 text-center"
    >
      <span
        class="inline-block size-5 border-2 border-border border-t-text-secondary rounded-full animate-spin mb-3"
      />
      <p class="text-xs text-text-muted">Scanning network… (3 s)</p>
    </div>

    <!-- Discovered devices -->
    <div v-else-if="discovered.length > 0" class="mb-6">
      <p class="text-sm font-medium text-text-secondary mb-3">
        {{ discovered.length }} device{{ discovered.length > 1 ? "s" : "" }}
        found on network
      </p>
      <div class="grid grid-cols-3 gap-3">
        <div
          v-for="d in discovered"
          :key="d.mac || d.ip"
          class="flex flex-col gap-3 p-4 border border-dashed border-border rounded-xl bg-surface"
        >
          <div class="flex items-center gap-2">
            <ServerIcon class="size-4 text-text-secondary shrink-0" />
            <span class="text-sm font-semibold text-text-primary truncate">{{
              d.identity || d.board || "MikroTik"
            }}</span>
          </div>
          <p class="text-xs font-mono text-text-secondary">
            {{ d.ip && d.ip !== "0.0.0.0" ? d.ip : "No IP yet" }}
            <span v-if="d.mac" class="text-text-muted"> · {{ d.mac }}</span>
          </p>
          <p v-if="d.board || d.version" class="text-xs text-text-muted">
            {{ [d.board, d.version].filter(Boolean).join(" · ") }}
          </p>
          <p v-if="!d.ip || d.ip === '0.0.0.0'" class="text-xs text-amber">
            Factory reset — assign an IP via WinBox first
          </p>
          <button
            class="btn btn-ghost self-start mt-auto"
            :class="
              isAlreadySaved(d.ip)
                ? 'cursor-default'
                : !d.ip || d.ip === '0.0.0.0'
                  ? 'opacity-40 cursor-not-allowed'
                  : ''
            "
            :style="
              isAlreadySaved(d.ip)
                ? 'color: var(--color-green); border-color: rgba(76,195,138,0.3)'
                : ''
            "
            :disabled="isAlreadySaved(d.ip) || !d.ip || d.ip === '0.0.0.0'"
            @click="openAdd(d)"
          >
            <CheckIcon v-if="isAlreadySaved(d.ip)" class="size-3.5" />
            <PlusIcon v-else class="size-3.5" />
            {{ isAlreadySaved(d.ip) ? "Saved" : "Add" }}
          </button>
        </div>
      </div>
    </div>

    <!-- No devices after scan -->
    <div
      v-else-if="scanned && discovered.length === 0"
      class="border border-dashed border-border rounded-xl py-12 text-center mb-6"
    >
      <ServerIcon class="size-6 text-text-muted mx-auto mb-2" />
      <p class="text-xs font-medium text-text-secondary">
        No MikroTik devices found
      </p>
      <p class="text-xs text-text-muted mt-0.5">
        Make sure the router is on the same subnet
      </p>
    </div>

    <div
      v-if="store.routers.length > 0"
      class="grid grid-cols-[repeat(auto-fill,minmax(260px,1fr))] gap-5"
    >
      <div
        v-for="r in store.routers"
        :key="r.id"
        class="flex flex-col rounded-xl border transition-all"
        :class="
          store.activeId === r.id
            ? 'border-primary bg-surface shadow-[0_0_0_2px_var(--color-primary)]'
            : 'border-border bg-surface hover:border-muted'
        "
      >
        <div class="flex items-center justify-between pt-3.5 px-3.5">
          <div>
            <span
              v-if="testResults[r.id] === 'ok'"
              class="inline-flex items-center gap-1 text-xs px-2 py-0.5 rounded-full text-green bg-green/10"
            >
              <CheckCircleIcon class="size-3" /> Reachable
            </span>
            <span
              v-else-if="testResults[r.id] === 'error'"
              class="inline-flex items-center gap-1 text-xs px-2 py-0.5 rounded-full text-red bg-red/10"
            >
              <XCircleIcon class="size-3" /> Unreachable
            </span>
          </div>
          <button
            class="btn btn-sm btn-ghost"
            :title="testing[r.id] ? 'Testing…' : 'Test connection'"
            :disabled="testing[r.id]"
            @click.stop="test(r.id)"
          >
            <ArrowPathIcon
              class="size-3.5"
              :class="{ 'animate-spin': testing[r.id] }"
            />
            {{ testing[r.id] ? "Pinging router..." : "Ping" }}
          </button>
        </div>

        <!-- Card body -->
        <div class="flex flex-col items-center gap-2 py-1.5">
          <!-- Router illustration -->
          <RouterArt
            :board-name="r.name"
            :size="96"
            :power-led="
              testResults[r.id] === 'ok'
                ? 'var(--color-green)'
                : testResults[r.id] === 'error'
                  ? 'var(--color-red)'
                  : 'var(--color-border)'
            "
            wifi-led="var(--color-border)"
            wan-led="var(--color-border)"
          />

          <!-- Name + active badge -->
          <div class="flex flex-wrap items-center justify-center gap-1.5 mt-1">
            <span class="font-semibold text-text-primary">{{ r.name }}</span>
            <span
              v-if="store.activeId === r.id"
              class="text-xs px-1.5 py-0.5 rounded font-medium"
              style="
                color: var(--color-primary-fg);
                background: var(--color-primary);
              "
              >Active</span
            >
          </div>

          <!-- Connection info -->
          <p
            class="text-xs font-mono text-text-secondary text-center leading-relaxed"
          >
            {{ r.host }}:{{ r.port }}<br />
            <span v-if="r.useTls"> · TLS</span>
          </p>

          <!-- Hotspot info -->
          <div
            v-if="
              r.hotspotSettings?.hotspotName ||
              r.hotspotSettings?.dnsName ||
              r.hotspotSettings?.currency
            "
            class="flex flex-wrap items-center justify-center gap-x-2 gap-y-0.5 text-xs text-text-secondary font-medium px-2"
          >
            <span
              v-if="r.hotspotSettings?.hotspotName"
              class="inline-flex items-center gap-1"
            >
              <WifiIcon class="size-4" />{{ r.hotspotSettings.hotspotName }}
              -
            </span>
            <span v-if="r.hotspotSettings?.dnsName">{{
              r.hotspotSettings.dnsName
            }} -</span>
            <span
              v-if="r.hotspotSettings?.currency"
              class="inline-flex items-center gap-1"
            >
              <BanknotesIcon class="size-4" />{{ r.hotspotSettings.currency }}
            </span>
          </div>
        </div>

        <!-- Divider -->
        <div class="h-px bg-border" />

        <!-- Action row -->
        <div class="flex flex-wrap gap-2 py-3 px-3.5">
          <button
            class="btn flex-1 min-w-0 justify-center whitespace-normal text-center"
            :class="
              store.activeId === r.id
                ? 'border-transparent opacity-40 cursor-default bg-transparent text-text-muted'
                : 'btn-primary'
            "
            :disabled="store.activeId === r.id"
            @click="selectRouter(r.id)"
          >
            <BoltIcon class="size-3.5 shrink-0" />
            {{ store.activeId === r.id ? "In use" : "Use" }}
          </button>

          <button
            class="btn btn-ghost flex-1 min-w-0 justify-center whitespace-normal text-center"
            @click="openEdit(r)"
          >
            <PencilSquareIcon class="size-3.5 shrink-0" />

            Edit
          </button>

          <button
            class="btn border-transparent bg-red-10 text-red hover:bg-red/50 hover:text-white flex-1 min-w-0 justify-center whitespace-normal text-center"
            @click="remove(r.id)"
          >
            <TrashIcon class="size-4 shrink-0" />
            Delete
          </button>
        </div>
      </div>
    </div>

    <EmptyState
      v-else-if="!scanning && !scanned"
      bordered
      size="lg"
      title="No routers configured"
      message="Click Scan to find devices on your network, or add one manually"
    />

    <AddRouterDialog
      v-model:open="showAdd"
      :prefill="prefill"
      @added="onAdded"
    />
    <EditRouterDialog
      v-model:open="showEdit"
      :router="editTarget"
      @saved="store.load()"
    />
  </PageLayout>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useRouter, useRoute } from "vue-router";
import {
  PlusIcon,
  ServerIcon,
  ArrowPathIcon,
  BoltIcon,
  TrashIcon,
  CheckCircleIcon,
  XCircleIcon,
  MagnifyingGlassIcon,
  CheckIcon,
  ExclamationTriangleIcon,
  PencilSquareIcon,
  WifiIcon,
  BanknotesIcon,
} from "@heroicons/vue/24/outline";
import { useRoutersStore } from "@/stores/routers";
import {
  testRouter,
  discoverRouters,
  type DiscoveredDevice,
  type RouterProfile,
} from "@/api";
import AddRouterDialog from "@/components/AddRouterDialog.vue";
import EditRouterDialog from "@/components/EditRouterDialog.vue";
import PageLayout from "@/components/PageLayout.vue";
import RouterArt from "@/components/router-art/RouterArt.vue";
import EmptyState from "@/components/EmptyState.vue";

const store = useRoutersStore();
const router = useRouter();
const route = useRoute();

onMounted(() => {
  if (route.query.scan === "1") scan();
});

const testing = ref<Record<string, boolean>>({});
const testResults = ref<Record<string, "ok" | "error">>({});

const scanning = ref(false);
const scanned = ref(false);
const scanError = ref("");
const discovered = ref<DiscoveredDevice[]>([]);
const addedIps = ref(new Set<string>());

const showAdd = ref(false);
const prefill = ref<Partial<DiscoveredDevice>>({});

const showEdit = ref(false);
const editTarget = ref<RouterProfile | null>(null);

function openEdit(r: RouterProfile) {
  editTarget.value = r;
  showEdit.value = true;
}

async function scan() {
  scanning.value = true;
  scanned.value = false;
  scanError.value = "";
  discovered.value = [];
  try {
    discovered.value = await discoverRouters();
    scanned.value = true;
  } catch (e: any) {
    scanError.value =
      e?.response?.data?.error ?? "Scan failed — check network/firewall";
  } finally {
    scanning.value = false;
  }
}

function isAlreadySaved(ip: string): boolean {
  return store.routers.some((r) => r.host === ip) || addedIps.value.has(ip);
}

function openAdd(d: DiscoveredDevice) {
  prefill.value = d;
  showAdd.value = true;
}

function onAdded(ip: string) {
  addedIps.value = new Set([...addedIps.value, ip]);
}

async function test(id: string) {
  testing.value[id] = true;
  try {
    await testRouter(id);
    testResults.value[id] = "ok";
  } catch {
    testResults.value[id] = "error";
  } finally {
    testing.value[id] = false;
  }
}

function selectRouter(id: string) {
  store.select(id);
  router.push("/dashboard");
}

async function remove(id: string) {
  if (!confirm("Remove this router?")) return;
  await store.remove(id);
}
</script>
