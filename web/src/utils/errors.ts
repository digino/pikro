export function friendlyError(e: any, fallback = 'An unexpected error occurred'): string {
  const raw: string = e?.response?.data?.error ?? e?.message ?? fallback

  // Connection timeout / unreachable
  const timeoutMatch = raw.match(/dial tcp ([\d.]+):\d+.*timeout/i)
  if (timeoutMatch) return `Cannot reach router at ${timeoutMatch[1]} — check the IP and that port 8728 is open`

  const refusedMatch = raw.match(/dial tcp ([\d.]+):\d+.*refused/i)
  if (refusedMatch) return `Connection refused by ${refusedMatch[1]} — RouterOS API may be disabled`

  const noRouteMatch = raw.match(/dial tcp ([\d.]+):\d+.*no route/i)
  if (noRouteMatch) return `No route to ${noRouteMatch[1]} — router may be offline or on a different network`

  // Auth failure
  if (/invalid user name or password|login failed|not permitted/i.test(raw))
    return 'Authentication failed — check username and password'

  // RouterOS API error (already a clean message)
  if (raw.length < 120) return raw

  return fallback
}
