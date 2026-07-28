<template>
  <PageLayout title="Hotspot" subtitle="Users">
    <template #actions>
      <button class="btn btn-ghost" @click="openBatch">
        <TicketIcon class="size-3.5" />
        Generate users
      </button>
      <button class="btn btn-primary" @click="openAdd">
        <PlusIcon class="size-3.5" />
        New user
      </button>
    </template>

    <!-- Tabs -->
    <div class="flex items-center gap-1 border-b border-border -mt-2">
      <button
        v-for="t in tabs"
        :key="t.key"
        class="px-3 py-2 text-sm font-medium border-b-2 transition-colors"
        :class="
          tab === t.key
            ? 'border-text-primary text-text-primary'
            : 'border-transparent text-text-muted hover:text-text-secondary'
        "
        @click="switchTab(t.key)"
      >
        {{ t.label }}
        <span
          v-if="t.key === 'users' || active.length > 0"
          class="ml-1 text-text-secondary"
        >
          ({{ t.key === "users" ? filteredUsers.length : active.length }})
        </span>
      </button>
    </div>

    <div v-if="loading" class="flex justify-center py-10">
      <span class="spinner" />
    </div>

    <div
      v-else-if="error"
      class="flex items-center gap-2 p-4 border rounded-xl text-sm bg-red/8 border-red/20 text-red"
    >
      <ExclamationTriangleIcon class="size-4 shrink-0" />
      {{ error }}
      <button
        class="ml-auto text-xs underline"
        @click="tab === 'users' ? loadUsers() : loadActive()"
      >
        Retry
      </button>
    </div>

    <!-- Users tab -->
    <div v-else-if="tab === 'users'" class="space-y-2">
      <!-- Bulk action bar -->
      <div
        v-if="selected.size > 0"
        class="flex items-center gap-3 px-4 py-2.5 border border-border rounded-xl text-sm font-medium bg-surface"
      >
        <span class="text-text-secondary">{{ selected.size }} selected</span>
        <button class="btn btn-ghost btn-sm ml-auto" @click="printSelected">
          <PrinterIcon class="size-3.5" />
          Print
        </button>
        <button
          class="btn btn-danger btn-sm"
          :disabled="bulkDeleting"
          @click="removeSelected"
        >
          <span
            v-if="bulkDeleting"
            class="size-3.5 border-2 border-white/30 border-t-white rounded-full animate-spin"
          />
          <TrashIcon v-else class="size-3.5" />
          Delete {{ selected.size }}
        </button>
        <button
          class="text-text-muted hover:text-text-secondary transition-colors"
          @click="selected = new Set()"
        >
          ✕
        </button>
      </div>

      <!-- Search & filters -->
      <div class="flex items-center gap-2.5 flex-wrap">
        <div class="relative flex-1 min-w-40 max-w-xs">
          <input
            v-model="searchQuery"
            class="input pl-9"
            placeholder="Search username or comment…"
          />
        </div>
        <select v-model="filterProfile" class="input-select">
          <option value="">All profiles</option>
          <option v-for="p in profiles" :key="p['.id']" :value="p.name">
            {{ p.name }}
          </option>
        </select>
        <select v-model="filterStatus" class="input-select">
          <option value="">All statuses</option>
          <option value="waiting">Waiting</option>
          <option value="limit-reached">Limit reached</option>
          <option value="expired">Expired</option>
          <option value="disabled">Disabled</option>
        </select>
        <button
          v-if="searchQuery || filterProfile || filterStatus"
          class="text-sm text-text-muted hover:text-text-secondary transition-colors px-1"
          @click="
            searchQuery = '';
            filterProfile = '';
            filterStatus = '';
          "
        >
          Clear
        </button>
      </div>

      <div class="border border-border rounded-xl overflow-hidden">
        <table class="w-full text-sm">
          <thead>
            <tr class="border-b border-border bg-surface">
              <th class="px-4 py-3 w-8">
                <button
                  class="size-4 rounded border flex items-center justify-center transition-colors"
                  :class="
                    allSelected
                      ? 'bg-text-primary border-text-primary text-base'
                      : 'border-border hover:border-muted'
                  "
                  @click="toggleSelectAll"
                >
                  <CheckIcon v-if="allSelected" class="size-3" />
                </button>
              </th>
              <th
                class="text-left px-4 py-3 text-xs font-semibold text-text-primary uppercase tracking-wide"
              >
                Username
              </th>
              <th
                class="text-left px-4 py-3 text-xs font-semibold text-text-primary uppercase tracking-wide"
              >
                Profile
              </th>
              <th
                class="text-left px-4 py-3 text-xs font-semibold text-text-primary uppercase tracking-wide"
              >
                Time limit
              </th>
              <th
                class="text-left px-4 py-3 text-xs font-semibold text-text-primary uppercase tracking-wide"
              >
                Data limit
              </th>
              <th
                class="text-left px-4 py-3 text-xs font-semibold text-text-primary uppercase tracking-wide"
              >
                Comment
              </th>
              <th
                class="text-left px-4 py-3 text-xs font-semibold text-text-primary uppercase tracking-wide"
              >
                Expires
              </th>
              <th
                class="text-left px-4 py-3 text-xs font-semibold text-text-primary uppercase tracking-wide"
              >
                Status
              </th>
              <th class="px-4 py-3"></th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="filteredUsers.length === 0">
              <td colspan="9" class="text-center text-text-muted text-sm py-10">
                {{
                  searchQuery || filterProfile
                    ? "No users match your search"
                    : "No users found"
                }}
              </td>
            </tr>
            <tr
              v-for="u in pagedUsers"
              :key="u['.id']"
              class="border-b border-border last:border-0 transition-colors"
              :class="selected.has(u['.id']) ? 'bg-surface' : ''"
            >
              <td class="px-4 py-3">
                <button
                  class="size-4 rounded border flex items-center justify-center transition-colors"
                  :class="
                    selected.has(u['.id'])
                      ? 'bg-text-primary border-text-primary text-base'
                      : 'border-border hover:border-muted'
                  "
                  @click="toggleSelect(u['.id'])"
                >
                  <CheckIcon v-if="selected.has(u['.id'])" class="size-3" />
                </button>
              </td>
              <td
                class="px-4 py-3 font-mono font-bold text-text-primary text-sm"
              >
                {{ u.name }}
              </td>
              <td class="px-4 py-3">
                <span
                  class="text-xs text-text-secondary border border-border px-2 py-0.5 rounded bg-base"
                >
                  {{ u.profile || "default" }}
                </span>
              </td>
              <td class="px-4 py-3 text-sm text-text-secondary">
                {{ u["limit-uptime"] || "—" }}
              </td>
              <td class="px-4 py-3 text-sm text-text-secondary">
                {{ formatBytes(u["limit-bytes-total"]) }}
              </td>
              <td class="px-4 py-3 text-sm text-text-muted">
                {{ displayComment(u.comment) }}
              </td>
              <td class="px-4 py-3 text-sm" :class="expiryClass(u)">
                {{ expiryLabel(u) }}
              </td>
              <td class="px-4 py-3">
                <span
                  v-if="userStatus(u) === 'active'"
                  class="text-xs px-2 py-0.5 rounded-full bg-green/10 text-green"
                  >Active</span
                >
                <span
                  v-else-if="userStatus(u) === 'waiting'"
                  class="text-xs px-2 py-0.5 rounded-full bg-blue-500/10 text-blue-400"
                  >Waiting</span
                >
                <span
                  v-else-if="userStatus(u) === 'disabled'"
                  class="text-xs px-2 py-0.5 rounded-full bg-red/10 text-red"
                  >Disabled</span
                >
                <span
                  v-else-if="userStatus(u) === 'expired'"
                  class="text-xs px-2 py-0.5 rounded-full bg-amber/10 text-amber"
                  >Expired</span
                >
                <span
                  v-else-if="userStatus(u) === 'limit-reached'"
                  class="text-xs px-2 py-0.5 rounded-full bg-amber/10 text-amber"
                  >Limit reached</span
                >
              </td>
              <td class="px-4 py-3 text-right">
                <div class="flex items-center justify-end gap-1">
                  <button
                    class="btn btn-ghost btn-sm"
                    @click="openEdit(u)"
                  >
                    <PencilIcon class="size-3.5" />
                    Edit
                  </button>
                  <button
                    class="btn btn-sm btn-ghost hover:text-red hover:bg-red/10"
                    @click="removeUser(u['.id'])"
                  >
                    <TrashIcon class="size-3.5" />
                    Delete
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
        <!-- Pagination footer -->
        <div
          v-if="filteredUsers.length > 0"
          class="flex items-center justify-between px-4 py-2.5 border-t border-border bg-surface"
        >
          <div class="flex items-center gap-1">
            <span class="text-xs text-text-muted mr-1">Per page</span>
            <button
              v-for="n in PAGE_SIZES"
              :key="n"
              class="px-2 py-0.5 text-xs rounded transition-colors"
              :class="
                usersPageSize === n
                  ? 'bg-muted text-text-primary'
                  : 'text-text-muted hover:text-text-secondary'
              "
              @click="setUsersPageSize(n)"
            >
              {{ n }}
            </button>
          </div>
          <div class="flex items-center gap-1 text-xs text-text-muted">
            <span
              >{{ (usersPage - 1) * usersPageSize + 1 }}–{{
                Math.min(usersPage * usersPageSize, filteredUsers.length)
              }}
              of {{ filteredUsers.length }}</span
            >
            <button
              class="p-1 rounded hover:bg-surface disabled:opacity-30 transition-colors"
              :disabled="usersPage === 1"
              @click="usersPage--"
            >
              <ChevronLeftIcon class="size-3.5" />
            </button>
            <button
              class="p-1 rounded hover:bg-surface disabled:opacity-30 transition-colors"
              :disabled="usersPage >= usersPageCount"
              @click="usersPage++"
            >
              <ChevronRightIcon class="size-3.5" />
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Active sessions tab -->
    <div v-else class="border border-border rounded-xl overflow-hidden">
      <table class="w-full text-sm">
        <thead>
          <tr class="border-b border-border bg-surface">
            <th class="text-left px-4 py-3 text-xs font-semibold text-text-primary uppercase tracking-wide">User</th>
            <th class="text-left px-4 py-3 text-xs font-semibold text-text-primary uppercase tracking-wide">IP</th>
            <th class="text-left px-4 py-3 text-xs font-semibold text-text-primary uppercase tracking-wide">MAC</th>
            <th class="text-right px-4 py-3 text-xs font-semibold text-text-primary uppercase tracking-wide">Uptime</th>
            <th class="text-right px-4 py-3 text-xs font-semibold text-text-primary uppercase tracking-wide">Down</th>
            <th class="text-right px-4 py-3 text-xs font-semibold text-text-primary uppercase tracking-wide">Up</th>
            <th class="text-right px-4 py-3 text-xs font-semibold text-text-primary uppercase tracking-wide">Left</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="active.length === 0">
            <td colspan="7" class="text-center text-text-muted text-sm py-10">No active sessions</td>
          </tr>
          <tr
            v-for="s in pagedActive"
            :key="s['.id']"
            class="border-b border-border last:border-0 transition-colors hover:bg-surface"
          >
            <td class="px-4 py-3 font-mono font-semibold text-text-primary text-xs">{{ s.user }}</td>
            <td class="px-4 py-3 font-mono text-xs text-text-secondary">{{ s.address }}</td>
            <td class="px-4 py-3 font-mono text-xs text-text-muted">{{ s["mac-address"] || "—" }}</td>
            <td class="px-4 py-3 text-right font-mono text-xs text-text-secondary">{{ s.uptime || "—" }}</td>
            <td class="px-4 py-3 text-right font-mono text-xs text-text-primary">{{ formatBytes(s["bytes-in"]) }}</td>
            <td class="px-4 py-3 text-right font-mono text-xs text-text-secondary">{{ formatBytes(s["bytes-out"]) }}</td>
            <td class="px-4 py-3 text-right font-mono text-xs text-text-secondary">{{ s["session-time-left"] || "—" }}</td>
          </tr>
        </tbody>
      </table>
      <!-- Active sessions pagination footer -->
      <div
        v-if="active.length > 0"
        class="flex items-center justify-between px-4 py-2.5 border-t border-border bg-surface"
      >
        <div class="flex items-center gap-1">
          <span class="text-xs text-text-muted mr-1">Per page</span>
          <button
            v-for="n in PAGE_SIZES"
            :key="n"
            class="px-2 py-0.5 text-xs rounded transition-colors"
            :class="
              activePageSize === n
                ? 'bg-muted text-text-primary'
                : 'text-text-muted hover:text-text-secondary'
            "
            @click="setActivePageSize(n)"
          >
            {{ n }}
          </button>
        </div>
        <div class="flex items-center gap-1 text-xs text-text-muted">
          <span
            >{{ (activePage - 1) * activePageSize + 1 }}–{{
              Math.min(activePage * activePageSize, active.length)
            }}
            of {{ active.length }}</span
          >
          <button
            class="p-1 rounded hover:bg-surface disabled:opacity-30 transition-colors"
            :disabled="activePage === 1"
            @click="activePage--"
          >
            <ChevronLeftIcon class="size-3.5" />
          </button>
          <button
            class="p-1 rounded hover:bg-surface disabled:opacity-30 transition-colors"
            :disabled="activePage >= activePageCount"
            @click="activePage++"
          >
            <ChevronRightIcon class="size-3.5" />
          </button>
        </div>
      </div>
    </div>

    <!-- Edit user dialog -->
    <AppDialog
      :open="showEdit"
      :title="editForm.name ? `Edit — ${editForm.name}` : 'Edit user'"
      @update:open="showEdit = $event"
    >
      <form @submit.prevent="submitEdit" class="space-y-4">
        <label class="flex flex-col gap-1">
          <span class="text-sm font-medium text-text-secondary">Password</span>
          <input
            v-model="editForm.password"
            class="input"
            placeholder="leave blank to keep current"
          />
        </label>

        <label class="flex flex-col gap-1">
          <span class="text-xs font-medium text-text-secondary">Profile</span>
          <select v-model="editForm.profile" class="input">
            <option value="">default</option>
            <option v-for="p in profiles" :key="p['.id']" :value="p.name">
              {{ p.name }}
            </option>
          </select>
        </label>

        <div class="border-t border-border pt-3 space-y-3">
          <p
            class="text-xs font-semibold text-text-muted uppercase tracking-wide"
          >
            Limits
          </p>

          <div class="flex flex-col gap-1">
            <span class="text-xs font-medium text-text-secondary"
              >Time limit</span
            >
            <input
              v-model="editForm.limitUptimeRaw"
              class="input font-mono"
              placeholder="e.g. 1h30m — blank for unlimited"
              @blur="normalizeEditUptime"
            />
            <p
              v-if="editForm.limitUptimeRaw && !editUptimePreview"
              class="text-xs text-red"
            >
              Invalid format — use: 30m, 2h, 1d, 1w or combinations like 1d12h
            </p>
            <p v-else-if="editUptimePreview" class="text-xs text-text-muted">
              Sends to router:
              <span class="font-mono">{{ editUptimePreview }}</span>
            </p>
          </div>

          <div class="flex flex-col gap-1">
            <span class="text-xs font-medium text-text-secondary"
              >Data limit</span
            >
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

        <label class="flex flex-col gap-1">
          <span class="text-xs font-medium text-text-secondary">Comment</span>
          <input
            v-model="editForm.comment"
            class="input"
            placeholder="optional note"
          />
        </label>

        <div
          class="flex items-center justify-between border-t border-border pt-3"
        >
          <div>
            <p class="text-xs font-medium text-text-secondary">Disable user</p>
            <p class="text-xs text-text-muted">Blocked from logging in</p>
          </div>
          <button
            type="button"
            class="relative inline-flex h-5 w-9 shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors focus:outline-none"
            :class="editForm.disabled ? 'bg-red' : 'bg-border'"
            @click="editForm.disabled = !editForm.disabled"
          >
            <span
              class="pointer-events-none inline-block size-4 rounded-full bg-white shadow transform transition-transform"
              :class="editForm.disabled ? 'translate-x-4' : 'translate-x-0'"
            />
          </button>
        </div>

        <p v-if="editError" class="text-xs text-red">{{ editError }}</p>

        <div class="flex justify-end gap-2 pt-1">
          <button type="button" class="btn btn-ghost" @click="showEdit = false">
            Cancel
          </button>
          <button
            type="submit"
            class="btn btn-primary"
            :disabled="
              editing || (!!editForm.limitUptimeRaw && !editUptimePreview)
            "
          >
            <span
              v-if="editing"
              class="size-4 border-2 border-black/20 border-t-black rounded-full animate-spin"
            />
            Save
          </button>
        </div>
      </form>
    </AppDialog>

    <!-- Add single user dialog -->
    <AppDialog
      :open="showAdd"
      title="New Hotspot User"
      @update:open="showAdd = $event"
    >
      <form @submit.prevent="submitAdd" class="space-y-4">
        <div class="grid grid-cols-2 gap-3">
          <label class="flex flex-col gap-1">
            <span class="text-xs font-medium text-red"
              >Username <span>*</span></span
            >
            <input v-model="form.name" class="input" required />
          </label>
          <label class="flex flex-col gap-1">
            <span class="text-xs font-medium text-text-secondary"
              >Password</span
            >
            <input
              v-model="form.password"
              class="input"
              placeholder="leave blank for none"
            />
          </label>
        </div>

        <label class="flex flex-col gap-1">
          <span class="text-xs font-medium text-text-secondary">Profile</span>
          <select v-model="form.profile" class="input">
            <option value="">default</option>
            <option v-for="p in profiles" :key="p['.id']" :value="p.name">
              {{ p.name }}
            </option>
          </select>
        </label>

        <div class="border-t border-border pt-3 space-y-3">
          <p
            class="text-xs font-semibold text-text-muted uppercase tracking-wide"
          >
            Limits (override profile)
          </p>

          <div class="flex flex-col gap-1">
            <span class="text-xs font-medium text-text-secondary"
              >Time limit</span
            >
            <input
              v-model="form.limitUptimeRaw"
              class="input font-mono"
              placeholder="e.g. 1h, 1d, 1w, 1d12h — blank for unlimited"
              @blur="normalizeFormUptime"
            />
            <p
              v-if="form.limitUptimeRaw && !formUptimePreview"
              class="text-xs text-red"
            >
              Invalid format — use: 30m, 2h, 1d, 1w or combinations like 1d12h
            </p>
            <p
              v-else-if="formUptimePreview"
              class="text-xs"
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
            <span class="text-xs font-medium text-text-secondary"
              >Data limit</span
            >
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
          <span class="text-xs font-medium text-text-secondary">Comment</span>
          <input
            v-model="form.comment"
            class="input"
            placeholder="optional note"
          />
        </label>

        <p v-if="addError" class="text-xs text-red">{{ addError }}</p>

        <div class="flex justify-end gap-2 pt-1">
          <button type="button" class="btn btn-ghost" @click="showAdd = false">
            Cancel
          </button>
          <button
            type="submit"
            class="btn btn-primary"
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

    <!-- Batch generate dialog -->
    <AppDialog
      :open="showBatch"
      title="Generate Users"
      @update:open="onBatchDialogUpdate"
    >
      <!-- Results screen -->
      <div v-if="batchDone" class="space-y-4">
        <div class="flex items-center gap-2 text-sm text-text-secondary">
          <span class="font-bold text-text-primary">{{
            batchResults.filter((r) => r.ok).length
          }}</span>
          created,
          <span class="font-bold text-red">{{
            batchResults.filter((r) => !r.ok).length
          }}</span>
          failed
        </div>
        <div
          class="max-h-64 overflow-y-auto border border-border rounded-lg text-xs font-mono bg-base"
        >
          <div
            v-for="r in batchResults"
            :key="r.name"
            class="flex items-center gap-3 px-3 py-2 border-b border-border last:border-0"
            :class="r.ok ? 'text-text-secondary' : 'text-red bg-red/5'"
          >
            <span class="flex-1">{{ r.name }}</span>
            <span class="text-text-muted">{{ r.password }}</span>
            <span v-if="!r.ok" class="text-xs ml-auto text-red">{{
              r.error
            }}</span>
          </div>
        </div>
        <div class="flex justify-end gap-2 pt-1 border-t border-border">
          <button
            type="button"
            class="btn btn-ghost"
            @click="showBatch = false"
          >
            Close
          </button>
          <button type="button" class="btn btn-primary" @click="printResults">
            <PrinterIcon class="size-4" /> Print
          </button>
        </div>
      </div>

      <!-- Progress screen -->
      <div v-else-if="batchRunning" class="space-y-4">
        <div
          class="flex items-center justify-between text-xs text-text-secondary"
        >
          <span>Creating users…</span>
          <span>{{ batchProgress }} / {{ batch.count }}</span>
        </div>
        <div class="w-full rounded-full h-2 border border-border bg-base">
          <div
            class="h-2 rounded-full transition-all bg-text-primary"
            :style="{ width: `${(batchProgress / batch.count) * 100}%` }"
          />
        </div>
        <p class="text-xs text-text-muted font-mono">{{ batchCurrentName }}</p>
      </div>

      <!-- Config screen -->
      <form v-else @submit.prevent="submitBatch" class="space-y-4">
        <div class="grid grid-cols-2 gap-3">
          <label class="flex flex-col gap-1">
            <span class="text-sm font-medium text-red"
              >Quantity <span>*</span></span
            >
            <input
              v-model.number="batch.count"
              type="number"
              min="1"
              max="500"
              class="input"
              required
            />
            <p v-if="batch.count > 100" class="text-xs text-amber">
              ⚠ Large batch — may take {{ Math.round(batch.count * 0.3) }}s+
            </p>
          </label>
          <label class="flex flex-col gap-1">
            <span class="text-sm font-medium text-text-secondary"
              >Name length</span
            >
            <select v-model.number="batch.nameLength" class="input">
              <option v-for="n in [3, 4, 5, 6]" :key="n" :value="n">
                {{ n }} characters
              </option>
            </select>
          </label>
        </div>

        <div class="flex flex-col gap-1">
          <span class="text-sm font-medium text-text-secondary"
            >Username characters</span
          >
          <div class="flex gap-2">
            <label
              class="flex items-center gap-1.5 text-sm text-text-secondary cursor-pointer"
            >
              <input
                type="checkbox"
                v-model="batch.charsLetters"
                class="rounded"
              />
              Letters (a-z)
            </label>
            <label
              class="flex items-center gap-1.5 text-sm text-text-secondary cursor-pointer"
            >
              <input
                type="checkbox"
                v-model="batch.charsDigits"
                class="rounded"
              />
              Digits (0-9)
            </label>
          </div>
          <p class="text-xs text-text-muted font-mono">
            Preview: <span>{{ batchNamePreview }}</span>
          </p>
        </div>

        <div class="flex flex-col gap-1">
          <span class="text-sm font-medium text-text-secondary">Password</span>
          <select v-model="batch.passwordMode" class="input">
            <option value="same">Same as username</option>
            <option value="random">Random (separate)</option>
            <option value="fixed">Fixed password</option>
          </select>
          <input
            v-if="batch.passwordMode === 'fixed'"
            v-model="batch.fixedPassword"
            class="input mt-1"
            placeholder="Enter fixed password"
          />
        </div>

        <label class="flex flex-col gap-1">
          <span class="text-sm font-medium text-text-secondary">Profile</span>
          <select v-model="batch.profile" class="input">
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
            <span class="text-sm font-medium text-text-secondary"
              >Time limit</span
            >
            <input
              v-model="batch.limitUptimeRaw"
              class="input font-mono"
              placeholder="e.g. 1h, 1d, 1w, 1d12h — blank for unlimited"
              @blur="normalizeBatchUptime"
            />
            <p
              v-if="batch.limitUptimeRaw && !batchUptimePreview"
              class="text-sm text-red"
            >
              Invalid format — use: 30m, 2h, 1d, 1w or combinations like 1d12h
            </p>
            <p
              v-else-if="batchUptimePreview"
              class="text-sm"
              :class="batchUptimeWarning ? 'text-amber' : 'text-text-muted'"
            >
              <span v-if="batchUptimeWarning"
                >⚠ Time limit exceeds the profile's validity — user may never
                hit this limit.</span
              >
              <span v-else class="font-mono"
                >Sends to router: <span>{{ batchUptimePreview }}</span></span
              >
            </p>
          </div>

          <div class="flex flex-col gap-1">
            <span class="text-sm font-medium text-text-secondary"
              >Data limit</span
            >
            <div
              class="flex w-full overflow-hidden rounded-lg border border-border focus-within:outline-2 focus-within:outline-accent focus-within:outline-offset-1"
            >
              <input
                v-model.number="batch.limitBytesTotalValue"
                type="number"
                min="0"
                class="input-bare flex-1 min-w-0"
                placeholder="0"
              />
              <select
                v-model="batch.limitBytesTotalUnit"
                class="input-bare border-l border-border shrink-0 w-16 text-xs"
              >
                <option value="M">MB</option>
                <option value="G">GB</option>
              </select>
            </div>
          </div>
        </div>

        <label class="flex flex-col gap-1">
          <span class="text-sm font-medium text-text-secondary">Comment</span>
          <input
            v-model="batch.comment"
            class="input"
            placeholder="optional note"
          />
        </label>

        <p v-if="batchError" class="text-xs text-red">{{ batchError }}</p>

        <div class="flex justify-end gap-2 pt-1 border-t border-border">
          <button
            type="button"
            class="btn btn-ghost"
            @click="showBatch = false"
          >
            Cancel
          </button>
          <button
            type="submit"
            class="btn btn-primary"
            :disabled="
              !batch.count ||
              batch.count > 500 ||
              !batchCharset ||
              (!!batch.limitUptimeRaw && !batchUptimePreview)
            "
          >
            Generate {{ batch.count || "" }} users
          </button>
        </div>
      </form>
    </AppDialog>
  </PageLayout>
