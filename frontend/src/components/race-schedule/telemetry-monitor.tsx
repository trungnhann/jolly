"use client";

import * as React from "react";
import { Award, Shield, Trophy, ChevronRight, Zap, Star, ShieldAlert, Cpu, Settings } from "lucide-react";
import { CountryFlag } from "./svg-flags";

export function TeamLogo({ code, className }: { code: string; className?: string }) {
  switch (code.toUpperCase()) {
    case "MCL":
      return (
        <svg viewBox="0 0 100 100" className={className} xmlns="http://www.w3.org/2000/svg">
          <path
            d="M 20 80 C 50 80 85 55 75 25 C 65 50 40 60 20 80 Z"
            fill="#FF8000"
          />
        </svg>
      );
    case "SF":
      return (
        <svg viewBox="0 0 100 100" className={className} xmlns="http://www.w3.org/2000/svg">
          <path
            d="M 25 15 L 75 15 L 75 60 C 75 80 50 90 50 90 C 50 90 25 80 25 60 Z"
            fill="#FFEB3B"
          />
          <rect x="25" y="15" width="16.6" height="6" fill="#4CAF50" />
          <rect x="41.6" y="15" width="16.6" height="6" fill="#FFF" />
          <rect x="58.2" y="15" width="16.8" height="6" fill="#F44336" />
          <path
            d="M 46 72 C 47 70 48 67 47 64 C 46 61 44 59 45 56 C 46 53 48 51 49 48 C 50 45 50 42 49 39 C 48 36 46 34 48 32 C 50 30 52 33 53 35 C 54 37 53 40 54 42 C 55 44 57 43 59 44 C 61 45 60 48 57 48 C 54 48 53 50 54 53 C 55 56 57 58 56 61 C 55 64 53 65 54 68 C 55 71 58 73 57 75 L 53 71 C 51 72 49 73 48 75 L 49 68 Z"
            fill="#000"
          />
        </svg>
      );
    case "RBR":
      return (
        <svg viewBox="0 0 100 100" className={className} xmlns="http://www.w3.org/2000/svg">
          <circle cx="50" cy="50" r="42" fill="#0A1C36" stroke="#FFCC00" strokeWidth="2.5" />
          <circle cx="50" cy="50" r="18" fill="#FFCC00" />
          <path
            d="M 20 54 C 24 51 29 48 35 48 C 41 48 46 51 52 49 C 58 47 62 42 66 38 C 70 34 74 34 76 36 C 78 38 75 42 72 46 C 69 50 63 54 58 56 C 53 58 50 62 51 66 L 47 60 C 44 61 41 63 39 66 L 40 58 C 36 58 32 59 28 61 Z"
            fill="#E10600"
          />
        </svg>
      );
    default:
      return null;
  }
}

type DriverStanding = {
  position: number;
  driverName: string;
  teamName: string;
  teamCode: string;
  points: number;
  colorCode: string;
  flagCountryCode: string;
  change: "up" | "down" | "none";
};

type TeamSpotlight = {
  id: string;
  code: string;
  fullName: string;
  chassis: string;
  powerUnit: string;
  headquarters: string;
  principal: string;
  colorCode: string;
  flagCountryCode: string;
  championships: number;
  wins: number;
  poles: number;
  attributes: {
    aero: number;
    powerUnit: number;
    chassisEfficiency: number;
    strategy: number;
  };
};

const INITIAL_STANDINGS: DriverStanding[] = [
  { position: 1, driverName: "Max Verstappen", teamName: "Red Bull Racing", teamCode: "RBR", points: 218, colorCode: "#3671C6", flagCountryCode: "NL", change: "none" },
  { position: 2, driverName: "Lando Norris", teamName: "McLaren", teamCode: "MCL", points: 185, colorCode: "#FF8000", flagCountryCode: "GB", change: "up" },
  { position: 3, driverName: "Charles Leclerc", teamName: "Ferrari", teamCode: "SF", points: 172, colorCode: "#E10600", flagCountryCode: "MC", change: "down" },
  { position: 4, driverName: "Oscar Piastri", teamName: "McLaren", teamCode: "MCL", points: 142, colorCode: "#FF8000", flagCountryCode: "AU", change: "up" },
  { position: 5, driverName: "Carlos Sainz", teamName: "Ferrari", teamCode: "SF", points: 138, colorCode: "#E10600", flagCountryCode: "ES", change: "down" }
];

