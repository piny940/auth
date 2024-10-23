import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

import "./specs/model.tsp";
import "./specs/approvals.tsp";
import "./specs/clients.tsp";
import "./specs/oauth.tsp";
import "./specs/session.tsp";
import "./specs/users.tsp";
import "./specs/account/clients.tsp";
import "./specs/healthz.tsp";

using TypeSpec.Http;
@service({
  title: "Auth",
})
@server("https://auth.piny940.com/api/v1", "Auth Server")
namespace Auth;
