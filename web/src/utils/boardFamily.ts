export type BoardFamily =
  | "hap"
  | "hex"
  | "ccr"
  | "rack"
  | "ap"
  | "outdoor"
  | "generic";

export function boardFamily(boardName: string): BoardFamily {
  const b = boardName.toLowerCase().trim();

  if (b.startsWith("ccr")) return "ccr";
  if (b.startsWith("crs")) return "rack";
  if (/^rb(1100|3011|4011|5009|1009)/.test(b)) return "rack";
  if (/^(lhg|sxt|ldf|metal|loco|groove|omnitik|disc|netmetal|basebox)/.test(b))
    return "outdoor";
  if (/^(wap|cap|rbcap)/.test(b)) return "ap";

  // Strict matching for specific legacy models using a clear prefix strategy
  // RB951G is hex-like/desktop case but has Wi-Fi (it's functionally a hAP).
  // Let's ensure 'rb951' matches accurately.
  if (/^(hap|map|audience|rb951|rb952|rb962|rb941|rb2011)/.test(b))
    return "hap";
  if (/^(hex|rb750|rb760)/.test(b)) return "hex";

  return "generic";
}
