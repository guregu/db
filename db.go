package db

type indexkey int

const (
	sqlIndex indexkey = iota + 1
	redisIndex
	mongoIndex
)
