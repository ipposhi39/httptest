package errof

// SkipErr :
type SkipErr string

// UserErr :
type UserErr string

// InternalErr :
type InternalErr string

func (e SkipErr) Error() string {
	return string(e)
}

func (e UserErr) Error() (msg string) {
	var ok bool
	if msg, ok = ErrCodeNames[e]; !ok {
		return string(e)
	}
	return msg
}

func (e InternalErr) Error() (msg string) {
	var ok bool
	if msg, ok = InternalErrCodeNames[e]; !ok {
		return string(e)
	}
	return msg
}

// InternalErrCodeNames :
var InternalErrCodeNames = map[InternalErr]string{
	ErrInternal: "内部エラーが発生しました",
	ErrServer:   "サーバエラーが発生しました",
	ErrFirebase: "認証システムでエラーが発生しました",
	ErrDatabase: "データベースでの不整合が発生しました",
}

// ErrCodeNames :
var ErrCodeNames = map[UserErr]string{
	ErrAuthentication:   "認証に失敗しました",
	ErrInvalidMethod:    "メソッド名が不正です",
	ErrInvalidParams:    "パラメータが不正です",
	ErrMethodNotFound:   "メソッドが存在しません",
	ErrParse:            "パラメータの変換に失敗しました",
	ErrInvalidRequest:   "リクエストが不正です",
	ErrParameter:        "パラメータの形式が違います",
	ErrTooLongParameter: "パラメータが長すぎます",
	ErrHTTP:             "HTTPでエラーが発生しました",
	ErrExpired:          "トークンの有効期限が切れています",

	ErrNoOrg: "オーガニゼーションが見つかりません",
}

// Error 定義
var (
	// JSON-RPC Error object
	// http://www.jsonrpc.org/specification#error_object

	ErrInternal InternalErr = "ErrInternal"
	ErrServer   InternalErr = "ErrServer"
	ErrFirebase InternalErr = "ErrFirebase"
	ErrDatabase InternalErr = "ErrDatabase"

	ErrParse            UserErr = "ErrParse"
	ErrInvalidRequest   UserErr = "ErrInvalidRequest"
	ErrMethodNotFound   UserErr = "ErrMethodNotFound"
	ErrInvalidParams    UserErr = "ErrInvalidParams"
	ErrAuthentication   UserErr = "ErrAuthentication"
	ErrInvalidMethod    UserErr = "ErrInvalidMethod"
	ErrParameter        UserErr = "ErrParameter"
	ErrTooLongParameter UserErr = "ErrTooLongParameter"
	ErrHTTP             UserErr = "ErrHTTP"
	ErrExpired          UserErr = "ErrExpired"

	ErrNoOrg UserErr = "ErrNoOrg"
)
