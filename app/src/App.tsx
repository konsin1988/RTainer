import { useState } from 'react'
import AppRouter from "./router/router"
import { AppAuthProvider } from "./auth/AuthProvider"; 
import './App.css'

function App() {
  return (
    <AppAuthProvider>
      <AppRouter />
    </AppAuthProvider>
  ); 
}

export default App;
