import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

import "./model.tsp";

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

@route("/login")
@post
@summary("Login")
op login(@body body: ReqLogin): {
  @statusCode statusCode: 204;
} | {
  @statusCode statusCode: 400;
  error: string;
  error_description: string;
};

@route("/me")
@get
@summary("Get me")
op me(): {
  @statusCode statusCode: 200;
  user: User;
};

@route("/account")
namespace Account {
  @route("/clients")
  namespace Clients {
    @route("/")
    @get
    @summary("Get all clients")
    op listClients(@header cookie: string): {
      @statusCode statusCode: 200;
      clients: Client[];
    };

    @route("/")
    @post
    @summary("Create a new client")
    op createClient(client: Client): {
      @statusCode statusCode: 201;
      client: Client;
    } | {
      @statusCode statusCode: 400;
      error: string;
    };

    @route("/:id")
    @post
    @summary("Update a client")
    op updateClient(@path id: int64, client: Client): {
      @statusCode statusCode: 200;
      client: Client;
    } | {
      @statusCode statusCode: 400;
      error: string;
    };

    @route("/:id")
    @delete
    @summary("Delete a client")
    op deleteClient(@path id: int64): {
      @statusCode statusCode: 204;
    } | {
      @statusCode statusCode: 400;
      error: string;
    };
  }
}
