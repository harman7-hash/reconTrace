import Link from "next/link";
import ThemeToggle from "./components/ThemeToggle";

const features = [
  {
    number: "01",
    title: "Real-time visibility",
    description:
      "Track CPU, memory, disk, network, uptime, and service health from one clean operational view.",
  },
  {
    number: "02",
    title: "Instant anomaly signals",
    description:
      "Surface unusual resource behavior before it becomes an incident with configurable thresholds and health states.",
  },
  {
    number: "03",
    title: "Trace every event",
    description:
      "Connect infrastructure metrics with logs and service events so your team can move from signal to root cause faster.",
  },
];

const metrics = [
  { label: "API Gateway", value: "24 ms", state: "Healthy", width: "82%" },
  { label: "Database", value: "61%", state: "Healthy", width: "61%" },
  { label: "Worker Pool", value: "78%", state: "Watch", width: "78%" },
];

export default function Home() {
  return (
    <main className="site-shell">
      <nav className="navbar">

        <Link href="/" className="brand">
          <span className="brand-mark">
            <span />
            <span />
            <span />
          </span>

          <span>
            recon<span className="brand-accent">Trace</span>
          </span>
        </Link>


        <div className="nav-links">
          <a href="#features">Features</a>
          <a href="#how-it-works">How it works</a>
          <a href="#stack">Architecture</a>
        </div>


        <div className="nav-actions">

          {/* Theme toggle */}
          <ThemeToggle />

          <Link href="/login" className="nav-login">
            Login
          </Link>

          <Link
            href="/signup"
            className="button button-small"
          >
            Get started
          </Link>

        </div>

      </nav>

      <section className="hero">
        <div className="hero-copy">
          <div className="eyebrow">
            <span className="status-dot" />
            Infrastructure observability, simplified
          </div>

          <h1>
            Know what your
            <span> servers are doing.</span>
          </h1>

          <p className="hero-description">
            reconTrace gives you a real-time operational view of your servers,
            services, and infrastructure — without drowning your team in noise.
          </p>

          <div className="hero-actions">
            <Link href="#get-started" className="button">
              Start monitoring <span>→</span>
            </Link>
            <Link href="#how-it-works" className="button button-ghost">
              Explore the system
            </Link>
          </div>

          <div className="hero-meta">
            <div>
              <strong>99.99%</strong>
              <span>monitoring uptime</span>
            </div>
            <div>
              <strong>&lt; 1 min</strong>
              <span>setup to first signal</span>
            </div>
            <div>
              <strong>24/7</strong>
              <span>infrastructure watch</span>
            </div>
          </div>
        </div>

        <div className="hero-visual" aria-label="reconTrace monitoring dashboard preview">
          <div className="glow" />
          <div className="dashboard-window">
            <div className="window-top">
              <div className="window-dots">
                <i />
                <i />
                <i />
              </div>
              <span>reconTrace / overview</span>
              <span className="live-pill">
                <b /> LIVE
              </span>
            </div>

            <div className="dashboard-body">
              <div className="dashboard-header">
                <div>
                  <span className="muted-label">INFRASTRUCTURE</span>
                  <h3>System overview</h3>
                </div>
                <span className="healthy-pill">● All systems operational</span>
              </div>

              <div className="metric-grid">
                <div className="metric-card">
                  <span>CPU USAGE</span>
                  <strong>42.8%</strong>
                  <small>−4.2% from last hour</small>
                  <div className="sparkline">
                    <span style={{ height: "32%" }} />
                    <span style={{ height: "48%" }} />
                    <span style={{ height: "42%" }} />
                    <span style={{ height: "66%" }} />
                    <span style={{ height: "54%" }} />
                    <span style={{ height: "73%" }} />
                    <span style={{ height: "58%" }} />
                    <span style={{ height: "44%" }} />
                    <span style={{ height: "62%" }} />
                    <span style={{ height: "50%" }} />
                  </div>
                </div>

                <div className="metric-card">
                  <span>MEMORY</span>
                  <strong>7.2 GB</strong>
                  <small>of 16 GB allocated</small>
                  <div className="progress-track">
                    <span style={{ width: "45%" }} />
                  </div>
                </div>

                <div className="metric-card">
                  <span>NETWORK</span>
                  <strong>18.4 MB/s</strong>
                  <small>↑ 2.8% traffic</small>
                  <div className="network-bars">
                    <span style={{ height: "25%" }} />
                    <span style={{ height: "40%" }} />
                    <span style={{ height: "33%" }} />
                    <span style={{ height: "58%" }} />
                    <span style={{ height: "48%" }} />
                    <span style={{ height: "72%" }} />
                    <span style={{ height: "62%" }} />
                    <span style={{ height: "84%" }} />
                  </div>
                </div>
              </div>

              <div className="server-panel">
                <div className="panel-heading">
                  <span>ACTIVE SERVICES</span>
                  <span>STATUS</span>
                </div>

                {metrics.map((metric) => (
                  <div className="server-row" key={metric.label}>
                    <div className="service-name">
                      <span className="service-icon">↗</span>
                      <div>
                        <strong>{metric.label}</strong>
                        <small>prod-us-east-01</small>
                      </div>
                    </div>
                    <div className="service-load">
                      <div className="mini-track">
                        <span style={{ width: metric.width }} />
                      </div>
                      <span>{metric.value}</span>
                    </div>
                    <span
                      className={`service-status ${
                        metric.state === "Watch" ? "watch" : ""
                      }`}
                    >
                      ● {metric.state}
                    </span>
                  </div>
                ))}
              </div>
            </div>
          </div>
        </div>
      </section>

      <section className="signal-strip">
        <span>MONITOR</span>
        <b>→</b>
        <span>DETECT</span>
        <b>→</b>
        <span>TRACE</span>
        <b>→</b>
        <span>RESPOND</span>
      </section>

      <section className="features section" id="features">
        <div className="section-heading">
          <div>
            <span className="section-kicker">BUILT FOR OPERATIONS</span>
            <h2>Less guessing.<br />More signal.</h2>
          </div>
          <p>
            A focused monitoring layer for engineers who need to understand
            infrastructure health quickly and act with confidence.
          </p>
        </div>

        <div className="feature-grid">
          {features.map((feature) => (
            <article className="feature-card" key={feature.number}>
              <span className="feature-number">{feature.number}</span>
              <div className="feature-icon">
                {feature.number === "01" && "◌"}
                {feature.number === "02" && "⌁"}
                {feature.number === "03" && "⌘"}
              </div>
              <h3>{feature.title}</h3>
              <p>{feature.description}</p>
              <span className="feature-arrow">↗</span>
            </article>
          ))}
        </div>
      </section>

      <section className="workflow section" id="how-it-works">
        <div className="workflow-copy">
          <span className="section-kicker">FROM HOST TO INSIGHT</span>
          <h2>One pipeline.<br />Every signal connected.</h2>
          <p>
            Deploy a lightweight agent on your servers, stream telemetry to
            the Go backend, and let the Next.js dashboard turn raw metrics into
            an operational picture.
          </p>
          <div className="workflow-list">
            <div><span>01</span><p><strong>Collect</strong> — system and service metrics</p></div>
            <div><span>02</span><p><strong>Process</strong> — normalize and evaluate signals</p></div>
            <div><span>03</span><p><strong>Visualize</strong> — surface health in real time</p></div>
          </div>
        </div>

        <div className="pipeline">
          <div className="pipeline-node">
            <span className="node-icon">⌁</span>
            <div><strong>Server Agent</strong><small>metrics + events</small></div>
          </div>
          <div className="pipeline-line"><span /></div>
          <div className="pipeline-node">
            <span className="node-icon">◆</span>
            <div><strong>Go Backend</strong><small>ingest + process</small></div>
          </div>
          <div className="pipeline-line"><span /></div>
          <div className="pipeline-node active-node">
            <span className="node-icon">▦</span>
            <div><strong>reconTrace</strong><small>observe + respond</small></div>
          </div>
        </div>
      </section>

      <section className="stack section" id="stack">
        <div className="stack-card">
          <div>
            <span className="section-kicker">YOUR STACK, CONNECTED</span>
            <h2>Designed around the tools<br />you already build with.</h2>
          </div>
          <div className="tech-list">
            <span><b>Go</b> Backend</span>
            <span><b>Next.js</b> Frontend</span>
            <span><b>REST / WS</b> Transport</span>
            <span><b>Linux</b> Hosts</span>
          </div>
        </div>
      </section>

      <section className="cta section" id="get-started">
        <div className="cta-inner">
          <span className="section-kicker">START WITH VISIBILITY</span>
          <h2>Build a calmer<br /><span>way to operate.</span></h2>
          <p>Connect your first server and see what your infrastructure is really doing.</p>
          <Link href="#login" className="button">
            Enter reconTrace <span>→</span>
          </Link>
        </div>
      </section>

      <footer className="footer">
        <div className="brand">
          <span className="brand-mark">
            <span />
            <span />
            <span />
          </span>
          <span>recon<span className="brand-accent">Trace</span></span>
        </div>
        <span>Infrastructure visibility for modern engineering teams.</span>
        <span>© {new Date().getFullYear()} reconTrace</span>
      </footer>
    </main>
  );
}
