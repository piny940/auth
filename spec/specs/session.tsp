import "@typespec/http";
import "@typespec/rest";
import "./model.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace Auth.Session;

model MeRes {
  user: User | null;
}
model LoginReq {
  name: string;
  password: string;
}
enum LoginErr {
  invalid_name_or_password,
}
enum LogoutErr {
  not_logged_in,
}
@route("/session")
interface SessionInterface {
  @route("")
  @get
  @summary("Get session")
  @useAuth(ApiKeyAuth<ApiKeyLocation.cookie, "auth">)
  me(): {
    @statusCode statusCode: 200;
    @body body: MeRes;
  };

  @route("")
  @post
  @summary("Login")
  login(@body body: LoginReq): {
    @statusCode statusCode: 204;
  } | {
    @statusCode statusCode: 400;
    error: LoginErr;
    error_description: string;
  };

  @route("")
  @delete
  @summary("Logout")
  @useAuth(ApiKeyAuth<ApiKeyLocation.cookie, "auth">)
  logout(): {
    @statusCode statusCode: 204;
  };
}
