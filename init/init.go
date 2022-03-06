package init

import "os"

func init() {
	_ = os.MkdirAll("./data/videos", os.ModePerm)
	_ = os.MkdirAll("./data/portrait", os.ModePerm)
	_ = os.MkdirAll("./data/database", os.ModePerm)
}
