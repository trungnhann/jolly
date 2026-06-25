"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import * as React from "react";

import { AuthCard } from "@/components/auth/auth-card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { apiRequest } from "@/lib/api";
import { useUserSession } from "@/lib/session";
import { useTranslation } from "@/lib/i18n";

export default function SigninPage() {
  const router = useRouter();
  const session = useUserSession();
  const { t } = useTranslation();

  const [email, setEmail] = React.useState("");
  const [password, setPassword] = React.useState("");
  const [error, setError] = React.useState<string | null>(null);
  const [isLoading, setIsLoading] = React.useState(false);

  return (
    <main className="flex flex-1 items-center justify-center px-6 py-16">
      <AuthCard
        title={t("auth.signInTitle")}
        description={t("auth.signInDesc")}
      >
        <form
          className="grid gap-5"
          onSubmit={async (e) => {
            e.preventDefault();
            if (!email.trim() || !password) {
              setError("Email and password are required");
              return;
            }
            setError(null);
            setIsLoading(true);

            try {
              const result = await apiRequest<{
                token: string;
                user_uuid: string;
              }>("/api/users/login", {
                method: "POST",
                body: JSON.stringify({ email: email.trim(), password }),
              });
              session.signin(result.user_uuid, result.token);
              router.push("/profile");
            } catch (err) {
              setError(err instanceof Error ? err.message : t("common.error"));
            } finally {
              setIsLoading(false);
            }
          }}
        >
          <div className="grid gap-2">
            <Label htmlFor="email">{t("auth.email")}</Label>
            <Input
              id="email"
              type="email"
              autoComplete="email"
              placeholder={t("auth.emailPlaceholder")}
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
            />
          </div>

          <div className="grid gap-2">
            <Label htmlFor="password">{t("auth.password")}</Label>
            <Input
              id="password"
              type="password"
              autoComplete="current-password"
              placeholder={t("auth.passwordPlaceholder")}
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />
          </div>

          {error ? (
            <div className="rounded-[calc(var(--radius)-10px)] border border-border bg-white/6 px-4 py-3 text-sm text-white/80">
              {error}
            </div>
          ) : null}

          <div className="grid gap-3">
            <Button type="submit" disabled={isLoading}>
              {isLoading ? t("auth.signingIn") : t("auth.submitSignIn")}
            </Button>
            <Button
              asChild
              variant="secondary"
              type="button"
              disabled={isLoading}
            >
              <Link href="/register">No account yet? Register</Link>
            </Button>
          </div>
        </form>
      </AuthCard>
    </main>
  );
}
