package gateway

import "auth/internal/domain/oauth"

var ScopeMap = map[int32]oauth.TypeScope{
	0: oauth.ScopeOpenID,
	1: oauth.ScopeEmail,
	2: oauth.ScopeProfile,
}
var ScopeMapReverse = map[oauth.TypeScope]int32{
	oauth.ScopeOpenID:  0,
	oauth.ScopeEmail:   1,
	oauth.ScopeProfile: 2,
}
