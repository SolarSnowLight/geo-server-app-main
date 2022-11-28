package user

// Модель для работы с экземпляром пользовательских данных из таблицы users
type UserModel struct {
	Id       int    `json:"id" db:"id"`
	Uuid     string `json:"uuid" binding:"required" db:"uuid"`
	Email    string `json:"email" binding:"required" db:"email"`
	Password string `json:"password" binding:"required" db:"password"`
}

// Модель для работы с данными при регистрации пользователя (парсинг JSON, etc.)
type UserRegisterModel struct {
	Id       int            `json:"-" db:"id"`
	Email    string         `json:"email" binding:"required"`
	Password string         `json:"password" binding:"required"`
	Data     UserJSONBModel `json:"data" binding:"required"`
}

/* Модель для хранения основных пользовательских данных */
type UserJSONBModel struct {
	Name       string `json:"name" binding:"required"`
	Surname    string `json:"surname" binding:"required"`
	Patronymic string `json:"patronymic"`
	Gender     bool   `json:"gender"`
	Phone      string `json:"phone"`
	Nickname   string `json:"nickname" binding:"required"`
	DateBirth  string `json:"date_birth"`
}

// Модель для регистрации через Google OAuth2
type UserRegisterOAuth2Model struct {
	Email      string `json:"email" binding:"required"`
	FamilyName string `json:"family_name" binding:"required"`
	GivenName  string `json:"given_name" binding:"required"`
	Name       string `json:"name" binding:"required"`
}

// Модель для работы с данными при авторизации пользователя (парсинг JSON, etc.)
type UserLoginModel struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

/* Модель для работы с данными при авторизации пользователя через Google OAuth2 */
type UserLoginOAuth2Model struct {
	Code string `json:"code" binding:"required"`
}

// Модель представляющая пользовательские авторизационные данные
type UserAuthDataModel struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// Модель представляющая активационные данные пользователя
type UserActivateModel struct {
	ActivationLink string `json:"activation_link" db:"activation_link"`
	IsActivated    bool   `json:"is_activated" db:"is_activated"`
}

/* Модель для представления типов авторизации */
type AuthTypeModel struct {
	Id    int    `json:"id" db:"id"`
	Uuid  string `json:"uuid" db:"uuid"`
	Value string `json:"value" db:"value"`
}

/* Модель для связки пользователей с конкретными типами авторизаций */
type UserAuthTypeModel struct {
	Id          int `json:"id" db:"id"`
	UsersId     int `json:"users_id" db:"users_id"`
	AuthTypesId int `json:"auth_types_id" db:"auth_types_id"`
}
