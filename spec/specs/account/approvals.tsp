import "@typespec/http";
import "@typespec/rest";
import "../model.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace Auth.Account.Approvals;

model ListApprovalsRes {
  approvals: Approval[];
}

@route("/account/approvals")
@useAuth(ApiKeyAuth<ApiKeyLocation.cookie, "auth">)
interface Approvals {
  @route("")
  @get
  @summary("Get all approvals")
  listApprovals(): {
    @statusCode statusCode: 200;
    @body body: ListApprovalsRes;
  };

  @route("/{id}")
  @delete
  @summary("Delete an approval")
  deleteApproval(@path id: int64): {
    @statusCode statusCode: 204;
  } | {
    @statusCode statusCode: 404;
    error: string;
  };
}
