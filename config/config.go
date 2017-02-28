package config

var DefaultKeys = map[string]map[string]string{
	"global": map[string]string{
		"CtrlC": "quit",
		"Tab":   "nextView",
		"CtrlJ": "nextView",
		"CtrlK": "prevView",
		"CtrlR": "find",
	},
	"file": map[string]string{
		"Enter": "find",
	},
}
