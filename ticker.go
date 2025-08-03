package yfinance_api

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/url"
)

type Ticker struct {
	Symbol string
	Client *Client
}

// InstantiateTicker creates a new Ticker instance with the provided symbol and exchange name.
// It is used to represent a financial instrument on a specific exchange.
func (c *YFinanceAPI) InstantiateTicker(symbol string) *Ticker {
	ticker := &Ticker{
		Symbol: symbol,
		Client: c.Client,
	}

	return ticker
}

// GetSymbol returns the symbol of the Ticker instance.
func (t *Ticker) GetSymbol() string {
	return t.Symbol
}

// SetSymbol sets the symbol of the Ticker instance.
func (t *Ticker) SetSymbol(symbol string) {
	t.Symbol = symbol
}

// FetchInformation retrieves detailed information about the ticker from Yahoo Finance.
// It uses the Yahoo Finance quoteSummary API to fetch the price module data.
// Returns a YahooTickerInfo struct containing the ticker's price information or an error if the request
// fails or if the response cannot be parsed.
func (t *Ticker) FetchInformation() (YahooTickerInfo, error) {
	// Prepare URL parameters to request the "price" module
	params := url.Values{}
	params.Add("modules", "price")

	// Build the endpoint URL for the Yahoo Finance quoteSummary API
	endpoint := fmt.Sprintf("%s/v10/finance/quoteSummary/%s", BaseUrl, t.Symbol)

	// Make the HTTP GET request using the client
	resp, err := t.Client.Get(endpoint, params)
	if err != nil {
		slog.Error("Failed to get ticker info", "err", err)
		return YahooTickerInfo{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			slog.Error("Failed to close response body", "err", err)
		}
	}(resp.Body)

	// Read the response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return YahooTickerInfo{}, fmt.Errorf("failed to read response body: %w", err)
	}

	// Unmarshal the JSON response into the YahooInfoResponse struct
	var infoResponse YahooInfoResponse
	if err := json.Unmarshal(bodyBytes, &infoResponse); err != nil {
		return YahooTickerInfo{}, fmt.Errorf("failed to decode info JSON: %w", err)
	}

	// Check if the result array is empty
	if len(infoResponse.QuoteSummary.Result) == 0 {
		return YahooTickerInfo{}, fmt.Errorf("no info found for symbol: %s", t.Symbol)
	}

	// Return the ticker price information
	return infoResponse.QuoteSummary.Result[0].Price, nil
}

// FetchPriceValue retrieves the current price value of the ticker from Yahoo Finance.
// It uses the FetchInformation method to get the ticker's price information.
// Returns a PriceValue struct containing the price information or an error if the request fails
// or if the response cannot be parsed.
func (t *Ticker) FetchPriceValue() (PriceValue, error) {
	info, err := t.FetchInformation()
	if err != nil {
		slog.Error("Failed to fetch ticker price value", "err", err)
		return PriceValue{}, err
	}

	// Check if the RegularMarketPrice is available and return it
	if info.RegularMarketPrice != nil {
		return *info.RegularMarketPrice, nil
	}

	// If RegularMarketPrice is not available, return an error
	return PriceValue{}, fmt.Errorf("regular market price not available for symbol: %s", t.Symbol)
}

