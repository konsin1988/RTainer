export const oidcConfig = {
  authority: "http://localhost:8080/realms/rtainer",
  client_id: "rtainer-frontend",
  redirect_uri: "http://localhost:5173/auth/callback",
  post_logout_redirect_uri: "http://localhost:5173/login",
  response_type: "code",
  scope: "openid profile email",
};
