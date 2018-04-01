package internal

type DefaultRequest struct {
	Action string `json:"action"`
	App    string `json:"app,omitempty"`
}

type DefaultResponse struct {
	Error string `json:"error"`
}

//--------------

type ProcessBlockRequest struct {
	Block string `json:"block"`
	DefaultRequest
}

type ProcessBlockResponse struct {
	Hash string `json:"hash"`
	DefaultResponse
}

//--------------

type RetrieveBlockRequest struct {
	Hash string `json:"hash"`
	DefaultRequest
}

type RetrieveBlockResponse struct {
	Contents string `json:"contents"`
	DefaultResponse
}
