import WebSidebar from "./fragments/WebSidebar";
import DataTable from "./fragments/DataTable";
export default function Home() {
  return (
    <div className="flex">
      <WebSidebar />
      <DataTable />
    </div>
  );
}
