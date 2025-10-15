# ts-unduck

A Tailscale-based server version of [Unduck](https://github.com/t3dotgg/unduck) that provides fast DuckDuckGo bang redirects on your private tailnet.

## What is this?

DuckDuckGo's bang redirects (!gh, !wiki, etc.) are slow because they happen server-side. The original Unduck solved this with client-side JavaScript. This version runs as a tsnet server on your Tailscale network, giving you:

- **Fast redirects** - Server runs on your tailnet, no external DNS lookups
- **Privacy** - All searches stay within your network
- **All DuckDuckGo bangs** - 13,000+ bang shortcuts supported
- **Easy access** - Use via MagicDNS: `http://unduck/?q=!gh+golang`

## Prerequisites

- [Tailscale](https://tailscale.com/) account and network
- Linux server (works on low-end boxes, Raspberry Pi, etc.)

## Quick Deploy (Linux)

One-line install on any Linux machine:

```bash
curl -fsSL https://raw.githubusercontent.com/t3dotgg/ts-unduck/main/quick-deploy.sh | bash
```

This will:
1. Install Go if needed
2. Clone the repo
3. Build the binary
4. Set up systemd service
5. Start the server

On first run, check logs for Tailscale auth URL:
```bash
sudo journalctl -u unduck -f
```

## Manual Install

### Building

```bash
go build -o unduck
```

Or with automatic Go version management:

```bash
GOTOOLCHAIN=auto go build -o unduck
```

### Running

```bash
./unduck
```

The server will:
1. Join your Tailscale network as a device named `unduck`
2. Listen on port 80
3. Be accessible via `http://unduck/` (or `http://unduck.<tailnet-name>.ts.net/`)

### Custom hostname

```bash
./unduck -hostname my-search
```

### Custom port

```bash
./unduck -addr :8080
```

## Low-End Linux Box Tips

This runs great on:
- Raspberry Pi (any model)
- Old laptops
- $5/month VPS
- Home servers

Memory usage: ~20-30MB  
CPU usage: Minimal (only during redirects)

The `install.sh` script handles everything automatically. Just make sure you have:
- Internet connection
- Tailscale account (free tier works great)

## Browser Setup

Add Unduck as a custom search engine in your browser:

**URL:** `http://unduck/?q=%s`

### Chrome/Brave/Edge
1. Settings → Search engine → Manage search engines
2. Add new search engine
3. Use `http://unduck/?q=%s` as the URL

### Firefox
1. Visit `http://unduck/` in your browser
2. Right-click the address bar
3. "Add Unduck" should appear in the context menu

## Usage Examples

- `!gh tailscale/tailscale` → Search GitHub
- `!wiki rust` → Search Wikipedia  
- `!yt golang tutorial` → Search YouTube
- `!ghr t3dotgg/unduck` → Go to GitHub repo
- `hello world` → Falls back to Google search (no bang specified)

See [all available bangs](https://duckduckgo.com/bang.html).

## How it works

1. Parses query parameter from URLs like `?q=!gh+search+term`
2. Matches bang shortcuts (e.g., `!gh`) against the bang database
3. Redirects to the appropriate search URL
4. Falls back to Google (`!g`) if no bang is specified
5. Shows a landing page when accessed without query parameters

## Project Structure

- `main.go` - HTTP server with tsnet integration and redirect logic
- `bangs.go` - Database of 13,000+ bang shortcuts (auto-generated)
- `convert-bangs.js` - Script to convert TypeScript bang data to Go
- `unduck/` - Original client-side Unduck web app

## Regenerating bang database

If you need to update the bangs from the original source:

```bash
node convert-bangs.js
```

This will regenerate `bangs.go` from `unduck/src/bang.ts`.

## Running as a service

### systemd (Linux)

Create `/etc/systemd/system/unduck.service`:

```ini
[Unit]
Description=Unduck Bang Redirect Server
After=network.target

[Service]
Type=simple
User=YOUR_USER
ExecStart=/path/to/unduck
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

Then:

```bash
sudo systemctl enable unduck
sudo systemctl start unduck
```

## Credits

- Original [Unduck](https://github.com/t3dotgg/unduck) by [Theo](https://twitter.com/theo)
- Bang data from [DuckDuckGo](https://duckduckgo.com/bang.js)
- Built with [tsnet](https://pkg.go.dev/tailscale.com/tsnet)

## License

MIT