</template>

<script setup lang="ts">
import { ref, computed, watch } from "vue";
import {
  PlusIcon,
  TrashIcon,
  ExclamationTriangleIcon,
  TicketIcon,
  PrinterIcon,
  CheckIcon,
  PencilIcon,
  MagnifyingGlassIcon,
  ChevronLeftIcon,
  ChevronRightIcon,
} from "@heroicons/vue/24/outline";
import { useRoutersStore } from "@/stores/routers";
import {
  listHotspotUsers,
  listHotspotActive,
  listHotspotProfiles,
  createHotspotUser,
  updateHotspotUser,
  toggleHotspotUser,
  deleteHotspotUser,
  getProfileMetas,
  getHotspotSettings,
  type ProfileMeta,
  type HotspotSettings,
} from "@/api";
import { friendlyError } from "@/utils/errors";
import AppDialog from "@/components/AppDialog.vue";
import PageLayout from "@/components/PageLayout.vue";

const store = useRoutersStore();

const tab = ref<"users" | "active">("users");
const tabs = [
  { key: "users" as const, label: "All Users" },
  { key: "active" as const, label: "Active Sessions" },
];
const users = ref<Record<string, string>[]>([]);
const active = ref<Record<string, string>[]>([]);

const searchQuery = ref("");
const filterProfile = ref("");
const filterStatus = ref("");

