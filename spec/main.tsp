import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

import "./model.tsp";
import "./approvals.tsp";
import "./clients.tsp";
import "./oauth.tsp";
import "./session.tsp";

using TypeSpec.Http;
@service({
  title: "Auth",
})
@server("https://auth.piny940.com/api/v1", "Auth Server")
namespace Auth;

model ReqSignup {
  name: string;
  password: string;
  password_confirmation: string;
}
@route("/signup")
@post
@summary("Signup")
op signup(@body body: User): {
  @statusCode statusCode: 204;
} | {
  @statusCode statusCode: 400;
  error: string;
  error_description: string;
};
