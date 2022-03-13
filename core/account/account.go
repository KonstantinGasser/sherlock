package account

type Serializer interface {
	Serialize() ([]byte, error)
}
