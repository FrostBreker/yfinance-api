package yfinance_api

type YahooInfoResponse struct {
	QuoteSummary struct {
		Result []struct {
			Price YahooTickerInfo `json:"price"`
		} `json:"result"`
		Error interface{} `json:"error"`
	} `json:"quoteSummary"`
}

// PriceValue represents a price value with raw and formatted representations
type PriceValue struct {
	Raw     float64 `json:"raw"`
	Fmt     string  `json:"fmt"`
	LongFmt string  `json:"longFmt,omitempty"`
}

// YahooTickerInfo --> Struct to hold key metadata about the ticker
type YahooTickerInfo struct {
	MaxAge                     int         `json:"maxAge"`
	PreMarketChange            *PriceValue `json:"preMarketChange"`
	PreMarketPrice             *PriceValue `json:"preMarketPrice"`
	PreMarketSource            string      `json:"preMarketSource"`
	PostMarketChangePercent    *PriceValue `json:"postMarketChangePercent"`
	PostMarketChange           *PriceValue `json:"postMarketChange"`
	PostMarketTime             int64       `json:"postMarketTime"`
	PostMarketPrice            *PriceValue `json:"postMarketPrice"`
	PostMarketSource           string      `json:"postMarketSource"`
	RegularMarketChangePercent *PriceValue `json:"regularMarketChangePercent"`
	RegularMarketChange        *PriceValue `json:"regularMarketChange"`
	RegularMarketTime          int64       `json:"regularMarketTime"`
	PriceHint                  *PriceValue `json:"priceHint"`
	RegularMarketPrice         *PriceValue `json:"regularMarketPrice"`
	RegularMarketDayHigh       *PriceValue `json:"regularMarketDayHigh"`
	RegularMarketDayLow        *PriceValue `json:"regularMarketDayLow"`
	RegularMarketVolume        *PriceValue `json:"regularMarketVolume"`
	AverageDailyVolume10Day    *PriceValue `json:"averageDailyVolume10Day"`
	AverageDailyVolume3Month   *PriceValue `json:"averageDailyVolume3Month"`
	RegularMarketPreviousClose *PriceValue `json:"regularMarketPreviousClose"`
	RegularMarketSource        string      `json:"regularMarketSource"`
	RegularMarketOpen          *PriceValue `json:"regularMarketOpen"`
	StrikePrice                *PriceValue `json:"strikePrice"`
	OpenInterest               *PriceValue `json:"openInterest"`
	Exchange                   string      `json:"exchange"`
	ExchangeName               string      `json:"exchangeName"`
	ExchangeDataDelayedBy      int         `json:"exchangeDataDelayedBy"`
	MarketState                string      `json:"marketState"`
	QuoteType                  string      `json:"quoteType"`
	Symbol                     string      `json:"symbol"`
	UnderlyingSymbol           *string     `json:"underlyingSymbol"`
	ShortName                  string      `json:"shortName"`
	LongName                   string      `json:"longName"`
	Currency                   string      `json:"currency"`
	QuoteSourceName            string      `json:"quoteSourceName"`
	CurrencySymbol             string      `json:"currencySymbol"`
	FromCurrency               *string     `json:"fromCurrency"`
	ToCurrency                 *string     `json:"toCurrency"`
	LastMarket                 *string     `json:"lastMarket"`
	Volume24Hr                 *PriceValue `json:"volume24Hr"`
	VolumeAllCurrencies        *PriceValue `json:"volumeAllCurrencies"`
	CirculatingSupply          *PriceValue `json:"circulatingSupply"`
	MarketCap                  *PriceValue `json:"marketCap"`
}

// PriceData represents historical price and volume data for a specific time period
type PriceData struct {
	Open   *float64 `json:"open"`
	High   *float64 `json:"high"`
	Low    *float64 `json:"low"`
	Close  *float64 `json:"close"`
	Volume *int64   `json:"volume"`
}

// Query represents the query parameters for historical data requests
type Query struct {
	Range    string `json:"range"`
	Interval string `json:"interval"`
	Start    string `json:"start"`
	End      string `json:"end"`
}

// SetDefault sets default values for the query parameters
func (q *Query) SetDefault() {
	if q.Interval == "" {
		q.Interval = "1d"
	}
	if q.Range == "" {
		q.Range = "1y"
	}
}

// History represents a historical data request handler
type History struct {
	Client *Client `json:"client"`
	Query  Query   `json:"query"`
}

// YahooHistoryResponse represents the response from Yahoo Finance historical data API
type YahooHistoryResponse struct {
	Chart struct {
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
	} `json:"chart"`
}

// NewsItem represents a single news article from Yahoo Finance
type NewsItem struct {
	UUID                string `json:"uuid"`
	Title               string `json:"title"`
	Publisher           string `json:"publisher"`
	Link                string `json:"link"`
	ProviderPublishTime int64  `json:"providerPublishTime"`
	Type                string `json:"type"`
	Thumbnail           *struct {
		Resolutions []struct {
			URL    string `json:"url"`
			Width  int    `json:"width"`
			Height int    `json:"height"`
			Tag    string `json:"tag"`
		} `json:"resolutions"`
	} `json:"thumbnail,omitempty"`
	RelatedTickers []string `json:"relatedTickers"`
}

