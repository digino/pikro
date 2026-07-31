<template>
  <span
    v-if="delta !== null"
    class="inline-flex items-center gap-0.5 text-xs font-semibold px-1.5 py-0.5 rounded-full"
    :class="isUp ? 'bg-green/10 text-green' : isDown ? 'bg-red/10 text-red' : 'bg-muted text-text-muted'"
  >
    <ArrowTrendingUpIcon v-if="isUp" class="size-3" />
    <ArrowTrendingDownIcon v-else-if="isDown" class="size-3" />
    {{ Math.abs(delta).toFixed(0) }}%
  </span>
</template>

<script setup lang="ts">
import { computed } from "vue";
import { ArrowTrendingUpIcon, ArrowTrendingDownIcon } from "@heroicons/vue/24/outline";

// delta is a percentage change vs. the previous period (e.g. 12.5 = +12.5%).
// null means there's no prior-period baseline to compare against.
const props = defineProps<{ delta: number | null }>();

const isUp = computed(() => (props.delta ?? 0) > 0);
const isDown = computed(() => (props.delta ?? 0) < 0);
</script>
