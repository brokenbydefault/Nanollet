// +build !js

package Util

import "runtime/debug"

func FreeMemory(){
	debug.FreeOSMemory()
}
