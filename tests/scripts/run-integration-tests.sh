#!/bin/bash

# Bitget Go SDK Integration Test Runner
# This script provides an easy way to run integration tests against real Bitget API

set -e

# Find project root directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

# Change to project root directory
cd "$PROJECT_ROOT"
echo "Running from project root: $PROJECT_ROOT"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default values
CONFIG_FILE="tests/configs/integration.json"
SUITE="account"
TIMEOUT="5m"
VERBOSE=true
GENERATE_REPORT=true

# Helper functions
print_header() {
    echo -e "${BLUE}====================================${NC}"
    echo -e "${BLUE}  Bitget Go SDK Integration Tests  ${NC}"
    echo -e "${BLUE}====================================${NC}"
    echo
}

print_success() {
    echo -e "${GREEN}$1${NC}"
}

print_warning() {
    echo -e "${YELLOW}$1${NC}"
}

print_error() {
    echo -e "${RED}$1${NC}"
}

show_usage() {
    echo "Usage: $0 [OPTIONS]"
    echo
    echo "Options:"
    echo "  -c, --config FILE     Use specific config file (default: tests/configs/integration.json)"
    echo "  -s, --suite SUITE     Run specific test suite (account, market, position, all)"
    echo "  -t, --timeout TIME    Test timeout (default: 5m)"
    echo "  -q, --quiet           Quiet mode - minimal output"
    echo "  -r, --no-report       Don't generate test report"
    echo "  -d, --demo-check      Check if demo trading is enabled"
    echo "  -m, --mock-mode       Enable mock mode for compilation testing"
    echo "  -h, --help            Show this help message"
    echo
    echo "API Credentials Setup:"
    echo "  1. Copy .env.example to .env and add your API keys"
    echo "  2. Export environment variables: BITGET_API_KEY, BITGET_SECRET_KEY, BITGET_PASSPHRASE"
    echo "  3. Use mock mode for compilation testing (no real API calls)"
    echo
    echo "Examples:"
    echo "  $0                                    # Run with default settings"
    echo "  $0 -s account                        # Run only account tests"
    echo "  $0 -c tests/configs/my_config.json   # Use custom config"
    echo "  $0 -s account -t 10m                 # Run account tests with 10m timeout"
    echo "  $0 -d                                # Check demo trading status"
    echo "  $0 -m                                # Run in mock mode (compilation testing)"
    echo
}

load_env_file() {
    # Load environment files in order of preference
    local env_files=(".env.local" ".env" ".env.development")
    
    for env_file in "${env_files[@]}"; do
        if [ -f "$env_file" ]; then
            print_success "$env_file file found - loading environment variables"
            # Load the .env file while preserving existing environment variables
            set -a  # automatically export all variables
            source "$env_file"
            set +a  # stop automatically exporting
            return 0
        fi
    done
    return 1
}

check_prerequisites() {
    echo "Checking prerequisites..."
    
    # Check if Go is installed
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed or not in PATH"
        exit 1
    fi
    print_success "Go is installed: $(go version)"
    
    # Load environment files
    if load_env_file; then
        # Environment file was loaded
        :
    else
        print_warning "No .env file found"
        echo "Options to set up API credentials:"
        echo "1. Copy .env.example to .env and fill in your API keys"
        echo "2. Set environment variables directly"
        echo "3. Use mock mode for compilation testing: export BITGET_MOCK_MODE=true"
        echo ""
    fi
    
    # Check environment variables
    if [ -n "$BITGET_API_KEY" ] && [ -n "$BITGET_SECRET_KEY" ] && [ -n "$BITGET_PASSPHRASE" ]; then
        print_success "API credentials found in environment variables"
    elif [ "$BITGET_MOCK_MODE" = "true" ]; then
        print_success "Mock mode enabled - using test credentials"
    else
        print_warning "API credentials not found in environment variables"
    fi
    
    # Check if config file exists
    if [ -f "$CONFIG_FILE" ]; then
        print_success "Config file found: $CONFIG_FILE"
    else
        print_warning "Config file not found: $CONFIG_FILE"
        if [ -f "tests/configs/integration.example.json" ]; then
            echo "Would you like to create config from example? (y/n)"
            read -r response
            if [[ "$response" =~ ^[Yy]$ ]]; then
                mkdir -p "$(dirname "$CONFIG_FILE")"
                cp tests/configs/integration.example.json "$CONFIG_FILE"
                print_success "Created config file from example"
                print_warning "Please edit $CONFIG_FILE with your settings before running tests"
                exit 0
            fi
        fi
    fi
    
    # Create reports directory if it doesn't exist
    mkdir -p tests/reports
    
    echo
}

