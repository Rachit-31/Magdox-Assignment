import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { afterEach, expect, it, vi } from "vitest";
import { App } from "./App";

const json = (body: unknown, status = 200) => new Response(JSON.stringify(body), { status });

// A stand-in for the API: two pages of products and one finding.
async function fakeApi(input: RequestInfo | URL) {
  const url = new URL(String(input), "http://localhost");
  if (url.pathname === "/api/session") {
    const user = { name: "Asha", email: "asha@kestrel.example", role: "analyst", org: "Kestrel Cameras" };
    return json({ user, csrfToken: "made-up" });
  }
  if (url.pathname === "/api/products") {
    const page = Number(url.searchParams.get("page"));
    const item = { id: `7c1e5a90-2d4f-4b68-9a13-5e8f0c6d2b0${page}`, name: `Dome camera ${page}`, sku: `K-${page}`, deviceType: "camera" };
    return json({ items: [item], page, pageSize: 5, hasNext: page < 2 });
  }
  if (url.pathname === "/api/findings") {
    const finding = { id: "e3a9d7c4-1f52-4e86-b0a7-4c2d9f8e6a55", title: "Telnet service enabled by default", severity: "critical", component: "", product: "Kestrel K200 Dome" };
    return json({ items: [finding] });
  }
  return json({ error: "not found" }, 404);
}

afterEach(() => vi.unstubAllGlobals());

it("signs in, pages through products, searches and shows findings", async () => {
  const fetchMock = vi.fn(fakeApi);
  vi.stubGlobal("fetch", fetchMock);
  const user = userEvent.setup();
  render(<App />);

  expect(screen.getByText(/practice app with made-up users/)).toBeInTheDocument();
  expect(screen.getAllByRole("button", { name: /^Sign in as/ })).toHaveLength(3);
  await user.click(screen.getByRole("button", { name: "Sign in as Asha" }));

  expect(await screen.findByText("Dome camera 1")).toBeInTheDocument();
  expect(await screen.findByRole("table", { name: "Findings" })).toBeInTheDocument();
  expect(screen.getByRole("button", { name: "Previous" })).toBeDisabled();

  await user.click(screen.getByRole("button", { name: "Next" }));
  expect(await screen.findByText("Dome camera 2")).toBeInTheDocument();
  expect(screen.getByRole("button", { name: "Next" })).toBeDisabled();
  expect(screen.getByRole("button", { name: "Previous" })).toBeEnabled();

  // A new search starts again from the first page.
  await user.type(screen.getByLabelText("Search products"), "d");
  expect(await screen.findByText("Dome camera 1")).toBeInTheDocument();
  expect(fetchMock).toHaveBeenLastCalledWith("/api/products?q=d&page=1", { credentials: "same-origin" });
});

it("shows the error from the API when sign-in fails", async () => {
  vi.stubGlobal("fetch", vi.fn(async () => json({ error: "unknown user" }, 400)));
  const user = userEvent.setup();
  render(<App />);
  await user.click(screen.getByRole("button", { name: "Sign in as Ravi" }));
  expect(await screen.findByRole("alert")).toHaveTextContent("unknown user");
});
