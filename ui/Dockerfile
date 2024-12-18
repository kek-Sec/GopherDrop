# Build stage
FROM node:18-alpine AS build

WORKDIR /app
COPY package.json package-lock.json ./
RUN npm install --legacy-peer-deps

ARG VITE_API_URL
ENV VITE_API_URL=${VITE_API_URL}

COPY . .
RUN npm run build

# Production stage
FROM nginx:alpine

# Add OCI Image Spec labels
ARG GIT_COMMIT_SHA
ARG GIT_VERSION

LABEL org.opencontainers.image.title="GopherDrop UI" \
      org.opencontainers.image.description="UI for GopherDrop, a secure one-time secret sharing service" \
      org.opencontainers.image.source="https://github.com/kek-Sec/gopherdrop" \
      org.opencontainers.image.revision="${GIT_COMMIT_SHA}" \
      org.opencontainers.image.version="${GIT_VERSION}" \
      org.opencontainers.image.url="https://github.com/kek-Sec/gopherdrop" \
      org.opencontainers.image.documentation="https://github.com/kek-Sec/gopherdrop" \
      org.opencontainers.image.licenses="MIT"

COPY --from=build /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
