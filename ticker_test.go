package yfinance_api

import (
	"testing"
	"time"
)

// TestNewClient tests the creation of a new client
func TestNewClient(t *testing.T) {
	client := NewClient()
	if client == nil {
		t.Fatal("NewClient() returned nil")
	}

	if client.Client == nil {
		t.Fatal("NewClient() returned YFinanceAPI with nil Client")
	}

	// Test that multiple calls return different YFinanceAPI instances but same underlying Client (singleton pattern)
	client2 := NewClient()
	if client == client2 {
		t.Error("NewClient() should return different YFinanceAPI instances")
	}

	// But the underlying Client should be the same (singleton)
	if client.Client != client2.Client {
		t.Error("NewClient() should return YFinanceAPI instances with the same underlying Client (singleton)")
	}
}

// TestNewTicker tests the creation of a new ticker
func TestNewTicker(t *testing.T) {
	symbol := "AAPL"
	ticker := NewTicker(symbol)

	if ticker == nil {
		t.Fatal("NewTicker() returned nil")
	}

	if ticker.Symbol != symbol {
		t.Errorf("Expected symbol %s, got %s", symbol, ticker.Symbol)
	}

	if ticker.Client == nil {
		t.Error("Ticker client should not be nil")
	}
}

// TestTickerGetSymbol tests the GetSymbol method
func TestTickerGetSymbol(t *testing.T) {
	symbol := "MSFT"
	ticker := NewTicker(symbol)

	if ticker.GetSymbol() != symbol {
		t.Errorf("Expected symbol %s, got %s", symbol, ticker.GetSymbol())
	}
}

// TestTickerSetSymbol tests the SetSymbol method
func TestTickerSetSymbol(t *testing.T) {
	ticker := NewTicker("AAPL")
	newSymbol := "GOOGL"

	ticker.SetSymbol(newSymbol)

	if ticker.GetSymbol() != newSymbol {
		t.Errorf("Expected symbol %s after SetSymbol, got %s", newSymbol, ticker.GetSymbol())
	}
}

// TestFetchInformation tests fetching ticker information
func TestFetchInformation(t *testing.T) {
	ticker := NewTicker("AAPL")

	info, err := ticker.FetchInformation()
	if err != nil {
		t.Skipf("Skipping test due to API error: %v", err)
		return
	}

	if info.Symbol == "" {
		t.Error("Expected non-empty symbol in ticker info")
	}

	if info.RegularMarketPrice == nil {
		t.Error("Expected regular market price to be available")
	}
}

// TestFetchPriceValue tests fetching the current price
func TestFetchPriceValue(t *testing.T) {
	ticker := NewTicker("AAPL")

	price, err := ticker.FetchPriceValue()
	if err != nil {
		t.Skipf("Skipping test due to API error: %v", err)
		return
	}

	if price.Raw <= 0 {
		t.Error("Expected positive price value")
	}

	if price.Fmt == "" {
		t.Error("Expected non-empty formatted price")
	}
}

// TestFetchHistoricalData tests fetching historical data with different parameters
func TestFetchHistoricalData(t *testing.T) {
	ticker := NewTicker("AAPL")

	testCases := []struct {
		name       string
		rangeParam string
		interval   string
		period1    string
		period2    string
		expectData bool
	}{
		{
			name:       "Default parameters",
			rangeParam: "",
			interval:   "",
			period1:    "",
			period2:    "",
			expectData: true,
		},
		{
			name:       "1 month daily",
			rangeParam: "1mo",
			interval:   "1d",
			period1:    "",
			period2:    "",
			expectData: true,
		},
		{
			name:       "1 day 5 minute intervals",
			rangeParam: "1d",
			interval:   "5m",
			period1:    "",
			period2:    "",
			expectData: true,
		},
		{
			name:       "5 days daily",
			rangeParam: "5d",
			interval:   "1d",
			period1:    "",
			period2:    "",
			expectData: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data, err := ticker.FetchHistoricalData(tc.rangeParam, tc.interval, tc.period1, tc.period2)
			if err != nil {
				t.Skipf("Skipping test due to API error: %v", err)
				return
			}

			if tc.expectData && len(data) == 0 {
				t.Error("Expected historical data but got empty result")
			}

			// Check data structure
			for date, priceData := range data {
				if date == "" {
					t.Error("Expected non-empty date key")
				}

				// Check that we have some price data (at least one field should be non-nil)
				if priceData.Open == nil && priceData.High == nil &&
					priceData.Low == nil && priceData.Close == nil {
					t.Error("Expected at least one price field to be non-nil")
				}

				// Only check first few entries to avoid long test runs
				break
			}
		})
	}
}

