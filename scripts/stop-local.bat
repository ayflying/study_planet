@echo off
rem Stop StudyPlanet local environment: server (8095) + portable MySQL (33061)
echo Stopping server on port 8095 ...
for /f "tokens=5" %%p in ('netstat -ano ^| findstr ":8095 " ^| findstr "LISTENING"') do (
    taskkill /f /pid %%p >nul 2>&1
    echo   stopped PID %%p
)
echo Stopping MySQL on port 33061 ...
for /f "tokens=5" %%p in ('netstat -ano ^| findstr ":33061 " ^| findstr "LISTENING"') do (
    taskkill /f /pid %%p >nul 2>&1
    echo   stopped PID %%p
)
echo Done.
timeout /t 2 /nobreak >nul
exit /b 0
