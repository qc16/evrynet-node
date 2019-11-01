package types

// FakeHeader update fake info to block
func (b *Block) FakeHeader(fakeHeader *Header) *Block {
	b.header = fakeHeader
	return b
}
