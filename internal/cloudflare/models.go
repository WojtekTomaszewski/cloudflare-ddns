package cloudflare

// Zone represents Cloudflare zone object, we need only id and name for further processing
type Zone struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Zones is a list of Cloudflare zones
type Zones struct {
	Result  []Zone  `json:"result"`
	Success bool    `json:"success"`
	Errors  []Error `json:"errors,omitempty"`
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

// ListRecords is response object for list DNS records call
type ListRecords struct {
	Result  []Record `json:"result"`
	Success bool     `json:"success"`
	Errors  []Error  `json:"errors,omitempty"`
}

// PutRecord is response object for update DNS record call
type PutRecord struct {
	Result  Record  `json:"result"`
	Success bool    `json:"success"`
	Errors  []Error `json:"errors,omitempty"`
}

// Error represents errors object in 4xx responses
type Error struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}
