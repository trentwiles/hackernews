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
import { TrendingUp, FileText, Users } from "lucide-react";
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
  const [activeUsers, setActiveUsers] = useState<number>(-1);
  const [totalUsers, setTotalUsers] = useState<number>(-1);
  const [isPending, setIsPending] = useState<boolean>(true);
  const [isError, setIsError] = useState<boolean>(false);

  useEffect(() => {
    fetch(import.meta.env.VITE_API_ENDPOINT + "/api/v1/adminMetrics")
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
        setActiveUsers(json.metrics.TotalActiveUsers);
        setTotalUsers(json.metrics.TotalAllTimeUsers);
        setIsPending(false);
        setIsError(false);
      })
      .catch((err) => {
        console.error(err);
        setIsError(true);
        setIsPending(false);
      });
  }, []);

  const chartConfig = {
    posts: {
      label: "Total Posts",
      color: "hsl(var(--chart-1))",
    },
  } satisfies ChartConfig;

  const totalWeekPosts = chartData.reduce((sum, day) => sum + day.posts, 0);

  return (
    <SidebarProvider>
      <div className="flex min-h-screen w-full">
        <WebSidebar />
        <SidebarInset>
          <div className="p-6">
            <h1 className="text-3xl font-bold mb-6">Admin Dashboard</h1>

            {!isPending && !isError && (
              <div className="grid grid-cols-2 gap-6">
                {/* Chart - Left Half */}
                <Card className="col-span-1">
                  <CardHeader>
                    <CardTitle>Post Activity</CardTitle>
                    <CardDescription>
                      Total posts over the last 7 days
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
                          fill="var(--color-posts)"
                          fillOpacity={0.4}
                          stroke="var(--color-posts)"
                        />
                      </AreaChart>
                    </ChartContainer>
                  </CardContent>
                  <CardFooter>
                    <div className="flex w-full items-start gap-2 text-sm">
                      <div className="grid gap-2">
                        <div className="flex items-center gap-2 leading-none font-medium">
                          {totalWeekPosts} total
                        </div>
                        <div className="text-muted-foreground flex items-center gap-2 leading-none">
                          Last 7 days
                        </div>
                      </div>
                    </div>
                  </CardFooter>
                </Card>

                {/* Metrics - Right Half */}
                <div className="col-span-1 flex flex-col gap-6">
                  {/* Metric Card 1 - Today's Posts */}
                  <Card className="flex-1">
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                      <CardTitle className="text-sm font-large">
                        Active Users
                      </CardTitle>
                      <FileText className="h-4 w-4 text-muted-foreground" />
                    </CardHeader>
                    <CardContent>
                      <div className="text-2xl font-bold">
                        {activeUsers || 0}
                      </div>
                      <p className="text-xs text-muted-foreground mt-1">
                        in the last 7 days,{" "}
                        {(
                          (activeUsers / (activeUsers + totalUsers)) *
                          100
                        ).toFixed(2)}
                        % of {totalUsers} total users
                      </p>
                    </CardContent>
                  </Card>

                  {/* Metric Card 2 - Weekly Total */}
                  <Card className="flex-1">
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                      <CardTitle className="text-sm font-large">
                        Weekly Total
                      </CardTitle>
                      <Users className="h-4 w-4 text-muted-foreground" />
                    </CardHeader>
                    <CardContent>
                      <div className="text-2xl font-bold">{totalWeekPosts}</div>
                      <p className="text-xs text-muted-foreground mt-1">
                        Average: {(totalWeekPosts / 7).toFixed(1)} posts/day
                      </p>
                    </CardContent>
                  </Card>
                </div>
              </div>
            )}

            {isError && (
              <Card>
                <CardHeader>
                  <CardTitle>Error</CardTitle>
                  <CardDescription>
                    Failed to load metrics. Please try again later.
                  </CardDescription>
                </CardHeader>
              </Card>
            )}
          </div>
        </SidebarInset>
      </div>
    </SidebarProvider>
  );
}
