<template>
  <SelectRoot
    :model-value="modelValue === '' || modelValue === undefined ? undefined : String(modelValue)"
    @update:model-value="onUpdate"
  >
    <SelectTrigger
      class="input-select flex items-center gap-1.5 cursor-pointer"
      :class="triggerClass"
    >
      <SelectValue :placeholder="placeholder" class="flex-1 text-left truncate" />
      <ChevronDownIcon class="size-3.5 text-text-secondary shrink-0" />
    </SelectTrigger>
    <SelectPortal>
      <SelectContent
        class="z-50 min-w-(--reka-select-trigger-width) bg-surface border border-border rounded-lg shadow-xl overflow-hidden"
        position="popper"
        :side-offset="3"
      >
        <SelectViewport class="p-1">
          <SelectItem
            v-for="opt in options"
            :key="opt.value"
            :value="String(opt.value)"
            class="flex items-center justify-between p-2 text-sm rounded-md cursor-pointer text-text-secondary transition-colors hover:bg-muted hover:text-text-primary data-highlighted:bg-muted data-highlighted:text-text-primary data-[state=checked]:text-text-primary data-[state=checked]:font-medium"
          >
            <SelectItemText>{{ opt.label }}</SelectItemText>
            <SelectItemIndicator><CheckCircleIcon class="size-4 text-green" /></SelectItemIndicator>
          </SelectItem>
        </SelectViewport>
      </SelectContent>
    </SelectPortal>
  </SelectRoot>
</template>

<script setup lang="ts">
import {
  SelectRoot,
  SelectTrigger,
  SelectValue,
  SelectPortal,
  SelectContent,
  SelectViewport,
  SelectItem,
  SelectItemText,
  SelectItemIndicator,
} from 'reka-ui'
import { ChevronDownIcon, CheckCircleIcon } from '@heroicons/vue/24/outline'

export interface SelectOption {
  value: string | number
  label: string
}

const props = defineProps<{
  modelValue: string | number | undefined
  options: SelectOption[]
  placeholder?: string
  /** When true, model-value/emitted values are coerced back to number. */
  numeric?: boolean
  triggerClass?: string
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string | number]
}>()

function onUpdate(value: string | undefined) {
  if (value === undefined) return
  emit('update:modelValue', props.numeric ? Number(value) : value)
}
</script>
