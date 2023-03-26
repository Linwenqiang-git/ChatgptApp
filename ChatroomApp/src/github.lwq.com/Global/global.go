package global

import (
	"context"

	. "github.lwq.com/Completion"
)

var (
	Configuration *GlobalConfig = LoadGlobalConfig()
	CtClient      *CompletionClient
	Ctx           context.Context
)
