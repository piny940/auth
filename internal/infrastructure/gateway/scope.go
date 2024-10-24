package gateway

import "auth/internal/domain/oauth"

var scopeMap = map[int32]oauth.TypeScope{
	0: oauth.ScopeOpenID,
}
var scopeMapReverse = map[oauth.TypeScope]int32{
	oauth.ScopeOpenID: 0,
}
