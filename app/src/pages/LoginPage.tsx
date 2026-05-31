import { useAuth } from "react-oidc-context";

export default function LoginPage() {

  const auth = useAuth();

  return (
    <button 
      className="z-10000"
      onClick={() => auth.signinRedirect()}>
      Login with Keycloak
    </button>
  );
}
