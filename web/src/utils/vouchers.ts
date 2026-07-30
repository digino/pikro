import type { VoucherSettings, ProfileMeta } from "@/api";

export interface VoucherEntry {
  name: string;
  password: string;
  profile?: string;
}

export interface PrintVouchersOptions {
  layout: NonNullable<VoucherSettings["layout"]>;
  businessName: string;
  showValidity: boolean;
  showPrice: boolean;
  currency: string;
  profileMetas: Record<string, ProfileMeta>;
}

export const VOUCHER_LAYOUTS: {
  key: NonNullable<VoucherSettings["layout"]>;
  icon: string;
  label: string;
  description: string;
}[] = [
  { key: "card", icon: "▦", label: "Card", description: "6 per page — compact grid, maximises A4 space" },
  { key: "ticket", icon: "▣", label: "Ticket", description: "2 per page — large text, easy to read and hand out" },
];

export function printVouchers(entries: VoucherEntry[], opts: PrintVouchersOptions) {
  const { layout, businessName, showValidity, showPrice, currency, profileMetas } = opts;

  const items = entries.map((r) => {
    const meta = profileMetas[r.profile || "default"];
    const validity = meta?.validity ?? "";
    const price = meta?.price ?? "";
    const priceStr = showPrice && price ? `${price}${currency ? " " + currency : ""}` : "";
    const uptimeStr = showValidity && validity ? validity : "";
    return { name: r.name, password: r.password, priceStr, uptimeStr };
  });

  let css = "";
  let body = "";

  if (layout === "ticket") {
    css = `
      *{box-sizing:border-box;margin:0;padding:0}
      body{font-family:ui-sans-serif,system-ui,sans-serif;background:#fff;padding:8mm;counter-reset:voucher}
      .grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(88mm,1fr));gap:5mm}
      .card{border:1px solid #d1d5db;border-radius:6px;overflow:hidden;page-break-inside:avoid;counter-increment:voucher}
      .header{background:#111827;color:#fff;padding:3mm 5mm;display:flex;align-items:center;justify-content:space-between}
      .biz{font-size:10pt;font-weight:700;letter-spacing:-.01em}
      .price{font-size:10pt;font-weight:700}
      .body{padding:4mm 5mm;display:flex;flex-direction:column;gap:0}
      .row{display:flex;align-items:center;justify-content:space-between;padding:2.5mm 0}
      .row+.row{border-top:1px solid #f3f4f6}
      .lbl{font-size:7pt;color:#9ca3af;text-transform:uppercase;letter-spacing:.05em}
      .val{font-size:14pt;font-weight:700;color:#111827;font-family:ui-monospace,monospace}
      .val-sm{font-size:9pt;font-weight:500;color:#374151;font-family:ui-sans-serif,system-ui,sans-serif}
      .num{font-size:6pt;color:#d1d5db;margin-left:auto;padding-left:3mm}
      @media print{body{padding:4mm}.grid{gap:4mm}}
    `;
    const cards = items
      .map(({ name, password, priceStr, uptimeStr }, i) => {
        const validityRow = uptimeStr
          ? `<div class="row"><span class="lbl">Valid for</span><span class="val-sm">${uptimeStr}</span></div>`
          : "";
        return `<div class="card"><div class="header"><span class="biz">${businessName || "WiFi Voucher"}</span>${priceStr ? `<span class="price">${priceStr}</span>` : ""}<span class="num">#${i + 1}</span></div><div class="body"><div class="row"><span class="lbl">Username</span><span class="val">${name}</span></div><div class="row"><span class="lbl">Password</span><span class="val">${password}</span></div>${validityRow}</div></div>`;
      })
      .join("");
    body = `<div class="grid">${cards}</div>`;
  } else {
    // card (default) — 6-up compact grid
    css = `
      *{box-sizing:border-box;margin:0;padding:0}
      body{font-family:ui-sans-serif,system-ui,sans-serif;background:#fff;padding:8mm}
      .grid{display:grid;grid-template-columns:repeat(6,1fr);gap:4mm}
      .card{border:1px solid #d1d5db;border-radius:6px;padding:3mm 3.5mm;display:flex;flex-direction:column;gap:2.5mm;page-break-inside:avoid}
      .header{display:flex;justify-content:space-between;align-items:baseline;border-bottom:1px solid #e5e7eb;padding-bottom:1.5mm}
      .validity{font-size:7pt;color:#6b7280}.price{font-size:8pt;font-weight:700;color:#111827}
      .creds{display:grid;grid-template-columns:1fr 1px 1fr;gap:2mm;align-items:center}
      .cred-col{display:flex;flex-direction:column;gap:.5mm;align-items:center}
      .divider{width:1px;background:#e5e7eb;align-self:stretch}
      .lbl{font-size:6pt;color:#9ca3af;text-transform:uppercase;letter-spacing:.04em}
      .val{font-size:9pt;font-weight:700;color:#111827;font-family:ui-monospace,monospace}
      .biz{font-size:7pt;color:#9ca3af;text-align:center;margin-top:auto;padding-top:1mm;border-top:1px solid #f3f4f6}
      .num{font-size:6pt;color:#d1d5db;text-align:right;margin-top:auto;padding-top:1mm}
      @media print{body{padding:4mm}.grid{gap:3mm}}
    `;
    const cards = items
      .map(({ name, password, priceStr, uptimeStr }, i) => {
        const headerLine =
          priceStr || uptimeStr
            ? `<div class="header"><span class="validity">${uptimeStr}</span><span class="price">${priceStr}</span></div>`
            : "";
        const bizLine = businessName ? `<div class="biz">${businessName}</div>` : "";
        return `<div class="card">${headerLine}<div class="creds"><div class="cred-col"><div class="lbl">Username</div><div class="val">${name}</div></div><div class="divider"></div><div class="cred-col"><div class="lbl">Password</div><div class="val">${password}</div></div></div>${bizLine}<div class="num">#${i + 1}</div></div>`;
      })
      .join("");
    body = `<div class="grid">${cards}</div>`;
  }

  const html = `<!DOCTYPE html><html><head><meta charset="UTF-8"/><title>Vouchers</title><style>${css}</style></head><body>${body}</body></html>`;

  const iframe = document.createElement("iframe");
  iframe.style.cssText = "position:fixed;top:-9999px;left:-9999px;width:1px;height:1px";
  document.body.appendChild(iframe);
  iframe.srcdoc = html;
  iframe.onload = () => {
    iframe.contentWindow?.print();
    setTimeout(() => iframe.remove(), 1000);
  };
}
