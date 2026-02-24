"use client";

import * as React from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import {
  Trophy,
  Bookmark,
  User,
  FunnelIcon,
} from "@phosphor-icons/react/dist/ssr";
import { cn } from "@/lib/utils";
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
} from "@/components/ui/command";

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
  const [open, setOpen] = React.useState(false);
  const pathname = usePathname();

  React.useEffect(() => {
    const handleScroll = () => {
      setScrolled(window.scrollY > 10);
    };

    handleScroll();

    window.addEventListener("scroll", handleScroll);
    return () => window.removeEventListener("scroll", handleScroll);
  }, []);

  React.useEffect(() => {
    const down = (e: KeyboardEvent) => {
      if (e.key === "k" && (e.metaKey || e.ctrlKey)) {
        e.preventDefault();
        setOpen((open) => !open);
      }
    };
    document.addEventListener("keydown", down);
    return () => document.removeEventListener("keydown", down);
  }, []);

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
          <Button
            variant="ghost"
            size="icon"
            className="h-8 w-8"
            onClick={() => setOpen(true)}
          >
            <FunnelIcon className="h-5 w-5" />
          </Button>

          <CommandDialog open={open} onOpenChange={setOpen}>
            <CommandInput placeholder="Filter hackathons..." />
            <CommandList>
              <CommandEmpty>No results found.</CommandEmpty>
              <CommandGroup heading="Filters">
                <CommandItem>All Hackathons</CommandItem>
                <CommandItem>Upcoming</CommandItem>
                <CommandItem>Ongoing</CommandItem>
                <CommandItem>Past</CommandItem>
              </CommandGroup>
              <CommandSeparator />
              <CommandGroup heading="Location">
                <CommandItem>Online</CommandItem>
                <CommandItem>In-Person</CommandItem>
                <CommandItem>Hybrid</CommandItem>
              </CommandGroup>
              <CommandSeparator />
              <CommandGroup heading="Prize Range">
                <CommandItem>$0 - $1,000</CommandItem>
                <CommandItem>$1,000 - $10,000</CommandItem>
                <CommandItem>$10,000+</CommandItem>
              </CommandGroup>
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
