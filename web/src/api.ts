export type Severity = "critical" | "high" | "medium" | "low";

export type User = { name: string; email: string; role: "viewer" | "analyst" | "admin"; org: string };
export type Product = { id: string; name: string; sku: string; deviceType: string };
export type Finding = { id: string; title: string; severity: Severity; component: string; product: string };

export type Session = { user: User; csrfToken: string };
export type ProductPage = { items: Product[]; page: number; pageSize: number; hasNext: boolean };

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(path, { credentials: "same-origin", ...init });
  const body = await response.json().catch(() => null);
  if (!response.ok) {
    throw new Error(body?.error ?? `Request failed with status ${response.status}`);
  }
  return body as T;
}

export function signIn(email: string): Promise<Session> {
  return request("/api/session", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email }),
  });
}

// Pages start at 1.
export function getProducts(q: string, page: number): Promise<ProductPage> {
  const query = new URLSearchParams({ q, page: String(page) });
  return request(`/api/products?${query}`);
}

export async function getFindings(): Promise<Finding[]> {
  const body = await request<{ items: Finding[] }>("/api/findings");
  return body.items;
}