const filteredUsers = computed(() => {
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
  if (filterStatus.value)
    list = list.filter((u) => userStatus(u) === filterStatus.value);
  return list;
});

const PAGE_SIZES = [20, 50, 100] as const;
const usersPageSize = ref<20 | 50 | 100>(20);
const usersPage = ref(1);
const activePageSize = ref<20 | 50 | 100>(20);
const activePage = ref(1);

const usersPageCount = computed(() =>
  Math.max(1, Math.ceil(filteredUsers.value.length / usersPageSize.value)),
);
const activePageCount = computed(() =>
  Math.max(1, Math.ceil(active.value.length / activePageSize.value)),
);

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
  dnsName: "",
  currency: "",
});
const loading = ref(false);
const error = ref("");

const selected = ref<Set<string>>(new Set());
const bulkDeleting = ref(false);

const allSelected = computed(
  () =>
    users.value.length > 0 &&
    users.value.every((u) => selected.value.has(u[".id"])),
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
    : new Set(users.value.map((u) => u[".id"]));
}

async function removeSelected() {
  if (!store.activeId || !selected.value.size) return;
  if (
    !confirm(
      `Delete ${selected.value.size} selected user${selected.value.size > 1 ? "s" : ""}?`,
    )
  )
    return;
  bulkDeleting.value = true;
  const ids = [...selected.value];
  await Promise.allSettled(
    ids.map((id) => deleteHotspotUser(store.activeId!, id)),
  );
  selected.value = new Set();
  bulkDeleting.value = false;
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

const showBatch = ref(false);
const batchRunning = ref(false);
const batchDone = ref(false);
const batchProgress = ref(0);
const batchCurrentName = ref("");
const batchError = ref("");
const batchResults = ref<
  { name: string; password: string; ok: boolean; error?: string }[]
>([]);

const emptyBatch = () => ({
  count: 10,
  nameLength: 6,
  charsLetters: true,
  charsDigits: true,
  passwordMode: "same" as "same" | "random" | "fixed",
  fixedPassword: "",
  profile: "",
  limitUptimeRaw: "",
  limitBytesTotalValue: 0,
  limitBytesTotalUnit: "M" as "M" | "G",
  comment: "",
});
const batch = ref(emptyBatch());

const batchCharset = computed(() => {
  let s = "";
  if (batch.value.charsLetters) s += "abcdefghijklmnopqrstuvwxyz";
  if (batch.value.charsDigits) s += "0123456789";
  return s;
});

const batchNamePreview = computed(() => {
  if (!batchCharset.value) return "—";
  return Array.from(
    { length: batch.value.nameLength },
    () =>
      batchCharset.value[Math.floor(Math.random() * batchCharset.value.length)],
  ).join("");
});

const batchUptimePreview = computed(() =>
  uptimePreviewFrom(batch.value.limitUptimeRaw),
);
const batchUptimeWarning = computed(() => {
  if (!batchUptimePreview.value || !batch.value.profile) return false;
  const validitySecs = selectedProfileValiditySeconds(batch.value.profile);
  if (!validitySecs) return false;
  return shorthandToSeconds(batch.value.limitUptimeRaw) > validitySecs;
});
function normalizeBatchUptime() {
  batch.value.limitUptimeRaw = normalizeShorthand(batch.value.limitUptimeRaw);
}

function generateName(
  charset: string,
  length: number,
  existing: Set<string>,
): string {
  let name = "",
    attempts = 0;
  do {
    name = Array.from(
      { length },
      () => charset[Math.floor(Math.random() * charset.length)],
    ).join("");
    attempts++;
  } while (existing.has(name) && attempts < 1000);
  return name;
}

function generatePassword(
  mode: "same" | "random" | "fixed",
  name: string,
  fixed: string,
  charset: string,
  length: number,
): string {
  if (mode === "same") return name;
  if (mode === "fixed") return fixed;
  return Array.from(
    { length },
    () => charset[Math.floor(Math.random() * charset.length)],
  ).join("");
}

function switchTab(key: "users" | "active") {
  tab.value = key;
  if (key === "active") loadActive();
}

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
          ({ hotspotName: "", dnsName: "", currency: "" }) as HotspotSettings,
      ),
    ]);
    users.value = u;
    profiles.value = p;
    profileMetas.value = m;
    hotspotSettings.value = s;
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
  if (!store.activeId || !confirm("Delete this user?")) return;
  await deleteHotspotUser(store.activeId, userId);
  await loadUsers();
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

