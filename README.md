# Hypabase

A starter project, describing a recommended structure for new applications.

## Getting Started

- copy project, make build, woooo
- make dev, edit, woooo

## Key Components

- Golang server backend (to write custom APIs/functions).
- Markdown generation (to make authoring blogs/pages easier for everyone).
- esbuild (to integrate client-side bundling without hard reliance on npm).
- htmx (minimalist approach to create web applications)
- plain css (you may choose Tailwind/SCSS/SASS later but start simple).
- SQLite (minimal, integrated SQL database).

## Add ons

_These are not essential but often required, so here's the default plan_

- TailwindCSS (for projects with multiple contributors where css gets weird).
- PostgreSQL (for scaling beyond a single machine).
- CMS? (tina? contentful? ghost? strapi?) not sure yet

## Deployment

This project is designed to compile to a single binary, which can be controlled
using `GOOS` and `GOARCH`, you can find a full list of options with `go tool
dist list`

For example, to build your application to run on a linux server:
```
GOOS=linux GOARCH=amd64 go build ./cmd/main.go
```
