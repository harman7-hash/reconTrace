"use client";

import { useState, useEffect, type FormEvent } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";

const API = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

type Host = {
  host: string;
  status: string;
  statusColor: string;
  label: string;
  value: string;
  metric: string;
  metricValue: string;
  bars: string[];
  barColor: string;
};

const HOSTS: Host[] = [
  {
    host: "host-07 / prod-east",
    status: "WARMUP",
    statusColor: "text-amber-400",
    label: "baseline",
    value: "Day 9 / 14",
    metric: "status",
    metricValue: "learning traffic",
    bars: ["h-3", "h-3", "h-4", "h-6", "h-3", "h-5"],
    barColor: "bg-amber-400",
  },
  {
    host: "host-12 / staging",
    status: "READY",
    statusColor: "text-emerald-400",
    label: "inference",
    value: "Active Model",
    metric: "anomaly score",
    metricValue: "0.02 (Normal)",
    bars: ["h-2", "h-3", "h-2", "h-4", "h-3", "h-5"],
    barColor: "bg-emerald-400",
  },
  {
    host: "host-03 / db-primary",
    status: "RETRAIN",
    statusColor: "text-sky-400",
    label: "pipeline",
    value: "Midnight Sync",
    metric: "daily update",
    metricValue: "Scheduled 00:00",
    bars: ["h-5", "h-6", "h-4", "h-6", "h-5", "h-6"],
    barColor: "bg-sky-400",
  },
];

function HostCard({
  host,
  className = "",
}: {
  host: Host;
  className?: string;
}) {
  return (
    <div
      className={`w-full max-w-88.75 bg-zinc-900/90 p-4 font-mono text-xs leading-relaxed transition-all duration-300 hover:rotate-0 ${className}`}
    >
      <div className="mb-3 flex justify-between border-b border-zinc-800 pb-3 text-zinc-400">
        <span>{host.host}</span>
        <span className={`font-medium ${host.statusColor}`}>{host.status}</span>
      </div>
      <div className="mb-1 flex justify-between text-zinc-400">
        <span>{host.label}</span>
        <span className="text-zinc-100">{host.value}</span>
      </div>
      <div className="flex justify-between text-zinc-400">
        <span>{host.metric}</span>
        <span className="text-zinc-100">{host.metricValue}</span>
      </div>
      <div className="mt-3 flex h-6 items-end gap-1" aria-hidden="true">
        {host.bars.map((height, barIndex) => (
          <i className={`flex-1 ${height} ${host.barColor}`} key={`${host.host}-${barIndex}`} />
        ))}
      </div>
    </div>
  );
}

