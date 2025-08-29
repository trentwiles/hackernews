import WebSidebar, { SidebarBreadcrumbHeader } from "./fragments/WebSidebar";
import DataTable from "./fragments/DataTable";
import { SidebarProvider, SidebarInset } from "@/components/ui/sidebar";

type props = {
  sortType?: string;
};
const breadcrumbs = [{ label: "Home", isCurrentPage: true }];
export default function Home(props: props) {
  return (
    <SidebarProvider>
      <WebSidebar />
      <SidebarInset>
        <SidebarBreadcrumbHeader breadcrumbs={breadcrumbs} />
        <div className="flex flex-1 flex-col gap-4 p-4">
          <DataTable key={props.sortType} {...props} />
        </div>
      </SidebarInset>
    </SidebarProvider>
  );
}
