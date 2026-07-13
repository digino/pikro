<template>
  <div class="flex flex-col items-center justify-center min-h-screen px-6 text-center bg-base">
    <div class="max-w-xl w-full">
      <h1 class="text-2xl font-bold tracking-tight text-text-primary mb-1">Pikro</h1>
      <p class="text-text-secondary mb-10">One app to rule your MikroTik.</p>

      <div v-if="store.routers.length === 0" class="border border-dashed rounded-xl p-8 mb-6">
        <ServerIcon class="size-7 text-text-muted mx-auto mb-3" />
        <p class="font-medium text-text-secondary mb-4">No routers configured yet</p>
        <button class="btn btn-primary" @click="router.push('/routers')">
          <PlusIcon class="size-4" />
          Add your first router
        </button>
      </div>

      <div v-else class="grid grid-cols-2 gap-3 text-left">
        <RouterLink
          v-for="item in nav"
          :key="item.to"
          :to="item.to"
          class="flex items-center gap-3 p-4 border border-border rounded-xl transition-colors hover:border-muted group bg-surface"
        >
          <div class="size-8 rounded-lg border border-border flex items-center justify-center transition-colors bg-base">
            <component :is="item.icon" class="size-4 text-text-muted" />
          </div>
          <div>
            <p class="font-semibold text-text-primary">{{ item.label }}</p>
            <p class="text-sm text-text-secondary">{{ item.desc }}</p>
          </div>
        </RouterLink>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { PlusIcon, ServerIcon, Squares2X2Icon, UsersIcon, BoltIcon, CircleStackIcon } from '@heroicons/vue/24/outline'
import { useRoutersStore } from "@/stores/routers";
import { useRouter } from "vue-router";

const store = useRoutersStore();
const router = useRouter();

const nav = [
  {
    to: "/dashboard",
    icon: Squares2X2Icon,
    label: "Dashboard",
    desc: "CPU, memory, uptime",
  },
  {
    to: "/hotspot",
    icon: UsersIcon,
    label: "Hotspot",
    desc: "Manage users & sessions",
  },
  {
    to: "/speedtest",
    icon: BoltIcon,
    label: "Speed Test",
    desc: "Bandwidth from router",
  },
  {
    to: "/routers",
    icon: CircleStackIcon,
    label: "Routers",
    desc: "Manage connections",
  },
];
</script>
