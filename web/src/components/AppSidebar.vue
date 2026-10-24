<template>
  <aside
    class="w-64 min-w-55 h-screen bg-sidebar text-sidebar-text border-r-2 border-ink flex flex-col gap-5.5 px-4 py-5.5 sticky top-0 overflow-y-auto"
  >
    <!-- Logo -->
    <RouterLink to="/dashboard" class="flex items-center gap-2.5 px-1.5 no-underline text-sidebar-text">
      <span
        class="type-display inline-flex items-center justify-center size-10 shrink-0 rounded-control border-2 border-sidebar-text bg-primary text-primary-fg text-2xl leading-none"
        aria-hidden="true"
        >p</span
      >
      <span class="type-display text-2xl leading-none">pikro</span>
      <span
        class="ml-auto max-w-20 truncate rounded-full border-[1.5px] border-sidebar-text px-2 py-0.5 text-xs font-semibold"
        :title="appVersion"
        >{{ appVersion }}</span
      >
    </RouterLink>

    <!-- Nav sections -->
    <nav :aria-label="t('sidebar.mainNav')" class="flex flex-col gap-5.5">
      <div v-for="section in navSections" :key="section.label" class="flex flex-col gap-0.5">
        <div class="px-2.5 pb-1.5 text-[11px] font-semibold uppercase tracking-widest text-sidebar-muted">
          {{ t(section.label) }}
        </div>
        <RouterLink
          v-for="item in section.items"
          :key="item.to"
          :to="item.to"
          class="flex items-center gap-3 min-h-11 px-2.5 py-1.25 rounded-control border-2 text-sm no-underline transition-colors"
          :class="
            isActive(item.to)
              ? 'bg-primary text-primary-fg font-bold border-sidebar-text shadow-[3px_3px_0_var(--color-sidebar-text)]'
              : 'font-medium text-sidebar-text border-transparent hover:bg-white/8'
          "
        >
          <img :src="inkIcon(item.icon)" alt="" width="30" height="30" class="shrink-0" />
          <span>{{ t(item.label) }}</span>
        </RouterLink>
      </div>
    </nav>

    <!-- Router switcher, language and theme -->
    <div class="mt-auto flex flex-col gap-2.5">
      <SelectRoot
        :model-value="store.activeId ?? ''"
        @update:model-value="onSelectRouter"
      >
        <SelectTrigger
          :aria-label="t('sidebar.switchRouter')"
          class="flex items-center gap-2.5 w-full p-3 text-left rounded-2xl border-[1.5px] border-sidebar-line text-sidebar-text cursor-pointer transition-colors hover:bg-white/8 focus-visible:outline-3 focus-visible:outline-primary focus-visible:outline-offset-2"
        >
          <span class="flex flex-col gap-0.5 min-w-0 flex-1">
            <span class="truncate text-sm font-semibold">
              {{ activeRouter?.name ?? t('sidebar.noRouterSelected') }}
            </span>
            <span v-if="activeRouter" class="truncate font-mono text-xs text-sidebar-muted">
              {{ activeRouter.host }}
            </span>
          </span>
          <ChevronUpDownIcon class="size-4 shrink-0 text-sidebar-muted" />
        </SelectTrigger>
        <SelectPortal>
          <SelectContent
            class="z-50 min-w-56 bg-paper text-ink border-2 border-ink rounded-control shadow-btn overflow-hidden"
            position="popper"
            side="top"
            :side-offset="6"
          >
            <SelectViewport class="p-1">
              <div
                v-if="store.routers.length === 0"
                class="px-3 py-2 text-sm text-text-muted"
              >
                {{ t('sidebar.noRoutersConfigured') }}
              </div>
              <SelectItem
                v-for="r in store.routers"
                :key="r.id"
                :value="r.id"
                class="flex items-center justify-between gap-3 min-h-11 px-3 text-sm rounded-lg cursor-pointer outline-none data-highlighted:bg-primary-wash data-[state=checked]:font-semibold"
              >
                <SelectItemText>{{ r.name }}</SelectItemText>
                <SelectItemIndicator>
                  <CheckIcon class="size-4" />
                </SelectItemIndicator>
              </SelectItem>
            </SelectViewport>
          </SelectContent>
        </SelectPortal>
      </SelectRoot>

      <div class="flex gap-2">
        <button
          type="button"
          class="flex-1 h-11 rounded-control border-[1.5px] border-sidebar-line text-xs font-semibold text-sidebar-muted cursor-pointer hover:bg-white/8 transition-colors"
          :aria-label="t('sidebar.language')"
          :title="t('sidebar.language')"
          @click="toggleLocale"
        >
          <span :class="{ 'text-sidebar-text': localeStore.locale === 'en' }">EN</span>
          ·
          <span :class="{ 'text-sidebar-text': localeStore.locale === 'fr' }">FR</span>
        </button>
        <button
          type="button"
          class="w-11 h-11 inline-flex items-center justify-center rounded-control border-[1.5px] border-sidebar-line text-sidebar-text cursor-pointer hover:bg-white/8 transition-colors"
          :aria-label="
            themeStore.theme === 'dark'
              ? t('sidebar.switchToLight')
              : t('sidebar.switchToDark')
          "
          :title="
            themeStore.theme === 'dark'
              ? t('sidebar.switchToLight')
              : t('sidebar.switchToDark')
          "
          @click="themeStore.toggle()"
        >
          <SunIcon v-if="themeStore.theme === 'dark'" class="size-4.5" />
          <MoonIcon v-else class="size-4.5" />
        </button>
      </div>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { useRoute, useRouter } from "vue-router";
