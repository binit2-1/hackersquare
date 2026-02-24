"use client";
import { cn } from "@/lib/utils";
import { HackathonProps } from "@/temp/hackathons";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import {
  CommandDialog,
  CommandInput,
  CommandList,
  CommandEmpty,
  CommandGroup,
  CommandItem,
  CommandSeparator,
  CommandShortcut,
} from "@/components/ui/command";
import Link from "next/link";
import { usePathname } from "next/navigation";
import {
  Trophy,
  Bookmark,
  User,
  FunnelIcon,
  MagnifyingGlass,
  Check,
  Calendar,
  MapPin,
  CurrencyDollar,
  ArrowRight,
} from "@phosphor-icons/react/dist/ssr";
import * as React from "react";

interface FilterState {
  status: string | null;
  location: string | null;
  prizeRange: string | null;
}

const STATUS_FILTERS = [
  { value: "all", label: "All Hackathons" },
  { value: "upcoming", label: "Upcoming" },
  { value: "ongoing", label: "Ongoing" },
  { value: "past", label: "Past" },
];

const LOCATION_FILTERS = [
  { value: "online", label: "Online" },
  { value: "in-person", label: "In-Person" },
  { value: "hybrid", label: "Hybrid" },
];

const PRIZE_FILTERS = [
  { value: "0-1000", label: "$0 - $1,000" },
  { value: "1000-10000", label: "$1,000 - $10,000" },
  { value: "10000+", label: "$10,000+" },
];

const navItems = [
  {
    name: "Hackathons",
    href: "/search",
    icon: Trophy,
  },
  {
    name: "Bookmarked",
    href: "/bookmarks",
    icon: Bookmark,
  },
];

