import { useEffect, useState } from "react";
import { getFindings, getProducts, signIn, type Finding, type ProductPage, type User } from "./api";
import { FindingsTable } from "./FindingsTable";

// The seeded users. This is a practice app, so there are no passwords.
const USERS = [
  { email: "asha@kestrel.example", name: "Asha", role: "analyst", org: "Kestrel Cameras" },
  { email: "meera@kestrel.example", name: "Meera", role: "viewer", org: "Kestrel Cameras" },
  { email: "ravi@tidewater.example", name: "Ravi", role: "analyst", org: "Tidewater Routers" },
];

const messageOf = (error: unknown) => (error instanceof Error ? error.message : "Something went wrong");

export function App() {
  const [user, setUser] = useState<User | null>(null);
  const [error, setError] = useState("");

  function choose(email: string) {
    setError("");
    signIn(email).then(
      (session) => setUser(session.user),
      (e) => setError(messageOf(e)),
    );
  }

  if (user) {
    return (
      <>
        <header className="island">
          <h1 className="brand">Devshelf</h1>
          <p className="who">
            Signed in as {user.name} ({user.role}, {user.org})
          </p>
          <button className="on-island" onClick={() => setUser(null)}>
            Switch user
          </button>
        </header>
        <main>
          <Products />
          <Findings />
        </main>
      </>
    );
  }

  return (
    <>
      <header className="island">
        <h1 className="brand">Devshelf</h1>
        <p className="who">Device firmware inventory</p>
      </header>
      <main className="narrow">
        <section className="card">
          <h2>Sign in</h2>
          <p className="muted">
            This is a practice app with made-up users. There are no passwords: choose a user to sign in.
          </p>
          <ul className="people">
            {USERS.map((u) => (
              <li key={u.email}>
                <span className="avatar" aria-hidden="true">
                  {u.name[0]}
                </span>
                <span className="person">
                  <strong>{u.name}</strong>
                  <span className="muted">
                    {u.role}, {u.org}, {u.email}
                  </span>
                </span>
                <button className="primary" onClick={() => choose(u.email)}>
                  Sign in as {u.name}
                </button>
              </li>
            ))}
          </ul>
          {error && (
            <p role="alert" className="notice">
              {error}
            </p>
          )}
        </section>
      </main>
    </>
  );
}

function Products() {
  const [q, setQ] = useState("");
  const [page, setPage] = useState(1);
  const [data, setData] = useState<ProductPage | null>(null);
  const [error, setError] = useState("");

  useEffect(() => {
    let current = true; // ignore a reply that arrives after a newer request
    getProducts(q, page).then(
      (result) => {
        if (current) {
          setData(result);
          setError("");
        }
      },
      (e) => {
        if (current) setError(messageOf(e));
      },
    );
    return () => {
      current = false;
    };
  }, [q, page]);

  return (
    <section className="card">
      <div className="card-head">
        <h2>Products</h2>
        <div className="field">
          <label htmlFor="product-search">Search products</label>
          <input
            id="product-search"
            type="search"
            value={q}
            onChange={(e) => {
              setQ(e.target.value);
              setPage(1);
            }}
          />
        </div>
      </div>
      {error ? (
        <p role="alert" className="notice">
          {error}
        </p>
      ) : !data ? (
        <p className="muted">Loading products.</p>
      ) : data.items.length === 0 ? (
        <p className="muted">No products found.</p>
      ) : (
        <ul className="rows">
          {data.items.map((p) => (
            <li key={p.id}>
              <span>{p.name}</span>
              <span className="meta muted">
                <code>{p.sku}</code>
                <span>{p.deviceType}</span>
              </span>
            </li>
          ))}
        </ul>
      )}
      <p className="pager">
        <button disabled={page <= 1} onClick={() => setPage(page - 1)}>
          Previous
        </button>
        <span>Page {page}</span>
        <button disabled={!data?.hasNext} onClick={() => setPage(page + 1)}>
          Next
        </button>
      </p>
    </section>
  );
}

function Findings() {
  const [findings, setFindings] = useState<Finding[] | null>(null);
  const [error, setError] = useState("");

  useEffect(() => {
    getFindings().then(setFindings, (e) => setError(messageOf(e)));
  }, []);

  return (
    <section className="card">
      <h2>Findings</h2>
      {error ? (
        <p role="alert" className="notice">
          {error}
        </p>
      ) : findings ? (
        <FindingsTable findings={findings} />
      ) : (
        <p className="muted">Loading findings.</p>
      )}
    </section>
  );
}
