export default function Footer() { 
  return ( 
    <footer className="h-6 px-5 flex items-center shadow-[0_-4px_12px_rgba(0,0,0,0.25)]
                    justify-center bg-[#595961] text-gray-300 text-xs">
      <span> 
         ©{new Date().getFullYear()}
      </span> 
    </footer> 
  ); }
