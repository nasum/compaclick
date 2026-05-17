@echo off
chcp 65001 > nul
echo Building Go program...
go build -ldflags="-H windowsgui" -o CompaClick.exe
if %ERRORLEVEL% NEQ 0 (
    echo Failed to build Go program.
    exit /b %ERRORLEVEL%
)

echo.
echo Building installer...
set ISCC_PATH="%LOCALAPPDATA%\Programs\Inno Setup 6\ISCC.exe"
if exist %ISCC_PATH% goto build_installer

set ISCC_PATH="%ProgramFiles(x86)%\Inno Setup 6\ISCC.exe"
if exist %ISCC_PATH% goto build_installer

set ISCC_PATH="%ProgramFiles%\Inno Setup 6\ISCC.exe"
if exist %ISCC_PATH% goto build_installer

echo Error: Inno Setup (ISCC.exe) not found. Please make sure Inno Setup is installed.
exit /b 1

:build_installer
%ISCC_PATH% installer.iss
if %ERRORLEVEL% NEQ 0 (
    echo Failed to build installer.
    exit /b %ERRORLEVEL%
)

echo.
echo Success! Installer has been created in the output folder.
