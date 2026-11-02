<template>
  <div
    v-if="updates.showBanner && updates.status"
    role="status"
    class="flex flex-wrap items-center gap-3 px-4 sm:px-6 lg:px-10 py-2.5 bg-primary text-primary-fg border-b-2 border-ink text-sm font-semibold"
  >
    <span class="flex-1 min-w-0">
      {{ t('update.available', { version: displayVersion(updates.status.latest) }) }}
    </span>
    <a
      :href="downloadPage(locale)"
      target="_blank"
      rel="noopener noreferrer"
      class="btn btn-sm btn-ink"
    >
      {{ t('update.download') }}
    </a>
    <button
      type="button"
      class="inline-flex items-center justify-center size-9 rounded-control cursor-pointer hover:bg-black/10"
      :aria-label="t('update.dismiss')"
      :title="t('update.dismiss')"
      @click="updates.dismiss()"
    >
      <XMarkIcon class="size-5" />
    </button>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { XMarkIcon } from '@heroicons/vue/24/outline'
import { useUpdateStore, displayVersion } from '@/stores/update'
import { downloadPage } from '@/utils/links'

const { t, locale } = useI18n()
const updates = useUpdateStore()
</script>
