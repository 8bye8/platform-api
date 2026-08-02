package types

type RegisterUserParams struct {
	Email    string
	Password string
	Username string
}

type LoginParams struct {
	Username string
	Password string
}

type HasherParams struct {
	SaltLength int
	KeyLength  uint32
	TimeCost   uint32
	MemoryCost uint32
	Threads    uint8
}
