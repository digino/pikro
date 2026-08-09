import type { ProfileMeta } from "@/api";
import { getVoucherTemplate, type VoucherTemplate, type VoucherItem } from "./voucherTemplates";

export interface VoucherEntry {
  name: string;
  password: string;
  profile?: string;
  /** Time-limit (RouterOS limit-uptime) set for this voucher at generation time, e.g. "24h". */
  timeLimit?: string;
}

export interface PrintVouchersOptions {
  template: VoucherTemplate["key"];
  businessName: string;
  showPrice: boolean;
  currency: string;
  profileMetas: Record<string, ProfileMeta>;
  /** Hotspot login page URL, e.g. "http://myspot.spot/login" — used for the Business template's QR code. */
  loginUrl?: string;
  /** See VoucherTemplateOptions.loginUrlSupportsCredentials. */
  loginUrlSupportsCredentials?: boolean;
}

export async function printVouchers(entries: VoucherEntry[], opts: PrintVouchersOptions) {
  const {
    template,
    businessName,
    showPrice,
    currency,
    profileMetas,
    loginUrl = "",
    loginUrlSupportsCredentials = false,
  } = opts;

  const items: VoucherItem[] = entries.map((r) => {
    const meta = profileMetas[r.profile || "default"];
    const price = showPrice && meta?.price ? `${meta.price}${currency ? " " + currency : ""}` : "";
    return {
      name: r.name,
      password: r.password,
      profile: r.profile ?? "",
      timeLimit: r.timeLimit ?? "",
      price,
    };
  });

  const html = await getVoucherTemplate(template).render(items, {
    businessName,
    loginUrl,
    loginUrlSupportsCredentials,
  });

  const iframe = document.createElement("iframe");
  iframe.style.cssText = "position:fixed;top:-9999px;left:-9999px;width:1px;height:1px";
  document.body.appendChild(iframe);
  iframe.srcdoc = html;
  iframe.onload = () => {
    iframe.contentWindow?.print();
    setTimeout(() => iframe.remove(), 1000);
  };
}
