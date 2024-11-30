import "@typespec/http";
import "@typespec/rest";
import "./model.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace Auth.Userinfo;

model UserinfoRes {
  sub: string;
  name?: string;
  email?: string;
}

@route("/resources/userinfo")
@useAuth(BearerAuth)
interface Userinfo {
  @route("")
  @get
  @summary("Get userinfo")
  getUserinfo(): {
    @statusCode statusCode: 200;
    @body body: UserinfoRes;
  };
}
