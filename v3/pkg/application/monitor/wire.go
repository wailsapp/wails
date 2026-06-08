package monitor

// Wire protocol. Every line the server sends is an Envelope (NDJSON). Every line
// the client sends is a Request (NDJSON). Wrapping the trace stream in an
// envelope lets the same connection carry on-demand snapshot replies.

// MsgType discriminates server-sent envelopes.
type MsgType string

const (
	MsgTrace    MsgType = "trace"
	MsgSnapshot MsgType = "snapshot"
	MsgSample   MsgType = "sample"
)

// Envelope is one server→client message.
type Envelope struct {
	Type     MsgType   `json:"t"`
	Trace    *Trace    `json:"trace,omitempty"`
	Snapshot *Snapshot `json:"snapshot,omitempty"`
	Sample   *Sample   `json:"sample,omitempty"`
}

// RequestType discriminates client-sent requests.
type RequestType string

const (
	ReqDescribe RequestType = "describe"
)

// Request is one client→server message.
type Request struct {
	Type RequestType `json:"req"`
}
