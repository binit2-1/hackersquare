"use client";

import { useEffect, useMemo, useState } from "react";
import { Sparkles } from "lucide-react";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";

type SearchAIOverviewProps = {
  query: string;
};

type OverviewResponse = {
  overview?: string;
};

function AIOverviewSkeleton() {
  return (
    <Card className="mb-6 border-muted-foreground/20 bg-linear-to-b from-muted/30 to-background py-4">
      <CardHeader className="gap-2 px-4 pb-2">
        <div className="flex items-center gap-2">
          <Skeleton className="h-6 w-6 rounded-full" />
          <Skeleton className="h-5 w-36" />
        </div>
        <Skeleton className="h-4 w-2/3" />
      </CardHeader>
      <CardContent className="space-y-2 px-4 pb-1">
        <Skeleton className="h-4 w-full" />
        <Skeleton className="h-4 w-[96%]" />
        <Skeleton className="h-4 w-[92%]" />
        <Skeleton className="h-4 w-[88%]" />
      </CardContent>
    </Card>
  );
}

export function SearchAIOverview({ query }: SearchAIOverviewProps) {
  const [loading, setLoading] = useState(true);
  const [overview, setOverview] = useState("");
  const [hideSection, setHideSection] = useState(false);

  const normalizedQuery = useMemo(() => query.trim(), [query]);

  useEffect(() => {
    if (!normalizedQuery) {
      setLoading(false);
      setHideSection(true);
      setOverview("");
      return;
    }

    const controller = new AbortController();

    const fetchOverview = async () => {
      setLoading(true);
      setHideSection(false);
      setOverview("");

      try {
        const response = await fetch(
          `/api/v1/search/overview?q=${encodeURIComponent(normalizedQuery)}`,
          {
            method: "GET",
            credentials: "include",
            signal: controller.signal,
          },
        );

        if (!response.ok) {
          setHideSection(true);
          return;
        }

        const data = (await response.json()) as OverviewResponse;
        const content = (data.overview || "").trim();

        if (!content) {
          setHideSection(true);
          return;
        }

        setOverview(content);
      } catch {
        setHideSection(true);
      } finally {
        setLoading(false);
      }
    };

    fetchOverview();

    return () => {
      controller.abort();
    };
  }, [normalizedQuery]);

  if (hideSection) return null;
  if (loading) return <AIOverviewSkeleton />;
  if (!overview) return null;

  return (
    <Card className="mb-6 border-blue-200/40 bg-linear-to-b from-blue-50/50 to-background py-4 dark:border-blue-400/20 dark:from-blue-950/30">
      <CardHeader className="gap-2 px-4 pb-1">
        <div className="flex items-center gap-2">
          <Sparkles className="size-4 text-blue-600 dark:text-blue-400" />
          <h3 className="text-base font-semibold tracking-tight">AI Overview</h3>
        </div>
        <p className="text-xs text-muted-foreground">Generated insight for "{normalizedQuery}"</p>
      </CardHeader>

      <CardContent className="px-4 pb-1">
        <div className="space-y-3 text-[15px] leading-7 text-foreground/90">
          <ReactMarkdown
            remarkPlugins={[remarkGfm]}
            components={{
              p: ({ ...props }) => <p {...props} className="leading-7" />,
              ul: ({ ...props }) => <ul {...props} className="list-disc space-y-2 pl-6" />,
              ol: ({ ...props }) => <ol {...props} className="list-decimal space-y-2 pl-6" />,
              strong: ({ ...props }) => <strong {...props} className="font-semibold text-foreground" />,
              a: ({ ...props }) => (
                <a
                  {...props}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="font-medium text-blue-600 underline underline-offset-2 dark:text-blue-400"
                />
              ),
            }}
          >
            {overview}
          </ReactMarkdown>
        </div>
      </CardContent>
    </Card>
  );
}
