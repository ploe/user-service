FROM golang:1.22


ADD ./ /user-service
WORKDIR /user-service

RUN <<EOF
    go mod download
    go test ./...
EOF

EXPOSE 8080

CMD ["go", "run", "main.go"]