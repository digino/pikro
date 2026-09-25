import type { VoucherTemplate } from "./types";
import { voucherDurationLabel } from "./types";

// Voucher — single-code ticket, matching the "Voucher" hotspot login template
// (one field, no separate username/password). Requires vouchers generated
// with password set to "Same as username" — prints item.name only.
const render: VoucherTemplate["render"] = async (items, { businessName }) => {
  const cards = items
    .map((item) => {
      const { name, price } = item;
      const duration = voucherDurationLabel(item);
      return `
        <div class="v">
          ${businessName ? `<div class="biz">${businessName}</div>` : ""}
          <div class="lbl">Voucher code</div>
          <div class="code">${name}</div>
          ${duration || price ? `<div class="sub">${[duration, price].filter(Boolean).join(" · ")}</div>` : ""}
        </div>`;
    })
    .join("");

  const css = `
    *{box-sizing:border-box;margin:0;padding:0}
    body{font-family:ui-sans-serif,system-ui,sans-serif;background:#fff;color:#000;padding:8mm}
    .grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(36mm,1fr));gap:2mm}
    .v{border:1.25pt solid #000;border-radius:3px;padding:3mm;display:flex;flex-direction:column;align-items:center;gap:1mm;page-break-inside:avoid;text-align:center}
    .biz{font-size:6.5pt;font-weight:700;text-transform:uppercase;letter-spacing:.02em;border-bottom:0.75pt solid #000;padding-bottom:1mm;width:100%}
    .lbl{font-size:5.5pt;text-transform:uppercase;letter-spacing:.05em;color:#333}
    .code{font-size:12pt;font-weight:800;letter-spacing:.06em}
    .sub{font-size:6pt;border-top:0.75pt solid #000;padding-top:1mm;width:100%}
    @media print{body{padding:5mm}.grid{gap:1.5mm}}
  `;

  return `<!DOCTYPE html><html><head><meta charset="UTF-8"/><title>Vouchers — Voucher</title><style>${css}</style></head><body><div class="grid">${cards}</div></body></html>`;
};

export const VOUCHER_PRINT_TEMPLATE: VoucherTemplate = {
  key: "voucher",
  label: "Voucher",
  description: "Single-code ticket — pairs with the Voucher login page, ~50 per page",
  render,
};
