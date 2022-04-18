

image=hff/webhook-demo:$1
echo build $image

rm -f bin/webhook-demo
GO111MODULE="on" go build -o bin/webhook-demo ./cmd/main.go


docker build -t $image .