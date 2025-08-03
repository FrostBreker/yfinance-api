package yfinance_api

import (
	"strings"
	"time"
)

// transformFinancialData converts Yahoo Finance API response into structured FinancialData
func (t *Ticker) transformFinancialData(result struct {
	DefaultKeyStatistics *FinancialSummary `json:"defaultKeyStatistics"`
	FinancialData        *FinancialRatios  `json:"financialData"`
	SummaryDetail        *struct {
		MarketCap                    *PriceValue `json:"marketCap"`
		ForwardPE                    *PriceValue `json:"forwardPE"`
		TrailingPE                   *PriceValue `json:"trailingPE"`
		PriceToSalesTrailing12Months *PriceValue `json:"priceToSalesTrailing12Months"`
		PriceToBook                  *PriceValue `json:"priceToBook"`
		Beta                         *PriceValue `json:"beta"`
		DividendRate                 *PriceValue `json:"dividendRate"`
		DividendYield                *PriceValue `json:"dividendYield"`
	} `json:"summaryDetail"`
	IncomeStatementHistory *struct {
		IncomeStatementHistory []struct {
			EndDate         *PriceValue `json:"endDate"`
			TotalRevenue    *PriceValue `json:"totalRevenue"`
			GrossProfit     *PriceValue `json:"grossProfit"`
			OperatingIncome *PriceValue `json:"operatingIncome"`
			NetIncome       *PriceValue `json:"netIncome"`
			Ebitda          *PriceValue `json:"ebitda"`
		} `json:"incomeStatementHistory"`
	} `json:"incomeStatementHistory"`
	BalanceSheetHistory *struct {
		BalanceSheetStatements []struct {
			EndDate                *PriceValue `json:"endDate"`
			TotalAssets            *PriceValue `json:"totalAssets"`
			TotalLiab              *PriceValue `json:"totalLiab"`
			TotalStockholderEquity *PriceValue `json:"totalStockholderEquity"`
			TotalDebt              *PriceValue `json:"totalDebt"`
			Cash                   *PriceValue `json:"cash"`
		} `json:"balanceSheetStatements"`
	} `json:"balanceSheetHistory"`
	CashflowStatementHistory *struct {
		CashflowStatements []struct {
			EndDate                          *PriceValue `json:"endDate"`
			TotalCashFromOperatingActivities *PriceValue `json:"totalCashFromOperatingActivities"`
			CapitalExpenditures              *PriceValue `json:"capitalExpenditures"`
			FreeCashFlow                     *PriceValue `json:"freeCashFlow"`
			DividendsPaid                    *PriceValue `json:"dividendsPaid"`
		} `json:"cashflowStatements"`
	} `json:"cashflowStatementHistory"`
}) FinancialData {
	return FinancialData{
		Ratios:          t.extractFinancialRatios(result),
		Summary:         t.extractFinancialSummary(result),
		IncomeStatement: t.extractIncomeStatement(result),
		BalanceSheet:    t.extractBalanceSheet(result),
		CashFlow:        t.extractCashFlow(result),
	}
}

