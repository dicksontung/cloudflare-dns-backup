# cloudflare-dns-backup
```
An application to back up cloudflare dns record. Configure using ENVIRONMENT_VARIABLES or yaml file.

Usage:
  cloudflare-dns-backup [flags]

Flags:
      --config string   config file (default "./cloudflare-dns-backup.yaml")
  -h, --help            help for cloudflare-dns-backup
  -p, --prefix string   prefix for the dns record file name (default "dns_record_")
  -t, --token string    bearer token to connect to Cloudflare account
  -z, --zone strings    id of the zone to backup
```
