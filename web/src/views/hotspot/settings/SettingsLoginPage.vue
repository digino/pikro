<template>
  <div
    class="grid grid-cols-1 xl:grid-cols-[minmax(0,1.3fr)_minmax(0,1fr)] gap-7 items-start"
  >
    <section class="card p-6 flex flex-col gap-5.5">
      <div>
        <h2 class="type-display m-0 text-2xl">
          {{ t("settings.loginPage.heading") }}
        </h2>
        <p class="m-0 mt-1.5 text-sm leading-normal text-text-secondary">
          {{ t("settings.loginPage.intro") }}
        </p>
      </div>

      <div class="flex flex-col gap-2">
        <span :id="modeLabelId" class="eyebrow">{{
          t("settings.loginPage.howGuestsLogIn")
        }}</span>
        <div
          role="group"
          :aria-labelledby="modeLabelId"
          class="grid grid-cols-1 sm:grid-cols-2 gap-3"
        >
          <button
            v-for="m in MODES"
            :key="m"
            type="button"
            :aria-pressed="form.mode === m"
            class="flex flex-col gap-1 p-3.5 text-left rounded-card border-2 border-ink cursor-pointer transition-colors"
            :class="
              form.mode === m
                ? 'bg-primary text-primary-fg shadow-btn'
                : 'bg-paper text-ink hover:bg-primary-wash'
            "
            @click="form.mode = m"
          >
            <span class="font-bold text-sm">{{
              t(`settings.loginPage.modes.${m}.label`)
            }}</span>
            <span class="text-xs leading-snug">{{
              t(`settings.loginPage.modes.${m}.note`)
            }}</span>
          </button>
        </div>
        <p class="m-0 text-xs text-text-muted">
          {{ t("settings.loginPage.modeHint") }}
        </p>
      </div>

      <div class="flex flex-col gap-2">
        <span :id="langLabelId" class="eyebrow">{{
          t("settings.loginPage.language")
        }}</span>
        <div
          role="group"
          :aria-labelledby="langLabelId"
          class="flex self-start border-2 border-ink rounded-control overflow-hidden bg-paper"
        >
          <button
            v-for="(l, i) in LANGUAGES"
            :key="l.key"
            type="button"
            :lang="l.key"
            :aria-pressed="form.language === l.key"
            class="min-h-10 px-4 text-sm font-semibold cursor-pointer"
            :class="[
              i > 0 && 'border-l-2 border-ink',
              form.language === l.key
                ? 'bg-primary text-primary-fg'
                : 'text-ink hover:bg-primary-wash',
            ]"
            @click="form.language = l.key"
          >
            {{ l.label }}
          </button>
        </div>
      </div>

      <div class="grid grid-cols-1 sm:grid-cols-2 gap-3.5">
        <label class="flex flex-col gap-1.5">
          <span class="eyebrow">{{ t("settings.loginPage.title") }}</span>
          <input
            v-model="form.title"
            class="input"
            :placeholder="settings.hotspotName || 'Wi-Fi'"
          />
          <span class="text-xs text-text-muted">{{
            t("settings.loginPage.titleHint")
          }}</span>
        </label>
        <label class="flex flex-col gap-1.5">
          <span class="eyebrow">{{ t("settings.loginPage.subtitle") }}</span>
          <input
            v-model="form.subtitle"
            class="input"
            placeholder="$(hostname)"
          />
          <span class="text-xs text-text-muted">{{
            t("settings.loginPage.subtitleHint")
          }}</span>
        </label>
      </div>

      <div class="flex flex-col gap-1.5">
        <div class="flex items-center justify-between gap-4">
          <label :for="pricesId" class="eyebrow">{{
            t("settings.loginPage.showPrices")
          }}</label>
          <SwitchRoot
            :id="pricesId"
            v-model="form.showPrices"
            class="relative inline-flex items-center w-15 h-8.5 shrink-0 rounded-full border-2 border-ink shadow-[2px_2px_0_var(--color-shadow)] cursor-pointer transition-colors data-[state=checked]:bg-primary data-[state=unchecked]:bg-paper focus-visible:outline-3 focus-visible:outline-primary focus-visible:outline-offset-2"
          >
            <SwitchThumb
              class="pointer-events-none block size-6 rounded-full bg-ink transition-transform translate-x-0.75 data-[state=checked]:translate-x-7.25"
            />
          </SwitchRoot>
        </div>
        <span class="text-xs text-text-muted">
          {{
            priceList.length
              ? t("settings.loginPage.showPricesHint")
              : t("settings.loginPage.noPrices")
          }}
        </span>
      </div>

      <label class="flex flex-col gap-1.5">
        <span class="eyebrow">{{ t("settings.loginPage.contact") }}</span>
        <input
          v-model="form.contact"
          class="input"
          type="tel"
          inputmode="tel"
          autocomplete="off"
          placeholder="+225 07 00 00 00 00"
        />
        <span class="text-xs text-text-muted">{{
          t("settings.loginPage.contactHint")
        }}</span>
      </label>

      <div
        class="flex items-start gap-2.5 p-3 rounded-control bg-ground text-sm leading-snug"
      >
        <img
          :src="inkIcon('routers')"
          alt=""
          width="28"
          height="28"
          class="shrink-0"
        />
        <span>
          {{ t("settings.loginPage.routerNote") }}
          <RouterLink to="/routers" class="font-semibold underline">{{
            t("settings.loginPage.routersLink")
          }}</RouterLink
          >.
        </span>
      </div>

      <p v-if="error" role="alert" class="m-0 text-sm text-danger-text">
        {{ error }}
      </p>

      <div class="flex flex-wrap items-center gap-2.5">
        <button
          type="button"
          class="btn btn-primary"
          :disabled="uploading"
          @click="confirming = true"
        >
          <span v-if="uploading" class="spinner spinner--sm" />
          {{ t("settings.loginPage.upload") }}
        </button>
        <button type="button" class="btn" :disabled="uploading" @click="reset">
          {{ t("settings.loginPage.reset") }}
        </button>
        <span class="text-[13px] text-text-muted sm:ml-auto">{{
          lastUpload
        }}</span>
      </div>
    </section>

    <section
      :aria-label="t('settings.loginPage.preview')"
      class="flex flex-col items-center gap-3"
    >
      <div v-if="showingLive" class="eyebrow">
        {{ t("settings.loginPage.livePreview") }}
      </div>
      <div
        class="w-[300px] h-[580px] p-2.5 rounded-[40px] border-[3px] border-black bg-black shadow-[8px_8px_0_var(--color-shadow)]"
      >
        <iframe
          :srcdoc="showingLive ? liveHTML : previewHtml(pages[previewPage])"
          class="w-full h-full rounded-[30px] border-0 bg-white"
          sandbox=""
          :title="t('settings.loginPage.preview')"
        />
      </div>
      <div
        v-if="!showingLive"
        role="group"
        :aria-label="t('settings.loginPage.previewPage')"
        class="flex border-2 border-ink rounded-control overflow-hidden bg-paper"
      >
        <button
          v-for="(p, i) in PREVIEW_PAGES"
          :key="p"
          type="button"
          :aria-pressed="previewPage === p"
          class="min-h-10 px-3.5 text-[13px] font-semibold cursor-pointer"
          :class="[
            i > 0 && 'border-l-2 border-ink',
            previewPage === p
              ? 'bg-ink text-paper'
              : 'text-ink hover:bg-primary-wash',
          ]"
          @click="previewPage = p"
        >
          {{ t(`settings.loginPage.pages.${p}`) }}
        </button>
      </div>
      <button
        type="button"
        class="btn btn-ghost btn-sm"
        :disabled="loadingLive"
        @click="toggleLive"
      >
        <span v-if="loadingLive" class="spinner spinner--sm" />
        {{
          showingLive
            ? t("settings.loginPage.backToPreview")
            : t("settings.loginPage.viewLive")
        }}
      </button>
      <p v-if="liveError" role="alert" class="m-0 text-xs text-danger-text">
        {{ liveError }}
      </p>
    </section>

    <AppDialog
      v-model:open="confirming"
      :title="t('settings.loginPage.confirmTitle')"
    >
      <p class="m-0 text-sm leading-normal text-text-secondary">
        {{ t("settings.loginPage.confirmBody") }}
      </p>
      <div class="flex justify-end gap-2.5 mt-6">
        <button type="button" class="btn" @click="confirming = false">
          {{ t("settings.loginPage.cancel") }}
        </button>
        <button type="button" class="btn btn-primary" @click="upload">
          {{ t("settings.loginPage.upload") }}
        </button>
      </div>
    </AppDialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, useId } from "vue";
