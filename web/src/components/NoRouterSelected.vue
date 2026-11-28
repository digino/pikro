<template>
  <div class="flex flex-col border border-border border-dashed rounded-lg items-center justify-center gap-4 py-16">
    <ServerIcon class="size-10 text-text-muted" />
    <div class="text-center">
      <p class="font-semibold text-text-primary">No router selected</p>
      <p class="text-sm text-text-secondary mt-1">
        Add a router or select one from the sidebar.
      </p>
    </div>
    <div class="flex items-center gap-2 mt-1">
      <button class="btn btn-ghost" @click="goScan">
        <MagnifyingGlassIcon class="size-3.5" />
        Scan network
      </button>
      <button class="btn btn-ink" @click="showAdd = true">
        <PlusIcon class="size-3.5" />
        Add router
      </button>
    </div>
  </div>

  <AddRouterDialog v-model:open="showAdd" @added="onAdded" />
</template>

<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";
import {
  ServerIcon,
  MagnifyingGlassIcon,
  PlusIcon,
} from "@heroicons/vue/24/outline";
import AddRouterDialog from "@/components/AddRouterDialog.vue";

const vueRouter = useRouter();
const showAdd = ref(false);

function goScan() {
  vueRouter.push({ path: '/routers', query: { scan: '1' } })
}

// store.add has already selected the new router, by id: routers often share
// an IP (every MikroTik ships as 192.168.88.1).
function onAdded() {
  showAdd.value = false;
}
</script>
