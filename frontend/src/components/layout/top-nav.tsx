"use client";

import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import { LogOut, User } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { useUserSession } from "@/lib/session";
import { cn } from "@/lib/utils";
import { ThemeToggle } from "@/components/theme-toggle";
import { useTranslation } from "@/lib/i18n";

function shortUUID(uuid: string) {
  if (uuid.length <= 12) {
    return uuid;
  }
  return `${uuid.slice(0, 8)}…${uuid.slice(-4)}`;
}

export function TopNav() {
  const pathname = usePathname();
  const router = useRouter();
  const session = useUserSession();
  const { t } = useTranslation();

  const isActive = (href: string) => pathname === href;

  return (
    <header className="relative z-50 border-b border-border/40 bg-background/80 px-6 backdrop-blur-md">
      <div className="mx-auto flex h-14 w-full max-w-6xl items-center justify-between gap-4">
        {/* Logo */}
        <div className="flex items-center gap-8">
          <Link
            href="/"
            className="flex items-center gap-2.5 transition-opacity hover:opacity-80"
          >
            <div className="flex h-6 w-6 items-center justify-center rounded-sm bg-primary shadow-[0_0_15px_rgba(225,6,0,0.5)]">
              <span className="text-[10px] font-black italic text-white">
                {t("common.f1")}
              </span>
            </div>
            <div className="hidden sm:block">
              <div className="text-xs font-black uppercase italic tracking-widest text-foreground">
                {t("common.jollyRacing")}
              </div>
            </div>
          </Link>

          {/* Main Navigation (F1 Style) */}
          <nav className="hidden items-center gap-6 sm:flex h-full">
            {[
              { name: t("nav.dashboard") as string, href: "/" },
              { name: t("nav.profile") as string, href: "/profile", authRequired: true },
              { name: t("nav.orders") as string, href: "/orders", authRequired: true },
              { name: t("nav.standings") as string, href: "/standings" },
            ].map((item) => {
              if (item.authRequired && !session.userUUID) return null;
              const active = isActive(item.href);

              return (
                <Link
                  key={item.name}
                  href={item.href}
                  className={cn(
                    "relative flex h-14 items-center text-[11px] font-bold uppercase tracking-[0.15em] transition-colors",
                    active
                      ? "text-primary"
                      : "text-foreground/70 hover:text-foreground",
                  )}
                >
                  {item.name}
                  {active && (
                    <div className="absolute bottom-0 left-0 h-0.5 w-full bg-primary shadow-[0_-2px_10px_rgba(225,6,0,0.5)]" />
                  )}
                </Link>
              );
            })}
          </nav>
        </div>

        {/* User Actions */}
        <div className="flex items-center gap-3">
          <ThemeToggle />
          {session.userUUID ? (
            <>
              <Badge
                className="hidden border-border/50 bg-white/5 font-mono text-[10px] sm:inline-flex"
                variant="outline"
              >
                <User className="mr-1.5 h-3 w-3 text-primary" />
              </Badge>
              <Button
                variant="ghost"
                size="icon"
                className="h-8 w-8 text-foreground/70 hover:text-foreground hover:bg-foreground/10 rounded-full"
                onClick={() => {
                  session.signout();
                  router.push("/signin");
                }}
                title={t("nav.signOut")}
              >
                <LogOut className="h-4 w-4" />
                <span className="sr-only">{t("nav.signOut")}</span>
              </Button>
            </>
          ) : (
            <>
              <Link
                href="/signin"
                className="text-[11px] font-bold uppercase tracking-widest text-foreground/70 transition-colors hover:text-foreground"
              >
                {t("nav.signIn")}
              </Link>
              <Button
                asChild
                size="sm"
                className="h-8 px-4 text-[10px] uppercase tracking-wider"
              >
                <Link href="/register">{t("nav.register")}</Link>
              </Button>
            </>
          )}
        </div>
      </div>
    </header>
  );
}
