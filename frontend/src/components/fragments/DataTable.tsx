import {
  Table,
  TableBody,
  TableCell,
  TableFooter,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { useEffect, useState } from "react";
import { Link } from "react-router-dom";

import { ChevronLeft, ChevronRight } from "lucide-react";
import { datePrettyPrint, getTimeAgo, truncate } from "@/utils";
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@radix-ui/react-tooltip";

type Submission = {
  Id: string;
  Title: string;
  Username: string;
  Link: string;
  Body: string;
  Created_at: string;
};


function buildNextFetch(filter: string, offset: number): string {
  return `/api/v1/all?sort=${filter}&offset=${offset}`;
}

type props = {
  sortType?: string;
};

export default function DataTable(props: props) {
  const [submission, setSubmission] = useState<Submission[]>([]);

  let tempFilter: string = "latest";
  if (props !== undefined && props.sortType !== undefined) {
    tempFilter = props.sortType;
  }

  const [filter /*, setFilter */] = useState<string>(tempFilter);
  const [offset, setOffset] = useState<number>(0);

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
  }, [offset, filter]);

  return (
    <div className="w-full p-6 max-w-6xl mx-auto">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead colSpan={4} className="h-auto">
              <div className="py-2">
                <p className="text-lg font-semibold">HackerNews</p>
                <p className="text-sm text-muted-foreground">
                  {new Date().toLocaleString("en-US", {
                    month: "2-digit",
                    day: "2-digit",
                    year: "numeric",
                    hour: "2-digit",
                    minute: "2-digit",
                    hour12: true,
                  })}
                </p>
              </div>
            </TableHead>
          </TableRow>
          <TableRow>
            <TableHead className="w-[30%]">Title</TableHead>
            <TableHead className="w-[40%]">Preview</TableHead>
            <TableHead className="w-[15%]">User</TableHead>
            <TableHead className="w-[15%] text-right">Date</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {isPending ? (
            Array.from({ length: 10 }).map((_, index) => (
              <TableRow key={`skeleton-${index}`}>
                <TableCell className="font-medium py-4">
                  <Skeleton className="h-5 w-[250px]" />
                </TableCell>
                <TableCell className="py-4">
                  <Skeleton className="h-5 w-[350px]" />
                </TableCell>
                <TableCell className="py-4">
                  <Skeleton className="h-5 w-[100px]" />
                </TableCell>
                <TableCell className="text-right py-4">
                  <Skeleton className="h-5 w-[120px] ml-auto" />
                </TableCell>
              </TableRow>
            ))
          ) : isError ? (
            <TableRow>
              <TableCell
                colSpan={4}
                className="text-center text-muted-foreground"
              >
                Error loading data. Please try again.
              </TableCell>
            </TableRow>
          ) : submission.length === 0 ? (
            <TableRow>
              <TableCell
                colSpan={4}
                className="text-center text-muted-foreground"
              >
                No submissions found.
              </TableCell>
            </TableRow>
          ) : (
            submission.map((s) => (
              <TableRow key={s.Id}>
                <TableCell className="font-medium py-4">
                  <a href={s.Link} target="_blank" className="hover:underline">
                    {truncate(s.Title)}
                  </a>
                </TableCell>
                <TableCell className="py-4">
                  <Link
                    to={"/submission/" + s.Id}
                    className="text-muted-foreground hover:text-foreground"
                  >
                    {truncate(s.Body, 60)}
                    {s.Body.length == 0 && "(no preview)"}
                  </Link>
                </TableCell>
                <TableCell className="py-4">
                  <Link to={"/u/" + s.Username} className="hover:underline">
                    {s.Username}
                  </Link>
                </TableCell>
                <TableCell className="text-right py-4">
                  <Tooltip>
                    <TooltipTrigger asChild>
                      <p>{getTimeAgo(s.Created_at)}</p>
                    </TooltipTrigger>
                    <TooltipContent>
                      <span>{datePrettyPrint(s.Created_at)}</span>
                    </TooltipContent>
                  </Tooltip>
                </TableCell>
              </TableRow>
            ))
          )}
        </TableBody>
        <TableFooter>
          <TableRow key={"end"}>
            <TableCell colSpan={4}>
              <div className="flex gap-2">
                <Button
                  variant="outline"
                  disabled={offset === 0 || isPending}
                  onClick={() => {
                    setOffset((prev) => Math.max(0, prev - 10));
                  }}
                >
                  <ChevronLeft />
                  Prev
                </Button>
                <Button
                  variant="outline"
                  disabled={submission.length < 10 || isPending}
                  onClick={() => {
                    setOffset((prev) => prev + 10);
                  }}
                >
                  Next
                  <ChevronRight />
                </Button>
              </div>
            </TableCell>
          </TableRow>
        </TableFooter>
      </Table>
    </div>
  );
}
