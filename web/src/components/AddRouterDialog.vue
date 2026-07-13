<template>
  <AppDialog :open="open" title="Add Router" @update:open="$emit('update:open', $event)">
    <form @submit.prevent="submit" class="space-y-3">
      <div class="grid grid-cols-2 gap-3">
        <label class="col-span-2 flex flex-col gap-1">
          <span class="text-xs font-medium text-text-secondary">Name</span>
          <input v-model="form.name" class="input" placeholder="Home router" required />
        </label>
        <label class="flex flex-col gap-1">
          <span class="text-xs font-medium text-text-secondary">Host / IP</span>
          <input v-model="form.host" class="input" placeholder="192.168.88.1" required />
        </label>
        <label class="flex flex-col gap-1">
          <span class="text-xs font-medium text-text-secondary">Port</span>
          <input v-model.number="form.port" type="number" class="input" placeholder="8728" />
        </label>
        <label class="flex flex-col gap-1">
          <span class="text-xs font-medium text-text-secondary">Username</span>
          <input v-model="form.username" class="input" placeholder="admin" required />
        </label>
        <label class="flex flex-col gap-1">
          <span class="text-xs font-medium text-text-secondary">Password</span>
          <input v-model="form.password" type="password" class="input" />
        </label>
      </div>

      <label class="flex items-center gap-2 cursor-pointer select-none">
        <span
          class="relative inline-flex items-center justify-center size-4 rounded border shrink-0 transition-colors"
          :style="form.useTls
            ? 'background: var(--color-accent); border-color: var(--color-accent)'
            : 'background: transparent; border-color: var(--color-border)'"
        >
          <svg v-if="form.useTls" viewBox="0 0 10 8" fill="none" class="size-2.5">
            <path d="M1 4l2.5 2.5L9 1" stroke="#09090b" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" />
          </svg>
        </span>
        <input v-model="form.useTls" type="checkbox" class="sr-only" />
        <span class="text-sm text-text-secondary">Use TLS (port 8729)</span>
      </label>

 <p v-if="submitError" class="text-sm text-red" >{{ submitError }}</p>

      <div class="flex justify-end gap-2 pt-1">
        <button type="button" class="btn btn-ghost" @click="$emit('update:open', false)">Cancel</button>
        <button type="submit" class="btn btn-primary" :disabled="loading">
          <span v-if="loading" class="spinner spinner--sm" />
          Add router
        </button>
      </div>
    </form>
  </AppDialog>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRoutersStore } from '@/stores/routers'
import type { DiscoveredDevice } from '@/api'
import AppDialog from '@/components/AppDialog.vue'

const props = defineProps<{ open: boolean; prefill?: Partial<DiscoveredDevice> }>()
const emit = defineEmits<{ 'update:open': [value: boolean]; added: [ip: string] }>()

const store = useRoutersStore()
const loading = ref(false)
const submitError = ref('')

const form = ref({ name: '', host: '', port: 8728, username: 'admin', password: '', useTls: false })

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
    await store.add(form.value)
    emit('added', form.value.host)
    emit('update:open', false)
    form.value = { name: '', host: '', port: 8728, username: 'admin', password: '', useTls: false }
  } catch (e: any) {
    submitError.value = e?.response?.data?.error ?? 'Failed to add router'
  } finally {
    loading.value = false
  }
}
</script>
