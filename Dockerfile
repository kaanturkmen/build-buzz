WORKDIR /src
COPY go.mod ./
COPY . .
RUN go mod download
RUN go build -o ./server cmd/server/main.go

FROM alpine:latest
COPY ./config/email_overrides.json ./config/email_overrides.json
COPY --from=builder /src/server /server
CMD ["/server"]