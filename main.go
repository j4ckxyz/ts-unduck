package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"tailscale.com/tsnet"
)

var (
	hostname = flag.String("hostname", "unduck", "hostname to use on the tailnet")
	addr     = flag.String("addr", ":80", "address to listen on")
)

const defaultBangTag = "g"

var bangRegex = regexp.MustCompile(`!(\S+)`)

func findBang(tag string) *Bang {
	tag = strings.ToLower(tag)
	for i := range bangs {
		if strings.ToLower(bangs[i].T) == tag {
			return &bangs[i]
		}
	}
	return nil
}

func getBangRedirectURL(query string) string {
	query = strings.TrimSpace(query)
	if query == "" {
		return ""
	}

	match := bangRegex.FindStringSubmatch(query)

	var selectedBang *Bang
	if match != nil && len(match) > 1 {
		bangCandidate := strings.ToLower(match[1])
		selectedBang = findBang(bangCandidate)
	}

	if selectedBang == nil {
		selectedBang = findBang(defaultBangTag)
	}

	if selectedBang == nil {
		return ""
	}

	cleanQuery := bangRegex.ReplaceAllString(query, "")
	cleanQuery = strings.TrimSpace(cleanQuery)

	if cleanQuery == "" {
		return "https://" + selectedBang.D
	}

	searchURL := strings.Replace(
		selectedBang.U,
		"{{{s}}}",
		url.QueryEscape(cleanQuery),
		-1,
	)
	searchURL = strings.ReplaceAll(searchURL, "%2F", "/")

	return searchURL
}

const landingPageHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Unduck - Tailscale Edition</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: #333;
            min-height: 100vh;
            display: flex;
            flex-direction: column;
            align-items: center;
            justify-content: center;
        }
        .content-container {
            background: white;
            padding: 3rem;
            border-radius: 1rem;
            box-shadow: 0 20px 60px rgba(0,0,0,0.3);
            max-width: 600px;
            width: 90%;
        }
        h1 {
            font-size: 3rem;
            margin-bottom: 1rem;
            color: #667eea;
        }
        p {
            font-size: 1.1rem;
            line-height: 1.6;
            margin-bottom: 1.5rem;
            color: #555;
        }
        a {
            color: #667eea;
            text-decoration: none;
        }
        a:hover {
            text-decoration: underline;
        }
        .url-container {
            display: flex;
            gap: 0.5rem;
            margin-top: 1.5rem;
        }
        .url-input {
            flex: 1;
            padding: 0.75rem;
            border: 2px solid #ddd;
            border-radius: 0.5rem;
            font-size: 1rem;
            font-family: monospace;
        }
        .copy-button {
            padding: 0.75rem 1.5rem;
            background: #667eea;
            color: white;
            border: none;
            border-radius: 0.5rem;
            cursor: pointer;
            font-size: 1rem;
            transition: background 0.2s;
        }
        .copy-button:hover {
            background: #5568d3;
        }
        .footer {
            margin-top: 2rem;
            color: rgba(255,255,255,0.9);
            font-size: 0.9rem;
        }
        .footer a {
            color: white;
            margin: 0 0.5rem;
        }
    </style>
</head>
<body>
    <div class="content-container">
        <h1>Und*ck</h1>
        <p>DuckDuckGo's bang redirects are too slow. This Tailscale version runs on your private tailnet. Add the following URL as a custom search engine to your browser:</p>
        <div class="url-container">
            <input type="text" class="url-input" value="http://unduck/?q=%s" readonly />
            <button class="copy-button" onclick="copyURL()">Copy</button>
        </div>
        <p style="margin-top: 1.5rem; font-size: 0.95rem;">
            Enables <a href="https://duckduckgo.com/bang.html" target="_blank">all of DuckDuckGo's bangs</a>, 
            but faster and running on your own tailnet.
        </p>
    </div>
    <div class="footer">
        <a href="https://t3.chat" target="_blank">t3.chat</a>
        •
        <a href="https://x.com/theo" target="_blank">theo</a>
        •
        <a href="https://github.com/t3dotgg/unduck" target="_blank">github</a>
    </div>
    <script>
        function copyURL() {
            const input = document.querySelector('.url-input');
            input.select();
            document.execCommand('copy');
            const button = document.querySelector('.copy-button');
            button.textContent = 'Copied!';
            setTimeout(() => {
                button.textContent = 'Copy';
            }, 2000);
        }
    </script>
</body>
</html>`

func handleRequest(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")

	if query == "" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, landingPageHTML)
		return
	}

	redirectURL := getBangRedirectURL(query)
	if redirectURL == "" {
		http.Error(w, "Could not find redirect URL", http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func main() {
	flag.Parse()

	srv := &tsnet.Server{
		Hostname: *hostname,
		Logf:     log.Printf,
	}
	defer srv.Close()

	ln, err := srv.Listen("tcp", *addr)
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	log.Printf("Unduck server starting on tailnet as %s", *hostname)
	log.Printf("Listening on %s", *addr)
	log.Printf("Access via: http://%s/?q=%%s", *hostname)

	http.HandleFunc("/", handleRequest)
	if err := http.Serve(ln, nil); err != nil {
		log.Fatal(err)
	}
}
