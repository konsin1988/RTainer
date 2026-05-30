import { useAuth } from "react-oidc-context";

export default function Dashboard() {

  const auth = useAuth();
  console.log(auth.user);
  
  return (
    <div>
      <h1>Dashboard</h1>

      <p>
        Hello {
          auth.user?.profile.preferred_username
        }
      </p>

      <button
        onClick={() => auth.signoutRedirect()}
      >
        Logout
      </button>
    </div>
  );
}
