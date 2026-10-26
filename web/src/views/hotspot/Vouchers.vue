<template>
  <div v-if="!store.activeId" class="p-6">
    <NoRouterSelected />
  </div>

  <main v-else class="flex flex-col gap-6 px-4 sm:px-6 lg:px-10 pt-7 pb-14">
    <header class="flex items-center gap-4">
      <img :src="inkIcon('vouchers')" alt="" width="64" height="64" />
      <div>
        <nav :aria-label="t('vouchers.breadcrumb')" class="eyebrow">
          {{ t('vouchers.breadcrumb') }} / {{ t('nav.vouchers') }}
        </nav>
        <h1 class="type-display text-3xl">{{ t('vouchers.title') }}</h1>
      </div>
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
      <button type="button" class="btn btn-sm ml-auto" @click="load">{{ t('vouchers.retry') }}</button>
    </div>

    <div v-else-if="profiles.length === 0" class="card border-dashed shadow-none p-10 text-center">
      <img :src="inkIcon('profiles')" alt="" width="64" height="64" class="mx-auto" />
      <p class="m-0 mt-3 font-semibold">{{ t('vouchers.noPlansTitle') }}</p>
      <p class="m-0 mt-1 mx-auto max-w-100 text-sm text-text-muted">{{ t('vouchers.noPlansBody') }}</p>
      <RouterLink to="/hotspot/profiles" class="btn btn-ink mt-4">
        <PlusIcon class="size-4" />
        {{ t('vouchers.newPlan') }}
      </RouterLink>
    </div>

    <div v-else class="grid grid-cols-1 xl:grid-cols-2 items-start gap-7">
      <section
        :aria-label="t('vouchers.settings')"
        class="card min-w-0 p-6 flex flex-col gap-6.5"
      >
        <!-- 1. Plan -->
        <div class="flex flex-col gap-3">
          <div class="flex items-center gap-2.5">
            <span :class="STEP">1</span>
            <h2 class="m-0 text-[17px] font-bold">{{ t('vouchers.step1') }}</h2>
            <RouterLink
              to="/hotspot/profiles"
              class="ml-auto text-sm font-semibold text-ink underline underline-offset-2"
            >
              {{ t('vouchers.managePlans') }}
            </RouterLink>
          </div>
          <div class="grid grid-cols-2 gap-3">
            <button
              v-for="p in profiles"
              :key="p['.id']"
              type="button"
              :aria-pressed="form.profile === p.name"
              class="flex flex-col gap-2 px-3.5 py-3 text-left rounded-[14px] border-2 border-ink cursor-pointer transition-transform"
              :class="
                form.profile === p.name
                  ? 'bg-primary text-primary-fg shadow-[4px_4px_0_var(--color-shadow)] -rotate-1'
                  : 'bg-paper text-ink hover:bg-primary-wash'
              "
              @click="form.profile = p.name"
            >
              <span class="flex items-baseline justify-between gap-2">
                <span class="type-display text-xl truncate">{{ p.name }}</span>
                <span class="shrink-0 font-mono text-xs">
                  {{ profileMetas[p.name]?.validity || t('vouchers.noValidity') }}
                </span>
              </span>
              <span class="-mx-3.5 border-t-2 border-dashed border-current" />
              <span class="text-[15px] font-bold">{{ priceLabel(p.name) }}</span>
            </button>
          </div>
        </div>

        <!-- 2. Quantity -->
        <div class="flex flex-col gap-3">
          <div class="flex items-center gap-2.5">
            <span :class="STEP">2</span>
            <h2 class="m-0 text-[17px] font-bold">{{ t('vouchers.step2') }}</h2>
          </div>
          <div class="flex flex-wrap items-center gap-3">
            <div class="flex items-center border-2 border-ink rounded-control overflow-hidden bg-paper">
              <button
                type="button"
                :aria-label="t('vouchers.fewer')"
                class="inline-flex items-center justify-center size-12 border-r-2 border-ink cursor-pointer hover:bg-primary-wash disabled:opacity-40 disabled:cursor-not-allowed"
                :disabled="form.count <= 1"
                @click="setCount(form.count - 10)"
              >
                <MinusIcon class="size-5" />
              </button>
              <input
                v-model.number="form.count"
                type="number"
                min="1"
                :max="MAX_COUNT"
                :aria-label="t('vouchers.quantity')"
                class="type-display w-20 border-0 bg-transparent text-center text-2xl outline-none [appearance:textfield] [&::-webkit-inner-spin-button]:appearance-none"
                @blur="setCount(form.count)"
              />
              <button
                type="button"
                :aria-label="t('vouchers.more')"
                class="inline-flex items-center justify-center size-12 border-l-2 border-ink bg-primary text-primary-fg cursor-pointer disabled:opacity-40 disabled:cursor-not-allowed"
                :disabled="form.count >= MAX_COUNT"
                @click="setCount(form.count + 10)"
              >
                <PlusIcon class="size-5" />
              </button>
            </div>
            <button
              v-for="n in QUICK_COUNTS"
              :key="n"
              type="button"
              :aria-pressed="form.count === n"
              class="min-h-11 min-w-13 px-3 rounded-full border-2 border-ink text-sm font-semibold cursor-pointer"
              :class="form.count === n ? 'bg-ink text-paper' : 'bg-paper text-ink hover:bg-primary-wash'"
              @click="setCount(n)"
            >
              {{ n }}
            </button>
          </div>
          <p class="m-0 text-xs text-text-muted">{{ t('vouchers.maxHint') }}</p>
        </div>

        <!-- 3. Code style -->
        <div class="flex flex-col gap-3">
          <div class="flex items-center gap-2.5">
            <span :class="STEP">3</span>
            <h2 class="m-0 text-[17px] font-bold">{{ t('vouchers.step3') }}</h2>
          </div>
          <div class="grid grid-cols-3 gap-2.5">
            <label class="flex flex-col gap-1.5 min-w-0">
              <span class="eyebrow">{{ t('vouchers.prefix') }}</span>
              <input v-model.trim="form.prefix" class="input font-mono" maxlength="8" placeholder="lum-" />
            </label>
            <label class="flex flex-col gap-1.5 min-w-0">
              <span class="eyebrow">{{ t('vouchers.length') }}</span>
              <select v-model.number="form.nameLength" class="input-select font-semibold">
                <option v-for="n in NAME_LENGTHS" :key="n" :value="n">{{ n }}</option>
              </select>
            </label>
            <label class="flex flex-col gap-1.5 min-w-0">
              <span class="eyebrow">{{ t('vouchers.login') }}</span>
              <select
                :value="form.passwordMode === 'same' ? 'same' : 'random'"
                class="input-select font-semibold"
                @change="form.passwordMode = ($event.target as HTMLSelectElement).value as 'same' | 'random'"
              >
                <option value="same">{{ t('vouchers.codeOnly') }}</option>
                <option value="random">{{ t('vouchers.codeAndPassword') }}</option>
              </select>
            </label>
          </div>
        </div>

        <!-- 4. Template -->
        <div class="flex flex-col gap-3">
          <div class="flex items-center gap-2.5">
            <span :class="STEP">4</span>
            <h2 class="m-0 text-[17px] font-bold">{{ t('vouchers.step4') }}</h2>
          </div>
          <div class="grid grid-cols-1 sm:grid-cols-3 gap-2.5">
            <button
              v-for="tpl in VOUCHER_TEMPLATES"
              :key="tpl.key"
              type="button"
              :aria-pressed="form.template === tpl.key"
              class="flex flex-col gap-1 min-h-16 p-3 text-left rounded-control border-2 border-ink cursor-pointer"
              :class="
                form.template === tpl.key
                  ? 'bg-primary text-primary-fg shadow-btn'
                  : 'bg-paper text-ink hover:bg-primary-wash'
              "
              @click="form.template = tpl.key"
            >
              <span class="text-sm font-bold">{{ t(`templates.${tpl.key}.label`) }}</span>
              <span class="text-xs opacity-80">{{ t(`templates.${tpl.key}.description`) }}</span>
            </button>
          </div>
        </div>

        <!-- 5. Options — always visible so no admin misses them -->
        <div class="flex flex-col gap-3">
          <div class="flex items-center gap-2.5">
            <span :class="STEP">5</span>
            <h2 class="m-0 text-[17px] font-bold">{{ t('vouchers.step5') }}</h2>
          </div>
          <div class="flex flex-col gap-4">
            <div class="flex flex-col gap-1.5">
              <span class="eyebrow">{{ t('vouchers.characters') }}</span>
              <div role="group" :aria-label="t('vouchers.characters')" class="flex border-2 border-ink rounded-control overflow-hidden">
                <button
                  v-for="(opt, i) in CHARSET_OPTS"
                  :key="opt.key"
                  type="button"
                  :aria-pressed="charsetKey === opt.key"
                  class="flex-1 min-h-11 px-2 text-xs font-semibold cursor-pointer"
                  :class="[
                    i > 0 && 'border-l-2 border-ink',
                    charsetKey === opt.key ? 'bg-ink text-paper' : 'bg-paper text-ink hover:bg-primary-wash',
                  ]"
                  @click="setCharset(opt.key)"
                >
                  {{ t(opt.label) }}
                </button>
              </div>
            </div>

            <label class="flex items-center gap-2.5 text-sm font-semibold cursor-pointer">
              <input
                type="checkbox"
                class="size-5 accent-ink"
                :checked="form.passwordMode === 'fixed'"
                @change="form.passwordMode = ($event.target as HTMLInputElement).checked ? 'fixed' : 'random'"
              />
              {{ t('vouchers.samePassword') }}
            </label>
            <input
              v-if="form.passwordMode === 'fixed'"
              v-model="form.fixedPassword"
              class="input"
              :placeholder="t('vouchers.samePasswordPlaceholder')"
              :aria-label="t('vouchers.samePassword')"
            />

            <label class="flex flex-col gap-1.5">
              <span class="eyebrow">{{ t('vouchers.timeLimit') }}</span>
              <input
                v-model="form.limitUptimeRaw"
                class="input font-mono"
                :placeholder="t('vouchers.timeLimitPlaceholder')"
                @blur="form.limitUptimeRaw = normalizeShorthand(form.limitUptimeRaw)"
              />
              <span v-if="form.limitUptimeRaw && !uptimePreview" class="text-xs text-red">
                {{ t('vouchers.timeLimitInvalid') }}
              </span>
              <span v-else-if="uptimeWarning" class="text-xs text-text-secondary">
                {{ t('vouchers.timeLimitWarning') }}
              </span>
            </label>

            <div class="flex flex-col gap-1.5">
              <span class="eyebrow">{{ t('vouchers.dataLimit') }}</span>
              <div class="flex gap-2">
                <input
                  v-model.number="form.limitBytesTotalValue"
                  type="number"
                  min="0"
                  class="input flex-1 min-w-0"
                  :placeholder="t('vouchers.dataLimitPlaceholder')"
                  :aria-label="t('vouchers.dataLimit')"
                />
                <select v-model="form.limitBytesTotalUnit" class="input-select w-24" :aria-label="t('vouchers.dataLimit')">
                  <option value="M">MB</option>
                  <option value="G">GB</option>
                </select>
              </div>
            </div>

            <label class="flex flex-col gap-1.5">
              <span class="eyebrow">{{ t('vouchers.comment') }}</span>
              <input v-model="form.comment" class="input" :placeholder="t('vouchers.commentPlaceholder')" />
            </label>
          </div>
        </div>

        <!-- Summary -->
        <div class="-mx-2 -mb-2 flex flex-col gap-3.5 p-4.5 rounded-[14px] bg-sidebar text-sidebar-text">
          <div class="flex flex-wrap items-baseline justify-between gap-3">
            <span class="text-sm">
              {{ t('vouchers.summary', { n: form.count, plan: form.profile || '—' }, form.count) }}
            </span>
            <span v-if="totalLabel" class="type-display text-2xl">{{ totalLabel }}</span>
          </div>
          <p v-if="batchError" role="alert" class="m-0 text-sm text-primary">{{ batchError }}</p>
          <div class="flex flex-wrap gap-2.5">
            <button
              type="button"
              class="btn btn-primary flex-1 min-h-13 text-[15px] border-sidebar-text shadow-[3px_3px_0_var(--color-sidebar-text)]"
              :disabled="!canGenerate"
              @click="submitBatch(true)"
            >
              <PrinterIcon class="size-5" />
              {{ t('vouchers.generatePrint') }}
            </button>
            <button
              type="button"
              class="btn min-h-13 bg-transparent text-sidebar-text border-sidebar-text shadow-none"
              :disabled="!canGenerate"
              @click="submitBatch(false)"
            >
              {{ t('vouchers.saveOnly') }}
            </button>
          </div>
        </div>
      </section>

      <!-- Live print preview -->
      <section :aria-label="t('vouchers.preview')" class="min-w-0 flex flex-col gap-3.5 xl:sticky xl:top-6">
        <div class="eyebrow">{{ t('vouchers.preview') }} · {{ t(`templates.${form.template}.label`) }}</div>
        <div class="h-160 bg-white border-2 border-ink rounded-md shadow-[8px_8px_0_var(--color-shadow)] overflow-hidden">
          <iframe
            :srcdoc="previewHtml"
            sandbox=""
            :title="t('vouchers.preview')"
            class="w-full h-full border-0"
          />
        </div>
        <p class="m-0 text-sm text-text-muted">{{ t('vouchers.previewNote') }}</p>
      </section>
    </div>

    <!-- Progress & results -->
    <AppDialog
      :open="showBatch"
      :title="t('vouchers.resultTitle')"
      @update:open="onBatchDialogUpdate"
    >
      <div v-if="batchRunning" class="space-y-4">
        <div class="flex items-center justify-between text-sm">
          <span>{{ t('vouchers.creating') }}</span>
          <span class="font-mono">{{ batchProgress }} / {{ batchTotal }}</span>
        </div>
        <div class="h-3 w-full rounded-full border-2 border-ink bg-paper overflow-hidden">
          <div
            class="h-full bg-ink transition-[width]"
            :style="{ width: `${(batchProgress / batchTotal) * 100}%` }"
          />
        </div>
        <p class="m-0 text-xs text-text-muted font-mono">{{ batchCurrentName }}</p>
      </div>

      <div v-else class="space-y-4">
        <div class="flex flex-wrap items-center gap-2">
          <span class="chip bg-online-wash">
            <span class="size-2 rounded-full bg-online" />
            {{ t('vouchers.created', { n: createdCount }) }}
          </span>
          <span v-if="failedCount" class="chip bg-danger-wash text-danger-text border-danger-text">
            <span class="size-2 rotate-45 bg-danger" />
            {{ t('vouchers.failed', { n: failedCount }) }}
          </span>
        </div>
        <div class="max-h-64 overflow-y-auto rounded-control border-2 border-ink font-mono text-xs">
          <div
            v-for="r in batchResults"
            :key="r.name"
            class="flex items-center gap-3 px-3 py-2 border-b border-hairline last:border-0"
            :class="r.ok ? '' : 'bg-danger-wash text-danger-text'"
          >
            <span class="flex-1 font-semibold">{{ r.name }}</span>
            <span class="text-text-muted">{{ r.password }}</span>
            <span v-if="!r.ok" class="ml-auto">{{ r.error }}</span>
          </div>
        </div>
        <div class="flex justify-end gap-2">
          <button type="button" class="btn" @click="showBatch = false">{{ t('vouchers.close') }}</button>
          <button type="button" class="btn btn-ink" :disabled="!createdCount" @click="printResults">
            <PrinterIcon class="size-4" /> {{ t('vouchers.print') }}
          </button>
        </div>
      </div>
    </AppDialog>
  </main>
