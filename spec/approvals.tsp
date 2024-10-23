import "@typespec/http";
import "@typespec/rest";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace Auth.Approvals;

model ApproveReq {
  client_id: string;
  scope: string;
}
enum ApproveErr {
  invalid_client,
  invalid_scope,
}

@route("/account/approvals")
@useAuth(ApiKeyAuth<ApiKeyLocation.cookie, "auth">)
interface ApprovalsInterface {
  @summary("Approve a auth request")
  @post
  approve(@body body: ApproveReq): {
    @statusCode statusCode: 204;
  } | {
    @statusCode statusCode: 400;
    error: ApproveErr;
    error_description: string;
  };
}
