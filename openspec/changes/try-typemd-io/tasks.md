## 1. Project Setup

- [ ] 1.1 Scaffold React project in `websites/try/` with Vite, TypeScript, and React
- [ ] 1.2 Install dependencies: shadcn/ui, tailwindcss, react-markdown, gray-matter, react-router
- [ ] 1.3 Configure shadcn/ui with base theme and components (Button, Input, Checkbox, Card)
- [ ] 1.4 Set up project structure: `src/lib/`, `src/components/`, `src/pages/`

## 2. VaultStorage Interface & GitHub Backend

- [ ] 2.1 Define TypeScript types: `TypeSchema`, `ObjectSummary`, `Object`, `ObjectFilter`, `WikiLink`, `DisplayProperty`
- [ ] 2.2 Define `VaultStorage` interface with `listTypes()`, `getTypeSchema()`, `listObjects()`, `getObject()`, `searchObjects()`
- [ ] 2.3 Implement `GitHubBackend.connect(repo, token?)` — fetch repo tree via Trees API, validate `.typemd/` exists
- [ ] 2.4 Implement `GitHubBackend.listTypes()` — fetch and parse `.typemd/types/*.yaml`
- [ ] 2.5 Implement `GitHubBackend.listObjects()` — build object index from file tree paths, support type filter
- [ ] 2.6 Implement `GitHubBackend.getObject(id)` — fetch markdown blob, parse frontmatter and body, extract wiki-links
- [ ] 2.7 Implement `GitHubBackend.searchObjects(keyword)` — case-insensitive substring match on name and cached body
- [ ] 2.8 Implement rate limit tracking from `X-RateLimit-Remaining` response headers

## 3. GitHub Connect Page

- [ ] 3.1 Create landing page with repo input field, PAT input field, "Remember token" checkbox, and Connect button
- [ ] 3.2 Implement connection flow: validate repo, check for `.typemd/` directory, show errors for invalid repo/token/non-vault
- [ ] 3.3 Implement token persistence: save to localStorage when opted in, pre-fill on revisit, clear on disconnect
- [ ] 3.4 Add connected state header: show repo name, rate limit status, disconnect button

## 4. Sidebar — Object List

- [ ] 4.1 Create sidebar component with grouped type list (emoji + type name + count as headers)
- [ ] 4.2 Implement expand/collapse toggle on group headers
- [ ] 4.3 Implement object selection — click object to load detail
- [ ] 4.4 Add search input at top of sidebar with instant filtering (flat list, no group headers)
- [ ] 4.5 Highlight currently selected object

## 5. Detail View — Title, Body, Properties

- [ ] 5.1 Create title panel component: `{emoji} {type} · {displayName}`
- [ ] 5.2 Create body panel with react-markdown rendering (headings, lists, code blocks, links)
- [ ] 5.3 Implement wiki-link rendering: parse `[[...]]` syntax, render as clickable links, navigate on click
- [ ] 5.4 Show "(empty)" placeholder when body is empty
- [ ] 5.5 Create properties panel: list schema properties with formatted values
- [ ] 5.6 Render relation properties as clickable links that navigate to the target object
- [ ] 5.7 Render backlinks section (computed from in-memory index)
- [ ] 5.8 Add properties panel toggle button (hidden by default)

## 6. Layout & Responsive

- [ ] 6.1 Implement CSS Grid three-panel layout: sidebar, body, properties
- [ ] 6.2 Add responsive breakpoint: collapse to single-column on narrow screens with list/detail navigation
- [ ] 6.3 Handle empty state: no object selected shows welcome message

## 7. Deployment

- [ ] 7.1 Configure Vite build for static output
- [ ] 7.2 Set up GitHub Actions workflow: build and deploy to GitHub Pages
- [ ] 7.3 Configure custom domain `try.typemd.io` with CNAME
