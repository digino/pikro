import i18n from "@/locales";

// Printed tickets follow the language the operator is using Pikro in.
export function ticketLabels() {
  const { t } = i18n.global;
  return {
    username: t("tickets.username"),
    password: t("tickets.password"),
    connectHint: t("tickets.connectHint"),
  };
}
