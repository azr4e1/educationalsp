package baseprotocol

// Base Types
type Integer int64
type UInteger uint64
type Decimal float64
type LSPAny interface {
	LSPObject | LSPArray | string | Integer | UInteger | Decimal | bool
}

type LSPArray = []LSPAny
type LSPObject map[string]LSPAny
