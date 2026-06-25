import * as React from "react";

import { TopNav } from "./top-nav";
import { useTranslation } from "@/lib/i18n";

export function AppShell({ children }: { children: React.ReactNode }) {
  const { t } = useTranslation();

  return (
    <div className="relative flex min-h-full flex-col font-sans selection:bg-primary selection:text-white">
      <div className="pointer-events-none fixed inset-0 z-[-1] overflow-hidden">
        {/* F1/AWS Dark Glow Background */}
        <div className="absolute top-0 right-0 h-[800px] w-[800px] -translate-y-1/2 translate-x-1/3 rounded-full bg-[#E10600]/15 blur-[120px]" />
        <div className="absolute bottom-0 left-0 h-[600px] w-[600px] translate-y-1/3 -translate-x-1/3 rounded-full bg-blue-900/10 blur-[100px]" />
        {/* Tech Grid Pattern */}
        <div className="absolute inset-0 bg-[linear-gradient(to_right,#8080800a_1px,transparent_1px),linear-gradient(to_bottom,#8080800a_1px,transparent_1px)] bg-[size:24px_24px] [mask-image:radial-gradient(ellipse_60%_50%_at_50%_0%,#000_70%,transparent_100%)]" />
      </div>

      <TopNav />

      <div className="relative flex flex-1 flex-col">{children}</div>

      <footer className="relative mt-auto">
        <div className="mx-auto w-full">
          <div className="rounded-t-[32px] bg-[#0f1115] text-white px-8 py-16 md:px-16 shadow-2xl dark:bg-[#0f1115]">
            <div className="flex flex-col md:flex-row justify-between items-start gap-12">
              <div className="flex flex-col gap-6 max-w-sm">
                <div className="flex items-center gap-2.5">
                  <div className="flex h-8 w-8 items-center justify-center rounded bg-primary">
                    <span className="text-xs font-black italic text-white">
                      {t("common.f1")}
                    </span>
                  </div>
                  <div className="text-lg font-black uppercase italic tracking-widest text-white">
                    {t("common.jollyRacing")}
                  </div>
                </div>
                <p className="text-sm text-gray-400 font-medium">
                  {t("footer.description")}
                </p>
                <div className="mt-2 flex gap-4">
                  <button className="h-10 px-6 rounded-full bg-white text-black text-sm font-bold hover:bg-gray-200 transition-colors">
                    {t("footer.createAccount")}
                  </button>
                </div>
              </div>

              <div className="grid grid-cols-2 md:grid-cols-3 gap-12 w-full md:w-auto">
                <div className="flex flex-col gap-4">
                  <h4 className="text-sm font-bold text-white mb-2">
                    {t("footer.learn")}
                  </h4>
                  <a
                    href="#"
                    className="text-sm text-gray-400 hover:text-white transition-colors"
                  >
                    {t("footer.whatIsJolly")}
                  </a>
                  <a
                    href="#"
                    className="text-sm text-gray-400 hover:text-white transition-colors"
                  >
                    {t("footer.architecture")}
                  </a>
                  <a
                    href="#"
                    className="text-sm text-gray-400 hover:text-white transition-colors"
                  >
                    {t("footer.security")}
                  </a>
                  <a
                    href="#"
                    className="text-sm text-gray-400 hover:text-white transition-colors"
                  >
                    {t("footer.whatsNew")}
                  </a>
                </div>
                <div className="flex flex-col gap-4">
                  <h4 className="text-sm font-bold text-white mb-2">
                    {t("footer.resources")}
                  </h4>
                  <a
                    href="#"
                    className="text-sm text-gray-400 hover:text-white transition-colors"
                  >
                    {t("footer.gettingStarted")}
                  </a>
                  <a
                    href="#"
                    className="text-sm text-gray-400 hover:text-white transition-colors"
                  >
                    {t("footer.documentation")}
                  </a>
                  <a
                    href="#"
                    className="text-sm text-gray-400 hover:text-white transition-colors"
                  >
                    {t("footer.apiReference")}
                  </a>
                  <a
                    href="#"
                    className="text-sm text-gray-400 hover:text-white transition-colors"
                  >
                    {t("footer.partners")}
                  </a>
                </div>
                <div className="flex flex-col gap-4 col-span-2 md:col-span-1">
                  <h4 className="text-sm font-bold text-white mb-2">
                    {t("footer.help")}
                  </h4>
                  <a
                    href="#"
                    className="text-sm text-gray-400 hover:text-white transition-colors"
                  >
                    {t("footer.contactUs")}
                  </a>
                  <a
                    href="#"
                    className="text-sm text-gray-400 hover:text-white transition-colors"
                  >
                    {t("footer.supportTicket")}
                  </a>
                  <a
                    href="#"
                    className="text-sm text-gray-400 hover:text-white transition-colors"
                  >
                    {t("footer.knowledgeCenter")}
                  </a>
                </div>
              </div>
            </div>

            <div className="mt-16 pt-8 border-t border-white/10 flex flex-col md:flex-row justify-between items-center gap-4 text-xs text-gray-500">
              <p>{t("footer.equalOpportunity")}</p>
              <div className="flex gap-6">
                <a href="#" className="hover:text-white transition-colors">
                  {t("footer.privacy")}
                </a>
                <a href="#" className="hover:text-white transition-colors">
                  {t("footer.terms")}
                </a>
                <span>{t("footer.copyright")}</span>
              </div>
            </div>
          </div>
        </div>
      </footer>
    </div>
  );
}
