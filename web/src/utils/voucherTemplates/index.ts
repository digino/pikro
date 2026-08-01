import { CLASSIC_TEMPLATE } from "./classic";
import { MODERN_TEMPLATE } from "./modern";
import { BUSINESS_TEMPLATE } from "./business";
import type { VoucherTemplate } from "./types";

export * from "./types";

export const VOUCHER_TEMPLATES: VoucherTemplate[] = [
  CLASSIC_TEMPLATE,
  MODERN_TEMPLATE,
  BUSINESS_TEMPLATE,
];

export function getVoucherTemplate(key: VoucherTemplate["key"]): VoucherTemplate {
  return VOUCHER_TEMPLATES.find((t) => t.key === key) ?? CLASSIC_TEMPLATE;
}
