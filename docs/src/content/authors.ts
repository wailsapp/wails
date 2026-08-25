import type { StarlightBlogUserConfig } from "starlight-blog";

// starlight-blog 0.28 accepts either a single blog config or an array of them,
// so the exported user-config type is a union and cannot be indexed directly.
// Pick the single-blog member, which is the one carrying the authors map.
type SingleBlogConfig = Extract<
  NonNullable<StarlightBlogUserConfig>,
  { authors?: unknown }
>;
type Authors = NonNullable<SingleBlogConfig["authors"]>;
export const authors: Authors = {
  leaanthony: {
    name: "Lea Anthony",
    title: "Maintainer of Wails",
    url: "https://github.com/leaanthony",
    picture: "https://github.com/leaanthony.png",
  },
};