// FetchHistoricalData retrieves historical price data for the ticker from Yahoo Finance.
// It accepts query parameters directly and handles all processing internally.
// Parameters:
//   - rangeParam: time range (e.g., "1d", "5d", "1mo", "3mo", "6mo", "1y", "2y", "5y", "10y", "ytd", "max")
//   - interval: data interval (e.g., "1m", "2m", "5m", "15m", "30m", "60m", "90m", "1h", "1d", "5d", "1wk", "1mo", "3mo")
//   - period1: start timestamp (optional, can be empty string)
//   - period2: end timestamp (optional, can be empty string)
func (t *Ticker) FetchHistoricalData(rangeParam, interval, period1, period2 string) (map[string]PriceData, error) {
	// Set default values if not provided
	if interval == "" {
		interval = "1d"
	}
	if rangeParam == "" {
		rangeParam = "1y"
	}

	// Build query parameters
	params := url.Values{}
	params.Add("range", rangeParam)

	params.Add("interval", interval)
	if period1 != "" {
		params.Add("period1", period1)
	}
	if period2 != "" {
		params.Add("period2", period2)
	}

	// Build the endpoint URL
	endpoint := fmt.Sprintf("%s/v8/finance/chart/%s", BaseUrl, t.Symbol)

	// Make the HTTP request
	resp, err := t.Client.Get(endpoint, params)
	if err != nil {
		slog.Error("Failed to get historical data", "err", err)
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			slog.Error("Failed to close response body", "err", err)
		}
	}(resp.Body)

	// Decode the JSON response
	var historyResponse YahooHistoryResponse
	if err := json.NewDecoder(resp.Body).Decode(&historyResponse); err != nil {
		return nil, fmt.Errorf("failed to decode history data JSON response: %v", err)
	}

	// Check if we have data
	if len(historyResponse.Chart.Result) == 0 {
		return nil, fmt.Errorf("no data found for symbol: %s", t.Symbol)
	}

	// Transform and return the data
	return transformHistoricalData(historyResponse, interval), nil
}

// FetchNews retrieves recent news articles related to the ticker from Yahoo Finance.
// Parameters:
//   - count: number of news articles to fetch (optional, defaults to 10)
//   - start: starting index for pagination (optional, defaults to 0)
//
// Returns a slice of NewsItem structs containing news articles related to the ticker
func (t *Ticker) FetchNews(count, start int) ([]NewsItem, error) {
	// Set default values if not provided
	if count <= 0 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	// Build query parameters
	params := url.Values{}
	params.Add("symbols", t.Symbol)
	params.Add("count", fmt.Sprintf("%d", count))
	params.Add("start", fmt.Sprintf("%d", start))
	params.Add("region", "US")
	params.Add("lang", "en-US")

	// Build the endpoint URL for Yahoo Finance news API
	endpoint := fmt.Sprintf("%s/v1/finance/search", BaseUrl)

	// Make the HTTP request
	resp, err := t.Client.Get(endpoint, params)
	if err != nil {
		slog.Error("Failed to get news", "err", err)
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			slog.Error("Failed to close response body", "err", err)
		}
	}(resp.Body)

	// Read response body for debugging if needed
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Try alternative news endpoint if the first one doesn't work
	var newsResponse struct {
		News []NewsItem `json:"news"`
	}

	if err := json.Unmarshal(bodyBytes, &newsResponse); err != nil {
		// If that fails, try a different structure
		var altResponse struct {
			Result struct {
				News []NewsItem `json:"news"`
			} `json:"result"`
		}
		if err := json.Unmarshal(bodyBytes, &altResponse); err != nil {
			return nil, fmt.Errorf("failed to decode news JSON response: %v", err)
		}
		return altResponse.Result.News, nil
	}

	return newsResponse.News, nil
}

// FetchNewsAlternative uses an alternative endpoint to fetch news for the ticker
// This method uses the quoteSummary API with recommendationTrend module which sometimes includes news
func (t *Ticker) FetchNewsAlternative() ([]NewsItem, error) {
	// Build query parameters for quoteSummary API
	params := url.Values{}
	params.Add("modules", "recommendationTrend,upgradeDowngradeHistory")

	// Build the endpoint URL
	endpoint := fmt.Sprintf("%s/v10/finance/quoteSummary/%s", BaseUrl, t.Symbol)

	// Make the HTTP request
	resp, err := t.Client.Get(endpoint, params)
	if err != nil {
		slog.Error("Failed to get alternative news", "err", err)
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			slog.Error("Failed to close response body", "err", err)
		}
	}(resp.Body)

	// For now, return empty slice as this is a fallback method
	// In a real implementation, you would parse the response for any news data
	return []NewsItem{}, nil
}

