<template>
  <PageLayout :title="t('about.title')">
    <section class="card p-6 flex flex-wrap items-center gap-5 max-w-2xl">
      <span
        aria-hidden="true"
        class="type-display inline-flex items-center justify-center size-16 shrink-0 rounded-card border-2 border-ink bg-primary text-primary-fg text-4xl leading-none -rotate-4"
      >p</span>
      <div class="flex-1 min-w-0 flex flex-col gap-1">
        <h2 class="type-display m-0 text-3xl">pikro</h2>
        <p class="m-0 text-sm text-text-secondary">{{ t('about.tagline') }}</p>
        <p class="m-0 mt-1 font-mono text-sm font-semibold">
          {{ isDev ? t('about.devBuild') : t('about.version', { version: displayVersion(version) }) }}
        </p>
      </div>
    </section>

    <section class="card p-6 flex flex-wrap items-center gap-4 max-w-2xl" aria-live="polite">
      <span v-if="updates.checking" class="spinner spinner--sm" />
      <template v-else-if="status?.available">
        <span class="chip bg-primary text-primary-fg">
          {{ t('about.updateAvailable', { version: displayVersion(status.latest) }) }}
        </span>
        <a :href="downloadPage(locale)" target="_blank" rel="noopener noreferrer" class="btn btn-primary">
          {{ t('about.download') }}
        </a>
      </template>
      <span v-else-if="status?.error" class="text-sm text-danger-text">{{ t('about.checkFailed') }}</span>
      <span v-else-if="status && !isDev" class="chip bg-paper">
        <span aria-hidden="true" class="size-2 rounded-full bg-online" />{{ t('about.upToDate') }}
      </span>
      <button
        type="button"
        class="btn ml-auto"
        :disabled="updates.checking"
        @click="updates.check(true)"
      >
        {{ t('about.checkAgain') }}
      </button>
    </section>

    <section class="card p-6 flex flex-col gap-3 max-w-2xl" aria-labelledby="links-title">
      <h2 id="links-title" class="eyebrow m-0">{{ t('about.linksTitle') }}</h2>
      <ul class="m-0 p-0 list-none flex flex-col">
        <li v-for="l in links" :key="l.href" class="border-b border-hairline last:border-0">
          <a
            :href="l.href"
            target="_blank"
            rel="noopener noreferrer"
            class="flex items-center justify-between gap-3 min-h-11 py-2 font-semibold text-ink hover:underline"
          >
            {{ l.label }}
            <ArrowTopRightOnSquareIcon class="size-4 shrink-0 text-text-muted" />
          </a>
        </li>
      </ul>
    </section>

    <section class="card p-6 flex flex-col gap-2 max-w-2xl" aria-labelledby="privacy-title">
      <h2 id="privacy-title" class="eyebrow m-0">{{ t('about.privacyTitle') }}</h2>
      <p class="m-0 text-sm leading-normal text-text-secondary">{{ t('about.privacy') }}</p>
    </section>

    <p class="m-0 text-sm text-text-muted">{{ t('about.license') }}</p>
  </PageLayout>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ArrowTopRightOnSquareIcon } from '@heroicons/vue/24/outline'
import PageLayout from '@/components/PageLayout.vue'
import { useUpdateStore, displayVersion } from '@/stores/update'
import { LINKS, downloadPage } from '@/utils/links'

const { t, locale } = useI18n()
const updates = useUpdateStore()

const status = computed(() => updates.status)
const version = computed(() => status.value?.current ?? '')
const isDev = computed(() => !/^v*\d/.test(version.value))

const links = computed(() => [
  { label: t('about.website'), href: locale.value === 'fr' ? LINKS.page.fr : LINKS.page.en },
  { label: t('about.releases'), href: LINKS.releases },
  { label: t('about.source'), href: LINKS.repo },
  { label: t('about.report'), href: LINKS.issues },
])

// The banner's check may not have run yet when the page opens first.
onMounted(() => {
  if (!updates.status) updates.check()
})
</script>
