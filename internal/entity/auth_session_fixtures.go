package entity

import (
	"github.com/photoprism/photoprism/pkg/authn"
	"github.com/photoprism/photoprism/pkg/clean"
)

type SessionMap map[string]Session

func (m SessionMap) Get(name string) Session {
	if result, ok := m[name]; ok {
		return result
	}

	return Session{}
}

func (m SessionMap) Pointer(name string) *Session {
	if result, ok := m[name]; ok {
		return &result
	}

	return &Session{}
}

var SessionFixtures = SessionMap{
	"alice": {
		ID:          "69be27ac5ca305b394046a83f6fda18167ca3d3f2dbe7ac0",
		RefID:       "sessxkkcabcd",
		SessTimeout: UnixDay * 3,
		SessExpires: UnixTime() + UnixWeek,
		user:        UserFixtures.Pointer("alice"),
		UserUID:     UserFixtures.Pointer("alice").UserUID,
		UserName:    UserFixtures.Pointer("alice").UserName,
	},
	"alice_token": {
		ID:           "bb8658e779403ae524a188712470060f050054324a8b104e",
		RefID:        "sess34q3hael",
		SessTimeout:  -1,
		SessExpires:  UnixTime() + UnixDay,
		AuthScope:    clean.Scope("*"),
		AuthProvider: authn.ProviderClient.String(),
		AuthMethod:   authn.MethodAccessToken.String(),
		LastActive:   -1,
		user:         UserFixtures.Pointer("alice"),
		UserUID:      UserFixtures.Pointer("alice").UserUID,
		UserName:     UserFixtures.Pointer("alice").UserName,
	},
	"alice_token_scope": {
		ID:            "778f0f7d80579a072836c65b786145d6e0127505194cc51e",
		RefID:         "sessjr0ge18d",
		SessTimeout:   0,
		SessExpires:   UnixTime() + UnixDay,
		AuthScope:     clean.Scope("metrics photos albums videos"),
		AuthProvider:  authn.ProviderClient.String(),
		AuthMethod:    authn.MethodAccessToken.String(),
		user:          UserFixtures.Pointer("alice"),
		UserUID:       UserFixtures.Pointer("alice").UserUID,
		UserName:      UserFixtures.Pointer("alice").UserName,
		PreviewToken:  "cdd3r0lr",
		DownloadToken: "64ydcbom",
	},
	"bob": {
		ID:          "69be27ac5ca305b394046a83f6fda18167ca3d3f2dbe7ac1",
		RefID:       "sessxkkcabce",
		SessTimeout: UnixDay * 3,
		SessExpires: UnixTime() + UnixWeek,
		user:        UserFixtures.Pointer("bob"),
		UserUID:     UserFixtures.Pointer("bob").UserUID,
		UserName:    UserFixtures.Pointer("bob").UserName,
	},
	"unauthorized": {
		ID:          "69be27ac5ca305b394046a83f6fda18167ca3d3f2dbe7ac2",
		RefID:       "sessxkkcabcf",
		SessTimeout: UnixDay * 3,
		SessExpires: UnixTime() + UnixWeek,
		user:        UserFixtures.Pointer("unauthorized"),
		UserUID:     UserFixtures.Pointer("unauthorized").UserUID,
		UserName:    UserFixtures.Pointer("unauthorized").UserName,
	},
	"visitor": {
		ID:          "69be27ac5ca305b394046a83f6fda18167ca3d3f2dbe7ac3",
		RefID:       "sessxkkcabcg",
		SessTimeout: UnixDay * 3,
		SessExpires: UnixTime() + UnixWeek,
		user:        &Visitor,
		UserUID:     Visitor.UserUID,
		UserName:    Visitor.UserName,
		DataJSON:    []byte(`{"tokens":["1jxf3jfn2k"],"shares":["as6sg6bxpogaaba8"]}`),
		data: &SessionData{
			Tokens: []string{"1jxf3jfn2k"},
			Shares: UIDs{"as6sg6bxpogaaba8"},
		},
	},
	"visitor_token_metrics": {
		ID:           "4ebe1048a7384e1e6af2930b5b6f29795ffab691df47a488",
		RefID:        "sessaae5cxun",
		SessTimeout:  0,
		SessExpires:  UnixTime() + UnixWeek,
		AuthScope:    clean.Scope("metrics"),
		AuthProvider: authn.ProviderClient.String(),
		AuthMethod:   authn.MethodAccessToken.String(),
		user:         &Visitor,
		UserUID:      Visitor.UserUID,
		UserName:     Visitor.UserName,
	},
	"friend": {
		ID:          "69be27ac5ca305b394046a83f6fda18167ca3d3f2dbe7ac4",
		RefID:       "sessxkkcabch",
		SessTimeout: UnixDay * 3,
		SessExpires: UnixTime() + UnixWeek,
		user:        UserFixtures.Pointer("friend"),
		UserUID:     UserFixtures.Pointer("friend").UserUID,
		UserName:    UserFixtures.Pointer("friend").UserName,
	},
	"token_metrics": {
		ID:            "9d8b8801ffa23eb52e08ca7766283799ddfd8dd368208a9b",
		RefID:         "sessgh6gjuo1",
		SessTimeout:   0,
		SessExpires:   UnixTime() + UnixWeek,
		AuthScope:     clean.Scope("metrics"),
		AuthProvider:  authn.ProviderClient.String(),
		AuthMethod:    authn.MethodAccessToken.String(),
		user:          nil,
		UserUID:       "",
		UserName:      "",
		PreviewToken:  "py2xrgr3",
		DownloadToken: "vgln2ffb",
	},
	"token_settings": {
		ID:            "3f9684f7d3dd3d5b84edd43289c7fb5ca32ee73bd0233237",
		RefID:         "sessyugn54so",
		SessTimeout:   0,
		SessExpires:   UnixTime() + UnixWeek,
		AuthScope:     clean.Scope("settings"),
		AuthProvider:  authn.ProviderClient.String(),
		AuthMethod:    authn.MethodAccessToken.String(),
		user:          nil,
		UserUID:       "",
		UserName:      "",
		PreviewToken:  "py2xrgr3",
		DownloadToken: "vgln2ffb",
	},
}

// CreateSessionFixtures inserts known entities into the database for testing.
func CreateSessionFixtures() {
	for _, entity := range SessionFixtures {
		Db().Create(&entity)
	}
}