// extractFinancialRatios extracts financial ratios from the API response
func (t *Ticker) extractFinancialRatios(result struct {
	DefaultKeyStatistics *FinancialSummary `json:"defaultKeyStatistics"`
	FinancialData        *FinancialRatios  `json:"financialData"`
	SummaryDetail        *struct {
		MarketCap                    *PriceValue `json:"marketCap"`
		ForwardPE                    *PriceValue `json:"forwardPE"`
		TrailingPE                   *PriceValue `json:"trailingPE"`
		PriceToSalesTrailing12Months *PriceValue `json:"priceToSalesTrailing12Months"`
		PriceToBook                  *PriceValue `json:"priceToBook"`
		Beta                         *PriceValue `json:"beta"`
		DividendRate                 *PriceValue `json:"dividendRate"`
		DividendYield                *PriceValue `json:"dividendYield"`
	} `json:"summaryDetail"`
	IncomeStatementHistory *struct {
		IncomeStatementHistory []struct {
			EndDate         *PriceValue `json:"endDate"`
			TotalRevenue    *PriceValue `json:"totalRevenue"`
			GrossProfit     *PriceValue `json:"grossProfit"`
			OperatingIncome *PriceValue `json:"operatingIncome"`
			NetIncome       *PriceValue `json:"netIncome"`
			Ebitda          *PriceValue `json:"ebitda"`
		} `json:"incomeStatementHistory"`
	} `json:"incomeStatementHistory"`
	BalanceSheetHistory *struct {
		BalanceSheetStatements []struct {
			EndDate                *PriceValue `json:"endDate"`
			TotalAssets            *PriceValue `json:"totalAssets"`
			TotalLiab              *PriceValue `json:"totalLiab"`
			TotalStockholderEquity *PriceValue `json:"totalStockholderEquity"`
			TotalDebt              *PriceValue `json:"totalDebt"`
			Cash                   *PriceValue `json:"cash"`
		} `json:"balanceSheetStatements"`
	} `json:"balanceSheetHistory"`
	CashflowStatementHistory *struct {
		CashflowStatements []struct {
			EndDate                          *PriceValue `json:"endDate"`
			TotalCashFromOperatingActivities *PriceValue `json:"totalCashFromOperatingActivities"`
			CapitalExpenditures              *PriceValue `json:"capitalExpenditures"`
			FreeCashFlow                     *PriceValue `json:"freeCashFlow"`
			DividendsPaid                    *PriceValue `json:"dividendsPaid"`
		} `json:"cashflowStatements"`
	} `json:"cashflowStatementHistory"`
}) FinancialRatios {
	ratios := FinancialRatios{}

	// Extract from SummaryDetail
	if result.SummaryDetail != nil {
		ratios.PriceToEarningsRatio = result.SummaryDetail.TrailingPE
		ratios.PriceToBookRatio = result.SummaryDetail.PriceToBook
		ratios.PriceToSalesRatio = result.SummaryDetail.PriceToSalesTrailing12Months
		ratios.DividendRate = result.SummaryDetail.DividendRate
		ratios.DividendYield = result.SummaryDetail.DividendYield
	}

	// Extract from FinancialData if available
	if result.FinancialData != nil {
		ratios.ReturnOnEquity = result.FinancialData.ReturnOnEquity
		ratios.ReturnOnAssets = result.FinancialData.ReturnOnAssets
		ratios.GrossMargins = result.FinancialData.GrossMargins
		ratios.EbitdaMargins = result.FinancialData.EbitdaMargins
		ratios.OperatingMargins = result.FinancialData.OperatingMargins
		ratios.ProfitMargins = result.FinancialData.ProfitMargins
		ratios.CurrentRatio = result.FinancialData.CurrentRatio
		ratios.QuickRatio = result.FinancialData.QuickRatio
		ratios.DebtToEquity = result.FinancialData.DebtToEquity
		ratios.TotalDebtToCapital = result.FinancialData.TotalDebtToCapital
		ratios.EarningsGrowth = result.FinancialData.EarningsGrowth
		ratios.RevenueGrowth = result.FinancialData.RevenueGrowth
		ratios.EarningsPerShare = result.FinancialData.EarningsPerShare
		ratios.BookValuePerShare = result.FinancialData.BookValuePerShare
	}

	// Extract from DefaultKeyStatistics if available
	if result.DefaultKeyStatistics != nil {
		if ratios.PriceToEarningsRatio == nil {
			ratios.PriceToEarningsRatio = result.DefaultKeyStatistics.TrailingPE
		}
		if ratios.PriceToBookRatio == nil {
			ratios.PriceToBookRatio = result.DefaultKeyStatistics.PriceToBook
		}
	}

	return ratios
}

