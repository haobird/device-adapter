FROM golang:1.16-stretch
ARG SVC
ARG GOARCH
ARG GOARM

WORKDIR /opt
COPY . .

# RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o server main.go &&\
#   upx --best server -o _upx_server && \
#   mv -f _upx_server server
RUN go build -mod=vendor -ldflags "-s -w" -o server boot/$SVC/main.go

CMD ["./server"]
# CMD ["/bin/sh", "-c", "while true;do echo hello docker;sleep 1;done"]