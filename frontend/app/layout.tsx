import type { Metadata } from "next";
import "./globals.css";
import { Anton } from "next/font/google";

export const metadata: Metadata = {
  title: "reconTrace — Infrastructure Monitoring",
  description:
    "Real-time server and infrastructure monitoring for modern engineering teams.",
};
const anton = Anton({
  weight: "400",
  subsets: ["latin"],
  variable: "--font-display",
});

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <head>
        <link rel="preconnect" href="https://fonts.googleapis.com" />
        <link rel="preconnect" href="https://fonts.gstatic.com" crossOrigin="anonymous" />
        <link
          href="https://fonts.googleapis.com/css2?family=Big+Shoulders+Display:wght@600;700;800&family=IBM+Plex+Sans:wght@400;500;600&family=IBM+Plex+Mono:wght@400;500&display=swap"
          rel="stylesheet"
        />
      </head>
      <body className={anton.variable}>{children}</body>
    </html>
  );
}
