import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "reconTrace — Infrastructure Monitoring",
  description:
    "Real-time server and infrastructure monitoring for modern engineering teams.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  );
}
