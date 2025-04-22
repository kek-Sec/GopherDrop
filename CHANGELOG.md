# Changelog

## [1.0.7] - 2025-04-22
### Added
- Customizable application title and description through environment variables
- New build arguments in Dockerfile: VITE_APP_TITLE and VITE_APP_DESCRIPTION
- Support for title and description customization in docker-compose configuration

### [1.0.6]
* `cmd/server/main.go`: Added rate limiting to POST requests.
* `ui/src/pages/Create.vue`: Added error handling for 429 responses.

### [1.0.5]
UI Enhancements:

* `ui/index.html`: Added a link to the Animate.css library for animation effects.
* [`ui/src/App.vue`: Integrated animation classes into various elements, including the header logo, buttons, and alerts.

Password Management Enhancements:
* `ui/src/pages/Create.vue`: Added functionality to toggle password visibility and generate random passwords, along with corresponding tooltips and animations.
404 Error Page Improvements:
* `ui/src/pages/Error404.vue`: Redesigned the 404 error page with a more user-friendly card layout, including an icon, message, and a button to navigate back home.

Other Enhancements:
* `ui/src/pages/Create.vue`, `ui/src/pages/View.vue`: Applied animation classes to various elements to provide a more dynamic user experience.

### [1.0.4]
- Reworked CD pipelines to follow semver tagging
- added version.yaml

### [1.0.1]

#### Added
- Initial unified Dockerfile combining backend and frontend.
- Support for Traefik reverse proxy configuration.
- Environment variables for configuring API URLs dynamically.
- Ability to copy the generated shareable link with improved formatting.
- New CORS configuration to support wildcard origins in development mode.
- Automatic HTTPS redirection using Traefik middlewares.

#### Fixed
- Corrected Vite's `VITE_API_URL` handling to avoid hardcoded URLs.
- Resolved 404 errors for static assets when accessed via Traefik.
- Fixed MIME type issues for serving CSS files behind Traefik.

---

### [1.0.0] - 2024-12-18

#### Added
- Secure one-time secret sharing with encrypted text and file support.
- Password protection for shared secrets.
- One-time retrieval mechanism to ensure secrets are accessed only once.
- Expiration settings for shared secrets.
- Responsive UI built with Vue.js and Vuetify.
- Dockerized deployment with `docker-compose.yml` for easy setup.
- Nginx reverse proxy configuration for serving the frontend and API.
- Health checks for PostgreSQL in `docker-compose.yml`.

#### Changed
- Updated Docker images to use multi-stage builds for backend and frontend.
- Improved project documentation and added sections for installation, configuration, and deployment.

---