import { useI18n } from "vue-i18n";
import { SwitchRoot, SwitchThumb } from "reka-ui";
import { useRoutersStore } from "@/stores/routers";
import { useToastStore } from "@/stores/toast";
import {
  getHotspotSettings,
  putHotspotSettings,
  uploadLoginPage,
  getLoginPageHTML,
  type HotspotSettings,
  type LoginPageSettings,
} from "@/api";
import AppDialog from "@/components/AppDialog.vue";
import { inkIcon } from "@/utils/inkIcons";
import {
  renderGuestPages,
  previewHtml,
  type LoginMode,
  type GuestLanguage,
  type GuestPages,
} from "./loginPageTemplates";

const props = defineProps<{ settings: HotspotSettings }>();
const emit = defineEmits<{ saved: [settings: HotspotSettings] }>();

const { t, locale } = useI18n();
const store = useRoutersStore();
const toast = useToastStore();

const modeLabelId = useId();
const langLabelId = useId();
const pricesId = useId();

const MODES: LoginMode[] = ["password", "code"];
const LANGUAGES: { key: GuestLanguage; label: string }[] = [
  { key: "en", label: "English" },
  { key: "fr", label: "Français" },
];
const PREVIEW_PAGES = ["login", "status", "logout"] as const;

// Pages saved before the login-mode choice only stored a design: the
// "voucher" design was the code-only one.
function fromSaved(lp: LoginPageSettings = {}) {
  return {
    mode:
      lp.loginMode ??
      ((lp.template && lp.template !== "voucher"
        ? "password"
        : "code") as LoginMode),
    language:
      lp.language ?? ((locale.value === "fr" ? "fr" : "en") as GuestLanguage),
    title: lp.title ?? "",
    subtitle: lp.subtitle ?? "",
    showPrices: lp.showPrices ?? false,
    contact: lp.contact ?? "",
  };
}

