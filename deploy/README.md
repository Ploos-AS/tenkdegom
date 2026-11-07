# Deployment – pending VPS selection

This is **preparation**, not an enabled deployment workflow.

## Target layout

- Restricted SSH account: `tenkdegom-deploy` (no sudo or root access).
- Root directory: `/srv/tenkdegom` (owned by deploy user).
- Releases: `/srv/tenkdegom/releases/<release-id>/`.
- Active symlink: `/srv/tenkdegom/current`.
- Caddy serves `/srv/tenkdegom/current` as static files, with read/traverse permissions.

Example Caddyfile (adapt to your server):

```caddyfile
tenkdegom.no {
    root * /srv/tenkdegom/current
    encode zstd gzip
    file_server
}
```

## Manual release

Run `npm run check && npm run build` locally. Configure SSH key authentication and host-key verification **before** deploying.

```bash
DEPLOY_HOST=server.example DEPLOY_USER=tenkdegom-deploy DEPLOY_PATH=/srv/tenkdegom bash deploy/deploy.sh
```

Do not place credentials in the repository. Never expose an SSH private key in build logs.

## Rollback

After selecting a previous known-good release name, switch `current` atomically using `ln -sfn` and `mv -Tf` (same pattern as deploy script). Keep prior releases until a retention policy is defined.

## Future GitHub Actions CD

Create a separate, protected production environment, use a dedicated deploy SSH key and pinned host key, require a green CI build, and deploy only from `main`. Choose VPS, DNS, Caddy setup, log retention and access model first. Do not enable automatic deployment before these decisions.
