module mosn.io/layotto/sdk/go-layotto-init

go 1.14

require (
	mosn.io/layotto v0.2.1-0.20211015040910-1cce41398cee
	mosn.io/layotto/components v0.0.0-20211020084508-6f5ee3cfeba0
	mosn.io/layotto/sdk/go-sdk v0.0.0-20211020084508-6f5ee3cfeba0
	mosn.io/layotto/spec v0.0.0-20211020084508-6f5ee3cfeba0
)

replace (
	mosn.io/layotto => ../../
	mosn.io/layotto/components => ../../components
	mosn.io/layotto/sdk/go-sdk => ../go-sdk
	mosn.io/layotto/spec => ../../spec
)
