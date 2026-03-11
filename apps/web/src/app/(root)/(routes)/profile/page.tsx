"use client";

import { useState } from "react";
import { useAuth } from "@/context/AuthContext";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Separator } from "@/components/ui/separator";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Card, CardContent } from "@/components/ui/card";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import {
  User,
  PencilSimple,
  MapPin,
  Briefcase,
  GithubLogo,
  LinkedinLogo,
  XLogo,
  Globe,
  Star,
  Code,
  Robot,
  LinkSimple,
  Check,
  SpinnerGap,
} from "@phosphor-icons/react/dist/ssr";

// Placeholder data — will be wired to Go backend
const mockPinnedRepos = [
  {
    id: 1,
    name: "hackersquare",
    desc: "A hackathon aggregator pulling from Devfolio, MLH, and Unstop into one clean feed.",
    language: "Go",
    stars: 34,
    url: "#",
  },
  {
    id: 2,
    name: "react-tiny-dom",
    desc: "A lightweight React clone built from scratch to understand the virtual DOM.",
    language: "JavaScript",
    stars: 12,
    url: "#",
  },
];

export default function ProfilePage() {
  const { user, isLoading, refreshUser } = useAuth();
  const [isEditing, setIsEditing] = useState(false);
  const [saving, setSaving] = useState(false);
  const [savingLinks, setSavingLinks] = useState(false);

  // Form state — seeded from user context
  const [headline, setHeadline] = useState("");
  const [location, setLocation] = useState("");
  const [website, setWebsite] = useState("");
  const [linkedin, setLinkedin] = useState("");
  const [twitter, setTwitter] = useState("");
  const [hydrated, setHydrated] = useState(false);

  // Seed form state from user once loaded
  if (user && !hydrated) {
    setHeadline(user.headline || "");
    setLocation(user.location || "");
    setWebsite(user.website_url || "");
    setLinkedin(user.linkedin_url || "");
    setTwitter(user.twitter_url || "");
    setHydrated(true);
  }

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <SpinnerGap className="size-6 animate-spin text-muted-foreground" />
      </div>
    );
  }

  if (!user) {
    return (
      <div className="min-h-screen flex items-center justify-center px-4">
        <div className="text-center space-y-3">
          <User className="size-10 mx-auto text-muted-foreground" />
          <p className="text-lg font-medium">Sign in to view your profile</p>
          <p className="text-sm text-muted-foreground">
            Your profile, projects, and connections live here.
          </p>
        </div>
      </div>
    );
  }

  const handleSave = async () => {
    setSaving(true);
    try {
      const res = await fetch("http://localhost:8080/v1/users/profile", {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({
          headline,
          location,
          website_url: website,
          linkedin_url: linkedin,
          twitter_url: twitter,
        }),
      });
      if (!res.ok) throw new Error("Failed to save profile");
      await refreshUser();
      setIsEditing(false);
    } catch (err) {
      console.error("Profile save error:", err);
    } finally {
      setSaving(false);
    }
  };

  const handleSaveLinks = async () => {
    setSavingLinks(true);
    try {
      const res = await fetch("http://localhost:8080/v1/users/profile", {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({
          headline,
          location,
          website_url: website,
          linkedin_url: linkedin,
          twitter_url: twitter,
        }),
      });
      if (!res.ok) throw new Error("Failed to save links");
      await refreshUser();
    } catch (err) {
      console.error("Links save error:", err);
    } finally {
      setSavingLinks(false);
    }
  };

  const handleConnectGitHub = () => {
    window.location.href = "http://localhost:8080/v1/auth/github/connect";
  };

  const initials = user.name
    .split(" ")
    .map((w) => w[0])
    .join("")
    .slice(0, 2)
    .toUpperCase();

  return (
    <div className="min-h-screen bg-background">
      <div className="max-w-196.5 mx-auto px-4 py-8 sm:py-12">
        {/* ── Identity header ── */}
        <div className="flex flex-col sm:flex-row sm:items-start gap-5 sm:gap-6">
          <Avatar className="size-20 sm:size-24 shrink-0 text-xl sm:text-2xl font-bold border">
            <AvatarFallback className="bg-muted text-foreground">
              {initials}
            </AvatarFallback>
          </Avatar>

          <div className="flex-1 min-w-0 space-y-1.5">
            <div className="flex items-start justify-between gap-3">
              <h1 className="text-2xl sm:text-3xl font-bold tracking-tight truncate">
                {user.name}
              </h1>
              <Button
                variant={isEditing ? "default" : "outline"}
                size="sm"
                className="shrink-0"
                onClick={() => (isEditing ? handleSave() : setIsEditing(true))}
                disabled={saving}
              >
                {saving ? (
                  <SpinnerGap className="size-4 animate-spin" />
                ) : isEditing ? (
                  <>
                    <Check className="size-4" />
                    <span className="hidden sm:inline ml-1.5">Save</span>
                  </>
                ) : (
                  <>
                    <PencilSimple className="size-4" />
                    <span className="hidden sm:inline ml-1.5">Edit</span>
                  </>
                )}
              </Button>
            </div>

            {isEditing ? (
              <div className="space-y-2 max-w-sm">
                <Input
                  value={headline}
                  onChange={(e) => setHeadline(e.target.value)}
                  placeholder="Your one-liner"
                  className="text-sm"
                />
                <Input
                  value={location}
                  onChange={(e) => setLocation(e.target.value)}
                  placeholder="City, Country"
                  className="text-sm"
                />
              </div>
            ) : (
              <>
                <p className="text-sm text-muted-foreground flex items-center gap-1.5">
                  <Briefcase className="size-3.5 shrink-0" />
                  {headline || "No headline set"}
                </p>
                <p className="text-sm text-muted-foreground flex items-center gap-1.5">
                  <MapPin className="size-3.5 shrink-0" />
                  {location || "No location set"}
                </p>
              </>
            )}

            <p className="text-xs text-muted-foreground/60 pt-0.5">{user.email}</p>
          </div>
        </div>

        <Separator className="my-6 sm:my-8" />

        {/* ── Tabs ── */}
        <Tabs defaultValue="portfolio" className="w-full">
          <TabsList variant="line" className="mb-6">
            <TabsTrigger value="portfolio">Portfolio</TabsTrigger>
            <TabsTrigger value="connections">Connections</TabsTrigger>
          </TabsList>

          {/* ── Portfolio tab ── */}
          <TabsContent value="portfolio" className="space-y-8">
            {/* AI Skill Summary */}
            <Card className="border-dashed">
              <CardContent className="p-5 space-y-2">
                <h3 className="text-sm font-semibold flex items-center gap-2">
                  <Robot className="size-4 text-muted-foreground" />
                  AI Skill Summary
                </h3>
                <p className="text-sm text-muted-foreground leading-relaxed">
                  Strong fundamentals in React and Next.js. Recently demonstrating
                  growth in backend systems using Go and PostgreSQL. Actively
                  participating in open-source tooling.
                </p>
              </CardContent>
            </Card>

            {/* Pinned Projects */}
            <div>
              <h2 className="text-sm font-semibold flex items-center gap-2 mb-4">
                <GithubLogo className="size-4" />
                Top Projects
              </h2>

              <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                {mockPinnedRepos.map((repo) => (
                  <a
                    key={repo.id}
                    href={repo.url}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="group block"
                  >
                    <Card className="h-full transition-shadow group-hover:shadow-md">
                      <CardContent className="p-4 space-y-2.5">
                        <div className="flex items-start justify-between gap-2">
                          <span className="text-sm font-semibold text-foreground group-hover:underline underline-offset-2 truncate">
                            {repo.name}
                          </span>
                          <span className="flex items-center gap-1 text-xs text-muted-foreground shrink-0">
                            <Star className="size-3.5" weight="fill" />
                            {repo.stars}
                          </span>
                        </div>
                        <p className="text-xs text-muted-foreground leading-relaxed line-clamp-2">
                          {repo.desc}
                        </p>
                        <span className="inline-flex items-center gap-1 text-[11px] font-medium text-muted-foreground bg-muted px-2 py-0.5 rounded">
                          <Code className="size-3" />
                          {repo.language}
                        </span>
                      </CardContent>
                    </Card>
                  </a>
                ))}
              </div>
            </div>
          </TabsContent>

          {/* ── Connections tab ── */}
          <TabsContent value="connections" className="space-y-5 max-w-lg">
            {/* GitHub connection */}
            <Card>
              <CardContent className="p-4 flex items-center justify-between gap-4">
                <div className="flex items-center gap-3 min-w-0">
                  <div className="size-9 rounded-full bg-muted flex items-center justify-center shrink-0">
                    <GithubLogo className="size-5" />
                  </div>
                  <div className="min-w-0">
                    <p className="text-sm font-medium">GitHub</p>
                    <p className="text-xs text-muted-foreground truncate">
                      Required for projects &amp; AI analysis
                    </p>
                  </div>
                </div>
                <Button size="sm" onClick={handleConnectGitHub}>
                  Connect
                </Button>
              </CardContent>
            </Card>

            {/* Social links */}
            <Card>
              <CardContent className="p-4 space-y-4">
                <h3 className="text-sm font-semibold flex items-center gap-2">
                  <LinkSimple className="size-4 text-muted-foreground" />
                  Social Links
                </h3>

                <div className="space-y-3">
                  <div className="flex items-center gap-2.5">
                    <Globe className="size-4 text-muted-foreground shrink-0" />
                    <Input
                      value={website}
                      onChange={(e) => setWebsite(e.target.value)}
                      placeholder="https://yourwebsite.com"
                      className="text-sm"
                    />
                  </div>
                  <div className="flex items-center gap-2.5">
                    <LinkedinLogo className="size-4 text-muted-foreground shrink-0" />
                    <Input
                      value={linkedin}
                      onChange={(e) => setLinkedin(e.target.value)}
                      placeholder="https://linkedin.com/in/username"
                      className="text-sm"
                    />
                  </div>
                  <div className="flex items-center gap-2.5">
                    <XLogo className="size-4 text-muted-foreground shrink-0" />
                    <Input
                      value={twitter}
                      onChange={(e) => setTwitter(e.target.value)}
                      placeholder="https://x.com/username"
                      className="text-sm"
                    />
                  </div>
                </div>

                <Button size="sm" className="w-full" onClick={handleSaveLinks} disabled={savingLinks}>
                  {savingLinks ? <SpinnerGap className="size-4 animate-spin" /> : "Save Links"}
                </Button>
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>
      </div>
    </div>
  );
}