// TestFetchNews tests fetching news articles
func TestFetchNews(t *testing.T) {
	ticker := NewTicker("AAPL")

	testCases := []struct {
		name  string
		count int
		start int
	}{
		{
			name:  "Default count",
			count: 0,
			start: 0,
		},
		{
			name:  "5 articles",
			count: 5,
			start: 0,
		},
		{
			name:  "Pagination",
			count: 3,
			start: 2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			news, err := ticker.FetchNews(tc.count, tc.start)
			if err != nil {
				t.Skipf("Skipping test due to API error: %v", err)
				return
			}

			// News might be empty, which is okay
			for _, article := range news {
				if article.UUID == "" && article.Title == "" {
					t.Error("Expected news article to have UUID or Title")
				}

				if article.Link == "" {
					t.Error("Expected news article to have a link")
				}

				// Only check first article to avoid long test runs
				break
			}
		})
	}
}

// TestFetchNewsAlternative tests the alternative news fetching method
func TestFetchNewsAlternative(t *testing.T) {
	ticker := NewTicker("AAPL")

	news, err := ticker.FetchNewsAlternative()
	if err != nil {
		t.Skipf("Skipping test due to API error: %v", err)
		return
	}

	// This method currently returns empty slice, so just check it's not nil
	if news == nil {
		t.Error("Expected non-nil news slice")
	}
}

