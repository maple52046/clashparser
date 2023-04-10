# Clash proxy converter

將 Clash config 中的 proxy 列表轉換成 Shadowsocks-rust 的 server 列表。

## 為什麼要開發這個東西？

[Clash](https://github.com/Dreamacro/clash) 作為一個 proxy platform 固然有其優點；然而鄙人並不需要 Clash 的功能，只需要供應商提供的 proxy address (shadowsocks server)。因此開發這個程式，讀取 clash config，提取 shadowsocks servers，並轉換成 shadowsocks-rust config。

程式功能：

1. 從供應商的訂閱網址下載並讀取 Clash config (yaml 格式)
2. 讀取本地的 Shadowsocks config
3. 將 Clash config 的 "proxies" 轉換成 Shadowsocks config 的 "servers"
4. 輸出新的 Shadowsocks config
5. 重新啟動 Shadowsocks client 服務

## 這是誰開發的？

特別提一下，這也是我初次嘗試讓 ChatGPT (3.5) 幫我寫程式，95% 的程式碼由 [ChatGPT](https://chat.openai.com/chat) 產生；我做的事情只有：更新 package 版本、用 `flag` 管理參數、Go module、將部分程式碼改成自己喜歡的風格。

## 怎麼用？

### 編譯

```shell
make
# 或者: go build -o main main.go
```

> 如果 clash config 的網址很長，且帶有`&`符號的話，可以考慮在 build 參數中加上 `-ldflags="-X main.DefaultClashUrl=http://other.source/xxx?abc=123&efg=456"`，將網址改成供應商提供的網址，這樣可以修改程式預設的 clash 網址 (請參考 Makefile 中註解的指令。)

### 運行

```shell
./main -s shadowsocks.json -o client.json
# sslocal -c client.json
```

> 這裡的 `sslocal` 是 [shadowsocks-rust](https://github.com/shadowsocks/shadowsocks-rust) 提供的；`shadowsocks.json`、`client.json` 的結構可以參考: https://github.com/shadowsocks/shadowsocks-rust#configuration