const showEdit = ref(false);
const editing = ref(false);
const editError = ref("");
const editingId = ref("");

const editForm = ref({
  name: "",
  password: "",
  profile: "",
  limitUptimeRaw: "",
  limitBytesTotalValue: 0,
  limitBytesTotalUnit: "M" as "M" | "G",
  comment: "",
  disabled: false,
});

const editUptimePreview = computed(() =>
  uptimePreviewFrom(editForm.value.limitUptimeRaw),
);
function normalizeEditUptime() {
  editForm.value.limitUptimeRaw = normalizeShorthand(
    editForm.value.limitUptimeRaw,
  );
}

function rosUptimeToShorthand(s: string): string {
  if (!s || s === "00:00:00") return "";
  const m = s.match(/^(?:(\d+)d\s*)?(\d+):(\d+):(\d+)$/);
  if (!m) return s;
  const d = parseInt(m[1] || "0"),
    h = parseInt(m[2]),
    min = parseInt(m[3]);
  return [d ? `${d}d` : "", h ? `${h}h` : "", min ? `${min}m` : ""]
    .filter(Boolean)
    .join("");
}

function openEdit(u: Record<string, string>) {
  editingId.value = u[".id"];
  editForm.value = {
    name: u.name ?? "",
    password: "",
    profile: u.profile ?? "",
    limitUptimeRaw: rosUptimeToShorthand(u["limit-uptime"] ?? ""),
    limitBytesTotalValue: u["limit-bytes-total"]
      ? Math.round(parseInt(u["limit-bytes-total"]) / 1024 ** 2)
      : 0,
    limitBytesTotalUnit: "M",
    comment: displayComment(u.comment) === "—" ? "" : displayComment(u.comment),
    disabled: u.disabled === "true",
  };
  editError.value = "";
  showEdit.value = true;
}

