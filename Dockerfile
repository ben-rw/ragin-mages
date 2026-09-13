FROM debian:stable-slim

#RUN "env GOOS=js GOARCH=wasm go build -o frontend/static/ragin-mages.wasm github.com/ben-rw/ragin-mages/cmd/game"
#RUN "go build -o ragin-mages-server ./cmd/server/"

COPY ragin-mages-server /bin/ragin-mages-server

CMD ["/bin/ragin-mages-server"]
