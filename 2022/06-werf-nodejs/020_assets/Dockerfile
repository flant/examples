FROM node:12-alpine@sha256:d4b15b3d48f42059a15bd659be60afe21762aae9d6cbea6f124440895c27db68 as builder
WORKDIR /app

# [<en>] Copy the files needed to install the application dependencies into the image.
# [<ru>] Копируем в образ файлы, нужные для установки зависимостей приложения.
COPY package.json package-lock.json ./

# [<en>] Install the application dependencies.
# [<ru>] Устанавливаем зависимости приложения.
RUN npm ci

# [<en>] Copy all other application files into the image.
# [<ru>] Копируем в образ все остальные файлы приложения.
COPY . .

# [<en>] Build static asseets.
# [<ru>] Собираем статические файлы.
RUN npm run build

#############################################################################

FROM node:12-alpine@sha256:d4b15b3d48f42059a15bd659be60afe21762aae9d6cbea6f124440895c27db68 as backend
WORKDIR /app

# [<en>] Copy the files needed to install the application dependencies into the image.
# [<ru>] Копируем в образ файлы, нужные для установки зависимостей приложения.
COPY package.json package-lock.json ./

# [<en>] Install the application dependencies.
# [<ru>] Устанавливаем зависимости приложения.
RUN npm ci --production

# [<en>] Copy app files.
# [<ru>] Копируем файлы приложения.
COPY app.js ./
RUN mkdir dist
COPY --from=builder /app/dist/*.html ./dist/
COPY --from=builder /app/bin         ./bin/
COPY --from=builder /app/routes      ./routes/


#############################################################################

# [<en>] Add an NIGINX image with the pre-built assets.
# [<ru>] NGINX-образ с собранными ранее ассетами.
FROM nginx:stable-alpine@sha256:d4d72ee8e6d028c5ad939454164d3645be2d38afb5c352277926a48e24abf4fa as frontend
WORKDIR /www

# [<en>] Copy the pre-built assets from the above image.
# [<ru>] Копируем собранные ассеты из предыдушего сборочного образа.
COPY --from=builder /app/dist /www/static

# [<en>] Copy the NGINX configuration.
# [<ru>] Копируем конфигурацию NGINX.
COPY .werf/nginx.conf /etc/nginx/nginx.conf
