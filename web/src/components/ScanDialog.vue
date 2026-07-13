<template>
  <AppDialog :open="open" title="Scan Network" @update:open="$emit('update:open', $event)">
    <div class="space-y-4">
      <!-- Scan trigger -->
      <div class="flex items-center justify-between">
        <p class="text-xs text-text-muted">
          Broadcasts an MNDP request to find MikroTik devices on this subnet.
        </p>
        <button class="btn btn-primary btn-sm shrink-0 ml-4" :disabled="scanning" @click="scan">
          <ArrowPathIcon class="size-3.5" :class="{ 'animate-spin': scanning }" />
          {{ scanning ? 'Scanning…' : 'Scan' }}
        </button>
      </div>

      <!-- Error -->
      <div v-if="error" class="flex items-center gap-2 p-3 rounded-lg border text-xs bg-red/8 border-red/20 text-red">
        <ExclamationTriangleIcon class="size-4 shrink-0" />
        {{ error }}
      </div>

      <!-- Not yet scanned -->
      <div v-else-if="!scanned" class="border border-dashed border-border rounded-xl py-10 text-center">
        <MagnifyingGlassIcon class="size-6 text-text-muted mx-auto mb-2" />
        <p class="text-xs text-text-muted">Click Scan to search for devices</p>
      </div>

      <!-- Scanning in progress -->
      <div v-else-if="scanning" class="border border-dashed border-border rounded-xl py-10 text-center">
        <span class="inline-block size-5 border-2 border-border border-t-text-secondary rounded-full animate-spin mb-2" />
        <p class="text-xs text-text-muted">Scanning… (3 s)</p>
      </div>

      <!-- No results -->
      <div v-else-if="devices.length === 0" class="border border-dashed border-border rounded-xl py-10 text-center">
        <ServerIcon class="size-6 text-text-muted mx-auto mb-2" />
        <p class="text-xs font-medium text-text-secondary">No devices found</p>
        <p class="text-xs text-text-muted mt-0.5">Make sure the router is on the same subnet</p>
      </div>

      <!-- Results -->
      <div v-else class="space-y-2">
        <p class="text-xs font-medium text-text-muted">{{ devices.length }} device{{ devices.length > 1 ? 's' : '' }} found</p>
        <div
          v-for="d in devices"
          :key="d.mac || d.ip"
          class="flex items-center gap-3 p-3 border border-border rounded-xl transition-colors bg-surface"
        >
 <div class="size-9 shrink-0 rounded-lg border border-border flex items-center justify-center bg-base" >
            <ServerIcon class="size-4 text-text-muted" />
          </div>
          <div class="flex-1 min-w-0">
            <p class="text-sm font-medium text-text-primary truncate">{{ d.identity || d.board || 'MikroTik' }}</p>
            <p class="text-xs text-text-secondary font-mono">{{ d.ip }} · {{ d.mac }}</p>
            <p v-if="d.board || d.version" class="text-xs text-text-muted mt-0.5">
              {{ [d.board, d.version].filter(Boolean).join(' · ') }}
            </p>
          </div>
          <button
            class="btn btn-sm shrink-0"
            :class="added.has(d.ip) ? 'btn-ghost cursor-default' : 'btn-ghost'"
            :style="added.has(d.ip) ? 'color: var(--color-green); border-color: rgba(76,195,138,0.3)' : ''"
            :disabled="added.has(d.ip)"
            @click="openAdd(d)"
          >
            <CheckIcon v-if="added.has(d.ip)" class="size-3.5" />
            <PlusIcon v-else class="size-3.5" />
            {{ added.has(d.ip) ? 'Added' : 'Add' }}
          </button>
        </div>
      </div>
    </div>
  </AppDialog>

  <AddRouterDialog v-model:open="showAdd" :prefill="prefill" @added="onAdded" />
</template>

<script setup lang="ts">
import { ref } from 'vue'
import {
  ArrowPathIcon, MagnifyingGlassIcon, ServerIcon,
  PlusIcon, CheckIcon, ExclamationTriangleIcon,
} from '@heroicons/vue/24/outline'
import { discoverRouters, type DiscoveredDevice } from '@/api'
import AppDialog from '@/components/AppDialog.vue'
import AddRouterDialog from '@/components/AddRouterDialog.vue'

defineProps<{ open: boolean }>()
defineEmits<{ 'update:open': [value: boolean] }>()

const scanning = ref(false)
const scanned = ref(false)
const error = ref('')
const devices = ref<DiscoveredDevice[]>([])
const added = ref(new Set<string>())

const showAdd = ref(false)
const prefill = ref<Partial<DiscoveredDevice>>({})

async function scan() {
  scanning.value = true
  scanned.value = true
  error.value = ''
  devices.value = []
  try {
    devices.value = await discoverRouters()
  } catch (e: any) {
    error.value = e?.response?.data?.error ?? 'Scan failed — check network/firewall'
    scanned.value = false
  } finally {
    scanning.value = false
  }
}

function openAdd(d: DiscoveredDevice) {
  prefill.value = d
  showAdd.value = true
}

function onAdded(ip: string) {
  added.value = new Set([...added.value, ip])
}
</script>
