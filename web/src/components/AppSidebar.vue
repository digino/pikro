<template>
  <aside
    class="w-64 min-w-55 h-screen bg-sidebar border-r border-border flex flex-col sticky top-0 overflow-y-auto"
  >
    <!-- Logo -->
    <div class="flex items-center gap-2 px-4 pt-4 pb-3.5">
      <span class="text-lg font-bold tracking-tight text-text-primary">pikro</span>
    </div>

    <!-- Router selector -->
    <div class="px-3 pb-3">
      <SelectRoot :model-value="store.activeId ?? ''" @update:model-value="onSelectRouter">
        <SelectTrigger
          class="flex items-center gap-1.5 w-full p-2 text-sm font-medium bg-surface border border-border rounded-lg text-text-primary cursor-pointer transition-colors hover:bg-muted hover:border-muted focus-visible:outline-2 focus-visible:outline-accent focus-visible:outline-offset-1"
        >
          <SelectValue placeholder="No router selected" class="flex-1 truncate text-left" />
          <ChevronDownIcon class="size-4 text-text-secondary shrink-0" />
        </SelectTrigger>
        <SelectPortal>
          <SelectContent
            class="z-50 min-w-48 bg-surface border border-border rounded-lg shadow-xl overflow-hidden"
            position="popper"
            :side-offset="3"
          >
            <SelectViewport class="p-1">
              <div v-if="store.routers.length === 0" class="px-3 py-2 text-sm text-text-muted">
                No routers configured
              </div>
              <SelectItem
                v-for="r in store.routers"
                :key="r.id"
                :value="r.id"
                class="flex items-center justify-between p-2 text-sm rounded-md cursor-pointer text-text-secondary transition-colors hover:bg-muted hover:text-text-primary data-highlighted:bg-muted data-highlighted:text-text-primary data-[state=checked]:text-text-primary data-[state=checked]:font-medium"
              >
                <SelectItemText>{{ r.name }}</SelectItemText>
                <SelectItemIndicator>
                  <CheckCircleIcon class="size-4 text-green" />
                </SelectItemIndicator>
              </SelectItem>
            </SelectViewport>
          </SelectContent>
        </SelectPortal>
      </SelectRoot>
    </div>

    <!-- Nav sections -->
    <nav class="flex flex-col px-2">
      <template v-for="(section, i) in navSections" :key="section.label">
        <div v-if="i > 0" class="h-px bg-border mx-1 my-2" />
        <div>
          <div class="px-2 pb-2 text-sm font-medium text-text-secondary">
            {{ section.label }}
          </div>
          <div class="flex flex-col gap-px">
            <RouterLink
              v-for="item in section.items"
              :key="item.to"
              :to="item.to"
              class="flex items-center font-semibold gap-2 p-1.75 rounded-lg transition-colors"
              :class="isActive(item.to)
                ? 'bg-muted text-text-primary'
                : 'text-text-secondary hover:bg-surface hover:text-text-primary'"
            >
              <component :is="item.icon" class="size-5 shrink-0" />
              <span>{{ item.label }}</span>
            </RouterLink>
          </div>
        </div>
      </template>
    </nav>

    <div class="flex-1" />

    <!-- Footer -->
    <div class="px-4 py-3 border-t border-border">
      <span class="text-[11px] text-text-muted font-mono">{{ appVersion }}</span>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router'
import { type FunctionalComponent } from 'vue'
import {
  SelectRoot, SelectTrigger, SelectValue, SelectPortal,
  SelectContent, SelectViewport, SelectItem, SelectItemText, SelectItemIndicator,
} from 'reka-ui'
import {
  ChevronDownIcon, CheckCircleIcon,
  Squares2X2Icon, CircleStackIcon,
  BoltIcon, WrenchScrewdriverIcon, UsersIcon, RectangleGroupIcon, Cog6ToothIcon,
} from '@heroicons/vue/24/outline'
import { useRoutersStore } from '@/stores/routers'

defineProps<{ appVersion: string }>()

const store = useRoutersStore()
const route = useRoute()
const router = useRouter()

function isActive(to: string) {
  return route.path === to || route.path.startsWith(to + '/')
}

interface NavItem {
  to: string
  icon: FunctionalComponent
  label: string
}

interface NavSection {
  label: string
  items: NavItem[]
}

const navSections: NavSection[] = [
  {
    label: 'General',
    items: [
      { to: '/dashboard', icon: Squares2X2Icon, label: 'Dashboard' },
      { to: '/routers',   icon: CircleStackIcon, label: 'Routers' },
    ],
  },
  {
    label: 'Hotspot',
    items: [
      { to: '/hotspot/setup',     icon: WrenchScrewdriverIcon, label: 'Setup' },
      { to: '/hotspot/users',     icon: UsersIcon,          label: 'Users' },
      { to: '/hotspot/profiles',  icon: RectangleGroupIcon, label: 'Profiles' },
      { to: '/hotspot/settings',  icon: Cog6ToothIcon,      label: 'Settings' },
    ],
  },
  {
    label: 'Tools',
    items: [
      { to: '/speedtest', icon: BoltIcon, label: 'Speed Test' },
    ],
  },
]

function onSelectRouter(id: string) {
  store.select(id)
  router.push('/dashboard')
}
</script>
