import type { VoucherTemplate } from "./types";
import { voucherDurationLabel, longestCode } from "./types";
import { ticketLabels } from "./labels";
import { loginQrSvg } from "./qr";

const QR_SIZE_PX = 52;

// Voucher — single-code ticket, matching the "Voucher" hotspot login template
// (one field, no separate username/password). Requires vouchers generated
// with password set to "Same as username" — prints item.name only.
// Ticket layout from the v1 design: details and code on the left, the login QR
// on a stub behind a dashed tear line. Black and white only, to print cheaply.
const render: VoucherTemplate["render"] = async (items, opts) => {
  const L = ticketLabels();
  const { businessName } = opts;

  const cards = await Promise.all(
    items.map(async (item) => {
      const { name, password, price } = item;
      const duration = voucherDurationLabel(item);
      const qr = await loginQrSvg(name, password, opts, QR_SIZE_PX);
      return `
        <div class="v">
          <div class="info">
            <div class="head"><span>${businessName}</span><span>${duration}</span></div>
            <div class="code">${name}</div>
            <div class="foot">${L.connectHint}${price ? ` <b>${price}</b>` : ""}</div>
          </div>
          ${qr ? `<div class="stub"><div class="qr">${qr}</div></div>` : ""}
        </div>`;
    }),
  );

  // Code at ~3mm per 14pt mono character, plus padding and the ~19mm QR stub.
  const minWidth = Math.max(60, longestCode(items) * 3 + 30);

  const css = `
    *{box-sizing:border-box;margin:0;padding:0}
    body{font-family:ui-sans-serif,system-ui,sans-serif;background:#fff;color:#000;padding:8mm}
    .grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(${minWidth}mm,1fr));gap:2.5mm}
    .v{display:flex;min-height:22mm;border:1.5pt solid #000;border-radius:2.5mm;overflow:hidden;page-break-inside:avoid}
    .info{flex:1;min-width:0;padding:2.5mm 3mm;display:flex;flex-direction:column;justify-content:center;gap:1mm}
    .head{display:flex;justify-content:space-between;gap:1.5mm;font-size:6.5pt;font-weight:700;letter-spacing:.06em;text-transform:uppercase}
    .code{font-family:ui-monospace,Menlo,Consolas,monospace;font-size:14pt;font-weight:600;letter-spacing:.02em}
    .foot{font-size:6.5pt;color:#333}
    .stub{flex:none;display:flex;align-items:center;justify-content:center;padding:0 2.5mm;border-left:1.5pt dashed #000}
    .qr{display:flex;padding:1mm;border:1pt solid #000;background:#fff}
    @media print{body{padding:5mm}.grid{gap:2mm}}
  `;

  return `<!DOCTYPE html><html><head><meta charset="UTF-8"/><title>Vouchers — Voucher</title><style>${css}</style></head><body><div class="grid">${cards.join("")}</div></body></html>`;
};

export const VOUCHER_PRINT_TEMPLATE: VoucherTemplate = {
  key: "voucher",
  render,
};
