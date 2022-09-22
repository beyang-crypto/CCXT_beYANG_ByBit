package v3

func (b *ByBitWS) processBookTicker(symbol string, data BookTicker) {
	b.Emit(ChannelBookTicker, symbol, data)
}
