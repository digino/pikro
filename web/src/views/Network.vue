<template>
  <PageLayout v-if="!store.activeId" title="Network" subtitle="Traffic, hosts, and DHCP">
    <NoRouterSelected />
  </PageLayout>

  <RouterView v-else />
</template>

<script setup lang="ts">
import { useRoutersStore } from '@/stores/routers'
import PageLayout from '@/components/PageLayout.vue'
import NoRouterSelected from '@/components/NoRouterSelected.vue'

const store = useRoutersStore()
</script>
