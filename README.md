## go-runapp
一个 `windows` 启动和停止程序的 app ，解决 `windows` 下脚本启动程序后找不到 `pid` 的痛点。  
`linux` 上使用也可以，没测试！

**报毒是正常现象，因为是无窗口程序**
**不要执行脚本文件（`bat` `ps1` 等），脚本文件执行后即退出，没有意义（进程立即结束）**
### 使用
#### 启动  
`myapp` 自定义程序的名字  
`start` 会检查数据文件，如果存在则不会启动程序。`starts`跳过检查数据文件是否存在，直接执行程序命令，建议使用 `starts`
```
  go-runapp.exe {AppName} start|starts [arg1] [arg2] [arg3]...
  # 例
  # go-runapp.exe myapp start javaw -jar myapp.jar
  # go-runapp.exe myapp starts javaw -jar myapp.jar 
```
### 停止
```
  go-runapp.exe {AppName} stop
  # 例
  # go-runapp.exe myapp stop
```
### 重启
`restart` 可以和 `start`|`starts` 一样后接启动命令。  
如果有启动命令，那么停止和启动时以参数为准。（推荐带启动参数）  
如果没有则已以pid文件为准。
```
  go-runapp.exe {AppName} restart [[arg1] [arg2] [arg3]]...
  # 例
  # go-runapp.exe myapp restart javaw -jar myapp.jar 
  # go-runapp.exe myapp restart 
```

### 构建
```
go build -ldflags "-s -w -H=windowsgui"
```

### 附带bat启动脚本
`run.bat`
``` bat
chcp 65001
@echo off
echo "usage: run.bat start|starts|stop|restart"
::: go-runapp 文件
set MAIN_EXE=go-runapp-1.0.0-win-amd64.exe
::: 名称（请修改为jar）
set APP_NAME=NingBoCollege-0.0.1-SNAPSHOT.jar
:::set REMOTE_DEBUG=-agentlib:jdwp=transport=dt_socket,server=y,suspend=n,address=15005
set REMOTE_DEBUG=
::: 启动参数（请修改）
set START_ARGS=javaw -Xms1024m -Xmx2048m %REMOTE_DEBUG% -jar %APP_NAME% --server.port=8316 --spring.profiles.active=dev

set RUN_EXEC="%MAIN_EXE%" "%APP_NAME%"

if "%1" == "stop" (
  set RUN_EXEC=%RUN_EXEC% %1
) else (
  set RUN_EXEC=%RUN_EXEC% %1 %START_ARGS%
)

::: 执行
:execStart
echo %RUN_EXEC%

%RUN_EXEC%

:END
```
