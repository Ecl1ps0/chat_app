FROM node:23-slim AS build
WORKDIR /app
COPY package.json package.json
COPY package-lock.json package-lock.json
RUN npm ci
COPY . .
RUN npm run build

FROM nginx:1.28-alpine
ENV VITE_DOMAIN_HTTPS = http://172.31.24.209
ENV VITE_DOMAIN_WSS = ws://172.31.24.209
COPY --from=build /app/dist /opt/site
COPY nginx.conf /etc/nginx/nginx.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
