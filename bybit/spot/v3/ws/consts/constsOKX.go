package consts

const (
	OKXHostPublicWebSocket     = "wss://ws.okx.com:8443/ws/v5/public"
	OKXHostPrivateWebSocket    = "wss://ws.okx.com:8443/ws/v5/private"
	OKXHostPublicWebSocketAWS  = "wss://wsaws.okx.com:8443/ws/v5/public"
	OKXHostPrivateWebSocketAWS = "wss://wsaws.okx.com:8443/ws/v5/private"
)

const (
	OKXChannelBalanceAndPosition = "balance_and_position"
	OKXChannelRfqs               = "rfqs"
	OKXChannelQuotes             = "quotes"
	OKXChannelTicker             = "tickers"
	OKXChannelOptSummary         = "opt-summary"
	OKXChannelAccount            = "account"
	OKXChannelAccountGreeks      = "account-greeks"
	OKXChannelGridPositions      = "grid-positions"
	OKXChannelGridSubOrders      = "grid-sub-orders"
	OKXChannelLiquidationWarning = "liquidation-warning"
	OKXChannelGridOrdersSpot     = "grid-orders-spot"
	OKXChannelGridOrdersContract = "grid-orders-contract"
	OKXChannelGridOrdersMoon     = "grid-orders-moon"
	OKXChannelAlgoAdvance        = "algo-advance"
	OKXChannelEstimatedPrice     = "estimated-price"
	OKXChannelPositions          = "positions"
	OKXChannelOrders             = "orders"
	OKXChannelOrdersAlgo         = "orders-algo"
)

const (
	OKXInstTypeSpot    = "SPOT"
	OKXInstTypeMargin  = "MARGIN"
	OKXInstTypeSwap    = "SWAP"
	OKXInstTypeFutures = "FUTURES"
	OKXInstTypeOption  = "OPTION"
	OKXInstTypeAny     = "ANY"
)
