import WebSidebar from "./fragments/WebSidebar";
import DataTable from "./fragments/DataTable";
import { SidebarProvider, SidebarInset } from "@/components/ui/sidebar";

export default function Home() {
  return (
    <SidebarProvider>
      <div className="flex min-h-screen w-full">
        <WebSidebar />
        <SidebarInset>
          <DataTable />
        </SidebarInset>
      </div>
    </SidebarProvider>
  );
}