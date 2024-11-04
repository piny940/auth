package gateway

import "auth/internal/domain/oauth"

var scopeMap = map[int32]oauth.TypeScope{
	0: oauth.ScopeOpenID,
	1: oauth.ScopeEmail,
}
var scopeMapReverse = map[oauth.TypeScope]int32{
	oauth.ScopeOpenID: 0,
	oauth.ScopeEmail:  1,
}
