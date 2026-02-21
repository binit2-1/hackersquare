"use client";

import { useState } from "react";
import {
  Bookmark,
  Calendar,
  MapPin,
  Trophy,
  Building,
} from "@phosphor-icons/react/dist/ssr";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { cn } from "@/lib/utils";
import type { HackathonProps } from "@/temp/hackathons";
import { useBookmarks } from "@/hooks/useBookmarks";
function formatDateRange(start: string, end: string): string {
  const s = new Date(start);
  const e = new Date(end);
  const opts: Intl.DateTimeFormatOptions = { month: "short", day: "numeric" };
  return `${s.toLocaleDateString("en-US", opts)} – ${e.toLocaleDateString("en-US", { ...opts, year: "numeric" })}`;
}

interface HackathonCardProps {
  hackathon: HackathonProps;
}

export function HackathonCard({ hackathon }: HackathonCardProps) {
  const {isBookmarked, toggleBookmark} = useBookmarks();
  const bookmarked =  isBookmarked(hackathon.title);

  return (
    <Card className="flex flex-col gap-0 py-0 overflow-hidden">
      <CardHeader className="flex flex-row items-start justify-between gap-3 px-5 pt-5 pb-3">
        <div className="flex flex-col gap-1 min-w-0">
          <CardTitle className="text-base leading-snug line-clamp-2">
            {hackathon.title}
          </CardTitle>
          <div className="flex items-center gap-1.5 text-muted-foreground text-xs">
            <Building className="size-3 shrink-0" />
            <span className="truncate">{hackathon.host}</span>
          </div>
        </div>
        <Button
          variant="ghost"
          size="icon"
          className="shrink-0 -mt-1 -mr-1"
          aria-label={bookmarked ? "Remove bookmark" : "Bookmark hackathon"}
          onClick={() => toggleBookmark(hackathon)}
        >
          <Bookmark
            className={cn(
              "size-4 transition-colors",
              bookmarked && "fill-primary"
            )}
            weight={bookmarked ? "fill" : "regular"}
          />
        </Button>
      </CardHeader>

      <CardContent className="flex flex-col gap-2.5 px-5 pb-4 text-sm">
        <div className="flex items-center gap-2 text-muted-foreground">
          <Calendar className="size-3.5 shrink-0" />
          <span>{formatDateRange(hackathon.startDate, hackathon.endDate)}</span>
        </div>

        <div className="flex items-center gap-2 text-muted-foreground">
          <MapPin className="size-3.5 shrink-0" />
          <span>{hackathon.location}</span>
        </div>

        <div className="flex items-center gap-2 font-medium">
          <Trophy className="size-3.5 shrink-0 text-amber-500" />
          <span>{hackathon.prize}</span>
        </div>

        <div className="flex flex-wrap gap-1.5 pt-1">
          {hackathon.tags.map((tag) => (
            <Badge key={tag} variant="secondary" className="text-xs px-2 py-0.5">
              {tag}
            </Badge>
          ))}
        </div>
      </CardContent>

      <CardFooter className="px-5 pb-5 pt-0">
        <Button asChild className="w-full" size="sm">
          <a href={hackathon.applyUrl} target="_blank" rel="noopener noreferrer">
            Apply Now
          </a>
        </Button>
      </CardFooter>
    </Card>
  );
}
