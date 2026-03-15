"use client";

import { useEffect, useMemo, useRef, useState } from "react";
import { ChevronDown } from "lucide-react";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { Button } from "@/components/ui/button";

type SearchAIOverviewProps = {
  query: string;
};

type OverviewResponse = {
  overview?: string;
};

function AIOverviewSkeleton() {
  return (
    <Card className="mb-6 border-muted-foreground/20 bg-card py-4">
      <CardHeader className="gap-2 px-4 pb-2">
        <Skeleton className="h-5 w-36" />
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
  const [expanded, setExpanded] = useState(false);
  const [showExpandButton, setShowExpandButton] = useState(false);

  const collapsedContentRef = useRef<HTMLDivElement | null>(null);
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
      setExpanded(false);
      setShowExpandButton(false);

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

  useEffect(() => {
    if (!overview || !collapsedContentRef.current) return;
    const element = collapsedContentRef.current;
    setShowExpandButton(element.scrollHeight > element.clientHeight + 1);
  }, [overview]);

  if (hideSection) return null;
  if (loading) return <AIOverviewSkeleton />;
  if (!overview) return null;

  return (
    <Card className="mb-6 border-muted-foreground/20 bg-card py-4">
      <CardHeader className="gap-2 px-4 pb-1">
        <h3 className="text-base font-semibold tracking-tight">AI Overview</h3>
        <p className="text-xs text-muted-foreground">Generated insight for "{normalizedQuery}"</p>
      </CardHeader>

      <CardContent className="px-4 pb-1">
        <div className="space-y-3">
          <div
            ref={collapsedContentRef}
            className={`relative text-[15px] leading-7 text-foreground/90 ${
              expanded ? "" : "max-h-52 overflow-hidden"
            }`}
          >
            <ReactMarkdown
              remarkPlugins={[remarkGfm]}
              components={{
                p: ({ ...props }) => <p {...props} className="mb-4 leading-7" />,
                ul: ({ ...props }) => <ul {...props} className="mb-4 list-disc space-y-2 pl-6" />,
                ol: ({ ...props }) => <ol {...props} className="mb-4 list-decimal space-y-2 pl-6" />,
                strong: ({ ...props }) => <strong {...props} className="font-semibold text-foreground" />,
                a: ({ ...props }) => (
                  <a
                    {...props}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="font-medium text-primary underline underline-offset-2"
                  />
                ),
              }}
            >
              {overview}
            </ReactMarkdown>

            {!expanded && showExpandButton && (
              <div className="pointer-events-none absolute inset-x-0 bottom-0 h-20 bg-linear-to-t from-card via-card/90 to-transparent" />
            )}
          </div>

          {showExpandButton && (
            <Button
              type="button"
              variant="outline"
              className="h-11 w-full rounded-full text-base"
              onClick={() => setExpanded((prev) => !prev)}
            >
              {expanded ? "Show less" : "Show more"}
              <ChevronDown
                className={`ml-2 size-4 transition-transform ${expanded ? "rotate-180" : ""}`}
              />
            </Button>
          )}
        </div>
      </CardContent>
    </Card>
  );
}
