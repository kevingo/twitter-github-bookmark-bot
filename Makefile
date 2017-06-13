BINARY=ptt
GOOS=linux
GOARCH=amd64

build:
	GOOS=${GOOS} GOARCH=${GOARCH} CGO_ENABLED=0 go build -o server *.go
deploy:
	scp -i ~/.ssh/google_compute_engine ./server mtk11018@35.197.28.88: