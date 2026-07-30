<template>
  <RouterHap
    v-if="family === 'hap' || family === 'generic'"
    :size="size"
    :power-led="powerLed"
    :wifi-led="wifiLed"
    :wan-led="wanLed"
  />
  <RouterHex
    v-else-if="family === 'hex'"
    :size="size"
    :power-led="powerLed"
    :wan-led="wanLed"
  />
  <RouterRack
    v-else-if="family === 'rack' || family === 'ccr'"
    :size="size"
    :power-led="powerLed"
    :act-led="wifiLed"
  />
  <RouterOutdoor
    v-else-if="family === 'outdoor'"
    :size="size"
    :power-led="powerLed"
  />
  <RouterHex
    v-else
    :size="size"
    :power-led="powerLed"
    :wan-led="wanLed"
  />
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { boardFamily, type BoardFamily } from '@/utils/boardFamily'
import RouterHap     from './RouterHap.vue'
import RouterHex     from './RouterHex.vue'
import RouterRack    from './RouterRack.vue'
import RouterOutdoor from './RouterOutdoor.vue'

const props = withDefaults(defineProps<{
  boardName: string
  size?: number
  powerLed?: string
  wifiLed?: string
  wanLed?: string
}>(), {
  size: 120,
  powerLed: 'var(--color-green)',
  wifiLed:  'var(--color-green)',
  wanLed:   'var(--color-amber)',
})

const family = computed<BoardFamily>(() => boardFamily(props.boardName))
</script>
