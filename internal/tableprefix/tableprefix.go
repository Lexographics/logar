package tableprefix

var tablePrefix string = ""

func Set(prefix string) {
	tablePrefix = prefix
}

func Get() string {
	return tablePrefix
}
