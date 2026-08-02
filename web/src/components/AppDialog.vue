<template>
  <DialogRoot :open="open" @update:open="$emit('update:open', $event)">
    <DialogPortal>
      <!-- Backdrop -->
      <DialogOverlay class="fixed inset-0 bg-black/40 backdrop-blur-sm z-40 data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0" />

      <!-- Panel -->
      <DialogContent
        class="fixed left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 z-50 w-full rounded-xl shadow-2xl p-6 focus:outline-none bg-surface border border-border data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95"
        :class="size === 'xl' ? 'max-w-3xl' : 'max-w-lg'"
      >
        <div class="flex items-start justify-between mb-5">
          <DialogTitle class="font-semibold text-text-primary" style="font-size: 1rem">{{ title }}</DialogTitle>
          <DialogClose class="ml-4 p-1 rounded-md transition-colors text-text-muted hover:bg-muted">
            <XMarkIcon class="size-4" />
          </DialogClose>
        </div>

        <slot />
      </DialogContent>
    </DialogPortal>
  </DialogRoot>
</template>

<script setup lang="ts">
import {
  DialogRoot, DialogPortal, DialogOverlay,
  DialogContent, DialogTitle, DialogClose,
} from 'reka-ui'
import { XMarkIcon } from '@heroicons/vue/20/solid'

defineProps<{ open: boolean; title: string; size?: 'md' | 'xl' }>()
defineEmits<{ 'update:open': [value: boolean] }>()
</script>
