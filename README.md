
<!-- bigone 交易所API文档 -->
https://open.bigonechina.com/docs/api.html

https://github.com/iam9k/bigone-sdk-go

https://github.com/go-resty/resty/


rm -rf one log.txt

go build -o one ./main.go

./one trade -priceSpace=0.01 -orderNum=0 &

ps -ef | grep one
