module github.com/fogfish/logger/examples

go 1.23.1

replace github.com/fogfish/logger/v3 => ../

replace github.com/fogfish/logger/x/xlog => ../x/xlog

require (
	github.com/fogfish/logger/v3 v3.2.0
	github.com/fogfish/logger/x/xlog v0.0.0-00010101000000-000000000000
)
