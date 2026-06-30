"use client";

import * as React from "react";
import { ShoppingCart, Check, Loader2, Tag } from "lucide-react";

type ProductVariant = {
  uuid: string;
  name: string;
  sku: string;
  priceCents: number;
  originalPriceCents: number; // For anchor discount UX
  label: string;
  shortLabel: string; // E.g., "M", "L", "STD", "PRO"
};

type Product = {
  uuid: string;
  name: string;
  description: string;
  priceCents: number;
  category: string;
  badge?: string;
  variants: ProductVariant[];
  svgType: "wheel" | "tshirt" | "gloves";
  accentColor: string; // For customized glow states
};

const MOCK_PRODUCTS: Product[] = [
  {
    uuid: "prod-1",
    name: "Carbon Sim Steering Wheel",
    description: "Autoclaved carbon fiber chassis with magnetic paddle shifters, dual-clutch configuration, and integrated shift light telemetry display.",
    priceCents: 89900,
    category: "Hardware",
    badge: "Free Shipping",
    svgType: "wheel",
    accentColor: "from-blue-500/10 to-violet-600/5",
    variants: [
      { uuid: "var-1a", name: "Standard Edition", sku: "WHEEL-CARB-STD", priceCents: 89900, originalPriceCents: 99900, label: "Standard", shortLabel: "STD" },
      { uuid: "var-1b", name: "Pro Dual-Clutch Edition", sku: "WHEEL-CARB-PRO", priceCents: 104900, originalPriceCents: 119900, label: "Pro Clutch", shortLabel: "PRO" }
    ]
  },
  {
    uuid: "prod-2",
    name: "Jolly Team Core T-Shirt",
    description: "Official team merchandise crafted from high-performance breathable organic cotton. Features classic racing cut-lines and logo graphics.",
    priceCents: 2900,
    category: "Apparel",
    badge: "15% OFF",
    svgType: "tshirt",
    accentColor: "from-[#e10600]/10 to-amber-500/5",
    variants: [
      { uuid: "var-2a", name: "Medium Size", sku: "TSHIRT-ORG-M", priceCents: 2900, originalPriceCents: 3500, label: "Size Medium", shortLabel: "M" },
      { uuid: "var-2b", name: "Large Size", sku: "TSHIRT-ORG-L", priceCents: 3200, originalPriceCents: 3800, label: "Size Large", shortLabel: "L" }
    ]
  },
  {
    uuid: "prod-3",
    name: "Titanium Track Day Gloves",
    description: "Pre-curved high-grip silicone racing gloves with integrated carbon knuckle shields and adjustable secure strap wraps.",
    priceCents: 5900,
    category: "Safety",
    badge: "Best Seller",
    svgType: "gloves",
    accentColor: "from-emerald-500/10 to-teal-600/5",
    variants: [
      { uuid: "var-3a", name: "Medium Size", sku: "GLOVES-TITAN-M", priceCents: 5900, originalPriceCents: 6900, label: "Size Medium", shortLabel: "M" },
      { uuid: "var-3b", name: "Large Size", sku: "GLOVES-TITAN-L", priceCents: 5900, originalPriceCents: 6900, label: "Size Large", shortLabel: "L" }
    ]
  }
];

