<template>
  <PageLayout title="Hotspot" subtitle="Profiles">
    <template #actions>
      <button
        v-if="store.activeId && !loading && !error"
        class="btn btn-primary"
        @click="openCreate"
      >
        <PlusIcon class="size-3.5" />
        New profile
      </button>
    </template>

    <NoRouterSelected v-if="!store.activeId" />

    <div v-else-if="loading" class="flex justify-center py-10">
      <span class="spinner" />
    </div>

    <div
      v-else-if="error"
      class="flex items-center gap-2 p-4 rounded-xl text-sm border bg-red/8 border-red/20 text-red"
    >
      <ExclamationTriangleIcon class="size-4 shrink-0" />
      {{ error }}
      <button class="ml-auto text-xs underline" @click="load">Retry</button>
    </div>

    <!-- Profile cards grid -->
    <div v-else-if="profiles.length > 0" class="space-y-3">
      <div class="grid grid-cols-[repeat(auto-fill,minmax(260px,1fr))] gap-5">
        <div
          v-for="p in profiles"
          :key="p['.id']"
          class="border border-border rounded-xl flex flex-col transition-colors hover:border-muted bg-surface"
        >
          <div class="p-4 flex flex-col gap-3">
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

            <div class="grid grid-cols-2 gap-x-3 gap-y-1 text-sm">
              <span class="text-text-secondary">Validity</span>
              <span class="text-text-primary">{{
                profileMetas[p.name]?.validity || "—"
              }}</span>
              <span class="text-text-secondary">Shared users</span>
              <span class="text-text-primary">{{
                p["shared-users"] || "1"
              }}</span>
              <span class="text-text-secondary">Price</span>
              <span class="text-text-primary">{{
                profileMetas[p.name]?.price
                  ? `${profileMetas[p.name].price} ${currency}`
                  : "—"
              }}</span>
            </div>
          </div>

          <div class="h-px bg-border" />

          <div class="flex flex-wrap gap-2 py-3 px-3.5">
            <button
              class="btn btn-ghost flex-1 min-w-0 justify-center whitespace-normal text-center"
              @click="openEdit(p)"
            >
              <PencilSquareIcon class="size-3.5 shrink-0" />
              Edit
            </button>
            <button
              class="btn border-transparent bg-red-10 text-red hover:bg-red/50 hover:text-white flex-1 min-w-0 justify-center whitespace-normal text-center"
              @click="remove(p)"
            >
              <TrashIcon class="size-3.5 shrink-0" />
              Delete
            </button>
          </div>
        </div>
      </div>
    </div>

    <div
      v-else
      class="border border-dashed border-border rounded-xl py-16 text-center"
    >
      <RectangleGroupIcon class="size-8 text-text-muted mx-auto mb-3" />
      <p class="text-sm font-medium text-text-secondary">No profiles yet</p>
      <p class="text-sm text-text-muted mt-1">
        Create your first profile to define bandwidth and pricing tiers.
      </p>
      <button class="btn btn-primary btn-sm mt-4" @click="openCreate">
        <PlusIcon class="size-3.5" />
        New profile
      </button>
    </div>

    <!-- Create / Edit dialog -->
    <AppDialog
      :open="showForm"
      :title="editing ? 'Edit Profile' : 'Create New profile'"
      @update:open="showForm = $event"
    >
      <TooltipProvider :delay-duration="200">
        <form class="space-y-4" @submit.prevent="submit">
          <label class="flex flex-col gap-1">
            <span class="font-medium">Name <span>*</span></span>
            <input
              v-model="form.name"
              class="input"
              required
              :disabled="!!editing"
              placeholder="e.g. 1hour, daily, weekly"
            />
          </label>

          <div class="flex flex-col gap-1">
            <FieldLabel
              label="Validity"
              tip="How long the account stays active after first login. Use: 30m, 2h, 1d, 1w, or combine: 1d12h. Leave blank for unlimited."
            />
            <input
              v-model="validityRaw"
              class="input font-mono"
              placeholder="e.g. 1h, 1d, 1w, 1d12h"
              @blur="normalizeValidity"
            />
            <p v-if="validityRaw && !validityPreview" class="text-xs text-red">
              Invalid format — use: 30m, 2h, 1d, 1w or combinations like 1d12h
            </p>
            <p v-else-if="validityPreview" class="text-xs text-text-secondary">
              Sends to router:
              <span class="font-mono">{{ validityPreview }}</span>
              <span class="ml-1 opacity-70"
                >(ROS {{ rosVersion || "…" }})</span
              >
            </p>
          </div>

          <div class="flex flex-col gap-1">
            <FieldLabel
              label="Rate limit"
              tip="Speed cap for this profile. Upload = client sending data. Download = client receiving data. Leave at 0 for unlimited."
            />
            <div class="flex gap-2">
              <div class="flex flex-col gap-1 flex-1 min-w-0">
                <span class="text-sm text-text-secondary">Upload</span>
                <div class="join">
                  <input
                    v-model.number="rateUp"
                    type="number"
                    min="0"
                    class="input min-w-0 flex-1"
                    placeholder="0"
                  />
                  <AppSelect
                    v-model="rateUpUnit"
                    :options="rateUnitOptions"
                    trigger-class="w-auto shrink-0"
                  />
                </div>
              </div>
              <div class="flex flex-col gap-1 flex-1 min-w-0">
                <span class="text-sm text-text-secondary">Download</span>
                <div class="join">
                  <input
                    v-model.number="rateDown"
                    type="number"
                    min="0"
                    class="input min-w-0 flex-1"
                    placeholder="0"
                  />
                  <AppSelect
                    v-model="rateDownUnit"
                    :options="rateUnitOptions"
                    trigger-class="w-auto shrink-0"
                  />
                </div>
              </div>
            </div>
          </div>

          <div class="grid grid-cols-2 gap-3">
            <div class="flex flex-col gap-1">
              <FieldLabel
                label="Shared users"
                tip="How many devices can be logged in simultaneously with the same credentials."
              />
              <input
                v-model="form.sharedUsers"
                type="number"
                min="1"
                class="input"
                placeholder="1"
              />
            </div>
            <label class="flex flex-col gap-1">
              <span class="text-sm font-medium text-text-secondary"
                >Price{{ currency ? ` (${currency})` : "" }}</span
              >
              <input
                v-model="form.price"
                class="input"
                placeholder="e.g. 500"
              />
            </label>
          </div>

          <div class="flex flex-col gap-1">
            <FieldLabel
              label="Address pool"
              tip="IP address pool to assign to users of this profile. Leave blank to use the hotspot's default pool."
            />
            <AppSelect
              v-model="form.addressPool"
              :options="
                addressPools.map((p) => ({
                  value: p.name,
                  label: p.ranges ? `${p.name} (${p.ranges})` : p.name,
                }))
              "
              placeholder="Use hotspot's default pool"
            />
          </div>

          <p v-if="formError" class="text-xs text-red">{{ formError }}</p>

          <div class="flex justify-end gap-2 pt-2">
            <button
              type="button"
              class="btn btn-ghost"
              @click="showForm = false"
            >
              Cancel
            </button>
            <button
              type="submit"
              class="btn btn-primary"
              :disabled="
                submitting ||
                (!!validityRaw && !validityPreview) ||
                (!rosVersion && !!validityPreview)
              "
            >
              <span
                v-if="submitting"
                class="size-4 border-2 border-black/20 border-t-black rounded-full animate-spin"
              />
              {{ editing ? "Save" : "Create" }}
            </button>
          </div>
        </form>
      </TooltipProvider>
    </AppDialog>
  </PageLayout>
