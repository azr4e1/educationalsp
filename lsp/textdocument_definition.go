package lsp

type DefinitionTextDocumentRequest struct {
	Request
	Param DefinitionParams `json:"params"`
}

type DefinitionParams struct {
	TextDocumentPositionParams
}

type DefinitionResponse struct {
	Response
	Result Location `json:"result"`
}
