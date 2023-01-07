# cloudflare-ddns

## Simple app to set DNS record to your public ip

### How to use 

Build

```bash
make build
```

Config is taken from env variables so make sure those are available:

```
TOKEN is a valid Cloudflare access token  
ZONE is name of Cloudflare zone to modify - usualy this is your top domain name  
(optional) SUBDOMAIN - provide if you want to update subdomain record  
TYPE - type of the record to update, most likely 'A'  
```

Run

```bash
TOKEN=123 ZONE=example.com TYPE=A ./cloudflare-ddns
```




