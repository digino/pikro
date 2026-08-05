// A RouterOS hotspot dns-name is a bare hostname (e.g. "myspot.spot") — it
// must never carry a scheme, path, or trailing slash.
export function normalizeDnsName(raw: string): string {
  return raw
    .trim()
    .replace(/^[a-z][a-z0-9+.-]*:\/\//i, "") // strip a leading scheme, e.g. "http://", "https://"
    .replace(/\/.*$/, ""); // strip any path/trailing slash
}
