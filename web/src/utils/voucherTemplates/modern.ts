import type { VoucherTemplate } from "./types";

// Modern — compact professional card grid, heavy borders, small gap between
// pieces (as opposed to Classic's zero-gap collapsed-border sheet).
const render: VoucherTemplate["render"] = async (items, { businessName }) => {
  const cards = items
    .map(
      ({ name, password, validity, price }, i) => `
        <div class="v">
          ${businessName ? `<div class="biz">${businessName}</div>` : ""}
          <div class="creds">
            <div class="col"><span class="lbl">Username</span><span class="val">${name}</span></div>
            <div class="divider"></div>
            <div class="col"><span class="lbl">Password</span><span class="val">${password}</span></div>
          </div>
          ${validity || price ? `<div class="sub">${[validity, price].filter(Boolean).join(" · ")}</div>` : ""}
          <div class="num">#${i + 1}</div>
        </div>`,
    )
    .join("");

  const css = `
    *{box-sizing:border-box;margin:0;padding:0}
    body{font-family:ui-sans-serif,system-ui,sans-serif;background:#fff;color:#000;padding:8mm}
    .grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(34mm,1fr));gap:2mm}
    .v{border:1.25pt solid #000;border-radius:2px;padding:2.5mm;display:flex;flex-direction:column;gap:1.5mm;page-break-inside:avoid}
    .biz{font-size:6.5pt;font-weight:700;text-align:center;text-transform:uppercase;letter-spacing:.02em;border-bottom:0.75pt solid #000;padding-bottom:1mm}
    .creds{display:grid;grid-template-columns:1fr 1px 1fr;gap:1.5mm;align-items:center}
    .col{display:flex;flex-direction:column;align-items:center;gap:0.4mm}
    .divider{width:1px;background:#000;align-self:stretch}
    .lbl{font-size:5.5pt;text-transform:uppercase;letter-spacing:.03em;color:#333}
    .val{font-size:9pt;font-weight:700;font-family:ui-monospace,monospace}
    .sub{font-size:6pt;text-align:center;border-top:0.75pt solid #000;padding-top:1mm}
    .num{font-size:5.5pt;color:#666;text-align:right}
    @media print{body{padding:5mm}.grid{gap:1.5mm}}
  `;

  return `<!DOCTYPE html><html><head><meta charset="UTF-8"/><title>Vouchers — Modern</title><style>${css}</style></head><body><div class="grid">${cards}</div></body></html>`;
};

export const MODERN_TEMPLATE: VoucherTemplate = {
  key: "modern",
  label: "Modern",
  description: "Compact professional cards — heavy borders, light spacing, ~60 per page",
  render,
};
