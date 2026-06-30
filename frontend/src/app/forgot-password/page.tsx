"use client";

import * as React from "react";
import Link from "next/link";
import { ArrowLeft, Loader2, Send } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { apiRequest } from "@/lib/api";

export default function ForgotPasswordPage() {
  const [email, setEmail] = React.useState("");
  const [isLoading, setIsLoading] = React.useState(false);
  const [message, setMessage] = React.useState<string | null>(null);
  const [error, setError] = React.useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!email) return;

    setIsLoading(true);
    setError(null);
    setMessage(null);

    try {
      await apiRequest("/api/users/password/forgot", {
        method: "POST",
        body: JSON.stringify({ email }),
      });
      // Neutral discovery check message: same message for success/failure
      setMessage("If your email is registered in our system, you will receive a password reset link shortly.");
      setEmail("");
    } catch (err) {
      setError(err instanceof Error ? err.message : "An error occurred");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <main className="flex flex-1 items-center justify-center px-6 py-16">
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
          <CardTitle className="text-2xl font-black italic tracking-wide uppercase">Forgot Password</CardTitle>
          <CardDescription className="text-xs text-muted-foreground mt-1">
            Enter your email address and we'll send you a link to reset your password.
          </CardDescription>
        </CardHeader>

        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="email" className="text-xs font-bold uppercase tracking-wider text-muted-foreground">
                Email Address
              </Label>
              <Input
                id="email"
                type="email"
                placeholder="name@example.com"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
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

            {message && (
              <div className="rounded-[calc(var(--radius)-10px)] border border-emerald-500/20 bg-emerald-500/10 px-4 py-3 text-xs text-emerald-500 font-medium">
                {message}
              </div>
            )}

            <Button
              type="submit"
              disabled={isLoading || !email}
              className="w-full bg-gradient-to-r from-primary to-violet-600 hover:from-primary/95 hover:to-violet-600/95 text-white font-bold uppercase tracking-wider h-10 mt-2"
            >
              {isLoading ? (
                <>
                  <Loader2 className="h-4 w-4 animate-spin mr-2" />
                  Sending Link...
                </>
              ) : (
                <>
                  <Send className="h-4 w-4 mr-2" />
                  Send Reset Link
                </>
              )}
            </Button>
          </form>
        </CardContent>
      </Card>
    </main>
  );
}
