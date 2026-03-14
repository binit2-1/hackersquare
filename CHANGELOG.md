# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2026-03-14

### Added
- **Core Engine:** Custom Go backend API with Colly-based web scrapers running on 12-hour cron jobs.
- **Data Pipelines:** MLH, Devfolio and Unstop scraper integrations.
- **Deep Scrape Router:** Proximity-based regex engine to extract unstructured prize pool data while ignoring false positives.
- **Frontend UI:** Next.js App Router interface with Shadcn UI and Tailwind CSS.
- **Authentication:** 1st-party GitHub OAuth integration via Next.js Rewrite proxies.
- **Telemetry:** Stealth Umami analytics proxy to capture 100% of pageviews while respecting user privacy.
- **Database:** PostgreSQL integration via Neon with automated schema migrations.