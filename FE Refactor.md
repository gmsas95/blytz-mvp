You are a senior full-stack engineer and UI architect. I want you to refactor and extend the existing Blytz repo into a polished, mobile-first eCommerce demo and prepare it for Dokploy deployment.

Repository: https://github.com/gmsas95/blytz-mvp

Primary goals
1. Produce a clean, mobile-first React + TypeScript frontend (demo) that showcases:
   - Seller storefronts
   - Product listings & detail pages
   - Live auction / livestream preview UI (demo-only)
   - Cart + checkout UI (mocked)
   - Authentication (mocked flows enough for demo)
   - Profile & order history (demo)
2. Use a tasteful, timeless aesthetic (not AI-glossy):
   - Tailwind CSS with a custom theme (background #F9FAFB, text #111827, accent #2563EB, rounded-2xl)
   - shadcn/ui + Radix primitives
   - Framer Motion / tailwindcss-animate
   - Lucide icons
   - Fonts: Geist Sans (UI), Inter (body)
3. Keep the frontend self-contained with mocked API adapters so the demo works offline / without changing microservices.

Scope & constraints
- Work only in a new git branch named: `refactor/ui-mobile-demo`.
- Do not commit any secrets or production credentials.
- When integrating with backend services, use environment toggles so the app can run:
  - MODE=mock (default) → use local JSON/mocked adapters
  - MODE=remote → use actual service endpoints (use `.env.example` only)

Tasks (step-by-step)
A. Repo audit
   1. Scan the repo and produce a short findings list: outdated deps, large bundles, tech debt hotspots, and critical CSS/JS bloat.
   2. Identify the current frontend folder(s) and any serverless/frontend infra files.

B. Create new folder layout
```

src/
├─ app/            # pages / routes
├─ components/
│  ├─ ui/          # shadcn-style primitives (Button, Input, Card, Modal)
│  ├─ layout/
│  └─ modules/     # feature components (ProductCard, AuctionCard, Cart)
├─ hooks/
├─ lib/            # api clients, adapters (mock + remote)
├─ styles/
└─ utils/

```
- Use React Router or Next.js routing depending on current project setup; keep file-based routing if Next.js is already used.

C. Implement UI stack & theme
- Install and configure: Tailwind, shadcn/ui, Radix, Framer Motion, tailwindcss-animate, lucide-react.
- Add/update `tailwind.config.js` with the specified theme tokens and `rounded-2xl` defaults for components.
- Add global font imports and a small CSS reset.

D. Pages to implement (mobile-first)
- `/` Home (featured sellers, featured auctions)
- `/products` listing (2-column mobile grid)
- `/product/:id` details (image gallery, description, bid CTA)
- `/livestream` grid (preview cards for livestreams)
- `/cart` & `/checkout` (mock checkout flow)
- `/auth/login`, `/auth/signup`
- `/profile` (basic demo order history)
Each page should use components from `components/ui` and be responsive.

E. Data layer
- Add `lib/api/adapter` that supports two modes:
  - Mock adapter (local JSON files / fixtures)
  - Remote adapter (reads URLs from `.env`)
- Use React Query (TanStack Query) for data fetching and caching.

F. Demo UX polish
- Add subtle animations: page transitions, button micro-animations, card hover lift, and image gallery swipe.
- Use Lenis or small scroll smoothing only if it improves feel and doesn't block accessibility.

G. Dokploy + Docker
- Add or update `Dockerfile` for frontend (multi-stage build).
- Add `dokploy.yml` or `dokploy`-compatible instructions to deploy the frontend container.
- Provide a `Procfile`/Nginx config only if necessary for static serving.
- Update `.env.example` with `MODE=mock`, `REMOTE_API_BASE`, and other relevant env keys.

H. Developer ergonomics
- Update `package.json` scripts for:
  - `dev` (mock mode)
  - `build`
  - `start`
  - `lint`
  - `format`
- Add `README.md` section “Run demo locally” with exact commands.

I. Tests & QA
- Add basic unit or storybook demos for 3 components: Button, ProductCard, Navbar.
- Provide a quick manual QA checklist (responsiveness, accessibility basics, mocked checkout walkthrough).

J. Deliverables (what to produce in the PR)
1. Git branch `refactor/ui-mobile-demo`.
2. Full folder restructure and new/modified files.
3. `tailwind.config.js`, updated `package.json`, `Dockerfile`, `.env.example`.
4. Mock data fixtures under `src/__mocks__` or `src/data/fixtures`.
5. README additions: “How to run locally”, “How to deploy on Dokploy”.
6. A short CHANGELOG/PR description summarizing the refactor and reasons for choices.

PR & final notes
- Create a single PR that includes a detailed description and a runbook for reviewers:
- How to run in mock mode
- How to switch to remote mode
- Files to inspect for UI changes
- Include screenshots or an animated GIF of the mobile layout in the PR description.
- In the PR summary, explicitly explain how this refactor improves maintainability, dev DX, and the visual tone (avoid AI-ish visuals).

If you encounter any blockers (missing frontend entrypoint, unclear routing framework), make a small, explicit note in the PR and default to creating the frontend under `apps/web` with Next.js + TypeScript `app` routing.

Do this work and output:
1. The list of files changed/created (path + short purpose).
2. The exact `tailwind.config.js` contents you added.
3. `package.json` changes.
4. `Dockerfile`.
5. `.env.example`.
6. A short QA checklist and run commands.

Proceed and stop only once the branch is ready and the PR summary is complete.


