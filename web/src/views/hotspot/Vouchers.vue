<template>
  <PageLayout title="Hotspot" subtitle="Vouchers">
    <template #actions>
      <button class="btn btn-primary" @click="openBatch()">
        <TicketIcon class="size-3.5" />
        Generate vouchers
      </button>
    </template>

    <div v-if="loading" class="flex justify-center py-10">
      <span class="spinner" />
    </div>

    <div
      v-else-if="error"
      class="flex items-center gap-2 p-4 border rounded-xl text-sm bg-red/8 border-red/20 text-red"
    >
      <ExclamationTriangleIcon class="size-4 shrink-0" />
      {{ error }}
      <button class="ml-auto text-xs underline" @click="load">Retry</button>
    </div>

    <template v-else>
      <!-- Profile cards -->
      <div v-if="profiles.length > 0" class="space-y-3">
        <p class="text-sm text-text-secondary">
          Pick a profile to quickly generate vouchers for it.
        </p>
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
          <button
            v-for="p in profiles"
            :key="p['.id']"
            class="border border-border rounded-xl p-4 flex flex-col gap-3 text-left transition-colors hover:border-muted bg-surface cursor-pointer"
            @click="openBatch(p.name)"
          >
            <div class="flex items-start justify-between gap-2">
              <div>
                <p class="text-sm font-semibold text-text-primary">
                  {{ p.name }}
                </p>
                <p
                  v-if="p['rate-limit']"
                  class="text-sm text-text-secondary font-mono mt-0.5"
                >
                  {{ p["rate-limit"] }}
                </p>
              </div>
              <div>
                <button class="btn btn-sm">
                  Generate
                  <TicketIcon class="size-4" />
                </button>
              </div>
            </div>

            <div class="grid grid-cols-2 gap-x-3 gap-y-1 text-sm">
              <span class="text-text-secondary">Validity</span>
              <span class="text-text-primary">{{
                profileMetas[p.name]?.validity || "—"
              }}</span>
              <span class="text-text-secondary">Price</span>
              <span class="text-text-primary">{{
                profileMetas[p.name]?.price
                  ? `${profileMetas[p.name].price} ${hotspotSettings.currency ?? ""}`
                  : "—"
              }}</span>
            </div>
          </button>
        </div>
      </div>

      <div
        v-else
        class="border border-dashed border-border rounded-xl py-16 text-center"
      >
        <TicketIcon class="size-8 text-text-muted mx-auto mb-3" />
        <p class="text-sm font-medium text-text-secondary">No profiles yet</p>
        <p class="text-sm text-text-muted mt-1">
          Create a profile first to define pricing and validity for vouchers.
        </p>
        <RouterLink to="/hotspot/profiles" class="btn btn-primary btn-sm mt-4">
          <PlusIcon class="size-3.5" />
          New profile
        </RouterLink>
      </div>
    </template>

    <!-- Batch generate dialog -->
    <AppDialog
      :open="showBatch"
      title="Generate Vouchers"
      @update:open="onBatchDialogUpdate"
    >
      <!-- Results screen -->
      <div v-if="batchDone" class="space-y-4">
        <div class="flex items-center gap-2 text-sm text-text-secondary">
          <span class="font-bold text-text-primary">{{
            batchResults.filter((r) => r.ok).length
          }}</span>
          created,
          <span class="font-bold text-red">{{
            batchResults.filter((r) => !r.ok).length
          }}</span>
          failed
        </div>
        <div
          class="max-h-64 overflow-y-auto border border-border rounded-lg text-xs font-mono bg-base"
        >
          <div
            v-for="r in batchResults"
            :key="r.name"
            class="flex items-center gap-3 px-3 py-2 border-b border-border last:border-0"
            :class="r.ok ? 'text-text-secondary' : 'text-red bg-red/5'"
          >
            <span class="flex-1">{{ r.name }}</span>
            <span class="text-text-muted">{{ r.password }}</span>
            <span v-if="!r.ok" class="text-xs ml-auto text-red">{{
              r.error
            }}</span>
          </div>
        </div>
        <div class="flex justify-end gap-2 pt-1 border-t border-border">
          <button
            type="button"
            class="btn btn-ghost"
            @click="showBatch = false"
          >
            Close
          </button>
          <button
            type="button"
            class="btn btn-primary"
            @click="openPrintResults"
          >
            <PrinterIcon class="size-4" /> Print
          </button>
        </div>
      </div>

      <!-- Progress screen -->
      <div v-else-if="batchRunning" class="space-y-4">
        <div
          class="flex items-center justify-between text-xs text-text-secondary"
        >
          <span>Creating vouchers…</span>
          <span>{{ batchProgress }} / {{ batchTotal }}</span>
        </div>
        <div class="w-full rounded-full h-2 border border-border bg-base">
          <div
            class="h-2 rounded-full transition-all bg-text-primary"
            :style="{ width: `${(batchProgress / batchTotal) * 100}%` }"
          />
        </div>
        <p class="text-xs text-text-muted font-mono">{{ batchCurrentName }}</p>
      </div>

      <!-- Wizard -->
      <BatchWizard
        v-else
        ref="wizardRef"
        :profiles="profiles"
        :profile-metas="profileMetas"
        :error="batchError"
        :initial-profile="pendingProfile"
        @cancel="showBatch = false"
        @submit="submitBatch"
      />
    </AppDialog>

    <PrintTemplateDialog
      :open="showPrintDialog"
      :entries="printEntries"
      :default-layout="hotspotSettings.voucher?.layout ?? 'card'"
      :business-name="hotspotSettings.hotspotName ?? ''"
      :show-validity="hotspotSettings.voucher?.showValidity ?? true"
      :show-price="hotspotSettings.voucher?.showPrice ?? true"
      :currency="hotspotSettings.currency ?? ''"
      :profile-metas="profileMetas"
      @update:open="showPrintDialog = $event"
    />
  </PageLayout>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { RouterLink, useRoute, useRouter } from "vue-router";