</template>

<script setup lang="ts">
import { ref, computed, watch } from "vue";
import { RouterLink, useRoute, useRouter } from "vue-router";
import { useI18n } from "vue-i18n";
import {
  PlusIcon,
  MinusIcon,
  ExclamationTriangleIcon,
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
import NoRouterSelected from "@/components/NoRouterSelected.vue";
import { printVouchers, toVoucherItems, type VoucherEntry } from "@/utils/vouchers";
import { VOUCHER_TEMPLATES, getVoucherTemplate, type VoucherTemplate } from "@/utils/voucherTemplates";
import { normalizeShorthand, shorthandToSeconds, uptimeForRouter } from "@/utils/duration";
import { inkIcon } from "@/utils/inkIcons";

const store = useRoutersStore();
const route = useRoute();
const router = useRouter();
const { t, locale } = useI18n();

const STEP =
  "type-display inline-flex items-center justify-center size-7.5 shrink-0 rounded-full bg-ink text-primary text-[15px]";
const MAX_COUNT = 500;
const QUICK_COUNTS = [20, 50, 100, 200];
const NAME_LENGTHS = [4, 5, 6, 8];
const PREVIEW_SAMPLES = 12;

const CHARSET_OPTS = [
  { key: "alphanumeric" as const, label: "vouchers.charsMixed", letters: true, digits: true },
  { key: "letters" as const, label: "vouchers.charsLetters", letters: true, digits: false },
  { key: "digits" as const, label: "vouchers.charsDigits", letters: false, digits: true },
];

const profiles = ref<Record<string, string>[]>([]);
const profileMetas = ref<Record<string, ProfileMeta>>({});
const hotspotSettings = ref<HotspotSettings>({
  hotspotName: "",
  dnsName: "",
  currency: "",
});
const loading = ref(false);
const error = ref("");

// Used by the Business voucher template's QR code — points at the
// hotspot's own login page, since that DNS name is what devices on the
// hotspot network actually resolve.
const loginUrl = computed(() =>
  hotspotSettings.value.dnsName ? `http://${hotspotSettings.value.dnsName}/login` : "",
);

const form = ref({
  profile: "",
  count: 50,
  prefix: "",
  nameLength: 6,
  charsLetters: true,
  charsDigits: true,
  passwordMode: "same" as "same" | "random" | "fixed",
  fixedPassword: "",
  template: "modern" as VoucherTemplate["key"],
  limitUptimeRaw: "",
  limitBytesTotalValue: null as number | null,
  limitBytesTotalUnit: "M" as "M" | "G",
  comment: "",
});

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
    if (!profiles.value.some((x) => x.name === form.value.profile)) {
      form.value.profile = profiles.value[0]?.name ?? "";
    }

    // A `?profile=` query param (e.g. linked from the Profiles page) opens
    // the wizard pre-filled for that profile, then clears the param so
    // back/forward navigation doesn't keep reopening it.
    const queryProfile = route.query.profile;
    if (typeof queryProfile === "string" && queryProfile) {
      form.value.profile = queryProfile;
      router.replace({ query: { ...route.query, profile: undefined } });
    }
  } catch (e: any) {
    error.value = friendlyError(e, t("vouchers.failedToLoad"));
  } finally {
    loading.value = false;
  }
}
watch(() => store.activeId, load, { immediate: true });

