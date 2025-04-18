import "@typespec/http";
import "@typespec/rest";
import "../model.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace Auth.Account.Clients;

model CreateClientReq {
  name: string;
  redirect_urls: string[];
}
model CreatedClient {
  id: string;
  name: string;
  secret: string;
  redirect_urls: string[];
}
model UpdateClientReq {
  name: string;
  redirect_urls: string[];
}

@route("/account/clients")
@useAuth(ApiKeyAuth<ApiKeyLocation.cookie, "auth">)
interface AccountClients {
  @route("")
  @get
  @summary("Get all clients")
  listClients(): {
    @statusCode statusCode: 200;
    clients: Client[];
  };

  @route("/{id}")
  @get
  @summary("Get a client")
  getClient(@path id: string): {
    @statusCode statusCode: 200;
    client: Client;
  } | {
    @statusCode statusCode: 404;
    error: string;
  };

  @route("")
  @post
  @summary("Create a new client")
  createClient(client: CreateClientReq): {
    @statusCode statusCode: 201;
    client: CreatedClient;
  } | {
    @statusCode statusCode: 400;
    error: string;
  };

  @route("/{id}")
  @post
  @summary("Update a client")
  updateClient(@path id: string, client: UpdateClientReq): {
    @statusCode statusCode: 204;
  } | {
    @statusCode statusCode: 400;
    error: string;
  };

  @route("/{id}")
  @delete
  @summary("Delete a client")
  deleteClient(@path id: string): {
    @statusCode statusCode: 204;
  } | {
    @statusCode statusCode: 400;
    error: string;
  };
}
