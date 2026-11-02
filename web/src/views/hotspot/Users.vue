<template>
  <main class="flex flex-col gap-6 px-4 sm:px-6 lg:px-10 pt-7 pb-14">
    <header class="flex flex-wrap items-end justify-between gap-4">
      <div class="flex items-center gap-4">
        <img :src="inkIcon('users')" alt="" width="50" height="50" />
        <div>
          <nav :aria-label="t('users.breadcrumb')" class="eyebrow">
            {{ t('users.breadcrumb') }} / {{ t('users.title') }}
          </nav>
          <h1 class="type-display flex items-center gap-3 mt-1 text-3xl">
            {{ t('users.title') }}
            <span class="chip px-2 text-sm bg-primary text-primary-fg">{{ users.length }}</span>
          </h1>
        </div>
      </div>
      <div class="flex flex-wrap gap-3">
        <button
          v-if="tab === 'users'"
          type="button"
          class="btn"
          :disabled="filteredUsers.length === 0"
          @click="openPrint()"
        >
          <PrinterIcon class="size-4" />
          {{ t('users.print') }}
        </button>
        <RouterLink to="/hotspot/vouchers" class="btn btn-primary">
          <img :src="inkIcon('vouchers')" alt="" width="28" height="28" class="-ml-1" />
          {{ t('users.generateVouchers') }}
        </RouterLink>
        <button type="button" class="btn btn-ink" @click="openAdd">
          <PlusIcon class="size-4" />
          {{ t('users.newUser') }}
        </button>
      </div>
    </header>

    <!-- All users / Online now -->
    <div
      role="group"
      :aria-label="t('users.viewLabel')"
      class="self-start flex border-2 border-ink rounded-control overflow-hidden"
    >
      <button
        v-for="(tb, i) in tabs"
        :key="tb.key"
        type="button"
        :aria-pressed="tab === tb.key"
        class="min-h-11 px-4 text-sm font-semibold cursor-pointer transition-colors"
        :class="[
          i > 0 && 'border-l-2 border-ink',
          tab === tb.key ? 'bg-ink text-paper' : 'bg-paper text-ink hover:bg-primary-wash',
        ]"
        @click="switchTab(tb.key)"
      >
        {{ t(tb.label) }}
        <span class="ml-1 font-mono text-xs opacity-75">
          {{ tb.key === 'users' ? users.length : active.length }}
        </span>
      </button>
    </div>

    <div v-if="loading" class="flex justify-center py-10">
      <span class="spinner" />
    </div>

    <div
      v-else-if="error"
      role="alert"
      class="flex items-center gap-2.5 p-4 rounded-control border-2 border-danger-text bg-danger-wash text-danger-text text-sm"
    >
      <ExclamationTriangleIcon class="size-5 shrink-0" />
      {{ error }}
      <button
        type="button"
        class="btn btn-sm ml-auto"
        @click="tab === 'users' ? loadUsers() : loadActive()"
      >
        {{ t('users.retry') }}
      </button>
    </div>

    <!-- Users tab -->
    <template v-else-if="tab === 'users'">
      <div class="flex flex-wrap items-center justify-between gap-3">
        <div role="group" :aria-label="t('users.filterLabel')" class="flex flex-wrap gap-2">
          <button
            v-for="f in statusFilters"
            :key="f.value"
            type="button"
            :aria-pressed="filterStatus === f.value"
            class="inline-flex items-center gap-2 min-h-11 px-4 rounded-full border-2 border-ink text-sm font-semibold cursor-pointer transition-colors"
            :class="
              filterStatus === f.value
                ? 'bg-ink text-paper shadow-[3px_3px_0_var(--color-primary)]'
                : 'bg-paper text-ink hover:bg-primary-wash'
            "
            @click="filterStatus = f.value"
          >
            {{ f.label }}
            <span class="font-mono text-xs opacity-75">{{ f.count }}</span>
          </button>
        </div>

        <div class="flex flex-wrap gap-2.5">
          <label
            class="flex items-center gap-2 min-h-11 min-w-60 px-3.5 rounded-control border-2 border-ink bg-paper focus-within:outline-3 focus-within:outline-primary focus-within:outline-offset-1"
          >
            <MagnifyingGlassIcon class="size-4.5 shrink-0" />
            <span class="sr-only">{{ t('users.search') }}</span>
            <input
              v-model="searchQuery"
              type="search"
              :placeholder="t('users.searchPlaceholder')"
              class="flex-1 min-w-0 border-0 bg-transparent text-sm outline-none"
            />
          </label>
          <SelectRoot
            :model-value="filterProfile || ALL_PLANS"
            @update:model-value="filterProfile = $event === ALL_PLANS ? '' : String($event ?? '')"
          >
            <SelectTrigger
              :aria-label="t('users.plan')"
              class="flex items-center gap-2 min-h-11 px-3 rounded-control border-2 border-ink bg-paper cursor-pointer focus-visible:outline-3 focus-visible:outline-primary focus-visible:outline-offset-1"
            >
              <span class="eyebrow">{{ t('users.plan') }}</span>
              <SelectValue class="text-sm font-semibold" />
              <ChevronDownIcon class="size-4 shrink-0" />
            </SelectTrigger>
            <SelectPortal>
              <SelectContent
                class="z-50 min-w-(--reka-select-trigger-width) bg-paper text-ink border-2 border-ink rounded-control shadow-btn overflow-hidden"
                position="popper"
                :side-offset="4"
              >
                <SelectViewport class="p-1">
                  <SelectItem
                    v-for="opt in planOptions"
                    :key="opt.value"
                    :value="opt.value"
                    class="flex items-center justify-between gap-3 min-h-10 px-3 text-sm rounded-lg cursor-pointer outline-none data-highlighted:bg-primary-wash data-[state=checked]:font-semibold"
                  >
                    <SelectItemText>{{ opt.label }}</SelectItemText>
                    <SelectItemIndicator><CheckIcon class="size-4" /></SelectItemIndicator>
                  </SelectItem>
                </SelectViewport>
              </SelectContent>
            </SelectPortal>
          </SelectRoot>
        </div>
      </div>

      <section class="card overflow-hidden">
        <div class="overflow-x-auto">
          <table class="w-full min-w-215 text-sm">
            <thead>
              <tr class="border-b-2 border-ink">
                <th class="w-5 px-4 py-3.5">
                  <CheckboxRoot
                    :model-value="someSelected ? 'indeterminate' : allSelected"
                    :aria-label="t('users.selectAll')"
                    class="size-5 rounded-md border-2 border-ink bg-paper flex items-center justify-center cursor-pointer data-[state=checked]:bg-ink data-[state=indeterminate]:bg-ink"
                    @update:model-value="toggleSelectAll"
                  >
                    <CheckboxIndicator>
                      <MinusIcon v-if="someSelected" class="size-3.5 text-paper" />
                      <CheckIcon v-else class="size-3.5 text-paper" />
                    </CheckboxIndicator>
                  </CheckboxRoot>
                </th>
                <th class="eyebrow text-left px-4 py-3.5">{{ t('users.colUsername') }}</th>
                <th class="eyebrow text-left px-4 py-3.5">{{ t('users.colPlan') }}</th>
                <th class="eyebrow text-left px-4 py-3.5">{{ t('users.colTimeLeft') }}</th>
                <th class="eyebrow text-left px-4 py-3.5">{{ t('users.colData') }}</th>
                <th class="eyebrow text-left px-4 py-3.5">{{ t('users.colExpires') }}</th>
                <th class="eyebrow text-left px-4 py-3.5">{{ t('users.colStatus') }}</th>
                <th class="px-4 py-3.5"><span class="sr-only">{{ t('users.colActions') }}</span></th>
              </tr>
            </thead>
            <tbody>
              <tr v-if="filteredUsers.length === 0">
                <td colspan="8" class="px-4 py-10 text-center text-text-muted">
                  {{ searchQuery || filterProfile || filterStatus ? t('users.noMatch') : t('users.noUsers') }}
                  <button
                    v-if="searchQuery || filterProfile || filterStatus"
                    type="button"
                    class="btn btn-sm ml-3"
                    @click="clearFilters"
                  >
                    {{ t('users.clearFilters') }}
                  </button>
                </td>
              </tr>
              <tr
                v-for="u in pagedUsers"
                :key="u['.id']"
                class="border-b border-hairline last:border-0"
                :class="{ 'bg-primary-wash': selected.has(u['.id']) }"
              >
                <td class="px-4 py-3">
                  <CheckboxRoot
                    :model-value="selected.has(u['.id'])"
                    :aria-label="t('users.selectUser', { user: u.name })"
                    class="size-5 rounded-md border-2 border-ink bg-paper flex items-center justify-center cursor-pointer data-[state=checked]:bg-ink"
                    @update:model-value="toggleSelect(u['.id'])"
                  >
                    <CheckboxIndicator>
                      <CheckIcon class="size-3.5 text-paper" />
                    </CheckboxIndicator>
                  </CheckboxRoot>
                </td>
                <td class="px-4 py-3">
                  <div class="font-mono font-semibold text-[15px]">{{ u.name }}</div>
                  <div v-if="displayComment(u.comment) !== '—'" class="mt-0.5 text-xs text-text-muted">
                    {{ displayComment(u.comment) }}
                  </div>
                </td>
                <td class="px-4 py-3">
                  <span class="chip">{{ u.profile || 'default' }}</span>
                </td>
                <td class="px-4 py-3">
                  <div v-if="timeLeft(u)" class="flex items-center gap-2.5">
                    <div class="w-18 h-2 shrink-0 rounded-full border-[1.5px] border-ink bg-paper overflow-hidden">
                      <div
                        class="h-full"
                        :class="timeLeft(u)!.pct < 40 ? 'bg-primary' : 'bg-ink'"
                        :style="{ width: `${timeLeft(u)!.pct}%` }"
                      />
                    </div>
                    <span class="text-[13px]">{{ timeLeft(u)!.label }}</span>
                  </div>
                  <span v-else class="text-text-muted">—</span>
                </td>
                <td class="px-4 py-3 font-mono text-[13px]">{{ dataUsage(u) }}</td>
                <td class="px-4 py-3 text-[13px] text-text-secondary">{{ expiresAt(u) }}</td>
                <td class="px-4 py-3">
                  <span class="chip" :class="STATUS_CHIP[userStatus(u)].chip">
                    <span class="size-2 shrink-0 rounded-full" :class="STATUS_CHIP[userStatus(u)].dot" />
                    {{ t(`users.status.${userStatus(u)}`) }}
                  </span>
                </td>
                <td class="px-4 py-3 text-right">
                  <DropdownMenuRoot>
                    <DropdownMenuTrigger
                      :aria-label="t('users.actionsFor', { user: u.name })"
                      class="inline-flex items-center justify-center size-10 rounded-[10px] border-2 border-ink bg-paper cursor-pointer hover:bg-primary-wash data-[state=open]:bg-primary-wash focus-visible:outline-3 focus-visible:outline-primary"
                    >
                      <EllipsisHorizontalIcon class="size-5" />
                    </DropdownMenuTrigger>
                    <DropdownMenuPortal>
                      <DropdownMenuContent
                        align="end"
                        :side-offset="6"
                        class="z-50 w-52 p-1.5 bg-paper text-ink border-2 border-ink rounded-[14px] shadow-[4px_4px_0_var(--color-shadow)]"
                      >
                        <DropdownMenuItem :class="MENU_ITEM" @select="openEdit(u)">
                          <PencilIcon class="size-4" />{{ t('users.edit') }}
                        </DropdownMenuItem>
                        <DropdownMenuItem :class="MENU_ITEM" @select="openPrint([u])">
                          <PrinterIcon class="size-4" />{{ t('users.printVoucher') }}
                        </DropdownMenuItem>
                        <DropdownMenuItem :class="MENU_ITEM" @select="toggleDisabled(u)">
                          <component :is="u.disabled === 'true' ? CheckCircleIcon : NoSymbolIcon" class="size-4" />
                          {{ u.disabled === 'true' ? t('users.enableUser') : t('users.disableUser') }}
                        </DropdownMenuItem>
                        <DropdownMenuSeparator class="h-px my-1 mx-1.5 bg-hairline" />
                        <DropdownMenuItem
                          :class="[MENU_ITEM, 'text-red font-semibold']"
                          @select="removeUser(u['.id'])"
                        >
                          <TrashIcon class="size-4" />{{ t('users.delete') }}
                        </DropdownMenuItem>
                      </DropdownMenuContent>
                    </DropdownMenuPortal>
                  </DropdownMenuRoot>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Pagination footer -->
        <div
          v-if="filteredUsers.length > 0"
          class="flex flex-wrap items-center justify-between gap-3 px-5 py-3.5 border-t-2 border-ink text-[13px]"
        >
          <span>
            {{ t('users.showing', {
              from: (usersPage - 1) * usersPageSize + 1,
              to: Math.min(usersPage * usersPageSize, filteredUsers.length),
              total: filteredUsers.length,
            }) }}
          </span>
          <div class="flex flex-wrap items-center gap-4">
            <div class="flex items-center gap-1.5">
              <span class="eyebrow mr-1">{{ t('users.perPage') }}</span>
              <button
                v-for="n in PAGE_SIZES"
                :key="n"
                type="button"
                class="btn btn-sm px-2.5"
                :class="{ 'btn-ink': usersPageSize === n }"
                :aria-pressed="usersPageSize === n"
                @click="setUsersPageSize(n)"
              >
                {{ n }}
              </button>
            </div>
            <PaginationRoot
              :page="usersPage"
              :items-per-page="usersPageSize"
              :total="filteredUsers.length"
              class="flex gap-2"
              @update:page="usersPage = $event"
            >
              <PaginationPrev class="btn btn-sm" :aria-label="t('users.previous')">
                <ChevronLeftIcon class="size-4" />
              </PaginationPrev>
              <PaginationNext class="btn btn-sm" :aria-label="t('users.next')">
                <ChevronRightIcon class="size-4" />
              </PaginationNext>
            </PaginationRoot>
          </div>
        </div>
      </section>

      <!-- Floating bulk bar -->
      <div
        v-if="selected.size > 0"
        role="region"
        :aria-label="t('users.bulkActions')"
        class="sticky bottom-5 z-10 self-center flex flex-wrap items-center gap-2 py-2.5 pr-2.5 pl-4.5 rounded-2xl bg-sidebar text-sidebar-text shadow-[5px_5px_0_var(--color-primary)]"
      >
        <span class="mr-2 font-bold">
          <span class="font-mono px-1.75 py-0.5 rounded-md bg-primary text-primary-fg">{{ selected.size }}</span>
          {{ t('users.selected') }}
        </span>
        <button type="button" :class="[BULK_BTN, 'bg-paper text-ink']" :disabled="bulkBusy" @click="setSelectedDisabled(false)">
          {{ t('users.enable') }}
        </button>
        <button type="button" :class="[BULK_BTN, 'bg-paper text-ink']" :disabled="bulkBusy" @click="setSelectedDisabled(true)">
          {{ t('users.disable') }}
        </button>
        <button type="button" :class="[BULK_BTN, 'bg-paper text-ink']" @click="openPrint(selectedUsers)">
          <PrinterIcon class="size-4" />{{ t('users.printSelected') }}
        </button>
        <button
          type="button"
          :class="[BULK_BTN, 'bg-danger text-white']"
          :disabled="bulkBusy"
          @click="removeSelected"
        >
          <span v-if="bulkBusy" class="spinner spinner--sm" />
          <TrashIcon v-else class="size-4" />
          {{ t('users.delete') }}
        </button>
        <button
          type="button"
          :aria-label="t('users.clearSelection')"
          class="inline-flex items-center justify-center size-10 rounded-[10px] cursor-pointer hover:bg-white/10"
          @click="selected = new Set()"
        >
          <XMarkIcon class="size-5" />
        </button>
      </div>
    </template>

    <!-- Online now tab -->
    <section v-else class="card overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full min-w-180 text-sm">
          <thead>
            <tr class="border-b-2 border-ink">
              <th class="eyebrow text-left px-4 py-3.5">{{ t('users.colUser') }}</th>
              <th class="eyebrow text-left px-4 py-3.5">{{ t('users.colIp') }}</th>
              <th class="eyebrow text-left px-4 py-3.5">{{ t('users.colMac') }}</th>
              <th class="eyebrow text-right px-4 py-3.5">{{ t('users.colUptime') }}</th>
              <th class="eyebrow text-right px-4 py-3.5">{{ t('users.colDown') }}</th>
              <th class="eyebrow text-right px-4 py-3.5">{{ t('users.colUp') }}</th>
              <th class="eyebrow text-right px-4 py-3.5">{{ t('users.colLeft') }}</th>
              <th class="px-4 py-3.5"><span class="sr-only">{{ t('users.colActions') }}</span></th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="active.length === 0">
              <td colspan="8" class="px-4 py-10 text-center text-text-muted">{{ t('users.noSessions') }}</td>
            </tr>
            <tr
              v-for="s in pagedActive"
              :key="s['.id']"
              class="border-b border-hairline last:border-0"
            >
              <td class="px-4 py-3 font-mono font-semibold">{{ s.user }}</td>
              <td class="px-4 py-3 font-mono">{{ s.address }}</td>
              <td class="px-4 py-3 font-mono text-text-muted">{{ s['mac-address'] || '—' }}</td>
              <td class="px-4 py-3 text-right font-mono">{{ s.uptime || '—' }}</td>
              <td class="px-4 py-3 text-right font-mono">{{ formatBytes(s['bytes-in']) }}</td>
              <td class="px-4 py-3 text-right font-mono text-text-secondary">{{ formatBytes(s['bytes-out']) }}</td>
              <td class="px-4 py-3 text-right font-mono text-text-secondary">{{ s['session-time-left'] || '—' }}</td>
              <td class="px-4 py-3">
                <div class="flex items-center justify-end gap-2">
                  <template v-if="userIdByName(s.user)">
                    <button type="button" class="btn btn-sm" @click="disableSessionUser(s.user)">
                      <NoSymbolIcon class="size-4" />{{ t('users.disable') }}
                    </button>
                    <button type="button" class="btn btn-sm btn-danger" @click="removeSessionUser(s.user)">
                      <TrashIcon class="size-4" />{{ t('users.delete') }}
                    </button>
                  </template>
                  <button v-else type="button" class="btn btn-sm btn-danger" @click="disconnectSession(s['.id'])">
                    {{ t('users.disconnect') }}
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Online now pagination footer -->
      <div
        v-if="active.length > 0"
        class="flex flex-wrap items-center justify-between gap-3 px-5 py-3.5 border-t-2 border-ink text-[13px]"
      >
        <span>
          {{ t('users.showing', {
            from: (activePage - 1) * activePageSize + 1,
            to: Math.min(activePage * activePageSize, active.length),
            total: active.length,
          }) }}
        </span>
        <div class="flex flex-wrap items-center gap-4">
          <div class="flex items-center gap-1.5">
            <span class="eyebrow mr-1">{{ t('users.perPage') }}</span>
            <button
              v-for="n in PAGE_SIZES"
              :key="n"
              type="button"
              class="btn btn-sm px-2.5"
              :class="{ 'btn-ink': activePageSize === n }"
              :aria-pressed="activePageSize === n"
              @click="setActivePageSize(n)"
            >
              {{ n }}
            </button>
          </div>
          <PaginationRoot
            :page="activePage"
            :items-per-page="activePageSize"
            :total="active.length"
            class="flex gap-2"
            @update:page="activePage = $event"
          >
            <PaginationPrev class="btn btn-sm" :aria-label="t('users.previous')">
              <ChevronLeftIcon class="size-4" />
            </PaginationPrev>
            <PaginationNext class="btn btn-sm" :aria-label="t('users.next')">
              <ChevronRightIcon class="size-4" />
            </PaginationNext>
          </PaginationRoot>
        </div>
      </div>
    </section>

    <!-- Add single user dialog -->
    <AppDialog
      :open="showAdd"
      title="New Hotspot User"
      @update:open="showAdd = $event"
    >
      <form @submit.prevent="submitAdd" class="space-y-4">
        <div class="grid grid-cols-2 gap-3">
          <label class="flex flex-col gap-1">
            <span class="font-medium text-red">Username <span>*</span></span>
            <input v-model="form.name" class="input" required />
          </label>
          <label class="flex flex-col gap-1">
            <span class="font-medium text-text-secondary">Password</span>
            <input
              v-model="form.password"
              class="input"
              placeholder="leave blank for none"
            />
          </label>
        </div>

        <label class="flex flex-col gap-1">
          <span class="font-medium text-text-secondary">Profile</span>
          <select v-model="form.profile" class="input">
            <option value="">default</option>
            <option v-for="p in profiles" :key="p['.id']" :value="p.name">
              {{ p.name }}
            </option>
          </select>
        </label>

        <div class="border-t border-border pt-3 space-y-3">
          <p
            class="text-sm font-semibold text-text-muted uppercase tracking-wide"
          >
            Limits (override profile)
          </p>

          <div class="flex flex-col gap-1">
            <span class="font-medium text-text-secondary">Time limit</span>
            <input
              v-model="form.limitUptimeRaw"
              class="input font-mono"
              placeholder="e.g. 1h, 1d, 1w, 1d12h — blank for unlimited"
              @blur="normalizeFormUptime"
            />
            <p
              v-if="form.limitUptimeRaw && !formUptimePreview"
              class="text-sm text-red"
            >
              Invalid format — use: 30m, 2h, 1d, 1w or combinations like 1d12h
            </p>
            <p
              v-else-if="formUptimePreview"
              class="text-sm"
              :class="formUptimeWarning ? 'text-amber' : 'text-text-muted'"
            >
              <span v-if="formUptimeWarning"
                >⚠ Time limit exceeds the profile's validity — user may never
                hit this limit.</span
              >
              <span v-else class="font-mono"
                >Sends to router: <span>{{ formUptimePreview }}</span></span
              >
            </p>
          </div>

          <div class="flex flex-col gap-1">
            <span class="font-medium text-text-secondary">Data limit</span>
            <div
              class="flex w-full overflow-hidden rounded-lg border border-border focus-within:outline-2 focus-within:outline-accent focus-within:outline-offset-1"
            >
              <input
                v-model.number="form.limitBytesTotalValue"
                type="number"
                min="0"
                class="input-bare flex-1 min-w-0"
                placeholder="0"
              />
              <select
                v-model="form.limitBytesTotalUnit"
                class="input-bare border-l border-border shrink-0 w-16 text-xs"
              >
                <option value="M">MB</option>
                <option value="G">GB</option>
              </select>
            </div>
          </div>
        </div>

        <label class="flex flex-col gap-1">
          <span class="font-medium text-text-secondary">Comment</span>
          <input
            v-model="form.comment"
            class="input"
            placeholder="optional note"
          />
        </label>

        <p v-if="addError" class="text-sm text-red">{{ addError }}</p>

        <div class="flex justify-end gap-2 pt-1">
          <button type="button" class="btn btn-ghost" @click="showAdd = false">
            Cancel
          </button>
          <button
            type="submit"
            class="btn btn-ink"
            :disabled="adding || (!!form.limitUptimeRaw && !formUptimePreview)"
          >
            <span
              v-if="adding"
              class="size-4 border-2 border-black/20 border-t-black rounded-full animate-spin"
            />
            Create
          </button>
        </div>
      </form>
    </AppDialog>

    <!-- Edit user dialog -->
    <AppDialog
      :open="showEdit"
      :title="`Edit ${editForm.name}`"
      @update:open="showEdit = $event"
    >
      <form @submit.prevent="submitEdit" class="space-y-4">
        <label class="flex flex-col gap-1">
          <span class="font-medium text-text-secondary">Password</span>
          <input
            v-model="editForm.password"
            class="input"
            placeholder="leave blank to keep current"
          />
        </label>

        <label class="flex flex-col gap-1">
          <span class="font-medium text-text-secondary">Profile</span>
          <select v-model="editForm.profile" class="input">
            <option value="">default</option>
            <option v-for="p in profiles" :key="p['.id']" :value="p.name">
              {{ p.name }}
            </option>
          </select>
        </label>

        <div class="border-t border-border pt-3 space-y-3">
          <p
            class="text-sm font-semibold text-text-muted uppercase tracking-wide"
          >
            Limits (override profile)
          </p>

          <div class="flex flex-col gap-1">
            <span class="font-medium text-text-secondary">Time limit</span>
            <input
              v-model="editForm.limitUptimeRaw"
              class="input font-mono"
              placeholder="e.g. 1h, 1d, 1w, 1d12h — blank for unlimited"
              @blur="normalizeEditUptime"
            />
            <p
              v-if="editForm.limitUptimeRaw && !editUptimePreview"
              class="text-sm text-red"
            >
              Invalid format — use: 30m, 2h, 1d, 1w or combinations like 1d12h
            </p>
            <p
              v-else-if="editUptimePreview"
              class="text-sm"
              :class="editUptimeWarning ? 'text-amber' : 'text-text-muted'"
            >
              <span v-if="editUptimeWarning"
                >⚠ Time limit exceeds the profile's validity — user may never
                hit this limit.</span
              >
              <span v-else class="font-mono"
                >Sends to router: <span>{{ editUptimePreview }}</span></span
              >
            </p>
          </div>

          <div class="flex flex-col gap-1">
            <span class="font-medium text-text-secondary">Data limit</span>
            <div
              class="flex w-full overflow-hidden rounded-lg border border-border focus-within:outline-2 focus-within:outline-accent focus-within:outline-offset-1"
            >
              <input
                v-model.number="editForm.limitBytesTotalValue"
                type="number"
                min="0"
                class="input-bare flex-1 min-w-0"
                placeholder="0"
              />
              <select
                v-model="editForm.limitBytesTotalUnit"
                class="input-bare border-l border-border shrink-0 w-16 text-xs"
              >
                <option value="M">MB</option>
                <option value="G">GB</option>
              </select>
            </div>
          </div>
        </div>

        <p v-if="editError" class="text-sm text-red">{{ editError }}</p>

        <div class="flex justify-end gap-2 pt-1">
          <button type="button" class="btn btn-ghost" @click="showEdit = false">
            Cancel
          </button>
          <button
            type="submit"
            class="btn btn-ink"
            :disabled="editing || (!!editForm.limitUptimeRaw && !editUptimePreview)"
          >
            <span
              v-if="editing"
              class="size-4 border-2 border-black/20 border-t-black rounded-full animate-spin"
            />
            Save changes
          </button>
        </div>
      </form>
    </AppDialog>

    <PrintTemplateDialog
      :open="showPrintDialog"
      :entries="printEntries"
      :business-name="hotspotSettings.hotspotName ?? ''"
      :currency="hotspotSettings.currency ?? ''"
      :profile-metas="profileMetas"
      :login-url="loginUrl"
      :login-url-supports-credentials="true"
      @update:open="showPrintDialog = $event"
    />
  </main>
