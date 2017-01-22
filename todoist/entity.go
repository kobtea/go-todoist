package todoist

type Identifier interface {
	getID() ID
	Equal(id Identifier) bool
}

type Entity struct {
	ID        ID `json:"id,omitempty"`
	IsDeleted IntBool `json:"is_deleted,omitempty"`
}

func (e Entity) getID() ID {
	return e.ID
}

func (e Entity) Equal(entity Identifier) bool {
	return e.ID == entity.getID()
}

type Resolver interface {
	Resolve(id ID) *Entity
}
