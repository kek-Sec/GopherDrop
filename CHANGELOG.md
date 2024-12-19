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

