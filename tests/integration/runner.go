//go:build integration
// +build integration

package integration

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
)

// TestRunner manages the execution of integration test suites
type TestRunner struct {
	config        *TestConfig
	logger        zerolog.Logger
	suites        map[string]TestSuiteRunner
	results       map[string][]TestResult
	overallReport *TestReport
}

// TestSuiteRunner interface for different test suites
type TestSuiteRunner interface {
	Name() string
	Description() string
	RunTests() ([]TestResult, error)
	IsEnabled(config *TestConfig) bool
}

// NewTestRunner creates a new integration test runner
func NewTestRunner(config *TestConfig) *TestRunner {
	return &TestRunner{
		config:  config,
		logger:  config.GetLogger(),
		suites:  make(map[string]TestSuiteRunner),
		results: make(map[string][]TestResult),
		overallReport: &TestReport{
			Config:          *config,
			StartTime:       time.Now(),
			TestResults:     make([]TestResult, 0),
			EnvironmentInfo: make(map[string]string),
		},
	}
}

// RegisterSuite registers a test suite with the runner
func (r *TestRunner) RegisterSuite(suite TestSuiteRunner) {
	r.suites[suite.Name()] = suite
	r.logger.Debug().Str("suite", suite.Name()).Msg("Registered test suite")
}

// RunAllSuites executes all registered and enabled test suites
func (r *TestRunner) RunAllSuites() error {
	r.logger.Info().
		Int("total_suites", len(r.suites)).
		Msg("Starting integration test execution")

	totalTests := 0
	totalPassed := 0
	totalFailed := 0

	// Execute each enabled suite
	for suiteName, suite := range r.suites {
		if !suite.IsEnabled(r.config) {
			r.logger.Info().Str("suite", suiteName).Msg("Suite disabled, skipping")
			continue
		}

		r.logger.Info().
			Str("suite", suiteName).
			Str("description", suite.Description()).
			Msg("Running test suite")

		suiteStart := time.Now()
		results, err := suite.RunTests()
		suiteDuration := time.Since(suiteStart)

		if err != nil {
			r.logger.Error().
				Err(err).
				Str("suite", suiteName).
				Msg("Suite execution failed")
			continue
		}

		r.results[suiteName] = results
		r.overallReport.TestResults = append(r.overallReport.TestResults, results...)

		// Count results
		suitePassed := 0
		suiteFailed := 0
		for _, result := range results {
			if result.Success {
				suitePassed++
			} else {
				suiteFailed++
			}
		}

		totalTests += len(results)
		totalPassed += suitePassed
		totalFailed += suiteFailed

		r.logger.Info().
			Str("suite", suiteName).
			Int("tests", len(results)).
			Int("passed", suitePassed).
			Int("failed", suiteFailed).
			Dur("duration", suiteDuration).
			Msg("Suite completed")
	}

	// Finalize overall report
	r.overallReport.EndTime = time.Now()
	r.overallReport.Duration = r.overallReport.EndTime.Sub(r.overallReport.StartTime)
	r.overallReport.TotalTests = totalTests
	r.overallReport.PassedTests = totalPassed
	r.overallReport.FailedTests = totalFailed
	r.overallReport.Summary = fmt.Sprintf(
		"Integration Tests Complete: %d passed, %d failed out of %d total across %d suites",
		totalPassed, totalFailed, totalTests, len(r.results))

	// Add environment info
	r.overallReport.EnvironmentInfo = map[string]string{
		"demo_trading":   fmt.Sprintf("%t", r.config.DemoTrading),
		"base_url":       r.config.BaseURL,
		"websocket_url":  r.config.WebSocketURL,
		"test_symbol":    r.config.TestSymbol,
		"test_coin":      r.config.TestCoin,
		"enabled_suites": fmt.Sprintf("%v", r.config.EnabledSuites),
		"max_retries":    fmt.Sprintf("%d", r.config.MaxRetries),
		"test_timeout":   fmt.Sprintf("%d", r.config.TestTimeout),
	}

	r.PrintOverallSummary()

	// Generate report if enabled
	if r.config.GenerateReport {
		if err := r.GenerateReport(); err != nil {
			r.logger.Error().Err(err).Msg("Failed to generate test report")
			return err
		}
	}

	return nil
}

// GenerateReport creates and saves the overall test report
func (r *TestRunner) GenerateReport() error {
	reportPath := r.config.ReportPath
	if reportPath == "" {
		reportPath = fmt.Sprintf("tests/reports/integration_report_%s.json",
			time.Now().Format("2006-01-02_15-04-05"))
	}

	return r.overallReport.SaveToFile(reportPath)
}

