import "@typespec/http";
import "@typespec/rest";
import "./model.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace Auth.Clients;

@route("/account/clients")
interface ClientInterface {
  @route("/")
  @get
  @summary("Get all clients")
  listClients(@header cookie: string): {
    @statusCode statusCode: 200;
    clients: Client[];
  };

  @route("/")
  @post
  @summary("Create a new client")
  createClient(client: Client): {
    @statusCode statusCode: 201;
    client: Client;
  } | {
    @statusCode statusCode: 400;
    error: string;
  };

  @route("/:id")
  @post
  @summary("Update a client")
  updateClient(@path id: int64, client: Client): {
    @statusCode statusCode: 200;
    client: Client;
  } | {
    @statusCode statusCode: 400;
    error: string;
  };

  @route("/:id")
  @delete
  @summary("Delete a client")
  deleteClient(@path id: int64): {
    @statusCode statusCode: 204;
  } | {
    @statusCode statusCode: 400;
    error: string;
  };
}
