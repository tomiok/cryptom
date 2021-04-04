package blocks

type Header struct {
	Version       int64
	PrevBlockId   string
	PayloadDigest string
	Transactions  string
}
