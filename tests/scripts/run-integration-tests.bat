@echo off
setlocal enabledelayedexpansion

REM Bitget Go SDK Integration Test Runner for Windows
REM This script provides an easy way to run integration tests against real Bitget API

REM Find project root directory
set SCRIPT_DIR=%~dp0
for %%i in ("%SCRIPT_DIR%..\..") do set PROJECT_ROOT=%%~fi

REM Change to project root directory
cd /d "%PROJECT_ROOT%"
echo Running from project root: %PROJECT_ROOT%

set CONFIG_FILE=tests\configs\integration.json
set SUITE=account
set TIMEOUT=5m
set VERBOSE=true
set GENERATE_REPORT=true

:parse_args
if "%1"=="" goto main
if "%1"=="-c" set CONFIG_FILE=%2& shift& shift& goto parse_args
if "%1"=="--config" set CONFIG_FILE=%2& shift& shift& goto parse_args
if "%1"=="-s" set SUITE=%2& shift& shift& goto parse_args
if "%1"=="--suite" set SUITE=%2& shift& shift& goto parse_args
if "%1"=="-t" set TIMEOUT=%2& shift& shift& goto parse_args
if "%1"=="--timeout" set TIMEOUT=%2& shift& shift& goto parse_args
if "%1"=="-q" set VERBOSE=false& shift& goto parse_args
if "%1"=="--quiet" set VERBOSE=false& shift& goto parse_args
if "%1"=="-r" set GENERATE_REPORT=false& shift& goto parse_args
if "%1"=="--no-report" set GENERATE_REPORT=false& shift& goto parse_args
if "%1"=="-d" goto demo_check
if "%1"=="--demo-check" goto demo_check
if "%1"=="-h" goto show_usage
if "%1"=="--help" goto show_usage
echo ERROR: Unknown option: %1
goto show_usage

:show_usage
echo Usage: %0 [OPTIONS]
echo.
echo Options:
echo   -c, --config FILE     Use specific config file ^(default: tests\configs\integration.json^)
echo   -s, --suite SUITE     Run specific test suite ^(account, market, position, all^)
echo   -t, --timeout TIME    Test timeout ^(default: 5m^)
echo   -q, --quiet           Quiet mode - minimal output
echo   -r, --no-report       Don't generate test report
echo   -d, --demo-check      Check if demo trading is enabled
echo   -h, --help            Show this help message
echo.
echo Examples:
echo   %0                                    # Run with default settings
echo   %0 -s account                        # Run only account tests
echo   %0 -c tests\configs\my_config.json   # Use custom config
echo   %0 -s account -t 10m                 # Run account tests with 10m timeout
echo   %0 -d                                # Check demo trading status
echo.
goto :eof

:print_header
echo ====================================
echo   Bitget Go SDK Integration Tests  
echo ====================================
echo.
goto :eof

