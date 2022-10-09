package ws

func (b *ByBitWS) processBookTicker(name string, symbol string, data BookTicker) {
	b.Emit(ChannelBookTicker, name, symbol, data)
}
