package appmodels

type Image struct {
	Max string `json:"max"`
}

type Product struct {
	RaecId      string `json:"productId" bson:"_id,omitempty"`
	Description string `json:"descriptionShort,omitempty" bson:"descriptionShort,omitempty" validate:"required,min=3,max=500"`
	Image       Image  `json:"image"`
}

type Remains struct {
	Name  string  `db:"name"`
	Code  string  `db:"code"`
	Cell  string  `db:"cell"`
	EH    string  `db:"eh"`
	Count float64 `db:"count"`
}

type Order struct {
	Executor string
	Zone     string
}