const SPOTLIGHT_TEAMS: TeamSpotlight[] = [
  {
    id: "redbull",
    code: "RBR",
    fullName: "Red Bull Racing",
    chassis: "RB22 Monolith",
    powerUnit: "Honda RBPT",
    headquarters: "Milton Keynes, UK",
    principal: "Christian Horner",
    colorCode: "#3671C6",
    flagCountryCode: "AT", // Austria (licence)
    championships: 6,
    wins: 120,
    poles: 102,
    attributes: {
      aero: 97,
      powerUnit: 94,
      chassisEfficiency: 92,
      strategy: 95
    }
  },
  {
    id: "mclaren",
    code: "MCL",
    fullName: "McLaren Racing",
    chassis: "MCL38 Papaya",
    powerUnit: "Mercedes",
    headquarters: "Woking, UK",
    principal: "Andrea Stella",
    colorCode: "#FF8000",
    flagCountryCode: "GB", // Great Britain
    championships: 8,
    wins: 186,
    poles: 160,
    attributes: {
      aero: 94,
      powerUnit: 98,
      chassisEfficiency: 95,
      strategy: 92
    }
  },
  {
    id: "ferrari",
    code: "SF",
    fullName: "Scuderia Ferrari HP",
    chassis: "SF-26 Scuderia",
    powerUnit: "Ferrari",
    headquarters: "Maranello, Italy",
    principal: "Frédéric Vasseur",
    colorCode: "#E10600",
    flagCountryCode: "IT", // Italy
    championships: 16,
    wins: 244,
    poles: 249,
    attributes: {
      aero: 92,
      powerUnit: 96,
      chassisEfficiency: 91,
      strategy: 88
    }
  }
];

