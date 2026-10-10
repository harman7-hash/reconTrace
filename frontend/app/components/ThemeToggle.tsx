"use client";

import { useEffect, useState } from "react";

export default function ThemeToggle() {
  const [isDark, setIsDark] = useState(true);
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    const savedTheme = localStorage.getItem("recontrace-theme");

    if (savedTheme === "light") {
      document.documentElement.dataset.theme = "light";
      setIsDark(false);
    } else {
      document.documentElement.removeAttribute("data-theme");
      setIsDark(true);
    }

    setMounted(true);
  }, []);

  const toggleTheme = () => {
    const nextIsDark = !isDark;

    setIsDark(nextIsDark);

    if (nextIsDark) {
      document.documentElement.removeAttribute("data-theme");
      localStorage.setItem("recontrace-theme", "dark");
    } else {
      document.documentElement.dataset.theme = "light";
      localStorage.setItem("recontrace-theme", "light");
    }
  };

  // Prevent hydration mismatch
  if (!mounted) {
    return <div className="theme-toggle-placeholder" />;
  }

  return (
    <button
      type="button"
      className={`theme-toggle ${isDark ? "dark" : "light"}`}
      onClick={toggleTheme}
      aria-label={`Switch to ${isDark ? "light" : "dark"} mode`}
      title={`Switch to ${isDark ? "light" : "dark"} mode`}
    >
      <span className="theme-icon">
        {isDark ? (
          <svg
            className="moon-icon"
            viewBox="0 0 24 24"
            aria-hidden="true"
          >
            <path d="M21 14.2A8.5 8.5 0 0 1 9.8 3a8.5 8.5 0 1 0 11.2 11.2Z" />
          </svg>
        ) : (
          <svg
            className="sun-icon"
            viewBox="0 0 24 24"
            aria-hidden="true"
          >
            <circle cx="12" cy="12" r="4" />
            <path d="M12 2v2.2M12 19.8V22M2 12h2.2M19.8 12H22M4.9 4.9l1.55 1.55M17.55 17.55l1.55 1.55M19.1 4.9l-1.55 1.55M6.45 17.55L4.9 19.1" />
          </svg>
        )}
      </span>
    </button>
  );
}