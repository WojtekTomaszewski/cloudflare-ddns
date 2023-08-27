# cloudflare-ddns updates Cloudflare DNS A record to specified or detected public ip

## Usage 

### Build binary or image

```bash
make build

or

make build-image
```

### Usage

Set `CLOUDFLARE_TOKEN` environment variable with api key having access to the zone you want to modify

Then you can invoke:

```bash
# Update zone example.com, domain example.com A record to ip x.x.x.x
bin/cloudflare-ddns --zone example.com --ip x.x.x.x

# Update zone exmaple.com, domain sub.exmaple.com A record to ip value detected with http://ifconfig.me
bin/cloudflare-ddns --zone example.com --domain sub.example.com

# Update zone example.com, domain exmaple.com A record to ip value detected with http://ifconfig.me but run in daemon mode and do update every 4h
bin/cloudflare-ddns --zone example.com --daemon --interval 4
```