:print_success
echo [92m%1[0m
goto :eof

:print_warning  
echo [93m%1[0m
goto :eof

:print_error
echo [91m%1[0m
goto :eof

:check_prerequisites
echo Checking prerequisites...

REM Check if Go is installed
go version >nul 2>&1
if errorlevel 1 (
    call :print_error "Go is not installed or not in PATH"
    exit /b 1
)
for /f "tokens=*" %%i in ('go version') do set GO_VERSION=%%i
call :print_success "Go is installed: !GO_VERSION!"

REM Check if .env file exists
if exist ".env" (
    call :print_success ".env file found"
) else (
    call :print_warning ".env file not found - make sure API credentials are set via environment variables"
)

REM Check if config file exists
if exist "%CONFIG_FILE%" (
    call :print_success "Config file found: %CONFIG_FILE%"
) else (
    call :print_warning "Config file not found: %CONFIG_FILE%"
    if exist "tests\configs\integration.example.json" (
        echo Would you like to create config from example? ^(y/n^)
        set /p response=
        if /i "!response!"=="y" (
            if not exist "tests\configs" mkdir "tests\configs"
            copy "tests\configs\integration.example.json" "%CONFIG_FILE%" >nul
            call :print_success "Created config file from example"
            call :print_warning "Please edit %CONFIG_FILE% with your settings before running tests"
            exit /b 0
        )
    )
)

REM Create reports directory if it doesn't exist
if not exist "tests\reports" mkdir "tests\reports"

echo.
goto :eof

:demo_check
call :print_header
echo Checking demo trading configuration...

if exist "%CONFIG_FILE%" (
    findstr "demo_trading.*true" "%CONFIG_FILE%" >nul 2>&1
    if !errorlevel!==0 (
        call :print_success "Demo trading is ENABLED in config"
    ) else (
        call :print_error "Demo trading is DISABLED in config"
        call :print_warning "It is STRONGLY RECOMMENDED to enable demo trading for safety"
        echo Would you like to continue anyway? ^(y/N^)
        set /p response=
        if /i not "!response!"=="y" (
            echo Exiting for safety. Please enable demo trading in %CONFIG_FILE%
            exit /b 1
        )
    )
)

if "%BITGET_DEMO_TRADING%"=="true" (
    call :print_success "Demo trading is ENABLED via environment variable"
) else if "%BITGET_DEMO_TRADING%"=="false" (
    call :print_error "Demo trading is DISABLED via environment variable"
)

echo.
goto :eof

:run_tests
set TEST_ARGS=
if "%VERBOSE%"=="false" (
    set TEST_ARGS=-q
) else (
    set TEST_ARGS=-v
)

if "%GENERATE_REPORT%"=="false" (
    set TEST_ARGS=%TEST_ARGS% -args -generate-report=false
)

echo Running integration tests...
echo Suite: %SUITE%
echo Timeout: %TIMEOUT%
echo Config: %CONFIG_FILE%
echo.

REM Convert config file to absolute path and set in environment
if "%CONFIG_FILE:~1,1%"==":" (
    REM Already absolute path
    set INTEGRATION_CONFIG_FILE=%CONFIG_FILE%
) else (
    REM Convert relative path to absolute
    set INTEGRATION_CONFIG_FILE=%PROJECT_ROOT%\%CONFIG_FILE%
)

echo Using config file: %INTEGRATION_CONFIG_FILE%
echo Current directory: %CD%

if "%SUITE%"=="account" (
    call :print_success "Running account endpoint tests..."
    go test -tags=integration ./tests/integration/suites %TEST_ARGS% -timeout=%TIMEOUT% -run TestAccountEndpoints
) else if "%SUITE%"=="all" (
    call :print_success "Running all test suites..."
    go test -tags=integration ./tests/integration/suites %TEST_ARGS% -timeout=%TIMEOUT%
) else (
    call :print_success "Running %SUITE% test suite..."
    go test -tags=integration ./tests/integration/suites %TEST_ARGS% -timeout=%TIMEOUT% -run Test%SUITE%
)

goto :eof

:main
call :print_header

echo Configuration:
echo   Config File: %CONFIG_FILE%
echo   Test Suite: %SUITE%
echo   Timeout: %TIMEOUT%
echo   Verbose: %VERBOSE%
echo   Generate Report: %GENERATE_REPORT%
echo.

call :check_prerequisites
if errorlevel 1 exit /b 1

call :demo_check
if errorlevel 1 exit /b 1

echo Starting integration tests...
echo.

call :run_tests
if errorlevel 1 (
    echo.
    call :print_error "Integration tests failed!"
    echo.
    echo Troubleshooting:
    echo 1. Check API credentials in .env file
    echo 2. Verify network connectivity
    echo 3. Ensure demo trading is enabled
    echo 4. Review configuration in %CONFIG_FILE%
    echo 5. Check logs for specific error messages
    exit /b 1
) else (
    echo.
    call :print_success "Integration tests completed!"
    
    if "%GENERATE_REPORT%"=="true" (
        if exist "tests\reports\integration_report.json" (
            call :print_success "Test report generated: tests\reports\integration_report.json"
        )
        if exist "tests\reports\integration_report.html" (
            call :print_success "HTML report generated: tests\reports\integration_report.html"
        )
    )
    
    echo.
    echo Next steps:
    echo 1. Review test results and any failed tests
    echo 2. Check generated reports for detailed analysis
    echo 3. Enable additional endpoints after verifying current tests
    echo 4. Consider running market data and position tests next
)

goto :eof