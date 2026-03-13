"use client";

import { KeyboardEvent, useEffect, useState } from "react";
import { useAuth } from "@/context/AuthContext";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Separator } from "@/components/ui/separator";
import { Card, CardContent } from "@/components/ui/card";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { GithubCalendar } from "@/components/ui/github-calendar";
import {
  User,
  PencilSimple,
  MapPin,
  CalendarDots,
  GithubLogo,
  LinkedinLogo,
  XLogo,
  Globe,
  Robot,
  X,
  SpinnerGap,
} from "@phosphor-icons/react/dist/ssr";

const TECH_STACK_STORAGE_PREFIX = "profile-tech-stack:";

export default function ProfilePage() {
  const { user, isLoading, refreshUser } = useAuth();
  const [isEditing, setIsEditing] = useState(false);
  const [saving, setSaving] = useState(false);

  const [headline, setHeadline] = useState("");
  const [location, setLocation] = useState("");
  const [website, setWebsite] = useState("");
  const [linkedin, setLinkedin] = useState("");
  const [twitter, setTwitter] = useState("");
  const [techStack, setTechStack] = useState<string[]>([]);
  const [newSkill, setNewSkill] = useState("");

  useEffect(() => {
    if (!user) return;

    setHeadline(user.headline || "");
    setLocation(user.location || "");
    setWebsite(user.website_url || "");
    setLinkedin(user.linkedin_url || "");
    setTwitter(user.twitter_url || "");

    try {
      const storedSkills = localStorage.getItem(`${TECH_STACK_STORAGE_PREFIX}${user.id}`);
      if (!storedSkills) {
        setTechStack([]);
        return;
      }

      const parsedSkills = JSON.parse(storedSkills);
      if (Array.isArray(parsedSkills)) {
        setTechStack(parsedSkills.filter((skill): skill is string => typeof skill === "string"));
      }
    } catch {
      setTechStack([]);
    }
  }, [user]);

  useEffect(() => {
    if (!user) return;
    localStorage.setItem(`${TECH_STACK_STORAGE_PREFIX}${user.id}`, JSON.stringify(techStack));
  }, [techStack, user]);

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
			Your profile and connections live here.
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

  const handleAddSkill = () => {
    const normalizedSkill = newSkill.trim();
    if (!normalizedSkill) return;

    const alreadyExists = techStack.some(
      (skill) => skill.toLowerCase() === normalizedSkill.toLowerCase(),
    );
    if (alreadyExists) return;

    setTechStack((prev) => [...prev, normalizedSkill]);
    setNewSkill("");
  };

  const handleRemoveSkill = (skillToRemove: string) => {
    setTechStack((prev) => prev.filter((skill) => skill !== skillToRemove));
  };

  const handleSkillInputKeyDown = (event: KeyboardEvent<HTMLInputElement>) => {
    if (event.key !== "Enter") return;
    event.preventDefault();
    handleAddSkill();
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

  const joinDate = new Date().toLocaleDateString("en-US", {
    day: "numeric",
    month: "short",
    year: "numeric",
  });

  return (
    <div className="min-h-screen bg-background">
      <div className="max-w-2xl mx-auto px-4 py-10 sm:py-14 space-y-8">

        {/* ── Identity hero ── */}
        <div className="flex flex-col items-center text-center gap-4">
          <Avatar className="size-24 sm:size-28 text-2xl sm:text-3xl font-bold border-2 border-border shadow-sm">
            <AvatarFallback className="bg-muted text-foreground">
              {initials}
            </AvatarFallback>
          </Avatar>

          <div className="space-y-1.5">
            <h1 className="text-2xl sm:text-3xl font-bold tracking-tight">
              {user.name}
            </h1>
            <p className="text-sm text-muted-foreground max-w-md">
              {headline || "No headline set"}
            </p>
            <div className="flex items-center justify-center gap-4 text-xs text-muted-foreground pt-1">
              <span className="flex items-center gap-1">
                <CalendarDots className="size-3.5" />
                Joined {joinDate}
              </span>
              <span className="flex items-center gap-1">
                <MapPin className="size-3.5" />
                {location || "No location"}
              </span>
            </div>
          </div>

          {/* Social icons */}
          <div className="flex items-center gap-1.5">
            {user.github_handle && (
              <a href={`https://github.com/${user.github_handle}`} target="_blank" rel="noopener noreferrer">
                <Button variant="ghost" size="icon" className="size-8 rounded-full">
                  <GithubLogo className="size-4" />
                </Button>
              </a>
            )}
            {linkedin && (
              <a href={linkedin} target="_blank" rel="noopener noreferrer">
                <Button variant="ghost" size="icon" className="size-8 rounded-full">
                  <LinkedinLogo className="size-4" />
                </Button>
              </a>
            )}
            {twitter && (
              <a href={twitter} target="_blank" rel="noopener noreferrer">
                <Button variant="ghost" size="icon" className="size-8 rounded-full">
                  <XLogo className="size-4" />
                </Button>
              </a>
            )}
            {website && (
              <a href={website} target="_blank" rel="noopener noreferrer">
                <Button variant="ghost" size="icon" className="size-8 rounded-full">
                  <Globe className="size-4" />
                </Button>
              </a>
            )}
          </div>

          {/* Skill badges */}
          <div className="flex flex-wrap items-center justify-center gap-1.5 max-w-md">
            {techStack.map((skill) => (
              <Badge key={skill} variant="secondary" className="text-xs font-normal px-2.5 py-0.5">
                {skill}
              </Badge>
            ))}
            {techStack.length === 0 && (
              <p className="text-xs text-muted-foreground">No tech stack added yet</p>
            )}
          </div>

          {/* Edit Profile button */}
          <Button
            variant="outline"
            size="sm"
            className="mt-1"
            onClick={() => setIsEditing(!isEditing)}
          >
            <PencilSimple className="size-3.5 mr-1.5" />
            Edit Profile
          </Button>
        </div>

        {/* ── Inline edit panel ── */}
        {isEditing && (
          <Card>
            <CardContent className="p-5 space-y-4">
              <div className="flex items-center justify-between">
                <h3 className="text-sm font-semibold">Edit Profile</h3>
                <Button
                  variant="ghost"
                  size="icon"
                  className="size-7"
                  onClick={() => setIsEditing(false)}
                >
                  <X className="size-4" />
                </Button>
              </div>

              <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                <div className="space-y-1.5">
                  <label className="text-xs font-medium text-muted-foreground">Headline</label>
                  <Input
                    value={headline}
                    onChange={(e) => setHeadline(e.target.value)}
                    placeholder="Frontend Developer & Builder"
                    className="text-sm"
                  />
                </div>
                <div className="space-y-1.5">
                  <label className="text-xs font-medium text-muted-foreground">Location</label>
                  <Input
                    value={location}
                    onChange={(e) => setLocation(e.target.value)}
                    placeholder="Bengaluru, India"
                    className="text-sm"
                  />
                </div>
              </div>

              <Separator />

              <div className="space-y-3">
                <p className="text-xs font-medium text-muted-foreground">Social Links</p>
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                  <div className="flex items-center gap-2">
                    <Globe className="size-4 text-muted-foreground shrink-0" />
                    <Input
                      value={website}
                      onChange={(e) => setWebsite(e.target.value)}
                      placeholder="https://yourwebsite.com"
                      className="text-sm"
                    />
                  </div>
                  <div className="flex items-center gap-2">
                    <LinkedinLogo className="size-4 text-muted-foreground shrink-0" />
                    <Input
                      value={linkedin}
                      onChange={(e) => setLinkedin(e.target.value)}
                      placeholder="https://linkedin.com/in/you"
                      className="text-sm"
                    />
                  </div>
                  <div className="flex items-center gap-2">
                    <XLogo className="size-4 text-muted-foreground shrink-0" />
                    <Input
                      value={twitter}
                      onChange={(e) => setTwitter(e.target.value)}
                      placeholder="https://x.com/you"
                      className="text-sm"
                    />
                  </div>
                </div>
              </div>

              <Separator />

              <div className="space-y-3">
                <p className="text-xs font-medium text-muted-foreground">Tech Stack</p>
                <div className="flex items-center gap-2">
                  <Input
                    value={newSkill}
                    onChange={(e) => setNewSkill(e.target.value)}
                    onKeyDown={handleSkillInputKeyDown}
                    placeholder="Add a skill (e.g. Next.js)"
                    className="text-sm"
                  />
                  <Button type="button" size="sm" variant="secondary" onClick={handleAddSkill}>
                    Add
                  </Button>
                </div>

                <div className="flex flex-wrap gap-1.5">
                  {techStack.map((skill) => (
                    <Badge key={skill} variant="secondary" className="text-xs font-normal px-2 py-0.5">
                      <span>{skill}</span>
                      <button
                        type="button"
                        onClick={() => handleRemoveSkill(skill)}
                        className="ml-1.5 inline-flex items-center"
                        aria-label={`Remove ${skill}`}
                      >
                        <X className="size-3" />
                      </button>
                    </Badge>
                  ))}
                  {techStack.length === 0 && (
                    <p className="text-xs text-muted-foreground">Add a few technologies you work with.</p>
                  )}
                </div>
              </div>

              <div className="flex items-center gap-2 pt-1">
                <Button size="sm" onClick={handleSave} disabled={saving} className="min-w-24">
                  {saving ? <SpinnerGap className="size-4 animate-spin" /> : "Save Changes"}
                </Button>
                <Button size="sm" variant="ghost" onClick={() => setIsEditing(false)}>
                  Cancel
                </Button>
              </div>
            </CardContent>
          </Card>
        )}

        <Separator />

        {/* ── AI Skill Summary ── */}
        <section className="space-y-3">
          <h2 className="text-xs font-semibold uppercase tracking-wider text-muted-foreground flex items-center gap-2">
            <Robot className="size-3.5" />
            AI Skill Summary
          </h2>
          <Card className="border-dashed">
            <CardContent className="p-4">
              <p className="text-sm text-muted-foreground leading-relaxed">
                Strong fundamentals in React and Next.js. Recently demonstrating
                growth in backend systems using Go and PostgreSQL. Actively
                participating in open-source tooling.
              </p>
            </CardContent>
          </Card>
        </section>

        <Separator />

        {/* GitHub graph without card container */}
        <section className="space-y-3">
          <h2 className="text-xs font-semibold uppercase tracking-wider text-muted-foreground flex items-center gap-2">
            <GithubLogo className="size-3.5" />
            GitHub
          </h2>

          {user.github_handle ? (
            <div className="overflow-x-auto py-1">
              <GithubCalendar
                username={user.github_handle}
                colorSchema="green"
                shape="rounded"
                showTotal
              />
            </div>
          ) : (
            <div className="space-y-2">
              <p className="text-xs text-muted-foreground">
                Connect GitHub to display your contribution graph.
              </p>
              <Button size="sm" onClick={handleConnectGitHub}>
                <GithubLogo className="size-4 mr-1.5" />
                Connect GitHub
              </Button>
            </div>
          )}
        </section>
      </div>
    </div>
  );
}
