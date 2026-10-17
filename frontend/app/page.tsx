import Link from "next/link";

const statusCards = [
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

const features = [
  {
    label: "/ isolation",
    title: "A dedicated ML model per server",
    description:
      "No global averages or blunt threshold rules. Every host runs its own tailored machine learning model trained exclusively on its unique workload.",
    accent: "border-l-emerald-400",
  },
  {
    label: "/ baseline",
    title: "14-day zero-noise warmup",
    description:
      "New servers undergo a 14-day observation phase to map normal operational patterns, eliminating false alarms before predictions go live.",
    accent: "border-l-amber-400",
  },
  {
    label: "/ adaptation",
    title: "Daily midnight retraining",
    description:
      "Models automatically retrain and recalibrate every single night at midnight, seamlessly absorbing gradual traffic growth and daily system changes.",
    accent: "border-l-sky-400",
  },
];

const steps = [
  ["01 — WARMUP", "14-day data collection", "Every new server spends 14 days gathering deep metrics telemetry to build its unique behavioral baseline."],
  ["02 — PREDICT", "Active inference", "Custom AI models run live checks against incoming telemetry, catching deviations and anomalies instantly."],
  ["03 — RETRAIN", "Daily midnight sync", "Models update automatically every night at midnight to continuously adapt to traffic evolution without manual intervention."],
];

export default function ReconTracePage() {
  return (
    <main className="min-h-screen bg-zinc-950 font-sans text-zinc-100 antialiased selection:bg-amber-400 selection:text-zinc-950">
      <div className="mx-auto max-w-6xl px-6">
        <nav className="flex h-20 items-center justify-between border-b border-zinc-800" aria-label="Main navigation">
          <Link className="font-bold tracking-wide text-2xl" href="#top">reconTrace</Link>
          <div className="hidden gap-12 text-sm text-zinc-400 md:flex">
            <Link className="transition-colors hover:text-zinc-100" href="#features">Features</Link>
            <Link className="transition-colors hover:text-zinc-100" href="#method">Method</Link>
            <Link className="transition-colors hover:text-zinc-100" href="#console">Get started</Link>
          </div>
          <Link className="justify-self-end bg-amber-400 px-5 py-3 text-sm font-medium text-zinc-950 transition hover:bg-amber-300 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-amber-400" href="/login">
            Login
          </Link>
        </nav>

        <header className="grid grid-cols-1 items-center gap-10 py-16 md:py-24 lg:grid-cols-2" id="top">
          <div>
            <h1 className="m-0 max-w-3xl text-5xl font-bold uppercase leading-none tracking-tight sm:text-6xl lg:text-5xl">
              AI infrastructure monitoring, trained per host.
            </h1>
            <p className="mb-8 mt-7 max-w-md text-lg text-zinc-400">
              Every server gets its own machine learning model. A 14-day warmup maps your normal traffic, followed by daily midnight retraining to prevent drift.
            </p>
            <div className="flex flex-wrap gap-4">
              <Link className="border border-zinc-100 bg-zinc-100 px-5 py-3.5 text-sm font-medium text-zinc-950 transition-colors hover:bg-white focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-amber-400" href="/signup">
                Enter the console
              </Link>
              <Link className="border border-zinc-800 px-5 py-3.5 text-sm font-medium transition-colors hover:border-zinc-500 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-amber-400" href="#method">
                See how it works
              </Link>
            </div>
          </div>

          <div className="relative mx-auto h-96 w-full max-w-sm" role="img" aria-label="Sample host status cards showing warmup, inference, and daily midnight retraining states">
            {statusCards.map((card, index) => (
              <div
                className={`absolute w-75 border border-zinc-800 bg-zinc-900 p-4 font-mono text-xs leading-relaxed shadow-xl ${index === 0 ? "left-10 -top-3.75 z-30 -rotate-2" : index === 1 ? "left-0 top-35 z-10 rotate-2" : "left-12 top-70 z-20 rotate-2"}`}
                key={card.host}
              >
                <div className="mb-3 flex justify-between border-b border-zinc-800 pb-3 text-zinc-400">
                  <span>{card.host}</span>
                  <span className={`font-medium ${card.statusColor}`}>{card.status}</span>
                </div>
                <div className="mb-1 flex justify-between text-zinc-400"><span>{card.label}</span><span className="text-zinc-100">{card.value}</span></div>
                <div className="flex justify-between text-zinc-400"><span>{card.metric}</span><span className="text-zinc-100">{card.metricValue}</span></div>
                <div className="mt-3 flex h-6 items-end gap-1" aria-hidden="true">
                  {card.bars.map((height, barIndex) => (
                    <i className={`flex-1 ${height} ${card.barColor}`} key={`${card.host}-${barIndex}`} />
                  ))}
                </div>
              </div>
            ))}
          </div>
        </header>

        <section className="border-t border-zinc-800 py-16 md:py-20" id="method">
          <p className="mb-3 font-mono text-xs uppercase tracking-widest text-amber-400">
            // model lifecycle
          </p>
          <h2 className="m-0 max-w-2xl text-4xl font-bold uppercase leading-tight sm:text-5xl">
            From 14-day baseline to daily midnight updates.
          </h2>
          <div className="mt-12 grid grid-cols-1 border border-zinc-800 bg-zinc-900 md:grid-cols-3">
            {steps.map(([label, title, description], index) => (
              <div className={`p-6 ${index < 2 ? "border-b border-zinc-800 md:border-b-0 md:border-r" : ""}`} key={label}>
                <p className="mb-4 font-mono text-xs text-zinc-500">{label}</p>
                <h3 className="mb-2 text-lg font-semibold">{title}</h3>
                <p className="text-sm leading-relaxed text-zinc-400">{description}</p>
              </div>
            ))}
          </div>
        </section>

        <section className="border-t border-zinc-800 py-16 md:py-20" id="features">
          <p className="mb-3 font-mono text-xs uppercase tracking-widest text-amber-400">
            // What you get
          </p>
          <h2 className="m-0 max-w-2xl text-4xl font-bold uppercase leading-tight sm:text-5xl">
            Zero configuration. Tailored server intelligence.
          </h2>
          <div className="mt-12 grid grid-cols-1 border border-zinc-800 md:grid-cols-3">
            {features.map((feature, index) => (
              <article className={`border-l-4 bg-zinc-950 p-7 ${feature.accent} ${index < 2 ? "border-b border-zinc-800 md:border-b-0 md:border-r" : ""}`} key={feature.label}>
                <p className="mb-4 font-mono text-xs text-zinc-500">{feature.label}</p>
                <h3 className="mb-2 text-lg font-semibold">{feature.title}</h3>
                <p className="text-sm leading-relaxed text-zinc-400">{feature.description}</p>
              </article>
            ))}
          </div>
        </section>

        <section className="border-t border-zinc-800 py-16 md:py-20" id="console">
          <p className="mb-3 font-mono text-xs uppercase tracking-widest text-amber-400">
            // Get started
          </p>
          <h2 className="m-0 text-4xl font-bold uppercase leading-tight sm:text-5xl">
            Deploy your first agent in under a minute.
          </h2>
          <div className="mt-8 flex flex-wrap items-center justify-between gap-7">
            <p className="max-w-sm text-sm leading-relaxed text-zinc-400">
              One install command starts telemetry collection. reconTrace automatically begins the 14-day baseline warmup for your model.
            </p>
            <Link className="border border-zinc-100 bg-zinc-100 px-5 py-4 text-sm font-medium text-zinc-950 transition-colors hover:bg-white focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-amber-400" href="/signup">
              Enter the console
            </Link>
          </div>
        </section>

        <footer className="h-20 border-t border-zinc-800" />
      </div>
    </main>
  );
}