import "@typespec/http";
import "@typespec/rest";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace Auth.Approvals;

model ApproveReq {
  client_id: string;
  scope: string;
}
model Res4xx {
  error: string;
  error_description: string;
}

@route("/account/approvals")
interface ApprovalsInterface {
  @summary("Approve a auth request")
  @post
  approve(@body body: ApproveReq): {
    @statusCode statusCode: 204;
  } | {
    @statusCode statusCode: 400;
    @body body: Res4xx;
  } | {
    @statusCode statusCode: 401;
    @body body: Res4xx;
  };
}
