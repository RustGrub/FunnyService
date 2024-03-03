package projects

type Repository interface {
	Create()
	Get()
	Update()
	Delete()
	GetByName()
	// etc...
}
