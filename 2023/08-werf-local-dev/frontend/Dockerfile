# этап сборки (build stage)
FROM node:lts-alpine@sha256:1a526b97cace6b4006256570efa1a29cd1fe4b96a5301f8d48e87c5139438a45 as build-stage
WORKDIR /app
COPY frontend/package*.json ./
RUN npm install
COPY frontend/. .
RUN npm run build

# этап production (production-stage)
FROM nginx:stable-alpine@sha256:d4d72ee8e6d028c5ad939454164d3645be2d38afb5c352277926a48e24abf4fa as production-stage
COPY --from=build-stage /app/dist /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]