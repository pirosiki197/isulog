package isulog

type Config struct {
	Filename string
}

var DefaultConfig = Config{
	Filename: "/home/isucon/isulog.out",
}
