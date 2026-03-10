"use client";

import React, { createContext, useContext, useEffect, useState } from "react";
import { AuthContextType, User } from "@/models/user";

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    checkAuth();
  }, []);

  const checkAuth = async () => {
    try {
      const res = await fetch("http://localhost:8080/v1/auth/me", {
        method: "GET",
        credentials: "include",
      });

      if (res.ok) {
        const data = await res.json();
        setUser(data);
      } else {
        setUser(null);
      }
    } catch (err) {
      setUser(null);
    } finally {
      setIsLoading(false);
    }
  };

  const logout = async () => {
    //TODO: Implement logout API call
    setUser(null);
  };
  return (
    <AuthContext.Provider value={{ user, isLoading, logout }}>
      {children}
    </AuthContext.Provider>
  );
}


export const useAuth = () =>{
    const context = useContext(AuthContext)
    if (context === undefined) throw new Error("useAuth must be used within an AuthProvider");
    return context;
}