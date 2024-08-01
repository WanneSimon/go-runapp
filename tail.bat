@echo off
:::echo "usage: tail.bat {file} [number]"

set FILE=%1
if "%FILE%" == "" (
  echo usage: tail.bat {file} [number]
  goto EOF
)

set SIZE=10
if "%2" neq "" (
  set SIZE=%2
)

powershell Get-Content -Path %FILE% -tail %SIZE% -wait

:EOF