// YahooNewsResponse represents the response from Yahoo Finance news API
type YahooNewsResponse struct {
	News struct {
		Stream []NewsItem `json:"stream"`
	} `json:"news"`
}

// FinancialRatios represents key financial ratios for fundamental analysis
type FinancialRatios struct {
	// Valuation Ratios
	PriceToEarningsRatio *PriceValue `json:"priceToEarningsRatio"`
	PriceToBookRatio     *PriceValue `json:"priceToBookRatio"`
	PriceToSalesRatio    *PriceValue `json:"priceToSalesRatio"`
	PriceToEbitda        *PriceValue `json:"priceToEbitda"`
	EnterpriseToRevenue  *PriceValue `json:"enterpriseToRevenue"`
	EnterpriseToEbitda   *PriceValue `json:"enterpriseToEbitda"`

	// Profitability Ratios
	ReturnOnEquity   *PriceValue `json:"returnOnEquity"`
	ReturnOnAssets   *PriceValue `json:"returnOnAssets"`
	GrossMargins     *PriceValue `json:"grossMargins"`
	EbitdaMargins    *PriceValue `json:"ebitdaMargins"`
	OperatingMargins *PriceValue `json:"operatingMargins"`
	ProfitMargins    *PriceValue `json:"profitMargins"`

	// Liquidity Ratios
	CurrentRatio *PriceValue `json:"currentRatio"`
	QuickRatio   *PriceValue `json:"quickRatio"`

	// Leverage Ratios
	DebtToEquity       *PriceValue `json:"debtToEquity"`
	TotalDebtToCapital *PriceValue `json:"totalDebtToCapital"`

	// Growth Ratios
	EarningsGrowth *PriceValue `json:"earningsGrowth"`
	RevenueGrowth  *PriceValue `json:"revenueGrowth"`

	// Per Share Data
	EarningsPerShare  *PriceValue `json:"earningsPerShare"`
	BookValuePerShare *PriceValue `json:"bookValuePerShare"`
	DividendRate      *PriceValue `json:"dividendRate"`
	DividendYield     *PriceValue `json:"dividendYield"`
}

// FinancialSummary represents key financial metrics summary
type FinancialSummary struct {
	MaxAge                       int         `json:"maxAge"`
	MarketCap                    *PriceValue `json:"marketCap"`
	EnterpriseValue              *PriceValue `json:"enterpriseValue"`
	ForwardPE                    *PriceValue `json:"forwardPE"`
	TrailingPE                   *PriceValue `json:"trailingPE"`
	PegRatio                     *PriceValue `json:"pegRatio"`
	PriceToSalesTrailing12Months *PriceValue `json:"priceToSalesTrailing12Months"`
	PriceToBook                  *PriceValue `json:"priceToBook"`
	Beta                         *PriceValue `json:"beta"`
	FiftyTwoWeekLow              *PriceValue `json:"fiftyTwoWeekLow"`
	FiftyTwoWeekHigh             *PriceValue `json:"fiftyTwoWeekHigh"`
	FiftyDayAverage              *PriceValue `json:"fiftyDayAverage"`
	TwoHundredDayAverage         *PriceValue `json:"twoHundredDayAverage"`
}

// IncomeStatement represents income statement data
type IncomeStatement struct {
	TotalRevenue     *PriceValue `json:"totalRevenue"`
	GrossProfit      *PriceValue `json:"grossProfit"`
	OperatingIncome  *PriceValue `json:"operatingIncome"`
	NetIncome        *PriceValue `json:"netIncome"`
	Ebitda           *PriceValue `json:"ebitda"`
	EarningsPerShare *PriceValue `json:"earningsPerShare"`
	DilutedEPS       *PriceValue `json:"dilutedEPS"`
}

// BalanceSheet represents balance sheet data
type BalanceSheet struct {
	TotalAssets       *PriceValue `json:"totalAssets"`
	TotalLiabilities  *PriceValue `json:"totalLiabilities"`
	TotalEquity       *PriceValue `json:"totalStockholderEquity"`
	TotalDebt         *PriceValue `json:"totalDebt"`
	Cash              *PriceValue `json:"cash"`
	ShortTermDebt     *PriceValue `json:"shortTermDebt"`
	LongTermDebt      *PriceValue `json:"longTermDebt"`
	BookValuePerShare *PriceValue `json:"bookValuePerShare"`
}

// CashFlow represents cash flow statement data
type CashFlow struct {
	OperatingCashFlow   *PriceValue `json:"operatingCashFlow"`
	FreeCashFlow        *PriceValue `json:"freeCashFlow"`
	CapitalExpenditures *PriceValue `json:"capitalExpenditures"`
	DividendsPaid       *PriceValue `json:"dividendsPaid"`
}

// FinancialData represents comprehensive financial data for a ticker
type FinancialData struct {
	Ratios          FinancialRatios  `json:"ratios"`
	Summary         FinancialSummary `json:"summary"`
	IncomeStatement IncomeStatement  `json:"incomeStatement"`
	BalanceSheet    BalanceSheet     `json:"balanceSheet"`
	CashFlow        CashFlow         `json:"cashFlow"`
}

// YahooFinancialResponse represents the response from Yahoo Finance financial APIs
type YahooFinancialResponse struct {
	QuoteSummary struct {
		Result []struct {
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
		} `json:"result"`
		Error interface{} `json:"error"`
	} `json:"quoteSummary"`
}
