package core

type Writer interface {
	Write(key string, s fs.Serializer) error
}

type Reader interface {
	// Read reads space-data into the passed in space.Space.
	Read(key string) ([]byte, error)
}

type SherlockFS interface {
	Writer
	Reader
}

type Initializer interface {
	Initialize(key string, sapce []byte) error
}