// extractFinancialSummary extracts financial summary data from the API response
func (t *Ticker) extractFinancialSummary(result struct {
	DefaultKeyStatistics *FinancialSummary `json:"defaultKeyStatistics"`
	FinancialData        *FinancialRatios  `json:"financialData"`
	SummaryDetail        *struct {
		MarketCap                    *PriceValue `json:"marketCap"`
		ForwardPE                    *PriceValue `json:"forwardPE"`
		TrailingPE                   *PriceValue `json:"trailingPE"`
		PriceToSalesTrailing12Months *PriceValue `json:"priceToSalesTrailing12Months"`
		PriceToBook                  *PriceValue `json:"priceToBook"`
		Beta                         *PriceValue `json:"beta"`
		DividendRate                 *PriceValue `json:"dividendRate"`
		DividendYield                *PriceValue `json:"dividendYield"`
	} `json:"summaryDetail"`
	IncomeStatementHistory *struct {
		IncomeStatementHistory []struct {
			EndDate         *PriceValue `json:"endDate"`
			TotalRevenue    *PriceValue `json:"totalRevenue"`
			GrossProfit     *PriceValue `json:"grossProfit"`
			OperatingIncome *PriceValue `json:"operatingIncome"`
			NetIncome       *PriceValue `json:"netIncome"`
			Ebitda          *PriceValue `json:"ebitda"`
		} `json:"incomeStatementHistory"`
	} `json:"incomeStatementHistory"`
	BalanceSheetHistory *struct {
		BalanceSheetStatements []struct {
			EndDate                *PriceValue `json:"endDate"`
			TotalAssets            *PriceValue `json:"totalAssets"`
			TotalLiab              *PriceValue `json:"totalLiab"`
			TotalStockholderEquity *PriceValue `json:"totalStockholderEquity"`
			TotalDebt              *PriceValue `json:"totalDebt"`
			Cash                   *PriceValue `json:"cash"`
		} `json:"balanceSheetStatements"`
	} `json:"balanceSheetHistory"`
	CashflowStatementHistory *struct {
		CashflowStatements []struct {
			EndDate                          *PriceValue `json:"endDate"`
			TotalCashFromOperatingActivities *PriceValue `json:"totalCashFromOperatingActivities"`
			CapitalExpenditures              *PriceValue `json:"capitalExpenditures"`
			FreeCashFlow                     *PriceValue `json:"freeCashFlow"`
			DividendsPaid                    *PriceValue `json:"dividendsPaid"`
		} `json:"cashflowStatements"`
	} `json:"cashflowStatementHistory"`
}) FinancialSummary {
	summary := FinancialSummary{}

	// Prioritize DefaultKeyStatistics
	if result.DefaultKeyStatistics != nil {
		summary = *result.DefaultKeyStatistics
	}

	// Override/supplement with SummaryDetail data
	if result.SummaryDetail != nil {
		if summary.MarketCap == nil {
			summary.MarketCap = result.SummaryDetail.MarketCap
		}
		if summary.ForwardPE == nil {
			summary.ForwardPE = result.SummaryDetail.ForwardPE
		}
		if summary.TrailingPE == nil {
			summary.TrailingPE = result.SummaryDetail.TrailingPE
		}
		if summary.PriceToSalesTrailing12Months == nil {
			summary.PriceToSalesTrailing12Months = result.SummaryDetail.PriceToSalesTrailing12Months
		}
		if summary.PriceToBook == nil {
			summary.PriceToBook = result.SummaryDetail.PriceToBook
		}
		if summary.Beta == nil {
			summary.Beta = result.SummaryDetail.Beta
		}
	}

	return summary
}