import BatchWizard, { type BatchConfig } from "@/components/BatchWizard.vue";
import {
  PlusIcon,
  ExclamationTriangleIcon,
  TicketIcon,
  PrinterIcon,
} from "@heroicons/vue/24/outline";
import { useRoutersStore } from "@/stores/routers";
import {
  listHotspotProfiles,
  createHotspotUser,
  getProfileMetas,
  getHotspotSettings,
  type ProfileMeta,
  type HotspotSettings,
} from "@/api";
import { friendlyError } from "@/utils/errors";
import AppDialog from "@/components/AppDialog.vue";
import PageLayout from "@/components/PageLayout.vue";
import PrintTemplateDialog from "@/components/PrintTemplateDialog.vue";
import { type VoucherEntry } from "@/utils/vouchers";

const store = useRoutersStore();
const route = useRoute();
const router = useRouter();

const profiles = ref<Record<string, string>[]>([]);
const profileMetas = ref<Record<string, ProfileMeta>>({});
const hotspotSettings = ref<HotspotSettings>({
  hotspotName: "",
  dnsName: "",
  currency: "",
});
const loading = ref(false);
const error = ref("");

async function load() {
  if (!store.activeId) return;
  loading.value = true;
  error.value = "";
  try {
    const [p, m, s] = await Promise.all([
      listHotspotProfiles(store.activeId),
      getProfileMetas(store.activeId).catch(
        () => ({}) as Record<string, ProfileMeta>,
      ),
      getHotspotSettings(store.activeId).catch(
        () =>
          ({ hotspotName: "", dnsName: "", currency: "" }) as HotspotSettings,
      ),
    ]);
    profiles.value = p;
    profileMetas.value = m;
    hotspotSettings.value = s;

    // A `?profile=` query param (e.g. linked from the Profiles page) opens
    // the wizard pre-filled for that profile, then clears the param so
    // back/forward navigation doesn't keep reopening it.
    const queryProfile = route.query.profile;
    if (typeof queryProfile === "string" && queryProfile) {
      openBatch(queryProfile);
      router.replace({ query: { ...route.query, profile: undefined } });
    }
  } catch (e: any) {
    error.value = friendlyError(e, "Failed to load profiles");
  } finally {
    loading.value = false;
  }
}
load();

