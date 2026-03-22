# Each stage/FROM is a different filesystem
# Docker caches each line
# Stage 1: install node and packages in /app with clean install
FROM node:20-alpine AS base
WORKDIR /app
COPY package*.json ./
RUN npm ci

# Stage 2: copy base, npm run dev. Note: exposing happens in compose, this is just documentation
FROM base AS dev
COPY . .
EXPOSE 5173
CMD ["npm", "run", "dev"]

# Stage 3: copy base, npm run build outputs to /app/dist
FROM base AS builder
COPY . .
RUN npm run build

# use nginx to run files without baggage from node since react just needs static files
FROM nginx:alpine AS prod
RUN rm -rf /usr/share/nginx/html/*
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]