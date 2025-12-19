# FB-15: mcpui-go Extraction & Documentation

**Priority**: Low
**Complexity**: Low-Medium
**Dependencies**: FB-12 ✅, FB-13 (integration battle-testing)

## Description
Extract `pkg/mcpui/` to standalone repository `github.com/ironystock/mcpui-go` with full documentation and examples following the official MCP Go SDK patterns.

## Current State
- Core SDK complete in `pkg/mcpui/` (84.7% test coverage)
- 12 source files, ~3,000 lines of code
- `example_test.go` with 13 runnable examples
- `doc.go` with package documentation

---

## Documentation Tasks (following go-sdk patterns)

### README.md
- [ ] Package overview and value proposition
- [ ] "Getting Started" with minimal example
- [ ] Package/feature documentation links
- [ ] Installation instructions
- [ ] License and acknowledgements

### docs/ Directory

| File | Description |
|------|-------------|
| `README.md` | Documentation index |
| `content-types.md` | HTMLContent, URLContent, RemoteDOMContent, BlobContent |
| `resources.md` | UIResource, UIResourceContents, validation |
| `actions.md` | UIAction types, payloads, parsing |
| `responses.md` | UIResponse builders, success/error handling |
| `handlers.md` | UIActionHandler, Router, typed wrappers |
| `integration.md` | How to integrate with MCP servers |

### CONTRIBUTING.md
- [ ] Contribution guidelines
- [ ] Code style requirements
- [ ] Testing requirements
- [ ] PR process

---

## Examples Directory (following go-sdk patterns)

```
examples/
├── basic/
│   └── main.go          # Minimal HTML resource example
├── router/
│   └── main.go          # Router with multiple handlers
├── remote-dom/
│   └── main.go          # Remote DOM with React framework
├── action-handling/
│   └── main.go          # Complete action→response flow
├── mcp-integration/
│   └── main.go          # Integration with MCP server
└── README.md            # Examples index
```

### Example: basic/main.go
```go
package main

import (
    "encoding/json"
    "fmt"
    "github.com/ironystock/mcpui-go"
)

func main() {
    // Create HTML content for a greeting card
    content := &mcpui.HTMLContent{
        HTML: `<div style="padding: 20px;">
            <h1>Hello, World!</h1>
            <p>This is rendered in a sandboxed iframe.</p>
        </div>`,
    }

    // Create resource contents
    rc, _ := mcpui.NewUIResourceContents("ui://greeting/hello", content)

    data, _ := json.MarshalIndent(rc, "", "  ")
    fmt.Println(string(data))
}
```

### Example: router/main.go
```go
package main

import (
    "context"
    "fmt"
    "github.com/ironystock/mcpui-go"
)

func main() {
    router := mcpui.NewRouter()

    // Handle tool actions
    router.HandleType(mcpui.ActionTypeTool, mcpui.WrapToolHandler(
        func(ctx context.Context, toolName string, params map[string]any) (any, error) {
            return map[string]string{"executed": toolName}, nil
        },
    ))

    // Handle specific resource
    router.HandleResource("ui://dashboard/main", func(ctx context.Context, req *mcpui.UIActionRequest) (*mcpui.UIActionResult, error) {
        return &mcpui.UIActionResult{Response: "dashboard handled"}, nil
    })

    fmt.Println("Router configured with", len(router.typeHandlers), "type handlers")
}
```

---

## Extraction Checklist

### Phase 1: Repository Setup
- [ ] Create `github.com/ironystock/mcpui-go` repository
- [ ] Initialize go.mod with `module github.com/ironystock/mcpui-go`
- [ ] Copy source files from `pkg/mcpui/`
- [ ] Update import paths in examples
- [ ] Add LICENSE (MIT)
- [ ] Configure GitHub Actions for CI

### Phase 2: Documentation
- [ ] Write README.md with getting started
- [ ] Create docs/ directory structure
- [ ] Write content-types.md
- [ ] Write resources.md
- [ ] Write actions.md
- [ ] Write responses.md
- [ ] Write handlers.md
- [ ] Write integration.md
- [ ] Write CONTRIBUTING.md

### Phase 3: Examples
- [ ] Create examples/ directory
- [ ] Write basic/ example
- [ ] Write router/ example
- [ ] Write remote-dom/ example
- [ ] Write action-handling/ example
- [ ] Write mcp-integration/ example
- [ ] Write examples/README.md index

### Phase 4: Polish
- [ ] Ensure all public APIs have godoc
- [ ] Add badges (Go Reference, CI status, coverage)
- [ ] Create GitHub release v0.1.0
- [ ] Submit to awesome-mcp list (if applicable)
- [ ] Update agentic-obs to use external module

---

## Timeline Recommendation
- Wait until FB-13 (MCP-UI Integration) is complete
- Real-world usage in agentic-obs will identify API improvements
- Extract after APIs are stable

## Reference: go-sdk Structure
The official MCP Go SDK has:
- `README.md` - Overview, getting started
- `docs/` - client.md, server.md, protocol.md, troubleshooting.md
- `examples/` - 14 server examples (basic, completion, hello, memory, etc.)
- `CONTRIBUTING.md` - Contribution guidelines
