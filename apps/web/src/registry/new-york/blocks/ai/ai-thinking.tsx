"use client";

import { Card, CardContent } from "@/components/ui/card";
import { Spinner } from "@/components/ui/spinner";

type AIThinkingProps = {
  spinner?: boolean;
  message?: string;
};

export default function AIThinking({
  spinner = true,
  message = "Let me think of the best answer...",
}: AIThinkingProps) {
  return (
    <Card className="border-dashed">
      <CardContent className="flex min-h-32 items-center gap-3 p-4">
        {spinner ? <Spinner className="size-4 text-muted-foreground" /> : null}
        <p className="text-sm text-muted-foreground">{message}</p>
      </CardContent>
    </Card>
  );
}
