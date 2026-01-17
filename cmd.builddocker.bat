@echo off
set /p VERSION=<VERSION.txt
echo Building version %VERSION%
CALL docker buildx build --platform linux/amd64,linux/arm64 -t git.cloud.zhishudali.ink/dicarne/vvorker:latest --push .