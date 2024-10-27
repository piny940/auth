import "@typespec/http";
import "@typespec/rest";
import "./model.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace Auth.OAuth;

model AuthorizeReq {
  response_type: string;
  client_id: string;
  redirect_uri: string;
  scope: string;
  state?: string;
}
enum AuthorizeErr {
  invalid_request,
  unauthorized_client,
  access_denied,
  unsupported_response_type,
  invalid_scope,
  server_error,
  temporarily_unavailable,
}
model TokenReq {
  grant_type: string;
  code: string;
  redirect_uri: string;
  client_id: string;
}
model TokenRes {
  access_token: string;
  token_type: TokenTokenType;
  expires_in: int32;

  @summary("Id token, if scope includes 'openid'")
  id_token?: string;
}
enum TokenTokenType {
  Bearer,
}
enum TokenCacheControlHeader {
  no_store: "no-store",
}
enum TokenPragmaHeader {
  no_store: "no-store",
}
enum TokenErr {
  invalid_request,
}

@route("/oauth")
interface OAuthInterface {
  @route("/authorize")
  @get
  @summary("Authorization Request")
  authorize(
    @query response_type: string,
    @query client_id: string,
    @query redirect_uri: string,
    @query scope: string,
    @query state?: string,
  ): {
    @statusCode statusCode: 204;
  } | {
    @statusCode statusCode: 302;
    @header location: string;
  } | {
    @statusCode statusCode: 400;
    error: AuthorizeErr;
    error_description: string;
    state?: string;
  };

  @route("/authorize")
  @post
  @summary("Authorization Request")
  postAuthorize(
    @header contentType: "multipart/form-data",
    @body body: AuthorizeReq,
  ): {
    @statusCode statusCode: 204;
  } | {
    @statusCode statusCode: 302;
    @header location: string;
  } | {
    @statusCode statusCode: 400;
    error: AuthorizeErr;
    error_description: string;
    state?: string;
  };

  @route("/token")
  @post
  @summary("Get token")
  @useAuth(BasicAuth)
  getToken(
    @header contentType: "application/x-www-form-urlencoded",
    @header authorization: string,
    @body body: TokenReq,
  ): {
    @statusCode statusCode: 200;
    @header cacheControl: TokenCacheControlHeader;
    @header pragma: TokenPragmaHeader;
    @body body: TokenRes;
  } | {
    @statusCode statusCode: 400;
    error: TokenErr;
    error_description: string;
  };
}
