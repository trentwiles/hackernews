import WebSidebar from "./fragments/WebSidebar";
import { SidebarInset, SidebarProvider } from "./ui/sidebar";

export default function Submit() {
  return (
    <SidebarProvider>
      <div className="flex min-h-screen w-full">
        <WebSidebar />
        <SidebarInset>
          <span>404</span>
        </SidebarInset>
      </div>
    </SidebarProvider>
  );
}
