// Calls fn every everyMs while the page is visible, never two calls at once:
// the next one waits for the previous to finish. A slow router (a call longer
// than a third of everyMs) gets three times its call's length between calls,
// so polling never keeps a small router's CPU busy. Returns a stop function.
export function startPolling(fn: () => Promise<unknown>, everyMs: number): () => void {
  let timer: ReturnType<typeof setTimeout> | undefined;
  let stopped = false;
  async function tick() {
    const started = Date.now();
    if (!document.hidden) await fn().catch(() => {});
    if (!stopped) timer = setTimeout(tick, Math.max(everyMs, (Date.now() - started) * 3));
  }
  timer = setTimeout(tick, everyMs);
  return () => {
    stopped = true;
    clearTimeout(timer);
  };
}
