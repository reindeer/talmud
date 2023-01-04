package try

func Catch(try func(), catch func(throwable error)) {
	defer func() {
		if exception := recover(); exception != nil {
			if throwable, ok := exception.(error); ok {
				if catch != nil {
					catch(throwable)
				}
				return
			}
			panic(exception)
		}
	}()

	try()
}

func Throw[T any](result T, err error) T {
	ThrowError(err)
	return result
}

func ThrowError(err error) {
	if err != nil {
		panic(err)
	}
}