// extractIncomeStatement extracts the latest income statement data
func (t *Ticker) extractIncomeStatement(result struct {
	DefaultKeyStatistics *FinancialSummary `json:"defaultKeyStatistics"`
	FinancialData        *FinancialRatios  `json:"financialData"`
	SummaryDetail        *struct {
		MarketCap                    *PriceValue `json:"marketCap"`
		ForwardPE                    *PriceValue `json:"forwardPE"`
		TrailingPE                   *PriceValue `json:"trailingPE"`
		PriceToSalesTrailing12Months *PriceValue `json:"priceToSalesTrailing12Months"`
		PriceToBook                  *PriceValue `json:"priceToBook"`
		Beta                         *PriceValue `json:"beta"`
		DividendRate                 *PriceValue `json:"dividendRate"`
		DividendYield                *PriceValue `json:"dividendYield"`
	} `json:"summaryDetail"`
	IncomeStatementHistory *struct {
		IncomeStatementHistory []struct {
			EndDate         *PriceValue `json:"endDate"`
			TotalRevenue    *PriceValue `json:"totalRevenue"`
			GrossProfit     *PriceValue `json:"grossProfit"`
			OperatingIncome *PriceValue `json:"operatingIncome"`
			NetIncome       *PriceValue `json:"netIncome"`
			Ebitda          *PriceValue `json:"ebitda"`
		} `json:"incomeStatementHistory"`
	} `json:"incomeStatementHistory"`
	BalanceSheetHistory *struct {
		BalanceSheetStatements []struct {
			EndDate                *PriceValue `json:"endDate"`
			TotalAssets            *PriceValue `json:"totalAssets"`
			TotalLiab              *PriceValue `json:"totalLiab"`
			TotalStockholderEquity *PriceValue `json:"totalStockholderEquity"`
			TotalDebt              *PriceValue `json:"totalDebt"`
			Cash                   *PriceValue `json:"cash"`
		} `json:"balanceSheetStatements"`
	} `json:"balanceSheetHistory"`
	CashflowStatementHistory *struct {
		CashflowStatements []struct {
			EndDate                          *PriceValue `json:"endDate"`
			TotalCashFromOperatingActivities *PriceValue `json:"totalCashFromOperatingActivities"`
			CapitalExpenditures              *PriceValue `json:"capitalExpenditures"`
			FreeCashFlow                     *PriceValue `json:"freeCashFlow"`
			DividendsPaid                    *PriceValue `json:"dividendsPaid"`
		} `json:"cashflowStatements"`
	} `json:"cashflowStatementHistory"`
}) IncomeStatement {
	income := IncomeStatement{}

	if result.IncomeStatementHistory != nil && len(result.IncomeStatementHistory.IncomeStatementHistory) > 0 {
		// Get the most recent income statement (first in the array)
		latest := result.IncomeStatementHistory.IncomeStatementHistory[0]
		income.TotalRevenue = latest.TotalRevenue
		income.GrossProfit = latest.GrossProfit
		income.OperatingIncome = latest.OperatingIncome
		income.NetIncome = latest.NetIncome
		income.Ebitda = latest.Ebitda
	}

	// Add EPS data from financial ratios if available
	if result.FinancialData != nil {
		income.EarningsPerShare = result.FinancialData.EarningsPerShare
	}

	return income
}

