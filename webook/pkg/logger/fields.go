package logger

func Error(err error) Field {
	return Field{
		Key:   "error",
		Value: err,
	}
}

func Int64(value int64) Field {
	return Field{
		Key:   "int64",
		Value: value,
	}
}