// TestTransformHistoricalData tests the historical data transformation
func TestTransformHistoricalData(t *testing.T) {
	// Create mock data
	mockResponse := YahooHistoryResponse{
		Chart: struct {
			Result []struct {
				Meta struct {
					Currency             string  `json:"currency"`
					Symbol               string  `json:"symbol"`
					ExchangeName         string  `json:"exchangeName"`
					InstrumentType       string  `json:"instrumentType"`
					FirstTradeDate       int64   `json:"firstTradeDate"`
					RegularMarketTime    int64   `json:"regularMarketTime"`
					Gmtoffset            int     `json:"gmtoffset"`
					Timezone             string  `json:"timezone"`
					ExchangeTimezoneName string  `json:"exchangeTimezoneName"`
					RegularMarketPrice   float64 `json:"regularMarketPrice"`
					ChartPreviousClose   float64 `json:"chartPreviousClose"`
					PriceHint            int     `json:"priceHint"`
					CurrentTradingPeriod struct {
						Pre struct {
							Timezone  string `json:"timezone"`
							Start     int64  `json:"start"`
							End       int64  `json:"end"`
							Gmtoffset int    `json:"gmtoffset"`
						} `json:"pre"`
						Regular struct {
							Timezone  string `json:"timezone"`
							Start     int64  `json:"start"`
							End       int64  `json:"end"`
							Gmtoffset int    `json:"gmtoffset"`
						} `json:"regular"`
						Post struct {
							Timezone  string `json:"timezone"`
							Start     int64  `json:"start"`
							End       int64  `json:"end"`
							Gmtoffset int    `json:"gmtoffset"`
						} `json:"post"`
					} `json:"currentTradingPeriod"`
					DataGranularity string   `json:"dataGranularity"`
					Range           string   `json:"range"`
					ValidRanges     []string `json:"validRanges"`
				} `json:"meta"`
				Timestamp  []int64 `json:"timestamp"`
				Indicators struct {
					Quote []struct {
						Open   []*float64 `json:"open"`
						High   []*float64 `json:"high"`
						Low    []*float64 `json:"low"`
						Close  []*float64 `json:"close"`
						Volume []*int64   `json:"volume"`
					} `json:"quote"`
				} `json:"indicators"`
			} `json:"result"`
			Error interface{} `json:"error"`
		}{
			Result: []struct {
				Meta struct {
					Currency             string  `json:"currency"`
					Symbol               string  `json:"symbol"`
					ExchangeName         string  `json:"exchangeName"`
					InstrumentType       string  `json:"instrumentType"`
					FirstTradeDate       int64   `json:"firstTradeDate"`
					RegularMarketTime    int64   `json:"regularMarketTime"`
					Gmtoffset            int     `json:"gmtoffset"`
					Timezone             string  `json:"timezone"`
					ExchangeTimezoneName string  `json:"exchangeTimezoneName"`
					RegularMarketPrice   float64 `json:"regularMarketPrice"`
					ChartPreviousClose   float64 `json:"chartPreviousClose"`
					PriceHint            int     `json:"priceHint"`
					CurrentTradingPeriod struct {
						Pre struct {
							Timezone  string `json:"timezone"`
							Start     int64  `json:"start"`
							End       int64  `json:"end"`
							Gmtoffset int    `json:"gmtoffset"`
						} `json:"pre"`
						Regular struct {
							Timezone  string `json:"timezone"`
							Start     int64  `json:"start"`
							End       int64  `json:"end"`
							Gmtoffset int    `json:"gmtoffset"`
						} `json:"regular"`
						Post struct {
							Timezone  string `json:"timezone"`
							Start     int64  `json:"start"`
							End       int64  `json:"end"`
							Gmtoffset int    `json:"gmtoffset"`
						} `json:"post"`
					} `json:"currentTradingPeriod"`
					DataGranularity string   `json:"dataGranularity"`
					Range           string   `json:"range"`
					ValidRanges     []string `json:"validRanges"`
				} `json:"meta"`
				Timestamp  []int64 `json:"timestamp"`
				Indicators struct {
					Quote []struct {
						Open   []*float64 `json:"open"`
						High   []*float64 `json:"high"`
						Low    []*float64 `json:"low"`
						Close  []*float64 `json:"close"`
						Volume []*int64   `json:"volume"`
					} `json:"quote"`
				} `json:"indicators"`
			}{
				{
					Timestamp: []int64{1640995200, 1641081600}, // Two timestamps
					Indicators: struct {
						Quote []struct {
							Open   []*float64 `json:"open"`
							High   []*float64 `json:"high"`
							Low    []*float64 `json:"low"`
							Close  []*float64 `json:"close"`
							Volume []*int64   `json:"volume"`
						} `json:"quote"`
					}{
						Quote: []struct {
							Open   []*float64 `json:"open"`
							High   []*float64 `json:"high"`
							Low    []*float64 `json:"low"`
							Close  []*float64 `json:"close"`
							Volume []*int64   `json:"volume"`
						}{
							{
								Open:   []*float64{floatPtr(150.0), floatPtr(151.0)},
								High:   []*float64{floatPtr(155.0), floatPtr(156.0)},
								Low:    []*float64{floatPtr(149.0), floatPtr(150.0)},
								Close:  []*float64{floatPtr(154.0), floatPtr(155.0)},
								Volume: []*int64{int64Ptr(1000000), int64Ptr(1100000)},
							},
						},
					},
				},
			},
		},
	}

	// Test daily interval (should format as date only)
	result := transformHistoricalData(mockResponse, "1d")

	if len(result) != 2 {
		t.Errorf("Expected 2 data points, got %d", len(result))
	}

	// Check date formatting for daily interval
	for date := range result {
		if len(date) != 10 { // Format: 2006-01-02
			t.Errorf("Expected date format YYYY-MM-DD, got %s", date)
		}
	}

	// Test minute interval (should format as datetime)
	result = transformHistoricalData(mockResponse, "5m")

	for datetime := range result {
		if len(datetime) != 19 { // Format: 2006-01-02 15:04:05
			t.Errorf("Expected datetime format YYYY-MM-DD HH:MM:SS, got %s", datetime)
		}
	}
}

// TestEmptyHistoricalData tests transformation with empty data
func TestEmptyHistoricalData(t *testing.T) {
	emptyResponse := YahooHistoryResponse{}

	result := transformHistoricalData(emptyResponse, "1d")

	if len(result) != 0 {
		t.Errorf("Expected empty result for empty response, got %d items", len(result))
	}
}

// Helper functions for creating pointers
func floatPtr(f float64) *float64 {
	return &f
}

func int64Ptr(i int64) *int64 {
	return &i
}

// Benchmark tests
func BenchmarkNewClient(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewClient()
	}
}

