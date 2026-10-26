import QRCode from "qrcode";
import type { VoucherTemplateOptions } from "./types";

// Renders a QR code as an inline SVG string (no network/canvas dependency —
// vouchers must be printable offline). Returns "" on failure so a broken QR
// never blocks printing the rest of the voucher.
export async function qrSvg(text: string, sizePx: number): Promise<string> {
  try {
    return await QRCode.toString(text, {
      type: "svg",
      margin: 0,
      width: sizePx,
    });
  } catch {
    return "";
  }
}

/**
 * QR for one voucher's login link (see business.ts for when credentials are
 * included). "" when the hotspot has no login URL, so templates just omit it.
 */
export async function loginQrSvg(
  name: string,
  password: string,
  opts: Pick<VoucherTemplateOptions, "loginUrl" | "loginUrlSupportsCredentials">,
  sizePx: number,
): Promise<string> {
  const { loginUrl, loginUrlSupportsCredentials } = opts;
  if (!loginUrl) return "";
  const target = loginUrlSupportsCredentials
    ? `${loginUrl}${loginUrl.includes("?") ? "&" : "?"}username=${encodeURIComponent(name)}&password=${encodeURIComponent(password)}`
    : loginUrl;
  return qrSvg(target, sizePx);
}
