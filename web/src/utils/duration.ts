// Operator-typed durations ("30m", "2h", "1d12h", "1w") for RouterOS limit-uptime.

export function parseShorthand(
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

export function shorthandToSeconds(s: string): number {
  const p = parseShorthand(s);
  if (!p) return 0;
  return p.w * 7 * 86400 + p.d * 86400 + p.h * 3600 + p.m * 60;
}

/** The value sent to the router: weeks folded into days ("1w2d" -> "9d"), "" when invalid. */
export function uptimeForRouter(raw: string): string {
  const p = parseShorthand(raw);
  if (!p) return "";
  const totalDays = p.w * 7 + p.d;
  if (!totalDays && !p.h && !p.m) return "";
  return [totalDays ? `${totalDays}d` : "", p.h ? `${p.h}h` : "", p.m ? `${p.m}m` : ""]
    .filter(Boolean)
    .join("");
}

/** Tidies what the operator typed ("1W02h" -> "1w2h"); leaves invalid input untouched. */
export function normalizeShorthand(raw: string): string {
  const p = parseShorthand(raw);
  if (!p) return raw;
  return [p.w ? `${p.w}w` : "", p.d ? `${p.d}d` : "", p.h ? `${p.h}h` : "", p.m ? `${p.m}m` : ""]
    .filter(Boolean)
    .join("");
}