// PrintOverallSummary prints a comprehensive test execution summary
func (r *TestRunner) PrintOverallSummary() {
	r.logger.Info().Msg("=== INTEGRATION TEST SUMMARY ===")

	r.logger.Info().
		Int("total_suites", len(r.results)).
		Int("total_tests", r.overallReport.TotalTests).
		Int("passed_tests", r.overallReport.PassedTests).
		Int("failed_tests", r.overallReport.FailedTests).
		Dur("total_duration", r.overallReport.Duration).
		Msg("Overall Results")

	// Print per-suite summary
	for suiteName, results := range r.results {
		passed := 0
		failed := 0
		for _, result := range results {
			if result.Success {
				passed++
			} else {
				failed++
			}
		}

		r.logger.Info().
			Str("suite", suiteName).
			Int("tests", len(results)).
			Int("passed", passed).
			Int("failed", failed).
			Msg("Suite Summary")
	}

	// Print failed tests details
	if r.overallReport.FailedTests > 0 {
		r.logger.Warn().Msg("=== FAILED TESTS DETAILS ===")
		for _, result := range r.overallReport.TestResults {
			if !result.Success {
				r.logger.Error().
					Str("test", result.TestName).
					Str("endpoint", result.Endpoint).
					Str("error", result.Error).
					Dur("duration", result.Duration).
					Msg("Failed Test")
			}
		}
	}

	if r.overallReport.FailedTests == 0 {
		r.logger.Info().Msg("🎉 ALL TESTS PASSED! 🎉")
	} else {
		r.logger.Warn().
			Int("failed_count", r.overallReport.FailedTests).
			Msg("⚠️  Some tests failed - check configuration and API status")
	}
}

// SaveToFile saves the test report to a JSON file
func (r *TestReport) SaveToFile(filename string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create report directory: %w", err)
	}

	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal report: %w", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write report file: %w", err)
	}

	fmt.Printf("Test report saved to: %s\n", filename)
	return nil
}

// LoadFromFile loads a test report from a JSON file
func (r *TestReport) LoadFromFile(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read report file: %w", err)
	}

	if err := json.Unmarshal(data, r); err != nil {
		return fmt.Errorf("failed to unmarshal report: %w", err)
	}

	return nil
}

// GenerateHTMLReport creates an HTML version of the test report
func (r *TestReport) GenerateHTMLReport(filename string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create report directory: %w", err)
	}

	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Bitget Integration Test Report</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; background: #f5f5f5; }
        .container { max-width: 1200px; margin: 0 auto; background: white; padding: 20px; border-radius: 8px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .header { text-align: center; margin-bottom: 30px; padding: 20px; background: linear-gradient(135deg, #0066cc, #004499); color: white; border-radius: 8px; }
        .summary { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 20px; margin-bottom: 30px; }
        .metric { background: #f8f9fa; padding: 20px; border-radius: 8px; text-align: center; border-left: 4px solid #0066cc; }
        .metric.failed { border-left-color: #dc3545; }
        .metric.passed { border-left-color: #28a745; }
        .metric h3 { margin: 0; font-size: 2em; }
        .metric p { margin: 5px 0 0 0; color: #666; }
        .test-results { margin-top: 20px; }
        .test-item { background: #fff; border: 1px solid #dee2e6; border-radius: 5px; margin-bottom: 10px; padding: 15px; }
        .test-item.success { border-left: 4px solid #28a745; }
        .test-item.failed { border-left: 4px solid #dc3545; }
        .test-name { font-weight: bold; font-size: 1.1em; margin-bottom: 5px; }
        .test-endpoint { color: #666; font-family: monospace; font-size: 0.9em; }
        .test-duration { color: #666; font-size: 0.9em; }
        .test-error { color: #dc3545; font-size: 0.9em; margin-top: 5px; background: #f8f9fa; padding: 10px; border-radius: 3px; }
        .environment { background: #f8f9fa; padding: 15px; border-radius: 5px; margin-top: 20px; }
        .timestamp { color: #666; font-size: 0.9em; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🧪 Bitget Integration Test Report</h1>
            <p class="timestamp">Generated: %s</p>
            <p>Duration: %v</p>
        </div>

        <div class="summary">
            <div class="metric">
                <h3>%d</h3>
                <p>Total Tests</p>
            </div>
            <div class="metric passed">
                <h3>%d</h3>
                <p>Passed</p>
            </div>
            <div class="metric failed">
                <h3>%d</h3>
                <p>Failed</p>
            </div>
            <div class="metric">
                <h3>%.1f%%</h3>
                <p>Success Rate</p>
            </div>
        </div>

        <div class="environment">
            <h3>Environment Information</h3>
            <p><strong>Demo Trading:</strong> %s</p>
            <p><strong>Base URL:</strong> %s</p>
            <p><strong>Test Symbol:</strong> %s</p>
            <p><strong>Test Coin:</strong> %s</p>
        </div>

        <div class="test-results">
            <h3>Test Results</h3>`

	successRate := 0.0
	if r.TotalTests > 0 {
		successRate = float64(r.PassedTests) / float64(r.TotalTests) * 100
	}

	html = fmt.Sprintf(html,
		r.EndTime.Format("2006-01-02 15:04:05"),
		r.Duration,
		r.TotalTests,
		r.PassedTests,
		r.FailedTests,
		successRate,
		r.EnvironmentInfo["demo_trading"],
		r.EnvironmentInfo["base_url"],
		r.EnvironmentInfo["test_symbol"],
		r.EnvironmentInfo["test_coin"],
	)

	// Add test results
	for _, result := range r.TestResults {
		statusClass := "success"
		if !result.Success {
			statusClass = "failed"
		}

		html += fmt.Sprintf(`
            <div class="test-item %s">
                <div class="test-name">%s</div>
                <div class="test-endpoint">%s</div>
                <div class="test-duration">Duration: %v</div>`,
			statusClass, result.TestName, result.Endpoint, result.Duration)

		if !result.Success && result.Error != "" {
			html += fmt.Sprintf(`<div class="test-error">Error: %s</div>`, result.Error)
		}

		html += `</div>`
	}

	html += `
        </div>
    </div>
</body>
</html>`

	return os.WriteFile(filename, []byte(html), 0644)
}
