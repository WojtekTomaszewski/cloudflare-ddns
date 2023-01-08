package cloudflare

// Zone represents Cloudflare zone object, we need only id and name for further processing
type Zone struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Zones is a list of Cloudflare zones
type Zones struct {
	Result []Zone `json:"result"`
}

// Record reprsents Cloudflare DNS record
type Record struct {
	ID       string `json:"id,omitempty"`
	ZoneID   string `json:"zone_id,omitempty"`
	ZoneName string `json:"zone_name,omitempty"`
	Name     string `json:"name,omitempty"`
	TTL      int    `json:"ttl,omitempty"`
	Type     string `json:"type,omitempty"`
	Content  string `json:"content,omitempty"`
	Proxied  bool   `json:"proxied,omitempty"`
}

// Record is a list of Cloudflare records
type Records struct {
	Result []Record `json:"result"`
}
