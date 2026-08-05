import { ref, computed, watch, type Ref } from "vue";

/**
 * Shared search-filter + pagination logic for simple client-side tables
 * (Logs, Hosts, DHCP leases, etc.) — a free-text search box that filters
 * `items` via `matchText`, paginated at `pageSize` per page, resetting to
 * page 1 whenever the search query changes.
 */
export function useSearchPagination<T>(
  items: Ref<T[]>,
  matchText: (item: T) => string,
  pageSize = 20,
) {
  const search = ref("");

  const filtered = computed(() => {
    const q = search.value.trim().toLowerCase();
    if (!q) return items.value;
    return items.value.filter((item) => matchText(item).toLowerCase().includes(q));
  });

  const page = ref(1);
  const pageCount = computed(() =>
    Math.max(1, Math.ceil(filtered.value.length / pageSize)),
  );
  const paged = computed(() => {
    const start = (page.value - 1) * pageSize;
    return filtered.value.slice(start, start + pageSize);
  });

  watch(search, () => {
    page.value = 1;
  });

  function reset() {
    search.value = "";
    page.value = 1;
  }

  return { search, filtered, page, pageCount, paged, pageSize, reset };
}