</template>

<script setup lang="ts">
import { ref, computed, watch, onUnmounted } from "vue";
import {
  SelectRoot,
  SelectTrigger,
  SelectValue,
  SelectPortal,
  SelectContent,
  SelectViewport,
  SelectItem,
  SelectItemText,
  SelectItemIndicator,
  CheckboxRoot,
  CheckboxIndicator,
  PaginationRoot,
  PaginationPrev,
  PaginationNext,
  DropdownMenuRoot,
  DropdownMenuTrigger,
  DropdownMenuPortal,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
} from "reka-ui";
import {
  PlusIcon,
  TrashIcon,
  ExclamationTriangleIcon,
  PrinterIcon,
  CheckIcon,
  MinusIcon,
  NoSymbolIcon,
  ChevronLeftIcon,
  ChevronRightIcon,
  ChevronDownIcon,
  CheckCircleIcon,
  PencilIcon,
  MagnifyingGlassIcon,
  EllipsisHorizontalIcon,
  XMarkIcon,
} from "@heroicons/vue/24/outline";
import { RouterLink, useRoute } from "vue-router";
import { useI18n } from "vue-i18n";
import { useRoutersStore } from "@/stores/routers";
import {
  listHotspotUsers,
  listHotspotActive,
  listHotspotProfiles,
  createHotspotUser,
  updateHotspotUser,
  toggleHotspotUser,
  deleteHotspotUser,
  disconnectHotspotActive,
  getProfileMetas,
  getHotspotSettings,
  type ProfileMeta,
  type HotspotSettings,
} from "@/api";
import { friendlyError } from "@/utils/errors";
import AppDialog from "@/components/AppDialog.vue";
import PrintTemplateDialog from "@/components/PrintTemplateDialog.vue";
import { hotspotLoginUrl, type VoucherEntry } from "@/utils/vouchers";
import { inkIcon } from "@/utils/inkIcons";

