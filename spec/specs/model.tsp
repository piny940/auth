model User {
  @visibility(Lifecycle.Read, Lifecycle.Update)
  id: int64;

  name: string;

  @visibility(Lifecycle.Create)
  password: string;

  @visibility(Lifecycle.Create)
  password_confirmation: string;
}
model Client {
  id: string;
  name: string;
  redirect_urls: string[];
}
model PublicClient {
  id: string;
  name: string;
}
model Approval {
  id: int64;
  client: PublicClient;
  scopes: string[];
}
