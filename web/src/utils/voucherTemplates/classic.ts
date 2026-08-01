import type { VoucherTemplate } from "./types";

// Classic — dense ticket-sheet layout, heavy black borders, zero gap.
// Borders are collapsed like a table: each cell overlaps its neighbor by
// its own border width (negative margin), so adjacent cells share a single
// heavy line instead of stacking into a double-thick seam.
const render: VoucherTemplate["render"] = async (items) => {
  const cards = items
    .map(
      ({ name, password, validity, price }) => `
        <div class="v">
          <div class="row name"><span class="lbl">User</span><span class="val">${name}</span></div>
          <div class="row"><span class="lbl">Pass</span><span class="val">${password}</span></div>
          ${validity || price ? `<div class="row sub">${[validity, price].filter(Boolean).join(" · ")}</div>` : ""}
        </div>`,
    )
    .join("");

  const css = `
    *{box-sizing:border-box;margin:0;padding:0}
    body{font-family:ui-monospace,Menlo,Consolas,monospace;background:#fff;color:#000;padding:10mm}
    .grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(32mm,1fr))}
    .v{border:1.5pt solid #000;margin:0 0 -1.5pt -1.5pt;padding:2mm;display:flex;flex-direction:column;gap:0.8mm;page-break-inside:avoid}
    .row{display:flex;justify-content:space-between;align-items:baseline;gap:1.5mm}
    .lbl{font-size:6pt;text-transform:uppercase;letter-spacing:.03em}
    .val{font-size:10pt;font-weight:700}
    .sub{font-size:6.5pt;border-top:0.75pt solid #000;padding-top:0.8mm;margin-top:0.3mm;justify-content:center}
    @media print{body{padding:5mm}}
  `;

  return `<!DOCTYPE html><html><head><meta charset="UTF-8"/><title>Vouchers — Classic</title><style>${css}</style></head><body><div class="grid">${cards}</div></body></html>`;
};

export const CLASSIC_TEMPLATE: VoucherTemplate = {
  key: "classic",
  label: "Classic",
  description: "Dense ticket sheet — heavy borders, zero gap, ~60 per page",
  render,
};
