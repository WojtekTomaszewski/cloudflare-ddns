package main

type Zone struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Zones struct {
	Result []Zone `json:"result"`
}

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

type Records struct {
	Result []Record `json:"result"`
}