func BenchmarkNewTicker(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewTicker("AAPL")
	}
}

func BenchmarkTransformHistoricalData(b *testing.B) {
	// Create mock response with more data points
	timestamps := make([]int64, 100)
	opens := make([]*float64, 100)
	highs := make([]*float64, 100)
	lows := make([]*float64, 100)
	closes := make([]*float64, 100)
	volumes := make([]*int64, 100)

	baseTime := time.Now().Unix()
	for i := 0; i < 100; i++ {
		timestamps[i] = baseTime + int64(i*3600) // Every hour
		opens[i] = floatPtr(150.0 + float64(i))
		highs[i] = floatPtr(155.0 + float64(i))
		lows[i] = floatPtr(149.0 + float64(i))
		closes[i] = floatPtr(154.0 + float64(i))
		volumes[i] = int64Ptr(1000000 + int64(i*1000))
	}

	mockResponse := YahooHistoryResponse{
		Chart: struct {
			Result []struct {
				Meta struct {
					Currency             string  `json:"currency"`
					Symbol               string  `json:"symbol"`
					ExchangeName         string  `json:"exchangeName"`
					InstrumentType       string  `json:"instrumentType"`
					FirstTradeDate       int64   `json:"firstTradeDate"`
					RegularMarketTime    int64   `json:"regularMarketTime"`
					Gmtoffset            int     `json:"gmtoffset"`
					Timezone             string  `json:"timezone"`
					ExchangeTimezoneName string  `json:"exchangeTimezoneName"`
					RegularMarketPrice   float64 `json:"regularMarketPrice"`
					ChartPreviousClose   float64 `json:"chartPreviousClose"`
					PriceHint            int     `json:"priceHint"`
					CurrentTradingPeriod struct {
						Pre struct {
							Timezone  string `json:"timezone"`
							Start     int64  `json:"start"`
							End       int64  `json:"end"`
							Gmtoffset int    `json:"gmtoffset"`
						} `json:"pre"`
						Regular struct {
							Timezone  string `json:"timezone"`
							Start     int64  `json:"start"`
							End       int64  `json:"end"`
							Gmtoffset int    `json:"gmtoffset"`
						} `json:"regular"`
						Post struct {
							Timezone  string `json:"timezone"`
							Start     int64  `json:"start"`
							End       int64  `json:"end"`
							Gmtoffset int    `json:"gmtoffset"`
						} `json:"post"`
					} `json:"currentTradingPeriod"`
					DataGranularity string   `json:"dataGranularity"`
					Range           string   `json:"range"`
					ValidRanges     []string `json:"validRanges"`
				} `json:"meta"`
				Timestamp  []int64 `json:"timestamp"`
				Indicators struct {
					Quote []struct {
						Open   []*float64 `json:"open"`
						High   []*float64 `json:"high"`
						Low    []*float64 `json:"low"`
						Close  []*float64 `json:"close"`
						Volume []*int64   `json:"volume"`
					} `json:"quote"`
				} `json:"indicators"`
			} `json:"result"`
			Error interface{} `json:"error"`
		}{
			Result: []struct {
				Meta struct {
					Currency             string  `json:"currency"`
					Symbol               string  `json:"symbol"`
					ExchangeName         string  `json:"exchangeName"`
					InstrumentType       string  `json:"instrumentType"`
					FirstTradeDate       int64   `json:"firstTradeDate"`
					RegularMarketTime    int64   `json:"regularMarketTime"`
					Gmtoffset            int     `json:"gmtoffset"`
					Timezone             string  `json:"timezone"`
					ExchangeTimezoneName string  `json:"exchangeTimezoneName"`
					RegularMarketPrice   float64 `json:"regularMarketPrice"`
					ChartPreviousClose   float64 `json:"chartPreviousClose"`
					PriceHint            int     `json:"priceHint"`
					CurrentTradingPeriod struct {
						Pre struct {
							Timezone  string `json:"timezone"`
							Start     int64  `json:"start"`
							End       int64  `json:"end"`
							Gmtoffset int    `json:"gmtoffset"`
						} `json:"pre"`
						Regular struct {
							Timezone  string `json:"timezone"`
							Start     int64  `json:"start"`
							End       int64  `json:"end"`
							Gmtoffset int    `json:"gmtoffset"`
						} `json:"regular"`
						Post struct {
							Timezone  string `json:"timezone"`
							Start     int64  `json:"start"`
							End       int64  `json:"end"`
							Gmtoffset int    `json:"gmtoffset"`
						} `json:"post"`
					} `json:"currentTradingPeriod"`
					DataGranularity string   `json:"dataGranularity"`
					Range           string   `json:"range"`
					ValidRanges     []string `json:"validRanges"`
				} `json:"meta"`
				Timestamp  []int64 `json:"timestamp"`
				Indicators struct {
					Quote []struct {
						Open   []*float64 `json:"open"`
						High   []*float64 `json:"high"`
						Low    []*float64 `json:"low"`
						Close  []*float64 `json:"close"`
						Volume []*int64   `json:"volume"`
					} `json:"quote"`
				} `json:"indicators"`
			}{
				{
					Timestamp: timestamps,
					Indicators: struct {
						Quote []struct {
							Open   []*float64 `json:"open"`
							High   []*float64 `json:"high"`
							Low    []*float64 `json:"low"`
							Close  []*float64 `json:"close"`
							Volume []*int64   `json:"volume"`
						} `json:"quote"`
					}{
						Quote: []struct {
							Open   []*float64 `json:"open"`
							High   []*float64 `json:"high"`
							Low    []*float64 `json:"low"`
							Close  []*float64 `json:"close"`
							Volume []*int64   `json:"volume"`
						}{
							{
								Open:   opens,
								High:   highs,
								Low:    lows,
								Close:  closes,
								Volume: volumes,
							},
						},
					},
				},
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		transformHistoricalData(mockResponse, "1d")
	}
}

// TestFetchFinancialData tests fetching comprehensive financial data
func TestFetchFinancialData(t *testing.T) {
	ticker := NewTicker("AAPL")

	data, err := ticker.FetchFinancialData()
	if err != nil {
		t.Skipf("Skipping test due to API error: %v", err)
		return
	}

	// Check that we have some financial data
	if data.Summary.MarketCap == nil && data.Summary.TrailingPE == nil {
		t.Error("Expected some financial summary data")
	}

	// Test ratios structure
	if data.Ratios.PriceToEarningsRatio == nil && data.Ratios.PriceToBookRatio == nil {
		t.Log("No P/E or P/B ratio available - this is normal for some stocks")
	}
}

// TestFetchFinancialRatios tests fetching financial ratios
func TestFetchFinancialRatios(t *testing.T) {
	ticker := NewTicker("AAPL")

	ratios, err := ticker.FetchFinancialRatios()
	if err != nil {
		t.Skipf("Skipping test due to API error: %v", err)
		return
	}

	// Check for at least one ratio being available
	hasRatio := ratios.PriceToEarningsRatio != nil ||
		ratios.PriceToBookRatio != nil ||
		ratios.PriceToSalesRatio != nil ||
		ratios.DividendYield != nil

	if !hasRatio {
		t.Log("No financial ratios available - this might be normal for some stocks")
	}
}

// TestFetchKeyStatistics tests fetching key financial statistics
func TestFetchKeyStatistics(t *testing.T) {
	ticker := NewTicker("AAPL")

	stats, err := ticker.FetchKeyStatistics()
	if err != nil {
		t.Skipf("Skipping test due to API error: %v", err)
		return
	}

	// Check for market cap (most stocks should have this)
	if stats.MarketCap == nil {
		t.Log("No market cap available")
	} else if stats.MarketCap.Raw <= 0 {
		t.Error("Expected positive market cap value")
	}
}

// TestFetchIncomeStatement tests fetching income statement data
func TestFetchIncomeStatement(t *testing.T) {
	ticker := NewTicker("AAPL")

	income, err := ticker.FetchIncomeStatement()
	if err != nil {
		t.Skipf("Skipping test due to API error: %v", err)
		return
	}

	// Check for revenue (most companies should have this)
	if income.TotalRevenue != nil && income.TotalRevenue.Raw <= 0 {
		t.Error("Expected positive total revenue")
	}
}

// TestFetchBalanceSheet tests fetching balance sheet data
func TestFetchBalanceSheet(t *testing.T) {
	ticker := NewTicker("AAPL")

	balance, err := ticker.FetchBalanceSheet()
	if err != nil {
		t.Skipf("Skipping test due to API error: %v", err)
		return
	}

	// Check for total assets
	if balance.TotalAssets != nil && balance.TotalAssets.Raw <= 0 {
		t.Error("Expected positive total assets")
	}
}

// TestFetchCashFlow tests fetching cash flow statement data
func TestFetchCashFlow(t *testing.T) {
	ticker := NewTicker("AAPL")

	cashflow, err := ticker.FetchCashFlow()
	if err != nil {
		t.Skipf("Skipping test due to API error: %v", err)
		return
	}

	// Check for operating cash flow
	if cashflow.OperatingCashFlow != nil {
		t.Logf("Operating Cash Flow: %s", cashflow.OperatingCashFlow.Fmt)
	}
}

// TestFinancialDataTypes tests the financial data type structures
func TestFinancialDataTypes(t *testing.T) {
	// Test PriceValue creation
	price := &PriceValue{
		Raw: 150.25,
		Fmt: "$150.25",
	}

	if price.Raw != 150.25 {
		t.Errorf("Expected Raw value 150.25, got %f", price.Raw)
	}

	if price.Fmt != "$150.25" {
		t.Errorf("Expected Fmt value '$150.25', got %s", price.Fmt)
	}

	// Test FinancialRatios structure
	ratios := FinancialRatios{
		PriceToEarningsRatio: price,
		DividendYield:        &PriceValue{Raw: 0.015, Fmt: "1.50%"},
	}

	if ratios.PriceToEarningsRatio.Raw != 150.25 {
		t.Error("Financial ratios structure not working correctly")
	}
}

// TestMultipleStockTypes tests financial data for different stock types
func TestMultipleStockTypes(t *testing.T) {
	testCases := []struct {
		symbol      string
		description string
	}{
		{"AAPL", "Large cap tech stock"},
		{"MSFT", "Large cap tech stock"},
		{"BRK-A", "Conglomerate with high stock price"},
		{"TSLA", "Growth stock"},
	}

	for _, tc := range testCases {
		t.Run(tc.symbol, func(t *testing.T) {
			ticker := NewTicker(tc.symbol)

			// Test key statistics
			stats, err := ticker.FetchKeyStatistics()
			if err != nil {
				t.Skipf("Skipping %s due to API error: %v", tc.symbol, err)
				return
			}

			t.Logf("%s (%s) - Market Cap available: %v",
				tc.symbol, tc.description, stats.MarketCap != nil)
		})
	}
}

// TestFetchDividendInfo tests fetching comprehensive dividend information
func TestFetchDividendInfo(t *testing.T) {
	testCases := []struct {
		symbol      string
		description string
	}{
		{"AAPL", "Large cap with dividend"},
		{"MSFT", "Tech stock with dividend"},
		{"JNJ", "Dividend aristocrat"},
		{"KO", "High dividend yield stock"},
	}

	for _, tc := range testCases {
		t.Run(tc.symbol, func(t *testing.T) {
			ticker := NewTicker(tc.symbol)

			dividendInfo, err := ticker.FetchDividendInfo()
			if err != nil {
				t.Skipf("Skipping %s due to API error: %v", tc.symbol, err)
				return
			}

			// Log dividend information for manual verification
			if dividendInfo.DividendRate != nil {
				t.Logf("%s - Dividend Rate: %s", tc.symbol, dividendInfo.DividendRate.Fmt)
			}
			if dividendInfo.DividendYield != nil {
				t.Logf("%s - Dividend Yield: %s", tc.symbol, dividendInfo.DividendYield.Fmt)
			}

			// Basic validation
			if dividendInfo.DividendRate != nil && dividendInfo.DividendRate.Raw < 0 {
				t.Error("Dividend rate should not be negative")
			}
			if dividendInfo.DividendYield != nil && dividendInfo.DividendYield.Raw < 0 {
				t.Error("Dividend yield should not be negative")
			}
		})
	}
}

// TestFetchCurrentDividendYield tests fetching just the dividend yield
func TestFetchCurrentDividendYield(t *testing.T) {
	ticker := NewTicker("AAPL")

	yield, err := ticker.FetchCurrentDividendYield()
	if err != nil {
		t.Skipf("Skipping test due to API error: %v", err)
		return
	}

	if yield < 0 {
		t.Error("Dividend yield should not be negative")
	}

	// AAPL typically has a dividend yield between 0% and 5%
	if yield > 0.1 { // 10% seems unreasonably high for AAPL
		t.Logf("Warning: Dividend yield seems unusually high: %.2f%%", yield*100)
	}
}

// TestFetchDividendRate tests fetching the annual dividend rate
func TestFetchDividendRate(t *testing.T) {
	ticker := NewTicker("AAPL")

	rate, err := ticker.FetchDividendRate()
	if err != nil {
		t.Skipf("Skipping test due to API error: %v", err)
		return
	}

	if rate < 0 {
		t.Error("Dividend rate should not be negative")
	}

	t.Logf("AAPL annual dividend rate: $%.2f", rate)
}

// TestIsDividendPaying tests checking if a stock pays dividends
func TestIsDividendPaying(t *testing.T) {
	testCases := []struct {
		symbol           string
		description      string
		expectedDividend bool // approximate expectation
	}{
		{"AAPL", "Apple - should pay dividends", true},
		{"MSFT", "Microsoft - should pay dividends", true},
		{"TSLA", "Tesla - historically no dividends", false},
		{"AMZN", "Amazon - historically no dividends", false},
		{"KO", "Coca-Cola - dividend aristocrat", true},
	}

	for _, tc := range testCases {
		t.Run(tc.symbol, func(t *testing.T) {
			ticker := NewTicker(tc.symbol)

			isPaying, err := ticker.IsDividendPaying()
			if err != nil {
				t.Skipf("Skipping %s due to API error: %v", tc.symbol, err)
				return
			}

			t.Logf("%s (%s) - Pays dividends: %v", tc.symbol, tc.description, isPaying)

			// Note: We don't assert the expected value because dividend policies can change
			// This test is mainly for functionality verification
		})
	}
}

// TestDividendInfoStructure tests the DividendInfo data structure
func TestDividendInfoStructure(t *testing.T) {
	// Test creating DividendInfo struct manually
	dividendInfo := DividendInfo{
		DividendRate:  &PriceValue{Raw: 0.88, Fmt: "$0.88"},
		DividendYield: &PriceValue{Raw: 0.015, Fmt: "1.50%"},
		PayoutRatio:   &PriceValue{Raw: 0.25, Fmt: "25.00%"},
	}

	if dividendInfo.DividendRate.Raw != 0.88 {
		t.Errorf("Expected dividend rate 0.88, got %f", dividendInfo.DividendRate.Raw)
	}

	if dividendInfo.DividendYield.Raw != 0.015 {
		t.Errorf("Expected dividend yield 0.015, got %f", dividendInfo.DividendYield.Raw)
	}

	if dividendInfo.PayoutRatio.Raw != 0.25 {
		t.Errorf("Expected payout ratio 0.25, got %f", dividendInfo.PayoutRatio.Raw)
	}
}

// TestNonDividendStock tests dividend functionality with non-dividend paying stock
func TestNonDividendStock(t *testing.T) {
	// Test with a stock that typically doesn't pay dividends
	ticker := NewTicker("AMZN")

	isPaying, err := ticker.IsDividendPaying()
	if err != nil {
		t.Skipf("Skipping test due to API error: %v", err)
		return
	}

	if isPaying {
		t.Log("AMZN appears to be paying dividends now - policy may have changed")
	} else {
		t.Log("AMZN confirmed as non-dividend paying stock")
	}

	// Try to get dividend rate for non-dividend stock
	_, err = ticker.FetchDividendRate()
	if err == nil {
		t.Log("AMZN has dividend rate available - this might be unexpected")
	} else {
		t.Logf("Expected error for non-dividend stock: %v", err)
	}
}

// Benchmark dividend functions
func BenchmarkFetchDividendInfo(b *testing.B) {
	ticker := NewTicker("AAPL")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ticker.FetchDividendInfo()
		if err != nil {
			b.Skipf("Skipping benchmark due to API error: %v", err)
			return
		}
	}
}

func BenchmarkFetchCurrentDividendYield(b *testing.B) {
	ticker := NewTicker("AAPL")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ticker.FetchCurrentDividendYield()
		if err != nil {
			b.Skipf("Skipping benchmark due to API error: %v", err)
			return
		}
	}
}