const store = useRoutersStore();
const route = useRoute();
const { t, locale } = useI18n();

type UserStatus = ReturnType<typeof userStatus>;

const STATUS_ORDER: UserStatus[] = [
  "active",
  "expires-soon",
  "waiting",
  "limit-reached",
  "expired",
  "disabled",
];

// Statuses differ by fill and dot shape, not colour alone (see the design system).
const STATUS_CHIP: Record<UserStatus, { chip: string; dot: string }> = {
  active: { chip: "bg-primary text-primary-fg", dot: "bg-primary-fg" },
  "expires-soon": { chip: "bg-primary text-primary-fg", dot: "bg-primary-fg" },
  waiting: { chip: "", dot: "border-2 border-ink" },
  "limit-reached": { chip: "bg-ink text-paper", dot: "bg-paper" },
  expired: { chip: "bg-ink text-paper", dot: "bg-paper" },
  disabled: { chip: "bg-hairline text-text-secondary border-dashed", dot: "bg-graphite" },
};

const MENU_ITEM =
  "flex items-center gap-2.5 w-full min-h-10 px-3 rounded-lg text-sm font-medium cursor-pointer outline-none data-highlighted:bg-ground";
const BULK_BTN =
  "inline-flex items-center gap-2 min-h-10 px-3.5 rounded-control border-2 border-sidebar-text text-sm font-semibold cursor-pointer disabled:opacity-40 disabled:cursor-not-allowed";

