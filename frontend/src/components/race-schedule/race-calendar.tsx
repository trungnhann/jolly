"use client";

import * as React from "react";
import { ChevronLeft, ChevronRight, Plus, X, Calendar as CalendarIcon, Clock, MapPin, Trophy } from "lucide-react";
import { CountryFlag } from "./svg-flags";
import { createPortal } from "react-dom";

type Session = {
  name: string;
  date: string;
  time: string;
};

type Race = {
  id: string;
  name: string;
  fullName: string;
  circuitName: string;
  location: string;
  dateRange: string;
  flagCountryCode: string;
  countdownTarget: string;
  sessions: Session[];
};

const MOCK_RACES: Race[] = [
  {
    id: "monaco",
    name: "MONACO",
    fullName: "MONACO GRAND PRIX",
    circuitName: "Circuit de Monaco",
    location: "Monte Carlo, Monaco",
    dateRange: "24 MAY 2026",
    flagCountryCode: "MC",
    countdownTarget: "2026-05-24T15:00:00Z",
    sessions: [
      { name: "Practice 1", date: "22 May 2026", time: "18:30 - 19:30" },
      { name: "Practice 2", date: "22 May 2026", time: "22:00 - 23:00" },
      { name: "Practice 3", date: "23 May 2026", time: "17:30 - 18:30" },
      { name: "Qualifying", date: "23 May 2026", time: "21:00 - 22:00" },
      { name: "Grand Prix", date: "24 May 2026", time: "20:00 - 22:00" }
    ]
  },
  {
    id: "barcelona",
    name: "BARCELONA",
    fullName: "SPANISH GRAND PRIX",
    circuitName: "Circuit de Barcelona-Catalunya",
    location: "Montmeló, Spain",
    dateRange: "14 JUN 2026",
    flagCountryCode: "ES",
    countdownTarget: "2026-06-14T15:00:00Z",
    sessions: [
      { name: "Practice 1", date: "12 Jun 2026", time: "18:30 - 19:30" },
      { name: "Practice 2", date: "12 Jun 2026", time: "22:00 - 23:00" },
      { name: "Practice 3", date: "13 Jun 2026", time: "17:30 - 18:30" },
      { name: "Qualifying", date: "13 Jun 2026", time: "21:00 - 22:00" },
      { name: "Grand Prix", date: "14 Jun 2026", time: "20:00 - 22:00" }
    ]
  },
  {
    id: "austria",
    name: "AUSTRIAN GP",
    fullName: "AUSTRIAN GRAND PRIX",
    circuitName: "Red Bull Ring",
    location: "Spielberg, Austria",
    dateRange: "28 JUN 2026",
    flagCountryCode: "AT",
    countdownTarget: new Date(Date.now() + (2 * 24 * 60 * 60 * 1000) + (6 * 60 * 60 * 1000) + (39 * 60 * 1000) + (16 * 1000)).toISOString(),
    sessions: [
      { name: "Practice 1", date: "26 Jun 2026", time: "18:30 - 19:30" },
      { name: "Practice 2", date: "26 Jun 2026", time: "22:00 - 23:00" },
      { name: "Practice 3", date: "27 Jun 2026", time: "17:30 - 18:30" },
      { name: "Qualifying", date: "27 Jun 2026", time: "21:00 - 22:00" },
      { name: "Grand Prix", date: "28 Jun 2026", time: "20:00 - 22:00" }
    ]
  },
  {
    id: "british",
    name: "BRITISH",
    fullName: "BRITISH GRAND PRIX",
    circuitName: "Silverstone Circuit",
    location: "Silverstone, Great Britain",
    dateRange: "05 JUL 2026",
    flagCountryCode: "GB",
    countdownTarget: new Date(Date.now() + (9 * 24 * 60 * 60 * 1000)).toISOString(),
    sessions: [
      { name: "Practice 1", date: "03 Jul 2026", time: "18:30 - 19:30" },
      { name: "Practice 2", date: "03 Jul 2026", time: "22:00 - 23:00" },
      { name: "Practice 3", date: "04 Jul 2026", time: "17:30 - 18:30" },
      { name: "Qualifying", date: "04 Jul 2026", time: "21:00 - 22:00" },
      { name: "Grand Prix", date: "05 Jul 2026", time: "20:00 - 22:00" }
    ]
  },
  {
    id: "belgium",
    name: "BELGIUM",
    fullName: "BELGIAN GRAND PRIX",
    circuitName: "Circuit de Spa-Francorchamps",
    location: "Stavelot, Belgium",
    dateRange: "19 JUL 2026",
    flagCountryCode: "BE",
    countdownTarget: new Date(Date.now() + (23 * 24 * 60 * 60 * 1000)).toISOString(),
    sessions: [
      { name: "Practice 1", date: "17 Jul 2026", time: "18:30 - 19:30" },
      { name: "Practice 2", date: "17 Jul 2026", time: "22:00 - 23:00" },
      { name: "Practice 3", date: "18 Jul 2026", time: "17:30 - 18:30" },
      { name: "Qualifying", date: "18 Jul 2026", time: "21:00 - 22:00" },
      { name: "Grand Prix", date: "19 Jul 2026", time: "20:00 - 22:00" }
    ]
  }
];

