model User {
  @visibility("read", "update")
  id: int64;

  name: string;

  @visibility("create")
  password: string;

  @visibility("create")
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
  id: string;
  client: PublicClient;
  scopes: string[];
}