export function Navbar() {
  const [scrolled, setScrolled] = React.useState(false);
  const [filterOpen, setFilterOpen] = React.useState(false);
  const [searchOpen, setSearchOpen] = React.useState(false);
  const [searchQuery, setSearchQuery] = React.useState("");
  const [debouncedQuery, setDebouncedQuery] = React.useState("");
  const [isLoading, setIsLoading] = React.useState(false);
  const [hackathons, setHackathons] = React.useState<HackathonProps[]>([]);
  const [activeFilters, setActiveFilters] = React.useState<FilterState>({
    status: null,
    location: null,
    prizeRange: null,
  });
  const pathname = usePathname();

  // Scroll effect
  React.useEffect(() => {
    const handleScroll = () => {
      setScrolled(window.scrollY > 10);
    };

    handleScroll();

    window.addEventListener("scroll", handleScroll);
    return () => window.removeEventListener("scroll", handleScroll);
  }, []);

  // Keyboard shortcuts: Cmd/Ctrl+K for search, Cmd/Ctrl+F for filter
  React.useEffect(() => {
    const down = (e: KeyboardEvent) => {
      if (e.key === "k" && (e.metaKey || e.ctrlKey)) {
        e.preventDefault();
        setSearchOpen((open) => !open);
        setFilterOpen(false);
      }
      if (e.key === "f" && (e.metaKey || e.ctrlKey)) {
        e.preventDefault();
        setFilterOpen((open) => !open);
        setSearchOpen(false);
      }
    };
    document.addEventListener("keydown", down);
    return () => document.removeEventListener("keydown", down);
  }, []);

  // Fetch hackathons
  React.useEffect(() => {
    const fetchHackathons = async () => {
      setIsLoading(true);
      try {
        const response = await fetch("http://localhost:8080/api/hackathons");
        if (response.ok) {
          const data = await response.json();
          setHackathons(data);
        }
      } catch (error) {
        console.error("Failed to fetch hackathons:", error);
      } finally {
        setIsLoading(false);
      }
    };
    if (searchOpen) {
      fetchHackathons();
    }
  }, [searchOpen]);

  // Debounce search query for better performance
  React.useEffect(() => {
    const timer = setTimeout(() => {
      setDebouncedQuery(searchQuery);
    }, 150);
    return () => clearTimeout(timer);
  }, [searchQuery]);

  // Memoized filtered results for performance
  const filteredHackathons = React.useMemo(() => {
    if (!debouncedQuery.trim()) return hackathons.slice(0, 10);
    const query = debouncedQuery.toLowerCase();
    return hackathons
      .filter(
        (h) =>
          h.title.toLowerCase().includes(query) ||
          h.host.toLowerCase().includes(query) ||
          h.tags.some((tag) => tag.toLowerCase().includes(query)),
      )
      .slice(0, 10);
  }, [debouncedQuery, hackathons]);

  // Count active filters
  const activeFilterCount = React.useMemo(() => {
    return Object.values(activeFilters).filter(Boolean).length;
  }, [activeFilters]);

  const handleFilterSelect = (
    type: keyof FilterState,
    value: string | null,
  ) => {
    setActiveFilters((prev) => ({
      ...prev,
      [type]: prev[type] === value ? null : value,
    }));
  };

  const clearFilters = () => {
    setActiveFilters({ status: null, location: null, prizeRange: null });
  };

  return (
    <header
      className={cn(
        "sticky top-0 z-50 w-full transition-all duration-200",
        scrolled
          ? "backdrop-blur-md border-b border-border/40 bg-background/80"
          : "bg-background",
      )}
    >
      <div className="flex h-14 items-center">
        {/* Logo - Far Left */}
        <div className="shrink-0 px-4 sm:px-6 lg:px-8">
          <Link href="/" className="block">
            <div className="h-8 w-8 rounded-full bg-primary" />
          </Link>
        </div>

        {/* Center - Navigation inside max-w container (left aligned) */}
        <div className="flex-1 max-w-196.5 mx-auto flex items-center justify-between">
          <nav className="flex items-center gap-px">
            {navItems.map((item) => {
              const Icon = item.icon;
              const isActive = pathname === item.href;
              return (
                <Link
                  key={item.href}
                  href={item.href}
                  className={cn(
                    "flex items-center gap-1 sm:gap-2 px-2 sm:px-3 py-2 rounded-md font-medium transition-colors",
                    "text-xs sm:text-sm lg:text-base",
                    isActive
                      ? "text-foreground"
                      : "text-muted-foreground hover:text-foreground",
                  )}
                >
                  <Icon className="h-5 w-5 shrink-0" weight="regular" />
                  <span className="hidden sm:inline">{item.name}</span>
                </Link>
              );
            })}
          </nav>
          <div className="flex items-center gap-1">
            <Button
              variant="ghost"
              size="icon"
              className="h-8 w-8 relative"
              onClick={() => setSearchOpen(true)}
              title="Search (Ctrl/Cmd + K)"
            >
              <MagnifyingGlass className="h-5 w-5" />
            </Button>
            <Button
              variant="ghost"
              size="icon"
              className="h-8 w-8 relative"
              onClick={() => setFilterOpen(true)}
              title="Filter (Ctrl/Cmd + F)"
            >
              <FunnelIcon className="h-5 w-5" />
              {activeFilterCount > 0 && (
                <span className="absolute -top-0.5 -right-0.5 h-4 w-4 rounded-full bg-primary text-[10px] font-medium text-primary-foreground flex items-center justify-center">
                  {activeFilterCount}
                </span>
              )}
            </Button>
          </div>

          {/* Search Dialog */}
          <CommandDialog
            open={searchOpen}
            onOpenChange={setSearchOpen}
            title="Search Hackathons"
            description="Search for hackathons by title, host, or tags"
          >
            <CommandInput
              placeholder="Search hackathons..."
              value={searchQuery}
              onValueChange={setSearchQuery}
            />
            <CommandList className="max-h-[400px]">
              {isLoading ? (
                <div className="py-8 text-center text-sm text-muted-foreground">
                  <div className="flex items-center justify-center gap-2">
                    <div className="h-4 w-4 animate-spin rounded-full border-2 border-primary border-t-transparent" />
                    Loading hackathons...
                  </div>
                </div>
              ) : (
                <>
                  <CommandEmpty>
                    {debouncedQuery.trim() ? (
                      <div className="flex flex-col items-center gap-2 py-6">
                        <span>No hackathons found for "{debouncedQuery}"</span>
                        <span className="text-xs text-muted-foreground">
                          Try searching by title, host, or tags
                        </span>
                      </div>
                    ) : (
                      "Start typing to search hackathons..."
                    )}
                  </CommandEmpty>
                  {filteredHackathons.length > 0 && (
                    <CommandGroup
                      heading={
                        debouncedQuery.trim()
                          ? `Results (${filteredHackathons.length} found)`
                          : "Recent Hackathons"
                      }
                    >
                      {filteredHackathons.map((hackathon) => (
                        <CommandItem
                          key={hackathon.id}
                          onSelect={() => {
                            setSearchOpen(false);
                            window.location.href = `/hackathon/${hackathon.id}`;
                          }}
                          className="cursor-pointer"
                        >
                          <div className="flex items-center gap-3 w-full">
                            <div className="flex-1 min-w-0">
                              <p className="font-medium truncate">
                                {hackathon.title}
                              </p>
                              <p className="text-xs text-muted-foreground truncate">
                                {hackathon.host}
                              </p>
                            </div>
                            <ArrowRight className="h-4 w-4 text-muted-foreground shrink-0 opacity-0 group-data-[selected=true]:opacity-100 transition-opacity" />
                          </div>
                        </CommandItem>
                      ))}
                    </CommandGroup>
                  )}
                </>
              )}
            </CommandList>
          </CommandDialog>

          {/* Filter Dialog */}
          <CommandDialog
            open={filterOpen}
            onOpenChange={setFilterOpen}
            title="Filter Hackathons"
            description="Filter hackathons by status, location, and prize range"
          >
            <CommandInput placeholder="Filter hackathons..." />
            <CommandList className="max-h-[400px]">
              <CommandEmpty>No matching filters.</CommandEmpty>

              <CommandGroup heading="Status">
                {STATUS_FILTERS.map((filter) => (
                  <CommandItem
                    key={filter.value}
                    onSelect={() =>
                      handleFilterSelect("status", filter.value === "all" ? null : filter.value)
                    }
                    className="cursor-pointer"
                  >
                    <div className="flex items-center gap-2 flex-1">
                      <Calendar className="h-4 w-4 text-muted-foreground" />
                      <span>{filter.label}</span>
                    </div>
                    {(activeFilters.status === filter.value ||
                      (filter.value === "all" && !activeFilters.status)) && (
                      <Check className="h-4 w-4 text-primary" />
                    )}
                  </CommandItem>
                ))}
              </CommandGroup>

              <CommandSeparator />

              <CommandGroup heading="Location">
                {LOCATION_FILTERS.map((filter) => (
                  <CommandItem
                    key={filter.value}
                    onSelect={() => handleFilterSelect("location", filter.value)}
                    className="cursor-pointer"
                  >
                    <div className="flex items-center gap-2 flex-1">
                      <MapPin className="h-4 w-4 text-muted-foreground" />
                      <span>{filter.label}</span>
                    </div>
                    {activeFilters.location === filter.value && (
                      <Check className="h-4 w-4 text-primary" />
                    )}
                  </CommandItem>
                ))}
              </CommandGroup>

              <CommandSeparator />

              <CommandGroup heading="Prize Range">
                {PRIZE_FILTERS.map((filter) => (
                  <CommandItem
                    key={filter.value}
                    onSelect={() => handleFilterSelect("prizeRange", filter.value)}
                    className="cursor-pointer"
                  >
                    <div className="flex items-center gap-2 flex-1">
                      <CurrencyDollar className="h-4 w-4 text-muted-foreground" />
                      <span>{filter.label}</span>
                    </div>
                    {activeFilters.prizeRange === filter.value && (
                      <Check className="h-4 w-4 text-primary" />
                    )}
                  </CommandItem>
                ))}
              </CommandGroup>

              {activeFilterCount > 0 && (
                <>
              <CommandSeparator />
              <CommandGroup>
                <CommandItem
                  onSelect={clearFilters}
                        className="cursor-pointer text-muted-foreground hover:text-foreground"
                >
                  <span className="flex-1 text-center text-xs">
                          Clear all filters ({activeFilterCount} active)
                  </span>
                </CommandItem>
              </CommandGroup>
            </>
              )}
            </CommandList>
          </CommandDialog>
        </div>

        {/* Profile Avatar - Extreme Right */}
        <div className="shrink-0 px-4 sm:px-6 lg:px-8">
          <Avatar className="h-8 w-8 cursor-pointer">
            <AvatarFallback className="bg-muted text-muted-foreground">
              <User className="h-4 w-4" weight="regular" />
            </AvatarFallback>
          </Avatar>
        </div>
      </div>
    </header>
  );
}