const ALL_PLANS = "__all__";

const tab = ref<"users" | "active">(
  route.query.tab === "active" ? "active" : "users",
);
const tabs = [
  { key: "users" as const, label: "users.tabAll" },
  { key: "active" as const, label: "users.tabActive" },
];
const users = ref<Record<string, string>[]>([]);
const active = ref<Record<string, string>[]>([]);

// Active sessions are keyed by username, not the user record's .id — this
// maps session.user -> the matching /ip/hotspot/user .id so the Active tab
// can offer a Disable action on the underlying account.
const userIdByNameMap = computed(() => {
  const map = new Map<string, string>();
  for (const u of users.value) map.set(u.name, u[".id"]);
  return map;
});
function userIdByName(name: string) {
  return userIdByNameMap.value.get(name);
}

const searchQuery = ref("");
const filterProfile = ref("");
const filterStatus = ref("");

// Search + plan filter; the status pills count within this set.
const searchedUsers = computed(() => {
  let list = users.value;
  const q = searchQuery.value.trim().toLowerCase();
  if (q)
    list = list.filter(
      (u) =>
        u.name?.toLowerCase().includes(q) ||
        displayComment(u.comment).toLowerCase().includes(q),
    );
  if (filterProfile.value)
    list = list.filter((u) => (u.profile || "default") === filterProfile.value);
  return list;
});

