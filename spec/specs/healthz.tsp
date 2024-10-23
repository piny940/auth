import "@typespec/http";
import "@typespec/rest";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace Auth.Healthz;

@route("/healthz")
interface Healthz {
  @get
  @summary("Check health")
  check(): {
    @statusCode statusCode: 200;
  };
}
