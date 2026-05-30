import { useAuth } from "react-oidc-context";

export default function LoginPage() {

  const auth = useAuth();

  return (
    <button onClick={() => auth.signinRedirect()}>
      Login with Keycloak
    </button>
  );
}
