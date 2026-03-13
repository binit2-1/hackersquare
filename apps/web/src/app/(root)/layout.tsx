import { AuthProvider } from "@/context/AuthContext";
import { Navbar } from "@/components/navbar";
import { Suspense } from "react";

export default function RootGroupLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <AuthProvider>
      <Suspense fallback={<div className="h-14 w-full" />}>
        <Navbar />
      </Suspense>
      {children}
    </AuthProvider>
  );
}