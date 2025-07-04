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

function findOffset(url: string): number | null {
  const params = new URLSearchParams(url.split("?")[1]);
  const offset = params.get("offset");
  return offset !== null ? Number(offset) : null;
}

export default function DataTable() {
  const [submission, setSubmission] = useState<Submission[]>([]);
  const [next, setNext] = useState<string>("/api/v1/all?sort=latest");
  const [offset, setOffset] = useState<number>(0);
  const [isPending, setIsPending] = useState<boolean>(true);
  const [isError, setIsError] = useState<boolean>(false);

  useEffect(() => {
    fetch("http://localhost:3000" + next)
      .then((res) => {
        if (!res.ok) {
          throw new Error("Network response was not ok");
        }
        return res.json();
      })
      .then((json) => {
        const res: Submission[] = json.results;
        setNext(json.next);
        let offset: number | null = findOffset(json.next);
        if (offset == null) {
          offset = 0;
        } else {
          offset = offset - 10;
        }
        setOffset(offset);
        setSubmission(res);
        setIsPending(false);
      })
      .catch((err: Error) => {
        console.error(err);
        setIsError(true);
        setIsPending(false);
      });
  }, []);

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
            {!isPending && !isError && (
              <>
                <Button variant="outline" disabled={offset == 0}>
                  <ChevronLeft />
                  Prev
                </Button>
                <Button variant="outline" disabled={submission.length != 10}>
                  Next
                  <ChevronRight />
                </Button>
              </>
            )}
          </TableFooter>
        </Table>
      )}
    </div>
    // future: next button
  );
}
