"use client";

import { useState, useEffect, type FormEvent } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";

const API = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

type Tone = "ok" | "watch" | "crit";
type Row = { label: string; filled: number; value: string };
type Host = { name: string; tone: Tone; status: string; rows: Row[] };

const HOSTS: Host[] = [
  {
    name: "host-07 / prod-east",
    tone: "ok",
    status: "OK",
    rows: [
      { label: "cpu", filled: 2, value: "34%" },
      { label: "mem", filled: 3, value: "4.1 / 16 GB" },
    ],
  },
  {
    name: "host-12 / api-gateway",
    tone: "watch",
    status: "WATCH",
    rows: [
      { label: "cpu", filled: 5, value: "81%" },
      { label: "net", filled: 4, value: "12.8 MB/s" },
    ],
  },
  {
    name: "host-03 / db-primary",
    tone: "crit",
    status: "CRIT",
    rows: [
      { label: "disk", filled: 6, value: "97.2% full" },
      { label: "iowait", filled: 5, value: "340ms" },
    ],
  },
];

const SEGMENTS = 6;

const toneStyles: Record<Tone, { badge: string; bar: string }> = {
  ok: {
    badge: "text-emerald-400 bg-emerald-950/60 border-emerald-800/60",
    bar: "bg-emerald-400",
  },
  watch: {
    badge: "text-amber-400 bg-amber-950/60 border-amber-800/60",
    bar: "bg-amber-400",
  },
  crit: {
    badge: "text-red-400 bg-red-950/60 border-red-800/60",
    bar: "bg-red-400",
  },
};

