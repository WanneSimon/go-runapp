## go-runapp
一个 `windows` 启动和停止程序的 app ，解决 `windows` 下脚本启动程序后找不到 `pid` 的痛点。  

**不要执行脚本文件（`bat` `ps1` 等），脚本文件执行后即退出，没有意义（进程立即结束）**
### 使用
#### 启动  
`myapp` 自定义程序的名字  
`start` 会检查数据文件，如果存在则不会启动程序。`starts`跳过检查数据文件是否存在，直接执行程序命令
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
```
  go-runapp.exe {AppName} restart
  # 例
  # go-runapp.exe myapp restart
```