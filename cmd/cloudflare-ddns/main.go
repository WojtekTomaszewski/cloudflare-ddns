package main

import (
	"log"
	"os"

	"github.com/WojtekTomaszewski/cloudflare-ddns/internal/cloudflare"
	publicip "github.com/WojtekTomaszewski/cloudflare-ddns/internal/public-ip"
)

var (
	token      = os.Getenv("CLOUDFLARE_TOKEN")
	zone       = os.Getenv("CLOUDFLARE_ZONE")
	subdomain  = os.Getenv("CLOUDFLARE_SUBDOMAIN")
	recordType = os.Getenv("CLOUDFLARE_RECORD_TYPE")
)

func main() {

	token, ok := os.LookupEnv("CLOUDFLARE_TOKEN")
	if !ok {
		log.Fatal("CLOUDFLARE_TOKEN env variable is not set")
	}

	zone, ok := os.LookupEnv("CLOUDFLARE_ZONE")
	if !ok {
		log.Fatal("CLOUDFLARE_ZONE env variable is not set")
	}

	subdomain, ok := os.LookupEnv("CLOUDFLARE_SUBDOMAIN")
	if !ok {
		log.Println("CLOUDFLARE_SUBDOMAIN is not set, using CLOUDFLARE_ZONE")
	}

	recordType, ok := os.LookupEnv("CLOUDFLARE_RECORD_TYPE")
	if !ok {
		recordType = "A"
		log.Println("CLOUDFLARE_RECORD_TYPE is not set, using 'A'")
	}

	ip, err := publicip.GetCurrentIP()
	if err != nil {
		log.Fatal("failed to get current ip address, ", err)
	}

	if !publicip.IsIpValid(ip) {
		log.Fatal("could not validat current public ip address", ip)
	}

	cf := cloudflare.NewClient(token)

	var zoneId string
	zones, err := cf.GetZones()
	if err != nil {
		log.Fatal("failed to get zones, ", err)
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

	records, err := cf.GetDnsRecord(zoneId, recordType, recordName)
	if err != nil {
		log.Fatal("failed to get dns records, ", err)
	}

	if len(records.Result) == 0 {
		log.Fatal("no dns records found")
	}

	if records.Result[0].Content == ip {
		log.Println("ip address has not changed: ", ip)
		return
	}

	updatedRecord := &cloudflare.Record{
		Name:    records.Result[0].Name,
		Type:    records.Result[0].Type,
		TTL:     records.Result[0].TTL,
		Proxied: records.Result[0].Proxied,
		Content: ip,
	}

	if err := cf.UpdateDnsRecord(zoneId, records.Result[0].ID, updatedRecord); err != nil {
		log.Fatal("failed to update dns record", err)
	}

	log.Printf("updated %s record, old ip %s, new ip %s", records.Result[0].Name, records.Result[0].Content, ip)

}
