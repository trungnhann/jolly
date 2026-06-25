"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import * as React from "react";

import { AuthCard } from "@/components/auth/auth-card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { apiRequest } from "@/lib/api";
import { useTranslation } from "@/lib/i18n";

type UserCreated = {
  user_uuid: string;
  role: string;
};

export default function RegisterPage() {
  const router = useRouter();
  const { t } = useTranslation();

  const [email, setEmail] = React.useState("");
  const [name, setName] = React.useState("");
  const [password, setPassword] = React.useState("");
  const [role, setRole] = React.useState<"customer" | "admin">("customer");
  const [isSubmitting, setIsSubmitting] = React.useState(false);
  const [error, setError] = React.useState<string | null>(null);

  return (
    <main className="flex flex-1 items-center justify-center px-6 py-16">
      <AuthCard
        title={t("auth.registerTitle")}
        description={t("auth.registerDesc")}
      >
        <form
          className="grid gap-5"
          onSubmit={async (e) => {
            e.preventDefault();
            setIsSubmitting(true);
            setError(null);
            try {
              await apiRequest<UserCreated>("/api/users", {
                method: "POST",
                body: JSON.stringify({
                  email,
                  name,
                  password,
                  role,
                }),
              });
              // Normally we could login automatically if CreateUser returns token,
              // but here we just redirect to sign in
              router.push("/signin");
            } catch (err) {
              setError(err instanceof Error ? err.message : "Register failed");
            } finally {
              setIsSubmitting(false);
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
            <Label htmlFor="name">{t("auth.name")}</Label>
            <Input
              id="name"
              type="text"
              autoComplete="name"
              placeholder={t("auth.namePlaceholder")}
              value={name}
              onChange={(e) => setName(e.target.value)}
              required
            />
          </div>

          <div className="grid gap-2">
            <Label htmlFor="password">{t("auth.password")}</Label>
            <Input
              id="password"
              type="password"
              autoComplete="new-password"
              placeholder={t("auth.passwordPlaceholder")}
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />
          </div>

          <div className="grid gap-2">
            <Label>{t("auth.role")}</Label>
            <Select
              value={role}
              onValueChange={(v) =>
                setRole(v === "admin" ? "admin" : "customer")
              }
            >
              <SelectTrigger>
                <SelectValue placeholder="Select a role" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="customer">
                  {t("auth.roleCustomer")}
                </SelectItem>
                <SelectItem value="admin">{t("auth.roleAdmin")}</SelectItem>
              </SelectContent>
            </Select>
          </div>

          {error ? (
            <div className="rounded-[calc(var(--radius)-10px)] border border-border bg-white/6 px-4 py-3 text-sm text-white/80">
              {error}
            </div>
          ) : null}

          <div className="grid gap-3">
            <Button type="submit" disabled={isSubmitting}>
              {isSubmitting ? t("auth.registering") : t("auth.submitRegister")}
            </Button>
            <Button asChild variant="secondary" type="button">
              <Link href="/signin">Already have an account? Sign in</Link>
            </Button>
          </div>
        </form>
      </AuthCard>
    </main>
  );
}
