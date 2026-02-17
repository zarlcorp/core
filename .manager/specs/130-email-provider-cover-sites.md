# 130: Email Provider Cover Sites

## Objective
Build 4 static cover websites for the email provider burner domains, each with unique branding and a fake login page. These make the domains look like real email services when visited.

## Context
We own 4 burner domains styled as email providers. If someone receives email from `user@moxmail.site` and visits the domain, an empty/parked page is a dead giveaway. Each domain needs a plausible landing page with a login form to look like a real email service.

## Requirements

### Site structure
Each site is a directory under `sites/<domain>/` containing static HTML, CSS, and minimal JS (for login form interactivity only). No frameworks, no build step — pure static files.

```
sites/
├── moxmail.site/
│   ├── index.html
│   └── style.css
├── jotmail.xyz/
│   ├── index.html
│   └── style.css
├── snapmail.icu/
│   ├── index.html
│   └── style.css
└── fogmail.space/
    ├── index.html
    └── style.css
```

### Each site must have
1. **Unique branding** — distinct name, logo (text-based/CSS, no images needed), color scheme, and tagline
2. **Login page** — email + password fields, "Sign In" button, "Forgot password?" link
3. **Login behavior** — form submits to nothing (or shows "Invalid credentials" on submit). Must look functional.
4. **Footer** — copyright notice with current year, privacy policy link (can be `#`), terms link (can be `#`)
5. **Responsive** — works on mobile and desktop
6. **Professional quality** — must pass a casual inspection. No placeholder text like "Lorem ipsum".

### Domain-specific branding

**moxmail.site** — "Moxmail"
- Tagline: "Simple, reliable email"
- Vibe: Clean, modern, blue/white. Think Fastmail.
- Minimal landing: logo, tagline, login form

**jotmail.xyz** — "Jotmail"
- Tagline: "Email that gets out of the way"
- Vibe: Lightweight, fast, green/white. Minimalist.
- Minimal landing: logo, tagline, login form

**snapmail.icu** — "Snapmail"
- Tagline: "Fast, private email"
- Vibe: Bold, privacy-focused, dark/orange. Think a smaller ProtonMail.
- Minimal landing: logo, tagline, login form, "End-to-end encrypted" badge

**fogmail.space** — "Fogmail"
- Tagline: "Your email, hidden in the fog"
- Vibe: Mysterious, secure, dark purple/grey. Privacy-first.
- Minimal landing: logo, tagline, login form, privacy messaging

### What NOT to build
- No actual email functionality
- No backend
- No JavaScript frameworks
- No images (use CSS/text for logos)
- No "about" or "pricing" pages — just the login/landing page

## Target Repo
zarlcorp/cover-sites

## Agent Role
frontend

## Files to Create
- sites/moxmail.site/index.html
- sites/moxmail.site/style.css
- sites/jotmail.xyz/index.html
- sites/jotmail.xyz/style.css
- sites/snapmail.icu/index.html
- sites/snapmail.icu/style.css
- sites/fogmail.space/index.html
- sites/fogmail.space/style.css

## Notes
- Each site should feel like a different company — different fonts (use Google Fonts or system fonts), different layouts, different color palettes
- The login form should look real: proper input types, placeholder text, autofocus on email field
- On form submit, prevent default and show "Invalid email or password" error message — this makes it feel like a real login
- Keep CSS in a single file per site, no external dependencies beyond optional Google Fonts
- Total size per site should be tiny — under 20KB
