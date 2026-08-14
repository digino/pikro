<template>
  <div v-if="installed !== null" class="p-3 rounded-lg border border-border bg-base space-y-3">
    <div class="flex items-center gap-3">
      <component
        :is="installed ? CheckCircleIcon : ExclamationTriangleIcon"
        class="size-6 shrink-0"
        :class="installed ? 'text-green' : 'text-amber'"
      />
      <div class="flex-1 min-w-0">
        <div class="font-semibold text-text-primary">
          Auto-cleanup
          {{ installed ? "active" : "not configured" }}
        </div>
        <div class="text-sm text-text-secondary mt-0.5">
          {{
            installed
              ? `Expired vouchers are removed automatically every ${intervalLabel}.`
              : "Expired users accumulate until removed manually."
          }}
        </div>
      </div>
      <SwitchRoot
        :model-value="installed === true"
        :disabled="toggling"
        class="relative inline-flex h-5 w-9 shrink-0 items-center rounded-full transition-colors data-[state=checked]:bg-green data-[state=unchecked]:bg-border disabled:opacity-50"
        @update:model-value="$emit('toggle', $event)"
      >
        <SwitchThumb
          class="pointer-events-none block size-4 rounded-full bg-white shadow transform transition-transform translate-x-0.5 data-[state=checked]:translate-x-4"
        />
      </SwitchRoot>
    </div>

    <div v-if="installed" class="flex items-center justify-end gap-2">
      <span>Run every</span>
      <div class="shrink-0">
        <AppSelect
          :model-value="interval"
          :options="INTERVAL_OPTIONS"
          :disabled="saving"
          @update:model-value="$emit('update:interval', String($event))"
        />
      </div>
    </div>

    <div v-if="deviceModeBlocked" class="p-3 border rounded-lg space-y-2 bg-amber/8 border-amber/20">
      <div class="flex items-start gap-2">
        <ExclamationTriangleIcon class="size-4 shrink-0 mt-0.5 text-amber" />
        <div class="space-y-1">
          <p class="font-semibold text-amber">Scheduler blocked by RouterOS device-mode</p>
          <p class="text-sm text-text-secondary">
            Your router runs RouterOS 7.17+ with the scheduler disabled. Run the command below in
            Winbox Terminal or SSH, then press the physical button on your router (or cold-reboot) to
            confirm.
          </p>
        </div>
      </div>
      <div class="border border-border rounded-lg px-3 py-2 font-mono text-xs text-text-primary select-all bg-surface">
        /system/device-mode/update scheduler=yes
      </div>
      <p class="text-xs text-text-muted">After confirming, re-run hotspot setup to install the scheduler.</p>
    </div>
    <p v-else-if="error" class="text-xs text-red pl-9">{{ error }}</p>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { SwitchRoot, SwitchThumb } from 'reka-ui'
import { CheckCircleIcon, ExclamationTriangleIcon } from '@heroicons/vue/24/outline'
import AppSelect from '@/components/AppSelect.vue'
import { CLEANUP_INTERVAL_OPTIONS as INTERVAL_OPTIONS } from '@/utils/cleanupIntervals'

const props = defineProps<{
  installed: boolean | null
  interval: string
  toggling: boolean
  saving: boolean
  error: string
}>()

defineEmits<{
  toggle: [enabled: boolean]
  'update:interval': [value: string]
}>()

const deviceModeBlocked = computed(() => props.error.toLowerCase().includes('device-mode'))

// Cleanup interval is stored as a router-shorthand duration (e.g. "7d", "12h").
// Expand it to a human phrase for the card copy.
const intervalLabel = computed(() => {
  const s = props.interval
  const m = s.match(/^(\d+)([wdhm])$/i)
  if (!m) return s || 'a regular interval'
  const n = parseInt(m[1])
  const unit = { w: 'week', d: 'day', h: 'hour', m: 'minute' }[m[2].toLowerCase()]!
  return `${n} ${unit}${n === 1 ? '' : 's'}`
})
</script>