function useCountdown(targetDate: string) {
  const [timeLeft, setTimeLeft] = React.useState({
    days: "00",
    hours: "00",
    minutes: "00",
    seconds: "00",
    isExpired: false
  });

  React.useEffect(() => {
    const calculateTime = () => {
      const now = Date.now();
      const target = new Date(targetDate).getTime();
      const difference = target - now;

      if (difference <= 0 || isNaN(target)) {
        setTimeLeft({ days: "00", hours: "00", minutes: "00", seconds: "00", isExpired: true });
        return;
      }

      const d = Math.floor(difference / (1000 * 60 * 60 * 24));
      const h = Math.floor((difference % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));
      const m = Math.floor((difference % (1000 * 60 * 60)) / (1000 * 60));
      const s = Math.floor((difference % (1000 * 60)) / 1000);

      setTimeLeft({
        days: String(d).padStart(2, "0"),
        hours: String(h).padStart(2, "0"),
        minutes: String(m).padStart(2, "0"),
        seconds: String(s).padStart(2, "0"),
        isExpired: false
      });
    };

    calculateTime();
    const timer = setInterval(calculateTime, 1000);

    return () => clearInterval(timer);
  }, [targetDate]);

  return timeLeft;
}

export function RaceCalendar() {
  const [activeIdx, setActiveIdx] = React.useState(2); // Austria defaults active
  const [isMobile, setIsMobile] = React.useState(false);
  const [mounted, setMounted] = React.useState(false);
  const [isModalOpen, setIsModalOpen] = React.useState(false);

  React.useEffect(() => {
    setMounted(true);
    const checkMobile = () => setIsMobile(window.innerWidth < 768);
    checkMobile();
    window.addEventListener("resize", checkMobile);
    return () => window.removeEventListener("resize", checkMobile);
  }, []);

  const activeRace = MOCK_RACES[activeIdx];
  const { days, hours, minutes, seconds, isExpired } = useCountdown(activeRace.countdownTarget);

  const prevRace = activeIdx > 0 ? MOCK_RACES[activeIdx - 1] : null;
  const nextRace = activeIdx < MOCK_RACES.length - 1 ? MOCK_RACES[activeIdx + 1] : null;

  const cardWidth = isMobile ? 280 : 360;
  const cardGap = isMobile ? 12 : 24;

  const handlePrev = () => {
    if (activeIdx > 0) {
      setActiveIdx(activeIdx - 1);
    }
  };

  const handleNext = () => {
    if (activeIdx < MOCK_RACES.length - 1) {
      setActiveIdx(activeIdx + 1);
    }
  };

  if (!mounted) {
    return (
      <div className="w-full max-w-5xl bg-card border border-border rounded-[24px] h-[220px] animate-pulse flex items-center justify-center text-muted-foreground font-mono uppercase text-xs tracking-widest">
        Loading Telemetry Engine...
      </div>
    );
  }

  return (
    <div className="w-full max-w-5xl bg-card border border-border rounded-[24px] overflow-hidden shadow-md dark:shadow-[0_30px_60px_rgba(0,0,0,0.8),0_0_50px_rgba(225,6,0,0.03)] flex flex-col select-none relative before:absolute before:inset-x-0 before:top-0 before:h-[2px] before:bg-gradient-to-r before:from-[#e10600] before:to-transparent z-10">
      
      {/* Background carbon texture grid */}
      <div className="absolute inset-0 bg-[linear-gradient(rgba(128,128,128,0.02)_1px,transparent_1px),linear-gradient(90deg,rgba(128,128,128,0.02)_1px,transparent_1px)] bg-[size:10px_10px] pointer-events-none opacity-40" />

      {/* Top Slider section */}
      <div className="relative h-44 w-full flex items-center overflow-hidden pt-4 pb-2 bg-radial from-neutral-100/50 to-neutral-200/10 dark:from-neutral-900/40 dark:to-black/80">
        
        {/* Sliding Track */}
        <div 
          className="relative left-1/2 flex flex-row flex-nowrap items-center transition-transform duration-500"
          style={{
            transform: `translateX(calc(-${activeIdx * (cardWidth + cardGap)}px - ${cardWidth / 2}px))`,
            width: `${MOCK_RACES.length * (cardWidth + cardGap)}px`,
            transitionTimingFunction: "cubic-bezier(0.25, 1, 0.5, 1)"
          }}
        >
          {MOCK_RACES.map((race, idx) => {
            const isActive = idx === activeIdx;
            const isClickable = !isActive && Math.abs(idx - activeIdx) === 1;

            return (
              <div
                key={race.id}
                onClick={() => isClickable && setActiveIdx(idx)}
                style={{
                  width: `${cardWidth}px`,
                  marginRight: `${cardGap}px`
                }}
                className={`relative rounded-[16px] overflow-hidden flex flex-col justify-between flex-shrink-0 transition-all duration-500 ${
                  isActive 
                    ? "bg-card border border-border shadow-[0_4px_20px_rgba(0,0,0,0.05)] dark:bg-[#101115] dark:border-neutral-800 dark:shadow-[0_0_25px_rgba(225,6,0,0.12),inset_0_1px_1px_rgba(255,255,255,0.05)] h-[130px] z-10 scale-100 opacity-100" 
                    : `bg-transparent border border-transparent h-[110px] scale-90 ${isClickable ? "opacity-40 dark:opacity-35 hover:opacity-75 dark:hover:opacity-60 cursor-pointer" : "opacity-15 pointer-events-none"}`
                }`}
              >
                {/* Active highlight top strip */}
                {isActive && (
                  <div className="absolute top-0 left-0 right-0 h-[1.5px] bg-gradient-to-r from-[#e10600] to-[#e10600]/20" />
                )}

                <div className="p-4 md:p-5 flex-1 flex flex-col justify-center">
                  <div className="flex items-center gap-3.5">
                    {/* Country Flag Flag Box */}
                    <div className={`rounded-[8px] overflow-hidden border border-white/10 flex-shrink-0 shadow-sm md:shadow-md transition-all duration-500 ${isActive ? "w-12 h-8 md:w-13 md:h-8.5" : "w-10 h-6.5"}`}>
                      <CountryFlag code={race.flagCountryCode} className="w-full h-full object-cover" />
                    </div>

                    {/* Race Title & Details */}
                    <div className="min-w-0 flex-1">
                      <h3 className={`font-black italic uppercase tracking-wider truncate leading-tight transition-all duration-500 ${isActive ? "text-base md:text-lg text-foreground" : "text-sm text-muted-foreground"}`}>
                        {isActive ? race.fullName : race.name}
                      </h3>
                      
                      {/* Sub-text changes dynamically */}
                      {isActive ? (
                        <div className="mt-1 flex items-center text-[10px] md:text-xs font-black tracking-wider text-muted-foreground uppercase font-mono">
                          <span className="text-[#e10600] animate-pulse">●</span>
                          <span className="ml-1.5 mr-2">START RACE:</span>
                          <span className="text-foreground font-extrabold text-[11px] md:text-xs tracking-normal">
                            {isExpired ? (
                              <span className="text-red-500 font-bold">RACE COMPLETED</span>
                            ) : (
                              <>
                                <span className="text-foreground">{days}</span><span className="text-muted-foreground font-normal lowercase text-[9px] mr-1">d</span>
                                <span className="text-foreground">{hours}</span><span className="text-muted-foreground font-normal text-[9px] mx-0.5">:</span>
                                <span className="text-foreground">{minutes}</span><span className="text-muted-foreground font-normal text-[9px] mx-0.5">:</span>
                                <span className="text-foreground">{seconds}</span>
                              </>
                            )}
                          </span>
                        </div>
                      ) : (
                        <span className="text-[10px] md:text-xs text-neutral-400 dark:text-neutral-600 font-mono font-bold mt-0.5 block">{race.dateRange}</span>
                      )}
                    </div>
                  </div>
                </div>

                {/* Session details bottom strip inside active card */}
                {isActive && (
                  <div className="bg-muted/30 border-t border-border/50 dark:bg-[#15161c] dark:border-neutral-900 py-1.5 px-4 md:px-5 flex justify-between items-center text-[9px] md:text-[10px] font-mono text-muted-foreground tracking-wider">
                    <div className="flex items-center gap-1.5">
                      <span className="text-neutral-500 font-bold uppercase">Quali.</span>
                      <span className="text-foreground/80 dark:text-neutral-300">{race.sessions[3].date.split(" ")[0] + " " + race.sessions[3].date.split(" ")[1]}</span>
                      <span className="text-neutral-400 dark:text-neutral-500 text-[8px]">{race.sessions[3].time.split(" ")[0]}</span>
                    </div>
                    <div className="w-[1px] h-2 bg-border dark:bg-neutral-800" />
                    <div className="flex items-center gap-1.5">
                      <span className="text-[#e10600] font-bold uppercase">Race</span>
                      <span className="text-foreground/80 dark:text-neutral-300">{race.sessions[4].date.split(" ")[0] + " " + race.sessions[4].date.split(" ")[1]}</span>
                      <span className="text-neutral-400 dark:text-neutral-500 text-[8px]">{race.sessions[4].time.split(" ")[0]}</span>
                    </div>
                  </div>
                )}
              </div>
            );
          })}
        </div>
      </div>

      {/* Bottom Nav section */}
      <div className="w-full bg-muted/40 border-t border-border/60 dark:bg-[#050507]/90 dark:border-neutral-900/50 h-12 flex items-center justify-between px-4 sm:px-6 relative z-10 backdrop-blur-xs">
        
        {/* Left Arrow Button */}
        <button
          onClick={handlePrev}
          disabled={!prevRace}
          className={`flex items-center gap-2 text-xs font-extrabold tracking-widest text-muted-foreground hover:text-foreground uppercase transition-all duration-300 focus:outline-hidden ${!prevRace ? "opacity-20 cursor-not-allowed" : "cursor-pointer group"}`}
        >
          <ChevronLeft className="h-4 w-4 transform group-hover:-translate-x-0.5 transition-transform text-[#e10600]" />
          {prevRace && (
            <div className="w-5 h-3.5 rounded overflow-hidden border border-white/10 flex-shrink-0 transition-transform group-hover:scale-105">
              <CountryFlag code={prevRace.flagCountryCode} className="w-full h-full object-cover" />
            </div>
          )}
        </button>

        {/* More Information Button */}
        <button 
          onClick={() => setIsModalOpen(true)}
          className="h-8 px-4 rounded-full border border-border hover:border-[#e10600]/40 dark:border-neutral-800 bg-background text-muted-foreground hover:text-foreground text-[9px] font-black tracking-widest uppercase flex items-center gap-1.5 transition-all duration-300 shadow-sm hover:shadow-[0_0_15px_rgba(225,6,0,0.06)] active:scale-95 group"
        >
          <Plus className="h-3 w-3 text-neutral-400 group-hover:text-[#e10600] transition-colors" />
          More Information
        </button>

        {/* Right Arrow Button */}
        <button
          onClick={handleNext}
          disabled={!nextRace}
          className={`flex items-center gap-2 text-xs font-extrabold tracking-widest text-muted-foreground hover:text-foreground uppercase transition-all duration-300 focus:outline-hidden ${!nextRace ? "opacity-20 cursor-not-allowed" : "cursor-pointer group"}`}
        >
          {nextRace && (
            <div className="w-5 h-3.5 rounded overflow-hidden border border-white/10 flex-shrink-0 transition-transform group-hover:scale-105">
              <CountryFlag code={nextRace.flagCountryCode} className="w-full h-full object-cover" />
            </div>
          )}
          <ChevronRight className="h-4 w-4 transform group-hover:translate-x-0.5 transition-transform text-[#e10600]" />
        </button>
      </div>

      {/* Details Modal Overlay */}
      {isModalOpen && mounted && createPortal(
        <div className="fixed inset-0 bg-black/50 dark:bg-black/85 backdrop-blur-xs dark:backdrop-blur-md z-50 flex items-center justify-center p-4 transition-all duration-300 animate-fade-in">
          <div className="w-full max-w-md bg-card border border-border rounded-[20px] overflow-hidden shadow-2xl relative animate-scale-up">
            
            {/* Red top line */}
            <div className="h-[3px] bg-[#e10600] w-full" />

            {/* Header */}
            <div className="p-6 pb-4 flex justify-between items-start">
              <div>
                <span className="text-[10px] font-black tracking-[0.25em] text-[#e10600] uppercase font-mono">Telemetry Data</span>
                <h2 className="text-xl font-black italic tracking-wide text-foreground uppercase mt-1 flex items-center gap-2">
                  {activeRace.fullName}
                </h2>
              </div>
              <button 
                onClick={() => setIsModalOpen(false)}
                className="p-1 rounded-lg border border-border hover:border-neutral-450 dark:border-neutral-800 bg-muted/40 dark:bg-neutral-900/50 text-muted-foreground hover:text-foreground transition-colors cursor-pointer"
              >
                <X className="h-4 w-4" />
              </button>
            </div>

            {/* Info body */}
            <div className="px-6 pb-6 space-y-5">
              
              {/* Location/Circuit Detail */}
              <div className="bg-muted/50 border border-border/50 dark:bg-[#101115] dark:border-neutral-900 rounded-[12px] p-4 space-y-2.5">
                <div className="flex items-center gap-2.5 text-xs text-foreground/80 dark:text-neutral-300 font-medium">
                  <MapPin className="h-4 w-4 text-[#e10600]" />
                  <span>{activeRace.location}</span>
                </div>
                <div className="flex items-center gap-2.5 text-xs text-foreground/80 dark:text-neutral-300 font-medium">
                  <Trophy className="h-4 w-4 text-amber-500" />
                  <span>{activeRace.circuitName}</span>
                </div>
              </div>

              {/* Race Schedule Timeline */}
              <div className="space-y-2">
                <h4 className="text-[10px] font-black tracking-wider text-muted-foreground uppercase font-mono">Race Weekend Schedule</h4>
                <div className="divide-y divide-border/60 dark:divide-neutral-900 border border-border/50 dark:border-neutral-900 rounded-[12px] bg-muted/20 dark:bg-[#101115]/50 overflow-hidden">
                  {activeRace.sessions.map((session, sidx) => {
                    const isMainRace = session.name === "Grand Prix";
                    const isQuali = session.name === "Qualifying";
                    return (
                      <div 
                        key={sidx} 
                        className={`p-3 flex justify-between items-center text-xs font-mono transition-colors ${
                          isMainRace ? "bg-[#e10600]/5 hover:bg-[#e10600]/10" : "hover:bg-muted/40 dark:hover:bg-white/2"
                        }`}
                      >
                        <div className="flex flex-col">
                          <span className={`font-bold uppercase tracking-wider ${isMainRace ? "text-[#e10600]" : isQuali ? "text-foreground/90 dark:text-neutral-200" : "text-muted-foreground"}`}>
                            {session.name}
                          </span>
                          <span className="text-[10px] text-muted-foreground/85 font-sans mt-0.5">{session.date}</span>
                        </div>
                        <div className="flex items-center gap-1.5 text-foreground/90 dark:text-neutral-300 font-bold">
                          <Clock className="h-3.5 w-3.5 text-muted-foreground" />
                          <span>{session.time}</span>
                        </div>
                      </div>
                    );
                  })}
                </div>
              </div>

            </div>

            {/* Modal Footer */}
            <div className="bg-muted/40 dark:bg-[#050507] px-6 py-4 border-t border-border/60 dark:border-neutral-900/60 flex justify-end">
              <button
                onClick={() => setIsModalOpen(false)}
                className="px-5 py-2 bg-muted hover:bg-neutral-200 border border-border dark:bg-neutral-900 dark:hover:bg-neutral-850 dark:border-neutral-800 rounded-lg text-foreground hover:text-black dark:text-neutral-300 dark:hover:text-white font-bold text-xs uppercase tracking-wider transition-all duration-200 cursor-pointer active:scale-95"
              >
                Close Data
              </button>
            </div>

          </div>
        </div>,
        document.body
      )}

      {/* Embed local animation frames inside global document context */}
      <style jsx global>{`
        @keyframes fadeIn {
          from { opacity: 0; }
          to { opacity: 1; }
        }
        @keyframes scaleUp {
          from { transform: scale(0.95); opacity: 0; }
          to { transform: scale(1); opacity: 1; }
        }
        .animate-fade-in {
          animation: fadeIn 0.2s ease-out forwards;
        }
        .animate-scale-up {
          animation: scaleUp 0.25s cubic-bezier(0.34, 1.56, 0.64, 1) forwards;
        }
      `}</style>
    </div>
  );
}