</template>

<script setup lang="ts">
import { ref, computed, watch, defineComponent, h } from "vue";
import {
  PlusIcon,
  TrashIcon,
  PencilSquareIcon,
  ExclamationTriangleIcon,
  RectangleGroupIcon,
  QuestionMarkCircleIcon,
} from "@heroicons/vue/24/outline";
import {
  TooltipProvider,
  TooltipRoot,
  TooltipTrigger,
  TooltipContent,
  TooltipPortal,
} from "reka-ui";

import { useRoutersStore } from "@/stores/routers";
import {
  listHotspotProfiles,
  createHotspotProfile,
  updateHotspotProfile,
  deleteHotspotProfile,
  getHotspotSettings,
  getSystemResource,
  getProfileMetas,
  listAddressPools,
  type ProfileMeta,
} from "@/api";
import { friendlyError } from "@/utils/errors";
import AppDialog from "@/components/AppDialog.vue";
import PageLayout from "@/components/PageLayout.vue";
import NoRouterSelected from "@/components/NoRouterSelected.vue";
import AppSelect from "@/components/AppSelect.vue";

const store = useRoutersStore();

const FieldLabel = defineComponent({
  props: { label: String, tip: String, iconOnly: Boolean },
  setup(props) {
    return () =>
      h("div", { class: "flex items-center gap-1" }, [
        props.iconOnly
          ? null
          : h("span", { class: "font-medium" }, props.label),
        h(
          TooltipRoot,
          {},
          {
            default: () => [
              h(
                TooltipTrigger,
                { asChild: true },
                {
                  default: () =>
                    h(
                      "button",
                      {
                        type: "button",
                        class:
                          "text-text-muted hover:text-text-secondary transition-colors",
                      },
                      h(QuestionMarkCircleIcon, { class: "size-3.5" }),
                    ),
                },
              ),
              h(
                TooltipPortal,
                {},
                {
                  default: () =>
                    h(
                      TooltipContent,
                      {
                        class:
                          "z-50 max-w-xs px-3 py-2 text-xs rounded-lg shadow-lg leading-relaxed",
                        style:
                          "background: var(--color-surface); color: var(--color-text-primary); border: 1px solid var(--color-border);",
                        side: "top",
                        sideOffset: 6,
                      },
                      { default: () => props.tip },
                    ),
                },
              ),
            ],
          },
        ),
      ]);
  },
});

