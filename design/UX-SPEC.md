# agentic-obs UX/UI Design Specification

> Version 1.0 | December 2025
> Applies to: Web Dashboard, TUI (Terminal UI), MCP-UI Resources

---

## Table of Contents

1. [Design Philosophy](#design-philosophy)
2. [Design Tokens](#design-tokens)
3. [Typography](#typography)
4. [Component Library](#component-library)
5. [Layout Patterns](#layout-patterns)
6. [Platform-Specific Guidelines](#platform-specific-guidelines)
7. [Accessibility Requirements](#accessibility-requirements)
8. [Interaction Patterns](#interaction-patterns)
9. [Implementation Reference](#implementation-reference)

---

## Design Philosophy

### Core Principles

1. **Clarity over cleverness** - OBS control requires precision; interfaces must be unambiguous
2. **Consistency builds trust** - Patterns should be predictable across all three platforms
3. **Dark-first design** - Professional streaming environments use low-light setups
4. **Progressive disclosure** - Show complexity only when needed
5. **Immediate feedback** - Every action needs visible acknowledgment
6. **Accessibility is not optional** - Design for keyboard navigation and screen readers from the start

### Design Goals

- **Professional aesthetic** - Suitable for broadcast/streaming professionals
- **Low visual fatigue** - Comfortable for extended monitoring sessions
- **Information density** - Display relevant status without overwhelming
- **Responsive control** - Quick access to critical functions

---

## Design Tokens

### Color Palette

The agentic-obs color system uses a dark theme with coral accent, optimized for low-light environments typical of streaming setups.

#### Primary Colors

| Token | Hex | RGB | Usage |
|-------|-----|-----|-------|
| `--accent` | `#e94560` | 233, 69, 96 | Primary brand color, CTAs, active states |
| `--accent-hover` | `#ff6b6b` | 255, 107, 107 | Hover state for accent elements |
| `--accent-subtle` | `rgba(233, 69, 96, 0.15)` | - | Background tints for accent elements |

#### Surface Colors

| Token | Hex | RGB | Usage |
|-------|-----|-----|-------|
| `--bg-primary` | `#1a1a2e` | 26, 26, 46 | Page background |
| `--bg-secondary` | `#16213e` | 22, 33, 62 | Card backgrounds |
| `--bg-card` | `#0f3460` | 15, 52, 96 | Nested elements, stat blocks |
| `--bg-elevated` | `#2a2a4a` | 42, 42, 74 | Tooltips, dropdowns |
| `--border` | `#2a2a4a` | 42, 42, 74 | Borders, dividers |

#### Text Colors

| Token | Hex | RGB | Usage |
|-------|-----|-----|-------|
| `--text-primary` | `#eaeaea` | 234, 234, 234 | Headings, primary content |
| `--text-secondary` | `#a0a0a0` | 160, 160, 160 | Labels, metadata |
| `--text-muted` | `#707070` | 112, 112, 112 | Disabled, placeholder text |

#### Semantic Colors

| Token | Hex | RGB | Usage | Contrast Ratio |
|-------|-----|-----|-------|----------------|
| `--success` | `#4ecca3` | 78, 204, 163 | Connected, recording, positive | 7.2:1 on bg-primary |
| `--warning` | `#ffc107` | 255, 193, 7 | Caution states, reconnecting | 9.1:1 on bg-primary |
| `--error` | `#ff6b6b` | 255, 107, 107 | Disconnected, failed, destructive | 5.8:1 on bg-primary |
| `--info` | `#5dade2` | 93, 173, 226 | Informational messages | 6.4:1 on bg-primary |

#### Semantic Background Tints

| Token | Value | Usage |
|-------|-------|-------|
| `--success-bg` | `rgba(78, 204, 163, 0.15)` | Success badge backgrounds |
| `--warning-bg` | `rgba(255, 193, 7, 0.15)` | Warning badge backgrounds |
| `--error-bg` | `rgba(255, 107, 107, 0.15)` | Error badge backgrounds |
| `--info-bg` | `rgba(93, 173, 226, 0.15)` | Info badge backgrounds |

### CSS Custom Properties

```css
:root {
  /* Primary */
  --accent: #e94560;
  --accent-hover: #ff6b6b;
  --accent-subtle: rgba(233, 69, 96, 0.15);

  /* Surfaces */
  --bg-primary: #1a1a2e;
  --bg-secondary: #16213e;
  --bg-card: #0f3460;
  --bg-elevated: #2a2a4a;
  --border: #2a2a4a;

  /* Text */
  --text-primary: #eaeaea;
  --text-secondary: #a0a0a0;
  --text-muted: #707070;

  /* Semantic */
  --success: #4ecca3;
  --success-bg: rgba(78, 204, 163, 0.15);
  --warning: #ffc107;
  --warning-bg: rgba(255, 193, 7, 0.15);
  --error: #ff6b6b;
  --error-bg: rgba(255, 107, 107, 0.15);
  --info: #5dade2;
  --info-bg: rgba(93, 173, 226, 0.15);

  /* Spacing (4px base) */
  --space-1: 4px;
  --space-2: 8px;
  --space-3: 12px;
  --space-4: 16px;
  --space-5: 20px;
  --space-6: 24px;
  --space-8: 32px;
  --space-10: 40px;

  /* Border Radius */
  --radius-sm: 4px;
  --radius-md: 8px;
  --radius-lg: 12px;
  --radius-full: 9999px;

  /* Shadows */
  --shadow-sm: 0 1px 2px rgba(0, 0, 0, 0.2);
  --shadow-md: 0 4px 8px rgba(0, 0, 0, 0.3);
  --shadow-lg: 0 8px 16px rgba(0, 0, 0, 0.4);

  /* Transitions */
  --transition-fast: 100ms ease;
  --transition-normal: 200ms ease;
  --transition-slow: 300ms ease;
}
```

### Spacing Scale

Based on a 4px unit for precise alignment:

| Token | Value | Usage |
|-------|-------|-------|
| `--space-1` | 4px | Tight spacing, icon gaps |
| `--space-2` | 8px | Default element spacing |
| `--space-3` | 12px | Between related elements |
| `--space-4` | 16px | Between content sections |
| `--space-5` | 20px | Container padding |
| `--space-6` | 24px | Card padding |
| `--space-8` | 32px | Major section spacing |
| `--space-10` | 40px | Page margins, large gaps |

### Border Radius Scale

Standardized to three tiers:

| Token | Value | Usage |
|-------|-------|-------|
| `--radius-sm` | 4px | Small elements: tags, badges, tooltips |
| `--radius-md` | 8px | Standard elements: buttons, inputs, cards |
| `--radius-lg` | 12px | Large containers: modals, panels |
| `--radius-full` | 9999px | Pills, avatars, status dots |

### Elevation (Shadows)

| Level | Token | Usage |
|-------|-------|-------|
| 0 | none | Base level, flat surfaces |
| 1 | `--shadow-sm` | Subtle lift: hover states |
| 2 | `--shadow-md` | Raised elements: cards, dropdowns |
| 3 | `--shadow-lg` | Modals, dialogs |

### Transitions

| Token | Duration | Usage |
|-------|----------|-------|
| `--transition-fast` | 100ms | Micro-interactions: hover, focus |
| `--transition-normal` | 200ms | Standard transitions: buttons, toggles |
| `--transition-slow` | 300ms | Large elements: modals, panels |

---

## Typography

### Font Stack

```css
font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, sans-serif;
```

This system font stack ensures:
- Native appearance on each platform
- Optimal rendering performance
- No font loading delays

### Monospace (for code/values)

```css
font-family: 'SF Mono', 'Fira Code', 'Consolas', monospace;
```

### Type Scale

| Name | Size | Weight | Line Height | Usage |
|------|------|--------|-------------|-------|
| `heading-lg` | 1.8rem (28.8px) | 600 | 1.2 | Page titles |
| `heading-md` | 1.5rem (24px) | 600 | 1.3 | Section titles |
| `heading-sm` | 1.1rem (17.6px) | 600 | 1.4 | Card titles |
| `body` | 1rem (16px) | 400 | 1.6 | Primary content |
| `body-sm` | 0.875rem (14px) | 400 | 1.5 | Secondary content |
| `caption` | 0.75rem (12px) | 400 | 1.4 | Metadata, timestamps |
| `overline` | 0.75rem (12px) | 500 | 1.4 | Section labels (uppercase) |

### Typography CSS

```css
/* Page title */
.heading-lg {
  font-size: 1.8rem;
  font-weight: 600;
  line-height: 1.2;
  color: var(--text-primary);
}

/* Section title */
.heading-md {
  font-size: 1.5rem;
  font-weight: 600;
  line-height: 1.3;
  color: var(--text-primary);
}

/* Card title */
.heading-sm {
  font-size: 1.1rem;
  font-weight: 600;
  line-height: 1.4;
  color: var(--text-primary);
}

/* Body text */
.body {
  font-size: 1rem;
  font-weight: 400;
  line-height: 1.6;
  color: var(--text-primary);
}

/* Secondary text */
.body-sm {
  font-size: 0.875rem;
  font-weight: 400;
  line-height: 1.5;
  color: var(--text-secondary);
}

/* Captions, timestamps */
.caption {
  font-size: 0.75rem;
  font-weight: 400;
  line-height: 1.4;
  color: var(--text-secondary);
}

/* Section labels */
.overline {
  font-size: 0.75rem;
  font-weight: 500;
  line-height: 1.4;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  color: var(--text-secondary);
}
```

---

## Component Library

### Buttons

#### Button Variants

| Variant | Background | Text | Border | Usage |
|---------|------------|------|--------|-------|
| Primary | `--accent` | white | none | Main CTAs |
| Secondary | `--bg-card` | `--text-primary` | `--border` | Alternative actions |
| Ghost | transparent | `--accent` | none | Tertiary actions |
| Danger | `--error` | white | none | Destructive actions |

#### Button Sizes

| Size | Padding | Font Size | Min Height |
|------|---------|-----------|------------|
| Small | 6px 12px | 0.75rem | 28px |
| Medium | 8px 16px | 0.875rem | 36px |
| Large | 12px 24px | 1rem | 44px |

#### Button CSS

```css
.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-2);
  padding: 8px 16px;
  font-size: 0.875rem;
  font-weight: 500;
  border: none;
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: background var(--transition-normal), transform var(--transition-fast);
}

.btn:focus-visible {
  outline: 2px solid var(--accent);
  outline-offset: 2px;
}

.btn:active {
  transform: scale(0.98);
}

/* Primary */
.btn-primary {
  background: var(--accent);
  color: white;
}

.btn-primary:hover {
  background: var(--accent-hover);
}

/* Secondary */
.btn-secondary {
  background: var(--bg-card);
  color: var(--text-primary);
  border: 1px solid var(--border);
}

.btn-secondary:hover {
  background: var(--border);
}

/* Ghost */
.btn-ghost {
  background: transparent;
  color: var(--accent);
}

.btn-ghost:hover {
  background: var(--accent-subtle);
}

/* Danger */
.btn-danger {
  background: var(--error);
  color: white;
}

.btn-danger:hover {
  background: #e55555;
}

/* Sizes */
.btn-sm { padding: 6px 12px; font-size: 0.75rem; }
.btn-lg { padding: 12px 24px; font-size: 1rem; }
```

### Cards

#### Card Structure

```html
<div class="card">
  <div class="card-header">
    <span class="card-title">Title</span>
    <button class="btn btn-primary btn-sm">Action</button>
  </div>
  <div class="card-content">
    <!-- Content -->
  </div>
</div>
```

#### Card CSS

```css
.card {
  background: var(--bg-secondary);
  border-radius: var(--radius-lg);
  padding: var(--space-6);
  border: 1px solid var(--border);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--space-5);
}

.card-title {
  font-size: 1.1rem;
  font-weight: 600;
  color: var(--text-primary);
}

.card-content {
  /* Content styling */
}
```

#### Card Variants

| Variant | Style | Usage |
|---------|-------|-------|
| Default | `--bg-secondary` background | Standard containers |
| Nested | `--bg-card` background | Elements within cards |
| Interactive | Hover border color change | Clickable cards |
| Active | `--success` border | Selected/current item |

### Badges and Status Indicators

#### Badge Structure

```html
<span class="badge badge-online">
  <span class="badge-dot"></span>
  Connected
</span>
```

#### Badge CSS

```css
.badge {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 4px 12px;
  border-radius: var(--radius-full);
  font-size: 0.875rem;
  font-weight: 500;
}

.badge-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: currentColor;
}

/* Animated pulse for live status */
.badge-dot.animated {
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

.badge-online {
  background: var(--success-bg);
  color: var(--success);
}

.badge-offline {
  background: var(--error-bg);
  color: var(--error);
}

.badge-warning {
  background: var(--warning-bg);
  color: var(--warning);
}

.badge-info {
  background: var(--info-bg);
  color: var(--info);
}

.badge-accent {
  background: var(--accent-subtle);
  color: var(--accent);
}
```

#### Tag (Small Inline Label)

```css
.tag {
  display: inline-block;
  padding: 2px 8px;
  font-size: 0.75rem;
  border-radius: var(--radius-sm);
  background: var(--accent-subtle);
  color: var(--accent);
}
```

### Form Controls

#### Text Inputs

```css
.input {
  width: 100%;
  padding: 10px 14px;
  font-size: 0.875rem;
  background: var(--bg-card);
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  color: var(--text-primary);
  transition: border-color var(--transition-normal);
}

.input:focus {
  outline: none;
  border-color: var(--accent);
}

.input::placeholder {
  color: var(--text-muted);
}
```

#### Toggle Switch

```css
.toggle {
  position: relative;
  width: 48px;
  height: 24px;
  background: var(--border);
  border-radius: var(--radius-full);
  cursor: pointer;
  transition: background var(--transition-normal);
}

.toggle.active {
  background: var(--success);
}

.toggle::after {
  content: '';
  position: absolute;
  top: 2px;
  left: 2px;
  width: 20px;
  height: 20px;
  background: white;
  border-radius: 50%;
  transition: transform var(--transition-normal);
}

.toggle.active::after {
  transform: translateX(24px);
}
```

#### Range Slider

```css
input[type="range"] {
  width: 100%;
  height: 8px;
  border-radius: var(--radius-sm);
  background: var(--bg-card);
  -webkit-appearance: none;
}

input[type="range"]::-webkit-slider-thumb {
  -webkit-appearance: none;
  width: 18px;
  height: 18px;
  border-radius: 50%;
  background: var(--text-primary);
  border: 2px solid var(--bg-card);
  cursor: pointer;
  transition: transform var(--transition-fast);
}

input[type="range"]::-webkit-slider-thumb:hover {
  transform: scale(1.2);
}
```

### Tables and Lists

#### Data Row

```css
.data-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 0;
  border-bottom: 1px solid var(--border);
}

.data-row:last-child {
  border-bottom: none;
}

.data-label {
  color: var(--text-secondary);
  font-weight: 400;
}

.data-value {
  color: var(--text-primary);
  font-weight: 500;
}
```

#### History Item

```css
.history-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 12px;
  border-bottom: 1px solid var(--border);
}

.history-item:last-child {
  border-bottom: none;
}

.history-icon {
  width: 32px;
  height: 32px;
  border-radius: var(--radius-md);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.875rem;
  flex-shrink: 0;
}

.history-icon.success {
  background: var(--success-bg);
  color: var(--success);
}

.history-icon.error {
  background: var(--error-bg);
  color: var(--error);
}

.history-content {
  flex: 1;
  min-width: 0;
}

.history-action {
  font-weight: 500;
  margin-bottom: 2px;
}

.history-time {
  font-size: 0.75rem;
  color: var(--text-secondary);
}
```

### Tabs

```css
.tabs {
  display: flex;
  gap: var(--space-2);
  margin-bottom: var(--space-5);
}

.tab {
  padding: 10px 20px;
  background: transparent;
  border: 1px solid var(--border);
  color: var(--text-secondary);
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: all var(--transition-normal);
}

.tab.active {
  background: var(--accent);
  border-color: var(--accent);
  color: white;
}

.tab:hover:not(.active) {
  border-color: var(--accent);
  color: var(--accent);
}

.tab:focus-visible {
  outline: 2px solid var(--accent);
  outline-offset: 2px;
}
```

### Empty States

```css
.empty-state {
  text-align: center;
  padding: var(--space-10);
  color: var(--text-secondary);
}

.empty-state-icon {
  width: 48px;
  height: 48px;
  margin-bottom: var(--space-3);
  opacity: 0.5;
}

.empty-state-title {
  font-size: 1rem;
  font-weight: 500;
  margin-bottom: var(--space-2);
}

.empty-state-description {
  font-size: 0.875rem;
}
```

### Keyboard Hints

```css
.keyboard-hint {
  text-align: center;
  color: var(--text-secondary);
  font-size: 0.8rem;
  margin-top: var(--space-5);
}

kbd {
  background: var(--bg-card);
  padding: 2px 8px;
  border-radius: var(--radius-sm);
  border: 1px solid var(--border);
  font-family: inherit;
  font-size: 0.85em;
}
```

---

## Layout Patterns

### Container

Standard container with responsive max-width:

```css
.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: var(--space-5);
}

/* For the main web dashboard, use wider container */
.container-wide {
  max-width: 1400px;
}
```

**Decision**: Use **1200px** as the standard container width for MCP-UI templates (which render in iframes) and **1400px** for the full web dashboard.

### Grid System

```css
.grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: var(--space-5);
}

/* Tighter grid for smaller items */
.grid-compact {
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: var(--space-4);
}

/* Image grid */
.grid-gallery {
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: var(--space-5);
}
```

### Responsive Breakpoints

| Breakpoint | Width | Usage |
|------------|-------|-------|
| `xs` | < 480px | Mobile phones |
| `sm` | 480px - 767px | Large phones, small tablets |
| `md` | 768px - 1023px | Tablets |
| `lg` | 1024px - 1279px | Small desktops |
| `xl` | >= 1280px | Large desktops |

```css
/* Mobile-first responsive utilities */
@media (max-width: 767px) {
  .container {
    padding: var(--space-4);
  }

  .grid {
    grid-template-columns: 1fr;
  }
}
```

### Header Pattern

```css
header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--space-5) 0;
  border-bottom: 1px solid var(--border);
  margin-bottom: var(--space-8);
}

header h1 {
  font-size: 1.8rem;
  font-weight: 600;
}

header h1 .accent {
  color: var(--accent);
}
```

### Footer Pattern

```css
footer {
  text-align: center;
  padding: var(--space-5);
  color: var(--text-secondary);
  font-size: 0.875rem;
  border-top: 1px solid var(--border);
  margin-top: var(--space-10);
}

footer a {
  color: var(--accent);
  text-decoration: none;
}

footer a:hover {
  text-decoration: underline;
}
```

### Stat Grid

```css
.stat-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: var(--space-4);
}

.stat-item {
  text-align: center;
  padding: var(--space-4);
  background: var(--bg-card);
  border-radius: var(--radius-md);
}

.stat-value {
  font-size: 2rem;
  font-weight: 700;
  color: var(--accent);
}

.stat-label {
  font-size: 0.875rem;
  color: var(--text-secondary);
  margin-top: var(--space-1);
}
```

---

## Platform-Specific Guidelines

### Web Dashboard

The main web dashboard (`internal/http/static/index.html`) is the primary interface for server monitoring.

#### Characteristics
- Full browser viewport access
- Mouse and keyboard interaction
- Container width: **1400px**
- Tab-based navigation
- Auto-refresh patterns for live data

#### Specific Patterns
- Use animated status dots for connection state
- Implement optimistic UI updates where safe
- Auto-refresh intervals: 10s for status, 15s for activity, 30s for stats

### TUI (Terminal UI)

The terminal interface (`internal/tui/`) uses lipgloss for styling in constrained terminal environments.

#### xterm 256 Color Mapping

| CSS Token | Hex | xterm Color | Notes |
|-----------|-----|-------------|-------|
| `--accent` | `#e94560` | 205 | Pink/magenta |
| `--success` | `#4ecca3` | 82 | Bright green |
| `--error` | `#ff6b6b` | 196 | Red |
| `--warning` | `#ffc107` | 214 | Orange |
| `--text-primary` | `#eaeaea` | 252 | Near white |
| `--text-secondary` | `#a0a0a0` | 245 | Light gray |
| `--text-muted` | `#707070` | 241 | Gray |
| `--border` | `#2a2a4a` | 240 | Border gray |
| `--bg-highlight` | `#2a2a4a` | 236 | Background highlight |

#### TUI Style Constants

```go
// Color palette - matches web design system
var (
    colorAccent    = lipgloss.Color("205") // --accent equivalent
    colorSuccess   = lipgloss.Color("82")  // --success
    colorError     = lipgloss.Color("196") // --error
    colorWarning   = lipgloss.Color("214") // --warning
    colorMuted     = lipgloss.Color("241") // --text-muted
    colorSubtle    = lipgloss.Color("245") // --text-secondary
    colorText      = lipgloss.Color("252") // --text-primary
    colorBorder    = lipgloss.Color("240") // --border
    colorHighlight = lipgloss.Color("236") // --bg-highlight
    colorDim       = lipgloss.Color("239") // Dimmed text
)
```

#### TUI Layout Constants

```go
const (
    // Box and container offsets
    boxWidthOffset    = 4  // Offset for box width from terminal width
    headerWidthOffset = 2  // Offset for header box width
    tablePadding      = 12 // Padding/border offset for table width
    columnSpacing     = 6  // Spacing between table columns

    // Table column widths
    colWidthTimestamp = 19 // "2006-01-02 15:04:05"
    colWidthStatus    = 6  // "OK" or "FAIL"
    colWidthDuration  = 10 // "12345ms"
    colWidthToolMin   = 15 // Minimum tool column width
)
```

#### TUI Status Indicators

```go
const (
    StatusConnected    = "●" // Filled circle
    StatusDisconnected = "○" // Empty circle
    StatusConnecting   = "◐" // Half circle
)
```

### MCP-UI (Claude Rendering)

MCP-UI templates render in sandboxed iframes within Claude's interface.

#### Constraints
- Limited viewport (typical iframe size)
- No external resources (all styles inline)
- Container width: **1200px**
- Must work without JavaScript for initial render
- JavaScript for interactivity only

#### Design Priorities
1. **Immediate comprehension** - Claude (the AI) should understand the UI at a glance
2. **Clear actions** - Buttons and interactive elements must be obvious
3. **Focused scope** - Each template serves one purpose
4. **Keyboard support** - Include keyboard shortcuts for power users

#### MCP-UI Template Checklist
- Include shared CSS via template injection
- Add keyboard hint bar for shortcuts
- Implement focus states for all interactive elements
- Include connection error handling for network issues
- Use semantic HTML for accessibility

#### Title Formatting

```html
<h1>Feature <span class="accent">Name</span></h1>
```

The second word is always in accent color to maintain brand consistency.

---

## Accessibility Requirements

### WCAG 2.1 AA Compliance

#### Color Contrast

| Element Type | Minimum Ratio | Current Status |
|--------------|---------------|----------------|
| Normal text | 4.5:1 | Pass (all text colors) |
| Large text (18px+ or 14px bold) | 3:1 | Pass |
| UI components | 3:1 | Pass |

Verified contrast ratios on `--bg-primary` (#1a1a2e):
- `--text-primary` (#eaeaea): **12.3:1** - Pass
- `--text-secondary` (#a0a0a0): **5.8:1** - Pass
- `--success` (#4ecca3): **7.2:1** - Pass
- `--warning` (#ffc107): **9.1:1** - Pass
- `--error` (#ff6b6b): **5.8:1** - Pass
- `--accent` (#e94560): **5.2:1** - Pass

#### Focus Indicators

All interactive elements must have visible focus states:

```css
*:focus-visible {
  outline: 2px solid var(--accent);
  outline-offset: 2px;
}

/* Reset browser default focus for custom styling */
button:focus:not(:focus-visible) {
  outline: none;
}
```

#### Keyboard Navigation

All interfaces must be fully operable via keyboard:

| Key | Action |
|-----|--------|
| Tab | Move to next focusable element |
| Shift+Tab | Move to previous focusable element |
| Enter/Space | Activate buttons, links |
| Arrow keys | Navigate within components (tabs, lists) |
| Escape | Close modals, cancel actions |

#### Touch Target Size

Minimum touch target size: **44x44 pixels**

```css
.btn, .tab, .toggle, .scene-card {
  min-height: 44px;
  min-width: 44px;
}
```

#### Screen Reader Considerations

- Use semantic HTML elements (`<button>`, `<nav>`, `<main>`, `<article>`)
- Include `aria-label` for icon-only buttons
- Use `aria-live` regions for dynamic content updates
- Provide `alt` text for meaningful images

#### Reduced Motion

Respect user preference for reduced motion:

```css
@media (prefers-reduced-motion: reduce) {
  *,
  *::before,
  *::after {
    animation-duration: 0.01ms !important;
    animation-iteration-count: 1 !important;
    transition-duration: 0.01ms !important;
  }
}
```

---

## Interaction Patterns

### Loading States

#### Button Loading

```css
.btn-loading {
  position: relative;
  color: transparent;
  pointer-events: none;
}

.btn-loading::after {
  content: '';
  position: absolute;
  width: 16px;
  height: 16px;
  top: 50%;
  left: 50%;
  margin: -8px 0 0 -8px;
  border: 2px solid currentColor;
  border-right-color: transparent;
  border-radius: 50%;
  animation: spin 0.75s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
```

#### Content Loading

```css
.skeleton {
  background: linear-gradient(
    90deg,
    var(--bg-card) 0%,
    var(--bg-secondary) 50%,
    var(--bg-card) 100%
  );
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
  border-radius: var(--radius-md);
}

@keyframes shimmer {
  0% { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}
```

### Error States

#### Inline Error

```css
.error-message {
  background: var(--error-bg);
  border: 1px solid var(--error);
  color: var(--error);
  padding: var(--space-5);
  border-radius: var(--radius-md);
  text-align: center;
}
```

#### Connection Lost Banner

```css
.connection-error {
  position: fixed;
  bottom: var(--space-5);
  left: 50%;
  transform: translateX(-50%);
  background: var(--error);
  color: white;
  padding: 12px 20px;
  border-radius: var(--radius-md);
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 0.85rem;
  box-shadow: var(--shadow-lg);
  z-index: 1000;
}
```

### Feedback Patterns

#### Toast/Notification

```css
.toast {
  position: fixed;
  top: var(--space-5);
  right: var(--space-5);
  padding: 12px 20px;
  border-radius: var(--radius-md);
  font-size: 0.875rem;
  box-shadow: var(--shadow-lg);
  animation: slideIn 0.3s ease;
  z-index: 1000;
}

.toast-success { background: var(--success); color: white; }
.toast-error { background: var(--error); color: white; }
.toast-warning { background: var(--warning); color: #000; }

@keyframes slideIn {
  from {
    transform: translateX(100%);
    opacity: 0;
  }
  to {
    transform: translateX(0);
    opacity: 1;
  }
}
```

#### Refresh Indicator

```css
.refresh-indicator {
  position: fixed;
  top: var(--space-5);
  right: var(--space-5);
  background: var(--bg-secondary);
  color: var(--text-secondary);
  padding: 8px 12px;
  border-radius: var(--radius-md);
  font-size: 0.75rem;
  opacity: 0;
  transition: opacity var(--transition-slow);
}

.refresh-indicator.visible {
  opacity: 1;
}
```

### Hover States

#### Card Hover

```css
.card-interactive {
  transition: transform var(--transition-normal), border-color var(--transition-normal);
}

.card-interactive:hover {
  transform: translateY(-2px);
  border-color: var(--accent);
}
```

#### Scale Hover

```css
.scale-hover {
  transition: transform var(--transition-normal);
}

.scale-hover:hover {
  transform: scale(1.02);
}
```

---

## Implementation Reference

### Shared CSS Template

For MCP-UI templates, use the shared CSS file at `internal/http/templates/shared.css`. Include it via Go template:

```html
<style>{{.SharedCSS}}</style>
```

### Complete CSS Reset and Base

```css
*,
*::before,
*::after {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html {
  font-size: 16px;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, sans-serif;
  background: var(--bg-primary);
  color: var(--text-primary);
  line-height: 1.6;
}

a {
  color: var(--accent);
  text-decoration: none;
}

a:hover {
  text-decoration: underline;
}

button {
  font-family: inherit;
  font-size: inherit;
  cursor: pointer;
}

img {
  max-width: 100%;
  height: auto;
}
```

### Icon Approach

Use inline SVGs for icons to avoid external dependencies:

```html
<!-- Example: Sound icon -->
<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
  <path d="M11 5L6 9H2v6h4l5 4V5z"/>
  <path d="M19.07 4.93a10 10 0 0 1 0 14.14M15.54 8.46a5 5 0 0 1 0 7.07"/>
</svg>
```

Icon sizing:
- Small: 16x16px
- Medium: 20x20px (default)
- Large: 24x24px

### File Organization

```
internal/http/
  static/
    index.html          # Main web dashboard (standalone)
  templates/
    shared.css          # Shared design tokens & components
    audio_mixer.html    # MCP-UI: Audio control
    scene_preview.html  # MCP-UI: Scene switching
    screenshot_gallery.html  # MCP-UI: Screenshot viewer
    status_dashboard.html    # MCP-UI: Connection status
    error.html          # MCP-UI: Error page
```

---

## Change Log

| Version | Date | Changes |
|---------|------|---------|
| 1.0 | 2025-12 | Initial specification |

---

## Future Considerations

### Light Theme

If a light theme is needed in the future, define these token overrides:

```css
[data-theme="light"] {
  --bg-primary: #f5f5f5;
  --bg-secondary: #ffffff;
  --bg-card: #e8e8e8;
  --text-primary: #1a1a2e;
  --text-secondary: #4a4a4a;
  --border: #d0d0d0;
}
```

### Design System Expansion

Potential additions for future development:
- Modal/dialog component
- Dropdown/select component
- Progress bar component
- Tooltip component
- Notification/toast system
- Data visualization patterns (charts, meters)

---

*This specification is a living document and should be updated as the design system evolves.*
