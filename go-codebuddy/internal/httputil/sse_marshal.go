package httputil

// SSEMarshaler is implemented by hot-path SSE payloads that can supply wire bytes
// without going through json.Marshal on the outer envelope.
type SSEMarshaler interface {
	MarshalSSE() ([]byte, error)
}
