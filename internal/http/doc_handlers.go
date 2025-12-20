package http

import (
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/ironystock/agentic-obs/internal/docs"
)

// docPageTemplate is the template for rendering documentation pages.
var docPageTemplate = template.Must(template.New("doc").Parse(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}} - agentic-obs</title>
    <link rel="icon" type="image/svg+xml" href="/favicon.svg">
    <style>
        :root {
            --bg-primary: #1a1a2e;
            --bg-secondary: #16213e;
            --bg-card: #0f3460;
            --text-primary: #eaeaea;
            --text-secondary: #a0a0a0;
            --accent: #e94560;
            --accent-hover: #ff6b6b;
            --success: #4ecca3;
            --border: #2a2a4a;
        }

        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background: var(--bg-primary);
            color: var(--text-primary);
            line-height: 1.6;
        }

        .container {
            max-width: 900px;
            margin: 0 auto;
            padding: 20px;
        }

        header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            padding: 20px 0;
            border-bottom: 1px solid var(--border);
            margin-bottom: 30px;
        }

        header h1 {
            font-size: 1.5rem;
            font-weight: 600;
        }

        header h1 span {
            color: var(--accent);
        }

        nav {
            display: flex;
            gap: 16px;
        }

        nav a {
            color: var(--text-secondary);
            text-decoration: none;
            font-size: 0.9rem;
            padding: 8px 16px;
            border-radius: 6px;
            transition: all 0.2s;
        }

        nav a:hover {
            color: var(--accent);
            background: rgba(233, 69, 96, 0.1);
        }

        nav a.active {
            color: var(--accent);
            background: rgba(233, 69, 96, 0.15);
        }

        .doc-list {
            list-style: none;
        }

        .doc-list li {
            margin-bottom: 16px;
        }

        .doc-list a {
            display: block;
            padding: 20px;
            background: var(--bg-secondary);
            border: 1px solid var(--border);
            border-radius: 8px;
            text-decoration: none;
            transition: all 0.2s;
        }

        .doc-list a:hover {
            border-color: var(--accent);
            transform: translateY(-2px);
        }

        .doc-list .title {
            color: var(--accent);
            font-size: 1.1rem;
            font-weight: 600;
            margin-bottom: 4px;
        }

        .doc-list .description {
            color: var(--text-secondary);
            font-size: 0.9rem;
        }

        .doc-content {
            background: var(--bg-secondary);
            border: 1px solid var(--border);
            border-radius: 12px;
            padding: 32px;
        }

        /* Markdown content styling */
        .doc-content h1 { font-size: 2rem; margin-bottom: 1rem; color: var(--accent); }
        .doc-content h2 { font-size: 1.5rem; margin: 2rem 0 1rem; padding-bottom: 0.5rem; border-bottom: 1px solid var(--border); }
        .doc-content h3 { font-size: 1.25rem; margin: 1.5rem 0 0.75rem; }
        .doc-content h4 { font-size: 1.1rem; margin: 1rem 0 0.5rem; color: var(--text-secondary); }

        .doc-content p { margin-bottom: 1rem; }

        .doc-content a { color: var(--accent); text-decoration: none; }
        .doc-content a:hover { text-decoration: underline; }

        .doc-content code {
            background: var(--bg-card);
            padding: 2px 6px;
            border-radius: 4px;
            font-family: 'SF Mono', 'Fira Code', Consolas, monospace;
            font-size: 0.9em;
        }

        .doc-content pre {
            background: var(--bg-card);
            padding: 16px;
            border-radius: 8px;
            overflow-x: auto;
            margin-bottom: 1rem;
        }

        .doc-content pre code {
            background: none;
            padding: 0;
        }

        .doc-content ul, .doc-content ol {
            margin-left: 1.5rem;
            margin-bottom: 1rem;
        }

        .doc-content li { margin-bottom: 0.5rem; }

        .doc-content blockquote {
            border-left: 4px solid var(--accent);
            padding-left: 16px;
            margin: 1rem 0;
            color: var(--text-secondary);
        }

        .doc-content table {
            width: 100%;
            border-collapse: collapse;
            margin-bottom: 1rem;
        }

        .doc-content th, .doc-content td {
            padding: 12px;
            text-align: left;
            border-bottom: 1px solid var(--border);
        }

        .doc-content th {
            background: var(--bg-card);
            font-weight: 600;
        }

        .doc-content tr:hover td {
            background: rgba(233, 69, 96, 0.05);
        }

        .doc-content img {
            max-width: 100%;
            height: auto;
            border-radius: 8px;
        }

        .doc-content hr {
            border: none;
            border-top: 1px solid var(--border);
            margin: 2rem 0;
        }

        footer {
            text-align: center;
            padding: 20px;
            color: var(--text-secondary);
            font-size: 0.875rem;
            margin-top: 40px;
        }

        footer a {
            color: var(--accent);
            text-decoration: none;
        }
    </style>
</head>
<body>
    <div class="container">
        <header>
            <h1><span>agentic</span>-obs</h1>
            <nav>
                <a href="/">Dashboard</a>
                <a href="/docs/" class="active">Documentation</a>
            </nav>
        </header>

        {{if .IsList}}
        <h2 style="margin-bottom: 20px;">Documentation</h2>
        <ul class="doc-list">
            {{range .Docs}}
            <li>
                <a href="/docs/{{.Name}}">
                    <div class="title">{{.Title}}</div>
                    <div class="description">{{.Description}}</div>
                </a>
            </li>
            {{end}}
        </ul>
        {{else}}
        <div class="doc-content">
            {{.Content}}
        </div>
        {{end}}

        <footer>
            <a href="/docs/">&larr; Back to Documentation</a>
        </footer>
    </div>
</body>
</html>`))

// handleDocsIndex serves the documentation index page.
func (s *Server) handleDocsIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	docList, err := docs.List()
	if err != nil {
		log.Printf("Failed to list docs: %v", err)
		http.Error(w, "Failed to list documentation", http.StatusInternalServerError)
		return
	}

	data := struct {
		Title   string
		IsList  bool
		Docs    []docs.Doc
		Content template.HTML
	}{
		Title:  "Documentation",
		IsList: true,
		Docs:   docList,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := docPageTemplate.Execute(w, data); err != nil {
		log.Printf("Failed to render docs index: %v", err)
	}
}

// handleDocView serves a specific documentation page.
func (s *Server) handleDocView(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract doc name from path: /docs/{name}
	path := strings.TrimPrefix(r.URL.Path, "/docs/")
	path = strings.TrimSuffix(path, "/")

	if path == "" {
		s.handleDocsIndex(w, r)
		return
	}

	// Validate and get document
	doc, err := docs.Get(path)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Render markdown to HTML
	htmlContent, err := docs.RenderHTML(path)
	if err != nil {
		log.Printf("Failed to render doc %s: %v", path, err)
		http.Error(w, "Failed to render documentation", http.StatusInternalServerError)
		return
	}

	data := struct {
		Title   string
		IsList  bool
		Docs    []docs.Doc
		Content template.HTML
	}{
		Title:   doc.Title,
		IsList:  false,
		Content: template.HTML(htmlContent),
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := docPageTemplate.Execute(w, data); err != nil {
		log.Printf("Failed to render doc page: %v", err)
	}
}

// handleDocsAPI returns the list of available docs as JSON.
func (s *Server) handleDocsAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	docList, err := docs.List()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to list documentation"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"count": len(docList),
		"docs":  docList,
	})
}
