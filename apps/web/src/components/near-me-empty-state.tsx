import { MapPin } from "@phosphor-icons/react/dist/ssr";

import {
  Empty,
  EmptyDescription,
  EmptyHeader,
  EmptyMedia,
  EmptyTitle,
} from "@/components/ui/empty";

export function NearMeEmptyState() {
  return (
    <Empty className="border border-border/60 bg-muted/20">
      <EmptyHeader>
        <EmptyMedia variant="icon">
          <MapPin className="size-5" />
        </EmptyMedia>
        <EmptyTitle>Location Not Detected</EmptyTitle>
        <EmptyDescription>
          We couldn&apos;t detect your location. Please type your city into the
          search bar to find local hackathons.
        </EmptyDescription>
      </EmptyHeader>
    </Empty>
  );
}
