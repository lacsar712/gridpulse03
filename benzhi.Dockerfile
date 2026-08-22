FROM golang:1.22
WORKDIR /app
COPY . .
ENV CGO_ENABLED=0 GOTOOLCHAIN=local
RUN go build -o /app/server ./cmd/server
EXPOSE 8080
CMD ["/app/server","-addr","0.0.0.0:8080","-db","/tmp/app.sqlite","-data","/tmp/data"]
