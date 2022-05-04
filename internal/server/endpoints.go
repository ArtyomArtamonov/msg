package server

type Endpoints struct {
	ApiService     apiServiceEndpoints
	AuthService    authServiceEndpoints
	MessageService messageServiceEndpoints
}

type apiServiceEndpoints struct {
	CreateRoom string
	ListRooms  string
}

type authServiceEndpoints struct {
	Login    string
	Register string
	Refresh  string
}

type messageServiceEndpoints struct {
	SendMessage string
	GetMessages string
}

func NewEndpoints() *Endpoints {
	apiServicePath := "/api.ApiService/"
	authServicePath := "/auth.AuthService/"
	messageServicePath := "/message.MessageService/"
	return &Endpoints{
		ApiService: apiServiceEndpoints{
			CreateRoom: apiServicePath + "CreateRoom",
			ListRooms:  apiServicePath + "ListRooms",
		},
		AuthService: authServiceEndpoints{
			Login:    authServicePath + "Login",
			Register: authServicePath + "Register",
			Refresh:  authServicePath + "Refresh",
		},
		MessageService: messageServiceEndpoints{
			SendMessage: messageServicePath + "SendMessage",
			GetMessages: messageServicePath + "GetMessages",
		},
	}
}
