"use client";

import * as React from "react";

interface FlagProps extends React.SVGProps<SVGSVGElement> {
  className?: string;
}

export function MonacoFlag({ className, ...props }: FlagProps) {
  return (
    <svg
      viewBox="0 0 30 20"
      className={className}
      xmlns="http://www.w3.org/2000/svg"
      {...props}
    >
      <rect width="30" height="20" fill="#FFFFFF" />
      <rect width="30" height="10" fill="#E10600" />
    </svg>
  );
}

export function SpainFlag({ className, ...props }: FlagProps) {
  return (
    <svg
      viewBox="0 0 750 500"
      className={className}
      xmlns="http://www.w3.org/2000/svg"
      {...props}
    >
      {/* Red bands */}
      <rect width="750" height="500" fill="#C10E24" />
      {/* Yellow band */}
      <rect y="125" width="750" height="250" fill="#FABA13" />
      {/* Simplified coat of arms */}
      <g transform="translate(180, 250) scale(0.65)">
        {/* Crown top */}
        <path d="M -30 -100 Q 0 -130 30 -100 Q 50 -100 40 -70 L -40 -70 Q -50 -100 -30 -100 Z" fill="#D32F2F" />
        <rect x="-40" y="-70" width="80" height="10" rx="3" fill="#FBC02D" />
        {/* Shield */}
        <path
          d="M -35 -50 L 35 -50 V -10 C 35 25 -35 25 -35 -10 Z"
          fill="#C10E24"
          stroke="#FBC02D"
          strokeWidth="6"
        />
        {/* Shield division line */}
        <line x1="0" y1="-50" x2="0" y2="10" stroke="#FBC02D" strokeWidth="6" />
        {/* Shield content colors */}
        <rect x="-29" y="-44" width="26" height="26" fill="#FBC02D" />
        <rect x="3" y="-18" width="26" height="26" fill="#FBC02D" />
        {/* Columns on sides */}
        <rect x="-65" y="-60" width="10" height="80" rx="2" fill="#E0E0E0" />
        <rect x="55" y="-60" width="10" height="80" rx="2" fill="#E0E0E0" />
      </g>
    </svg>
  );
}

export function AustriaFlag({ className, ...props }: FlagProps) {
  return (
    <svg
      viewBox="0 0 30 20"
      className={className}
      xmlns="http://www.w3.org/2000/svg"
      {...props}
    >
      {/* Red background */}
      <rect width="30" height="20" fill="#E10600" />
      {/* White middle band */}
      <rect y="6.67" width="30" height="6.67" fill="#FFFFFF" />
    </svg>
  );
}

export function UKFlag({ className, ...props }: FlagProps) {
  return (
    <svg
      viewBox="0 0 60 30"
      className={className}
      xmlns="http://www.w3.org/2000/svg"
      {...props}
    >
      {/* Blue field */}
      <rect width="60" height="30" fill="#012169" />
      {/* White diagonals */}
      <line x1="0" y1="0" x2="60" y2="30" stroke="#FFFFFF" strokeWidth="6" />
      <line x1="60" y1="0" x2="0" y2="30" stroke="#FFFFFF" strokeWidth="6" />
      {/* Red diagonals */}
      <line x1="0" y1="0" x2="60" y2="30" stroke="#C8102E" strokeWidth="2" />
      <line x1="60" y1="0" x2="0" y2="30" stroke="#C8102E" strokeWidth="2" />
      {/* White cross */}
      <line x1="30" y1="0" x2="30" y2="30" stroke="#FFFFFF" strokeWidth="10" />
      <line x1="0" y1="15" x2="60" y2="15" stroke="#FFFFFF" strokeWidth="10" />
      {/* Red cross */}
      <line x1="30" y1="0" x2="30" y2="30" stroke="#C8102E" strokeWidth="6" />
      <line x1="0" y1="15" x2="60" y2="15" stroke="#C8102E" strokeWidth="6" />
    </svg>
  );
}

export function BelgiumFlag({ className, ...props }: FlagProps) {
  return (
    <svg
      viewBox="0 0 15 10"
      className={className}
      xmlns="http://www.w3.org/2000/svg"
      {...props}
    >
      <rect width="5" height="10" fill="#000000" />
      <rect x="5" width="5" height="10" fill="#FFE300" />
      <rect x="10" width="5" height="10" fill="#FF0000" />
    </svg>
  );
}

export function CountryFlag({
  code,
  className,
  ...props
}: { code: string; className?: string } & React.SVGProps<SVGSVGElement>) {
  switch (code.toUpperCase()) {
    case "MC":
      return <MonacoFlag className={className} {...props} />;
    case "ES":
      return <SpainFlag className={className} {...props} />;
    case "AT":
      return <AustriaFlag className={className} {...props} />;
    case "GB":
      return <UKFlag className={className} {...props} />;
    case "BE":
      return <BelgiumFlag className={className} {...props} />;
    default:
      return (
        <div className={`bg-neutral-800 rounded border border-neutral-700 flex items-center justify-center text-[10px] text-neutral-400 font-mono ${className}`}>
          {code}
        </div>
      );
  }
}
