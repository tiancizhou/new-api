@echo off
setlocal EnableExtensions

REM ========== Config ==========
set "REGISTRY=192.168.180.146:8082"
set "IMAGE_NAME=humen-token-gateway/new-api"
REM Optional auto-login credentials.
REM Prefer setting these as environment variables instead of saving a password here.
if not defined REGISTRY_USERNAME set "REGISTRY_USERNAME="
if not defined REGISTRY_PASSWORD set "REGISTRY_PASSWORD="
REM ============================

set "ENV=%~1"
if "%ENV%"=="" set "ENV=test"
if "%ENV%"=="test" goto :env_ok
if "%ENV%"=="prod" goto :env_ok
echo [ERROR] Unsupported env: %ENV% (only test or prod)
goto :end

:env_ok
set "IMAGE_TAG=%ENV%"
set "REMOTE_IMAGE=%REGISTRY%/%IMAGE_NAME%:%IMAGE_TAG%"

echo ========================================
echo   Env:    %ENV%
echo   Remote: %REMOTE_IMAGE%
echo ========================================

echo [1/3] Logging in to registry ...
if "%REGISTRY_USERNAME%"=="" (
    set /p "REGISTRY_USERNAME=Registry username (leave empty for Docker prompt): "
)
if "%REGISTRY_USERNAME%"=="" (
    docker login "%REGISTRY%"
) else (
    if not "%REGISTRY_PASSWORD%"=="" (
        powershell -NoProfile -ExecutionPolicy Bypass -Command "[Console]::Out.Write($env:REGISTRY_PASSWORD)" | docker login "%REGISTRY%" --username "%REGISTRY_USERNAME%" --password-stdin
    ) else (
        powershell -NoProfile -ExecutionPolicy Bypass -Command "$p = Read-Host 'Registry password' -AsSecureString; $b = [Runtime.InteropServices.Marshal]::SecureStringToBSTR($p); try { [Runtime.InteropServices.Marshal]::PtrToStringBSTR($b) } finally { [Runtime.InteropServices.Marshal]::ZeroFreeBSTR($b) }" | docker login "%REGISTRY%" --username "%REGISTRY_USERNAME%" --password-stdin
    )
)
if errorlevel 1 (
    echo [ERROR] Login failed
    goto :end
)

echo [2/3] Building image for linux/amd64 ...
docker buildx build --no-cache --platform linux/amd64 -t "%REMOTE_IMAGE%" --load .
if errorlevel 1 (
    echo [ERROR] Build failed
    goto :end
)

echo [3/3] Pushing image ...
docker push "%REMOTE_IMAGE%"
if errorlevel 1 (
    echo [ERROR] Push failed
    goto :end
)

echo ========================================
echo [DONE] %REMOTE_IMAGE%
echo ========================================

:end
pause
