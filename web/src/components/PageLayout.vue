<template>
  <div class="flex flex-col gap-6 min-h-full px-4 sm:px-6 lg:px-10 pt-7 pb-14">
    <!-- Topbar -->
    <header class="flex flex-wrap items-end justify-between gap-4">
      <div class="flex items-center gap-4 min-w-0">
        <img v-if="icon" :src="inkIcon(icon)" alt="" width="50" height="50" class="shrink-0" />
        <div class="min-w-0">
          <nav v-if="subtitle" :aria-label="title" class="eyebrow">{{ title }} / {{ subtitle }}</nav>
          <h1 class="type-display mt-1 text-3xl truncate">{{ subtitle ?? title }}</h1>
        </div>
      </div>
      <div class="flex flex-wrap items-center gap-3">
        <slot name="actions" />
      </div>
    </header>

    <!-- Page body -->
    <div class="flex-1 min-w-0 flex flex-col gap-6">
      <slot />
    </div>
  </div>
</template>

<script setup lang="ts">
import { inkIcon, type InkIcon } from "@/utils/inkIcons";

defineProps<{ title: string; subtitle?: string; icon?: InkIcon }>();
</script>