const filteredUsers = computed(() =>
  filterStatus.value
    ? searchedUsers.value.filter((u) => userStatus(u) === filterStatus.value)
    : searchedUsers.value,
);

// "All" plus every status that has users (or is the current filter).
const statusFilters = computed(() => {
  const counts = new Map<string, number>();
  for (const u of searchedUsers.value) {
    const s = userStatus(u);
    counts.set(s, (counts.get(s) ?? 0) + 1);
  }
  return [
    { value: "", label: t("users.all"), count: searchedUsers.value.length },
    ...STATUS_ORDER.filter((s) => counts.has(s) || filterStatus.value === s).map((s) => ({
      value: s,
      label: t(`users.status.${s}`),
      count: counts.get(s) ?? 0,
    })),
  ];
});

const planOptions = computed(() => [
  { value: ALL_PLANS, label: t("users.allPlans") },
  ...profiles.value.map((p) => ({ value: p.name, label: p.name })),
]);

function clearFilters() {
  searchQuery.value = "";
  filterProfile.value = "";
  filterStatus.value = "";
}

const PAGE_SIZES = [20, 50, 100] as const;
const usersPageSize = ref<20 | 50 | 100>(20);
const usersPage = ref(1);
const activePageSize = ref<20 | 50 | 100>(100);
const activePage = ref(1);

