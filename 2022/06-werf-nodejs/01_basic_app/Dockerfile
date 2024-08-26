FROM node:12-alpine@sha256:d4b15b3d48f42059a15bd659be60afe21762aae9d6cbea6f124440895c27db68
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
