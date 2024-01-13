package util

func ErrorHandler(err error) {
	if err != nil {
		panic(err)
	}
}
