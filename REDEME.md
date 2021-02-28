

https://github.com/iam9k/bigone-sdk-go

https://github.com/go-resty/resty/

### 网格交易策略原理介绍

#### 名词解释

最低价：最低价（Pl）是执行网格交易的最低买入价，低于该价格时就不会再继续买入。

最高价：最高价（Ph）是执行网格交易的最高卖出价，高于该价格时就不会再继续卖出。

用户需注意，最高价需大于等于最低价的1.1倍。

进场价：进场价格（P0）是您在创建网格时所选币种的实时价格。

价格间距：Bibox网格交易工具为等比网格。价格间距（width）为网格间的距离，取值在0.2%-20%之间。

每格买入卖出金额：系统会按照当前市场价格把用户投入的资产均分在所有网格上。


go build -o auto-trade ./main.go

nohup ./auto-trade &

./auto-trade &

ps -ef | grep auto-trade
