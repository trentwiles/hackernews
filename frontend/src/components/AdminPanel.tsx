import WebSidebar from "./fragments/WebSidebar";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  type ChartConfig,
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
} from "@/components/ui/chart";
import { Area, AreaChart, CartesianGrid, XAxis } from "recharts";
import { SidebarInset, SidebarProvider } from "./ui/sidebar";
import { TrendingUp } from "lucide-react";
import { useEffect, useState } from "react";

export default function AdminPanel() {
  const daysOfTheWeek = [
    "Sunday",
    "Monday",
    "Tuesday",
    "Wednesday",
    "Thursday",
    "Friday",
    "Saturday",
  ];
  const [chartData, setChartData] = useState<{ day: string; posts: number }[]>(
    []
  );
  const [isPending, setIsPending] = useState<boolean>(true);
  const [isError, setIsError] = useState<boolean>(false);

  useEffect(() => {
    fetch(import.meta.env.VITE_API_ENDPOINT + "/api/v1/adminMetrics") // example API
      .then((response) => {
        if (!response.ok) throw new Error("Network response was not ok");
        return response.json();
      })
      .then((json) => {
        const data = [
          {
            day: daysOfTheWeek[((new Date().getDay() % 7) + 7) % 7],
            posts: json.metrics.TodayPosts,
          },
          {
            day: daysOfTheWeek[(((new Date().getDay() - 1) % 7) + 7) % 7],
            posts: json.metrics.TodayMinusOnePosts,
          },
          {
            day: daysOfTheWeek[(((new Date().getDay() - 2) % 7) + 7) % 7],
            posts: json.metrics.TodayMinusTwoPosts,
          },
          {
            day: daysOfTheWeek[(((new Date().getDay() - 3) % 7) + 7) % 7],
            posts: json.metrics.TodayMinusThreePosts,
          },
          {
            day: daysOfTheWeek[(((new Date().getDay() - 4) % 7) + 7) % 7],
            posts: json.metrics.TodayMinusFourPosts,
          },
          {
            day: daysOfTheWeek[(((new Date().getDay() - 5) % 7) + 7) % 7],
            posts: json.metrics.TodayMinusFivePosts,
          },
          {
            day: daysOfTheWeek[(((new Date().getDay() - 6) % 7) + 7) % 7],
            posts: json.metrics.TodayMinusSixPosts,
          },
        ].reverse();

        setChartData(data);
        setIsPending(false);
        setIsError(false);
      })
      .catch((err) => {
        console.error(err);
        setIsError(true);
      });
  }, []);
  // daysOfTheWeek[(OFFSET % 7 + 7) % 7]
  // starts at sunday because javascript starts the weekday on sunday (wtf???)

  const chartConfig = {
    posts: {
      label: "Total Posts",
      color: "hsl(var(--chart-1))",
    },
  } satisfies ChartConfig;

  return (
    <SidebarProvider>
      <div className="flex min-h-screen w-full">
        <WebSidebar />
        <SidebarInset>
          {!isPending && (
            <Card>
              <CardHeader>
                <CardTitle>Area Chart</CardTitle>
                <CardDescription>
                  Total Posts over the last 7 days
                </CardDescription>
              </CardHeader>
              <CardContent>
                <ChartContainer config={chartConfig}>
                  <AreaChart
                    accessibilityLayer
                    data={chartData}
                    margin={{
                      left: 12,
                      right: 12,
                    }}
                  >
                    <CartesianGrid vertical={false} />
                    <XAxis
                      dataKey="day"
                      tickLine={false}
                      axisLine={false}
                      tickMargin={8}
                      tickFormatter={(value) => value.slice(0, 3)}
                    />
                    <ChartTooltip
                      cursor={false}
                      content={<ChartTooltipContent indicator="line" />}
                    />
                    <Area
                      dataKey="posts"
                      type="linear"
                      fill="var(--color-desktop)"
                      fillOpacity={0.4}
                      stroke="var(--color-desktop)"
                    />
                  </AreaChart>
                </ChartContainer>
              </CardContent>
              <CardFooter>
                <div className="flex w-full items-start gap-2 text-sm">
                  <div className="grid gap-2">
                    <div className="flex items-center gap-2 leading-none font-medium">
                      Trending up by 5.2% this month{" "}
                      <TrendingUp className="h-4 w-4" />
                    </div>
                    <div className="text-muted-foreground flex items-center gap-2 leading-none">
                      January - June 2024
                    </div>
                  </div>
                </div>
              </CardFooter>
            </Card>
          )}
        </SidebarInset>
      </div>
    </SidebarProvider>
  );
}
