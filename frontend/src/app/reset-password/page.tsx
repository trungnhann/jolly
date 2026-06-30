"use client";

import * as React from "react";
import Link from "next/link";
import { useRouter, useSearchParams } from "next/navigation";
import { ArrowLeft, KeyRound, Loader2 } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { apiRequest } from "@/lib/api";

function ResetPasswordForm() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const token = searchParams.get("token");

  const [password, setPassword] = React.useState("");
  const [confirmPassword, setConfirmPassword] = React.useState("");
  const [isLoading, setIsLoading] = React.useState(false);
  const [error, setError] = React.useState<string | null>(null);
  const [success, setSuccess] = React.useState(false);

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!token) {
      setError("Reset token is missing or invalid. Please request a new password reset link.");
      return;
    }
    if (password.length < 8) {
      setError("Password must be at least 8 characters long.");
      return;
    }
    if (password !== confirmPassword) {
      setError("Passwords do not match.");
      return;
    }

    setIsLoading(true);
    setError(null);

    try {
      await apiRequest("/api/users/password/reset", {
        method: "POST",
        body: JSON.stringify({
          token,
          new_password: password,
        }),
      });
      setSuccess(true);
      setTimeout(() => {
        router.push("/signin");
      }, 3000);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to reset password");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Card className="w-full max-w-md bg-card/80 backdrop-blur-xl border border-border/50 shadow-lg relative overflow-hidden before:absolute before:inset-x-0 before:top-0 before:h-[2px] before:bg-gradient-to-r before:from-primary/80 before:to-transparent">
      {/* Glow backdrop */}
      <div className="absolute top-0 right-0 w-24 h-24 bg-primary/5 rounded-full blur-2xl pointer-events-none" />

      <CardHeader>
        <div className="flex items-center gap-2 mb-2">
          <Link
            href="/signin"
            className="inline-flex items-center gap-1 text-xs text-muted-foreground hover:text-primary transition-colors font-semibold uppercase tracking-wider group"
          >
            <ArrowLeft className="h-3.5 w-3.5 transform group-hover:-translate-x-0.5 transition-transform" />
            Back to Sign In
          </Link>
        </div>
        <CardTitle className="text-2xl font-black italic tracking-wide uppercase">Reset Password</CardTitle>
        <CardDescription className="text-xs text-muted-foreground mt-1">
          {!token ? "Reset token is invalid or missing." : "Enter and confirm your new password below."}
        </CardDescription>
      </CardHeader>

      <CardContent>
        {!token ? (
          <div className="space-y-4">
            <div className="rounded-[calc(var(--radius)-10px)] border border-destructive bg-destructive/10 px-4 py-3 text-xs text-destructive">
              The reset link you followed is invalid. Please generate a new link from the forgot password page.
            </div>
            <Button asChild className="w-full uppercase font-bold tracking-wider">
              <Link href="/forgot-password">Go to Forgot Password</Link>
            </Button>
          </div>
        ) : success ? (
          <div className="space-y-4 text-center">
            <div className="rounded-[calc(var(--radius)-10px)] border border-emerald-500/20 bg-emerald-500/10 px-4 py-4 text-xs text-emerald-500 font-medium">
              Your password has been successfully reset! Redirecting you to the sign-in page...
            </div>
            <Button asChild className="w-full uppercase font-bold tracking-wider">
              <Link href="/signin">Sign In Now</Link>
            </Button>
          </div>
        ) : (
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="password" className="text-xs font-bold uppercase tracking-wider text-muted-foreground">
                New Password
              </Label>
              <Input
                id="password"
                type="password"
                placeholder="••••••••"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
                disabled={isLoading}
                className="bg-background/50 border-border/50 focus:border-primary/50"
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="confirmPassword" className="text-xs font-bold uppercase tracking-wider text-muted-foreground">
                Confirm New Password
              </Label>
              <Input
                id="confirmPassword"
                type="password"
                placeholder="••••••••"
                value={confirmPassword}
                onChange={(e) => setConfirmPassword(e.target.value)}
                required
                disabled={isLoading}
                className="bg-background/50 border-border/50 focus:border-primary/50"
              />
            </div>

            {error && (
              <div className="rounded-[calc(var(--radius)-10px)] border border-destructive bg-destructive/10 px-4 py-3 text-xs text-destructive">
                {error}
              </div>
            )}

            <Button
              type="submit"
              disabled={isLoading || !password || !confirmPassword}
              className="w-full bg-gradient-to-r from-primary to-violet-600 hover:from-primary/95 hover:to-violet-600/95 text-white font-bold uppercase tracking-wider h-10 mt-2"
            >
              {isLoading ? (
                <>
                  <Loader2 className="h-4 w-4 animate-spin mr-2" />
                  Resetting Password...
                </>
              ) : (
                <>
                  <KeyRound className="h-4 w-4 mr-2" />
                  Reset Password
                </>
              )}
            </Button>
          </form>
        )}
      </CardContent>
    </Card>
  );
}

export default function ResetPasswordPage() {
  return (
    <main className="flex flex-1 items-center justify-center px-6 py-16">
      <React.Suspense fallback={
        <Card className="w-full max-w-md p-8 text-center bg-card/80 backdrop-blur-xl">
          <Loader2 className="h-8 w-8 animate-spin mx-auto text-primary" />
          <p className="text-xs text-muted-foreground mt-4 font-semibold uppercase tracking-wider">Loading reset form...</p>
        </Card>
      }>
        <ResetPasswordForm />
      </React.Suspense>
    </main>
  );
}
