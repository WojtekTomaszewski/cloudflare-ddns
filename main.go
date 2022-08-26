package main

import (
	"log"
	"os"
)

var (
	token      = os.Getenv("TOKEN")
	zone       = os.Getenv("ZONE")
	subdomain  = os.Getenv("SUBDOMAIN")
	recordType = os.Getenv("TYPE")
)

func main() {

	ip, err := getCurrentIP()
	if err != nil {
		log.Fatal("failed to get current ip address", err)
	}

	if !isIpValid(ip) {
		log.Fatal("invalid ip address", ip)
	}

	cf := newClient(token)

	var zoneId string
	zones, err := cf.getZones()
	if err != nil {
		log.Fatal("failed to get zones", err)
	}
	for _, z := range zones.Result {
		if z.Name == zone {
			zoneId = z.ID
			break
		}
	}

	recordName := zone
	if subdomain != "" {
		recordName = subdomain
	}

	records, err := cf.getDnsRecord(zoneId, recordType, recordName)
	if err != nil {
		log.Fatal("failed to get dns records", err)
	}

	if len(records.Result) == 0 {
		log.Fatal("no dns records found")
	}

	if records.Result[0].Content == ip {
		log.Println("ip address has not changed", ip)
		return
	}

	updatedRecord := &Record{
		Name:    records.Result[0].Name,
		Type:    records.Result[0].Type,
		TTL:     records.Result[0].TTL,
		Proxied: records.Result[0].Proxied,
		Content: ip,
	}

	if err := cf.updateDnsRecord(zoneId, records.Result[0].ID, updatedRecord); err != nil {
		log.Fatal("failed to update dns record", err)
	}

	log.Printf("updated %s record, old ip %s, new ip %s", records.Result[0].Name, records.Result[0].Content, ip)

}
