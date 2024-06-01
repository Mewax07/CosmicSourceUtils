@echo off

set GOOS=linux
go build -o bin/LINUX_CosmicSRCUtil .

set GOOS=windows
go build -o bin/WINDOWS_CosmicSRCUtil.exe .
