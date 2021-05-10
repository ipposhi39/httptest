package jsonrpc

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/friendsofgo/errors"
	"github.com/omeroid/wdc/backend/pkg/errof"
	"github.com/volatiletech/null/v8"
)

const (
	jsonrpc = "2.0"
)

var (
	batchRequestOpenToken  = '['
	batchRequestCloseToken = ']'
)

// Request : request
type Request struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Headers http.Header     `json:"_"`
	Params  json.RawMessage `json:"params"`
	ID      interface{}     `json:"id"`
}

// Return : JSONRPCReturn
type Return struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Error   *Error      `json:"error,omitempty"`
	Result  interface{} `json:"result,omitempty"`
}

// Credential :
type Credential struct {
	UID           string      `json:"uid"`
	Email         string      `json:"email"`
	EmailVerified bool        `json:"email_verified"`
	Picture       null.String `json:"picture"`
	Name          null.String `json:"name"`
}

// Parse : ParseJSONRPC
func Parse(r *http.Request) (requests []*Request, err error) {
	buf := bytes.NewBuffer(make([]byte, 0, r.ContentLength))
	if _, err := buf.ReadFrom(r.Body); err != nil {
		return nil, errors.WithStack(errof.ErrInvalidRequest)
	}
	defer r.Body.Close()

	if buf.Len() == 0 {
		return nil, errors.WithStack(errof.ErrInvalidRequest)
	}
	return parse(buf)
}

func parse(buf *bytes.Buffer) (requests []*Request, err error) {
	// read first rune
	f, _, err := buf.ReadRune()
	if err != nil {
		return nil, errors.WithStack(errof.ErrInvalidRequest)
	}
	if err := buf.UnreadRune(); err != nil {
		return nil, errors.WithStack(errof.ErrInvalidRequest)
	}

	// not batch
	if f != batchRequestOpenToken {
		r := &Request{}
		if err := json.Unmarshal(buf.Bytes(), r); err != nil {
			switch err.(type) {
			case *json.SyntaxError:
				return nil, errors.WithStack(errof.ErrParse)
			case *json.UnmarshalTypeError:
				return nil, errors.WithStack(errof.ErrInvalidRequest)
			}
		}
		return []*Request{r}, nil
	}

	// batch
	d := json.NewDecoder(buf)

	// read open bracket
	t, err := d.Token()
	if err != nil {
		return nil, errors.Wrap(errof.ErrParse, "Failed to read open bracket")
	}

	if t != json.Delim(batchRequestOpenToken) {
		return nil, errors.Wrap(errof.ErrParse, "Invalid open token")
	}

	for d.More() {
		r := &Request{}
		if err = d.Decode(r); err != nil {
			switch err.(type) {
			case *json.SyntaxError:
				return nil, errors.Wrap(errof.ErrParse, "Failed to decode batch request")
			case *json.UnmarshalTypeError:
				return nil, errors.Wrap(errof.ErrInvalidRequest, err.Error())
			}
		}
		if r.JSONRPC == "" || r.Method == "" || r.ID == nil {
			return nil, errors.Wrap(errof.ErrInvalidRequest, "'jsonrpc', 'method' and 'id' are required")
		}
		requests = append(requests, r)
	}
	// read closing bracket
	t, err = d.Token()
	if err != nil || t != json.Delim(batchRequestCloseToken) {
		return nil, errors.Wrap(errof.ErrParse, "Invalid closing token")
	}
	return requests, err
}

// WriteResponses writes responses
func WriteResponses(w io.Writer, returns ...*Return) (err error) {
	for _, r := range returns {
		r.JSONRPC = jsonrpc
	}

	if len(returns) == 1 {
		if err = json.NewEncoder(w).Encode(returns[0]); err != nil {
			return err
		}
		return nil
	}

	if 1 < len(returns) {
		if err = json.NewEncoder(w).Encode(returns); err != nil {
			return err
		}
		return nil
	}

	result := &Return{JSONRPC: jsonrpc, Error: ErrInternal()}
	if err = json.NewEncoder(w).Encode(result); err != nil {
		return err
	}
	return nil
}