export function TelemetryMonitor() {
  const [mounted, setMounted] = React.useState(false);
  const [activeTeamId, setActiveTeamId] = React.useState("redbull");

  React.useEffect(() => {
    setMounted(true);
  }, []);

  if (!mounted) return null;

  const team = SPOTLIGHT_TEAMS.find(t => t.id === activeTeamId) || SPOTLIGHT_TEAMS[0];

  return (
    <div className="w-full max-w-5xl mt-16 grid gap-8 md:grid-cols-5 relative z-10">
      
      {/* Grid Profile Spotlight Card (Spans 3 Columns on desktop) */}
      <div className="md:col-span-3 bg-card border border-border rounded-[24px] overflow-hidden p-6 shadow-md dark:shadow-xl relative before:absolute before:inset-x-0 before:top-0 before:h-[2px] before:bg-gradient-to-r before:from-[#e10600] before:to-transparent flex flex-col justify-between">
        
        {/* Background grid */}
        <div className="absolute inset-0 bg-[linear-gradient(rgba(128,128,128,0.015)_1px,transparent_1px),linear-gradient(90deg,rgba(128,128,128,0.015)_1px,transparent_1px)] bg-[size:16px_16px] pointer-events-none opacity-40" />

        {/* Card Header */}
        <div className="relative z-10 flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-6">
          <div>
            <span className="text-[9px] font-black tracking-[0.3em] text-[#e10600] uppercase font-mono flex items-center gap-1.5">
              <Shield className="h-3.5 w-3.5 text-[#e10600]" />
              Team Spotlight
            </span>
            <h3 className="text-lg font-black italic tracking-wide text-foreground uppercase mt-1">
              Constructor & Chassis Insights
            </h3>
          </div>

          {/* Constructor Selector Tabs */}
          <div className="flex items-center gap-1.5 bg-muted p-1.5 rounded-xl border border-border/50">
            {SPOTLIGHT_TEAMS.map((t) => (
              <button
                key={t.id}
                onClick={() => setActiveTeamId(t.id)}
                className={`px-3 py-1.5 rounded-lg text-[9px] font-black font-mono tracking-wider uppercase transition-all duration-200 cursor-pointer flex items-center gap-1.5 ${
                  t.id === activeTeamId
                    ? "bg-card text-foreground shadow-xs border border-border/60"
                    : "text-muted-foreground hover:text-foreground"
                }`}
              >
                <div className="w-3.5 h-3.5 flex items-center justify-center">
                  <TeamLogo code={t.code} className="w-full h-full object-contain" />
                </div>
                <span>{t.code}</span>
              </button>
            ))}
          </div>
        </div>

        {/* Driver Profile Layout */}
        <div className="relative z-10 grid grid-cols-1 sm:grid-cols-2 gap-6 items-stretch flex-1">
          
          {/* Left Column: Team Core Identity & Trophies */}
          <div className="bg-muted/30 border border-border/40 dark:bg-[#101115]/50 dark:border-neutral-900 rounded-xl p-5 flex flex-col justify-between relative overflow-hidden">
            {/* Large Watermark Team Code */}
            <div className="absolute right-2 bottom-0 text-[100px] font-black italic text-neutral-200/20 dark:text-neutral-800/10 select-none leading-none tracking-tighter">
              {team.code}
            </div>

            <div className="flex justify-between items-start gap-4">
              <div>
                <div className="flex items-center gap-2">
                  <span className="text-xs font-bold text-muted-foreground font-mono uppercase tracking-wider">
                    {team.headquarters.split(",")[1].trim()}
                  </span>
                  <div className="w-5 h-3 rounded overflow-hidden border border-white/5 flex-shrink-0">
                    <CountryFlag code={team.flagCountryCode} className="w-full h-full object-cover" />
                  </div>
                </div>
                <h2 className="text-2xl font-black italic text-foreground uppercase leading-tight mt-1">
                  {team.fullName.split(" ")[0]} <br />
                  <span className="text-3xl font-extrabold tracking-wide" style={{ color: team.colorCode }}>
                    {team.fullName.split(" ").slice(1).join(" ")}
                  </span>
                </h2>
              </div>

              {/* Team Logo Badge */}
              <div className="w-14 h-14 rounded-xl bg-card border border-border/40 flex items-center justify-center p-2 shadow-xs dark:bg-black/20">
                <TeamLogo code={team.code} className="w-full h-full object-contain" />
              </div>
            </div>
              <span className="text-[10px] text-muted-foreground font-bold uppercase tracking-widest mt-2 block">
                Chassis: {team.chassis}
              </span>
              <span className="text-[9px] text-neutral-400 dark:text-neutral-500 font-bold uppercase tracking-widest block">
                Power Unit: {team.powerUnit}
              </span>

            {/* Career stats indicators */}
            <div className="mt-8 grid grid-cols-3 gap-2 border-t border-border/50 dark:border-neutral-900/60 pt-4">
              <div className="flex flex-col">
                <span className="text-[8px] font-black tracking-widest text-muted-foreground uppercase font-mono">Titles</span>
                <span className="text-lg font-black font-mono text-foreground mt-0.5">{team.championships}</span>
              </div>
              <div className="flex flex-col">
                <span className="text-[8px] font-black tracking-widest text-muted-foreground uppercase font-mono">Wins</span>
                <span className="text-lg font-black font-mono text-foreground mt-0.5">{team.wins}</span>
              </div>
              <div className="flex flex-col">
                <span className="text-[8px] font-black tracking-widest text-muted-foreground uppercase font-mono">Poles</span>
                <span className="text-lg font-black font-mono text-foreground mt-0.5">{team.poles}</span>
              </div>
            </div>
          </div>

          {/* Right Column: Driver Stats Attributes */}
          <div className="bg-muted/30 border border-border/40 dark:bg-[#101115]/50 dark:border-neutral-900 rounded-xl p-5 flex flex-col justify-between">
            <div>
              <div className="flex justify-between items-center mb-4">
                <span className="text-[10px] font-black tracking-wider text-muted-foreground uppercase font-mono">Chassis telemetry</span>
                <div className="flex items-center gap-1 bg-[#e10600]/10 text-[#e10600] px-2.5 py-0.5 rounded-full text-xs font-black font-mono">
                  <span>SPEC</span>
                  <span>2026</span>
                </div>
              </div>

              {/* Progress attributes list */}
              <div className="space-y-3.5">
                {/* Aerodynamics */}
                <div className="space-y-1">
                  <div className="flex justify-between text-[10px] font-mono font-bold tracking-wide uppercase text-foreground/80">
                    <span>Aerodynamics</span>
                    <span>{team.attributes.aero}</span>
                  </div>
                  <div className="h-1.5 w-full bg-muted dark:bg-neutral-900 rounded-full overflow-hidden">
                    <div 
                      className="h-full rounded-full transition-all duration-500 shadow-sm"
                      style={{ width: `${team.attributes.aero}%`, backgroundColor: team.colorCode }}
                    />
                  </div>
                </div>

                {/* Power Unit Performance */}
                <div className="space-y-1">
                  <div className="flex justify-between text-[10px] font-mono font-bold tracking-wide uppercase text-foreground/80">
                    <span>Power Unit</span>
                    <span>{team.attributes.powerUnit}</span>
                  </div>
                  <div className="h-1.5 w-full bg-muted dark:bg-neutral-900 rounded-full overflow-hidden">
                    <div 
                      className="h-full rounded-full transition-all duration-500 shadow-sm"
                      style={{ width: `${team.attributes.powerUnit}%`, backgroundColor: team.colorCode }}
                    />
                  </div>
                </div>

                {/* Chassis/Weight Efficiency */}
                <div className="space-y-1">
                  <div className="flex justify-between text-[10px] font-mono font-bold tracking-wide uppercase text-foreground/80">
                    <span>Weight Efficiency</span>
                    <span>{team.attributes.chassisEfficiency}</span>
                  </div>
                  <div className="h-1.5 w-full bg-muted dark:bg-neutral-900 rounded-full overflow-hidden">
                    <div 
                      className="h-full rounded-full transition-all duration-500 shadow-sm"
                      style={{ width: `${team.attributes.chassisEfficiency}%`, backgroundColor: team.colorCode }}
                    />
                  </div>
                </div>

                {/* Strategic Execution */}
                <div className="space-y-1">
                  <div className="flex justify-between text-[10px] font-mono font-bold tracking-wide uppercase text-foreground/80">
                    <span>Pit Strategy</span>
                    <span>{team.attributes.strategy}</span>
                  </div>
                  <div className="h-1.5 w-full bg-muted dark:bg-neutral-900 rounded-full overflow-hidden">
                    <div 
                      className="h-full rounded-full transition-all duration-500 shadow-sm"
                      style={{ width: `${team.attributes.strategy}%`, backgroundColor: team.colorCode }}
                    />
                  </div>
                </div>
              </div>
            </div>

            {/* Subtle bio accent footer */}
            <div className="text-[9px] text-muted-foreground/80 italic font-medium pt-3.5 border-t border-border/40 mt-4 flex items-center justify-between">
              <span>Principal: {team.principal}</span>
              <span className="text-[8px] font-mono font-bold uppercase tracking-wider">Hq: {team.headquarters.split(",")[0]}</span>
            </div>
          </div>

        </div>

      </div>

      {/* Driver Standings Preview (Spans 2 Columns on desktop) */}
      <div className="md:col-span-2 bg-card border border-border rounded-[24px] overflow-hidden p-6 shadow-md dark:shadow-xl relative before:absolute before:inset-x-0 before:top-0 before:h-[2px] before:bg-gradient-to-r before:from-[#e10600] before:to-transparent flex flex-col justify-between">
        
        {/* Background grid */}
        <div className="absolute inset-0 bg-[linear-gradient(rgba(128,128,128,0.015)_1px,transparent_1px),linear-gradient(90deg,rgba(128,128,128,0.015)_1px,transparent_1px)] bg-[size:16px_16px] pointer-events-none opacity-40" />

        {/* Header */}
        <div className="relative z-10 flex justify-between items-center mb-6">
          <div>
            <span className="text-[9px] font-black tracking-[0.3em] text-[#e10600] uppercase font-mono flex items-center gap-1.5">
              <Trophy className="h-3.5 w-3.5 text-amber-500" />
              Standings Preview
            </span>
            <h3 className="text-lg font-black italic tracking-wide text-foreground uppercase mt-1">
              2026 Championship
            </h3>
          </div>
        </div>

        {/* Standings List */}
        <div className="relative z-10 space-y-3 flex-1 flex flex-col justify-center">
          {INITIAL_STANDINGS.map((d) => (
            <div 
              key={d.position}
              className="flex items-center justify-between p-2.5 rounded-xl border border-border/40 bg-muted/20 dark:bg-[#101115]/40 hover:bg-muted dark:hover:bg-[#101115]/80 hover:border-border dark:hover:border-neutral-800 transition-all duration-200"
            >
              <div className="flex items-center gap-3">
                <div className="flex items-center gap-2">
                  <span className="text-xs font-black font-mono text-muted-foreground w-4 text-center">
                    {d.position}
                  </span>
                  <div className="w-1.5 h-6 rounded-full" style={{ backgroundColor: d.colorCode }} />
                  {/* Team Logo Badge */}
                  <div className="w-6 h-6 rounded-md bg-card border border-border/40 flex items-center justify-center p-0.5 shadow-xs dark:bg-black/10 flex-shrink-0">
                    <TeamLogo code={d.teamCode} className="w-full h-full object-contain" />
                  </div>
                </div>

                <div>
                  <div className="flex items-center gap-1.5">
                    <span className="text-xs font-extrabold text-foreground">{d.driverName}</span>
                    <div className="w-4 h-2.5 rounded overflow-hidden border border-white/5 flex-shrink-0">
                      <CountryFlag code={d.flagCountryCode} className="w-full h-full object-cover" />
                    </div>
                  </div>
                  <span className="text-[9px] text-muted-foreground font-medium font-sans uppercase tracking-wider">{d.teamName}</span>
                </div>
              </div>

              <div className="flex items-center gap-3">
                {d.change === "up" && (
                  <span className="text-[9px] font-black text-green-500 font-mono">▲</span>
                )}
                {d.change === "down" && (
                  <span className="text-[9px] font-black text-[#e10600] font-mono">▼</span>
                )}
                {d.change === "none" && (
                  <span className="text-[9px] font-black text-neutral-400 dark:text-neutral-600 font-mono">─</span>
                )}
                <span className="text-xs font-black font-mono text-foreground">
                  {d.points} pts
                </span>
              </div>
            </div>
          ))}
        </div>

        {/* View All Standings Button */}
        <button className="relative z-10 w-full mt-6 h-9 rounded-xl border border-border hover:border-border bg-background hover:bg-muted text-foreground text-[9px] font-black tracking-widest uppercase flex items-center justify-center gap-1.5 transition-all duration-300 cursor-pointer active:scale-95 group shadow-xs">
          View Complete Standings
          <ChevronRight className="h-3.5 w-3.5 transform group-hover:translate-x-0.5 transition-transform text-[#e10600]" />
        </button>

      </div>

    </div>
  );
}
