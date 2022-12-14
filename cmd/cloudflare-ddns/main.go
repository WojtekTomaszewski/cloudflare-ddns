package main

import (
	"log"

	"github.com/WojtekTomaszewski/cloudflare-ddns/internal/cloudflare"
	"github.com/WojtekTomaszewski/cloudflare-ddns/internal/iptools"
	"github.com/spf13/viper"
)

// var (
// 	token      = os.Getenv("CLOUDFLARE_TOKEN")
// 	zone       = os.Getenv("CLOUDFLARE_ZONE")
// 	subdomain  = os.Getenv("CLOUDFLARE_SUBDOMAIN")
// 	recordType = os.Getenv("CLOUDFLARE_RECORD_TYPE")
// )

func init() {
	viper.SetConfigName("config.cfg")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.config/cloudflare-ddns")
	viper.SetDefault("subdomain", "")
	viper.SetDefault("type", "A")
	viper.SetEnvPrefix("cloudflare")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("failed to read config", err)
	}

}

func main() {

	log.Printf("running config: zone: %s, subdomain: %s, record: %s", viper.Get("zone"), viper.Get("subdomain"), viper.Get("type"))

	ip, err := iptools.GetCurrentIP()
	if err != nil {
		log.Fatal("failed to get current ip address, ", err)
	}
	log.Println("detected ip: ", ip)

	if !iptools.IsIPValid(ip) {
		log.Fatal("could not validat current public ip address", ip)
	}

	cf := cloudflare.NewClient(viper.Get("token").(string))

	var zoneID string
	zones, err := cf.GetZones()
	if err != nil {
		log.Fatal("failed to get zones, ", err)
	}
	for _, z := range zones.Result {
		if z.Name == viper.Get("zone").(string) {
			zoneID = z.ID
			break
		}
	}

	recordName := viper.Get("zone").(string)
	if viper.Get("subdomain") != "" {
		recordName = viper.Get("subdomain").(string)
	}

	records, err := cf.GetDNSRecord(zoneID, viper.Get("type").(string), recordName)
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

	if err := cf.UpdateDNSRecord(zoneID, records.Result[0].ID, updatedRecord); err != nil {
		log.Fatal("failed to update dns record", err)
	}

	log.Printf("updated %s record, old ip %s, new ip %s", records.Result[0].Name, records.Result[0].Content, ip)

}
