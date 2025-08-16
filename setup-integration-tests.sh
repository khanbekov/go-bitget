#!/bin/bash

# Setup script for Bitget Go SDK Integration Tests
# This script helps developers configure integration testing environment

set -e

echo "🚀 Setting up Bitget Go SDK Integration Tests"
echo "=============================================="
echo

# Check if .env.example exists
if [ ! -f ".env.example" ]; then
    echo "❌ .env.example not found. This script must be run from the project root."
    exit 1
fi

echo "📁 Current setup status:"

# Check for existing .env files
if [ -f ".env" ]; then
    echo "  ✅ .env file exists"
    echo "  📄 Current configuration will be preserved"
else
    echo "  ⚠️  No .env file found"
fi

echo

# Ask what the user wants to do
echo "🔧 Setup options:"
echo "1. Copy .env.example to .env (for custom API keys)"
echo "2. Create .env with mock credentials (for compilation testing)"
echo "3. Show environment variable setup"
echo "4. Test current configuration"
echo "5. Exit"
echo

while true; do
    read -p "Choose option [1-5]: " choice
    case $choice in
        1)
            if [ -f ".env" ]; then
                read -p "⚠️  .env file exists. Overwrite? (y/N): " confirm
                if [[ ! "$confirm" =~ ^[Yy]$ ]]; then
                    echo "Skipping .env creation"
                    break
                fi
            fi
            cp .env.example .env
            echo "✅ Created .env from example"
            echo "📝 Please edit .env and add your API credentials"
            echo "   BITGET_API_KEY=your_actual_api_key"
            echo "   BITGET_SECRET_KEY=your_actual_secret"
            echo "   BITGET_PASSPHRASE=your_actual_passphrase"
            break
            ;;
        2)
            if [ -f ".env" ]; then
                read -p "⚠️  .env file exists. Overwrite? (y/N): " confirm
                if [[ ! "$confirm" =~ ^[Yy]$ ]]; then
                    echo "Skipping .env creation"
                    break
                fi
            fi
            cat > .env << 'EOF'
# Mock credentials for compilation testing
BITGET_API_KEY=mock_test_key
BITGET_SECRET_KEY=mock_test_secret
BITGET_PASSPHRASE=mock_test_passphrase
BITGET_DEMO_TRADING=true
BITGET_MOCK_MODE=true
EOF
            echo "✅ Created .env with mock credentials"
            echo "🧪 This setup is for compilation testing only"
            echo "   Tests will fail with 'Apikey does not exist' (expected)"
            break
            ;;
        3)
            echo
            echo "🌍 Environment Variable Setup:"
            echo "  Export these variables in your shell:"
            echo
            echo "    export BITGET_API_KEY=your_api_key"
            echo "    export BITGET_SECRET_KEY=your_secret_key"
            echo "    export BITGET_PASSPHRASE=your_passphrase"
            echo "    export BITGET_DEMO_TRADING=true"
            echo
            echo "  Or add them to your shell profile (.bashrc, .zshrc, etc.)"
            echo
            break
            ;;
        4)
            echo
            echo "🔍 Testing current configuration..."
            echo
            if [ -f "tests/scripts/run-integration-tests.sh" ]; then
                bash tests/scripts/run-integration-tests.sh --demo-check
            else
                echo "❌ Integration test script not found"
                echo "   Expected: tests/scripts/run-integration-tests.sh"
            fi
            echo
            break
            ;;
        5)
            echo "👋 Exiting setup"
            exit 0
            ;;
        *)
            echo "❌ Invalid option. Please choose 1-5."
            ;;
    esac
done

echo
echo "🎯 Next steps:"
echo "1. Run tests: bash tests/scripts/run-integration-tests.sh"
echo "2. Run in mock mode: bash tests/scripts/run-integration-tests.sh --mock-mode"
echo "3. Check demo trading: bash tests/scripts/run-integration-tests.sh --demo-check"
echo "4. See all options: bash tests/scripts/run-integration-tests.sh --help"
echo
echo "📚 For more information, see:"
echo "   - tests/README.md"
echo "   - tests/INTEGRATION_TESTING.md"
echo
echo "✨ Setup complete!"