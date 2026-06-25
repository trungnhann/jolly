"use client";

import * as React from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";

import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { apiRequest } from "@/lib/api";
import { useUserSession } from "@/lib/session";
import { useTranslation } from "@/lib/i18n";

type User = {
  user_uuid: string;
  email: string;
  name: string;
  role: string;
  created_at: string;
  updated_at: string;
};

function formatDateTime(input: string) {
  const date = new Date(input);
  if (Number.isNaN(date.getTime())) {
    return input;
  }
  return new Intl.DateTimeFormat("en-US", {
    dateStyle: "medium",
    timeStyle: "short",
  }).format(date);
}

export default function ProfilePage() {
  const router = useRouter();
  const session = useUserSession();
  const { t } = useTranslation();

  const [user, setUser] = React.useState<User | null>(null);
  const [isLoading, setIsLoading] = React.useState(false);
  const [error, setError] = React.useState<string | null>(null);

  const userUUID = session.userUUID;

  const loadUser = React.useCallback(async (uuid: string) => {
    setIsLoading(true);
    setError(null);
    try {
      const u = await apiRequest<User>(`/api/users/${uuid}`, { method: "GET" });
      setUser(u);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load profile");
      setUser(null);
    } finally {
      setIsLoading(false);
    }
  }, []);

  React.useEffect(() => {
    if (!userUUID) {
      return;
    }

    let cancelled = false;
    void (async () => {
      await loadUser(userUUID);
      if (cancelled) return;
    })();
    return () => {
      cancelled = true;
    };
  }, [loadUser, userUUID]);

  return (
    <main className="flex flex-1 items-center justify-center px-6 py-16">
      <Card className="w-full max-w-2xl">
        <CardHeader className="flex flex-row items-center justify-between gap-4">
          <div>
            <p className="text-xs font-semibold tracking-[0.2em] text-muted-foreground uppercase">
              {t("profile.title")}
            </p>
            <CardTitle className="mt-2">{t("profile.description")}</CardTitle>
          </div>
          <div className="flex items-center gap-2">
            <Button
              variant="outline"
              onClick={() => {
                session.signout();
                router.push("/signin");
              }}
            >
              Sign out
            </Button>
          </div>
        </CardHeader>

        <CardContent className="grid gap-6">
          {!userUUID ? (
            <div className="rounded-[calc(var(--radius)-10px)] border border-border bg-background/35 px-5 py-4">
              <div className="text-sm font-semibold">
                {t("profile.notSignedIn")}
              </div>
              <div className="mt-1 text-xs text-muted-foreground">
                You must be signed in to view your profile.
              </div>
              <Button asChild className="mt-4">
                <Link href="/signin">{t("auth.goToSignIn")}</Link>
              </Button>
            </div>
          ) : (
            <>
              <div className="flex items-center justify-between gap-4 rounded-[calc(var(--radius)-10px)] border border-border bg-background/35 px-5 py-4">
                <div className="flex items-center gap-4">
                  <Avatar>
                    <AvatarFallback>
                      {(user?.name ?? "U").trim().slice(0, 1).toUpperCase()}
                    </AvatarFallback>
                  </Avatar>
                  <div>
                    <div className="text-sm font-semibold">
                      {user?.name ??
                        (isLoading ? t("profile.loading") : "Unknown")}
                    </div>
                    <div className="text-xs text-muted-foreground">
                      {user?.email ?? "—"}
                    </div>
                  </div>
                </div>
              </div>

              {error ? (
                <div className="rounded-[calc(var(--radius)-10px)] border border-destructive bg-destructive/10 px-4 py-3 text-sm text-destructive">
                  {error}
                </div>
              ) : null}

              <div className="grid gap-3 sm:grid-cols-2">
                <div className="rounded-[calc(var(--radius)-10px)] border border-border bg-background/20 p-4">
                  <div className="text-xs font-semibold tracking-[0.2em] text-muted-foreground uppercase">
                    {t("profile.createdAt")}
                  </div>
                  <div className="mt-2 text-sm font-semibold">
                    {user?.created_at ? formatDateTime(user.created_at) : "—"}
                  </div>
                </div>
                <div className="rounded-[calc(var(--radius)-10px)] border border-border bg-background/20 p-4">
                  <div className="text-xs font-semibold tracking-[0.2em] text-muted-foreground uppercase">
                    UPDATED AT
                  </div>
                  <div className="mt-2 text-sm font-semibold">
                    {user?.updated_at ? formatDateTime(user.updated_at) : "—"}
                  </div>
                </div>
              </div>
            </>
          )}
        </CardContent>
      </Card>
    </main>
  );
}
