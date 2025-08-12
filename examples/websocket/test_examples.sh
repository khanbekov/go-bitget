#!/bin/bash

# Test script to verify all WebSocket examples compile correctly

echo "🧪 Testing WebSocket Examples Compilation..."
echo "=============================================="

examples=(
    "basic_public_channels.go"
    "multiple_symbols.go" 
    "private_channels.go"
    "advanced_usage.go"
    "mixed_channels.go"
)

failed=0
passed=0

for example in "${examples[@]}"; do
    echo -n "Testing $example... "
    
    if go build "$example" > /dev/null 2>&1; then
        echo "✅ PASS"
        ((passed++))
        # Clean up binary
        rm -f "${example%.go}"
    else
        echo "❌ FAIL"
        echo "Compilation errors:"
        go build "$example"
        ((failed++))
    fi
done

echo "=============================================="
echo "Results: $passed passed, $failed failed"

if [ $failed -eq 0 ]; then
    echo "🎉 All examples compile successfully!"
    exit 0
else
    echo "💥 Some examples failed to compile"
    exit 1
fi