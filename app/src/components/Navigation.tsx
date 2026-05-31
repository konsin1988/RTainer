import { NavLink, useParams } from "react-router-dom";
import { useAuth } from "react-oidc-context";

export default function Navigation() {
  const auth = useAuth();
  const navClass = ({ isActive }) =>
    `
    min-w-[100px]
    h-7
    rounded-lg
    text-[2.2vh]
    flex
    items-center
    justify-center
    text-[#ba9477]
    transition-colors
    ${isActive
        ? "bg-[#ba9477] text-white font-bold hover:bg-[#a18775] active:bg-white"
        : "bg-[#717178] hover:bg-[#8a8a92] active:bg-white"
    }
  `;

  return (
      <nav className="flex gap-4 bg-tranparant shadow">
        {!auth.isAuthenticated ? (
        <>
          <button
            onClick={() => auth.signinRedirect()}
            className="
              min-w-[100px]
              h-7
              rounded-lg
              text-[2.2vh]
              flex
              items-center
              justify-center
              text-[#ba9477]
              transition-colors
              bg-[#ba9477] text-white font-bold hover:bg-[#a18775] active:bg-white
            "
          >
            Log In
          </button>
        </>
        ) : (
        <>
        <NavLink 
            to={`/user/dashboard`}
            className={navClass}
          >
            Main 
          </NavLink> 

          <NavLink 
            to={`/user/costs`} 
            className={navClass}
          >
            Затраты 
          </NavLink>  

          <NavLink 
            to={`/user/dup`} 
            className={navClass}
          >
            ДУП 
          </NavLink> 

          <button
            onClick={() => auth.signoutRedirect()}
            className={`min-w-[100px] h-7 rounded-lg flex items-center text-[2.2vh]
                        justify-center text-[#ba9477] transition-colors 
                        bg-[#717178] hover:bg-[#8a8a92]`}
          >
            Logout
          </button>
          </>
        )}
      </nav>
  );
}
