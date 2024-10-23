import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

import "./model.tsp";
import "./approvals.tsp";
import "./clients.tsp";
import "./oauth.tsp";
import "./session.tsp";
import "./users.tsp";
import "./account/clients.tsp";

using TypeSpec.Http;
@service({
  title: "Auth",
})
@server("https://auth.piny940.com/api/v1", "Auth Server")
namespace Auth;
