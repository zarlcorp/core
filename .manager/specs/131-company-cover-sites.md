# 131: Company Cover Sites

## Objective
Build 6 static marketing landing pages for the company burner domains, each with unique branding that matches its cover story. These make the domains look like real businesses when visited.

## Context
We own 6 burner domains styled as various companies. If someone checks one of these domains and finds nothing, it's suspicious. Each needs a single-page marketing site that looks like a real (small/startup) company.

## Requirements

### Site structure
Each site is a directory under `sites/<domain>/` containing static HTML and CSS. No frameworks, no build step.

```
sites/
├── nordvik.work/
│   ├── index.html
│   └── style.css
├── hawkpoint.live/
│   ├── index.html
│   └── style.css
├── axontel.network/
│   ├── index.html
│   └── style.css
├── clearwire.digital/
│   ├── index.html
│   └── style.css
├── nexacom.online/
│   ├── index.html
│   └── style.css
└── calvera.run/
    ├── index.html
    └── style.css
```

### Each site must have
1. **Unique branding** — company name, text-based logo (CSS), color scheme, tagline
2. **Hero section** — company name, tagline, brief description (1-2 sentences)
3. **Services/features section** — 3-4 bullet points or cards describing what the company does
4. **Contact section** — "Get in touch" with a generic email (e.g. info@domain) or contact form (non-functional)
5. **Footer** — copyright, year, location hint (city/country)
6. **Responsive** — works on mobile and desktop
7. **Professional quality** — must pass a casual inspection

### Domain-specific branding

**nordvik.work** — "Nordvik Industries"
- Cover: Oil/energy services, offshore operations
- Tagline: "Integrated energy solutions for the North Sea"
- Color: Dark navy, steel grey, orange accents
- Sections: Offshore services, Pipeline maintenance, Environmental compliance, Safety consulting
- Footer location: Stavanger, Norway

**hawkpoint.live** — "Hawkpoint Advisory"
- Cover: Management consulting
- Tagline: "Strategic clarity for complex decisions"
- Color: Clean white, charcoal, gold accents
- Sections: Strategy, Operations, Digital transformation, M&A advisory
- Footer location: London, UK

**axontel.network** — "Axontel"
- Cover: Telecom/network infrastructure provider
- Tagline: "Connecting businesses to what matters"
- Color: Deep blue, cyan accents, dark backgrounds
- Sections: Enterprise networking, VOIP solutions, Managed connectivity, Data centre services
- Footer location: Amsterdam, Netherlands

**clearwire.digital** — "Clearwire Digital"
- Cover: Digital marketing platform/agency
- Tagline: "Data-driven growth for modern brands"
- Color: Vibrant gradient (purple to blue), white text
- Sections: Campaign management, Analytics & insights, SEO & content, Social media
- Footer location: Dublin, Ireland

**nexacom.online** — "Nexacom"
- Cover: Internal communications & sentiment analysis platform
- Tagline: "Understand your team, elevate your culture"
- Color: Soft blue, warm grey, teal accents
- Sections: Real-time sentiment, Team pulse surveys, Analytics dashboard, Integration hub
- Footer location: Berlin, Germany

**calvera.run** — "Calvera"
- Cover: Running/exercise platform
- Tagline: "Every run tells a story"
- Color: Energetic — coral/orange, white, dark accents
- Sections: Route tracking, Training plans, Community challenges, Performance analytics
- Footer location: Barcelona, Spain

### What NOT to build
- No actual functionality behind any buttons/forms
- No JavaScript frameworks
- No images (use CSS gradients, shapes, icons via Unicode or CSS)
- No multi-page sites — single page each
- No blog or news sections

## Target Repo
zarlcorp/cover-sites

## Agent Role
frontend

## Files to Create
- sites/nordvik.work/index.html
- sites/nordvik.work/style.css
- sites/hawkpoint.live/index.html
- sites/hawkpoint.live/style.css
- sites/axontel.network/index.html
- sites/axontel.network/style.css
- sites/clearwire.digital/index.html
- sites/clearwire.digital/style.css
- sites/nexacom.online/index.html
- sites/nexacom.online/style.css
- sites/calvera.run/index.html
- sites/calvera.run/style.css

## Notes
- Each site must look like a DIFFERENT company — varied layouts, different fonts, different visual styles
- Use system fonts or Google Fonts (linked via CDN) — no font files to host
- Keep it minimal but professional — a well-designed single page is more convincing than a sloppy multi-page site
- No "Under construction" or "Coming soon" — these are meant to look like established (if small) companies
- Contact forms can have `action="#"` with `onsubmit="return false"` or just show "Thank you" on submit
- Total size per site should be tiny — under 20KB