import { computed } from "vue";
import {
  SelectRoot,
  SelectTrigger,
  SelectValue,
  SelectPortal,
  SelectContent,
  SelectViewport,
  SelectItem,
  SelectItemText,
  SelectItemIndicator,
} from "reka-ui";
import {
  ChevronUpDownIcon,
  CheckIcon,
  SunIcon,
  MoonIcon,
} from "@heroicons/vue/24/outline";
import { useI18n } from "vue-i18n";
import { useRoutersStore } from "@/stores/routers";
import { useThemeStore } from "@/stores/theme";
import { useLocaleStore } from "@/stores/locale";
import { inkIcon, type InkIcon } from "@/utils/inkIcons";

const { t } = useI18n();
const themeStore = useThemeStore();
const localeStore = useLocaleStore();

defineProps<{ appVersion: string }>();

const store = useRoutersStore();
const route = useRoute();
const router = useRouter();

const activeRouter = computed(() => store.active());

function toggleLocale() {
  localeStore.set(localeStore.locale === "en" ? "fr" : "en");
}

function isActive(to: string) {
  return route.path === to || route.path.startsWith(to + "/");
}

interface NavItem {
  to: string;
  icon: InkIcon;
  label: string;
}

interface NavSection {
  label: string;
  items: NavItem[];
}

const navSections: NavSection[] = [
  {
    label: "nav.general",
    items: [
      { to: "/dashboard", icon: "dashboard", label: "nav.dashboard" },
      { to: "/routers", icon: "routers", label: "nav.routers" },
    ],
  },
  {
    label: "nav.hotspot",
    items: [
      { to: "/hotspot/users", icon: "users", label: "nav.users" },
      { to: "/hotspot/vouchers", icon: "vouchers", label: "nav.vouchers" },
      { to: "/hotspot/profiles", icon: "profiles", label: "nav.profiles" },
      { to: "/hotspot/reports", icon: "reports", label: "nav.reports" },
      { to: "/hotspot/logs", icon: "logs", label: "nav.logs" },
      { to: "/hotspot/settings", icon: "settings", label: "nav.settings" },
    ],
  },
  {
    label: "nav.network",
    items: [
      { to: "/network/monitor", icon: "monitor", label: "nav.monitor" },
      { to: "/network/hosts", icon: "hosts", label: "nav.hosts" },
      { to: "/network/dhcp", icon: "dhcp", label: "nav.dhcp" },
    ],
  },
  {
    label: "nav.tools",
    items: [{ to: "/speedtest", icon: "speed-test", label: "nav.speedTest" }],
  },
];

function onSelectRouter(id: string) {
  store.select(id);
  router.push("/dashboard");
}
</script>
