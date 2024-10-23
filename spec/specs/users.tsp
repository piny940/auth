import "@typespec/http";
import "@typespec/rest";
import "./model.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace Auth.Users;

model ReqSignup {
  name: string;
  password: string;
  password_confirmation: string;
}

enum SignupErr {
  name_length_not_enough,
  name_already_used,
  password_length_not_enough,
  password_confirmation_not_match,
}

@route("/users")
interface UsersInterface {
  @route("/signup")
  @post
  @summary("Signup")
  signup(@body body: ReqSignup): {
    @statusCode statusCode: 204;
    @header setCookie: string;
  } | {
    @statusCode statusCode: 400;
    error: SignupErr;
    error_description: string;
  };
}
