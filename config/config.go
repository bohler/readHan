package config

var DefaultKeys = map[string]map[string]string{
	"global": map[string]string{
		"CtrlC": "quit",
		"Tab":   "nextView",
		"CtrlJ": "nextView",
		"CtrlK": "prevView",
		"CtrlF": "find",
		"Enter": "find",
	},
	"file": map[string]string{
		"Enter": "find",
	},
}

// 的 都 到
