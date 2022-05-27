package server

type Endpoints struct {
	ApiService     apiServiceEndpoints
	AuthService    authServiceEndpoints
	MessageService messageServiceEndpoints
}

type apiServiceEndpoints struct {
	CreateRoom   string
	ListRooms    string
	SendMessage  string
	ListMessages string
}

type authServiceEndpoints struct {
	Login    string
	Register string
	Refresh  string
}

type messageServiceEndpoints struct {
	GetMessages string
}

func NewEndpoints() *Endpoints {
	apiServicePath := "/api.ApiService/"
	authServicePath := "/auth.AuthService/"
	messageServicePath := "/message.MessageService/"
	return &Endpoints{
		ApiService: apiServiceEndpoints{
			CreateRoom:   apiServicePath + "CreateRoom",
			ListRooms:    apiServicePath + "ListRooms",
			SendMessage:  messageServicePath + "SendMessage",
			ListMessages: messageServicePath + "ListMessages",
		},
		AuthService: authServiceEndpoints{
			Login:    authServicePath + "Login",
			Register: authServicePath + "Register",
			Refresh:  authServicePath + "Refresh",
		},
		MessageService: messageServiceEndpoints{
			GetMessages: messageServicePath + "GetMessages",
		},
	}
}
