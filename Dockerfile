FROM golang:1.18 AS builder

WORKDIR /code

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
  go build -ldflags "-s -w -extldflags '-static'" -o build/ley-manager cmd/manager/main.go


FROM scratch

WORKDIR /

COPY --from=builder ./code/build/ley-manager .

EXPOSE 8080

ENTRYPOINT ["./ley-manager"]
