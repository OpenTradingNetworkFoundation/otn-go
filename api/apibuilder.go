package api

import "log"

// Builder interface provides a clever and convinient way to perform Bitshares API connection initial setup routine
type Builder interface {
	Node(endpoint string) Builder
	Credentials(login, password string) Builder
	BlockHandler(handler BlockAppliedHandler) Builder
	LoginHandler(login func()) Builder
	LogoutHandler(close func()) Builder
	Build() (BitsharesAPIConnection, BitsharesAPI)
}

type builder struct {
	endpoint       string
	conn           BitsharesAPIConnection
	login          string
	password       string
	blockHandlers  []BlockAppliedHandler
	connection     BitsharesAPIConnection
	loginHandlers  []func()
	logoutHandlers []func()
}

// NewBuilder constructs Builder object
func NewBuilder() Builder {
	return &builder{}
}

// Credentials call is used to setup login information
func (b *builder) Credentials(login, password string) Builder {
	b.login = login
	b.password = password
	return b
}

// Node call specifies blockchain node API address
func (b *builder) Node(endpoint string) Builder {
	b.endpoint = endpoint
	return b
}

// BlockHandler call sets up subscription to new block info
func (b *builder) BlockHandler(handler BlockAppliedHandler) Builder {
	b.blockHandlers = append(b.blockHandlers, handler)
	return b
}

// LoginHandler sets up reaction to login API event handler
func (b *builder) LoginHandler(login func()) Builder {
	b.loginHandlers = append(b.loginHandlers, login)
	return b
}

// LogoutHandler sets up reaction to logout API event handler
func (b *builder) LogoutHandler(logout func()) Builder {
	b.logoutHandlers = append(b.logoutHandlers, logout)
	return b
}

// Build performs all necessary settings and results with ready to use API objects
func (b *builder) Build() (BitsharesAPIConnection, BitsharesAPI) {
	conn := NewConnection(b.endpoint)
	api := New(conn, Params{b.login, b.password})
	for _, bh := range b.blockHandlers {
		api.OnLogin(func() {
			dbAPI, err := api.DatabaseAPI()
			if err != nil {
				log.Printf("Unable to get database API: %v", err)
			}
			err = dbAPI.SubscribeBlockApplied(bh)
			if err != nil {
				log.Printf("Unable to subscribe to applied blocks: %v", err)
			}
		})
	}
	for _, cb := range b.loginHandlers {
		api.OnLogin(cb)
	}
	for _, cb := range b.logoutHandlers {
		api.OnLogout(cb)
	}
	return conn, api
}
