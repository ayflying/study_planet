@echo off
rem StudyPlanet one-click local launcher: portable MySQL (33061) + Go server (8095)
rem Safe to re-run: skips services that are already listening.
setlocal enabledelayedexpansion
set "ROOT=%~dp0.."
set "MYSQL=%USERPROFILE%\.workbuddy\binaries\mysql\mysql-8.0.29-winx64"
set "SERVER_EXE=%ROOT%\server\studyplanet_new.exe"
set "PORT=8095"
set "MYSQL_PORT=33061"

echo [1/2] Checking MySQL on port %MYSQL_PORT% ...
netstat -ano | findstr ":%MYSQL_PORT% " | findstr "LISTENING" >nul
if not errorlevel 1 goto mysql_ready

if not exist "%MYSQL%\bin\mysqld.exe" (
    echo [ERROR] Portable MySQL not found: %MYSQL%\bin\mysqld.exe
    pause
    exit /b 1
)
echo       Starting MySQL ...
start "studyplanet-mysql" /min "%MYSQL%\bin\mysqld.exe" --defaults-file="%MYSQL%\my.ini" --console

set /a tries=0
:wait_mysql
timeout /t 2 /nobreak >nul
netstat -ano | findstr ":%MYSQL_PORT% " | findstr "LISTENING" >nul
if not errorlevel 1 goto mysql_ready
set /a tries+=1
if %tries% lss 10 goto wait_mysql
echo [ERROR] MySQL start timeout. Check %MYSQL%\mysqld.log
pause
exit /b 1

:mysql_ready
echo       MySQL OK

echo [2/2] Checking server on port %PORT% ...
netstat -ano | findstr ":%PORT% " | findstr "LISTENING" >nul
if not errorlevel 1 goto server_ready

if not exist "%SERVER_EXE%" (
    echo [ERROR] Server exe not found: %SERVER_EXE%
    echo         Build it first inside server folder:
    echo         go build -o studyplanet_new.exe .
    pause
    exit /b 1
)
echo       Starting server ...
cd /d "%ROOT%\server"
set "SERVER_PORT=%PORT%"
set "DB_DSN=sp:sp123456@tcp(127.0.0.1:%MYSQL_PORT%)/study_planet?charset=utf8mb4&parseTime=true&loc=Local"
set "JWT_SECRET=local-dev-secret-0123456789abcdef"
start "studyplanet-server" "%SERVER_EXE%"
timeout /t 4 /nobreak >nul

:server_ready
echo       Server OK
echo.
echo ============================================
echo   StudyPlanet:  http://127.0.0.1:%PORT%/
echo   First time:   http://127.0.0.1:%PORT%/local-login.html
echo   Stop:         scripts\stop-local.bat
echo ============================================
echo Opening browser ...
start "" "http://127.0.0.1:%PORT%/local-login.html"
timeout /t 3 /nobreak >nul
exit /b 0