const pagedUsers = computed(() => {
  const start = (usersPage.value - 1) * usersPageSize.value;
  return filteredUsers.value.slice(start, start + usersPageSize.value);
});
const pagedActive = computed(() => {
  const start = (activePage.value - 1) * activePageSize.value;
  return active.value.slice(start, start + activePageSize.value);
});

function setUsersPageSize(n: 20 | 50 | 100) {
  usersPageSize.value = n;
  usersPage.value = 1;
}
function setActivePageSize(n: 20 | 50 | 100) {
  activePageSize.value = n;
  activePage.value = 1;
}

watch(users, () => {
  usersPage.value = 1;
});
watch(active, () => {
  activePage.value = 1;
});
watch([searchQuery, filterProfile, filterStatus], () => {
  usersPage.value = 1;
});

const profiles = ref<Record<string, string>[]>([]);
const profileMetas = ref<Record<string, ProfileMeta>>({});
const hotspotSettings = ref<HotspotSettings>({
  hotspotName: "",
  currency: "",
});

// Used by the Business voucher template's QR code — points at the
// hotspot's own login page, since that DNS name is what devices on the
// hotspot network actually resolve. normalizeDnsName guards against a
// dnsName saved with a scheme already in it (e.g. from before this was
// validated on input) producing a doubled "http://http://...".
const loginUrl = ref("");

const loading = ref(false);
const error = ref("");

const selected = ref<Set<string>>(new Set());
const bulkBusy = ref(false);

const selectedUsers = computed(() =>
  users.value.filter((u) => selected.value.has(u[".id"])),
);

const allSelected = computed(
  () =>
    filteredUsers.value.length > 0 &&
    filteredUsers.value.every((u) => selected.value.has(u[".id"])),
);
const someSelected = computed(
  () =>
    !allSelected.value &&
    filteredUsers.value.some((u) => selected.value.has(u[".id"])),
);

function toggleSelect(id: string) {
  const s = new Set(selected.value);
  if (s.has(id)) s.delete(id);
  else s.add(id);
  selected.value = s;
}

function toggleSelectAll() {
  selected.value = allSelected.value
    ? new Set()
    : new Set(filteredUsers.value.map((u) => u[".id"]));
}

async function removeSelected() {
  if (!store.activeId || !selected.value.size) return;
  if (!confirm(t("users.confirmDeleteMany", selected.value.size))) return;
  bulkBusy.value = true;
  const ids = [...selected.value];
  await Promise.allSettled(
    ids.map((id) => deleteHotspotUser(store.activeId!, id)),
  );
  selected.value = new Set();
  bulkBusy.value = false;
  await loadUsers();
}

// Keeps the selection so the operator can follow up (e.g. print after enabling).
async function setSelectedDisabled(disabled: boolean) {
  if (!store.activeId || !selected.value.size) return;
  bulkBusy.value = true;
  await Promise.allSettled(
    [...selected.value].map((id) => toggleHotspotUser(store.activeId!, id, disabled)),
  );
  bulkBusy.value = false;
  await loadUsers();
}

const showAdd = ref(false);
const adding = ref(false);
const addError = ref("");

const emptyForm = () => ({
  name: "",
  password: "",
  profile: "",
  limitUptimeRaw: "",
  limitBytesTotalValue: 0,
  limitBytesTotalUnit: "M" as "M" | "G",
  comment: "",
});
const form = ref(emptyForm());

const showEdit = ref(false);
const editing = ref(false);
const editError = ref("");
const editUserId = ref("");

const emptyEditForm = () => ({
  name: "",
  password: "",
  profile: "",
  limitUptimeRaw: "",
  limitBytesTotalValue: 0,
  limitBytesTotalUnit: "M" as "M" | "G",
});
const editForm = ref(emptyEditForm());

function bytesToValueUnit(val: string | undefined): { value: number; unit: "M" | "G" } {
  const n = parseInt(val ?? "0");
  if (!n) return { value: 0, unit: "M" };
  if (n % 1024 ** 3 === 0) return { value: n / 1024 ** 3, unit: "G" };
  return { value: Math.round(n / 1024 ** 2), unit: "M" };
}

function openEdit(u: Record<string, string>) {
  editUserId.value = u[".id"];
  const { value, unit } = bytesToValueUnit(u["limit-bytes-total"]);
  editForm.value = {
    name: u.name,
    password: "",
    profile: u.profile ?? "",
    limitUptimeRaw: normalizeShorthand(u["limit-uptime"] ?? ""),
    limitBytesTotalValue: value,
    limitBytesTotalUnit: unit,
  };
  editError.value = "";
  showEdit.value = true;
}

const editUptimePreview = computed(() =>
  uptimePreviewFrom(editForm.value.limitUptimeRaw),
);
const editUptimeWarning = computed(() => {
  if (!editUptimePreview.value || !editForm.value.profile) return false;
  const validitySecs = selectedProfileValiditySeconds(editForm.value.profile);
  if (!validitySecs) return false;
  return shorthandToSeconds(editForm.value.limitUptimeRaw) > validitySecs;
});
function normalizeEditUptime() {
  editForm.value.limitUptimeRaw = normalizeShorthand(editForm.value.limitUptimeRaw);
}

