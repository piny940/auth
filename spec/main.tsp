import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

import "./model.tsp";
import "./approvals.tsp";
import "./clients.tsp";

using TypeSpec.Http;
@service({
  title: "Auth",
})
@server("https://auth.piny940.com/api/v1", "Auth Server")
namespace Auth;

model ReqLogin {
  name: string;
  password: string;
}
@route("/authorize")
@get
@summary("Authorization Request")
op authorize(
  @query response_type: string,
  @query client_id: string,
  @query redirect_uri: string,
  @query scope: string,
  @query state?: string,
): {
  @statusCode statusCode: 204;
} | {
  @statusCode statusCode: 302;
  error: string;
  error_description: string;
  state?: string;
} | {
  @statusCode statusCode: 400;
  error: string;
  error_description: string;
  state?: string;
};

@route("/authorize")
@post
@summary("Authorization Request")
op postAuthorize(
  @header contentType: "multipart/form-data",
  response_type: string,
  client_id: string,
  redirect_uri: string,
  scope: string,
  state?: string,
):
  | {
      @statusCode statusCode: 204;
    }
  | {
      @statusCode statusCode: 302;
      error: string;
      error_description: string;
      state?: string;
    }
  | {
      @statusCode statusCode: 400;
      error: string;
      error_description: string;
      state?: string;
    }
  | {
      @statusCode statusCode: 401;
      error: string;
      error_description: string;
      state?: string;
    };

@route("/token")
namespace Token {
  @route("/")
  @post
  @summary("Get a token")
  op getToken(client_id: string, client_secret: string): {
    @statusCode statusCode: 200;
    token: string;
  } | {
    @statusCode statusCode: 400;
    error: int32;
    error_description: string;
  };
}

@route("/session")
@post
@summary("Login")
op login(@body body: ReqLogin): {
  @statusCode statusCode: 204;
} | {
  @statusCode statusCode: 400;
  error: "invalid_name_or_password";
  error_description: string;
};
@route("/session")
@delete
@summary("Logout")
op logout(): {
  @statusCode statusCode: 204;
} | {
  @statusCode statusCode: 400;
  error: "not_logged_in";
  error_description: string;
};

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
  error:
    | "name_already_used"
    | "name_length_not_enough"
    | "password_length_not_enough"
    | "password_confirmation";
  error_description: string;
};

@route("/me")
@get
@summary("Get me")
op me(): {
  @statusCode statusCode: 200;
  user: User;
};
