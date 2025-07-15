import WebSidebar from "./fragments/WebSidebar";
import { SidebarProvider, SidebarInset } from "@/components/ui/sidebar";

export default function Forbidden() {

  return (
    <SidebarProvider>
      <div className="flex min-h-screen w-full">
        <WebSidebar />
        <SidebarInset>
          403 Forbidden
        </SidebarInset>
      </div>
    </SidebarProvider>
  );
}
