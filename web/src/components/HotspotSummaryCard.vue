<template>
  <div
    v-if="loading && totalUsers === 0"
    class="flex items-center justify-center py-4"
  >
    <span class="spinner spinner--sm" />
  </div>
  <div
    v-else
    class="grid grid-cols-4 gap-3 tracking-tight text-sm font-semibold text-secondary"
  >
    <RouterLink
      to="/hotspot/users"
      class="rounded-lg border border-border p-3 space-y-1 bg-base transition-colors hover:border-text-muted"
    >
      <div class="text-text-secondary">Users</div>
      <div class="font-mono text-3xl font-bold">
        {{ totalUsers }}
      </div>
      <div class="text-right">
        <span class="btn btn-ghost btn-sm">View users</span>
      </div>
    </RouterLink>
    <RouterLink
      to="/hotspot/users?tab=active"
      class="rounded-lg border border-border p-3 space-y-1 bg-base transition-colors hover:border-text-muted"
    >
      <div class="text-text-secondary">Active</div>
      <div class="font-mono text-3xl font-bold">{{ activeSessions }}</div>
      <div class="text-right">
        <span class="btn btn-ghost btn-sm">View sessions</span>
      </div>
    </RouterLink>
    <div class="rounded-lg border border-border p-3 space-y-1 bg-base">
      <div class="text-text-secondary">Disabled</div>
      <div class="font-mono text-3xl font-bold text-text-secondary">
        {{ disabledUsers }}
      </div>
    </div>
    <div class="flex flex-col gap-2.5">
      <RouterLink
        to="/hotspot/vouchers"
        class="btn flex-1 justify-center border-transparent text-primary-fg bg-primary hover:bg-primary-hover"
      >
        <TicketIcon class="size-5" />
        Generate users
      </RouterLink>
      <RouterLink
        to="/hotspot/users"
        class="btn btn-primary flex-1 justify-center"
      >
        <UserPlusIcon class="size-5" />
        Add user
      </RouterLink>
    </div>
  </div>
</template>

<script setup lang="ts">
import { RouterLink } from "vue-router";
import { TicketIcon, UserPlusIcon } from "@heroicons/vue/24/outline";

defineProps<{
  loading: boolean;
  totalUsers: number;
  activeSessions: number;
  disabledUsers: number;
}>();
</script>
