import {
  Table,
  TableBody,
  TableCaption,
  TableCell,
  TableFooter,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import { useEffect, useState } from "react";
import { Link } from "react-router-dom";

import { ChevronLeft, ChevronRight } from "lucide-react";

type Submission = {
  Id: string;
  Title: string;
  Username: string;
  Link: string;
  Body: string;
  Created_at: string;
};

function truncate(
  str: string,
  maxLength: number = 40,
  suffix: string = "..."
): string {
  if (str.length <= maxLength) {
    return str;
  }
  return str.slice(0, maxLength - suffix.length) + suffix;
}

function buildNextFetch(filter: string, offset: number): string {
  return `/api/v1/all?sort=${filter}&offset=${offset}`;
}

export default function DataTable() {
  const [submission, setSubmission] = useState<Submission[]>([]);

  // use both of these states to build the URL to fetch from: /api/v1/all?sort=<FILTER>&offset=<OFFSET>
  const [filter, setFilter] = useState<string>("latest");
  const [offset, setOffset] = useState<number>(0); // Start at 0

  const [isPending, setIsPending] = useState<boolean>(true);
  const [isError, setIsError] = useState<boolean>(false);

  useEffect(() => {
    setIsPending(true);
    setIsError(false);
    
    fetch("http://localhost:3000" + buildNextFetch(filter, offset))
      .then((res) => {
        if (!res.ok) {
          throw new Error("Network response was not ok");
        }
        return res.json();
      })
      .then((json) => {
        const res: Submission[] = json.results;
        setSubmission(res);
        setIsPending(false);
      })
      .catch((err: Error) => {
        console.error(err);
        setIsError(true);
        setIsPending(false);
      });
  }, [offset, filter]); // Also depend on filter in case you add filter controls later

  return (
    <div>
      {!isPending && !isError && (
        <Table>
          <TableHeader>
            <p className="text-lg">HackerNews</p>
            <p className="text-base text-muted-foreground">
              {new Date().toLocaleString("en-US", {
                month: "2-digit",
                day: "2-digit",
                year: "numeric",
                hour: "2-digit",
                minute: "2-digit",
                hour12: true,
              })}
            </p>
          </TableHeader>
          <TableBody>
            {submission.map((s) => (
              <TableRow key={s.Id}>
                <TableCell className="font-medium">
                  <a href={s.Link} target="_blank">
                    {truncate(s.Title)}
                  </a>
                </TableCell>
                <TableCell>
                  <Link to={"/submission/" + s.Id}>
                    {" "}
                    {truncate(s.Body, 60)}
                  </Link>
                </TableCell>
                <TableCell>
                  <Link to={"/u/" + s.Username}>{s.Username}</Link>
                </TableCell>
                <TableCell className="text-right">{s.Created_at}</TableCell>
              </TableRow>
            ))}
          </TableBody>
          <TableFooter>
            <TableRow key={"end"}>
              <TableCell>
                <Button
                  variant="outline"
                  disabled={offset === 0}
                  onClick={() => {
                    setOffset((prev) => Math.max(0, prev - 10));
                  }}
                >
                  <ChevronLeft />
                  Prev
                </Button>
                <Button
                  variant="outline"
                  disabled={submission.length < 10}
                  onClick={() => {
                    setOffset((prev) => prev + 10);
                  }}
                >
                  Next
                  <ChevronRight />
                </Button>
              </TableCell>
            </TableRow>
          </TableFooter>
        </Table>
      )}
      {isPending && <div>Loading...</div>}
      {isError && <div>Error loading data</div>}
    </div>
  );
}