FROM golang:1.21.5-bullseye AS build

RUN apt-get update

WORKDIR /app

COPY . .

RUN go mod download

WORKDIR /app/cmd

RUN go build -o project-service

FROM busybox:latest

WORKDIR /project-service/cmd

COPY --from=build /app/cmd/project-service .

COPY --from=build /app/.env /project-service

EXPOSE 4002

CMD ["./project-service"]