async function submitEdit() {
  if (!store.activeId || !editingId.value) return;
  editing.value = true;
  editError.value = "";
  try {
    const f = editForm.value;
    const limitUptime = editUptimePreview.value;
    const mul = f.limitBytesTotalUnit === "G" ? 1024 ** 3 : 1024 ** 2;
    const limitBytesTotal = f.limitBytesTotalValue
      ? String(f.limitBytesTotalValue * mul)
      : "";
    await updateHotspotUser(store.activeId, editingId.value, {
      password: f.password,
      profile: f.profile,
      limitUptime,
      limitBytesTotal,
      comment: f.comment,
    });
    await toggleHotspotUser(store.activeId, editingId.value, f.disabled);
    showEdit.value = false;
    await loadUsers();
  } catch (e: any) {
    editError.value = friendlyError(e, "Failed to update user");
  } finally {
    editing.value = false;
  }
}

function openBatch() {
  batch.value = emptyBatch();
  batchDone.value = false;
  batchRunning.value = false;
  batchResults.value = [];
  batchError.value = "";
  showBatch.value = true;
}

function onBatchDialogUpdate(open: boolean) {
  if (!batchRunning.value) showBatch.value = open;
}

async function submitBatch() {
  if (!store.activeId) return;
  const charset = batchCharset.value;
  if (!charset) {
    batchError.value = "Select at least one character type.";
    return;
  }
  batchRunning.value = true;
  batchProgress.value = 0;
  batchResults.value = [];
  batchError.value = "";
  const b = batch.value;
  const limitUptime = batchUptimePreview.value;
  const mul = b.limitBytesTotalUnit === "G" ? 1024 ** 3 : 1024 ** 2;
  const limitBytesTotal = b.limitBytesTotalValue
    ? String(b.limitBytesTotalValue * mul)
    : "";
  const usedNames = new Set<string>();
  for (let i = 0; i < b.count; i++) {
    const name = generateName(charset, b.nameLength, usedNames);
    usedNames.add(name);
    const password = generatePassword(
      b.passwordMode,
      name,
      b.fixedPassword,
      charset,
      b.nameLength,
    );
    batchCurrentName.value = name;
    try {
      await createHotspotUser(store.activeId, {
        name,
        password,
        profile: b.profile,
        limitUptime,
        limitBytesTotal,
        rateLimit: "",
        comment: b.comment,
        expiryComment: "",
        price: profileMetas.value[b.profile]?.price ?? "",
        currency: hotspotSettings.value.currency ?? "",
      });
      batchResults.value.push({ name, password, ok: true });
    } catch (e: any) {
      batchResults.value.push({
        name,
        password,
        ok: false,
        error: e?.response?.data?.error ?? e?.message ?? "error",
      });
    }
    batchProgress.value = i + 1;
  }
  batchRunning.value = false;
  batchDone.value = true;
  await loadUsers();
}