const rateUnitOptions = [
  { value: "k", label: "kbps" },
  { value: "M", label: "Mbps" },
];

const profiles = ref<Record<string, string>[]>([]);
const profileMetas = ref<Record<string, ProfileMeta>>({});
const addressPools = ref<Record<string, string>[]>([]);
const loading = ref(false);
const error = ref("");
const currency = ref("");
const rosVersion = ref("");

const showForm = ref(false);
const editing = ref<string | null>(null);
const submitting = ref(false);
const formError = ref("");

const rateUp = ref(0);
const rateUpUnit = ref<"k" | "M">("M");
const rateDown = ref(0);
const rateDownUnit = ref<"k" | "M">("M");
const validityRaw = ref("");

const emptyForm = () => ({
  name: "",
  addressPool: "",
  sharedUsers: "1",
  price: "",
});
const form = ref(emptyForm());

const isMajorV7 = computed(() => rosVersion.value.startsWith("7"));

function parseShorthand(
  s: string,
): { w: number; d: number; h: number; m: number } | null {
  if (!s.trim()) return null;
  const re = /^(?:(\d+)w)?(?:(\d+)d)?(?:(\d+)h)?(?:(\d+)m)?$/i;
  const m = s.trim().match(re);
  if (!m || !m[0]) return null;
  const [, w, d, h, min] = m;
  if (!w && !d && !h && !min) return null;
  return {
    w: parseInt(w || "0"),
    d: parseInt(d || "0"),
    h: parseInt(h || "0"),
    m: parseInt(min || "0"),
  };
}

const validityPreview = computed(() => {
  const p = parseShorthand(validityRaw.value);
  if (!p) return validityRaw.value ? "" : "";
  const totalDays = p.w * 7 + p.d;
  const totalHours = p.h;
  const totalMins = p.m;
  if (!totalDays && !totalHours && !totalMins) return "";
  if (isMajorV7.value) {
    let s = "P";
    if (totalDays) s += `${totalDays}D`;
    if (totalHours || totalMins) {
      s += "T";
      if (totalHours) s += `${totalHours}H`;
      if (totalMins) s += `${totalMins}M`;
    }
    return s;
  }
  return [
    totalDays ? `${totalDays}d` : "",
    totalHours ? `${totalHours}h` : "",
    totalMins ? `${totalMins}m` : "",
  ]
    .filter(Boolean)
    .join("");
});

function normalizeValidity() {
  const p = parseShorthand(validityRaw.value);
  if (!p) return;
  validityRaw.value = [
    p.w ? `${p.w}w` : "",
    p.d ? `${p.d}d` : "",
    p.h ? `${p.h}h` : "",
    p.m ? `${p.m}m` : "",
  ]
    .filter(Boolean)
    .join("");
}

const builtRateLimit = computed(() => {
  const up = rateUp.value ? `${rateUp.value}${rateUpUnit.value}` : "";
  const down = rateDown.value ? `${rateDown.value}${rateDownUnit.value}` : "";
  if (!up && !down) return "";
  return `${up || "0"}/${down || "0"}`;
});

