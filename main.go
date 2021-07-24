// eco4go
package main

import (
	"eco4go/dbmod"
	"eco4go/router"
)

func main() {
	dbmod.OpenDB()
	defer dbmod.CloseDB()
	router.OpenEcho()
}
