rm -rf one log.txt

go build -o one ./main.go

./one trade -priceSpace=0.01 -orderNum=0 &

ps -ef | grep one