check_demo_mode() {
    echo "Checking demo trading configuration..."
    
    if [ -f "$CONFIG_FILE" ]; then
        demo_trading=$(grep -o '"demo_trading":\s*true' "$CONFIG_FILE" || echo "")
        if [ -n "$demo_trading" ]; then
            print_success "Demo trading is ENABLED in config"
        else
            print_error "Demo trading is DISABLED in config"
            print_warning "It is STRONGLY RECOMMENDED to enable demo trading for safety"
            echo "Would you like to continue anyway? (y/N)"
            read -r response
            if [[ ! "$response" =~ ^[Yy]$ ]]; then
                echo "Exiting for safety. Please enable demo trading in $CONFIG_FILE"
                exit 1
            fi
        fi
    fi
    
    # Check environment variable
    if [ "$BITGET_DEMO_TRADING" = "true" ]; then
        print_success "Demo trading is ENABLED via environment variable"
    elif [ "$BITGET_DEMO_TRADING" = "false" ]; then
        print_error "Demo trading is DISABLED via environment variable"
    fi
    
    echo
}

run_tests() {
    local test_args=""
    
    if [ "$VERBOSE" = false ]; then
        test_args=""  # No verbose flag for quiet mode
    else
        test_args="-v"  # Verbose mode
    fi
    
    if [ "$GENERATE_REPORT" = false ]; then
        test_args="$test_args -args -generate-report=false"
    fi
    
    echo "Running integration tests..."
    echo "Suite: $SUITE"
    echo "Timeout: $TIMEOUT"
    echo "Config: $CONFIG_FILE"
    echo
    
    # Convert config file to absolute path and set in environment
    if [[ "$CONFIG_FILE" = /* ]]; then
        # Already absolute path
        export INTEGRATION_CONFIG_FILE="$CONFIG_FILE"
    else
        # Convert relative path to absolute
        export INTEGRATION_CONFIG_FILE="$PROJECT_ROOT/$CONFIG_FILE"
    fi
    
    echo "Using config file: $INTEGRATION_CONFIG_FILE"
    echo "Current directory: $(pwd)"
    
    case $SUITE in
        "account")
            print_success "Running account endpoint tests..."
            go test -tags=integration ./tests/integration/suites $test_args -timeout="$TIMEOUT" -run TestAccountEndpoints
            ;;
        "all")
            print_success "Running all test suites..."
            go test -tags=integration ./tests/integration/suites $test_args -timeout="$TIMEOUT"
            ;;
        *)
            print_success "Running $SUITE test suite..."
            go test -tags=integration ./tests/integration/suites $test_args -timeout="$TIMEOUT" -run "Test${SUITE^}"
            ;;
    esac
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -c|--config)
            CONFIG_FILE="$2"
            shift 2
            ;;
        -s|--suite)
            SUITE="$2"
            shift 2
            ;;
        -t|--timeout)
            TIMEOUT="$2"
            shift 2
            ;;
        -q|--quiet)
            VERBOSE=false
            shift
            ;;
        -r|--no-report)
            GENERATE_REPORT=false
            shift
            ;;
        -d|--demo-check)
            print_header
            check_demo_mode
            exit 0
            ;;
        -m|--mock-mode)
            export BITGET_MOCK_MODE=true
            print_success "Mock mode enabled"
            shift
            ;;
        -h|--help)
            show_usage
            exit 0
            ;;
        *)
            print_error "Unknown option: $1"
            show_usage
            exit 1
            ;;
    esac
done

# Main execution
print_header

echo "Configuration:"
echo "  Config File: $CONFIG_FILE"
echo "  Test Suite: $SUITE"
echo "  Timeout: $TIMEOUT"
echo "  Verbose: $VERBOSE"
echo "  Generate Report: $GENERATE_REPORT"
echo

check_prerequisites
check_demo_mode

echo "Starting integration tests..."
echo

if run_tests; then
    echo
    print_success "Integration tests completed!"
    
    if [ "$GENERATE_REPORT" = true ]; then
        if [ -f "tests/reports/integration_report.json" ]; then
            print_success "Test report generated: tests/reports/integration_report.json"
        fi
        if [ -f "tests/reports/integration_report.html" ]; then
            print_success "HTML report generated: tests/reports/integration_report.html"
        fi
    fi
    
    echo
    echo "Next steps:"
    echo "1. Review test results and any failed tests"
    echo "2. Check generated reports for detailed analysis"
    echo "3. Enable additional endpoints after verifying current tests"
    echo "4. Consider running market data and position tests next"
    
else
    echo
    print_error "Integration tests failed!"
    echo
    echo "Troubleshooting:"
    echo "1. Check API credentials in .env file"
    echo "2. Verify network connectivity"
    echo "3. Ensure demo trading is enabled"
    echo "4. Review configuration in $CONFIG_FILE"
    echo "5. Check logs for specific error messages"
    
    exit 1
fi