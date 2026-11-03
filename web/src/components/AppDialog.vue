<template>
  <DialogRoot :open="open" @update:open="$emit('update:open', $event)">
    <DialogPortal>
      <!-- Backdrop -->
      <DialogOverlay class="fixed inset-0 bg-black/40 backdrop-blur-sm z-40 data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0" />

      <!-- Panel -->
      <DialogContent
        class="card fixed left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 z-50 w-[calc(100%-2rem)] max-h-[calc(100vh-2rem)] overflow-y-auto p-6 focus:outline-none data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95"
        :class="size === 'xl' ? 'max-w-3xl' : 'max-w-lg'"
      >
        <div class="flex items-start justify-between mb-5">
          <DialogTitle class="type-display m-0 text-xl text-text-primary">{{ title }}</DialogTitle>
          <DialogClose
            :aria-label="t('common.close')"
            class="ml-4 -mt-2 -mr-2 size-11 shrink-0 inline-flex items-center justify-center rounded-control text-text-muted transition-colors hover:bg-ground hover:text-ink"
          >
            <XMarkIcon class="size-5" />
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
import { useI18n } from 'vue-i18n'

defineProps<{ open: boolean; title: string; size?: 'md' | 'xl' }>()
defineEmits<{ 'update:open': [value: boolean] }>()

const { t } = useI18n()
</script>
