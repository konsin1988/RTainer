import Navigation from "./Navigation"

export default function Header() { 
  return ( 
    <header className="sticky top-0 z-50 h-[60px] px-8 flex items-center 
                    justify-between shadow-lg border-b border-slate-200 bg-[#595961]"> 
      <div className="flex items-center gap-3"> 
        <img
          src="/logo.svg"
          alt="Rtainer"
          className="h-[6.5vh]"
        />
      </div> <Navigation /> 
    </header> 
  ); }
