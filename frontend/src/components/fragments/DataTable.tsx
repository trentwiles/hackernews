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
import { useEffect, useState } from "react";

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

export default function DataTable() {
  const [submission, setSubmission] = useState<Submission[]>([]);
  const [next, setNext] = useState<string>("");
  const [isPending, setIsPending] = useState<boolean>(true);
  const [isError, setIsError] = useState<boolean>(false);

  useEffect(() => {
    fetch("http://localhost:3000/api/v1/all?sort=latest")
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
  }, []);

  return (
    <div>
      {!isPending && !isError && (
        <Table>
          <TableBody>
            {submission.map((s) => (
              <TableRow key={s.Id}>
                <TableCell className="font-medium">
                  <a href={s.Link} target="_blank">
                    {truncate(s.Title)}
                  </a>
                </TableCell>
                <TableCell>
                  <a href={"/submission/" + s.Id}> {truncate(s.Body, 60)}</a>
                </TableCell>
                <TableCell><a href={"/u/" + s.Username}>{s.Username}</a></TableCell>
                <TableCell className="text-right">{s.Created_at}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      )}
    </div>
    // future: next button
  );
}
