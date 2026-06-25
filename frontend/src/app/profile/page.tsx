"use client";

import * as React from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";

import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { apiRequest } from "@/lib/api";
import { useUserSession } from "@/lib/session";
import { useTranslation } from "@/lib/i18n";
import { Camera, Loader2 } from "lucide-react";

type User = {
  user_uuid: string;
  email: string;
  name: string;
  role: string;
  avatar_url?: string;
  created_at: string;
  updated_at: string;
};

function formatDateTime(input: string) {
  const date = new Date(input);
  if (Number.isNaN(date.getTime())) {
    return input;
  }
  return new Intl.DateTimeFormat("en-US", {
    dateStyle: "medium",
    timeStyle: "short",
  }).format(date);
}

export default function ProfilePage() {
  const router = useRouter();
  const session = useUserSession();
  const { t } = useTranslation();

  const [user, setUser] = React.useState<User | null>(null);
  const [isLoading, setIsLoading] = React.useState(false);
  const [error, setError] = React.useState<string | null>(null);

  const userUUID = session.userUUID;

  const fileInputRef = React.useRef<HTMLInputElement>(null);
  const [isUploading, setIsUploading] = React.useState(false);
  const [uploadError, setUploadError] = React.useState<string | null>(null);

  const [selectedFile, setSelectedFile] = React.useState<File | null>(null);
  const [selectedImageSrc, setSelectedImageSrc] = React.useState<string | null>(null);
  const [zoom, setZoom] = React.useState(1);
  const [pan, setPan] = React.useState({ x: 0, y: 0 });
  const [isDragging, setIsDragging] = React.useState(false);
  const dragStart = React.useRef({ x: 0, y: 0 });
  const [imageDimensions, setImageDimensions] = React.useState({ width: 0, height: 0 });

  const handleMouseDown = (e: React.MouseEvent<HTMLDivElement>) => {
    e.preventDefault();
    setIsDragging(true);
    dragStart.current = { x: e.clientX - pan.x, y: e.clientY - pan.y };
  };

  const handleMouseMove = (e: React.MouseEvent<HTMLDivElement>) => {
    if (!isDragging) return;
    setPan({
      x: e.clientX - dragStart.current.x,
      y: e.clientY - dragStart.current.y,
    });
  };

  const handleMouseUp = () => {
    setIsDragging(false);
  };

  const handleTouchStart = (e: React.TouchEvent<HTMLDivElement>) => {
    if (e.touches.length !== 1) return;
    setIsDragging(true);
    const touch = e.touches[0];
    dragStart.current = { x: touch.clientX - pan.x, y: touch.clientY - pan.y };
  };

  const handleTouchMove = (e: React.TouchEvent<HTMLDivElement>) => {
    if (!isDragging || e.touches.length !== 1) return;
    const touch = e.touches[0];
    setPan({
      x: touch.clientX - dragStart.current.x,
      y: touch.clientY - dragStart.current.y,
    });
  };

  const handleImageLoad = (e: React.SyntheticEvent<HTMLImageElement>) => {
    const { naturalWidth, naturalHeight } = e.currentTarget;
    setImageDimensions({ width: naturalWidth, height: naturalHeight });
  };

  const getDisplayDimensions = () => {
    if (!imageDimensions.width || !imageDimensions.height) {
      return { w: 256, h: 256 };
    }
    const aspectRatio = imageDimensions.width / imageDimensions.height;
    if (aspectRatio > 1) {
      return { w: 256 * aspectRatio, h: 256 };
    } else {
      return { w: 256, h: 256 / aspectRatio };
    }
  };
  const { w: wd, h: hd } = getDisplayDimensions();

  const handleFileSelect = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (!file) return;

    if (!file.type.startsWith("image/")) {
      setUploadError("Please select a valid image file.");
      return;
    }
    if (file.size > 10 * 1024 * 1024) {
      setUploadError("Image size must be less than 10MB.");
      return;
    }

    setUploadError(null);
    setSelectedFile(file);
    const objectUrl = URL.createObjectURL(file);
    setSelectedImageSrc(objectUrl);
  };

  const handleCloseCropper = () => {
    if (selectedImageSrc) {
      URL.revokeObjectURL(selectedImageSrc);
    }
    setSelectedImageSrc(null);
    setSelectedFile(null);
    setZoom(1);
    setPan({ x: 0, y: 0 });
    setImageDimensions({ width: 0, height: 0 });
  };

  const handleCropAndUpload = async () => {
    if (!selectedImageSrc || !userUUID) return;

    setIsUploading(true);
    setUploadError(null);

    const img = new Image();
    img.crossOrigin = "anonymous";
    img.src = selectedImageSrc;
    img.onload = () => {
      const canvas = document.createElement("canvas");
      canvas.width = 256;
      canvas.height = 256;
      const ctx = canvas.getContext("2d");

      if (!ctx) {
        setUploadError("Failed to initialize canvas context.");
        setIsUploading(false);
        return;
      }

      ctx.clearRect(0, 0, 256, 256);
      ctx.translate(128 + pan.x, 128 + pan.y);
      ctx.scale(zoom, zoom);
      ctx.drawImage(img, -wd / 2, -hd / 2, wd, hd);

      canvas.toBlob(
        async (blob) => {
          if (!blob) {
            setUploadError("Failed to generate cropped image blob.");
            setIsUploading(false);
            return;
          }

          const formData = new FormData();
          const fileType = selectedFile?.type || "image/png";
          const fileExtension = fileType.split("/")[1] || "png";
          formData.append("file", blob, `avatar.${fileExtension}`);

          try {
            const updatedUser = await apiRequest<User>(`/api/users/${userUUID}/avatar`, {
              method: "POST",
              body: formData,
            });
            setUser(updatedUser);
            handleCloseCropper();
          } catch (err) {
            setUploadError(err instanceof Error ? err.message : "Failed to upload avatar");
          } finally {
            setIsUploading(false);
          }
        },
        selectedFile?.type || "image/png",
        0.95
      );
    };

    img.onerror = () => {
      setUploadError("Failed to load image for cropping.");
      setIsUploading(false);
    };
  };

  const loadUser = React.useCallback(async (uuid: string) => {
    setIsLoading(true);
    setError(null);
    try {
      const u = await apiRequest<User>(`/api/users/${uuid}`, { method: "GET" });
      setUser(u);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load profile");
      setUser(null);
    } finally {
      setIsLoading(false);
    }
  }, []);

  React.useEffect(() => {
    if (!userUUID) {
      return;
    }

    let cancelled = false;
    void (async () => {
      await loadUser(userUUID);
      if (cancelled) return;
    })();
    return () => {
      cancelled = true;
    };
  }, [loadUser, userUUID]);

  return (
    <main className="flex flex-1 items-center justify-center px-6 py-16">
      <Card className="w-full max-w-2xl">
        <CardHeader className="flex flex-row items-center justify-between gap-4">
          <div>
            <p className="text-xs font-semibold tracking-[0.2em] text-muted-foreground uppercase">
              {t("profile.title")}
            </p>
            <CardTitle className="mt-2">{t("profile.description")}</CardTitle>
          </div>
          <div className="flex items-center gap-2">
            <Button
              variant="outline"
              onClick={() => {
                session.signout();
                router.push("/signin");
              }}
            >
              Sign out
            </Button>
          </div>
        </CardHeader>

        <CardContent className="grid gap-6">
          {!userUUID ? (
            <div className="rounded-[calc(var(--radius)-10px)] border border-border bg-background/35 px-5 py-4">
              <div className="text-sm font-semibold">
                {t("profile.notSignedIn")}
              </div>
              <div className="mt-1 text-xs text-muted-foreground">
                You must be signed in to view your profile.
              </div>
              <Button asChild className="mt-4">
                <Link href="/signin">{t("auth.goToSignIn")}</Link>
              </Button>
            </div>
          ) : (
            <>
              <div className="flex flex-col items-center sm:flex-row sm:items-center justify-between gap-6 rounded-[calc(var(--radius)-10px)] border border-border bg-background/35 px-6 py-6 backdrop-blur-md relative overflow-hidden">
                <div className="absolute top-0 right-0 w-32 h-32 bg-primary/10 rounded-full blur-3xl pointer-events-none" />
                <div className="absolute bottom-0 left-0 w-32 h-32 bg-primary/5 rounded-full blur-3xl pointer-events-none" />

                <div className="flex flex-col items-center sm:flex-row gap-6 relative z-10 w-full">
                  <div className="relative group cursor-pointer" onClick={() => !isUploading && fileInputRef.current?.click()}>
                    <div className="absolute -inset-0.5 bg-gradient-to-r from-primary to-violet-600 rounded-full blur opacity-30 group-hover:opacity-75 transition duration-500" />
                    <div className="relative">
                      <Avatar className="h-24 w-24 rounded-full border-2 border-background shadow-xl">
                        {user?.avatar_url && (
                          <AvatarImage
                            src={user.avatar_url}
                            alt={user.name}
                            className="object-cover rounded-full"
                          />
                        )}
                        <AvatarFallback className="h-24 w-24 rounded-full text-lg bg-gradient-to-br from-primary/10 to-primary/20 text-primary font-bold">
                          {(user?.name ?? "U").trim().slice(0, 1).toUpperCase()}
                        </AvatarFallback>
                      </Avatar>

                      <div className="absolute inset-0 bg-black/60 rounded-full flex flex-col items-center justify-center text-white opacity-0 group-hover:opacity-100 transition-opacity duration-300">
                        <Camera className="h-6 w-6 mb-1 text-primary-foreground/90 scale-90 group-hover:scale-100 transition-transform duration-300" />
                        <span className="text-[10px] uppercase font-bold tracking-wider text-primary-foreground/80">Change</span>
                      </div>

                      {isUploading && (
                        <div className="absolute inset-0 bg-black/75 rounded-full flex items-center justify-center text-white">
                          <Loader2 className="h-8 w-8 animate-spin text-primary" />
                        </div>
                      )}
                    </div>
                  </div>

                  <input
                    type="file"
                    ref={fileInputRef}
                    onChange={handleFileSelect}
                    accept="image/*"
                    className="hidden"
                    disabled={isUploading}
                  />

                  <div className="text-center sm:text-left flex-1">
                    <div className="text-xl font-bold flex flex-wrap items-center justify-center sm:justify-start gap-2">
                      {user?.name ?? (isLoading ? t("profile.loading") : "Unknown")}
                      {user?.role && (
                        <Badge variant="outline" className="capitalize text-xs font-semibold bg-primary/5 text-primary border-primary/20">
                          {user.role}
                        </Badge>
                      )}
                    </div>
                    <div className="text-sm text-muted-foreground mt-1">
                      {user?.email ?? "—"}
                    </div>
                    <p className="text-xs text-muted-foreground/60 mt-2">
                      Click avatar to upload a new profile picture.
                    </p>
                  </div>
                </div>
              </div>

              {uploadError ? (
                <div className="rounded-[calc(var(--radius)-10px)] border border-destructive bg-destructive/10 px-4 py-3 text-sm text-destructive">
                  {uploadError}
                </div>
              ) : null}

              {error ? (
                <div className="rounded-[calc(var(--radius)-10px)] border border-destructive bg-destructive/10 px-4 py-3 text-sm text-destructive">
                  {error}
                </div>
              ) : null}

              <div className="grid gap-3 sm:grid-cols-2">
                <div className="rounded-[calc(var(--radius)-10px)] border border-border bg-background/20 p-4">
                  <div className="text-xs font-semibold tracking-[0.2em] text-muted-foreground uppercase">
                    {t("profile.createdAt")}
                  </div>
                  <div className="mt-2 text-sm font-semibold">
                    {user?.created_at ? formatDateTime(user.created_at) : "—"}
                  </div>
                </div>
                <div className="rounded-[calc(var(--radius)-10px)] border border-border bg-background/20 p-4">
                  <div className="text-xs font-semibold tracking-[0.2em] text-muted-foreground uppercase">
                    UPDATED AT
                  </div>
                  <div className="mt-2 text-sm font-semibold">
                    {user?.updated_at ? formatDateTime(user.updated_at) : "—"}
                  </div>
                </div>
              </div>
            </>
          )}
        </CardContent>
      </Card>
      {selectedImageSrc && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-md p-4 animate-in fade-in duration-200">
          <Card className="w-full max-w-md bg-background/95 border-border shadow-2xl relative overflow-hidden">
            <div className="absolute top-0 right-0 w-24 h-24 bg-primary/10 rounded-full blur-2xl pointer-events-none" />

            <CardHeader>
              <CardTitle className="text-lg font-bold">Crop Profile Picture</CardTitle>
              <p className="text-xs text-muted-foreground mt-1">
                Drag the image to adjust position and use the slider to zoom.
              </p>
            </CardHeader>

            <CardContent className="flex flex-col items-center gap-6">
              <div
                className="w-64 h-64 rounded-full border-2 border-primary/30 bg-muted overflow-hidden relative cursor-move select-none shadow-inner group"
                onMouseDown={handleMouseDown}
                onMouseMove={handleMouseMove}
                onMouseUp={handleMouseUp}
                onMouseLeave={handleMouseUp}
                onTouchStart={handleTouchStart}
                onTouchMove={handleTouchMove}
                onTouchEnd={handleMouseUp}
              >
                <div className="absolute inset-0 border-4 border-background/20 rounded-full pointer-events-none z-10" />

                <img
                  src={selectedImageSrc}
                  alt="Crop preview"
                  onLoad={handleImageLoad}
                  style={{
                    transform: `translate(${pan.x}px, ${pan.y}px) scale(${zoom})`,
                    transformOrigin: "center center",
                    width: `${wd}px`,
                    height: `${hd}px`,
                    maxWidth: "none",
                  }}
                  className="absolute pointer-events-none select-none"
                />
              </div>

              <div className="w-full space-y-2">
                <div className="flex items-center justify-between text-xs text-muted-foreground font-medium">
                  <span>Zoom</span>
                  <span>{Math.round(zoom * 100)}%</span>
                </div>
                <div className="flex items-center gap-3">
                  <span className="text-muted-foreground text-xs font-bold">A-</span>
                  <input
                    type="range"
                    min="1"
                    max="3"
                    step="0.01"
                    value={zoom}
                    onChange={(e) => setZoom(parseFloat(e.target.value))}
                    className="flex-1 h-1 bg-muted rounded-lg appearance-none cursor-pointer accent-primary"
                  />
                  <span className="text-muted-foreground text-xs font-bold">A+</span>
                </div>
              </div>

              {uploadError && (
                <div className="w-full text-xs text-destructive bg-destructive/10 border border-destructive/20 rounded-[calc(var(--radius)-10px)] p-3">
                  {uploadError}
                </div>
              )}

              <div className="flex w-full gap-3 justify-end mt-2">
                <Button
                  variant="outline"
                  onClick={handleCloseCropper}
                  disabled={isUploading}
                  className="flex-1 sm:flex-none"
                >
                  Cancel
                </Button>
                <Button
                  onClick={handleCropAndUpload}
                  disabled={isUploading}
                  className="flex-1 sm:flex-none min-w-[120px] bg-gradient-to-r from-primary to-violet-600 hover:from-primary/90 hover:to-violet-600/90 text-white font-medium"
                >
                  {isUploading ? (
                    <>
                      <Loader2 className="h-4 w-4 animate-spin mr-2" />
                      Saving...
                    </>
                  ) : (
                    "Save & Upload"
                  )}
                </Button>
              </div>
            </CardContent>
          </Card>
        </div>
      )}
    </main>
  );
}