function setCount(n: number) {
  form.value.count = Math.min(MAX_COUNT, Math.max(1, Math.round(n) || 1));
}

// No 0/o, 1/l/i: codes are read off paper tickets and typed on phones.
const charset = computed(() =>
  [
    form.value.charsLetters ? "abcdefghjkmnpqrstuvwxyz" : "",
    form.value.charsDigits ? "23456789" : "",
  ].join(""),
);

const charsetKey = computed(() => {
  if (form.value.charsLetters && form.value.charsDigits) return "alphanumeric";
  if (form.value.charsDigits) return "digits";
  return "letters";
});

function setCharset(key: (typeof CHARSET_OPTS)[number]["key"]) {
  const opt = CHARSET_OPTS.find((o) => o.key === key)!;
  form.value.charsLetters = opt.letters;
  form.value.charsDigits = opt.digits;
}

const uptimePreview = computed(() => uptimeForRouter(form.value.limitUptimeRaw));
const uptimeWarning = computed(() => {
  const validity = profileMetas.value[form.value.profile]?.validity;
  if (!uptimePreview.value || !validity) return false;
  return shorthandToSeconds(form.value.limitUptimeRaw) > shorthandToSeconds(validity);
});

const canGenerate = computed(
  () =>
    !!form.value.profile &&
    form.value.count >= 1 &&
    form.value.count <= MAX_COUNT &&
    !!charset.value &&
    (!form.value.limitUptimeRaw || !!uptimePreview.value) &&
    (form.value.passwordMode !== "fixed" || !!form.value.fixedPassword),
);

