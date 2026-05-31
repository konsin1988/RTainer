import React from "react";
import { Outlet } from "react-router-dom";
import Header from "../components/Header";
import Footer from "../components/Footer";


export default function Layout() { 
  return ( 
    <div className="min-h-screen flex flex-col"> 
      <Header /> 
      <main className="
        flex-1 p-8 
        bg-slate-50 
        bg-[radial-gradient(circle_at_1px_1px,rgba(100,116,139,0.15)_1px,transparent_0)] 
        bg-[length:24px_24px]
      "> 
        <Outlet />
      </main> 
      <Footer /> 
    </div> 
  ); }

