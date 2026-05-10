FROM mcr.microsoft.com/devcontainers/go:1-1.22-bookworm

# The base image may include 3rd-party apt sources (e.g. Yarn) without keys.
# Remove any such sources to keep apt-get update reliable in Codespaces.
RUN rm -f /etc/apt/sources.list.d/yarn.list /etc/apt/sources.list.d/yarn.list.save \
  && apt-get update \
  && apt-get install -y --no-install-recommends ca-certificates curl gnupg lsb-release \
  && rm -rf /var/lib/apt/lists/*

# Install Docker CLI + Compose plugin (daemon comes from host socket mount if available)
RUN set -eux; \
  install -m 0755 -d /etc/apt/keyrings; \
  curl -fsSL https://download.docker.com/linux/debian/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg; \
  chmod a+r /etc/apt/keyrings/docker.gpg; \
  echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/debian $(. /etc/os-release && echo $VERSION_CODENAME) stable" > /etc/apt/sources.list.d/docker.list; \
  apt-get update; \
  apt-get install -y --no-install-recommends docker-ce-cli docker-compose-plugin; \
  rm -rf /var/lib/apt/lists/*
