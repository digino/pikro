<template>
  <div class="grid grid-cols-1 xl:grid-cols-[minmax(0,1.3fr)_minmax(0,1fr)] gap-7 items-start">
    <section class="card p-7 flex flex-col gap-4.5">
      <div class="flex flex-wrap items-center gap-3">
        <h2 class="type-display m-0 text-3xl">{{ t('settings.migration.heading') }}</h2>
        <span class="chip bg-primary text-primary-fg">{{ t('settings.beta') }}</span>
      </div>
      <p class="m-0 text-[15px] leading-relaxed text-text-secondary max-w-155">{{ t('settings.migration.intro') }}</p>
      <ul class="m-0 p-0 list-none flex flex-col gap-2.5 text-sm leading-snug">
        <li v-for="(step, i) in tm('settings.migration.steps')" :key="i" class="flex items-start gap-2.5">
          <CheckIcon class="size-5 shrink-0 mt-px stroke-[2.75]" aria-hidden="true" />
          <span>{{ rt(step) }}</span>
        </li>
      </ul>
      <p class="m-0 text-[13px] font-semibold">{{ t('settings.migration.safe') }}</p>
    </section>

    <section class="card p-6 flex flex-col gap-4">
      <div>
        <h3 class="m-0 font-bold">{{ t('settings.migration.actionTitle') }}</h3>
        <p class="m-0 mt-1 text-sm text-text-secondary">{{ t('settings.migration.actionBody') }}</p>
      </div>

      <template v-if="result">
        <div class="flex items-center gap-2.5 p-3.5 rounded-control border-2 border-ink bg-online-wash font-semibold">
          <CheckIcon class="size-5 shrink-0 stroke-3 text-online" aria-hidden="true" />
          <span>
            {{ t('settings.migration.done') }}
            <template v-if="result.cleanupInstalled">{{ t('settings.migration.cleanupOn') }}</template>
          </span>
        </div>
        <p v-if="!result.cleanupInstalled" role="alert" class="m-0 text-sm text-danger-text">
          {{ t('settings.migration.cleanupFailed') }}
        </p>
        <dl class="m-0 grid grid-cols-2 gap-2.5">
          <div v-for="s in stats" :key="s.label" class="rounded-control border-2 border-ink p-3">
            <dd class="m-0 type-display text-[28px] leading-none">{{ s.value }}</dd>
            <dt class="mt-1 text-xs">{{ s.label }}</dt>
          </div>
        </dl>
        <p class="m-0 text-xs text-text-muted">
          {{ t('settings.migration.scanned', { users: result.usersScanned, profiles: result.profilesScanned, skipped: result.usersSkipped }) }}
        </p>
      </template>
      <p v-else-if="error" role="alert" class="m-0 text-sm text-danger-text">{{ error }}</p>

      <button
        type="button"
        class="btn min-h-12.5 text-[15px]"
        :class="result ? '' : 'btn-primary'"
        :disabled="migrating"
        @click="confirming = true"
      >
        <span v-if="migrating" class="spinner spinner--sm" />
        {{ migrating ? t('settings.migration.running') : result ? t('settings.migration.runAgain') : t('settings.migration.run') }}
      </button>
    </section>

    <AppDialog v-model:open="confirming" :title="t('settings.migration.confirmTitle')">
      <p class="m-0 text-sm leading-normal text-text-secondary">{{ t('settings.migration.confirmBody') }}</p>
      <div class="flex justify-end gap-2.5 mt-6">
        <button type="button" class="btn" @click="confirming = false">{{ t('settings.migration.cancel') }}</button>
        <button type="button" class="btn btn-primary" @click="runMigration">{{ t('settings.migration.run') }}</button>
      </div>
    </AppDialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { CheckIcon } from '@heroicons/vue/24/outline'
import { useRoutersStore } from '@/stores/routers'
import { migrateFromMikhmon, type MikhmonMigrationResult } from '@/api'
import AppDialog from '@/components/AppDialog.vue'

const { t, tm, rt } = useI18n()
const store = useRoutersStore()

const confirming = ref(false)
const migrating = ref(false)
const result = ref<MikhmonMigrationResult | null>(null)
const error = ref('')

const stats = computed(() => {
  const r = result.value
  if (!r) return []
  return [
    { label: t('settings.migration.usersConverted'), value: r.usersConverted },
    { label: t('settings.migration.usersUnused'), value: r.usersUnused },
    { label: t('settings.migration.profilesConverted'), value: r.profilesConverted },
    { label: t('settings.migration.leftoversRemoved'), value: r.schedulersRemoved + r.scriptsRemoved },
  ]
})

async function runMigration() {
  confirming.value = false
  if (!store.activeId) return
  migrating.value = true
  error.value = ''
  result.value = null
  try {
    result.value = await migrateFromMikhmon(store.activeId)
  } catch (e: any) {
    error.value = e?.response?.data?.error ?? e?.message ?? t('settings.migration.failed')
  } finally {
    migrating.value = false
  }
}
</script>