function formatMoney(amount: number): string {
  const currency = hotspotSettings.value.currency;
  return `${new Intl.NumberFormat(locale.value).format(amount)}${currency ? " " + currency : ""}`;
}

function priceLabel(profile: string): string {
  const price = parseFloat(profileMetas.value[profile]?.price ?? "");
  return price ? formatMoney(price) : t("vouchers.noPrice");
}

const totalLabel = computed(() => {
  const price = parseFloat(profileMetas.value[form.value.profile]?.price ?? "");
  return price ? formatMoney(price * form.value.count) : "";
});

const currentTemplate = computed(() => getVoucherTemplate(form.value.template));

// ── Live preview: the real print template, fed with sample codes ─────────
const previewHtml = ref("");
let previewRun = 0;

watch(
  () => [
    form.value.profile,
    form.value.count,
    form.value.prefix,
    form.value.nameLength,
    charset.value,
    form.value.passwordMode,
    form.value.fixedPassword,
    form.value.template,
    uptimePreview.value,
    hotspotSettings.value,
    profileMetas.value,
  ],
  async () => {
    const run = ++previewRun;
    const used = new Set<string>();
    const samples: VoucherEntry[] = Array.from(
      { length: Math.min(PREVIEW_SAMPLES, form.value.count) },
      () => {
        const name = generateName(charset.value || "abc", form.value.nameLength, used);
        used.add(name);
        return {
          name,
          password: generatePassword(name),
          profile: form.value.profile,
          timeLimit: uptimePreview.value,
        };
      },
    );
    const html = await currentTemplate.value.render(
      toVoucherItems(samples, {
        showPrice: true,
        currency: hotspotSettings.value.currency ?? "",
        profileMetas: profileMetas.value,
      }),
      {
        businessName: hotspotSettings.value.hotspotName ?? "",
        // A sample QR when the hotspot has no DNS name yet, so the preview still shows the layout.
        loginUrl: loginUrl.value || "http://hotspot.local/login",
        loginUrlSupportsCredentials: true,
      },
    );
    if (run === previewRun) previewHtml.value = html;
  },
  { immediate: true },
);

