"use client";

import React, { createContext, useCallback, useContext, useEffect, useState } from "react";
import { AuthContextType, User } from "@/models/user";
import { useRouter } from "next/navigation";

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const router = useRouter();

  const checkAuth = useCallback(async () => {
    setIsLoading(true);
    try {
      const res = await fetch(`/api/v1/auth/me`, {
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
  }, []);

  useEffect(() => {
    checkAuth();
  }, [checkAuth]);

  const logout = async () => {
    try {
      await fetch(`/api/v1/auth/logout`, {
        method: "POST",
        credentials: "include",
      });

    } catch (error) {
        console.error("Logout failed:", error);     
    } finally{
        setUser(null);
        router.push("/");
        router.refresh();
    }
    
  };

  const refreshUser = async () => {
    await checkAuth();
  };
  return (
    <AuthContext.Provider value={{ user, isLoading, logout, refreshUser }}>
      {children}
    </AuthContext.Provider>
  );
}


export const useAuth = () =>{
    const context = useContext(AuthContext)
    if (context === undefined) throw new Error("useAuth must be used within an AuthProvider");
    return context;
}
