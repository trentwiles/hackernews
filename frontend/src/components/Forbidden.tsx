import { SidebarInset, SidebarProvider } from "@/components/ui/sidebar";
import WebSidebar, { SidebarBreadcrumbHeader } from "./fragments/WebSidebar";

export default function Forbidden() {
  const breadcrumbs = [
    { label: "Building Your Application", href: "/docs" },
    { label: "Data Fetching", isCurrentPage: true }
  ];

  return (
    <SidebarProvider>
      <WebSidebar />
      <SidebarInset>
        <SidebarBreadcrumbHeader breadcrumbs={breadcrumbs} />
        <div className="flex flex-1 flex-col gap-4 p-4">
          {/* START CUSTOM PAGE CONTENT */}
          <p>"403 Forbidden"</p>
        </div>
      </SidebarInset>
    </SidebarProvider>
  );
}