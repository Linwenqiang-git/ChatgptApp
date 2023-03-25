package global

import (
	"context"
	. "github.lwq.com/Completion"
)

var (
	Chatgpt_token = "your api key"
	CtClient      *CompletionClient
	Ctx           context.Context
)
