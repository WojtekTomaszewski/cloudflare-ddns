package main

import (
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/WojtekTomaszewski/cloudflare-ddns/internal/cloudflare"
	"github.com/WojtekTomaszewski/cloudflare-ddns/internal/iptools"
)

// dynDNS represent DNS record updater
type dynDNS struct {
	daemon   bool
	domain   string
	interval int64
	ip       string
	log      *slog.Logger
	token    string
	zone     string
}

// newDynDNS return minimal dynDNS instance with just logger
func newDynDNS() *dynDNS {
	d := &dynDNS{
		log: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}

	d.token = os.Getenv("CLOUDFLARE_TOKEN")
	if d.token == "" {
		d.log.Error("CLOUDFLARE_TOKEN environment variable not set")
		os.Exit(1)
	}

	flag.StringVar(&d.zone, "zone", "", "set zone name where A record resides")
	flag.StringVar(&d.domain, "domain", "", "set domain to update if name different from zone name")
	flag.StringVar(&d.ip, "ip", "", "set A record to this IP value instead of trying to detect public ip")
	flag.BoolVar(&d.daemon, "daemon", false, "run in daemon mode")
	flag.Int64Var(&d.interval, "interval", 12, "interval in hours between record update in daemon mode")
	flag.Parse()
	if d.zone == "" {
		d.log.Error("setting zone flag is mandatory")
		os.Exit(1)
	}
	if d.domain == "" {
		d.domain = d.zone
	}

	if d.ip != "" && !iptools.IsIPValid(d.ip) {
		d.log.Error("provided ip is not valid")
		os.Exit(1)
	}

	d.log = d.log.With("zone", d.zone, "domain", d.domain)

	return d
}

// start runs record A update on specified zone and domain
func (d *dynDNS) start() {
	if d.daemon {
		duration := time.Duration(d.interval)
		ticker := time.NewTicker(duration * time.Hour)
		sigs := make(chan os.Signal, 1)
		done := make(chan bool, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			for {
				select {
				case <-ticker.C:
					d.updateRecord()
				case <-sigs:
					d.log.Info("shutting down")
					ticker.Stop()
					done <- true
				}
			}
		}()
		d.log.Info("started in daemon mode", "interval", fmt.Sprintf("%dh", d.interval))
		<-done

	} else {
		d.updateRecord()
	}
}

// updateRecord updates A record for specified zone and domain
func (d *dynDNS) updateRecord() {
	cf := cloudflare.NewClient(d.token)

	var zoneID string
	zones, err := cf.GetZones()
	if err != nil {
		d.log.Error("failed to fetch DNS zones", "err", err.Error())
		return
	}

	for _, z := range zones.Result {
		if z.Name == d.zone {
			zoneID = z.ID
			break
		}
	}

	records, err := cf.GetDNSRecord(zoneID, "A", d.domain)
	if err != nil {
		d.log.Error("failed to get DNS A record(s)", "err", err.Error())
		return
	}
	if len(records.Result) == 0 {
		d.log.Warn("no matching DNS A record defined")
		return
	}

	if d.ip == "" {
		if err := d.getIP(); err != nil {
			d.log.Info("failed to detect public ip", "err", err.Error())
			return
		}

	}

	if records.Result[0].Content == d.ip {
		d.log.Info("detected IP matches already defined in A record")
		return
	}

	updatedRecord := &cloudflare.Record{
		Name:    records.Result[0].Name,
		Type:    records.Result[0].Type,
		TTL:     records.Result[0].TTL,
		Proxied: records.Result[0].Proxied,
		Content: d.ip,
	}

	if err := cf.UpdateDNSRecord(zoneID, records.Result[0].ID, updatedRecord); err != nil {
		d.log.Error("failed to update dns record", "err", err.Error())
		return
	}

	d.log.Info("record updated", "ip", d.ip)
}

func (d *dynDNS) getIP() error {
	ip, err := iptools.GetCurrentIP()
	if err != nil {
		return err
	}

	if !iptools.IsIPValid(ip) {
		return errors.New("detected ip address is not valid")
	}

	d.ip = ip
	return nil
}
