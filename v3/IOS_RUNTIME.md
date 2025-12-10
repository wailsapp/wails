# iOS Runtime Feature Plan

This document outlines proposed iOS-only runtime features for Wails v3, the initial milestones, and method shapes exposed to the frontend runtime as `IOS.*`.

## Goals
- Provide a first-class iOS runtime namespace: `IOS`.
- Expose UX-critical features with a small, well-defined, promise-based API.
- Follow the existing runtime pattern: JS -> /wails/runtime -> Go -> ObjC.

## Object: IOS
- Object ID: 11 (reserved in runtime objectNames)

## Milestone 1 (MVP)
- Haptics
  - `IOS.Haptics.Impact(style: "light"|"medium"|"heavy"|"soft"|"rigid"): Promise<void>`
- Device
  - `IOS.Device.Info(): Promise<{ model: string; systemName: string; systemVersion: string; isSimulator: boolean }>`

## Milestone 2
- Permissions
  - `IOS.Permissions.Request("camera"|"microphone"|"photos"|"notifications"): Promise<"granted"|"denied"|"limited">`
  - `IOS.Permissions.Status(kind): Promise<"granted"|"denied"|"limited"|"restricted"|"not_determined">`
- Camera
  - `IOS.Camera.PickPhoto(options?): Promise<{ uri: string }>`
  - `IOS.Camera.PickVideo(options?): Promise<{ uri: string, duration?: number }>`
- Photos
  - `IOS.Photos.SaveImage(dataURL|blob, options?): Promise<void>`
  - `IOS.Photos.SaveVideo(fileURI, options?): Promise<void>`

## Milestone 3
- Share
  - `IOS.Share.Sheet({ text?, url?, imageDataURL? }): Promise<void>`
- Files
  - `IOS.Files.Pick({ types?, multiple? }): Promise<Array<{ uri: string, name: string, size?: number }>>`
- Biometric
  - `IOS.Biometric.CanAuthenticate(): Promise<boolean>`
  - `IOS.Biometric.Authenticate(reason: string): Promise<boolean>`
- Notifications
  - `IOS.Notifications.RequestPermission(): Promise<boolean>`
  - `IOS.Notifications.Schedule(localNotification): Promise<string /* id */>`

## Notes
- All APIs should be safe no-ops on other platforms (reject with a meaningful error) or be tree-shaken by frontend bundlers.
- UI-affecting APIs must ensure main-thread execution in ObjC.
- File/Photo APIs will use security-scoped bookmarks where relevant.

## Implementation Status
- [x] Define plan (this document)
- [ ] JS runtime: add IOS object ID + IOS module exports
- [ ] Go: message dispatcher for IOS object
- [ ] iOS: Haptics.Impact(style) native bridge
- [ ] JS->Go->ObjC wiring for Haptics
- [ ] Device.Info() basic implementation
