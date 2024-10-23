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
  @summary("Get a token")
  getToken(client_id: string, client_secret: string): {
    @statusCode statusCode: 200;
    token: string;
  } | {
    @statusCode statusCode: 400;
    error: int32;
    error_description: string;
  };
}
