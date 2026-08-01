export interface VoucherItem {
  name: string;
  password: string;
  /** Pre-formatted validity string (e.g. "1 day"), empty if hidden. */
  validity: string;
  /** Pre-formatted price string (e.g. "500 XOF"), empty if hidden. */
  price: string;
}

export interface VoucherTemplateOptions {
  businessName: string;
  /**
   * Login URL to encode in a QR code (e.g. "http://myspot.spot/login"),
   * or empty if the hotspot's auth mode doesn't support it (see loginUrlSupportsCredentials).
   */
  loginUrl: string;
  /**
   * True when the hotspot profile uses an auth mode (http-pap) that accepts
   * plaintext username/password as URL query parameters — required for a
   * "scan to connect" QR code to actually log the user in. http-chap
   * (the default) rejects this, so the QR falls back to linking the login
   * page only, without prefilled credentials.
   */
  loginUrlSupportsCredentials: boolean;
}

export interface VoucherTemplate {
  key: "classic" | "modern" | "business";
  label: string;
  description: string;
  /** Renders the full printable HTML document for a batch of vouchers. */
  render(items: VoucherItem[], opts: VoucherTemplateOptions): Promise<string>;
}