// extractBalanceSheet extracts the latest balance sheet data
func (t *Ticker) extractBalanceSheet(result struct {
	DefaultKeyStatistics *FinancialSummary `json:"defaultKeyStatistics"`
	FinancialData        *FinancialRatios  `json:"financialData"`
	SummaryDetail        *struct {
		MarketCap                    *PriceValue `json:"marketCap"`
		ForwardPE                    *PriceValue `json:"forwardPE"`
		TrailingPE                   *PriceValue `json:"trailingPE"`
		PriceToSalesTrailing12Months *PriceValue `json:"priceToSalesTrailing12Months"`
		PriceToBook                  *PriceValue `json:"priceToBook"`
		Beta                         *PriceValue `json:"beta"`
		DividendRate                 *PriceValue `json:"dividendRate"`
		DividendYield                *PriceValue `json:"dividendYield"`
	} `json:"summaryDetail"`
	IncomeStatementHistory *struct {
		IncomeStatementHistory []struct {
			EndDate         *PriceValue `json:"endDate"`
			TotalRevenue    *PriceValue `json:"totalRevenue"`
			GrossProfit     *PriceValue `json:"grossProfit"`
			OperatingIncome *PriceValue `json:"operatingIncome"`
			NetIncome       *PriceValue `json:"netIncome"`
			Ebitda          *PriceValue `json:"ebitda"`
		} `json:"incomeStatementHistory"`
	} `json:"incomeStatementHistory"`
	BalanceSheetHistory *struct {
		BalanceSheetStatements []struct {
			EndDate                *PriceValue `json:"endDate"`
			TotalAssets            *PriceValue `json:"totalAssets"`
			TotalLiab              *PriceValue `json:"totalLiab"`
			TotalStockholderEquity *PriceValue `json:"totalStockholderEquity"`
			TotalDebt              *PriceValue `json:"totalDebt"`
			Cash                   *PriceValue `json:"cash"`
		} `json:"balanceSheetStatements"`
	} `json:"balanceSheetHistory"`
	CashflowStatementHistory *struct {
		CashflowStatements []struct {
			EndDate                          *PriceValue `json:"endDate"`
			TotalCashFromOperatingActivities *PriceValue `json:"totalCashFromOperatingActivities"`
			CapitalExpenditures              *PriceValue `json:"capitalExpenditures"`
			FreeCashFlow                     *PriceValue `json:"freeCashFlow"`
			DividendsPaid                    *PriceValue `json:"dividendsPaid"`
		} `json:"cashflowStatements"`
	} `json:"cashflowStatementHistory"`
}) BalanceSheet {
	balance := BalanceSheet{}

	if result.BalanceSheetHistory != nil && len(result.BalanceSheetHistory.BalanceSheetStatements) > 0 {
		// Get the most recent balance sheet (first in the array)
		latest := result.BalanceSheetHistory.BalanceSheetStatements[0]
		balance.TotalAssets = latest.TotalAssets
		balance.TotalLiabilities = latest.TotalLiab
		balance.TotalEquity = latest.TotalStockholderEquity
		balance.TotalDebt = latest.TotalDebt
		balance.Cash = latest.Cash
	}

	// Add book value per share from financial ratios if available
	if result.FinancialData != nil {
		balance.BookValuePerShare = result.FinancialData.BookValuePerShare
	}

	return balance
}

// extractCashFlow extracts the latest cash flow statement data
func (t *Ticker) extractCashFlow(result struct {
	DefaultKeyStatistics *FinancialSummary `json:"defaultKeyStatistics"`
	FinancialData        *FinancialRatios  `json:"financialData"`
	SummaryDetail        *struct {
		MarketCap                    *PriceValue `json:"marketCap"`
		ForwardPE                    *PriceValue `json:"forwardPE"`
		TrailingPE                   *PriceValue `json:"trailingPE"`
		PriceToSalesTrailing12Months *PriceValue `json:"priceToSalesTrailing12Months"`
		PriceToBook                  *PriceValue `json:"priceToBook"`
		Beta                         *PriceValue `json:"beta"`
		DividendRate                 *PriceValue `json:"dividendRate"`
		DividendYield                *PriceValue `json:"dividendYield"`
	} `json:"summaryDetail"`
	IncomeStatementHistory *struct {
		IncomeStatementHistory []struct {
			EndDate         *PriceValue `json:"endDate"`
			TotalRevenue    *PriceValue `json:"totalRevenue"`
			GrossProfit     *PriceValue `json:"grossProfit"`
			OperatingIncome *PriceValue `json:"operatingIncome"`
			NetIncome       *PriceValue `json:"netIncome"`
			Ebitda          *PriceValue `json:"ebitda"`
		} `json:"incomeStatementHistory"`
	} `json:"incomeStatementHistory"`
	BalanceSheetHistory *struct {
		BalanceSheetStatements []struct {
			EndDate                *PriceValue `json:"endDate"`
			TotalAssets            *PriceValue `json:"totalAssets"`
			TotalLiab              *PriceValue `json:"totalLiab"`
			TotalStockholderEquity *PriceValue `json:"totalStockholderEquity"`
			TotalDebt              *PriceValue `json:"totalDebt"`
			Cash                   *PriceValue `json:"cash"`
		} `json:"balanceSheetStatements"`
	} `json:"balanceSheetHistory"`
	CashflowStatementHistory *struct {
		CashflowStatements []struct {
			EndDate                          *PriceValue `json:"endDate"`
			TotalCashFromOperatingActivities *PriceValue `json:"totalCashFromOperatingActivities"`
			CapitalExpenditures              *PriceValue `json:"capitalExpenditures"`
			FreeCashFlow                     *PriceValue `json:"freeCashFlow"`
			DividendsPaid                    *PriceValue `json:"dividendsPaid"`
		} `json:"cashflowStatements"`
	} `json:"cashflowStatementHistory"`
}) CashFlow {
	cashflow := CashFlow{}

	if result.CashflowStatementHistory != nil && len(result.CashflowStatementHistory.CashflowStatements) > 0 {
		// Get the most recent cash flow statement (first in the array)
		latest := result.CashflowStatementHistory.CashflowStatements[0]
		cashflow.OperatingCashFlow = latest.TotalCashFromOperatingActivities
		cashflow.CapitalExpenditures = latest.CapitalExpenditures
		cashflow.FreeCashFlow = latest.FreeCashFlow
		cashflow.DividendsPaid = latest.DividendsPaid
	}

	return cashflow
}

