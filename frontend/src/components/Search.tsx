import WebSidebar from "./fragments/WebSidebar";
import { SidebarProvider, SidebarInset } from "@/components/ui/sidebar";
import { Search as SearchIcon } from "lucide-react";
import { Link } from "react-router-dom";

export default function Search() {
  return (
    <SidebarProvider>
      <div className="flex min-h-screen w-full">
        <WebSidebar />
        <SidebarInset>
          {/* Add a container div with padding */}
          <div className="px-6 py-6 md:px-8 lg:px-12">
            {/* Search Bar */}
            <div className="relative mb-8">
              <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <SearchIcon className="h-5 w-5 text-gray-400" />
              </div>
              <input
                type="text"
                className="block w-full pl-10 pr-3 py-3 border border-gray-300 rounded-lg bg-white text-gray-900 placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="Search for anything..."
                defaultValue="React components"
              />
            </div>

            {/* Sample Results */}
            <div className="space-y-6">
              <h2 className="text-lg font-semibold text-gray-700 mb-4">
                Search Results
              </h2>

              {/* Result 1, clone this in a .map or .forEach */}
              <div className="border rounded-lg p-4 hover:bg-accent transition-colors cursor-pointer">
                <div className="flex justify-between items-start">
                  <div className="space-y-1">
                    <h3 className="font-semibold text-lg hover:underline">
                      <a
                        href="https://react.dev/learn/your-first-component"
                        target="_blank"
                      >
                        Getting Started with React Components
                      </a>
                    </h3>
                    <p className="text-sm text-muted-foreground mb-2">
                      Components are one of the core concepts of React. They are
                      the foundation upon which you build user interfaces (UI),
                      which makes them the perfect place to start your React
                      journey...
                    </p>
                    <div className="flex gap-3 text-sm text-muted-foreground">
                      <span className="inline-flex items-center px-2 py-1 rounded-full text-xs bg-secondary">
                        <a
                          href="https://react.dev/learn/your-first-component"
                          target="_blank"
                        >
                          react.dev
                        </a>
                      </span>
                      <span className="inline-flex items-center px-2 py-1 rounded-full text-xs bg-secondary">
                        <Link to="/u/trent">u/trent</Link>
                      </span>
                      <span>2 days ago</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </SidebarInset>
      </div>
    </SidebarProvider>
  );
}