// FetchFinancialData retrieves comprehensive financial data including ratios, fundamentals, and financial statements
// Returns a FinancialData struct containing all financial metrics for fundamental analysis
func (t *Ticker) FetchFinancialData() (FinancialData, error) {
	// Build query parameters to request multiple financial modules
	params := url.Values{}
	params.Add("modules", "defaultKeyStatistics,financialData,summaryDetail,incomeStatementHistory,balanceSheetHistory,cashflowStatementHistory")

	// Build the endpoint URL
	endpoint := fmt.Sprintf("%s/v10/finance/quoteSummary/%s", BaseUrl, t.Symbol)

	// Make the HTTP request
	resp, err := t.Client.Get(endpoint, params)
	if err != nil {
		slog.Error("Failed to get financial data", "err", err)
		return FinancialData{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			slog.Error("Failed to close response body", "err", err)
		}
	}(resp.Body)

	// Decode the JSON response
	var financialResponse YahooFinancialResponse
	if err := json.NewDecoder(resp.Body).Decode(&financialResponse); err != nil {
		return FinancialData{}, fmt.Errorf("failed to decode financial data JSON response: %v", err)
	}

	// Check if we have data
	if len(financialResponse.QuoteSummary.Result) == 0 {
		return FinancialData{}, fmt.Errorf("no financial data found for symbol: %s", t.Symbol)
	}

	result := financialResponse.QuoteSummary.Result[0]

	// Transform and return the financial data
	return t.transformFinancialData(result), nil
}

// FetchFinancialRatios retrieves only the financial ratios for quick analysis
func (t *Ticker) FetchFinancialRatios() (FinancialRatios, error) {
	params := url.Values{}
	params.Add("modules", "defaultKeyStatistics,financialData,summaryDetail")

	endpoint := fmt.Sprintf("%s/v10/finance/quoteSummary/%s", BaseUrl, t.Symbol)

	resp, err := t.Client.Get(endpoint, params)
	if err != nil {
		slog.Error("Failed to get financial ratios", "err", err)
		return FinancialRatios{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			slog.Error("Failed to close response body", "err", err)
		}
	}(resp.Body)

	var financialResponse YahooFinancialResponse
	if err := json.NewDecoder(resp.Body).Decode(&financialResponse); err != nil {
		return FinancialRatios{}, fmt.Errorf("failed to decode financial ratios JSON response: %v", err)
	}

	if len(financialResponse.QuoteSummary.Result) == 0 {
		return FinancialRatios{}, fmt.Errorf("no financial ratios found for symbol: %s", t.Symbol)
	}

	result := financialResponse.QuoteSummary.Result[0]
	return t.extractFinancialRatios(result), nil
}

// FetchKeyStatistics retrieves key financial statistics and metrics
func (t *Ticker) FetchKeyStatistics() (FinancialSummary, error) {
	params := url.Values{}
	params.Add("modules", "defaultKeyStatistics,summaryDetail")

	endpoint := fmt.Sprintf("%s/v10/finance/quoteSummary/%s", BaseUrl, t.Symbol)

	resp, err := t.Client.Get(endpoint, params)
	if err != nil {
		slog.Error("Failed to get key statistics", "err", err)
		return FinancialSummary{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			slog.Error("Failed to close response body", "err", err)
		}
	}(resp.Body)

	var financialResponse YahooFinancialResponse
	if err := json.NewDecoder(resp.Body).Decode(&financialResponse); err != nil {
		return FinancialSummary{}, fmt.Errorf("failed to decode key statistics JSON response: %v", err)
	}

	if len(financialResponse.QuoteSummary.Result) == 0 {
		return FinancialSummary{}, fmt.Errorf("no key statistics found for symbol: %s", t.Symbol)
	}

	result := financialResponse.QuoteSummary.Result[0]
	return t.extractFinancialSummary(result), nil
}

