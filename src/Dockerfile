FROM node:22-alpine AS frontend

WORKDIR /app/web/frontend
COPY web/frontend/package.json web/frontend/vite.config.js web/frontend/index.html ./
COPY web/frontend/src ./src
RUN npm install
RUN npm run build

FROM golang:1.26-alpine AS backend

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /out/ice-web ./web/backend

FROM alpine:3.22

WORKDIR /app
COPY --from=backend /out/ice-web /app/ice-web
COPY --from=frontend /app/web/frontend/dist /app/web/frontend/dist

ENV FRONTEND_DIST=/app/web/frontend/dist
ENV PORT=8080

EXPOSE 8080
CMD ["/app/ice-web"]