async function submitEdit() {
  if (!store.activeId || !editUserId.value) return;
  editing.value = true;
  editError.value = "";
  try {
    const f = editForm.value;
    const limitUptime = editUptimePreview.value;
    const mul = f.limitBytesTotalUnit === "G" ? 1024 ** 3 : 1024 ** 2;
    const limitBytesTotal = f.limitBytesTotalValue
      ? String(f.limitBytesTotalValue * mul)
      : "";
    await updateHotspotUser(store.activeId, editUserId.value, {
      ...(f.password ? { password: f.password } : {}),
      profile: f.profile,
      limitUptime,
      limitBytesTotal,
    });
    showEdit.value = false;
    await loadUsers();
  } catch (e: any) {
    editError.value = friendlyError(e, "Failed to update user");
  } finally {
    editing.value = false;
  }
}

function parseShorthand(
  s: string,
): { w: number; d: number; h: number; m: number } | null {
  if (!s.trim()) return null;
  const re = /^(?:(\d+)w)?(?:(\d+)d)?(?:(\d+)h)?(?:(\d+)m)?$/i;
  const match = s.trim().match(re);
  if (!match || !match[0]) return null;
  const [, w, d, h, min] = match;
  if (!w && !d && !h && !min) return null;
  return {
    w: parseInt(w || "0"),
    d: parseInt(d || "0"),
    h: parseInt(h || "0"),
    m: parseInt(min || "0"),
  };
}

function shorthandToSeconds(s: string): number {
  const p = parseShorthand(s);
  if (!p) return 0;
  return p.w * 7 * 86400 + p.d * 86400 + p.h * 3600 + p.m * 60;
}

function uptimePreviewFrom(raw: string): string {
  const p = parseShorthand(raw);
  if (!p) return raw ? "" : "";
  const totalDays = p.w * 7 + p.d;
  if (!totalDays && !p.h && !p.m) return "";
  return [
    totalDays ? `${totalDays}d` : "",
    p.h ? `${p.h}h` : "",
    p.m ? `${p.m}m` : "",
  ]
    .filter(Boolean)
    .join("");
}

function normalizeShorthand(raw: string): string {
  const p = parseShorthand(raw);
  if (!p) return raw;
  return [
    p.w ? `${p.w}w` : "",
    p.d ? `${p.d}d` : "",
    p.h ? `${p.h}h` : "",
    p.m ? `${p.m}m` : "",
  ]
    .filter(Boolean)
    .join("");
}

function selectedProfileValiditySeconds(profileName: string): number {
  if (!profileName) return 0;
  const meta = profileMetas.value[profileName];
  if (!meta?.validity) return 0;
  return shorthandToSeconds(meta.validity);
}

const formUptimePreview = computed(() =>
  uptimePreviewFrom(form.value.limitUptimeRaw),
);
const formUptimeWarning = computed(() => {
  if (!formUptimePreview.value || !form.value.profile) return false;
  const validitySecs = selectedProfileValiditySeconds(form.value.profile);
  if (!validitySecs) return false;
  return shorthandToSeconds(form.value.limitUptimeRaw) > validitySecs;
});
function normalizeFormUptime() {
  form.value.limitUptimeRaw = normalizeShorthand(form.value.limitUptimeRaw);
}

let activeSessionsTimer: ReturnType<typeof setInterval> | undefined;

function switchTab(key: "users" | "active") {
  tab.value = key;
  clearInterval(activeSessionsTimer);
  if (key === "active") {
    loadActive();
    activeSessionsTimer = setInterval(() => {
      if (!document.hidden) loadActive();
    }, 5 * 60_000);
  }
}

onUnmounted(() => clearInterval(activeSessionsTimer));

async function loadUsers() {
  if (!store.activeId) return;
  loading.value = true;
  error.value = "";
  try {
    const [u, p, m, s] = await Promise.all([
      listHotspotUsers(store.activeId),
      listHotspotProfiles(store.activeId),
      getProfileMetas(store.activeId).catch(
        () => ({}) as Record<string, ProfileMeta>,
      ),
      getHotspotSettings(store.activeId).catch(
        () =>
          ({ hotspotName: "", currency: "" }) as HotspotSettings,
      ),
    ]);
    users.value = u;
    profiles.value = p;
    profileMetas.value = m;
    hotspotSettings.value = s;
    loginUrl.value = await hotspotLoginUrl(store.activeId);
  } catch (e: any) {
    error.value = friendlyError(e, "Failed to load users");
  } finally {
    loading.value = false;
  }
}

async function loadActive() {
  if (!store.activeId) return;
  loading.value = true;
  error.value = "";
  try {
    active.value = await listHotspotActive(store.activeId);
  } catch (e: any) {
    error.value = friendlyError(e, "Failed to load sessions");
  } finally {
    loading.value = false;
  }
}

function openAdd() {
  form.value = emptyForm();
  addError.value = "";
  showAdd.value = true;
}

async function removeUser(userId: string) {
  if (!store.activeId || !confirm(t("users.confirmDelete"))) return;
  await deleteHotspotUser(store.activeId, userId);
  await loadUsers();
}

// Deletes the underlying user record for an active session — used instead of
// a plain session kick because RouterOS's mac-cookie auth (part of the
// default login-by chain) silently re-authenticates the same device by MAC
// as soon as it's disconnected, so a kick alone doesn't actually remove
// access. Deleting the user (which cascades to kick any active session, see
// DeleteHotspotUser) is the only way to truly revoke it.
async function removeSessionUser(username: string) {
  const userId = userIdByName(username);
  if (
    !store.activeId ||
    !userId ||
    !confirm("Delete this user? This also disconnects their active session.")
  )
    return;
  await deleteHotspotUser(store.activeId, userId);
  await Promise.all([loadUsers(), loadActive()]);
}

async function disconnectSession(sessionId: string) {
  if (
    !store.activeId ||
    !confirm("Disconnect this session? The user can log back in.")
  )
    return;
  await disconnectHotspotActive(store.activeId, sessionId);
  await loadActive();
}

async function submitAdd() {
  if (!store.activeId) return;
  adding.value = true;
  addError.value = "";
  try {
    const f = form.value;
    const limitUptime = formUptimePreview.value;
    const mul = f.limitBytesTotalUnit === "G" ? 1024 ** 3 : 1024 ** 2;
    const limitBytesTotal = f.limitBytesTotalValue
      ? String(f.limitBytesTotalValue * mul)
      : "";
    await createHotspotUser(store.activeId, {
      name: f.name,
      password: f.password,
      profile: f.profile,
      limitUptime,
      limitBytesTotal,
      rateLimit: "",
      comment: f.comment,
      expiryComment: "",
    });
    showAdd.value = false;
    form.value = emptyForm();
    await loadUsers();
  } catch (e: any) {
    addError.value = friendlyError(e, "Failed to create user");
  } finally {
    adding.value = false;
  }
}

