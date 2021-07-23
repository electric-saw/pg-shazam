package definitions

type Error struct {
	Err error
}

func (e Error) Frontend()                {}
func (e Error) Backend()                 {}
func (e Error) Decode(data []byte) error { return nil }
func (e Error) Encode(dst []byte) []byte { return dst }
