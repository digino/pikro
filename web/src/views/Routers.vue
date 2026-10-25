<template>
  <main class="flex flex-col gap-7 px-4 sm:px-6 lg:px-10 pt-7 pb-14">
    <!-- Discovery hero — on amber, contents use primary-fg so they stay dark in both themes -->
    <section
      class="card bg-primary text-primary-fg p-[clamp(24px,4vw,44px)] flex flex-wrap items-center gap-7 overflow-hidden"
    >
      <div class="flex-[1_1_380px] flex flex-col gap-4">
        <span
          v-if="scanning"
          class="chip self-start bg-primary-fg text-primary border-primary-fg"
        >
          <span class="size-2 rounded-full bg-primary motion-safe:animate-pulse" />
          {{ t('routers.scanning') }}
        </span>
        <h1 class="type-display m-0 text-[clamp(30px,3.5vw,48px)] leading-none">
          {{ t('routers.heroTitle') }}
        </h1>
        <p class="m-0 max-w-120 text-[16px] leading-normal">{{ t('routers.heroBody') }}</p>
        <div class="flex flex-wrap gap-3">
          <button type="button" class="btn btn-ink" @click="openAdd({})">
            <PlusIcon class="size-4" />
            {{ t('routers.addByIp') }}
          </button>
          <button type="button" class="btn" :disabled="scanning" @click="scan">
            <MagnifyingGlassIcon class="size-4" />
            {{ scanned ? t('routers.scanAgain') : t('routers.scan') }}
          </button>
        </div>
      </div>
      <div
        aria-hidden="true"
        class="relative flex-[0_1_320px] h-70 mx-auto flex items-center justify-center"
      >
        <div
          class="absolute size-70 rounded-full border-2 border-dashed border-primary-fg opacity-35"
          :class="{ 'motion-safe:animate-pulse': scanning }"
        />
        <div
          class="absolute size-50 rounded-full border-2 border-dashed border-primary-fg opacity-60"
          :class="{ 'motion-safe:animate-pulse': scanning }"
        />
        <div class="absolute size-32.5 rounded-full bg-white border-2 border-primary-fg" />
        <img :src="inkIcon('routers')" alt="" width="190" height="190" class="relative -rotate-6" />
      </div>
    </section>

    <!-- Scan error -->
    <div
      v-if="scanError"
      role="alert"
      class="flex items-start gap-2.5 p-4 rounded-control border-2 border-danger-text bg-danger-wash text-danger-text text-sm"
    >
      <ExclamationTriangleIcon class="size-5 shrink-0" />
      {{ scanError }}
    </div>

    <!-- Discovered devices -->
    <section v-if="scanning || scanned" class="flex flex-col gap-3.5">
      <div class="flex flex-wrap items-baseline gap-3">
        <h2 class="type-display m-0 text-2xl">{{ t('routers.foundTitle') }}</h2>
        <span v-if="!scanning" class="font-mono text-sm text-text-muted">
          {{ t('routers.foundCount', discovered.length) }}
        </span>
      </div>

      <div v-if="scanning" class="card border-dashed shadow-none py-12 text-center">
        <span class="spinner" />
        <p class="m-0 mt-3 text-sm text-text-muted">{{ t('routers.scanningHint') }}</p>
      </div>

      <div v-else-if="discovered.length === 0" class="card border-dashed shadow-none p-8 text-center">
        <img :src="inkIcon('routers')" alt="" width="64" height="64" class="mx-auto" />
        <p class="m-0 mt-3 font-semibold">{{ t('routers.noneTitle') }}</p>
        <p class="m-0 mt-1 mx-auto max-w-120 text-sm text-text-muted">{{ t('routers.noneBody') }}</p>
      </div>

      <div v-else class="grid grid-cols-[repeat(auto-fill,minmax(280px,1fr))] gap-5">
        <article
          v-for="d in discovered"
          :key="d.mac || d.ip"
          class="card p-5 flex flex-col gap-3.5"
        >
          <div class="flex items-center gap-3">
            <img :src="inkIcon('routers')" alt="" width="60" height="60" class="shrink-0" />
            <div class="flex-1 min-w-0">
              <h3 class="type-display m-0 text-xl truncate">{{ deviceName(d) }}</h3>
              <div v-if="d.board || d.version" class="text-sm text-text-secondary">
                {{ [d.board, d.version && `RouterOS ${d.version}`].filter(Boolean).join(' · ') }}
              </div>
            </div>
          </div>
          <dl class="m-0 grid grid-cols-[auto_1fr] gap-x-3.5 gap-y-1.5 text-sm">
            <dt class="eyebrow self-center">{{ t('routers.ip') }}</dt>
            <dd class="m-0 font-mono">{{ hasIp(d) ? d.ip : '—' }}</dd>
            <dt class="eyebrow self-center">{{ t('routers.mac') }}</dt>
            <dd class="m-0 font-mono">{{ d.mac || '—' }}</dd>
          </dl>
          <p v-if="!hasIp(d)" class="m-0 text-sm text-text-secondary">{{ t('routers.noIpHint') }}</p>
          <div class="mt-auto flex items-center justify-between gap-2.5">
            <span v-if="savedRouter(d)" class="chip bg-ink text-paper">{{ t('routers.connected') }}</span>
            <span v-else-if="hasIp(d)" class="chip bg-online-wash">
              <span class="size-2 rounded-full bg-online" />{{ t('routers.ready') }}
            </span>
            <span v-else class="chip border-dashed">{{ t('routers.noIp') }}</span>

            <button
              v-if="savedRouter(d)"
              type="button"
              class="btn btn-sm"
              @click="selectRouter(savedRouter(d)!.id)"
            >
              {{ t('routers.open') }}
            </button>
            <button
              v-else
              type="button"
              class="btn btn-sm"
              :disabled="!hasIp(d)"
              :aria-label="t('routers.connectDevice', { name: deviceName(d) })"
              @click="openAdd(d)"
            >
              {{ t('routers.connect') }}
            </button>
          </div>
        </article>
      </div>
    </section>

    <!-- Saved routers -->
    <section class="flex flex-col gap-3.5">
      <h2 class="type-display m-0 text-2xl">{{ t('routers.yoursTitle') }}</h2>

      <div v-if="store.routers.length === 0" class="card border-dashed shadow-none p-8 text-center">
        <p class="m-0 font-semibold">{{ t('routers.noneYetTitle') }}</p>
        <p class="m-0 mt-1 text-sm text-text-muted">{{ t('routers.noneYetBody') }}</p>
      </div>

      <div v-else class="grid grid-cols-[repeat(auto-fill,minmax(280px,1fr))] gap-5">
        <article
          v-for="r in store.routers"
          :key="r.id"
          class="card p-5 flex flex-col gap-3.5"
          :class="{ 'bg-primary-wash': store.activeId === r.id }"
        >
          <div class="flex items-center gap-3">
            <img :src="inkIcon('routers')" alt="" width="56" height="56" class="shrink-0" />
            <div class="flex-1 min-w-0">
              <div class="flex flex-wrap items-center gap-2">
                <h3 class="type-display m-0 text-xl truncate">{{ r.name }}</h3>
                <span v-if="store.activeId === r.id" class="chip bg-primary text-primary-fg">
                  {{ t('routers.active') }}
                </span>
              </div>
              <div class="font-mono text-sm text-text-secondary">
                {{ r.host }}:{{ r.port }}<template v-if="r.useTls"> · {{ t('routers.apiSsl') }}</template>
              </div>
            </div>
          </div>

          <div v-if="hotspotLine(r)" class="flex items-center gap-1.5 text-sm text-text-secondary">
            <WifiIcon class="size-4 shrink-0" />
            <span class="truncate">{{ hotspotLine(r) }}</span>
          </div>

          <div class="flex flex-wrap items-center gap-2">
            <span v-if="testResults[r.id] === 'ok'" class="chip bg-online-wash">
              <span class="size-2 rounded-full bg-online" />{{ t('routers.reachable') }}
            </span>
            <span
              v-else-if="testResults[r.id] === 'error'"
              class="chip bg-danger-wash text-danger-text border-danger-text"
            >
              <span class="size-2 rotate-45 bg-danger" />{{ t('routers.unreachable') }}
            </span>
            <button
              type="button"
              class="btn btn-sm ml-auto"
              :disabled="testing[r.id]"
              :aria-label="t('routers.pingRouter', { name: r.name })"
              @click="test(r.id)"
            >
              <ArrowPathIcon class="size-4" :class="{ 'animate-spin': testing[r.id] }" />
              {{ testing[r.id] ? t('routers.pinging') : t('routers.ping') }}
            </button>
          </div>

          <div class="mt-auto pt-3.5 border-t border-hairline flex flex-wrap gap-2">
            <button
              type="button"
              class="btn btn-ink flex-1"
              :disabled="store.activeId === r.id"
              @click="selectRouter(r.id)"
            >
              {{ store.activeId === r.id ? t('routers.inUse') : t('routers.use') }}
            </button>
            <button
              type="button"
              class="btn flex-1"
              :aria-label="t('routers.editRouter', { name: r.name })"
              @click="openEdit(r)"
            >
              <PencilSquareIcon class="size-4" />
              {{ t('routers.edit') }}
            </button>
            <button
              type="button"
              class="btn btn-danger"
              :aria-label="t('routers.deleteRouter', { name: r.name })"
              :title="t('routers.delete')"
              @click="remove(r.id)"
            >
              <TrashIcon class="size-4" />
            </button>
          </div>
        </article>
      </div>
    </section>

    <AddRouterDialog v-model:open="showAdd" :prefill="prefill" />
    <EditRouterDialog
      v-model:open="showEdit"
      :router="editTarget"
      @saved="store.load()"
    />
  </main>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useRouter, useRoute } from "vue-router";
import { useI18n } from "vue-i18n";
import {
  PlusIcon,
  ArrowPathIcon,
  TrashIcon,
  MagnifyingGlassIcon,
  ExclamationTriangleIcon,
  PencilSquareIcon,
  WifiIcon,
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
import { inkIcon } from "@/utils/inkIcons";

const { t } = useI18n();
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

// A factory-reset router announces itself before it has an IP address.
function hasIp(d: DiscoveredDevice): boolean {
  return !!d.ip && d.ip !== "0.0.0.0";
}

function deviceName(d: DiscoveredDevice): string {
  return d.identity || d.board || "MikroTik";
}

function savedRouter(d: DiscoveredDevice): RouterProfile | undefined {
  return store.routers.find((r) => r.host === d.ip);
}

function hotspotLine(r: RouterProfile): string {
  const h = r.hotspotSettings;
  return [h?.hotspotName, h?.dnsName, h?.currency].filter(Boolean).join(" · ");
}

function openAdd(d: Partial<DiscoveredDevice>) {
  prefill.value = d;
  showAdd.value = true;
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
  if (!confirm(t("routers.confirmDelete"))) return;
  await store.remove(id);
}
</script>
