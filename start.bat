@echo off


echo Starting servers...

REM Start the pocketbase serve
start cmd /k "/pb/pocketbase.exe serve"
echo Pocketbase server is running...

REM Change to the directory where the backend resides and start the server
cd UPLINK/server
start cmd /k "go run uplink.go"
echo Backend server is running...

REM Change to the directory where the frontend resides and start the server
cd ../../uplink-fe
start cmd /k "npm run dev"
echo Frontend server is running...



echo All servers are running.

