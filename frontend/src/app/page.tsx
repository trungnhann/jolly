import Link from "next/link";
import { getTranslation } from "@/lib/i18n";

export default function Home() {
  const t = getTranslation;

  return (
    <main className="flex flex-1 items-center justify-center px-6 py-16">
      <div className="w-full max-w-3xl rounded-[var(--radius)] border border-border/50 bg-card/80 p-8 shadow-lg backdrop-blur-xl relative overflow-hidden before:absolute before:inset-x-0 before:top-0 before:h-[2px] before:bg-gradient-to-r before:from-primary/80 before:to-transparent">
        <div className="flex items-center justify-between gap-6 relative z-10">
          <div>
            <p className="text-xs font-bold tracking-[0.3em] text-primary uppercase">
              {t("home.architecture")}
            </p>
            <h1 className="mt-3 text-4xl font-black italic tracking-wider text-foreground uppercase">
              {t("home.platform")}
            </h1>
            <p className="mt-3 text-sm text-muted-foreground font-medium">
              {t("home.description")}
            </p>
          </div>
        </div>

        <div className="mt-10 grid gap-4 sm:grid-cols-2 relative z-10">
          <Link
            href="/register"
            className="group relative overflow-hidden rounded-[calc(var(--radius)-6px)] border border-border/50 bg-background/50 p-6 transition-all hover:border-primary/50 hover:bg-background/80 hover:shadow-[0_0_20px_rgba(225,6,0,0.15)]"
          >
            <div className="absolute top-0 left-0 h-1 w-full bg-gradient-to-r from-primary to-transparent opacity-0 transition-opacity group-hover:opacity-100" />
            <div className="text-base font-bold uppercase italic tracking-wide">
              {t("home.registerTitle")}
            </div>
            <div className="mt-1 text-xs text-muted-foreground font-medium">
              {t("home.registerDesc")}
            </div>
            <div className="mt-5 h-0.5 w-8 bg-border transition-all group-hover:w-16 group-hover:bg-primary" />
          </Link>
          <Link
            href="/signin"
            className="group relative overflow-hidden rounded-[calc(var(--radius)-6px)] border border-border/50 bg-background/50 p-6 transition-all hover:border-primary/50 hover:bg-background/80 hover:shadow-[0_0_20px_rgba(225,6,0,0.15)]"
          >
            <div className="absolute top-0 left-0 h-1 w-full bg-gradient-to-r from-primary to-transparent opacity-0 transition-opacity group-hover:opacity-100" />
            <div className="text-base font-bold uppercase italic tracking-wide">
              {t("home.signInTitle")}
            </div>
            <div className="mt-1 text-xs text-muted-foreground font-medium">
              {t("home.signInDesc")}
            </div>
            <div className="mt-5 h-0.5 w-8 bg-border transition-all group-hover:w-16 group-hover:bg-primary" />
          </Link>
        </div>
      </div>
    </main>
  );
}
