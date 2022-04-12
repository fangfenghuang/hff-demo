tag=$1

rm -f bin/go-demo
GO111MODULE="on" go build -o bin/go-demo ./main.go

if [ $? -ne 0 ]; then
exit 1
fi

if [ $tag ]
then
chmod +x go-demo
img=go-demo:$tag
echo build image: $img
docker build -t $img .
fi
