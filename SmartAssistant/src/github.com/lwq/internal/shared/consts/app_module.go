package consts

type AppModule int

const (
	OpenaiKey AppModule = iota
	HelpCenter
	ReqSort
	LiveChat
)

var moduleOptions = []AppModule{OpenaiKey,HelpCenter, ReqSort, LiveChat}

func GetModuleOption() []AppModule {
    return moduleOptions
}