const form = ref(fromSaved(props.settings.loginPage));
watch(
  () => props.settings,
  (s) => {
    form.value = fromSaved(s.loginPage);
  },
);

function reset() {
  form.value = fromSaved(props.settings.loginPage);
}

const priceList = computed(() =>
  Object.entries(props.settings.profileMetas ?? {})
    .filter(([, m]) => Number(m.price) > 0)
    .sort(([, a], [, b]) => Number(a.price) - Number(b.price))
    .map(([name, m]) => ({
      name,
      price: `${m.price} ${props.settings.currency}`.trim(),
    })),
);

const pages = computed<GuestPages>(() =>
  renderGuestPages({
    title: form.value.title || props.settings.hotspotName || "Wi-Fi",
    subtitle: form.value.subtitle || "$(hostname)",
    mode: form.value.mode,
    language: form.value.language,
    prices: form.value.showPrices ? priceList.value : [],
    contact: form.value.contact.trim(),
  }),
);
const previewPage = ref<(typeof PREVIEW_PAGES)[number]>("login");

const lastUpload = computed(() => {
  const lp = props.settings.loginPage;
  if (!lp?.uploadedAt) return t("settings.loginPage.neverUploaded");
  const date = new Date(lp.uploadedAt).toLocaleDateString(locale.value, {
    day: "numeric",
    month: "short",
    year: "numeric",
  });
  return t("settings.loginPage.lastUpload", {
    date,
    mode: t(`settings.loginPage.modes.${fromSaved(lp).mode}.label`),
  });
});

const showingLive = ref(false);
const loadingLive = ref(false);
const liveError = ref("");
const liveHTML = ref("");

async function toggleLive() {
  if (showingLive.value) {
    showingLive.value = false;
    return;
  }
  if (!store.activeId) return;
  loadingLive.value = true;
  liveError.value = "";
  try {
    liveHTML.value = await getLoginPageHTML(store.activeId);
    showingLive.value = true;
  } catch (e: any) {
    liveError.value =
      e?.response?.data?.error ??
      e?.message ??
      t("settings.loginPage.liveFailed");
  } finally {
    loadingLive.value = false;
  }
}

const confirming = ref(false);
const uploading = ref(false);
const error = ref("");

async function upload() {
  confirming.value = false;
  if (!store.activeId) return;
  uploading.value = true;
  error.value = "";
  try {
    const p = pages.value;
    await uploadLoginPage(store.activeId, {
      html: p.login,
      statusHtml: p.status,
      logoutHtml: p.logout,
      redirectHtml: p.redirect,
    });
    const f = form.value;
    const loginPage: LoginPageSettings = {
      title: f.title.trim(),
      subtitle: f.subtitle.trim(),
      loginMode: f.mode,
      language: f.language,
      showPrices: f.showPrices,
      contact: f.contact.trim(),
      uploadedAt: new Date().toISOString(),
    };
    // Re-read so a Profiles or Routers change made meanwhile isn't overwritten.
    const latest = await getHotspotSettings(store.activeId);
    emit(
      "saved",
      await putHotspotSettings(store.activeId, { ...latest, loginPage }),
    );
    showingLive.value = false;
    toast.success(
      t("settings.loginPage.uploaded"),
      t("settings.loginPage.uploadedBody"),
    );
  } catch (e: any) {
    const message =
      e?.response?.data?.error ??
      e?.message ??
      t("settings.loginPage.uploadFailed");
    error.value = message;
    toast.error(t("settings.loginPage.uploadFailed"), message);
  } finally {
    uploading.value = false;
  }
}
</script>