function printVouchers(
  entries: { name: string; password: string; profile?: string }[],
) {
  const v = hotspotSettings.value.voucher;
  const businessName = v?.businessName ?? hotspotSettings.value.hotspotName ?? "";
  const showValidity = v?.showValidity ?? true;
  const showPrice = v?.showPrice ?? true;
  const currency = hotspotSettings.value.currency ?? "";
  const layout = v?.layout ?? "card";

  const items = entries.map((r) => {
    const meta = profileMetas.value[r.profile || "default"];
    const validity = meta?.validity ?? "";
    const price = meta?.price ?? "";
    const priceStr = showPrice && price ? `${price}${currency ? " " + currency : ""}` : "";
    const uptimeStr = showValidity && validity ? validity : "";
    return { name: r.name, password: r.password, priceStr, uptimeStr };
  });

  let css = "";
  let body = "";

  if (layout === "ticket") {
    css = `
      *{box-sizing:border-box;margin:0;padding:0}
      body{font-family:ui-sans-serif,system-ui,sans-serif;background:#fff;padding:8mm;counter-reset:voucher}
      .grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(88mm,1fr));gap:5mm}
      .card{border:1px solid #d1d5db;border-radius:6px;overflow:hidden;page-break-inside:avoid;counter-increment:voucher}
      .header{background:#111827;color:#fff;padding:3mm 5mm;display:flex;align-items:center;justify-content:space-between}
      .biz{font-size:10pt;font-weight:700;letter-spacing:-.01em}
      .price{font-size:10pt;font-weight:700}
      .body{padding:4mm 5mm;display:flex;flex-direction:column;gap:0}
      .row{display:flex;align-items:center;justify-content:space-between;padding:2.5mm 0}
      .row+.row{border-top:1px solid #f3f4f6}
      .lbl{font-size:7pt;color:#9ca3af;text-transform:uppercase;letter-spacing:.05em}
      .val{font-size:14pt;font-weight:700;color:#111827;font-family:ui-monospace,monospace}
      .val-sm{font-size:9pt;font-weight:500;color:#374151;font-family:ui-sans-serif,system-ui,sans-serif}
      .num{font-size:6pt;color:#d1d5db;margin-left:auto;padding-left:3mm}
      @media print{body{padding:4mm}.grid{gap:4mm}}
    `;
    const cards = items.map(({ name, password, priceStr, uptimeStr }, i) => {
      const validityRow = uptimeStr
        ? `<div class="row"><span class="lbl">Valid for</span><span class="val-sm">${uptimeStr}</span></div>`
        : "";
      return `<div class="card"><div class="header"><span class="biz">${businessName || "WiFi Voucher"}</span>${priceStr ? `<span class="price">${priceStr}</span>` : ""}<span class="num">#${i + 1}</span></div><div class="body"><div class="row"><span class="lbl">Username</span><span class="val">${name}</span></div><div class="row"><span class="lbl">Password</span><span class="val">${password}</span></div>${validityRow}</div></div>`;
    }).join("");
    body = `<div class="grid">${cards}</div>`;

  } else {
    // card (default) — 6-up compact grid
    css = `
      *{box-sizing:border-box;margin:0;padding:0}
      body{font-family:ui-sans-serif,system-ui,sans-serif;background:#fff;padding:8mm}
      .grid{display:grid;grid-template-columns:repeat(6,1fr);gap:4mm}
      .card{border:1px solid #d1d5db;border-radius:6px;padding:3mm 3.5mm;display:flex;flex-direction:column;gap:2.5mm;page-break-inside:avoid}
      .header{display:flex;justify-content:space-between;align-items:baseline;border-bottom:1px solid #e5e7eb;padding-bottom:1.5mm}
      .validity{font-size:7pt;color:#6b7280}.price{font-size:8pt;font-weight:700;color:#111827}
      .creds{display:grid;grid-template-columns:1fr 1px 1fr;gap:2mm;align-items:center}
      .cred-col{display:flex;flex-direction:column;gap:.5mm;align-items:center}
      .divider{width:1px;background:#e5e7eb;align-self:stretch}
      .lbl{font-size:6pt;color:#9ca3af;text-transform:uppercase;letter-spacing:.04em}
      .val{font-size:9pt;font-weight:700;color:#111827;font-family:ui-monospace,monospace}
      .biz{font-size:7pt;color:#9ca3af;text-align:center;margin-top:auto;padding-top:1mm;border-top:1px solid #f3f4f6}
      .num{font-size:6pt;color:#d1d5db;text-align:right;margin-top:auto;padding-top:1mm}
      @media print{body{padding:4mm}.grid{gap:3mm}}
    `;
    const cards = items.map(({ name, password, priceStr, uptimeStr }, i) => {
      const headerLine = priceStr || uptimeStr
        ? `<div class="header"><span class="validity">${uptimeStr}</span><span class="price">${priceStr}</span></div>`
        : "";
      const bizLine = businessName ? `<div class="biz">${businessName}</div>` : "";
      return `<div class="card">${headerLine}<div class="creds"><div class="cred-col"><div class="lbl">Username</div><div class="val">${name}</div></div><div class="divider"></div><div class="cred-col"><div class="lbl">Password</div><div class="val">${password}</div></div></div>${bizLine}<div class="num">#${i + 1}</div></div>`;
    }).join("");
    body = `<div class="grid">${cards}</div>`;
  }

  const html = `<!DOCTYPE html><html><head><meta charset="UTF-8"/><title>Vouchers</title><style>${css}</style></head><body>${body}</body></html>`;

  const iframe = document.createElement("iframe");
  iframe.style.cssText = "position:fixed;top:-9999px;left:-9999px;width:1px;height:1px";
  document.body.appendChild(iframe);
  iframe.srcdoc = html;
  iframe.onload = () => {
    iframe.contentWindow?.print();
    setTimeout(() => iframe.remove(), 1000);
  };
}

