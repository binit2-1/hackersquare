"use client";
import { cn } from "@/lib/utils";
import { HackathonProps } from "@/models/hackathon";
import { Facehash } from "facehash";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuTrigger,
  DropdownMenuContent,
  DropdownMenuItem,
} from "@/components/ui/dropdown-menu";
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
import Image from "next/image";
import { usePathname, useSearchParams, useRouter } from "next/navigation";
import {
  Trophy,
  Bookmark,
  UserIcon,
  FunnelIcon,
  MagnifyingGlass,
  Check,
  Calendar,
  MapPin,
  CurrencyDollar,
  ArrowRight,
} from "@phosphor-icons/react/dist/ssr";
import * as React from "react";
import AuthModal from "@/components/auth-modal";
import { useAuth } from "@/context/AuthContext";

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
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();

  const [scrolled, setScrolled] = React.useState(false);
  const [filterOpen, setFilterOpen] = React.useState(false);
  const [searchOpen, setSearchOpen] = React.useState(false);
  const [showAuthButton, setShowAuthButton] = React.useState(true);
  const [searchQuery, setSearchQuery] = React.useState("");
  const [showAuthModal, setShowAuthModal] = React.useState(false);
  const [authInitialTab, setAuthInitialTab] = React.useState<
    "signin" | "signup"
  >("signin");
  const [activeFilters, setActiveFilters] = React.useState<FilterState>({
    status: searchParams.get("status") || null,
    location: searchParams.get("location") || null,
    prizeRange: searchParams.get("prizeRange") || null,
  });

  const { user, isLoading, logout } = useAuth();
  const isAuthed = Boolean(user);
  const authReady = !isLoading;
  const showAuthActions = authReady && !isAuthed;

  // Scroll effect
  React.useEffect(() => {
    const handleScroll = () => {
      setScrolled(window.scrollY > 10);
    };

    handleScroll();

    window.addEventListener("scroll", handleScroll);
    return () => window.removeEventListener("scroll", handleScroll);
  }, []);

  // Show auth button until 1300px width; enable avatar dropdown only when hidden
  React.useEffect(() => {
    if (typeof window === "undefined") return;
    const m = window.matchMedia("(min-width: 1300px)");
    const handler = (e: any) => setShowAuthButton(e.matches);
    // set initial value from media query
    setShowAuthButton(m.matches);
    // prefer addEventListener if available
    if (m.addEventListener) {
      m.addEventListener("change", handler);
    } else {
      // older browsers
      // @ts-ignore
      m.addListener(handler);
    }

    return () => {
      if (m.removeEventListener) {
        m.removeEventListener("change", handler);
      } else {
        // @ts-ignore
        m.removeListener(handler);
      }
    };
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

  // Count active filters
  const activeFilterCount = React.useMemo(() => {
    return Object.values(activeFilters).filter(Boolean).length;
  }, [activeFilters]);

  const handleFilterSelect = (
    type: keyof FilterState,
    value: string | null,
  ) => {
    const newValue = activeFilters[type] === value ? null : value;

    setActiveFilters((prev) => ({ ...prev, [type]: newValue }));

    const current = new URLSearchParams(Array.from(searchParams.entries()));

    if (newValue) {
      current.set(type, newValue);
    } else {
      current.delete(type);
    }

    current.set("page", "1");

    const search = current.toString();

    router.push(`/search${search ? `?${search}` : ""}`);
  };

  const clearFilters = () => {
    setActiveFilters({ status: null, location: null, prizeRange: null });
    router.push(pathname);
  };

  // Handles global search routing
  const handleSearchSubmit = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "Enter" && searchQuery.trim()) {
      setSearchOpen(false);
      const current = new URLSearchParams(Array.from(searchParams.entries()));
      current.set("q", searchQuery.trim());
      current.delete("page");
      router.push(`/search?${current.toString()}`);
    }
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
      <div className="relative flex h-14 items-center">
        <div className="shrink-0 px-4 sm:px-6 lg:px-8">
          <Link href="/" className="block">
            <Image
              src="/logoHS.svg"
              alt="HackerSquare logo"
              width={28}
              height={28}
              className="h-8 w-8 rounded-full"
              priority
            />
          </Link>
        </div>

        <div className="absolute left-1/2 top-1/2 z-0 flex w-[calc(100%-7.5rem)] max-w-196.5 -translate-x-1/2 -translate-y-1/2 items-center justify-between px-1 sm:w-[calc(100%-8.5rem)] sm:px-2 lg:w-full">
          <nav className="flex items-center gap-px">
            {navItems.map((item) => {
              const Icon = item.icon;
              const isActive = pathname === item.href;
              return (
                <Link
                  key={item.href}
                  href={item.href}
                  className={cn(
                    "flex items-center gap-1 sm:gap-2 px-2 sm:px-3 py-2 rounded-md font-medium transition-colors text-xs sm:text-sm lg:text-base",
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

          <CommandDialog
            open={searchOpen}
            onOpenChange={setSearchOpen}
            title="Search Hackathons"
          >
            <CommandInput
              placeholder="Search hackathons and press Enter..."
              value={searchQuery}
              onValueChange={setSearchQuery}
              onKeyDown={handleSearchSubmit}
            />
            <CommandList className="max-h-100">
              <CommandEmpty>
                {searchQuery.trim()
                  ? `Press Enter to search for "${searchQuery}"`
                  : "Start typing to search..."}
              </CommandEmpty>
            </CommandList>
          </CommandDialog>

          <CommandDialog
            open={filterOpen}
            onOpenChange={setFilterOpen}
            title="Filter Hackathons"
          >
            <CommandInput placeholder="Search filters..." />
            <CommandList className="max-h-100">
              <CommandEmpty>No matching filters.</CommandEmpty>

              <CommandGroup heading="Status">
                {STATUS_FILTERS.map((filter) => (
                  <CommandItem
                    key={filter.value}
                    onSelect={() =>
                      handleFilterSelect(
                        "status",
                        filter.value === "all" ? null : filter.value,
                      )
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
                    onSelect={() =>
                      handleFilterSelect("location", filter.value)
                    }
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
                    onSelect={() =>
                      handleFilterSelect("prizeRange", filter.value)
                    }
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

        <div className="absolute right-4 top-1/2 z-10 flex -translate-y-1/2 flex-row items-center gap-4 sm:right-6 lg:right-8">
          {showAuthButton && showAuthActions ? (
            <Button
              className="inline-flex min-w-36 rounded-4xl px-4 justify-center"
              size="sm"
              onClick={() => setShowAuthModal(true)}
            >
              Sign In / Sign Up
            </Button>
          ) : user ? (
            <div className="flex items-center gap-2">
              <span className="hidden sm:inline text-sm font-medium">
                {user.name}
              </span>
            </div>
          ) : null}

          {/* Avatar: only act as a dropdown trigger when auth button is hidden (mobile/smaller screens) */}
          {!showAuthButton && showAuthActions ? (
            <DropdownMenu>
              <DropdownMenuTrigger className="h-8 w-8 rounded-full p-0">
                <Avatar className="h-8 w-8">
                  <AvatarFallback className="bg-muted text-muted-foreground">
                    <UserIcon className="h-4 w-4" weight="regular" />
                  </AvatarFallback>
                </Avatar>
              </DropdownMenuTrigger>

              <DropdownMenuContent align="end" className="w-40">
                <DropdownMenuItem
                  onSelect={() => {
                    setAuthInitialTab("signin");
                    setShowAuthModal(true);
                  }}
                >
                  Sign In
                </DropdownMenuItem>
                <DropdownMenuItem
                  onSelect={() => {
                    setAuthInitialTab("signup");
                    setShowAuthModal(true);
                  }}
                >
                  Sign Up
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          ) : isAuthed ? (
            <DropdownMenu>
              <DropdownMenuTrigger className="rounded-full p-0 focus:outline-none focus-visible:ring-2 focus-visible:ring-ring">
                <Facehash
                  name={user!.name}
                  size={32}
                  intensity3d="subtle"
                  showInitial={true}
                  variant="solid"
                  enableBlink
                  className="rounded-full text-black"
                  colors={["#FFFFFF", "#FFFFFF", "#FFFFFF"]}
                />
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end" className="w-40">
                <DropdownMenuItem
                  onSelect={() => {
                    logout();
                  }}
                >
                  Logout
                </DropdownMenuItem>
                <DropdownMenuItem
                  onSelect={() => {
                    router.push("/profile");
                  }}
                >
                  Profile
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          ) : (
            <Avatar className="h-8 w-8 cursor-default">
              <AvatarFallback className="bg-muted text-muted-foreground">
                <UserIcon className="h-4 w-4" weight="regular" />
              </AvatarFallback>
            </Avatar>
          )}
        </div>
      </div>
      <AuthModal
        open={showAuthModal}
        onOpenChange={setShowAuthModal}
        initialTab={authInitialTab}
      />
    </header>
  );
}
