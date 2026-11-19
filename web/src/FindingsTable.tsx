import { useState } from "react";
import type { Finding } from "./api";

export function FindingsTable({ findings }: { findings: Finding[] }) {
  const [severity, setSeverity] = useState<"all" | Finding["severity"]>("all");
  const visibleFindings = severity === "all" ? findings : findings.filter((finding) => finding.severity === severity);

  return (
    <div className="table-scroll">
      <label htmlFor="finding-severity">Severity</label>
      <select
        id="finding-severity"
        value={severity}
        onChange={(event) => setSeverity(event.target.value as "all" | Finding["severity"])}
      >
        <option value="all">All</option>
        <option value="critical">Critical</option>
        <option value="high">High</option>
        <option value="medium">Medium</option>
        <option value="low">Low</option>
      </select>
      <p role="status">Showing {visibleFindings.length} of {findings.length} findings</p>
      <table>
        <caption>Findings</caption>
        <thead>
          <tr>
            <th scope="col">Title</th>
            <th scope="col">Severity</th>
            <th scope="col">Component</th>
            <th scope="col">Product</th>
          </tr>
        </thead>
        <tbody>
          {visibleFindings.length === 0 ? (
            <tr>
              <td colSpan={4}>No findings match.</td>
            </tr>
          ) : (
            visibleFindings.map((f) => (
              <tr key={f.id}>
                <td>{f.title}</td>
                <td>
                  <span className={`pill pill-${f.severity}`}>{f.severity}</span>
                </td>
                <td>{f.component || "Not identified"}</td>
                <td>{f.product}</td>
              </tr>
            ))
          )}
        </tbody>
      </table>
    </div>
  );
}
