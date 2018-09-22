package db

type Query interface {
	Raw(dst interface{}, sql string, args ...interface{}) (err error)
	Exec(sql string, args ...interface{}) (err error)
	Migrate() error
}
