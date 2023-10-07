FROM node:18-alpine AS js-builder

WORKDIR /src

COPY package.json .
COPY package-lock.json .

RUN npm install

COPY . .

RUN npm run build

FROM golang:1.20-bullseye AS bin-builder

WORKDIR /src

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o bin/app

FROM gcr.io/distroless/static-debian11:nonroot

WORKDIR /my

COPY --from=js-builder /src/dist dist
COPY --from=bin-builder /src/bin/app app
COPY public public

CMD ["/my/app"]