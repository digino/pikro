<template>
  <div v-if="!store.activeId" class="p-6">
    <NoRouterSelected />
  </div>

  <main v-else class="flex flex-col gap-6 px-4 sm:px-6 lg:px-10 pt-7 pb-14">
    <header class="flex flex-wrap items-end justify-between gap-4">
      <div class="flex items-center gap-4">
        <img :src="inkIcon('profiles')" alt="" width="64" height="64" class="shrink-0" />
        <div>
          <nav :aria-label="t('profiles.breadcrumb')" class="eyebrow">
            {{ t('profiles.breadcrumb') }} / {{ t('profiles.title') }}
          </nav>
          <h1 class="type-display text-3xl">{{ t('profiles.title') }}</h1>
          <p class="text-sm text-text-secondary">{{ t('profiles.subtitle') }}</p>
        </div>
      </div>
      <button
        v-if="!loading && !error"
        type="button"
        class="btn btn-primary"
        @click="openCreate"
      >
        <PlusIcon class="size-4" />
        {{ t('profiles.newProfile') }}
      </button>
    </header>

    <div v-if="loading" class="flex justify-center py-10">
      <span class="spinner" />
    </div>

    <div
      v-else-if="error"
      role="alert"
      class="flex items-center gap-2.5 p-4 rounded-control border-2 border-danger-text bg-danger-wash text-danger-text text-sm"
    >
      <ExclamationTriangleIcon class="size-5 shrink-0" />
      {{ error }}
      <button type="button" class="btn btn-sm ml-auto" @click="load">{{ t('profiles.retry') }}</button>
    </div>

    <div v-else class="grid grid-cols-1 xl:grid-cols-[2fr_1fr] items-start gap-7">
      <!-- Profile cards grid -->
      <section :aria-label="t('profiles.title')">
        <div v-if="profiles.length === 0" class="card border-dashed shadow-none p-10 text-center">
          <img :src="inkIcon('profiles')" alt="" width="64" height="64" class="mx-auto" />
          <p class="m-0 mt-3 font-semibold">{{ t('profiles.noProfilesTitle') }}</p>
          <p class="m-0 mt-1 text-sm text-text-muted">{{ t('profiles.noProfilesBody') }}</p>
        </div>

        <div v-else class="grid grid-cols-1 lg:grid-cols-2 gap-5.5">
          <article
            v-for="c in cards"
            :key="c.id"
            class="flex flex-col overflow-hidden rounded-card border-2 border-ink"
            :class="c.cardClass"
          >
            <div class="flex items-start justify-between gap-3 px-4.5 pt-4.5 pb-3.5">
              <div class="min-w-0">
                <div class="eyebrow" :class="{ 'text-current opacity-75': c.free }">{{ c.kind }}</div>
                <h2 class="type-display m-0 mt-1 text-[28px] leading-none truncate">{{ c.name }}</h2>
              </div>
              <span
                class="type-display shrink-0 whitespace-nowrap rounded-[10px] border-2 px-2.5 py-1.5 text-[15px]"
                :class="c.priced ? 'bg-primary text-primary-fg border-ink' : 'border-current'"
              >
                {{ c.price }}
              </span>
            </div>
            <div class="mx-4.5 border-t-2 border-dashed border-current opacity-90" />
            <dl class="m-0 grid grid-cols-[auto_1fr] gap-x-4 gap-y-2 px-4.5 py-3.5 text-sm">
              <dt class="eyebrow self-center" :class="{ 'text-current opacity-75': c.free }">{{ t('profiles.validity') }}</dt>
              <dd class="m-0 font-semibold">{{ c.validity }}</dd>
              <dt class="eyebrow self-center" :class="{ 'text-current opacity-75': c.free }">{{ t('profiles.speed') }}</dt>
              <dd class="m-0 font-mono text-[13px]">{{ c.speed }}</dd>
              <dt class="eyebrow self-center" :class="{ 'text-current opacity-75': c.free }">{{ t('profiles.devices') }}</dt>
              <dd class="m-0">{{ c.devices }}</dd>
            </dl>
            <div
              class="mt-auto flex flex-wrap items-center justify-between gap-2.5 px-4.5 py-3"
              :class="c.footClass"
            >
              <span class="text-[13px]">{{ c.stats }}</span>
              <div class="flex gap-2">
                <RouterLink
                  v-if="!c.system"
                  :to="{ path: '/hotspot/vouchers', query: { profile: c.name } }"
                  class="btn btn-sm"
                  :aria-label="t('profiles.generateFor', { name: c.name })"
                >
                  {{ t('profiles.generate') }}
                </RouterLink>
                <button
                  type="button"
                  class="btn btn-sm"
                  :class="{ 'btn-primary': c.editing }"
                  :aria-pressed="c.editing"
                  :aria-label="t('profiles.editProfile', { name: c.name })"
                  @click="openEdit(c.raw)"
                >
                  {{ t('profiles.edit') }}
                </button>
              </div>
            </div>
          </article>
        </div>
      </section>

      <!-- Create / Edit form -->
      <aside
        :aria-label="formTitle"
        class="card p-6 flex flex-col gap-4.5 xl:sticky xl:top-6"
      >
        <TooltipProvider :delay-duration="200">
          <form class="flex flex-col gap-4.5" @submit.prevent="submit">
            <div class="flex items-center justify-between gap-3">
              <h2 class="type-display m-0 text-[26px] truncate">{{ formTitle }}</h2>
              <span class="chip bg-primary text-primary-fg">{{ formChip }}</span>
            </div>

            <label class="flex flex-col gap-1.5">
              <span class="eyebrow">{{ t('profiles.name') }}</span>
              <input
                ref="nameInput"
                v-model="form.name"
                class="input"
                required
                :disabled="!!editing"
                :placeholder="t('profiles.namePlaceholder')"
              />
            </label>

            <div class="flex flex-col gap-1.5">
              <FieldLabel
                :label="t('profiles.validity')"
                tip="How long the account stays active after first login. Use: 30m, 2h, 1d, 1w, or combine: 1d12h. Leave blank for unlimited."
              />
              <div class="flex flex-wrap gap-2">
                <button
                  v-for="v in VALIDITY_CHIPS"
                  :key="v"
                  type="button"
                  :aria-pressed="validityRaw === v"
                  class="min-h-10 min-w-13 px-3 rounded-full border-2 border-ink text-[13px] font-semibold cursor-pointer"
                  :class="validityRaw === v ? 'bg-ink text-paper' : 'bg-paper text-ink hover:bg-primary-wash'"
                  @click="validityRaw = v"
                >
                  {{ v }}
                </button>
                <button
                  type="button"
                  :aria-pressed="isCustomValidity"
                  class="min-h-10 px-3 rounded-full border-2 border-ink text-[13px] font-semibold cursor-pointer"
                  :class="isCustomValidity ? 'bg-ink text-paper' : 'bg-paper text-ink hover:bg-primary-wash'"
                  @click="validityInput?.focus()"
                >
                  {{ t('profiles.custom') }}
                </button>
              </div>
              <input
                ref="validityInput"
                v-model="validityRaw"
                class="input font-mono"
                :aria-label="t('profiles.validityCustom')"
                :placeholder="t('profiles.validityPlaceholder')"
                @blur="normalizeValidity"
              />
              <p v-if="validityRaw && !validityPreview" class="m-0 text-xs text-red">
                {{ t('profiles.validityInvalid') }}
              </p>
              <p v-else-if="validityPreview" class="m-0 text-xs text-text-secondary">
                {{ t('profiles.sendsToRouter') }}
                <span class="font-mono">{{ validityPreview }}</span>
                <span class="ml-1 opacity-70">(ROS {{ rosVersion || "…" }})</span>
              </p>
            </div>

            <label class="flex flex-col gap-1.5">
              <span class="eyebrow">{{ t('profiles.price') }}</span>
              <span class="join">
                <input v-model="form.price" class="input min-w-0 flex-1 font-semibold" inputmode="decimal" placeholder="500" />
                <span
                  v-if="currency"
                  class="inline-flex items-center px-3 border-2 border-ink rounded-control bg-primary text-primary-fg text-[13px] font-bold"
                >
                  {{ currency }}
                </span>
              </span>
            </label>

            <div class="flex flex-col gap-1.5">
              <FieldLabel
                :label="t('profiles.speedLimit')"
                tip="Speed cap for this profile. Upload = client sending data. Download = client receiving data. Leave at 0 for unlimited."
              />
              <div class="grid grid-cols-2 gap-3">
                <div class="flex flex-col gap-1">
                  <label class="text-[13px] text-text-muted" for="rate-down">{{ t('profiles.down') }}</label>
                  <div class="join">
                    <input
                      id="rate-down"
                      v-model.number="rateDown"
                      type="number"
                      min="0"
                      class="input min-w-0 flex-1"
                      placeholder="0"
                    />
                    <AppSelect v-model="rateDownUnit" :options="rateUnitOptions" trigger-class="w-auto shrink-0" />
                  </div>
                </div>
                <div class="flex flex-col gap-1">
                  <label class="text-[13px] text-text-muted" for="rate-up">{{ t('profiles.up') }}</label>
                  <div class="join">
                    <input
                      id="rate-up"
                      v-model.number="rateUp"
                      type="number"
                      min="0"
                      class="input min-w-0 flex-1"
                      placeholder="0"
                    />
                    <AppSelect v-model="rateUpUnit" :options="rateUnitOptions" trigger-class="w-auto shrink-0" />
                  </div>
                </div>
              </div>
            </div>

            <div class="grid grid-cols-2 gap-3">
              <div class="flex flex-col gap-1.5">
                <FieldLabel
                  :label="t('profiles.devicesPerVoucher')"
                  tip="How many devices can be logged in simultaneously with the same credentials."
                />
                <div class="flex items-center min-h-11 border-2 border-ink rounded-control overflow-hidden bg-paper">
                  <button
                    type="button"
                    :aria-label="t('profiles.fewerDevices')"
                    class="inline-flex items-center justify-center size-11 border-r-2 border-ink cursor-pointer hover:bg-primary-wash disabled:opacity-40"
                    :disabled="sharedUsers <= 1"
                    @click="sharedUsers = sharedUsers - 1"
                  >
                    <MinusIcon class="size-4.5" />
                  </button>
                  <span class="type-display flex-1 text-center text-xl" aria-live="polite">{{ sharedUsers }}</span>
                  <button
                    type="button"
                    :aria-label="t('profiles.moreDevices')"
                    class="inline-flex items-center justify-center size-11 border-l-2 border-ink bg-primary text-primary-fg cursor-pointer"
                    @click="sharedUsers = sharedUsers + 1"
                  >
                    <PlusIcon class="size-4.5" />
                  </button>
                </div>
              </div>
              <div class="flex flex-col gap-1.5">
                <FieldLabel
                  :label="t('profiles.addressPool')"
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
                  :placeholder="t('profiles.defaultPool')"
                />
              </div>
            </div>

            <div class="flex items-start gap-2.5 p-3 rounded-control bg-ground text-[13px] leading-snug">
              <img :src="inkIcon('settings')" alt="" width="28" height="28" class="shrink-0" />
              <span>{{ t('profiles.note') }}</span>
            </div>

            <p v-if="formError" role="alert" class="m-0 text-sm text-red">{{ formError }}</p>

            <div class="flex flex-wrap gap-2.5">
              <button
                type="submit"
                class="btn btn-primary flex-1 min-h-12.5 text-[15px]"
                :disabled="
                  submitting ||
                  (!!validityRaw && !validityPreview) ||
                  (!rosVersion && !!validityPreview)
                "
              >
                <span v-if="submitting" class="spinner spinner--sm" />
                {{ t('profiles.save') }}
              </button>
              <button type="button" class="btn min-h-12.5" @click="openCreate">{{ t('profiles.cancel') }}</button>
              <button
                v-if="editingProfile && editingProfile.name !== 'default'"
                type="button"
                class="btn btn-danger min-h-12.5"
                @click="remove(editingProfile)"
              >
                <TrashIcon class="size-4" />
                {{ t('profiles.delete') }}
              </button>
            </div>
          </form>
        </TooltipProvider>
      </aside>
    </div>
  </main>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick, defineComponent, h } from "vue";