// FetchIncomeStatement retrieves the latest income statement data
func (t *Ticker) FetchIncomeStatement() (IncomeStatement, error) {
	params := url.Values{}
	params.Add("modules", "incomeStatementHistory")

	endpoint := fmt.Sprintf("%s/v10/finance/quoteSummary/%s", BaseUrl, t.Symbol)

	resp, err := t.Client.Get(endpoint, params)
	if err != nil {
		slog.Error("Failed to get income statement", "err", err)
		return IncomeStatement{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			slog.Error("Failed to close response body", "err", err)
		}
	}(resp.Body)

	var financialResponse YahooFinancialResponse
	if err := json.NewDecoder(resp.Body).Decode(&financialResponse); err != nil {
		return IncomeStatement{}, fmt.Errorf("failed to decode income statement JSON response: %v", err)
	}

	if len(financialResponse.QuoteSummary.Result) == 0 {
		return IncomeStatement{}, fmt.Errorf("no income statement found for symbol: %s", t.Symbol)
	}

	result := financialResponse.QuoteSummary.Result[0]
	return t.extractIncomeStatement(result), nil
}

// FetchBalanceSheet retrieves the latest balance sheet data
func (t *Ticker) FetchBalanceSheet() (BalanceSheet, error) {
	params := url.Values{}
	params.Add("modules", "balanceSheetHistory")

	endpoint := fmt.Sprintf("%s/v10/finance/quoteSummary/%s", BaseUrl, t.Symbol)

	resp, err := t.Client.Get(endpoint, params)
	if err != nil {
		slog.Error("Failed to get balance sheet", "err", err)
		return BalanceSheet{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			slog.Error("Failed to close response body", "err", err)
		}
	}(resp.Body)

	var financialResponse YahooFinancialResponse
	if err := json.NewDecoder(resp.Body).Decode(&financialResponse); err != nil {
		return BalanceSheet{}, fmt.Errorf("failed to decode balance sheet JSON response: %v", err)
	}

	if len(financialResponse.QuoteSummary.Result) == 0 {
		return BalanceSheet{}, fmt.Errorf("no balance sheet found for symbol: %s", t.Symbol)
	}

	result := financialResponse.QuoteSummary.Result[0]
	return t.extractBalanceSheet(result), nil
}

// FetchCashFlow retrieves the latest cash flow statement data
func (t *Ticker) FetchCashFlow() (CashFlow, error) {
	params := url.Values{}
	params.Add("modules", "cashflowStatementHistory")

	endpoint := fmt.Sprintf("%s/v10/finance/quoteSummary/%s", BaseUrl, t.Symbol)

	resp, err := t.Client.Get(endpoint, params)
	if err != nil {
		slog.Error("Failed to get cash flow", "err", err)
		return CashFlow{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			slog.Error("Failed to close response body", "err", err)
		}
	}(resp.Body)

	var financialResponse YahooFinancialResponse
	if err := json.NewDecoder(resp.Body).Decode(&financialResponse); err != nil {
		return CashFlow{}, fmt.Errorf("failed to decode cash flow JSON response: %v", err)
	}

	if len(financialResponse.QuoteSummary.Result) == 0 {
		return CashFlow{}, fmt.Errorf("no cash flow found for symbol: %s", t.Symbol)
	}

	result := financialResponse.QuoteSummary.Result[0]
	return t.extractCashFlow(result), nil
}

// DividendInfo represents dividend-related information for a stock
type DividendInfo struct {
	DividendRate             *PriceValue `json:"dividendRate"`             // Annual dividend per share
	DividendYield            *PriceValue `json:"dividendYield"`            // Dividend yield percentage
	DividendsPaid            *PriceValue `json:"dividendsPaid"`            // Total dividends paid (from cash flow)
	PayoutRatio              *PriceValue `json:"payoutRatio"`              // Dividend payout ratio
	ExDividendDate           *PriceValue `json:"exDividendDate"`           // Ex-dividend date
	DividendDate             *PriceValue `json:"dividendDate"`             // Dividend payment date
	FiveYearAvgDividendYield *PriceValue `json:"fiveYearAvgDividendYield"` // 5-year average dividend yield
}