// transformHistoricalData converts YahooHistoryResponse into a map of PriceData keyed by date/time
func transformHistoricalData(data YahooHistoryResponse, interval string) map[string]PriceData {
	d := make(map[string]PriceData)
	if len(data.Chart.Result) == 0 {
		return d
	}

	result := data.Chart.Result[0]
	for i, timestamp := range result.Timestamp {
		t := time.Unix(timestamp, 0)
		var key string
		if strings.HasSuffix(interval, "d") || strings.HasSuffix(interval, "wk") || strings.HasSuffix(interval, "mo") {
			key = t.Format("2006-01-02")
		} else {
			key = t.Format("2006-01-02 15:04:05")
		}

		// Ensure we have quote data
		if len(result.Indicators.Quote) > 0 {
			quote := result.Indicators.Quote[0]

			// Check bounds to avoid index out of range
			if i < len(quote.Open) && i < len(quote.High) && i < len(quote.Low) && i < len(quote.Close) && i < len(quote.Volume) {
				d[key] = PriceData{
					Open:   quote.Open[i],
					High:   quote.High[i],
					Low:    quote.Low[i],
					Close:  quote.Close[i],
					Volume: quote.Volume[i],
				}
			}
		}
	}
	return d
}

// extractDividendInfo extracts dividend information from the API response
func (t *Ticker) extractDividendInfo(result struct {
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
}) DividendInfo {
	dividend := DividendInfo{}

	// Prioritize SummaryDetail data
	if result.SummaryDetail != nil {
		dividend.DividendRate = result.SummaryDetail.DividendRate
		dividend.DividendYield = result.SummaryDetail.DividendYield
		dividend.ExDividendDate = result.SummaryDetail.ExDividendDate
		dividend.DividendDate = result.SummaryDetail.DividendDate
		dividend.PayoutRatio = result.SummaryDetail.PayoutRatio
		dividend.FiveYearAvgDividendYield = result.SummaryDetail.FiveYearAvgDividendYield
	}

	// Fallback to DefaultKeyStatistics if SummaryDetail is missing data
	if result.DefaultKeyStatistics != nil {
		if dividend.DividendRate == nil {
			dividend.DividendRate = result.DefaultKeyStatistics.DividendRate
		}
		if dividend.DividendYield == nil {
			dividend.DividendYield = result.DefaultKeyStatistics.DividendYield
		}
		if dividend.PayoutRatio == nil {
			dividend.PayoutRatio = result.DefaultKeyStatistics.PayoutRatio
		}
		if dividend.FiveYearAvgDividendYield == nil {
			dividend.FiveYearAvgDividendYield = result.DefaultKeyStatistics.FiveYearAvgDividendYield
		}
	}

	// Get dividends paid from cash flow statements
	if result.CashflowStatementHistory != nil && len(result.CashflowStatementHistory.CashflowStatements) > 0 {
		// Get the most recent cash flow statement
		latest := result.CashflowStatementHistory.CashflowStatements[0]
		dividend.DividendsPaid = latest.DividendsPaid
	}

	return dividend
}
