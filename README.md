# LogSync

查询日志信息

#### 注册/卸载服务/启动/停止/重启

```shell script

LogSync.exe install 
LogSync.exe remove 
LogSync.exe start
LogSync.exe stop
LogSync.exe restart
LogSync.exe version 查看版本号

```

#### 编译

```shell script
go build -ldflags "-w -s -X main.Version=v1.9  -o ./cmd/RemoteSync.exe"
```

#### 版本更新

- v1.1 增加 pscp 输入参数 y,重启程序后不检查超时15分钟

