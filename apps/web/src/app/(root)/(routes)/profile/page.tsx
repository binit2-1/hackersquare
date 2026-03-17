"use client";

import { KeyboardEvent, useEffect, useRef, useState } from "react";
import { useAuth } from "@/context/AuthContext";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Separator } from "@/components/ui/separator";
import { Card, CardContent } from "@/components/ui/card";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Facehash } from "facehash";
import { Badge } from "@/components/ui/badge";
import { GithubCalendar } from "@/components/ui/github-calendar";
import AIThinking from "@/registry/new-york/blocks/ai/ai-thinking";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
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
  TextB,
  TextItalic,
  LinkSimple,
  Quotes,
  Code,
  ListBullets,
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
  const [profileReadme, setProfileReadme] = useState("");
  const [generatingReadme, setGeneratingReadme] = useState(false);
  const [savingReadme, setSavingReadme] = useState(false);
  const [readmeError, setReadmeError] = useState("");
  const [savedReadme, setSavedReadme] = useState("");
  const [showFullReadme, setShowFullReadme] = useState(false);
  const [isEditingReadme, setIsEditingReadme] = useState(false);
  const readmeRef = useRef<HTMLTextAreaElement | null>(null);

  useEffect(() => {
    if (!user) return;

    setHeadline(user.headline || "");
    setLocation(user.location || "");
    setWebsite(user.website_url || "");
    setLinkedin(user.linkedin_url || "");
    setTwitter(user.twitter_url || "");
    setProfileReadme(user.profileReadme || "");
    setSavedReadme(user.profileReadme || "");
    setIsEditingReadme(!(user.profileReadme || "").trim());

    const normalizedTags = Array.isArray(user.tech_tags)
      ? user.tech_tags.filter((skill): skill is string => typeof skill === "string")
      : [];
    setTechStack(normalizedTags);
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
      const res = await fetch(`/api/v1/users/profile`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({
          headline,
          location,
          website_url: website,
          linkedin_url: linkedin,
          twitter_url: twitter,
          tech_tags: techStack,
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
    window.location.href = "/api/v1/auth/github/connect";
  };

  const wrapSelectionWith = (
    prefix: string,
    suffix = "",
    placeholder = "text",
  ) => {
    const element = readmeRef.current;
    if (!element) return;

    const start = element.selectionStart || 0;
    const end = element.selectionEnd || 0;
    const selected = profileReadme.slice(start, end);
    const replacement = `${prefix}${selected || placeholder}${suffix}`;
    const updated = `${profileReadme.slice(0, start)}${replacement}${profileReadme.slice(end)}`;

    setProfileReadme(updated);

    requestAnimationFrame(() => {
      element.focus();
      const caret = start + replacement.length;
      element.setSelectionRange(caret, caret);
    });
  };

  const handleGenerateReadme = async () => {
    setReadmeError("");
    setGeneratingReadme(true);
    try {
      const response = await fetch("/api/v1/users/profile/generate-summary", {
        method: "POST",
        credentials: "include",
      });

      const payload = await response.json().catch(() => null);
      if (!response.ok) {
        throw new Error(payload?.message || "Failed to generate README");
      }

      if (typeof payload?.summary === "string" && payload.summary.trim()) {
        setProfileReadme(payload.summary);
        setSavedReadme(payload.summary);
        setShowFullReadme(false);
        setIsEditingReadme(false);
      }
      await refreshUser();
    } catch (error) {
      const message = error instanceof Error ? error.message : "Failed to generate README";
      setReadmeError(message);
    } finally {
      setGeneratingReadme(false);
    }
  };

  const handleSaveReadme = async () => {
    setReadmeError("");
    setSavingReadme(true);
    try {
      const response = await fetch("/api/v1/users/profile/readme", {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({ readme: profileReadme }),
      });

      if (!response.ok) {
        const payload = await response.json().catch(() => null);
        throw new Error(payload?.message || "Failed to save README");
      }

      setSavedReadme(profileReadme);
      setShowFullReadme(false);
      setIsEditingReadme(false);
      await refreshUser();
    } catch (error) {
      const message = error instanceof Error ? error.message : "Failed to save README";
      setReadmeError(message);
    } finally {
      setSavingReadme(false);
    }
  };



  const joinDate = new Date().toLocaleDateString("en-US", {
    day: "numeric",
    month: "short",
    year: "numeric",
  });
  const hasSavedReadme = savedReadme.trim().length > 0;
  const shouldCollapseReadme = savedReadme.length > 480;

  return (
    <div className="min-h-screen bg-background">
      <div className="max-w-2xl mx-auto px-4 py-10 sm:py-14 space-y-8">

        {/* ── Identity hero ── */}
        <div className="flex flex-col items-center text-center gap-4">
          <Facehash
            name={user.name}
            size={96}
            intensity3d="subtle"
            showInitial={true}
            variant="solid"
            enableBlink
            className="rounded-full text-black"
            colors={["#FFFFFF", "#FFFFFF", "#FFFFFF"]}
          />

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
            Skill Summary
          </h2>
          <Card>
            <CardContent className="p-4 space-y-3">
              <div className="flex flex-wrap items-center justify-between gap-2">
                <div>
                  <p className="text-sm font-medium">Developer README</p>
                  
                </div>
                {isEditingReadme ? (
                  <Button
                    type="button"
                    size="sm"
                    onClick={handleGenerateReadme}
                    disabled={generatingReadme || !user.github_handle}
                    title={user.github_handle ? "Generate README" : "Connect GitHub to generate README"}
                  >
                    <Robot className="size-4 mr-1.5" />
                    Generate README
                  </Button>
                ) : (
                  <Button
                    type="button"
                    size="sm"
                    variant="outline"
                    onClick={() => {
                      setProfileReadme(savedReadme);
                      setReadmeError("");
                      setIsEditingReadme(true);
                    }}
                  >
                    <PencilSimple className="size-4 mr-1.5" />
                    Edit README
                  </Button>
                )}
              </div>

              {isEditingReadme ? (
                <>
                  <div className="flex flex-wrap items-center gap-1 rounded-md border p-1">
                    <Button type="button" variant="ghost" size="sm" onClick={() => wrapSelectionWith("**", "**") }>
                      <TextB className="size-4" />
                    </Button>
                    <Button type="button" variant="ghost" size="sm" onClick={() => wrapSelectionWith("*", "*") }>
                      <TextItalic className="size-4" />
                    </Button>
                    <Button type="button" variant="ghost" size="sm" onClick={() => wrapSelectionWith("[", "](https://)", "link") }>
                      <LinkSimple className="size-4" />
                    </Button>
                    <Button type="button" variant="ghost" size="sm" onClick={() => wrapSelectionWith("\n> ", "", "quote") }>
                      <Quotes className="size-4" />
                    </Button>
                    <Button type="button" variant="ghost" size="sm" onClick={() => wrapSelectionWith("`", "`") }>
                      <Code className="size-4" />
                    </Button>
                    <Button type="button" variant="ghost" size="sm" onClick={() => wrapSelectionWith("\n- ", "", "item") }>
                      <ListBullets className="size-4" />
                    </Button>
                  </div>

                  {generatingReadme ? (
                    <AIThinking spinner={true} message="Generating README..." />
                  ) : (
                    <Textarea
                      ref={readmeRef}
                      value={profileReadme}
                      onChange={(e) => setProfileReadme(e.target.value)}
                      placeholder="Tell us about yourself..."
                      className="min-h-44 resize-y border-muted-foreground/20"
                    />
                  )}
                </>
              ) : hasSavedReadme ? (
                <div
                  className={`rounded-md border border-muted-foreground/20 bg-muted/20 p-4 text-sm ${
                    shouldCollapseReadme && !showFullReadme
                      ? "max-h-48 overflow-y-auto pr-2"
                      : ""
                  }`}
                >
                  <ReactMarkdown
                    remarkPlugins={[remarkGfm]}
                    components={{
                      h1: ({ ...props }) => (
                        <h1 {...props} className="mb-4 text-3xl font-semibold tracking-tight" />
                      ),
                      h2: ({ ...props }) => (
                        <h2 {...props} className="mb-3 mt-6 text-2xl font-semibold tracking-tight" />
                      ),
                      h3: ({ ...props }) => (
                        <h3 {...props} className="mb-2 mt-5 text-xl font-semibold" />
                      ),
                      p: ({ ...props }) => (
                        <p {...props} className="mb-4 leading-7 text-foreground/90" />
                      ),
                      ul: ({ ...props }) => (
                        <ul {...props} className="mb-4 list-disc space-y-2 pl-6" />
                      ),
                      ol: ({ ...props }) => (
                        <ol {...props} className="mb-4 list-decimal space-y-2 pl-6" />
                      ),
                      li: ({ ...props }) => (
                        <li {...props} className="leading-7 text-foreground/90" />
                      ),
                      blockquote: ({ ...props }) => (
                        <blockquote
                          {...props}
                          className="mb-4 border-l-2 border-muted-foreground/30 pl-3 italic text-muted-foreground"
                        />
                      ),
                      code: ({ className, children, ...props }) => {
                        const isBlock = Boolean(className && className.includes("language-"));
                        if (isBlock) {
                          return (
                            <code
                              {...props}
                              className="mb-4 block overflow-x-auto rounded bg-muted px-3 py-2 font-mono text-xs"
                            >
                              {children}
                            </code>
                          );
                        }
                        return (
                          <code
                            {...props}
                            className="rounded bg-muted px-1.5 py-0.5 font-mono text-xs"
                          >
                            {children}
                          </code>
                        );
                      },
                      strong: ({ ...props }) => (
                        <strong {...props} className="font-semibold text-foreground" />
                      ),
                      em: ({ ...props }) => <em {...props} className="italic text-foreground/90" />,
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
                    {savedReadme}
                  </ReactMarkdown>
                </div>
              ) : (
                <p className="text-xs text-muted-foreground">No README yet. Click Edit README to add your bio.</p>
              )}

              {readmeError && (
                <p className="text-xs text-red-500">{readmeError}</p>
              )}

              {!user.github_handle && (
                <p className="text-xs text-muted-foreground">
                  Connect GitHub to unlock AI README generation.
                </p>
              )}

              {isEditingReadme && (
                <>
                  <div className="flex items-center justify-between text-xs text-muted-foreground">
                    <span>Markdown supported</span>
                    <span>{profileReadme.length} chars</span>
                  </div>

                  <div className="flex justify-end gap-2">
                    <Button
                      type="button"
                      size="sm"
                      variant="ghost"
                      onClick={() => {
                        setProfileReadme(savedReadme);
                        setReadmeError("");
                        setIsEditingReadme(false);
                      }}
                    >
                      Cancel
                    </Button>
                    <Button
                      type="button"
                      size="sm"
                      onClick={handleSaveReadme}
                      disabled={generatingReadme || savingReadme}
                    >
                      {savingReadme ? <SpinnerGap className="size-4 animate-spin" /> : "Save README"}
                    </Button>
                  </div>
                </>
              )}

              {!isEditingReadme && shouldCollapseReadme && hasSavedReadme && (
                <div>
                  <Button
                    type="button"
                    variant="ghost"
                    size="sm"
                    className="h-7 px-2 text-xs"
                    onClick={() => setShowFullReadme((prev) => !prev)}
                  >
                    {showFullReadme ? "Show less" : "Show more"}
                  </Button>
                </div>
              )}
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
