package sandbox_plugin

import "os"

func Execute() {
	os.Remove("test.txt")
}
