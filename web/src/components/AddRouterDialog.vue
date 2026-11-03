<template>
  <AppDialog :open="open" :title="t('routers.addDialog.title')" @update:open="emit('update:open', $event)">
    <form class="flex flex-col gap-5" @submit.prevent="submit">
      <fieldset class="flex flex-col gap-3">
        <legend class="eyebrow mb-3">{{ t('routers.addDialog.connection') }}</legend>
        <div class="grid grid-cols-2 gap-3">
          <label class="col-span-2 flex flex-col gap-1">
            <span class="font-medium">{{ t('routers.addDialog.name') }}</span>
            <input v-model="form.name" class="input" :placeholder="t('routers.addDialog.namePlaceholder')" required />
          </label>
          <label class="flex flex-col gap-1">
            <span class="font-medium">{{ t('routers.addDialog.host') }}</span>
            <input v-model="form.host" class="input" placeholder="192.168.88.1" required />
          </label>
          <label class="flex flex-col gap-1">
            <span class="font-medium">{{ t('routers.addDialog.port') }}</span>
            <input v-model.number="form.port" type="number" class="input" placeholder="8728" />
          </label>
          <label class="flex flex-col gap-1">
            <span class="font-medium">{{ t('routers.addDialog.username') }}</span>
            <input v-model="form.username" class="input" placeholder="admin" autocomplete="off" required />
          </label>
          <label class="flex flex-col gap-1">
            <span class="font-medium">{{ t('routers.addDialog.password') }}</span>
            <input v-model="form.password" type="password" class="input" autocomplete="new-password" />
          </label>
        </div>

        <label class="flex items-center gap-2.5 cursor-pointer select-none">
          <CheckboxRoot
            v-model="form.useTls"
            class="size-5 rounded-md border-2 border-ink bg-paper flex items-center justify-center cursor-pointer data-[state=checked]:bg-ink"
          >
            <CheckboxIndicator>
              <CheckIcon class="size-3.5 text-paper" />
            </CheckboxIndicator>
          </CheckboxRoot>
          <span class="text-sm">{{ t('routers.addDialog.useTls') }}</span>
        </label>
      </fieldset>

      <fieldset class="flex flex-col pt-5 border-t-2 border-hairline">
        <legend class="eyebrow float-left mb-3">{{ t('routers.addDialog.hotspot') }}</legend>
        <RouterHotspotInfoFields
          v-model:hotspot-name="form.hotspotName"
          v-model:currency="form.currency"
        />
      </fieldset>

      <p v-if="submitError" role="alert" class="m-0 text-sm text-danger-text">{{ submitError }}</p>

      <div class="flex justify-end gap-2.5">
        <button type="button" class="btn" @click="emit('update:open', false)">
          {{ t('routers.addDialog.cancel') }}
        </button>
        <button
          type="submit"
          class="btn btn-primary"
          :disabled="loading || !form.host || !form.username || !form.name"
        >
          <span v-if="loading" class="spinner spinner--sm" />
          {{ t('routers.addDialog.add') }}
        </button>
      </div>
    </form>
  </AppDialog>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { CheckboxRoot, CheckboxIndicator } from 'reka-ui'
import { CheckIcon } from '@heroicons/vue/24/outline'
import { useRoutersStore } from '@/stores/routers'
import type { DiscoveredDevice } from '@/api'
import AppDialog from '@/components/AppDialog.vue'
import RouterHotspotInfoFields from '@/components/RouterHotspotInfoFields.vue'

const props = defineProps<{ open: boolean; prefill?: Partial<DiscoveredDevice> }>()
const emit = defineEmits<{ 'update:open': [value: boolean]; added: [ip: string] }>()

const { t } = useI18n()
const store = useRoutersStore()
const loading = ref(false)
const submitError = ref('')

function emptyForm() {
  return {
    name: '', host: '', port: 8728, username: 'admin', password: '', useTls: false,
    hotspotName: '', currency: '',
  }
}
const form = ref(emptyForm())

// Pre-fill form when opened from the scan dialog
watch(() => props.prefill, (d) => {
  if (!d) return
  form.value.name = d.identity || d.board || 'MikroTik'
  form.value.host = d.ip ?? ''
}, { immediate: true })

async function submit() {
  loading.value = true
  submitError.value = ''
  try {
    const f = form.value
    await store.add({
      name: f.name, host: f.host, port: f.port, username: f.username,
      password: f.password, useTls: f.useTls,
      // Only trust the scanned MAC if the operator kept the scanned address.
      mac: props.prefill?.ip === f.host ? props.prefill.mac : undefined,
      hotspotSettings: {
        hotspotName: f.hotspotName,
        currency: f.currency,
      },
    })
    emit('added', f.host)
    emit('update:open', false)
    form.value = emptyForm()
  } catch (e: any) {
    submitError.value = e?.response?.data?.error ?? t('routers.addDialog.failed')
  } finally {
    loading.value = false
  }
}
</script>
