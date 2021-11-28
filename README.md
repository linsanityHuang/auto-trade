
<!-- bigone 交易所API文档 -->
https://open.bigonechina.com/docs/api.html

https://github.com/iam9k/bigone-sdk-go

https://github.com/go-resty/resty/


rm -rf one log.txt

go build -o one ./main.go

./one trade -priceSpace=0.01 -orderNum=0 &

ps -ef | grep one

<!-- 价格查询 -->
./one price -assetPair btc-usdt
./one price -assetPair ens-usdt

<!-- 交易账户余额查询 -->
./one balance -asset btc

<!-- 订单相关操作 -->
./one order -list

./one order -show=21510328802

./one order -cancel=21510328802

<!-- redis -->
docker run -p 6379:6379 -v redis-data:/data -d redis:6.2.6-alpine redis-server --save 60 1 --loglevel warning