// ── Generation ───────────────────────────────────────────────────────────
const showBatch = ref(false);
const batchRunning = ref(false);
const batchProgress = ref(0);
const batchTotal = ref(0);
const batchCurrentName = ref("");
const batchError = ref("");
const batchResults = ref<
  { name: string; password: string; ok: boolean; error?: string }[]
>([]);
const lastBatchProfile = ref("");
const lastBatchTimeLimit = ref("");

const createdCount = computed(() => batchResults.value.filter((r) => r.ok).length);
const failedCount = computed(() => batchResults.value.length - createdCount.value);

function generateName(
  charset: string,
  length: number,
  existing: Set<string>,
): string {
  let name = "",
    attempts = 0;
  do {
    name =
      form.value.prefix +
      Array.from(
        { length },
        () => charset[Math.floor(Math.random() * charset.length)],
      ).join("");
    attempts++;
  } while (existing.has(name) && attempts < 1000);
  return name;
}

function generatePassword(name: string): string {
  const { passwordMode, fixedPassword, nameLength } = form.value;
  if (passwordMode === "same") return name;
  if (passwordMode === "fixed") return fixedPassword;
  return Array.from(
    { length: nameLength },
    () => charset.value[Math.floor(Math.random() * charset.value.length)],
  ).join("");
}

