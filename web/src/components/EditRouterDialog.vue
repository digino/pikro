<template>
  <AppDialog :open="open" title="Edit Router" @update:open="$emit('update:open', $event)">
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
          <input v-model="form.password" type="password" class="input" placeholder="Leave blank to keep current" />
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
          Save changes
        </button>
      </div>
    </form>
  </AppDialog>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import type { RouterProfile } from '@/api'
import AppDialog from '@/components/AppDialog.vue'

const props = defineProps<{ open: boolean; router: RouterProfile | null }>()
const emit = defineEmits<{ 'update:open': [value: boolean]; saved: [] }>()

const loading = ref(false)
const submitError = ref('')
const form = ref({ name: '', host: '', port: 8728, username: '', password: '', useTls: false })

watch(() => props.router, (r) => {
  if (!r) return
  form.value = { name: r.name, host: r.host, port: r.port, username: r.username, password: '', useTls: r.useTls }
}, { immediate: true })

async function submit() {
  if (!props.router) return
  loading.value = true
  submitError.value = ''
  try {
    const body: Record<string, unknown> = {
      name: form.value.name,
      host: form.value.host,
      port: form.value.port,
      username: form.value.username,
      useTls: form.value.useTls,
    }
    if (form.value.password) body.password = form.value.password

    const res = await fetch(`/api/routers/${props.router.id}`, {
      method: 'PATCH',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(body),
    })
    if (!res.ok) {
      const data = await res.json().catch(() => ({}))
      throw new Error(data.error ?? 'Failed to save')
    }
    emit('saved')
    emit('update:open', false)
  } catch (e: any) {
    submitError.value = e?.message ?? 'Failed to save'
  } finally {
    loading.value = false
  }
}
</script>