import { RouterLink } from "vue-router";
import { useI18n } from "vue-i18n";
import {
  PlusIcon,
  MinusIcon,
  TrashIcon,
  ExclamationTriangleIcon,
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
  listHotspotUsers,
  listHotspotActive,
  getSalesLedger,
  type ProfileMeta,
} from "@/api";
import { friendlyError } from "@/utils/errors";
import { inkIcon } from "@/utils/inkIcons";
import NoRouterSelected from "@/components/NoRouterSelected.vue";
import AppSelect from "@/components/AppSelect.vue";

const store = useRoutersStore();
const { t, locale } = useI18n();

const FieldLabel = defineComponent({
  props: { label: String, tip: String, iconOnly: Boolean },
  setup(props) {
    return () =>
      h("div", { class: "flex items-center gap-1" }, [
        props.iconOnly
          ? null
          : h("span", { class: "eyebrow" }, props.label),
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
                        "aria-label": props.tip,
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
                          "z-50 max-w-xs px-3 py-2 text-xs rounded-lg leading-relaxed bg-paper text-ink border-2 border-ink shadow-btn",
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

const VALIDITY_CHIPS = ["1h", "1d", "7d", "30d"];

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
// Per-profile counters for the card footers: vouchers sold this year, users online now.
const soldByProfile = ref<Record<string, number>>({});
const onlineByProfile = ref<Record<string, number>>({});

const editing = ref<string | null>(null);
const submitting = ref(false);
const formError = ref("");
const nameInput = ref<HTMLInputElement | null>(null);
const validityInput = ref<HTMLInputElement | null>(null);

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

const sharedUsers = computed({
  get: () => Math.max(1, parseInt(form.value.sharedUsers) || 1),
  set: (n: number) => {
    form.value.sharedUsers = String(Math.max(1, n));
  },
});

const isCustomValidity = computed(() => !VALIDITY_CHIPS.includes(validityRaw.value));

const editingProfile = computed(() => profiles.value.find((p) => p[".id"] === editing.value) ?? null);

const formTitle = computed(() =>
  editingProfile.value
    ? t("profiles.formEdit", { name: editingProfile.value.name })
    : t("profiles.formNew"),
);
const formChip = computed(() =>
  editingProfile.value
    ? t("profiles.soldChip", { n: soldByProfile.value[editingProfile.value.name] ?? 0 })
    : t("profiles.draft"),
);

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

// ── Cards ────────────────────────────────────────────────────────────────
const cards = computed(() =>
  profiles.value.map((p) => {
    const meta = profileMetas.value[p.name];
    const price = parseFloat(meta?.price ?? "");
    const system = p.name === "default";
    const free = !price && !system;
    const isEditing = p[".id"] === editing.value;
    const [up, down] = (p["rate-limit"] ?? "").split("/");
    const shared = parseInt(p["shared-users"] ?? "1") || 1;
    return {
      id: p[".id"],
      raw: p,
      name: p.name,
      system,
      free,
      priced: !!price,
      editing: isEditing,
      kind: system ? t("profiles.kindSystem") : free ? t("profiles.kindFree") : t("profiles.kindForSale"),
      price: price
        ? `${new Intl.NumberFormat(locale.value).format(price)}${currency.value ? " " + currency.value : ""}`
        : system ? t("profiles.noPrice") : t("profiles.free"),
      validity: validityToShorthand(meta?.validity ?? "") || (system ? t("profiles.routerDefault") : t("profiles.unlimited")),
      speed: up || down ? t("profiles.speedValue", { down: down || "—", up: up || "—" }) : t("profiles.noLimit"),
      devices: t("profiles.devicesCount", { n: shared }, shared),
      stats: t("profiles.soldOnline", {
        sold: soldByProfile.value[p.name] ?? 0,
        online: onlineByProfile.value[p.name] ?? 0,
      }),
      // Free plans read as ink tickets; RouterOS's own default is a dashed outline.
      cardClass: [
        system ? "border-dashed" : "",
        free ? "bg-ink text-paper" : "bg-paper text-ink",
        isEditing
          ? "shadow-[6px_6px_0_var(--color-primary),6px_6px_0_2px_var(--color-ink)]"
          : system ? "" : "shadow-card",
      ],
      footClass: [
        "border-t-2",
        system ? "border-dashed" : "",
        free ? "border-paper bg-white/5" : "border-ink bg-ground/60",
      ],
    };
  }),
);

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
    error.value = friendlyError(e, t("profiles.failedToLoad"));
  } finally {
    loading.value = false;
  }
  loadCounters();
}

// Non-critical: the cards render without these and fill in when they arrive.
async function loadCounters() {
  if (!store.activeId) return;
  const [ledger, users, active] = await Promise.all([
    getSalesLedger(store.activeId, new Date().getFullYear()).catch(() => []),
    listHotspotUsers(store.activeId).catch(() => []),
    listHotspotActive(store.activeId).catch(() => []),
  ]);
  const sold: Record<string, number> = {};
  for (const e of ledger) sold[e.profile] = (sold[e.profile] ?? 0) + e.count;
  const profileOf = new Map(users.map((u) => [u.name, u.profile || "default"]));
  const online: Record<string, number> = {};
  for (const s of active) {
    const prof = profileOf.get(s.user);
    if (prof) online[prof] = (online[prof] ?? 0) + 1;
  }
  soldByProfile.value = sold;
  onlineByProfile.value = online;
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
  nextTick(() => nameInput.value?.focus());
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
    await load();
    openCreate();
  } catch (e: any) {
    formError.value = friendlyError(e, t("profiles.failedToSave"));
  } finally {
    submitting.value = false;
  }
}

async function remove(p: Record<string, string>) {
  if (!store.activeId || !confirm(t("profiles.confirmDelete"))) return;
  try {
    await deleteHotspotProfile(store.activeId, p[".id"], p.name);
    openCreate();
    await load();
  } catch (e: any) {
    error.value = friendlyError(e, t("profiles.failedToDelete"));
  }
}

watch(() => store.activeId, load, { immediate: true });
</script>
