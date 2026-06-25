import * as React from "react";

import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

export function AuthCard({
  title,
  description,
  children,
}: {
  title: string;
  description: string;
  children: React.ReactNode;
}) {
  return (
    <Card className="w-full max-w-md">
      <CardHeader>
        <p className="text-xs font-bold tracking-[0.3em] text-primary uppercase">
          AUTHENTICATION
        </p>
        <CardTitle className="mt-2 text-2xl font-black italic tracking-wide uppercase">
          {title}
        </CardTitle>
        <CardDescription className="mt-2 font-medium">
          {description}
        </CardDescription>
      </CardHeader>
      <CardContent className="pt-6">{children}</CardContent>
    </Card>
  );
}
