package logar

import "sadk.dev/logar/internal/tableprefix"

func SetTablePrefix(prefix string) {
	tableprefix.Set(prefix)
}

func GetTablePrefix() string {
	return tableprefix.Get()
}
