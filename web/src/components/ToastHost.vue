<template>
  <ToastProvider swipe-direction="right">
    <ToastRoot
      v-for="t in store.toasts"
      :key="t.id"
      class="rounded-lg border p-3.5 pr-4 shadow-lg bg-surface flex items-start gap-2.5 transition-all"
      :class="borderClass(t.variant)"
      :duration="4000"
      @update:open="(open) => !open && store.dismiss(t.id)"
    >
      <component :is="iconFor(t.variant)" class="size-5 shrink-0 mt-0.5" :class="iconClass(t.variant)" />
      <div class="min-w-0 flex-1">
        <ToastTitle class="text-sm font-semibold text-text-primary">{{ t.title }}</ToastTitle>
        <ToastDescription v-if="t.description" class="text-xs text-text-secondary mt-0.5">
          {{ t.description }}
        </ToastDescription>
      </div>
      <ToastClose class="text-text-muted hover:text-text-primary transition-colors shrink-0 cursor-pointer">
        <XMarkIcon class="size-4" />
      </ToastClose>
    </ToastRoot>
    <ToastViewport class="fixed top-4 right-4 z-100 flex flex-col gap-2 w-full max-w-sm outline-none" />
  </ToastProvider>
</template>

<script setup lang="ts">
import {
  ToastProvider, ToastRoot, ToastTitle, ToastDescription, ToastClose, ToastViewport,
} from 'reka-ui'
import { CheckCircleIcon, ExclamationTriangleIcon, InformationCircleIcon, XMarkIcon } from '@heroicons/vue/24/outline'
import { useToastStore, type ToastVariant } from '@/stores/toast'

const store = useToastStore()

function iconFor(v: ToastVariant) {
  if (v === 'success') return CheckCircleIcon
  if (v === 'error') return ExclamationTriangleIcon
  return InformationCircleIcon
}
function iconClass(v: ToastVariant) {
  if (v === 'success') return 'text-green'
  if (v === 'error') return 'text-red'
  return 'text-text-secondary'
}
function borderClass(v: ToastVariant) {
  if (v === 'success') return 'border-green/30'
  if (v === 'error') return 'border-red/30'
  return 'border-border'
}
</script>
