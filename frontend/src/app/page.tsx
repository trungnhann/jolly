import Link from "next/link";
import { getTranslation } from "@/lib/i18n";
import { Calendar, Clock, ArrowRight } from "lucide-react";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { RaceCalendar } from "@/components/race-schedule/race-calendar";
import { TelemetryMonitor } from "@/components/race-schedule/telemetry-monitor";
import { ProductShowroom } from "@/components/race-schedule/product-showroom";

type Article = {
  id: string;
  title: string;
  excerpt: string;
  category: string;
  date: string;
  readTime: string;
  gradient: string;
  author: {
    name: string;
    role: string;
    initials: string;
  };
};

const SAMPLE_ARTICLES: Article[] = [
  {
    id: "1",
    title: "Unveiling the Jolly Modular Monolith",
    excerpt: "Explore how we decoupled our domain boundaries without separating databases, achieving sub-millisecond database response times and strict boundary validation.",
    category: "Architecture",
    date: "June 25, 2026",
    readTime: "5 min read",
    gradient: "from-red-500 to-violet-600",
    author: {
      name: "Alex Mercer",
      role: "Lead Architect",
      initials: "AM"
    }
  },
  {
    id: "2",
    title: "F1 Real-Time Telemetry & Data Processing",
    excerpt: "How our new telemetry engine orchestrates sub-millisecond racing simulation logs using high-performance Redis caches and optimized Postgres transaction blocks.",
    category: "Engineering",
    date: "June 23, 2026",
    readTime: "8 min read",
    gradient: "from-amber-500 to-orange-600",
    author: {
      name: "Sarah Chen",
      role: "Data Engineer",
      initials: "SC"
    }
  },
  {
    id: "3",
    title: "Next-Gen Static Assets & Local Storage",
    excerpt: "Detailing our migration to context-specific file references (like user avatars and product variant images) without polymorphic database table coupling.",
    category: "Infrastructure",
    date: "June 22, 2026",
    readTime: "4 min read",
    gradient: "from-emerald-500 to-teal-600",
    author: {
      name: "Marcus Vance",
      role: "DevOps Engineer",
      initials: "MV"
    }
  }
];

export default function Home() {
  const t = getTranslation;

  return (
    <main className="flex-1 flex flex-col items-center justify-start px-6 py-16">
      {/* F1 Hero Race Calendar Banner */}
      <RaceCalendar />

      {/* Featured Articles Section */}
      <div className="mt-16 w-full max-w-5xl relative z-10">
        <div className="flex items-center justify-between border-b border-border/40 pb-4 mb-8">
          <div>
            <p className="text-xs font-bold tracking-[0.3em] text-primary uppercase">LATEST UPDATES</p>
            <h2 className="text-2xl font-black italic tracking-wider text-foreground uppercase mt-1">Featured Articles</h2>
          </div>
          <div className="flex items-center gap-2 text-xs font-semibold text-muted-foreground hover:text-primary transition-colors cursor-pointer group">
            <span>View all updates</span>
            <ArrowRight className="h-4 w-4 transform group-hover:translate-x-1 transition-transform" />
          </div>
        </div>

        <div className="grid gap-6 md:grid-cols-3">
          {SAMPLE_ARTICLES.map((article) => (
            <div
              key={article.id}
              className="group relative flex flex-col rounded-[calc(var(--radius)-6px)] border border-border/50 bg-card/60 backdrop-blur-md overflow-hidden transition-all duration-300 hover:-translate-y-1 hover:border-primary/40 hover:shadow-[0_10px_30px_rgba(225,6,0,0.08)]"
            >
              {/* Glowing top line matching gradient */}
              <div className={`absolute top-0 left-0 h-1 w-full bg-gradient-to-r ${article.gradient} opacity-80`} />

              <div className="p-6 flex-1 flex flex-col justify-between">
                <div>
                  {/* Header: Category & Date */}
                  <div className="flex items-center justify-between gap-2 text-[10px] uppercase font-bold tracking-wider text-muted-foreground mb-4">
                    <Badge variant="outline" className="px-2 py-0.5 bg-primary/5 text-primary border-primary/10 text-[9px]">
                      {article.category}
                    </Badge>
                    <span className="flex items-center gap-1 font-mono text-[9px] text-muted-foreground/80">
                      <Calendar className="h-3 w-3" />
                      {article.date}
                    </span>
                  </div>

                  {/* Title */}
                  <h3 className="text-base font-bold text-foreground leading-snug group-hover:text-primary transition-colors duration-200 uppercase tracking-wide">
                    {article.title}
                  </h3>

                  {/* Excerpt */}
                  <p className="mt-3 text-xs text-muted-foreground font-medium line-clamp-3">
                    {article.excerpt}
                  </p>
                </div>

                {/* Footer: Author & Read time */}
                <div className="mt-6 pt-4 border-t border-border/30 flex items-center justify-between gap-4">
                  <div className="flex items-center gap-2">
                    <Avatar className="h-7 w-7 rounded-full bg-primary/10 text-primary border border-primary/20">
                      <AvatarFallback className="h-7 w-7 rounded-full text-[9px] font-bold">
                        {article.author.initials}
                      </AvatarFallback>
                    </Avatar>
                    <div>
                      <div className="text-[10px] font-bold text-foreground leading-none">
                        {article.author.name}
                      </div>
                      <div className="text-[8px] text-muted-foreground/80 leading-none mt-0.5 uppercase tracking-wider font-semibold">
                        {article.author.role}
                      </div>
                    </div>
                  </div>

                  <span className="flex items-center gap-1 text-[9px] font-mono text-muted-foreground/80 font-bold uppercase">
                    <Clock className="h-3 w-3" />
                    {article.readTime}
                  </span>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* Live Telemetry & Standings Widget Section */}
      <TelemetryMonitor />

      {/* Official Racing Gear Showroom Section */}
      <ProductShowroom />
    </main>
  );
}
