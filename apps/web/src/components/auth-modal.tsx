"use client";

import * as React from "react";

import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from "@/components/ui/dialog";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { cn } from "@/lib/utils";
import { useRouter } from "next/navigation";
import { useAuth } from "@/context/AuthContext";


/* ── inline SVG logos (avoids extra icon deps) ── */

function GitHubLogo({ className }: { className?: string }) {
  return (
    <svg className={className} viewBox="0 0 24 24" fill="currentColor">
      <path d="M12 .297c-6.63 0-12 5.373-12 12 0 5.303 3.438 9.8 8.205 11.385.6.113.82-.258.82-.577 0-.285-.01-1.04-.015-2.04-3.338.724-4.042-1.61-4.042-1.61C4.422 18.07 3.633 17.7 3.633 17.7c-1.087-.744.084-.729.084-.729 1.205.084 1.838 1.236 1.838 1.236 1.07 1.835 2.809 1.305 3.495.998.108-.776.417-1.305.76-1.605-2.665-.3-5.466-1.332-5.466-5.93 0-1.31.465-2.38 1.235-3.22-.135-.303-.54-1.523.105-3.176 0 0 1.005-.322 3.3 1.23.96-.267 1.98-.399 3-.405 1.02.006 2.04.138 3 .405 2.28-1.552 3.285-1.23 3.285-1.23.645 1.653.24 2.873.12 3.176.765.84 1.23 1.91 1.23 3.22 0 4.61-2.805 5.625-5.475 5.92.42.36.81 1.096.81 2.22 0 1.606-.015 2.896-.015 3.286 0 .315.21.69.825.57C20.565 22.092 24 17.592 24 12.297c0-6.627-5.373-12-12-12" />
    </svg>
  );
}

/* ── types ── */

type Props = {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  /** Pre-select sign-up tab when opening from a "Sign Up" action */
  initialTab?: "signin" | "signup";
};

export default function AuthModal({
  open,
  onOpenChange,
  initialTab = "signin",
}: Props) {
  const router = useRouter();
  const { refreshUser } = useAuth();

  const [tab, setTab] = React.useState<"signin" | "signup">(initialTab);
  const [isLoading, setIsLoading] = React.useState(false);
  const [error, setError] = React.useState<string | null>(null);

  // Sync tab when initialTab prop changes (e.g. dropdown click)
  React.useEffect(() => {
    if (open) setTab(initialTab);
  }, [open, initialTab]);

  const [name, setName] = React.useState("");
  const [email, setEmail] = React.useState("");
  const [password, setPassword] = React.useState("");

  const isSignUp = tab === "signup";

  const reset = () => {
    setName("");
    setEmail("");
    setPassword("");
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    setError(null);

    let endpoint = "";
    let payload: Record<string, string> = {};
    if (isSignUp) {
      endpoint = "/v1/auth/register";
      payload = { name, email, password };
    } else {
      endpoint = "/v1/auth/login";
      payload = {email, password}
    }

    try{
      const res = await fetch(`/api${endpoint}`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(payload),
        credentials: "include", 
      });

      if (!res.ok){
        const err = await res.text();
        throw new Error(err || "Authentication failed");
      }
      await refreshUser();
      reset();
      onOpenChange(false);
    } catch (err) {
      setError((err as Error).message);
    } finally{
      setIsLoading(false);
    }
  };

  const handleOAuth = () => {
    window.location.href = "/api/v1/auth/github/login";
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-105 p-0 gap-0 overflow-hidden">
        {/* ── header ── */}
        <DialogHeader className="px-6 pt-6 pb-2 space-y-1">
          <DialogTitle className="text-xl font-semibold tracking-tight text-center">
            {isSignUp ? "Create an account" : "Welcome back"}
          </DialogTitle>
          <DialogDescription className="text-center text-sm text-muted-foreground">
            {isSignUp
              ? "Enter your details to get started"
              : "Sign in to continue to HackerSquare"}
          </DialogDescription>
        </DialogHeader>

        {/* ── tab toggle ── */}
        <div className="px-6 pt-3">
          <div className="flex items-center rounded-lg bg-muted p-1">
            <button
              type="button"
              className={cn(
                "flex-1 rounded-md py-1.5 text-sm font-medium transition-all",
                !isSignUp
                  ? "bg-background text-foreground shadow-sm"
                  : "text-muted-foreground hover:text-foreground",
              )}
              onClick={() => setTab("signin")}
            >
              Sign In
            </button>
            <button
              type="button"
              className={cn(
                "flex-1 rounded-md py-1.5 text-sm font-medium transition-all",
                isSignUp
                  ? "bg-background text-foreground shadow-sm"
                  : "text-muted-foreground hover:text-foreground",
              )}
              onClick={() => setTab("signup")}
            >
              Sign Up
            </button>
          </div>
        </div>

        {/* ── form ── */}
        <form onSubmit={handleSubmit} className="px-6 pt-4 pb-6 space-y-4">
          {isSignUp && (
            <div className="space-y-1.5">
              <Label htmlFor="auth-name" className="text-sm">
                Name
              </Label>
              <Input
                id="auth-name"
                value={name}
                onChange={(e) => setName(e.target.value)}
                placeholder="John Doe"
                autoComplete="name"
              />
            </div>
          )}
          <div className="space-y-1.5">
            <Label htmlFor="auth-email" className="text-sm">Email</Label>
            <Input
              id="auth-email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="you@example.com"
              type="email"
              required
            />
          </div>

          <div className="space-y-1.5">
            <Label htmlFor="auth-password" className="text-sm">
              Password
            </Label>
            <Input
              id="auth-password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              type="password"
              placeholder="••••••••"
              autoComplete={isSignUp ? "new-password" : "current-password"}
            />
          </div>

          {error && <p className="text-sm font-medium text-red-500">{error}</p>}

          <Button type="submit" className="w-full h-10 text-sm font-medium">
            {isSignUp ? "Create account" : "Sign in"}
          </Button>

          {/* ── divider ── */}
          <div className="relative flex items-center py-1">
            <div className="h-px flex-1 bg-border" />
            <span className="px-3 text-xs text-muted-foreground select-none">
              or continue with
            </span>
            <div className="h-px flex-1 bg-border" />
          </div>

          {/* ── OAuth buttons ── */}
          <div className="grid gap-3">
            <Button
              type="button"
              variant="outline"
              className="h-10 gap-2 text-sm font-medium"
              onClick={handleOAuth}
            >
              <GitHubLogo className="h-4 w-4" />
              GitHub
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
}
