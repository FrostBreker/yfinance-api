# YFinance API Go Package

[![Go Version](https://img.shields.io/badge/Go-1.19+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Tests](https://github.com/FrostBreker/yfinance-api/workflows/Tests/badge.svg)](https://github.com/FrostBreker/yfinance-api/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/FrostBreker/yfinance-api)](https://goreportcard.com/report/github.com/FrostBreker/yfinance-api)

A comprehensive Go package for accessing Yahoo Finance data, providing easy-to-use methods for fetching stock prices, financial ratios, fundamentals, dividend information, historical data, and news.

## Features

- 📈 **Real-time Stock Prices** - Current market prices and trading information
- 📊 **Historical Data** - OHLCV data with flexible time ranges and intervals
- 💰 **Dividend Information** - Comprehensive dividend data and yield calculations
- 📰 **Financial News** - Latest news articles related to stocks
- 🔢 **Financial Ratios** - P/E, P/B, ROE, ROA, and 20+ other key ratios
- 📋 **Financial Statements** - Income statement, balance sheet, and cash flow data
- 🏢 **Company Fundamentals** - Market cap, beta, 52-week highs/lows, and more
- 🚀 **Easy Integration** - Simple API with comprehensive error handling
- ⚡ **High Performance** - Efficient HTTP client with singleton pattern
- 🧪 **Well Tested** - Comprehensive test suite with benchmarks

## Installation

```bash
go get github.com/FrostBreker/yfinance-api
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    yfinance "github.com/FrostBreker/yfinance-api"
)

func main() {
    // Create a ticker for Apple Inc.
    ticker := yfinance.NewTicker("AAPL")
    
    // Get current price
    price, err := ticker.FetchPriceValue()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("AAPL Current Price: %s\n", price.Fmt)
    
    // Get dividend information
    dividend, err := ticker.FetchDividendInfo()
    if err != nil {
        log.Fatal(err)
    }
    if dividend.DividendYield != nil {
        fmt.Printf("AAPL Dividend Yield: %s\n", dividend.DividendYield.Fmt)
    }
}
```

## Usage Examples

### Basic Stock Information

```go
ticker := yfinance.NewTicker("MSFT")

// Get comprehensive stock information
info, err := ticker.FetchInformation()
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Symbol: %s\n", info.Symbol)
fmt.Printf("Company: %s\n", info.LongName)
fmt.Printf("Current Price: %s\n", info.RegularMarketPrice.Fmt)
fmt.Printf("Market Cap: %s\n", info.MarketCap.Fmt)
```

### Historical Data

```go
ticker := yfinance.NewTicker("GOOGL")

// Get 1 year of daily data
historicalData, err := ticker.FetchHistoricalData("1y", "1d", "", "")
if err != nil {
    log.Fatal(err)
}

for date, data := range historicalData {
    if data.Close != nil {
        fmt.Printf("%s: $%.2f\n", date, *data.Close)
    }
}

// Get 1 day of 5-minute intervals
intradayData, err := ticker.FetchHistoricalData("1d", "5m", "", "")
```

### Dividend Analysis

```go
ticker := yfinance.NewTicker("KO") // Coca-Cola

// Check if stock pays dividends
isPaying, err := ticker.IsDividendPaying()
if err != nil {
    log.Fatal(err)
}

if isPaying {
    // Get detailed dividend information
    dividend, err := ticker.FetchDividendInfo()
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Annual Dividend: %s\n", dividend.DividendRate.Fmt)
    fmt.Printf("Dividend Yield: %s\n", dividend.DividendYield.Fmt)
    fmt.Printf("Payout Ratio: %s\n", dividend.PayoutRatio.Fmt)
    
    // Quick access methods
    yield, _ := ticker.FetchCurrentDividendYield()
    fmt.Printf("Yield (float): %.2f%%\n", yield*100)
}
```

### Financial Ratios and Fundamentals

```go
ticker := yfinance.NewTicker("TSLA")

// Get comprehensive financial data
financialData, err := ticker.FetchFinancialData()
if err != nil {
    log.Fatal(err)
}

// Access financial ratios
ratios := financialData.Ratios
if ratios.PriceToEarningsRatio != nil {
    fmt.Printf("P/E Ratio: %s\n", ratios.PriceToEarningsRatio.Fmt)
}
if ratios.PriceToBookRatio != nil {
    fmt.Printf("P/B Ratio: %s\n", ratios.PriceToBookRatio.Fmt)
}

// Access financial statements
income := financialData.IncomeStatement
if income.TotalRevenue != nil {
    fmt.Printf("Revenue: %s\n", income.TotalRevenue.Fmt)
}
```

### Financial Statements

```go
ticker := yfinance.NewTicker("AMZN")

// Get individual financial statements
income, err := ticker.FetchIncomeStatement()
balance, err := ticker.FetchBalanceSheet()
cashflow, err := ticker.FetchCashFlow()

// Income Statement
if income.TotalRevenue != nil {
    fmt.Printf("Revenue: %s\n", income.TotalRevenue.Fmt)
}
if income.NetIncome != nil {
    fmt.Printf("Net Income: %s\n", income.NetIncome.Fmt)
}

// Balance Sheet
if balance.TotalAssets != nil {
    fmt.Printf("Total Assets: %s\n", balance.TotalAssets.Fmt)
}
if balance.TotalDebt != nil {
    fmt.Printf("Total Debt: %s\n", balance.TotalDebt.Fmt)
}

// Cash Flow
if cashflow.OperatingCashFlow != nil {
    fmt.Printf("Operating Cash Flow: %s\n", cashflow.OperatingCashFlow.Fmt)
}
```

### News and Market Data

```go
ticker := yfinance.NewTicker("NVDA")

// Get latest news
news, err := ticker.FetchNews(5, 0) // Get 5 recent articles
if err != nil {
    log.Fatal(err)
}

for _, article := range news {
    fmt.Printf("Title: %s\n", article.Title)
    fmt.Printf("Publisher: %s\n", article.Publisher)
    fmt.Printf("Link: %s\n", article.Link)
    fmt.Println("---")
}
```

### Multiple Tickers

```go
symbols := []string{"AAPL", "MSFT", "GOOGL", "AMZN", "TSLA"}

for _, symbol := range symbols {
    ticker := yfinance.NewTicker(symbol)
    
    price, err := ticker.FetchPriceValue()
    if err != nil {
        fmt.Printf("Error fetching %s: %v\n", symbol, err)
        continue
    }
    
    fmt.Printf("%s: %s\n", symbol, price.Fmt)
}
```

## API Reference

### Core Functions

| Function | Description | Returns |
|----------|-------------|---------|
| `NewClient()` | Create a new YFinance API client | `*YFinanceAPI` |
| `NewTicker(symbol)` | Create a ticker instance | `*Ticker` |

### Ticker Methods

#### Price & Information
| Method | Description | Returns |
|--------|-------------|---------|
| `FetchInformation()` | Get comprehensive ticker info | `YahooTickerInfo` |
| `FetchPriceValue()` | Get current stock price | `PriceValue` |

#### Historical Data
| Method | Parameters | Description |
|--------|------------|-------------|
| `FetchHistoricalData()` | `range, interval, period1, period2` | Get OHLCV historical data |

**Range Options**: `1d`, `5d`, `1mo`, `3mo`, `6mo`, `1y`, `2y`, `5y`, `10y`, `ytd`, `max`

**Interval Options**: `1m`, `2m`, `5m`, `15m`, `30m`, `60m`, `90m`, `1h`, `1d`, `5d`, `1wk`, `1mo`, `3mo`

#### Dividend Information
| Method | Description | Returns |
|--------|-------------|---------|
| `FetchDividendInfo()` | Complete dividend information | `DividendInfo` |
| `FetchCurrentDividendYield()` | Current dividend yield | `float64` |
| `FetchDividendRate()` | Annual dividend per share | `float64` |
| `IsDividendPaying()` | Check if stock pays dividends | `bool` |

#### Financial Analysis
| Method | Description | Returns |
|--------|-------------|---------|
| `FetchFinancialData()` | Complete financial analysis | `FinancialData` |
| `FetchFinancialRatios()` | Financial ratios only | `FinancialRatios` |
| `FetchKeyStatistics()` | Key financial metrics | `FinancialSummary` |
| `FetchIncomeStatement()` | Income statement data | `IncomeStatement` |
| `FetchBalanceSheet()` | Balance sheet data | `BalanceSheet` |
| `FetchCashFlow()` | Cash flow statement | `CashFlow` |

#### News
| Method | Parameters | Description |
|--------|------------|-------------|
| `FetchNews()` | `count, start` | Get news articles |

## Data Structures

### PriceValue
```go
type PriceValue struct {
    Raw     float64 `json:"raw"`     // Raw numeric value
    Fmt     string  `json:"fmt"`     // Formatted string (e.g., "$150.25")
    LongFmt string  `json:"longFmt"` // Long format (optional)
}
```

### PriceData (Historical)
```go
type PriceData struct {
    Open   *float64 `json:"open"`
    High   *float64 `json:"high"`
    Low    *float64 `json:"low"`
    Close  *float64 `json:"close"`
    Volume *int64   `json:"volume"`
}
```

### DividendInfo
```go
type DividendInfo struct {
    DividendRate             *PriceValue `json:"dividendRate"`
    DividendYield            *PriceValue `json:"dividendYield"`
    DividendsPaid            *PriceValue `json:"dividendsPaid"`
    PayoutRatio              *PriceValue `json:"payoutRatio"`
    ExDividendDate           *PriceValue `json:"exDividendDate"`
    DividendDate             *PriceValue `json:"dividendDate"`
    FiveYearAvgDividendYield *PriceValue `json:"fiveYearAvgDividendYield"`
}
```

### FinancialRatios
Key financial ratios including:
- **Valuation**: P/E, P/B, P/S, Enterprise ratios
- **Profitability**: ROE, ROA, Profit margins
- **Liquidity**: Current ratio, Quick ratio
- **Leverage**: Debt-to-equity ratios
- **Growth**: Earnings and revenue growth

## Error Handling

The package provides comprehensive error handling. All methods return an error as the second return value:

```go
ticker := yfinance.NewTicker("INVALID")
_, err := ticker.FetchPriceValue()
if err != nil {
    fmt.Printf("Error: %v\n", err)
    // Handle error appropriately
}
```

Common error scenarios:
- Invalid ticker symbols
- Network connectivity issues
- API rate limiting
- Missing data for specific metrics

## Performance

- **Singleton HTTP Client**: Efficient connection reuse and cookie management
- **Concurrent Safe**: Thread-safe operations
- **Memory Efficient**: Pointer-based optional fields to minimize memory usage
- **Fast JSON Parsing**: Optimized JSON unmarshaling

## Testing

Run the test suite:

```bash
# Run all tests
go test

# Run tests with verbose output
go test -v

# Run specific tests
go test -run TestFetchDividendInfo

# Run benchmarks
go test -bench=.

# Run tests with coverage
go test -cover
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Disclaimer

This package is for educational and research purposes. Always verify financial data from official sources before making investment decisions. The authors are not responsible for any financial losses incurred from using this package.

## Acknowledgments

- Yahoo Finance for providing the data API
- The Go community for excellent tooling and libraries

## Support

If you find this package useful, please consider:
- ⭐ Starring the repository
- 🐛 Reporting bugs
- 💡 Suggesting new features
- 📖 Improving documentation

Huge thanks to [@oscarli916](https://github.com/oscarli916) for the original idea and implementation of this project with [oscarli916/yahoo-finance-api](https://github.com/oscarli916/yahoo-finance-api).

---

**Happy Trading! 📈**
