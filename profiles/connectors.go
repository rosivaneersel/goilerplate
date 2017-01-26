package profiles

type UserTestDB struct {
	DB *TestDB
}

func (o *UserTestDB) Get(id uint) (*User, error) {
	user := &User{}
	err := o.DB.First(user, id)

	return user, err
}

func NewPageManagerTestDB(db *TestDB) *UserTestDB {
	return &UserTestDB{DB: db}
}

