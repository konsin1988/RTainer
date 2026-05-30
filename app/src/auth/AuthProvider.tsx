import { AuthProvider } from "react-oidc-context";
import { oidcConfig } from "./authConfig";

export function AppAuthProvider({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <AuthProvider {...oidcConfig}>
      {children}
    </AuthProvider>
  );
}
