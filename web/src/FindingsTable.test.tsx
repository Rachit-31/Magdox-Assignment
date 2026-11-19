import { render, screen, within } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import type { Finding } from "./api";
import { FindingsTable } from "./FindingsTable";

const findings: Finding[] = [
  {
    id: "0b9f6a1e-3c47-4d1a-9e52-6f0a7c2d8b11",
    title: "Telnet service enabled by default",
    severity: "critical",
    component: "netd 2.1.0",
    product: "Kestrel K200 Dome",
  },
  {
    id: "5d2c8e70-91ab-4f36-8c0d-2e7b4a6f9c22",
    title: "Update package accepted without a signature check",
    severity: "high",
    component: "",
    product: "Kestrel K410 Bullet",
  },
  {
    id: "a47e1f93-6b28-4c05-b7d9-8c3e5f1a0d33",
    title: "Web console session never expires",
    severity: "medium",
    component: "webui 3.2.1",
    product: "Kestrel K200 Dome",
  },
  {
    id: "c81d4b26-f05e-4a79-a3c1-9d6e2b7f4e44",
    title: "Hard-coded maintenance account",
    severity: "high",
    component: "shelld 0.9.1",
    product: "Kestrel K410 Bullet",
  },
];

describe("FindingsTable", () => {
  it("renders a table named Findings with four columns", () => {
    render(<FindingsTable findings={findings} />);
    const table = screen.getByRole("table", { name: "Findings" });
    const headers = within(table)
      .getAllByRole("columnheader")
      .map((th) => th.textContent);
    expect(headers).toEqual(["Title", "Severity", "Component", "Product"]);
  });

  it("renders a row for every finding", () => {
    render(<FindingsTable findings={findings} />);
    for (const finding of findings) {
      expect(screen.getByText(finding.title)).toBeInTheDocument();
    }
    const row = screen.getByRole("row", { name: /^Telnet service enabled by default/ });
    expect(within(row).getByText("critical")).toBeInTheDocument();
    expect(within(row).getByText("netd 2.1.0")).toBeInTheDocument();
    expect(within(row).getByText("Kestrel K200 Dome")).toBeInTheDocument();
  });

  it("shows Not identified when the component is empty", () => {
    render(<FindingsTable findings={findings} />);
    const row = screen.getByRole("row", { name: /^Update package accepted without a signature check/ });
    expect(within(row).getByText("Not identified")).toBeInTheDocument();
  });
});