export default function LoginPage() {
  const router = useRouter();
  const [mode, setMode] = useState<"login" | "signup">("login");
  const [username, setUsername] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [remember, setRemember] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");
  const [loading, setLoading] = useState(false);

  const isLogin = mode === "login";

  useEffect(() => {
    if (typeof window !== "undefined") {
      const params = new URLSearchParams(window.location.search);
      if (params.get("mode") === "signup") {
        setMode("signup");
      }
    }
  }, []);

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    setError("");
    setSuccess("");

    if (!isLogin && username.trim().length < 2) {
      return setError("Username must be at least 2 characters.");
    }
    if (!/^\S+@\S+\.\S+$/.test(email)) {
      return setError("Enter a valid work email address.");
    }
    if (password.length < 8) {
      return setError("Password must be at least 8 characters.");
    }

    setLoading(true);

    try {
      const endpoint = isLogin ? `${API}/api/login` : `${API}/api/signupr`;
      const payload = isLogin
        ? { email, password }
        : { username: username.trim(), email, password };

      const res = await fetch(endpoint, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify(payload),
      });

      const data = await res.json().catch(() => null);

      if (!res.ok) {
        throw new Error(
          data?.error || data?.message || "Authentication failed. Please verify your credentials."
        );
      }

      if (data?.user?.api_key) {
        try {
          localStorage.setItem("recontrace_api_key", data.user.api_key);
          localStorage.setItem("recontrace_user", JSON.stringify(data.user));
        } catch {
          // Ignore localStorage write errors in private browsing
        }
      }

      if (!isLogin) {
        setSuccess("Account created successfully! Redirecting...");
        setTimeout(() => {
          router.push("/dashboard");
        }, 1200);
      } else {
        router.push("/dashboard");
      }
    } catch (err) {
      setError(
        err instanceof Error && err.message
          ? err.message
          : "Something went wrong connecting to the auth service. Try again."
      );
    } finally {
      setLoading(false);
    }
  }

  return (
    <main className="h-screen w-full bg-zinc-950 font-sans text-zinc-100 antialiased selection:bg-amber-400 selection:text-zinc-950 lg:grid lg:grid-cols-2 overflow-hidden">
      {/* Left Column: Telemetry & Host Showcase */}
      <section
        className="relative hidden flex-col justify-between overflow-hidden border-b border-zinc-800/80 p-8 lg:flex lg:border-b-0 lg:border-r lg:p-14"
        aria-hidden="true"
      >
        {/* Subtle ambient blur glow
        <div className="pointer-events-none absolute -left-20 -top-20 h-96 w-96 rounded-full bg-amber-500/5 blur-3xl" />
        <div className="pointer-events-none absolute -bottom-20 -right-20 h-96 w-96 rounded-full bg-emerald-500/5 blur-3xl" /> */}

        {/* Top Logo */}
        <div className="relative -top-7.5 -left-5 z-10">
          <Link
            href="/"
            className="inline-flex items-center gap-2.5 text-xl font-bold tracking-wide text-zinc-100 transition-colors hover:text-white"
          >
            <span>reconTrace</span>
          </Link>
        </div>

        {/* Center: Live Host Cards Stack */}
        <div className="relative -top-3 z-10 my-auto flex flex-col items-center justify-center py-10">
          <div className="relative w-full max-w-sm space-y-4">
            <HostCard host={HOSTS[0]} className="rotate-2" />
            <HostCard host={HOSTS[1]} className="rotate-2 -mt-5 relative -left-5" />
            <HostCard host={HOSTS[2]} className="-rotate-2 -mt-2 relative left-5 -top-4" />
          </div>
        </div>

        {/* Bottom Tagline */}
        <div className="relative z-10 max-w-md border-l-2 border-amber-400 pl-4">
          <p className="font-mono text-xs uppercase tracking-widest text-zinc-400 leading-relaxed">
            Only the hosts that need attention rise to the top.
          </p>
        </div>
      </section>

      {/* Right Column: Authentication Panel */}
      <section className="flex min-h-screen items-center justify-center px-6 py-12 sm:px-12 lg:px-16">
        <div className="w-full max-w-md">
          {/* Mobile Header */}
          <div className="mb-8 lg:hidden">
            <Link
              href="/"
              className="inline-flex items-center gap-2 text-xl font-bold tracking-wide text-zinc-100"
            >
              <span>reconTrace</span>
              <span className="h-2 w-2 rounded-full bg-emerald-400 shadow-[0_0_8px_rgba(52,211,153,0.8)]" />
            </Link>
          </div>

          <h1 className="text-3xl font-bold uppercase tracking-tight text-zinc-100 sm:text-4xl">
            {isLogin ? "Back to the control room." : "Create your account."}
          </h1>
          <p className="mb-8 mt-3 text-sm text-zinc-400 leading-relaxed">
            {isLogin
              ? "Sign in to see which hosts need your attention right now."
              : "Start monitoring your fleet in under sixty seconds."}
          </p>

          {/* Mode Switcher Tabs */}
          <div
            className="mb-7 grid grid-cols-2 border border-zinc-800 bg-zinc-900/60 p-1 font-mono text-xs"
            role="tablist"
          >
            <button
              type="button"
              role="tab"
              aria-selected={isLogin}
              onClick={() => {
                setMode("login");
                setError("");
                setSuccess("");
              }}
              className={`py-2.5 font-medium transition-colors ${
                isLogin
                  ? "bg-zinc-100 text-zinc-950 shadow-sm"
                  : "text-zinc-400 hover:text-zinc-200"
              }`}
            >
              Log in
            </button>
            <button
              type="button"
              role="tab"
              aria-selected={!isLogin}
              onClick={() => {
                setMode("signup");
                setError("");
                setSuccess("");
              }}
              className={`py-2.5 font-medium transition-colors ${
                !isLogin
                  ? "bg-zinc-100 text-zinc-950 shadow-sm"
                  : "text-zinc-400 hover:text-zinc-200"
              }`}
            >
              Create account
            </button>
          </div>

          <form onSubmit={onSubmit} noValidate>
            {/* 
              Username Field 
              Uses a CSS grid transition to smoothly animate from 0 height to full height.
              This prevents sudden layout shifts and scrollbar popping.
            */}
            <div
              className={`grid transition-[grid-template-rows,opacity] duration-300 ease-in-out ${
                isLogin ? "grid-rows-[0fr] opacity-0" : "grid-rows-[1fr] opacity-100"
              }`}
              aria-hidden={isLogin}
            >
              <div className="overflow-hidden">
                <div className="mb-5">
                  <label
                    htmlFor="username"
                    className="mb-2 block font-mono text-xs uppercase tracking-wider text-zinc-400"
                  >
                    Username
                  </label>
                  <input
                    id="username"
                    type="text"
                    autoComplete="username"
                    placeholder="e.g. dev_operator"
                    value={username}
                    onChange={(e) => setUsername(e.target.value)}
                    disabled={isLogin}
                    className="h-11 w-full border border-zinc-800 bg-zinc-900 px-3.5 font-sans text-sm text-zinc-100 placeholder:text-zinc-600 focus:border-amber-400 focus:outline-none focus:ring-1 focus:ring-amber-400 transition-colors disabled:opacity-50"
                  />
                </div>
              </div>
            </div>

            <div className="mb-5">
              <label
                htmlFor="email"
                className="mb-2 block font-mono text-xs uppercase tracking-wider text-zinc-400"
              >
                Work email
              </label>
              <input
                id="email"
                type="email"
                autoComplete="email"
                placeholder="you@company.com"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                className="h-11 w-full border border-zinc-800 bg-zinc-900 px-3.5 font-sans text-sm text-zinc-100 placeholder:text-zinc-600 focus:border-amber-400 focus:outline-none focus:ring-1 focus:ring-amber-400 transition-colors"
              />
            </div>

            <div className="mb-5">
              <label
                htmlFor="password"
                className="mb-2 block font-mono text-xs uppercase tracking-wider text-zinc-400"
              >
                Password
              </label>
              <div className="relative">
                <input
                  id="password"
                  type={showPassword ? "text" : "password"}
                  autoComplete={isLogin ? "current-password" : "new-password"}
                  placeholder="At least 8 characters"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  className="h-11 w-full border border-zinc-800 bg-zinc-900 pl-3.5 pr-11 font-sans text-sm text-zinc-100 placeholder:text-zinc-600 focus:border-amber-400 focus:outline-none focus:ring-1 focus:ring-amber-400 transition-colors"
                />
                <button
                  type="button"
                  onClick={() => setShowPassword((s) => !s)}
                  className="absolute right-0 top-0 flex h-11 items-center justify-center px-3 text-zinc-500 transition-colors hover:text-zinc-300"
                  aria-label={showPassword ? "Hide password" : "Show password"}
                >
                  {showPassword ? (
                    <svg
                      className="h-4 w-4"
                      viewBox="0 0 24 24"
                      fill="none"
                      stroke="currentColor"
                      strokeWidth="2"
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      aria-hidden="true"
                    >
                      <path d="M9.88 9.88a3 3 0 1 0 4.24 4.24" />
                      <path d="M10.73 5.08A10.43 10.43 0 0 1 12 5c7 0 10 7 10 7a13.16 13.16 0 0 1-1.67 2.68" />
                      <path d="M6.61 6.61A13.526 13.526 0 0 0 2 12s3 7 10 7a9.74 9.74 0 0 0 5.39-1.61" />
                      <line x1="2" x2="22" y1="2" y2="22" />
                    </svg>
                  ) : (
                    <svg
                      className="h-4 w-4"
                      viewBox="0 0 24 24"
                      fill="none"
                      stroke="currentColor"
                      strokeWidth="2"
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      aria-hidden="true"
                    >
                      <path d="M2 12s3-7 10-7 10 7 10 7-3 7-10 7-10-7-10-7Z" />
                      <circle cx="12" cy="12" r="3" />
                    </svg>
                  )}
                </button>
              </div>
            </div>

            <div className="mb-5 flex items-center justify-between font-mono text-xs">
              <label className="flex cursor-pointer select-none items-center gap-2 text-zinc-400">
                <input
                  type="checkbox"
                  checked={remember}
                  onChange={(e) => setRemember(e.target.checked)}
                  className="h-4 w-4 border-zinc-700 bg-zinc-900 text-amber-400 accent-amber-400 focus:ring-0 focus:ring-offset-0"
                />
                <span>Keep me signed in</span>
              </label>
              {isLogin && (
                <Link
                  href="/forgot-password"
                  className="text-zinc-400 underline underline-offset-2 transition-colors hover:text-zinc-200"
                >
                  Forgot password?
                </Link>
              )}
            </div>

            {error && (
              <div
                className="mb-5 flex items-start gap-2.5 border border-red-500/40 bg-red-950/40 px-3.5 py-2.5 font-mono text-xs text-red-300"
                role="alert"
              >
                <span className="font-bold text-red-400">!</span>
                <span>{error}</span>
              </div>
            )}

            {success && (
              <div
                className="mb-5 flex items-start gap-2.5 border border-emerald-500/40 bg-emerald-950/40 px-3.5 py-2.5 font-mono text-xs text-emerald-300"
                role="status"
              >
                <span className="font-bold text-emerald-400">✓</span>
                <span>{success}</span>
              </div>
            )}

            <button
              type="submit"
              disabled={loading}
              className="mb-6 flex h-11 w-full items-center justify-center bg-amber-400 font-sans text-sm font-medium tracking-wide text-zinc-950 transition hover:bg-amber-300 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-amber-400 disabled:cursor-not-allowed disabled:opacity-60"
            >
              {loading ? (
                <span className="inline-flex items-center gap-2">
                  <svg
                    className="h-4 w-4 animate-spin text-zinc-950"
                    viewBox="0 0 24 24"
                    fill="none"
                    aria-hidden="true"
                  >
                    <circle
                      className="opacity-25"
                      cx="12"
                      cy="12"
                      r="10"
                      stroke="currentColor"
                      strokeWidth="4"
                    />
                    <path
                      className="opacity-75"
                      fill="currentColor"
                      d="M4 12a8 8 0 018-8v8H4z"
                    />
                  </svg>
                  Please wait...
                </span>
              ) : isLogin ? (
                "Log in"
              ) : (
                "Create account"
              )}
            </button>

            

            
          </form>
        </div>
      </section>
    </main>
  );
}