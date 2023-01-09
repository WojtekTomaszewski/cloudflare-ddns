# cloudflare-ddns updates Cloudflare DNS address record to detected public ip value

## Usage 

### Build

(Optional) update of GOOS and GOARCH in Makefile

```bash
make build
```

### Provide config

Provide `config.cfg` file in current dir or `$HOME/.config/cloudflare-ddns`

```bash
TOKEN="123"
ZONE="exmaple.com"
TYPE="A"
SUBDOMAIN=""
```

You can override config options with env variables prefixed with `CLOUDFLARE_`

### Run

```bash
bin/cloudflare-ddns
```




