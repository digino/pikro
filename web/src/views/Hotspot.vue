<template>
  <PageLayout title="Hotspot" subtitle="Captive portal management">
    <div v-if="!store.activeId" class="border border-dashed border-border rounded-xl py-16 text-center">
      <ServerIcon class="size-8 text-text-muted mx-auto mb-3" />
      <p class="text-sm font-medium text-text-secondary">No router selected</p>
      <p class="text-xs text-text-muted mt-1">Select a router from the sidebar dropdown</p>
    </div>

    <RouterView v-else />
  </PageLayout>
</template>

<script setup lang="ts">
import { ServerIcon } from '@heroicons/vue/24/outline'
import { useRoutersStore } from '@/stores/routers'
import PageLayout from '@/components/PageLayout.vue'

const store = useRoutersStore()
</script>
