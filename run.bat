chcp 65001
@echo off
::: go-runapp  https://gitee.com/wanneme/go-runapp
echo "usage: run.bat start|starts|stop|restart"
::: go-runapp 文件
set MAIN_EXE=go-runapp-1.0.0-win-amd64.exe
::: 名称（请修改为jar）
set APP_NAME=NingBoCollege-0.0.1-SNAPSHOT.jar
:::set REMOTE_DEBUG=-agentlib:jdwp=transport=dt_socket,server=y,suspend=n,address=15005
set REMOTE_DEBUG=
::: 启动参数（请修改）
set START_ARGS=javaw -Xms1024m -Xmx2048m -jar %APP_NAME% -Dlogging.file=E:\GCALI\NingBoDaXue\log\out.log --server.port=8316 --spring.profiles.active=dev

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