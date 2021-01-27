# RemoteSync

查询日志信息

#### 注册/卸载服务/启动/停止/重启

```shell script
RemoteSync.exe --action=install 
RemoteSync.exe --action=remove 
RemoteSync.exe --action=start
RemoteSync.exe --action=stop
RemoteSync.exe --action=restart
RemoteSync.exe --action=version 
RemoteSync.exe --action=loc_sync  // 科室同步
RemoteSync.exe --action=remote_sync  // 探视数据同步
RemoteSync.exe --action=user_type_sync  // 职称同步
RemoteSync.exe --action=cache_clear  // 清除 token 缓存

```

#### 编译

```shell script
CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=/usr/local/bin/x86_64-w64-mingw32-gcc CXX=/usr/local/bin/x86_64-w64-mingw32-g+ go build -ldflags "-w -s -X main.Version=v1.1" -race  -o ./cmd/RemoteSync.exe main.go
```

#### 版本更新

- v1.0 完成本基础功能
- v1.1 修复并发问题
- v1.1.1 修复长时间运行token会失效问题

