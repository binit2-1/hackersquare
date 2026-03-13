import type { Metadata } from "next";
import { Providers } from "@/components/providers";
import "./globals.css";
import Script from "next/script";

export const metadata: Metadata = {
  title: "HackerSquare",
  description: "Discover and apply to hackathons worldwide",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body className="antialiased">
        {/* Umami Analytics */}
        <script
          defer
          src="https://cloud.umami.is/script.js"
          data-website-id="42840191-d0a5-42e1-8adc-43acc2e59b61"
        ></script>
        <Providers>{children}</Providers>
      </body>
    </html>
  );
}
