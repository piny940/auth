import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

using TypeSpec.Http;
@service({
  title: "Auth",
})
namespace Auth;

model Client {
  @visibility("read", "update")
  id: int64;

  name: string;
  redirect_urls: string[];
}

@route("/api")
namespace Api {
  @route("/v1")
  namespace V1 {
    @route("/token")
    namespace Token {
      @route("/")
      @post
      @summary("Get a token")
      op getToken(client_id: string, client_secret: string): {
        @statusCode statusCode: 200;
        @body token: string;
      } | {
        @statusCode statusCode: 400;
        @body error: string;
      };
    }

    @route("/account")
    namespace Account {
      @route("/clients")
      namespace Clients {
        @route("/")
        @get
        @summary("Get all clients")
        op listClients(): {
          @statusCode statusCode: 200;
          @body clients: Client[];
        };

        @route("/")
        @post
        @summary("Create a new client")
        op createClient(client: Client): {
          @statusCode statusCode: 201;
          @body client: Client;
        } | {
          @statusCode statusCode: 400;
          @body error: string;
        };

        @route("/:id")
        @post
        @summary("Update a client")
        op updateClient(@path id: int64, client: Client): {
          @statusCode statusCode: 200;
          @body client: Client;
        } | {
          @statusCode statusCode: 400;
          @body error: string;
        };

        @route("/:id")
        @delete
        @summary("Delete a client")
        op deleteClient(@path id: int64): {
          @statusCode statusCode: 204;
        } | {
          @statusCode statusCode: 400;
          @body error: string;
        };
      }
    }
  }
}