export function ProductShowroom() {
  const [selectedVariants, setSelectedVariants] = React.useState<Record<string, string>>({
    "prod-1": "var-1a",
    "prod-2": "var-2a",
    "prod-3": "var-3a"
  });
  const [cartState, setCartState] = React.useState<Record<string, "idle" | "loading" | "added">>({
    "prod-1": "idle",
    "prod-2": "idle",
    "prod-3": "idle"
  });

  const handleVariantChange = (productUuid: string, variantUuid: string) => {
    setSelectedVariants(prev => ({ ...prev, [productUuid]: variantUuid }));
  };

  const handleAddToCart = (productUuid: string) => {
    setCartState(prev => ({ ...prev, [productUuid]: "loading" }));
    setTimeout(() => {
      setCartState(prev => ({ ...prev, [productUuid]: "added" }));
      setTimeout(() => {
        setCartState(prev => ({ ...prev, [productUuid]: "idle" }));
      }, 2000);
    }, 1000);
  };

  const renderProductSVG = (type: "wheel" | "tshirt" | "gloves") => {
    switch (type) {
      case "wheel":
        return (
          <svg viewBox="0 0 200 200" className="w-full h-full text-foreground/80 overflow-visible">
            <defs>
              <radialGradient id="wheelGlow" cx="50%" cy="50%" r="50%">
                <stop offset="0%" stopColor="#3b82f6" stopOpacity="0.15" />
                <stop offset="100%" stopColor="#3b82f6" stopOpacity="0" />
              </radialGradient>
            </defs>
            <circle cx="100" cy="100" r="85" fill="url(#wheelGlow)" />
            {/* Steering Wheel Rim */}
            <circle cx="100" cy="100" r="64" fill="none" stroke="currentColor" strokeWidth="12" className="text-neutral-350 dark:text-neutral-800" />
            <path d="M 47 62 A 64 64 0 0 1 153 62" fill="none" stroke="#e10600" strokeWidth="12" strokeLinecap="round" />
            {/* Central carbon hub */}
            <rect x="70" y="80" width="60" height="40" rx="8" fill="currentColor" className="text-neutral-450 dark:text-neutral-900" />
            {/* Grip handles */}
            <rect x="30" y="75" width="14" height="50" rx="6" fill="#151619" />
            <rect x="156" y="75" width="14" height="50" rx="6" fill="#151619" />
            {/* Control Buttons */}
            <circle cx="85" cy="90" r="4.5" fill="#3b82f6" className="animate-pulse" />
            <circle cx="115" cy="90" r="4.5" fill="#22c55e" />
            <circle cx="85" cy="110" r="4.5" fill="#eab308" />
            <circle cx="115" cy="110" r="4.5" fill="#a855f7" />
          </svg>
        );
      case "tshirt":
        return (
          <svg viewBox="0 0 200 200" className="w-full h-full text-foreground/80 overflow-visible">
            <defs>
              <radialGradient id="tshirtGlow" cx="50%" cy="50%" r="50%">
                <stop offset="0%" stopColor="#e10600" stopOpacity="0.1" />
                <stop offset="100%" stopColor="#e10600" stopOpacity="0" />
              </radialGradient>
            </defs>
            <circle cx="100" cy="100" r="85" fill="url(#tshirtGlow)" />
            {/* Sleeves outline */}
            <path 
              d="M 50 50 L 70 32 L 130 32 L 150 50 L 136 63 L 125 56 L 125 168 L 75 168 L 75 56 L 64 63 Z" 
              fill="currentColor" 
              className="text-neutral-100 dark:text-[#13141a] stroke-border" 
              strokeWidth="2" 
            />
            {/* Collar line */}
            <path d="M 90 32 C 90 40 110 40 110 32" fill="none" stroke="#e10600" strokeWidth="2.5" />
            {/* Team Racing stripe */}
            <path d="M 75 92 L 125 92" fill="none" stroke="#e10600" strokeWidth="5" />
            <path d="M 75 102 L 125 102" fill="none" stroke="currentColor" strokeWidth="2" className="text-neutral-300 dark:text-neutral-700" />
          </svg>
        );
      case "gloves":
        return (
          <svg viewBox="0 0 200 200" className="w-full h-full text-foreground/80 overflow-visible">
            <defs>
              <radialGradient id="glovesGlow" cx="50%" cy="50%" r="50%">
                <stop offset="0%" stopColor="#10b981" stopOpacity="0.1" />
                <stop offset="100%" stopColor="#10b981" stopOpacity="0" />
              </radialGradient>
            </defs>
            <circle cx="100" cy="100" r="85" fill="url(#glovesGlow)" />
            {/* Main Hand glove shape */}
            <path 
              d="M 70 166 L 70 85 C 70 72 80 72 80 85 L 80 60 C 80 50 90 50 90 60 L 90 52 C 90 40 100 40 100 52 L 100 57 C 100 45 110 45 110 57 L 110 90 L 124 97 C 132 101 132 113 124 116 L 115 118 L 115 166 Z" 
              fill="currentColor" 
              className="text-neutral-150 dark:text-[#17181f] stroke-border" 
              strokeWidth="2" 
            />
            {/* Wrist strap */}
            <rect x="68" y="142" width="49" height="13" rx="2" fill="#e10600" />
            {/* Knuckles */}
            <circle cx="85" cy="98" r="4" fill="#0c0d10" />
            <circle cx="95" cy="95" r="4" fill="#0c0d10" />
            <circle cx="105" cy="98" r="4" fill="#0c0d10" />
          </svg>
        );
    }
  };

  return (
    <div className="w-full max-w-5xl mt-16 relative z-10">
      
      {/* Section Header */}
      <div className="flex items-center justify-between border-b border-border pb-4 mb-8">
        <div>
          <span className="text-[10px] font-black tracking-[0.25em] text-primary uppercase font-mono flex items-center gap-1.5">
            <Tag className="h-3.5 w-3.5 text-primary" />
            Official Racing Store
          </span>
          <h2 className="text-2xl font-black italic tracking-wider text-foreground uppercase mt-1">
            Gear & Simulator Hardware
          </h2>
        </div>
      </div>

      {/* Products Grid */}
      <div className="grid gap-6 md:grid-cols-3">
        {MOCK_PRODUCTS.map((product) => {
          const activeVarUuid = selectedVariants[product.uuid];
          const activeVar = product.variants.find(v => v.uuid === activeVarUuid) || product.variants[0];
          
          const price = (activeVar.priceCents / 100).toLocaleString("en-US", { style: "currency", currency: "USD", minimumFractionDigits: 0 });
          const originalPrice = (activeVar.originalPriceCents / 100).toLocaleString("en-US", { style: "currency", currency: "USD", minimumFractionDigits: 0 });
          const cart = cartState[product.uuid];

          const isApparel = product.category === "Apparel" || product.category === "Safety";

          return (
            <div 
              key={product.uuid}
              className="group relative flex flex-col h-[490px] rounded-[24px] border border-border bg-card overflow-hidden transition-all duration-500 hover:-translate-y-1.5 hover:border-primary/40 hover:shadow-[0_15px_40px_rgba(0,0,0,0.06)] dark:hover:shadow-[0_20px_45px_rgba(225,6,0,0.04)]"
            >
              {/* Premium Glassmorphic Badge */}
              {product.badge && (
                <span className="absolute top-4 left-4 z-10 px-3 py-1 rounded-full text-[9px] font-black uppercase tracking-wider bg-black/60 dark:bg-white/10 text-white backdrop-blur-md border border-white/10">
                  {product.badge}
                </span>
              )}

              {/* Product illustration container */}
              <div className="h-[190px] w-full flex items-center justify-center p-8 bg-gradient-to-b from-muted/30 to-muted/10 dark:from-[#0d0e12] dark:to-[#08090c] border-b border-border/50 relative overflow-hidden">
                <div className="absolute inset-0 bg-[linear-gradient(rgba(128,128,128,0.01)_1px,transparent_1px),linear-gradient(90deg,rgba(128,128,128,0.01)_1px,transparent_1px)] bg-[size:10px_10px]" />
                <div className="w-28 h-28 transform group-hover:scale-110 group-hover:-translate-y-1 transition-all duration-500 ease-out">
                  {renderProductSVG(product.svgType)}
                </div>
              </div>

              {/* Content body layout container (flex-1 flex flex-col justify-between guarantees aligned footer button row) */}
              <div className="p-5 flex-1 flex flex-col justify-between">
                
                {/* Text Description & Options */}
                <div className="space-y-3.5">
                  <div>
                    <span className="text-[9px] font-black tracking-widest text-primary uppercase font-mono">{product.category}</span>
                    <h3 className="text-base font-black italic tracking-wide text-foreground uppercase mt-1 truncate group-hover:text-primary transition-colors">
                      {product.name}
                    </h3>
                    <p className="mt-1.5 text-xs text-muted-foreground leading-relaxed font-medium line-clamp-3">
                      {product.description}
                    </p>
                  </div>

                  {/* Redesigned Variant Pickers */}
                  <div className="space-y-1.5">
                    <span className="text-[8px] font-black tracking-widest text-muted-foreground uppercase font-mono block">
                      {isApparel ? "Select Size" : "Select Specification"}
                    </span>
                    
                    {isApparel ? (
                      /* Apparel: Circular Size Bubbles Selector */
                      <div className="flex items-center gap-2">
                        {product.variants.map((v) => {
                          const isSelected = v.uuid === activeVarUuid;
                          return (
                            <button
                              key={v.uuid}
                              onClick={() => handleVariantChange(product.uuid, v.uuid)}
                              className={`w-8 h-8 rounded-full border text-[10px] font-black font-mono flex items-center justify-center transition-all duration-300 cursor-pointer ${
                                isSelected
                                  ? "bg-primary border-primary text-white shadow-xs"
                                  : "bg-background border-border text-muted-foreground hover:bg-muted hover:text-foreground"
                              }`}
                              title={v.name}
                            >
                              {v.shortLabel}
                            </button>
                          );
                        })}
                      </div>
                    ) : (
                      /* Hardware/Standard Segmented Pill Switcher Selector */
                      <div className="flex bg-muted/60 p-1 rounded-xl border border-border/40 max-w-xs">
                        {product.variants.map((v) => {
                          const isSelected = v.uuid === activeVarUuid;
                          return (
                            <button
                              key={v.uuid}
                              onClick={() => handleVariantChange(product.uuid, v.uuid)}
                              className={`flex-1 px-3 py-1 rounded-lg text-[9px] font-bold font-mono tracking-wide transition-all duration-300 cursor-pointer ${
                                isSelected
                                  ? "bg-card text-foreground shadow-xs border border-border/40"
                                  : "text-muted-foreground hover:text-foreground"
                              }`}
                            >
                              {v.label}
                            </button>
                          );
                        })}
                      </div>
                    )}
                  </div>
                </div>

                {/* Foot Pricing and Action Row - Guaranteed alignment */}
                <div className="pt-3 border-t border-border/40 flex items-center justify-between">
                  <div className="flex flex-col">
                    <span className="text-[8px] font-black tracking-widest text-muted-foreground uppercase font-mono leading-none">Price</span>
                    <div className="flex items-baseline gap-1.5 mt-1">
                      <span className="text-lg font-black font-mono text-foreground leading-none">{price}</span>
                      <span className="text-[10px] font-bold font-mono text-muted-foreground line-through opacity-75">{originalPrice}</span>
                    </div>
                  </div>

                  <button
                    onClick={() => handleAddToCart(product.uuid)}
                    disabled={cart !== "idle"}
                    className={`h-9 px-4 rounded-xl font-bold text-xs uppercase tracking-wider flex items-center gap-1.5 transition-all duration-300 cursor-pointer active:scale-95 shadow-sm ${
                      cart === "added"
                        ? "bg-green-500 text-white"
                        : "bg-primary hover:bg-primary/95 text-white shadow-md shadow-primary/10"
                    }`}
                  >
                    {cart === "loading" && <Loader2 className="h-3.5 w-3.5 animate-spin" />}
                    {cart === "added" && <Check className="h-3.5 w-3.5" />}
                    {cart === "loading" ? "Processing" : cart === "added" ? "Added" : "Add To Cart"}
                  </button>
                </div>

              </div>

            </div>
          );
        })}
      </div>

    </div>
  );
}
