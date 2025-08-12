@echo off
REM Test script to verify all WebSocket examples compile correctly

echo 🧪 Testing WebSocket Examples Compilation...
echo ==============================================

set examples=basic_public_channels.go multiple_symbols.go private_channels.go advanced_usage.go mixed_channels.go
set passed=0
set failed=0

for %%e in (%examples%) do (
    echo Testing %%e...
    
    go build %%e >nul 2>&1
    if errorlevel 1 (
        echo ❌ FAIL: %%e
        echo Compilation errors:
        go build %%e
        set /a failed+=1
    ) else (
        echo ✅ PASS: %%e
        set /a passed+=1
        REM Clean up binary
        if exist "%%~ne.exe" del "%%~ne.exe"
    )
)

echo ==============================================
echo Results: %passed% passed, %failed% failed

if %failed%==0 (
    echo 🎉 All examples compile successfully!
    exit /b 0
) else (
    echo 💥 Some examples failed to compile
    exit /b 1
)