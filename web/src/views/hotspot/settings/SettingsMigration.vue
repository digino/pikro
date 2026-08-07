<template>
  <div class="max-w-xl space-y-5">
    <div>
      <div class="flex items-center gap-2">
        <h3 class="font-semibold text-text-primary">Migrate from Mikhmonv3</h3>
        <StatusBadge label="Beta" color="amber" variant="pill" />
      </div>
      <p class="text-sm text-text-secondary mt-0.5">
        If this router was previously managed with Mikhmon, its hotspot users and
        profiles use a different expiry format that Pikro can't read yet — they'll
        show as "Waiting" even if connected, and expiry columns stay blank.
      </p>
    </div>

    <div class="p-4 border border-border rounded-xl bg-surface space-y-2 font-medium">
      <p class="font-medium text-text-primary underline">What this does</p>
      <ul class="list-disc list-inside space-y-1 text-sm text-text-secondary">
        <li>Unused vouchers (never logged in) — comment cleared, same as a fresh Pikro voucher</li>
        <li>Already-activated users — expiry converted to Pikro's format (already-expired ones will be removed by auto-cleanup on its normal schedule)</li>
        <li>Hotspot profiles — on-login script replaced so future logins use Pikro's format</li>
        <li>Mikhmon's own leftover per-user schedulers and scripts (created by its old on-login script at every login) are removed</li>
        <li>Pikro's own cleanup scheduler is installed, so expired vouchers keep getting removed automatically</li>
        <li>Anything already in Pikro's format, or not recognized, is left untouched</li>
      </ul>
      <p class="text-sm">Safe to run more than once.</p>
    </div>

    <div v-if="result" class="p-4 border border-green/20 bg-green/8 rounded-xl space-y-1 text-sm">
      <p class="font-medium text-green">Migration complete</p>
      <p class="text-text-secondary">
        {{ result.usersConverted }} user{{ result.usersConverted !== 1 ? 's' : '' }} converted,
        {{ result.usersUnused }} unused voucher{{ result.usersUnused !== 1 ? 's' : '' }} cleared,
        {{ result.usersSkipped }} left untouched
        ({{ result.usersScanned }} scanned)
      </p>
      <p class="text-text-secondary">
        {{ result.profilesConverted }} profile{{ result.profilesConverted !== 1 ? 's' : '' }} converted
        ({{ result.profilesScanned }} scanned)
      </p>
      <p class="text-text-secondary">
        {{ result.schedulersRemoved }} leftover scheduler{{ result.schedulersRemoved !== 1 ? 's' : '' }} and
        {{ result.scriptsRemoved }} leftover script{{ result.scriptsRemoved !== 1 ? 's' : '' }} removed
      </p>
      <p class="text-text-secondary">
        {{ result.cleanupInstalled ? 'Pikro cleanup scheduler installed' : 'Could not install Pikro cleanup scheduler — set it up manually from the Dashboard' }}
      </p>
    </div>
    <p v-else-if="error" class="text-xs text-red">{{ error }}</p>

    <button
      type="button"
      class="btn btn-primary"
      :disabled="migrating"
      @click="runMigration"
    >
      <span
        v-if="migrating"
        class="size-4 border-2 border-black/20 border-t-black rounded-full animate-spin"
      />
      {{ migrating ? 'Migrating…' : 'Migrate from Mikhmon' }}
    </button>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRoutersStore } from '@/stores/routers'
import { migrateFromMikhmon, type MikhmonMigrationResult } from '@/api'
import StatusBadge from '@/components/StatusBadge.vue'

const store = useRoutersStore()

const migrating = ref(false)
const result = ref<MikhmonMigrationResult | null>(null)
const error = ref('')

async function runMigration() {
  if (!store.activeId) return
  if (!confirm(
    'This will rewrite hotspot user comments and profile on-login scripts, remove Mikhmon\'s ' +
    'own leftover schedulers/scripts, and install Pikro\'s cleanup scheduler on this router. ' +
    'Already-migrated data is left untouched. Continue?',
  )) return
  migrating.value = true
  error.value = ''
  result.value = null
  try {
    result.value = await migrateFromMikhmon(store.activeId)
  } catch (e: any) {
    error.value = e?.response?.data?.error ?? e?.message ?? 'Migration failed'
  } finally {
    migrating.value = false
  }
}
</script>
