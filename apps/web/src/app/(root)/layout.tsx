import { AuthProvider } from "@/context/AuthContext";
import { Navbar } from "@/components/navbar";

export default function RootGroupLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <AuthProvider>
      <Navbar />
      {children}
    </AuthProvider>
  );
}