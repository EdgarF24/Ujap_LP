# ─────────────────────────────────────────────────────────────────────────────
# Nginx API Gateway — Dockerfile
# ─────────────────────────────────────────────────────────────────────────────
FROM nginx:alpine

# Install curl for healthchecks and remove the default Nginx config
RUN apk --no-cache add curl && \
    rm /etc/nginx/conf.d/default.conf

# Copy our custom gateway config
COPY nginx.conf /etc/nginx/conf.d/gateway.conf

# Expose HTTP port
EXPOSE 80

# Use the default nginx entrypoint (starts nginx in the foreground)
CMD ["nginx", "-g", "daemon off;"]
