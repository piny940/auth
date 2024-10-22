import "@typespec/http";
import "@typespec/rest";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace Auth.Clients;

enum GetClientErr {
  client_not_found,
}

model GetClientRes {
  client: PublicClient;
}

@route("/clients")
interface ClientsInterface {
  @route("/{id}")
  @get
  @summary("Get a client")
  getClient(@path id: string): {
    @statusCode statusCode: 200;
    @body body: GetClientRes;
  } | {
    @statusCode statusCode: 400;
    error: GetClientErr;
    error_description: string;
  };
}
