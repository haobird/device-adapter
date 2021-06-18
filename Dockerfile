FROM haobird/golang:develop-1.15
LABEL maintainer="xyqbgn@126.com"

# COPY . .

# RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o server main.go &&\
#   upx --best server -o _upx_server && \
#   mv -f _upx_server server
# RUN go build -ldflags "-s -w" -o server main.go

# CMD ["./server"]
CMD ["/bin/sh", "-c", "while true;do echo hello docker;sleep 1;done"]