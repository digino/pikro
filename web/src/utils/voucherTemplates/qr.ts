import QRCode from "qrcode";

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
