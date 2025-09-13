FROM node:22.15.0-alpine3.21 AS builder

WORKDIR /build

COPY package.json .

RUN npm i

COPY . .
COPY .env.docker .env

RUN npm run build

FROM nginx:1.27.5

WORKDIR /app

COPY --from=builder /build/dist/index.html /app/index.html
COPY --from=builder /build/dist/assets /app/assets
COPY --from=builder /build/dist/fonts /app/fonts
COPY --from=builder /build/dist/img /app/img

COPY ./nginx.conf /etc/nginx/conf.d/default.conf


EXPOSE 80