function printSelected() {
  const entries = users.value
    .filter((u) => selected.value.has(u[".id"]))
    .map((u) => ({
      name: u.name,
      password: u.password ?? "",
      profile: u.profile,
    }));
  printVouchers(entries);
}

function printResults() {
  const profileName = batch.value.profile || "default";
  printVouchers(
    batchResults.value
      .filter((r) => r.ok)
      .map((r) => ({
        name: r.name,
        password: r.password,
        profile: profileName,
      })),
  );
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
): "disabled" | "expired" | "limit-reached" | "waiting" | "active" {
  if (u.disabled === "true") return "disabled";
  const epoch = extractExpEpoch(u.comment);
  if (epoch !== null && epoch < Math.floor(Date.now() / 1000)) return "expired";
  if (isUptimeExhausted(u)) return "limit-reached";
  if (epoch === null) return "waiting";
  return "active";
}

function expiryLabel(u: Record<string, string>): string {
  const epoch = extractExpEpoch(u.comment);
  if (epoch === null) return "—";
  const now = Math.floor(Date.now() / 1000);
  const diff = epoch - now;
  if (diff <= 0) return "Expired";
  if (diff < 3600) return `${Math.floor(diff / 60)}m left`;
  if (diff < 86400) return `${Math.floor(diff / 3600)}h left`;
  return `${Math.floor(diff / 86400)}d left`;
}

function expiryClass(u: Record<string, string>): string {
  const epoch = extractExpEpoch(u.comment);
  if (epoch === null) return "text-text-muted";
  const diff = epoch - Math.floor(Date.now() / 1000);
  if (diff <= 0) return "text-red font-medium";
  if (diff < 3600) return "text-amber";
  return "text-text-secondary";
}

function parseUptimeSeconds(uptime: string | undefined): number {
  if (!uptime) return 0;
  const m = uptime.match(
    /(?:(\d+)d\s*)?(?:(\d+)h\s*)?(?:(\d+)m\s*)?(?:(\d+)s)?/,
  );
  if (!m) return 0;
  return (
    parseInt(m[1] || "0") * 86400 +
    parseInt(m[2] || "0") * 3600 +
    parseInt(m[3] || "0") * 60 +
    parseInt(m[4] || "0")
  );
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
</script>
