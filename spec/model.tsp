model User {
  @visibility("read", "update")
  id: int64;

  name: string;
}
model Client {
  @visibility("read", "update")
  id: int64;

  name: string;
  redirect_urls: string[];
}
