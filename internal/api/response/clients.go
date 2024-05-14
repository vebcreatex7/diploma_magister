package response

type User struct {
	UID        string
	Surname    string
	Name       string
	Patronymic string
	Login      string
	Email      string
	Role       string
	Approved   bool
}
