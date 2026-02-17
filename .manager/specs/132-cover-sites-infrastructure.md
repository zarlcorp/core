# 132: Cover Sites Infrastructure

## Objective
Containerize all 10 cover sites into a single nginx Docker image, create k8s manifests for the existing k3s cluster, write a deploy script, and set DNS A records via Namecheap API.

## Context
Specs 130 and 131 created 10 static sites under `sites/<domain>/`. Spec 133 redesigned them. This spec packages them into a Docker container and deploys to the k3s cluster at 161.35.175.61.

The cluster runs Traefik as ingress controller with a `letsencrypt` certResolver (HTTP-01 challenge). Existing services use Traefik IngressRoute CRDs.

Deploy once and forget — no CI/CD pipeline.

## Requirements

### Nginx configuration
Create `nginx.conf` that serves each site as a virtual host:

```nginx
server {
    listen 80;
    server_name moxmail.site www.moxmail.site;
    root /usr/share/nginx/html/moxmail.site;
    index index.html;
    location / {
        try_files $uri $uri/ /index.html;
    }
}
# ... repeat for all 10 domains
```

All 10 domains:
- moxmail.site
- jotmail.xyz
- snapmail.icu
- fogmail.space
- nordvik.work
- hawkpoint.live
- axontel.network
- clearwire.digital
- nexacom.online
- calvera.run

Also include a default server block that returns 444 for unmatched hosts.

### Dockerfile
```dockerfile
FROM nginx:alpine
COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY sites/ /usr/share/nginx/html/
```

### Kubernetes manifests
Create `k8s/` directory with:

**namespace.yaml** — `cover-sites` namespace

**deployment.yaml** — single replica nginx pod
- Image: `cover-sites:latest`
- Port 80
- Resource limits: 64Mi memory, 100m CPU
- Liveness probe: HTTP GET / on port 80

**service.yaml** — ClusterIP service on port 80

**ingress.yaml** — Traefik IngressRoute CRDs for all 10 domains:
- Each domain gets an IngressRoute with `Host()` match rule
- Both `domain.tld` and `www.domain.tld` in the match (OR'd)
- `entryPoints: [websecure]`
- `tls.certResolver: letsencrypt`
- One IngressRoute per domain (10 total) in a single YAML file separated by `---`

Also create a **redirect.yaml** — Traefik middleware + IngressRoute for HTTP→HTTPS redirect:
- Middleware: `redirectScheme` with `scheme: https` and `permanent: true`
- IngressRoute on `web` entrypoint matching all 10 domains, applying the redirect middleware

Reference the existing pattern from the cluster:
```yaml
apiVersion: traefik.io/v1alpha1
kind: IngressRoute
metadata:
  name: cover-sites-nordvik
  namespace: cover-sites
spec:
  entryPoints:
    - websecure
  routes:
    - kind: Rule
      match: Host(`nordvik.work`) || Host(`www.nordvik.work`)
      services:
        - name: cover-sites
          port: 80
  tls:
    certResolver: letsencrypt
```

### Deploy script
Create `deploy.sh` that:
1. `docker build -t cover-sites:latest .`
2. `docker save cover-sites:latest -o cover-sites.tar`
3. `scp cover-sites.tar root@161.35.175.61:/tmp/`
4. `ssh root@161.35.175.61 'k3s ctr images import /tmp/cover-sites.tar && rm /tmp/cover-sites.tar'`
5. `kubectl apply -f k8s/`
6. `kubectl -n cover-sites rollout restart deployment/cover-sites`

### DNS configuration
Create `dns-setup.sh` that sets A records for all 10 domains via Namecheap API.

Namecheap API details:
- Endpoint: `https://api.namecheap.com/xml.response`
- ApiUser: `8run0`
- ApiKey: `dd7dc79b749c4ebebd4c07884b2e119d`
- ClientIp: `81.187.253.175`
- Command: `namecheap.domains.dns.setHosts`

For each domain, set:
- `@` A record → 161.35.175.61
- `www` CNAME → `@` (or A record → 161.35.175.61)

**CRITICAL**: When setting DNS records via `namecheap.domains.dns.setHosts`, you MUST preserve existing MX records for email forwarding. Always GET existing records first (`namecheap.domains.dns.getHosts`), then SET with all records (existing MX + new A/CNAME records).

The script should:
1. For each domain, GET existing host records
2. Parse out any MX records
3. SET new A + www records while preserving MX records
4. Print confirmation for each domain

### What NOT to do
- Do not modify any site content
- Do not set up cert-manager — Traefik handles TLS natively
- Do not create separate pods per domain — single nginx pod serves all
- Do not create CI/CD pipelines — this is deploy once and forget
- Do not use Helm or kustomize — plain manifests only
- Do not touch zarlcorp.com or zarl.dev DNS/infrastructure

## Target Repo
zarlcorp/cover-sites

## Agent Role
backend

## Files to Create
- nginx.conf
- Dockerfile
- k8s/namespace.yaml
- k8s/deployment.yaml
- k8s/service.yaml
- k8s/ingress.yaml
- k8s/redirect.yaml
- deploy.sh
- dns-setup.sh

## Notes
- VPS: 161.35.175.61 (Digital Ocean, k3s)
- Traefik certResolver name: `letsencrypt` (confirmed from cluster)
- Traefik CRD API version: `traefik.io/v1alpha1` (confirmed from cluster)
- Traefik entrypoints: `web` (HTTP) and `websecure` (HTTPS)
- Local container registry exists at registry:5000 but we're using scp instead
- The agent should NOT actually run deploy.sh or dns-setup.sh — just create the scripts
- The agent should NOT run kubectl commands — just create the manifest files
