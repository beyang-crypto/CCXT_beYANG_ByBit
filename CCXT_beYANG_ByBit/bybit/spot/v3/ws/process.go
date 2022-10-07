package ws

func (b *ByBitWS) processBookTicker(symbol string, data BookTicker) {
	b.Emit(ChannelBookTicker, symbol, data)
}
