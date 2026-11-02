<template>
  <div v-if="installed !== null" class="card p-6 flex flex-col gap-4">
    <div class="flex items-center gap-3">
      <img :src="inkIcon('settings')" alt="" width="52" height="52" class="shrink-0" />
      <h2 class="type-display m-0 flex-1 text-2xl">{{ t('dashboard.autoCleanup.title') }}</h2>
      <SwitchRoot
        :model-value="installed === true"
        :disabled="toggling"
        :aria-label="t('dashboard.autoCleanup.title')"
        class="relative inline-flex items-center w-15 h-8.5 shrink-0 rounded-full border-2 border-ink shadow-[2px_2px_0_var(--color-shadow)] cursor-pointer transition-colors data-[state=checked]:bg-primary data-[state=unchecked]:bg-paper disabled:opacity-50 disabled:cursor-not-allowed focus-visible:outline-3 focus-visible:outline-primary focus-visible:outline-offset-2"
        @update:model-value="$emit('toggle', $event)"
      >
        <SwitchThumb
          class="pointer-events-none block size-6 rounded-full bg-ink transition-transform translate-x-0.75 data-[state=checked]:translate-x-7.25"
        />
      </SwitchRoot>
    </div>

    <p class="m-0 text-sm leading-normal text-text-secondary">
      {{ t('dashboard.autoCleanup.description') }}
    </p>

    <div v-if="installed">
      <div :id="intervalLabelId" class="eyebrow mb-2">{{ t('dashboard.autoCleanup.runsEvery') }}</div>
      <div
        role="group"
        :aria-labelledby="intervalLabelId"
        class="flex border-2 border-ink rounded-control overflow-hidden"
      >
        <button
          v-for="(value, i) in INTERVALS"
          :key="value"
          type="button"
          :aria-pressed="selected === value"
          :disabled="saving"
          class="flex-1 min-h-11 px-1 text-xs font-semibold cursor-pointer transition-colors disabled:cursor-not-allowed"
          :class="[
            i > 0 && 'border-l-2 border-ink',
            selected === value ? 'bg-primary text-primary-fg' : 'bg-paper text-ink hover:bg-primary-wash',
          ]"
          @click="selected !== value && $emit('update:interval', value)"
        >
          {{ t(`dashboard.autoCleanup.every.${value}`) }}
        </button>
      </div>
    </div>

    <div v-if="deviceModeBlocked" class="p-3 rounded-control border-2 border-ink bg-primary-wash space-y-2">
      <div class="flex items-start gap-2">
        <ExclamationTriangleIcon class="size-4 shrink-0 mt-0.5" />
        <div class="space-y-1">
          <p class="m-0 font-semibold">{{ t('dashboard.autoCleanup.deviceModeTitle') }}</p>
          <p class="m-0 text-sm text-text-secondary">{{ t('dashboard.autoCleanup.deviceModeBody') }}</p>
        </div>
      </div>
      <div class="border-2 border-ink rounded-lg px-3 py-2 font-mono text-xs select-all bg-paper">
        /system/device-mode/update scheduler=yes
      </div>
      <p class="m-0 text-xs text-text-muted">{{ t('dashboard.autoCleanup.deviceModeAfter') }}</p>
    </div>
    <p v-else-if="error" class="m-0 text-xs text-red">{{ error }}</p>

    <div class="mt-auto flex items-center gap-2.5 p-3 rounded-control bg-ground text-sm">
      <CheckIcon v-if="installed" class="size-4.5 shrink-0 text-green" />
      <ExclamationTriangleIcon v-else class="size-4.5 shrink-0" />
      <span>{{ installed ? t('dashboard.autoCleanup.statusOn') : t('dashboard.autoCleanup.statusOff') }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, useId } from 'vue'
import { useI18n } from 'vue-i18n'
import { SwitchRoot, SwitchThumb } from 'reka-ui'
import { CheckIcon, ExclamationTriangleIcon } from '@heroicons/vue/24/outline'
import { CLEANUP_INTERVAL_OPTIONS } from '@/utils/cleanupIntervals'
import { inkIcon } from '@/utils/inkIcons'
import { shorthandToSeconds } from '@/utils/duration'

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

const { t } = useI18n()
const intervalLabelId = useId()

const INTERVALS = CLEANUP_INTERVAL_OPTIONS.map(o => String(o.value))

// RouterOS reports the interval in its own form ("7d" comes back as "1w",
// "30d" as "4w2d"), so compare durations rather than strings.
const selected = computed(() => {
  const secs = shorthandToSeconds(props.interval)
  return INTERVALS.find(v => shorthandToSeconds(v) === secs)
})

const deviceModeBlocked = computed(() => props.error.toLowerCase().includes('device-mode'))
</script>