function onBatchDialogUpdate(open: boolean) {
  if (!batchRunning.value) showBatch.value = open;
}

async function submitBatch(andPrint: boolean) {
  if (!store.activeId || !canGenerate.value) return;
  const cfg = form.value;
  showBatch.value = true;
  batchRunning.value = true;
  batchTotal.value = cfg.count;
  lastBatchProfile.value = cfg.profile;
  lastBatchTimeLimit.value = uptimePreview.value;
  batchProgress.value = 0;
  batchResults.value = [];
  batchError.value = "";
  const mul = cfg.limitBytesTotalUnit === "G" ? 1024 ** 3 : 1024 ** 2;
  const limitBytesTotal = cfg.limitBytesTotalValue
    ? String(cfg.limitBytesTotalValue * mul)
    : "";
  const usedNames = new Set<string>();
  for (let i = 0; i < cfg.count; i++) {
    const name = generateName(charset.value, cfg.nameLength, usedNames);
    usedNames.add(name);
    const password = generatePassword(name);
    batchCurrentName.value = name;
    try {
      await createHotspotUser(store.activeId, {
        name,
        password,
        profile: cfg.profile,
        limitUptime: uptimePreview.value,
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
  if (andPrint && createdCount.value) printResults();
}

function printResults() {
  const profileName = lastBatchProfile.value || "default";
  printVouchers(
    batchResults.value
      .filter((r) => r.ok)
      .map((r) => ({
        name: r.name,
        password: r.password,
        profile: profileName,
        timeLimit: lastBatchTimeLimit.value,
      })),
    {
      template: form.value.template,
      businessName: hotspotSettings.value.hotspotName ?? "",
      showPrice: true,
      currency: hotspotSettings.value.currency ?? "",
      profileMetas: profileMetas.value,
      loginUrl: loginUrl.value,
      loginUrlSupportsCredentials: true,
    },
  );
}
</script>
