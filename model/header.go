package model

type Header struct {
	version       int64
	prevBlockId   string
	payloadDigest string
	transactions  string
}
