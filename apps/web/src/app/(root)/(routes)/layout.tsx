import * as React from "react";
import { Navbar } from "@/components/navbar";

export default function RoutesLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <>
      <Navbar />
      <div className="mx-auto min-h-screen max-w-196.5">{children}</div>
    </>
  );
}