async function toggleDisabled(u: Record<string, string>) {
  if (!store.activeId) return;
  const nextDisabled = u.disabled !== "true";
  try {
    await toggleHotspotUser(store.activeId, u[".id"], nextDisabled);
    await loadUsers();
  } catch (e: any) {
    error.value = friendlyError(e, "Failed to update user");
  }
}

// Disables the account behind an active session — this does not disconnect
// the session itself (RouterOS keeps it live until it naturally expires);
// use disconnectSession for that.
async function disableSessionUser(username: string) {
  if (!store.activeId) return;
  const userId = userIdByName(username);
  if (!userId) return;
  try {
    await toggleHotspotUser(store.activeId, userId, true);
    await loadUsers();
  } catch (e: any) {
    error.value = friendlyError(e, "Failed to update user");
  }
}

const showPrintDialog = ref(false);
const printEntries = ref<VoucherEntry[]>([]);

// With no selection, prints every currently filtered/visible user instead —
// so the header Print button always does something useful.
function openPrint(source: Record<string, string>[] = filteredUsers.value) {
  printEntries.value = source.map((u) => ({
    name: u.name,
    password: u.password ?? "",
    profile: u.profile,
    timeLimit: u["limit-uptime"] ?? "",
  }));
  showPrintDialog.value = true;
}

function isUptimeExhausted(u: Record<string, string>): boolean {
  const limit = u["limit-uptime"],
    used = u["uptime"];
  if (!limit || !used) return false;
  return parseUptimeSeconds(used) >= parseUptimeSeconds(limit);
}

function extractExpEpoch(comment: string | undefined): number | null {
  if (!comment) return null;
  const m = comment.match(/^exp:(\d+)/);
  if (!m) return null;
  return parseInt(m[1]);
}

function userStatus(
  u: Record<string, string>,
):
  | "disabled"
  | "expired"
  | "expires-soon"
  | "limit-reached"
  | "waiting"
  | "active" {
  if (u.disabled === "true") return "disabled";
  // Uptime quota is enforced by RouterOS in real time — once hit, the user is
  // already cut off regardless of what the (possibly stale) exp: comment says.
  // Check it before calendar expiry so the badge reflects the reason that
  // actually applies, not whichever the code happens to check first.
  if (isUptimeExhausted(u)) return "limit-reached";
  const epoch = extractExpEpoch(u.comment);
  const now = Math.floor(Date.now() / 1000);
  if (epoch !== null && epoch < now) return "expired";
  if (epoch !== null && epoch - now < 86400) return "expires-soon";
  if (epoch === null) return "waiting";
  return "active";
}

// Remaining time and how much of the allowance it is: from the uptime quota
// when the user has one, otherwise from the voucher's calendar expiry.
function timeLeft(u: Record<string, string>): { label: string; pct: number } | null {
  let left: number, total: number;
  const limit = parseUptimeSeconds(u["limit-uptime"]);
  const epoch = extractExpEpoch(u.comment);
  const validity = selectedProfileValiditySeconds(u.profile);
  if (limit) {
    left = limit - parseUptimeSeconds(u.uptime);
    total = limit;
  } else if (epoch !== null && validity) {
    left = epoch - Math.floor(Date.now() / 1000);
    total = validity;
  } else {
    return null;
  }
  left = Math.max(0, left);
  return {
    label: left > 0 ? formatSeconds(left) : "0m",
    pct: Math.min(100, (left / total) * 100),
  };
}

function dataUsage(u: Record<string, string>): string {
  const used = (parseInt(u["bytes-in"] ?? "0") || 0) + (parseInt(u["bytes-out"] ?? "0") || 0);
  const usedLabel = used ? formatBytes(String(used)) : "0";
  return u["limit-bytes-total"] && u["limit-bytes-total"] !== "0"
    ? `${usedLabel} / ${formatBytes(u["limit-bytes-total"])}`
    : usedLabel;
}

function expiresAt(u: Record<string, string>): string {
  const epoch = extractExpEpoch(u.comment);
  if (epoch === null) return "—";
  return new Date(epoch * 1000).toLocaleString(locale.value, {
    day: "numeric",
    month: "short",
    hour: "2-digit",
    minute: "2-digit",
  });
}

// RouterOS durations come as "1w2d3h4m5s" or, depending on version and field,
// clock style ("00:22:10", "1d 03:04:05").
function parseUptimeSeconds(uptime: string | undefined): number {
  if (!uptime) return 0;
  const UNIT: Record<string, number> = { w: 604800, d: 86400, h: 3600, m: 60, s: 1 };
  let secs = 0;
  const clock = uptime.match(/(\d+):(\d+):(\d+)/);
  if (clock) secs += +clock[1] * 3600 + +clock[2] * 60 + +clock[3];
  for (const [, n, unit] of uptime.matchAll(/(\d+)([wdhms])/g)) secs += +n * UNIT[unit];
  return secs;
}

function formatSeconds(secs: number): string {
  if (secs <= 0) return "expired";
  const d = Math.floor(secs / 86400),
    h = Math.floor((secs % 86400) / 3600),
    m = Math.floor((secs % 3600) / 60);
  if (d > 0) return `${d}d ${h}h`;
  if (h > 0) return `${h}h ${m}m`;
  return `${m}m`;
}

function displayComment(comment: string | undefined): string {
  if (!comment) return "—";
  return comment.replace(/^exp:\d+\s*/, "") || "—";
}

function formatBytes(val: string | undefined): string {
  const n = parseInt(val ?? "0");
  if (!n) return "—";
  if (n >= 1024 ** 3) return (n / 1024 ** 3).toFixed(1) + " GB";
  if (n >= 1024 ** 2) return (n / 1024 ** 2).toFixed(1) + " MB";
  if (n >= 1024) return (n / 1024).toFixed(1) + " KB";
  return n + " B";
}

watch(() => store.activeId, loadUsers, { immediate: true });

// Deep-linked from Dashboard's "Active sessions" card via ?tab=active.
if (tab.value === "active") switchTab("active");
</script>
