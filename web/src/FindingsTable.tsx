import type { Finding } from "./api";

export function FindingsTable({ findings }: { findings: Finding[] }) {
  return (
    <div className="table-scroll">
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
          {findings.map((f) => (
            <tr key={f.id}>
              <td>{f.title}</td>
              <td>
                <span className={`pill pill-${f.severity}`}>{f.severity}</span>
              </td>
              <td>{f.component || "Not identified"}</td>
              <td>{f.product}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
