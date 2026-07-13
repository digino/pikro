package assets

import _ "embed"

//go:embed hotspot_login.html
var HotspotLoginHTML []byte

//go:embed hotspot_logout.html
var HotspotLogoutHTML []byte

//go:embed hotspot_status.html
var HotspotStatusHTML []byte

//go:embed hotspot_alogin.html
var HotspotAloginHTML []byte

//go:embed hotspot_error.html
var HotspotErrorHTML []byte

//go:embed hotspot_redirect.html
var HotspotRedirectHTML []byte
