# Contact form deployment (Balder)

The Astro site is static. The Go service is **not** deployed by the website GitHub Actions workflow. Do not merge the contact page into production until the service is installed and tested.

1. Install Go on Balder and build: `cd contact-service && go build -trimpath -o /tmp/tenkdegom-contact .`
2. Install binary to `/usr/local/bin/tenkdegom-contact` with root ownership and mode 0755.
3. Create locked service user: `sudo useradd --system --no-create-home --shell /usr/sbin/nologin tenkdegom-contact` (if absent).
4. Create `/etc/tenkdegom/contact.env` owned by root, mode **0600** with:
   ```
   SMTP_HOST=smtp.protonmail.ch
   SMTP_USER=<the Proton address associated with your SMTP token>
   SMTP_TOKEN=<the token; never commit>
   CONTACT_TO=tenkdegom@ploos.no
   ```
   For characters requiring systemd EnvironmentFile quoting, follow systemd.syntax(7).
5. Copy `tenkdegom-contact.service` to `/etc/systemd/system/`, run `sudo systemctl daemon-reload && sudo systemctl enable --now tenkdegom-contact`.
6. Inside the **tenkdegom.no** Caddy site block, add **before** `file_server`:
   ```
   handle /api/contact {
       reverse_proxy 127.0.0.1:8787
   }
   handle {
       root * /srv/tenkdegom/current
       encode zstd gzip
       file_server
   }
   ```
   Replace existing root/encode/file_server directives in that block with the handles above. Do not change Forenio or the www redirect. Validate with `sudo caddy validate --config /etc/caddy/Caddyfile` then reload.
7. Test the endpoint with an invalid form (expect 400), a valid test message (expect 200 and arrival in Proton), and rate limiting (expect 429). Do not expose the service port in UFW.

Security: one global in-memory rate limit of 10 accepted submissions/minute per process, plus a honeypot; no database or persistent request log. This is a starter protection, not a substitute for abuse monitoring. SMTP errors are logged without message content. The sender's email is included in Reply-To and the message body; it is processed by Proton for delivery. Configure retention policy for the mailbox and update the privacy notice before production.

Do not share SMTP tokens, received messages or IP addresses in tickets or chat.
