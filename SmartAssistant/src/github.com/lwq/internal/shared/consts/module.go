package consts

type AppModule int

const (
	Unknown AppModule = iota
	HelpCenter
	ReqSort
	LiveChat
)

var moduleOptions = []AppModule{HelpCenter, ReqSort, LiveChat}

func GetModuleOption() []AppModule {
    return moduleOptions
}
