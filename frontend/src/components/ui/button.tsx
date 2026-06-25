"use client";

import * as React from "react";
import { Slot } from "@radix-ui/react-slot";
import { cva, type VariantProps } from "class-variance-authority";

import { cn } from "@/lib/utils";

const buttonVariants = cva(
  "inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-[calc(var(--radius)-10px)] text-sm font-semibold transition disabled:pointer-events-none disabled:opacity-50 outline-none ring-offset-background focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2",
  {
    variants: {
      variant: {
        default:
          "bg-primary text-primary-foreground shadow-[0_4px_14px_0_rgba(225,6,0,0.39)] hover:bg-[#cc0500] hover:shadow-[0_6px_20px_rgba(225,6,0,0.23)] border border-transparent hover:border-primary/20",
        secondary:
          "bg-white/5 text-foreground hover:bg-white/10 border border-white/10 shadow-sm",
        ghost: "hover:bg-white/10 text-foreground",
        outline:
          "border border-border bg-transparent hover:bg-white/5 hover:border-white/20",
        link: "underline-offset-4 hover:underline text-white/85 hover:text-white",
      },
      size: {
        default: "h-10 px-4",
        sm: "h-9 px-3",
        lg: "h-11 px-6",
        icon: "h-10 w-10",
      },
    },
    defaultVariants: {
      variant: "default",
      size: "default",
    },
  },
);

export type ButtonProps = React.ButtonHTMLAttributes<HTMLButtonElement> &
  VariantProps<typeof buttonVariants> & {
    asChild?: boolean;
  };

export const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, variant, size, asChild = false, ...props }, ref) => {
    const Comp = asChild ? Slot : "button";
    return (
      <Comp
        ref={ref}
        className={cn(buttonVariants({ variant, size, className }))}
        {...props}
      />
    );
  },
);

Button.displayName = "Button";

export { buttonVariants };