// FetchDividendInfo retrieves comprehensive dividend information for the ticker
// Returns dividend rate, yield, payment history, and related metrics
func (t *Ticker) FetchDividendInfo() (DividendInfo, error) {
	// Get comprehensive financial data including dividend information
	params := url.Values{}
	params.Add("modules", "summaryDetail,defaultKeyStatistics,cashflowStatementHistory,calendarEvents")

	endpoint := fmt.Sprintf("%s/v10/finance/quoteSummary/%s", BaseUrl, t.Symbol)

	resp, err := t.Client.Get(endpoint, params)
	if err != nil {
		slog.Error("Failed to get dividend info", "err", err)
		return DividendInfo{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			slog.Error("Failed to close response body", "err", err)
		}
	}(resp.Body)

	// Decode the JSON response
	var financialResponse struct {
		QuoteSummary struct {
			Result []struct {
				SummaryDetail *struct {
					DividendRate             *PriceValue `json:"dividendRate"`
					DividendYield            *PriceValue `json:"dividendYield"`
					ExDividendDate           *PriceValue `json:"exDividendDate"`
					DividendDate             *PriceValue `json:"dividendDate"`
					PayoutRatio              *PriceValue `json:"payoutRatio"`
					FiveYearAvgDividendYield *PriceValue `json:"fiveYearAvgDividendYield"`
				} `json:"summaryDetail"`
				DefaultKeyStatistics *struct {
					DividendRate             *PriceValue `json:"dividendRate"`
					DividendYield            *PriceValue `json:"dividendYield"`
					PayoutRatio              *PriceValue `json:"payoutRatio"`
					FiveYearAvgDividendYield *PriceValue `json:"fiveYearAvgDividendYield"`
				} `json:"defaultKeyStatistics"`
				CashflowStatementHistory *struct {
					CashflowStatements []struct {
						DividendsPaid *PriceValue `json:"dividendsPaid"`
					} `json:"cashflowStatements"`
				} `json:"cashflowStatementHistory"`
			} `json:"result"`
			Error interface{} `json:"error"`
		} `json:"quoteSummary"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&financialResponse); err != nil {
		return DividendInfo{}, fmt.Errorf("failed to decode dividend info JSON response: %v", err)
	}

	if len(financialResponse.QuoteSummary.Result) == 0 {
		return DividendInfo{}, fmt.Errorf("no dividend info found for symbol: %s", t.Symbol)
	}

	result := financialResponse.QuoteSummary.Result[0]
	return t.extractDividendInfo(result), nil
}

// FetchCurrentDividendYield retrieves just the current dividend yield for quick access
func (t *Ticker) FetchCurrentDividendYield() (float64, error) {
	dividendInfo, err := t.FetchDividendInfo()
	if err != nil {
		return 0, err
	}

	if dividendInfo.DividendYield == nil {
		return 0, fmt.Errorf("dividend yield not available for symbol: %s", t.Symbol)
	}

	return dividendInfo.DividendYield.Raw, nil
}

// FetchDividendRate retrieves the annual dividend rate per share
func (t *Ticker) FetchDividendRate() (float64, error) {
	dividendInfo, err := t.FetchDividendInfo()
	if err != nil {
		return 0, err
	}

	if dividendInfo.DividendRate == nil {
		return 0, fmt.Errorf("dividend rate not available for symbol: %s", t.Symbol)
	}

	return dividendInfo.DividendRate.Raw, nil
}

// IsDividendPaying checks if the stock currently pays dividends
func (t *Ticker) IsDividendPaying() (bool, error) {
	dividendInfo, err := t.FetchDividendInfo()
	if err != nil {
		return false, err
	}

	// A stock is considered dividend-paying if it has a positive dividend rate
	return dividendInfo.DividendRate != nil && dividendInfo.DividendRate.Raw > 0, nil
}