function validityToShorthand(s: string): string {
  if (!s) return "";
  const iso = s.match(/^P(?:(\d+)D)?(?:T(?:(\d+)H)?(?:(\d+)M)?)?$/);
  if (iso) {
    const d = parseInt(iso[1] || "0");
    const h = parseInt(iso[2] || "0");
    const m = parseInt(iso[3] || "0");
    const weeks = Math.floor(d / 7);
    const days = d % 7;
    return [
      weeks ? `${weeks}w` : "",
      days ? `${days}d` : "",
      h ? `${h}h` : "",
      m ? `${m}m` : "",
    ]
      .filter(Boolean)
      .join("");
  }
  const d = s.match(/(\d+)d/i);
  const h = s.match(/(\d+)h/i);
  const m = s.match(/(\d+)m/i);
  const days = d ? parseInt(d[1]) : 0;
  const hours = h ? parseInt(h[1]) : 0;
  const mins = m ? parseInt(m[1]) : 0;
  const weeks = Math.floor(days / 7);
  const remDays = days % 7;
  return [
    weeks ? `${weeks}w` : "",
    remDays ? `${remDays}d` : "",
    hours ? `${hours}h` : "",
    mins ? `${mins}m` : "",
  ]
    .filter(Boolean)
    .join("");
}

async function load() {
  if (!store.activeId) return;
  loading.value = true;
  error.value = "";
  try {
    const [p, s, res, m, pools] = await Promise.all([
      listHotspotProfiles(store.activeId),
      getHotspotSettings(store.activeId).catch(() => ({ currency: "" })),
      getSystemResource(store.activeId).catch(() => ({})),
      getProfileMetas(store.activeId).catch(
        () => ({}) as Record<string, ProfileMeta>,
      ),
      listAddressPools(store.activeId).catch(() => []),
    ]);
    profiles.value = p;
    profileMetas.value = m;
    addressPools.value = pools;
    currency.value = (s as any).currency ?? "";
    rosVersion.value =
      (res as any)["ros-version"] ?? (res as any)["version"] ?? "";
  } catch (e: any) {
    error.value = friendlyError(e, "Failed to load profiles");
  } finally {
    loading.value = false;
  }
}

function openCreate() {
  editing.value = null;
  form.value = emptyForm();
  rateUp.value = 0;
  rateUpUnit.value = "M";
  rateDown.value = 0;
  rateDownUnit.value = "M";
  validityRaw.value = "";
  formError.value = "";
  showForm.value = true;
}

function openEdit(p: Record<string, string>) {
  editing.value = p[".id"];
  const meta = profileMetas.value[p.name] ?? {};
  form.value = {
    name: p.name ?? "",
    addressPool: p["address-pool"] ?? "",
    sharedUsers: p["shared-users"] ?? "",
    price: (meta as ProfileMeta).price ?? "",
  };
  const rl = p["rate-limit"] ?? "";
  const [upStr, downStr] = rl.split("/");
  parseRate(upStr ?? "", rateUp, rateUpUnit);
  parseRate(downStr ?? "", rateDown, rateDownUnit);
  validityRaw.value = validityToShorthand((meta as ProfileMeta).validity ?? "");
  formError.value = "";
  showForm.value = true;
}

function parseRate(s: string, val: typeof rateUp, unit: typeof rateUpUnit) {
  const m = s.match(/^(\d+)(k|M)$/i);
  if (m) {
    val.value = parseInt(m[1]);
    unit.value = m[2].toUpperCase() === "M" ? "M" : "k";
  } else {
    val.value = 0;
  }
}

async function submit() {
  if (!store.activeId) return;
  submitting.value = true;
  formError.value = "";
  try {
    const params = {
      name: form.value.name,
      addressPool: form.value.addressPool,
      sharedUsers: form.value.sharedUsers,
      rateLimit: builtRateLimit.value,
      validity: validityRaw.value,
      price: form.value.price,
    };
    if (editing.value) {
      await updateHotspotProfile(store.activeId, editing.value, params);
    } else {
      await createHotspotProfile(store.activeId, params);
    }
    showForm.value = false;
    await load();
  } catch (e: any) {
    formError.value = friendlyError(e, "Failed to save profile");
  } finally {
    submitting.value = false;
  }
}

async function remove(p: Record<string, string>) {
  if (
    !store.activeId ||
    !confirm(
      "Delete this profile? Users assigned to it will fall back to the default profile.",
    )
  )
    return;
  try {
    await deleteHotspotProfile(store.activeId, p[".id"], p.name);
    await load();
  } catch (e: any) {
    error.value = friendlyError(e, "Failed to delete profile");
  }
}

watch(() => store.activeId, load, { immediate: true });
</script>
