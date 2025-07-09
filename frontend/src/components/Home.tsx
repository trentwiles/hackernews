import WebSidebar from "./fragments/WebSidebar";
import DataTable from "./fragments/DataTable";
import { SidebarProvider, SidebarInset } from "@/components/ui/sidebar";

type props = {
  sortType?: string;
};

export default function Home(props: props) {

  return (
    <SidebarProvider>
      <div className="flex min-h-screen w-full">
        <WebSidebar />
        <SidebarInset>
          <DataTable key={props.sortType} {...props} />
        </SidebarInset>
      </div>
    </SidebarProvider>
  );
}