function HostCard({
  host,
  className = "",
}: {
  host: Host;
  className?: string;
}) {
  const styles = toneStyles[host.tone];

  return (
    <div
      className={`w-full max-w-md border border-zinc-800 bg-zinc-900/90 p-4 font-mono text-xs leading-relaxed shadow-2xl backdrop-blur-md transition-all duration-300 hover:rotate-0 hover:border-zinc-700 ${className}`}
    >
      <div className="mb-3 flex items-center justify-between border-b border-zinc-800/80 pb-2.5">
        <span className="font-medium tracking-tight text-zinc-300">{host.name}</span>
        <span
          className={`border px-2 py-0.5 text-[10px] font-semibold uppercase tracking-wider ${styles.badge}`}
        >
          {host.status}
        </span>
      </div>

      <div className="space-y-2.5">
        {host.rows.map((r) => (
          <div key={r.label} className="grid grid-cols-[44px_1fr_96px] items-center gap-3">
            <span className="font-medium text-[11px] uppercase text-zinc-500">{r.label}</span>
            <div className="flex items-center gap-1.5" aria-hidden="true">
              {Array.from({ length: SEGMENTS }, (_, i) => (
                <i
                  key={i}
                  className={`h-2 flex-1 -skew-x-12 transition-colors duration-200 ${
                    i < r.filled ? styles.bar : "bg-zinc-800/80"
                  }`}
                />
              ))}
            </div>
            <span className="text-right tabular-nums text-zinc-300">{r.value}</span>
          </div>
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
      const endpoint = isLogin ? `${API}/api/login` : `${API}/api/register`;
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
    <main className="min-h-screen bg-zinc-950 font-sans text-zinc-100 antialiased selection:bg-amber-400 selection:text-zinc-950 lg:grid lg:grid-cols-2">
      {/* Left Column: Telemetry & Host Showcase */}
      <section
        className="relative hidden flex-col justify-between overflow-hidden border-b border-zinc-800/80 p-8 lg:flex lg:border-b-0 lg:border-r lg:p-14"
        aria-hidden="true"
      >
        {/* Subtle ambient blur glow */}
        <div className="pointer-events-none absolute -left-20 -top-20 h-96 w-96 rounded-full bg-amber-500/5 blur-3xl" />
        <div className="pointer-events-none absolute -bottom-20 -right-20 h-96 w-96 rounded-full bg-emerald-500/5 blur-3xl" />

        {/* Top Logo */}
        <div className="relative z-10">
          <Link
            href="/"
            className="inline-flex items-center gap-2.5 text-xl font-bold tracking-wide text-zinc-100 transition-colors hover:text-white"
          >
            <span>reconTrace</span>
            {/* <span className="h-2 w-2 rounded-full bg-emerald-400 shadow-[0_0_8px_rgba(52,211,153,0.8)]" /> */}
          </Link>
        </div>

        {/* Center: Live Host Cards Stack */}
        <div className="relative z-10 my-auto flex flex-col items-center justify-center py-10">
          <div className="relative w-full max-w-md space-y-4">
            <HostCard host={HOSTS[0]} className="-rotate-2" />
            <HostCard host={HOSTS[1]} className="rotate-1 -mt-2" />
            <HostCard host={HOSTS[2]} className="-rotate-1 -mt-2" />
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

          <form onSubmit={onSubmit} noValidate className="space-y-5">
            {!isLogin && (
              <div>
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
                  className="h-11 w-full border border-zinc-800 bg-zinc-900 px-3.5 font-sans text-sm text-zinc-100 placeholder:text-zinc-600 focus:border-amber-400 focus:outline-none focus:ring-1 focus:ring-amber-400 transition-colors"
                />
              </div>
            )}

            <div>
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

            <div>
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
                  className="absolute right-0 top-0 flex h-11 items-center justify-center px-3 text-zinc-500 hover:text-zinc-300 transition-colors"
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

            <div className="flex items-center justify-between font-mono text-xs">
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
                  className="text-zinc-400 underline underline-offset-2 hover:text-zinc-200 transition-colors"
                >
                  Forgot password?
                </Link>
              )}
            </div>

            {error && (
              <div
                className="flex items-start gap-2.5 border border-red-500/40 bg-red-950/40 px-3.5 py-2.5 font-mono text-xs text-red-300"
                role="alert"
              >
                <span className="font-bold text-red-400">!</span>
                <span>{error}</span>
              </div>
            )}

            {success && (
              <div
                className="flex items-start gap-2.5 border border-emerald-500/40 bg-emerald-950/40 px-3.5 py-2.5 font-mono text-xs text-emerald-300"
                role="status"
              >
                <span className="font-bold text-emerald-400">✓</span>
                <span>{success}</span>
              </div>
            )}

            <button
              type="submit"
              disabled={loading}
              className="flex h-11 w-full items-center justify-center bg-amber-400 font-sans text-sm font-medium tracking-wide text-zinc-950 transition hover:bg-amber-300 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-amber-400 disabled:cursor-not-allowed disabled:opacity-60"
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

            <div className="relative my-6 text-center">
              <div className="absolute inset-0 flex items-center">
                <div className="w-full border-t border-zinc-800" />
              </div>
              <span className="relative bg-zinc-950 px-3 font-mono text-xs uppercase text-zinc-600">
                or
              </span>
            </div>

            <button
              type="button"
              onClick={() => {
                window.location.href = `${API}/auth/sso`;
              }}
              className="flex h-11 w-full items-center justify-center gap-2 border border-zinc-800 bg-zinc-900/60 font-sans text-sm font-medium text-zinc-200 transition hover:border-zinc-700 hover:bg-zinc-900 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-amber-400"
            >
              <svg
                className="h-4 w-4 text-zinc-400"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                strokeWidth="2"
                strokeLinecap="round"
                strokeLinejoin="round"
                aria-hidden="true"
              >
                <rect width="18" height="11" x="3" y="11" rx="2" ry="2" />
                <path d="M7 11V7a5 5 0 0 1 10 0v4" />
              </svg>
              <span>Continue with SSO</span>
            </button>
          </form>
        </div>
      </section>
    </main>
  );
}