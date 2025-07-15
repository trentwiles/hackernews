import { getDomain, truncate } from "@/utils";
import WebSidebar from "./fragments/WebSidebar";
import { SidebarProvider, SidebarInset } from "@/components/ui/sidebar";
import { Search as SearchIcon } from "lucide-react";
import { useEffect, useState } from "react";
import { Link } from "react-router-dom";

type submission = {
  id: string;
  username: string;
  title: string;
  link: string;
  body: string;
  created_at: string;
};

export default function Search() {
  const [query, setQuery] = useState<string>("");
  const [results, setResults] = useState<submission[]>([]);

  useEffect(() => {
    const params = new URLSearchParams(window.location.search);
    const q = params.get("q");
    if (q !== null && q != "") {
      setQuery(q)
      document.title = `Search for '${q}' | ${import.meta.env.VITE_SERVICE_NAME}`
    }
  }, [])

  useEffect(() => {
    if (query == "") {
      console.log("Note: no query in search box. Exiting useEffect");
      window.history.replaceState({}, '', '/search');
      document.title = `Search | ${import.meta.env.VITE_SERVICE_NAME}`
      return;
    }

    window.history.replaceState({}, '', '?q=' + query);
    document.title = `Search for '${query}' | ${import.meta.env.VITE_SERVICE_NAME}`

    fetch(import.meta.env.VITE_API_ENDPOINT + "/api/v1/searchSubmissions?q=" + query)
      .then((response) => {
        if (response.status != 200) {
          throw new Error("non-200 HTTP status");
        }
        return response.json();
      })
      .then((data) => {
        if (data.results == null) {
          setResults([]);
          return;
        }

        const ls: submission[] = [];
        data.results.map((res) => {
          console.log(res);
          const tempResult: submission = {
            id: res.Id,
            body: res.Body,
            created_at: res.Created_at,
            link: res.Link,
            title: res.Title,
            username: res.Username,
          };

          ls.push(tempResult);
        });

        setResults(ls);
      })
      .catch((error) => {
        console.error(error);
      });
  }, [query]);

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
                value={query}
                onChange={(e) => setQuery(e.target.value)}
              />
            </div>

            {/* Sample Results */}
            <div className="space-y-6">
              <h2 className="text-lg font-semibold text-gray-700 mb-4">
                Search Results
              </h2>
              <p className="text-sm text-muted-foreground mt-5">
                {results.length !== 0 &&
                  `${results.length} result${
                    results.length === 1 ? "" : "s found"
                  }`}
                {results.length === 0 &&
                  query !== "" &&
                  `No results found for query '${query}'`}
              </p>

              {/* Result list */}
              {results.length !== 0 &&
                results.map((res, idx) => (
                  <div className="border rounded-lg p-4 hover:bg-accent transition-colors cursor-pointer" key={idx}>
                    <div className="flex justify-between items-start">
                      <div className="space-y-1">
                        <h3 className="font-semibold text-lg hover:underline">
                          <Link
                            to={"/submission/" + res.id}
                          >
                            {res.title}
                          </Link>
                        </h3>
                        <p className="text-sm text-muted-foreground mb-2">
                          {truncate(res.body)}
                        </p>
                        <div className="flex gap-3 text-sm text-muted-foreground">
                          <span className="inline-flex items-center px-2 py-1 rounded-full text-xs bg-secondary">
                            <a
                              href={res.link}
                              target="_blank"
                            >
                              {getDomain(res.link)}
                            </a>
                          </span>
                          <span className="inline-flex items-center px-2 py-1 rounded-full text-xs bg-secondary">
                            <Link to={`/u/${res.username}`}>u/{res.username}</Link>
                          </span>
                          <span>2 days ago</span>
                        </div>
                      </div>
                    </div>
                  </div>
                ))}
            </div>
          </div>
        </SidebarInset>
      </div>
    </SidebarProvider>
  );
}
