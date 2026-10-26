import { MODERN_TEMPLATE } from "./modern";
import { BUSINESS_TEMPLATE } from "./business";
import { VOUCHER_PRINT_TEMPLATE } from "./voucher";
import type { VoucherTemplate } from "./types";

export * from "./types";

export const VOUCHER_TEMPLATES: VoucherTemplate[] = [
  MODERN_TEMPLATE,
  BUSINESS_TEMPLATE,
  VOUCHER_PRINT_TEMPLATE,
];

export function getVoucherTemplate(key: VoucherTemplate["key"]): VoucherTemplate {
  return VOUCHER_TEMPLATES.find((t) => t.key === key) ?? MODERN_TEMPLATE;
}
