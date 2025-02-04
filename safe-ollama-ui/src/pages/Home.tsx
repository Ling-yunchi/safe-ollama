import { useState, useEffect } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  ChartConfig,
  ChartContainer,
  ChartLegend,
  ChartLegendContent,
  ChartTooltip,
  ChartTooltipContent,
} from "@/components/ui/chart";
import { CartesianGrid, XAxis, Area, AreaChart } from "recharts";
import { useAtomValue } from "jotai";
import { userToken } from "@/storage/user";
import { LoaderCircle } from "lucide-react";

const chartConfig = {
  promptTokens: {
    label: "Prompt",
    color: "hsl(var(--chart-1))",
  },
  responseTokens: {
    label: "Response",
    color: "hsl(var(--chart-2))",
  },
} satisfies ChartConfig;

interface TokenUsage {
  date: string;
  promptTokens: number;
  responseTokens: number;
}

function Home() {
  const [loading, setLoading] = useState(true);
  const [tokenUsage, setTokenUsage] = useState<TokenUsage[]>([]);
  const token = useAtomValue(userToken);

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    setLoading(true);
    const _endDate = new Date();
    const endDateStr = _endDate.toLocaleDateString("en-CA");
    // one week
    const _startDate = new Date(_endDate);
    _startDate.setDate(_endDate.getDate() - 6);
    const startDateStr = new Date(_startDate).toLocaleDateString("en-CA");
    console.log(startDateStr, endDateStr);

    const response = await fetch(
      `/api/token_usage/user/usage/daily?start=${startDateStr}&end=${endDateStr}`,
      {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
          Token: token,
        },
      }
    );

    if (!response.ok) {
      // toast.error("Failed to fetch tokens data");
      setLoading(false);
      return;
    }
    const data = (await response.json()) as TokenUsage[];
    console.log(data);

    const startDate = new Date(startDateStr);
    const endDate = new Date(endDateStr);
    const oneDay = 24 * 60 * 60 * 1000; // 一天的毫秒数
    const dateRange: TokenUsage[] = [];
    let currentDate = startDate;

    while (currentDate <= endDate) {
      const currentDateString = currentDate.toLocaleDateString("en-CA");
      const foundData = data.find((d) => d.date === currentDateString);

      dateRange.push({
        date: currentDateString,
        promptTokens: foundData ? foundData.promptTokens : 0,
        responseTokens: foundData ? foundData.responseTokens : 0,
      });
      currentDate = new Date(currentDate.getTime() + oneDay);
    }
    console.log(dateRange);

    setTokenUsage(dateRange);
    setLoading(false);
  };

  return (
    <div className="p-6 w-full">
      <Card className="w-full">
        <CardHeader>
          <CardTitle>Token Usage</CardTitle>
        </CardHeader>
        <CardContent className="px-2 pt-4 sm:px-6 sm:pt-6">
          {loading ? (
            <div className="flex justify-center items-center h-[250px]">
              <div className="flex items-center gap-2 text-md">
                <LoaderCircle className="h-4 w-4 animate-spin" />
                Loading ...
              </div>
            </div>
          ) : (
            <ChartContainer
              config={chartConfig}
              className="aspect-auto h-[250px] w-full"
            >
              <AreaChart data={tokenUsage}>
                <defs>
                  <linearGradient
                    id="fillPromptTokens"
                    x1="0"
                    y1="0"
                    x2="0"
                    y2="1"
                  >
                    <stop
                      offset="5%"
                      stopColor="var(--color-promptTokens)"
                      stopOpacity={0.8}
                    />
                    <stop
                      offset="95%"
                      stopColor="var(--color-promptTokens)"
                      stopOpacity={0.1}
                    />
                  </linearGradient>
                  <linearGradient
                    id="fillResponseTokens"
                    x1="0"
                    y1="0"
                    x2="0"
                    y2="1"
                  >
                    <stop
                      offset="5%"
                      stopColor="var(--color-responseTokens)"
                      stopOpacity={0.8}
                    />
                    <stop
                      offset="95%"
                      stopColor="var(--color-responseTokens)"
                      stopOpacity={0.1}
                    />
                  </linearGradient>
                </defs>
                <CartesianGrid vertical={false} />
                <XAxis
                  dataKey="date"
                  tickLine={true}
                  axisLine={true}
                  tickMargin={8}
                  minTickGap={32}
                  tickFormatter={(value) => {
                    const date = new Date(value);
                    return date.toLocaleDateString("en-US", {
                      month: "short",
                      day: "numeric",
                    });
                  }}
                />
                <ChartTooltip
                  cursor={false}
                  content={
                    <ChartTooltipContent
                      labelFormatter={(value) => {
                        return new Date(value).toLocaleDateString("en-US", {
                          month: "short",
                          day: "numeric",
                        });
                      }}
                      indicator="dot"
                    />
                  }
                />
                <Area
                  dataKey="promptTokens"
                  type="natural"
                  fill="url(#fillPromptTokens)"
                  stroke="var(--color-promptTokens)"
                  stackId="a"
                />
                <Area
                  dataKey="responseTokens"
                  type="natural"
                  fill="url(#fillResponseTokens)"
                  stroke="var(--color-responseTokens)"
                  stackId="a"
                />
                <ChartLegend content={<ChartLegendContent />} />
              </AreaChart>
            </ChartContainer>
          )}
        </CardContent>
      </Card>
    </div>
  );
}

export default Home;