const showBatch = ref(false);
const batchRunning = ref(false);
const batchDone = ref(false);
const batchProgress = ref(0);
const batchTotal = ref(0);
const batchCurrentName = ref("");
const batchError = ref("");
const batchResults = ref<
  { name: string; password: string; ok: boolean; error?: string }[]
>([]);
const wizardRef = ref<InstanceType<typeof BatchWizard> | null>(null);
const lastBatchProfile = ref("");
const pendingProfile = ref("");

function generateName(
  charset: string,
  length: number,
  existing: Set<string>,
): string {
  let name = "",
    attempts = 0;
  do {
    name = Array.from(
      { length },
      () => charset[Math.floor(Math.random() * charset.length)],
    ).join("");
    attempts++;
  } while (existing.has(name) && attempts < 1000);
  return name;
}

function generatePassword(
  mode: "same" | "random" | "fixed",
  name: string,
  fixed: string,
  charset: string,
  length: number,
): string {
  if (mode === "same") return name;
  if (mode === "fixed") return fixed;
  return Array.from(
    { length },
    () => charset[Math.floor(Math.random() * charset.length)],
  ).join("");
}

function openBatch(profileName = "") {
  batchDone.value = false;
  batchRunning.value = false;
  batchResults.value = [];
  batchError.value = "";
  pendingProfile.value = profileName;
  showBatch.value = true;
  wizardRef.value?.reset(profileName);
}

function onBatchDialogUpdate(open: boolean) {
  if (!batchRunning.value) showBatch.value = open;
}

async function submitBatch(cfg: BatchConfig) {
  if (!store.activeId) return;
  const charset = [
    cfg.charsLetters ? "abcdefghijklmnopqrstuvwxyz" : "",
    cfg.charsDigits ? "0123456789" : "",
  ].join("");
  if (!charset) {
    batchError.value = "Select at least one character type.";
    return;
  }
  batchRunning.value = true;
  batchTotal.value = cfg.count;
  lastBatchProfile.value = cfg.profile;
  batchProgress.value = 0;
  batchResults.value = [];
  batchError.value = "";
  const mul = cfg.limitBytesTotalUnit === "G" ? 1024 ** 3 : 1024 ** 2;
  const limitBytesTotal = cfg.limitBytesTotalValue
    ? String(cfg.limitBytesTotalValue * mul)
    : "";
  const usedNames = new Set<string>();
  for (let i = 0; i < cfg.count; i++) {
    const name = generateName(charset, cfg.nameLength, usedNames);
    usedNames.add(name);
    const password = generatePassword(
      cfg.passwordMode,
      name,
      cfg.fixedPassword,
      charset,
      cfg.nameLength,
    );
    batchCurrentName.value = name;
    try {
      await createHotspotUser(store.activeId, {
        name,
        password,
        profile: cfg.profile,
        limitUptime: cfg.uptimePreview,
        limitBytesTotal,
        rateLimit: "",
        comment: cfg.comment,
        expiryComment: "",
        price: profileMetas.value[cfg.profile]?.price ?? "",
        currency: hotspotSettings.value.currency ?? "",
      });
      batchResults.value.push({ name, password, ok: true });
    } catch (e: any) {
      batchResults.value.push({
        name,
        password,
        ok: false,
        error: e?.response?.data?.error ?? e?.message ?? "error",
      });
    }
    batchProgress.value = i + 1;
  }
  batchRunning.value = false;
  batchDone.value = true;
}

const showPrintDialog = ref(false);
const printEntries = ref<VoucherEntry[]>([]);

function openPrintResults() {
  const profileName = lastBatchProfile.value || "default";
  printEntries.value = batchResults.value
    .filter((r) => r.ok)
    .map((r) => ({
      name: r.name,
      password: r.password,
      profile: profileName,
    }));
  showPrintDialog.value = true;
}
</script>
