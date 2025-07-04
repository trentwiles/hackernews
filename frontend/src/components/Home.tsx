import WebSidebar from "./fragments/WebSidebar";
import DataTable from "./fragments/DataTable";
export default function Home() {
  return (
    <div className="flex min-h-screen">
      <WebSidebar />
      <DataTable />
    </div>
  );
}
