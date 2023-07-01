FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o /mini-project-evermos

FROM alpine:latest
COPY --from=builder /mini-project-evermos /mini-project-evermos
CMD ["/mini-project-evermos"]