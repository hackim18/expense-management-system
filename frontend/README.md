# Expense Management System - Frontend

Nuxt frontend using TailwindCSS + DaisyUI.

## Features
- Login and register
- Expense dashboard with status filters and pagination
- Expense submission form with IDR formatting and approval warning
- Manager approval queue with approve/reject + notes
- Responsive layout in Indonesian

## Setup
Install dependencies:
```bash
npm install
```

Run dev server:
```bash
npm run dev
```

Open http://localhost:3000

## Environment
- `NUXT_PUBLIC_API_BASE` (default `http://localhost:8080`)

## Production
```bash
npm run build
npm run preview
```

## Docker
Use the root `docker-compose.yml` to run full stack:
```bash
docker compose up --build
```

## Notes
- API base URL can be overridden with `NUXT_PUBLIC_API_BASE`.
- Make sure backend is running before using the UI.
