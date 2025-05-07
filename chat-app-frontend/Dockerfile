FROM node:23-slim AS build
WORKDIR /app
COPY package.json package.json
COPY package-lock.json package-lock.json
ARG VITE_DOMAIN_HTTPS
ARG VITE_DOMAIN_WSS

ENV VITE_DOMAIN_HTTPS=$VITE_DOMAIN_HTTPS
ENV VITE_DOMAIN_WSS=$VITE_DOMAIN_WSS
RUN npm ci
COPY . .
RUN npm run build

FROM nginx:1.28-alpine
COPY --from=build /app/dist /opt/site
COPY nginx.conf /etc/nginx/nginx.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
