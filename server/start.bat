@echo off
set DB_DSN=root:@tcp(127.0.0.1:33061)/study_planet?parseTime=true&loc=Local&charset=utf8mb4
set SERVER_PORT=8095
start /B studyplanet.exe > server.log 2>&1